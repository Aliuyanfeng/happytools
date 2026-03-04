package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	apiKey       = "2fbdb6b98badd2ca351e76f371e593b4ef281ea58b2a4c63fdc7caecad28be10"
	uploadDir    = "./samples"
	baseURL      = "https://www.virustotal.com/api/v3"
	workerCount  = 5
	waitLimitSec = 100 // 每个文件最多等待分析100秒
)

type ScanResult struct {
	FileName   string
	SHA256     string
	Malicious  string
	Suspicious string
}

func main() {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Fatal("读取目录失败:", err)
	}

	// 任务通道和结果通道
	tasks := make(chan string, len(files))
	results := make(chan ScanResult, len(files))

	// 启动 worker
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// 投放任务
	for _, f := range files {
		if !f.IsDir() {
			tasks <- filepath.Join(uploadDir, f.Name())
		}
	}
	close(tasks)

	// 等待所有 worker 完成
	wg.Wait()
	close(results)

	// 写 CSV
	var allResults []ScanResult
	for r := range results {
		allResults = append(allResults, r)
	}
	if err := writeCSV(allResults); err != nil {
		log.Println("写入 CSV 失败:", err)
	} else {
		fmt.Println("扫描结果已保存到 vt_result.csv")
	}

	// 阻塞退出
	fmt.Println("所有任务完成，按 Ctrl+C 退出...")
	waitExit()
}

// worker 负责上传+分析
func worker(id int, tasks <-chan string, results chan<- ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range tasks {
		filename := filepath.Base(path)
		fmt.Printf("[Worker %d] 处理: %s\n", id, filename)

		analysisID, err := uploadFile(path)
		if err != nil {
			fmt.Printf("[Worker %d] 上传失败: %s\n", id, err)
			results <- ScanResult{FileName: filename}
			continue
		}

		sha256, timedOut := waitForAnalysis(analysisID, waitLimitSec)
		if timedOut {
			fmt.Printf("[Worker %d] 超时: %s\n", id, filename)
			results <- ScanResult{FileName: filename}
			continue
		}

		mal, susp, err := getStats(sha256)
		if err != nil {
			fmt.Printf("[Worker %d] 获取分析失败: %s\n", id, err)
			results <- ScanResult{FileName: filename, SHA256: sha256}
			continue
		}

		results <- ScanResult{
			FileName:   filename,
			SHA256:     sha256,
			Malicious:  mal,
			Suspicious: susp,
		}
	}
}

func uploadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", baseURL+"/files", &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Printf("结果集合: %+v\n", resp)
	var respData struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	return respData.Data.ID, nil
}

func waitForAnalysis(id string, timeoutSec int) (sha256 string, timedOut bool) {
	start := time.Now()
	for time.Since(start).Seconds() < float64(timeoutSec) {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest("GET", baseURL+"/analyses/"+id, nil)
		req.Header.Set("x-apikey", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		
		var result struct {
			Data struct {
				Attributes struct {
					Status string `json:"status"`
				} `json:"attributes"`
			} `json:"data"`
			Meta struct {
				FileInfo struct {
					SHA256 string `json:"sha256"`
				} `json:"file_info"`
			} `json:"meta"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			continue
		}

		if result.Data.Attributes.Status == "completed" {
			return result.Meta.FileInfo.SHA256, false
		}
	}
	return "", true
}

func getStats(sha256 string) (malicious string, suspicious string, err error) {
	req, _ := http.NewRequest("GET", baseURL+"/files/"+sha256, nil)
	req.Header.Set("x-apikey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	attr := result["data"].(map[string]interface{})["attributes"].(map[string]interface{})
	stats := attr["last_analysis_stats"].(map[string]interface{})
	mal := fmt.Sprintf("%v", stats["malicious"])
	susp := fmt.Sprintf("%v", stats["suspicious"])

	return mal, susp, nil
}

func writeCSV(results []ScanResult) error {
	file, err := os.Create("vt_result.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"FileName", "SHA256", "Malicious", "Suspicious"})

	// 写入数据
	for _, r := range results {
		writer.Write([]string{r.FileName, r.SHA256, r.Malicious, r.Suspicious})
	}
	return nil
}

func waitExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("程序退出")
}
