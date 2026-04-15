// Package ncm
// @description: NCM 转 MP3/FLAC 服务
// NCM 文件格式解析参考: https://github.com/FLCYR/ncmToMp3
package ncm

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AES 密钥
var coreKey = []byte{0x68, 0x7A, 0x48, 0x52, 0x41, 0x6D, 0x73, 0x6F, 0x35, 0x6B, 0x49, 0x6E, 0x62, 0x61, 0x78, 0x57}
var metaKey = []byte{0x23, 0x31, 0x34, 0x6C, 0x6A, 0x6B, 0x5F, 0x21, 0x5C, 0x5D, 0x26, 0x30, 0x55, 0x3C, 0x27, 0x28}

// NCM 文件魔数
var ncmMagic = []byte{0x43, 0x54, 0x45, 0x4E, 0x46, 0x44, 0x41, 0x4D}

// ConvertResult 单个文件转换结果
type ConvertResult struct {
	Input      string `json:"input"`
	Output     string `json:"output"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	LrcCopied  bool   `json:"lrc_copied"`
}

// CheckLrc 检查 NCM 文件同目录下是否存在同名 .lrc 文件，返回 lrc 路径（不存在则返回空）
func (s *NCMService) CheckLrc(ncmPath string) string {
	return findLrc(ncmPath)
}

// CheckLrcBatch 批量检查，返回 map[ncmPath]lrcPath
func (s *NCMService) CheckLrcBatch(ncmPaths []string) map[string]string {
	result := make(map[string]string, len(ncmPaths))
	for _, p := range ncmPaths {
		result[p] = findLrc(p)
	}
	return result
}

// findLrc 查找同名 lrc 文件路径，不存在返回空字符串
func findLrc(ncmPath string) string {
	base := strings.TrimSuffix(ncmPath, filepath.Ext(ncmPath))
	lrc := base + ".lrc"
	if _, err := os.Stat(lrc); err == nil {
		return lrc
	}
	return ""
}

// NCMService NCM 转换服务
type NCMService struct {
	app *application.App
}

func NewNCMService(app *application.App) *NCMService {
	return &NCMService{app: app}
}

// SelectOutputDir 选择输出目录
func (s *NCMService) SelectOutputDir() (string, error) {
	dir, err := s.app.Dialog.OpenFile().
		SetTitle("选择输出目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// SelectFiles 选择 NCM 文件
func (s *NCMService) SelectFiles() ([]string, error) {
	files, err := s.app.Dialog.OpenFile().
		SetTitle("选择 NCM 文件").
		AddFilter("NCM 文件", "*.ncm").
		CanChooseFiles(true).
		PromptForMultipleSelection()
	if err != nil {
		return nil, err
	}
	return files, nil
}

// SelectDir 选择目录并返回其中所有 NCM 文件路径
func (s *NCMService) SelectDir() ([]string, error) {
	dir, err := s.app.Dialog.OpenFile().
		SetTitle("选择包含 NCM 文件的目录").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()
	if err != nil || dir == "" {
		return nil, err
	}
	return scanNCMFiles(dir)
}

// ScanDir 扫描指定目录下所有 NCM 文件（供前端拖拽目录时调用）
func (s *NCMService) ScanDir(dir string) ([]string, error) {
	return scanNCMFiles(dir)
}

// scanNCMFiles 递归扫描目录下所有 .ncm 文件
func scanNCMFiles(dir string) ([]string, error) {
	var result []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无权限目录
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".ncm" {
			result = append(result, path)
		}
		return nil
	})
	return result, err
}

// OpenOutputDir 在系统文件管理器中打开指定目录
func (s *NCMService) OpenOutputDir(dir string) error {
	if dir == "" {
		return errors.New("目录路径为空")
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", filepath.FromSlash(dir))
	case "darwin":
		cmd = exec.Command("open", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	return cmd.Start()
}

// ConvertOne 转换单个 NCM 文件
func (s *NCMService) ConvertOne(file string, outputDir string) (*ConvertResult, error) {
	outDir := outputDir
	if outDir == "" {
		outDir = filepath.Dir(file)
	}
	outPath, err := convertNCM(file, outDir)
	if err != nil {
		return &ConvertResult{Input: file, Success: false, Error: err.Error()}, nil
	}
	r := &ConvertResult{Input: file, Output: outPath, Success: true}
	if lrcSrc := findLrc(file); lrcSrc != "" {
		baseName := strings.TrimSuffix(filepath.Base(outPath), filepath.Ext(outPath))
		lrcDst := filepath.Join(outDir, baseName+".lrc")
		if copyErr := copyFile(lrcSrc, lrcDst); copyErr == nil {
			r.LrcCopied = true
		}
	}
	return r, nil
}
func (s *NCMService) ConvertFiles(files []string, outputDir string) ([]ConvertResult, error) {
	results := make([]ConvertResult, 0, len(files))
	for _, f := range files {
		outDir := outputDir
		if outDir == "" {
			outDir = filepath.Dir(f)
		}
		outPath, err := convertNCM(f, outDir)
		if err != nil {
			results = append(results, ConvertResult{
				Input:   f,
				Success: false,
				Error:   err.Error(),
			})
			continue
		}

		r := ConvertResult{Input: f, Output: outPath, Success: true}

		// 复制同名 lrc 文件
		if lrcSrc := findLrc(f); lrcSrc != "" {
			baseName := strings.TrimSuffix(filepath.Base(outPath), filepath.Ext(outPath))
			lrcDst := filepath.Join(outDir, baseName+".lrc")
			if copyErr := copyFile(lrcSrc, lrcDst); copyErr == nil {
				r.LrcCopied = true
			}
		}

		results = append(results, r)
	}
	return results, nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// ─── 核心转换逻辑 ───────────────────────────────────────────────

func convertNCM(inputPath, outputDir string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	offset := 0

	// 1. 校验魔数 (8 bytes)
	if len(data) < 10 || string(data[:8]) != string(ncmMagic) {
		return "", errors.New("不是有效的 NCM 文件")
	}
	offset += 10 // magic(8) + gap(2)

	// 2. 解密 RC4 密钥
	if offset+4 > len(data) {
		return "", errors.New("文件格式错误: key length")
	}
	keyLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	if offset+keyLen > len(data) {
		return "", errors.New("文件格式错误: key data")
	}
	keyData := make([]byte, keyLen)
	copy(keyData, data[offset:offset+keyLen])
	offset += keyLen

	// 按字节 XOR 0x64
	for i := range keyData {
		keyData[i] ^= 0x64
	}

	// AES-128-ECB 解密
	rc4Key, err := aesECBDecrypt(keyData, coreKey)
	if err != nil {
		return "", fmt.Errorf("解密 RC4 密钥失败: %w", err)
	}
	// 去除 "neteasecloudmusic" 前缀 (17 bytes)
	if len(rc4Key) < 17 {
		return "", errors.New("RC4 密钥格式错误")
	}
	rc4Key = rc4Key[17:]

	// 3. 解析音乐信息
	if offset+4 > len(data) {
		return "", errors.New("文件格式错误: meta length")
	}
	metaLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	format := "mp3" // 默认
	if metaLen > 0 && offset+metaLen <= len(data) {
		metaData := make([]byte, metaLen)
		copy(metaData, data[offset:offset+metaLen])

		// XOR 0x63
		for i := range metaData {
			metaData[i] ^= 0x63
		}

		// 去除前 22 字节 ("163 key(Don't modify):")
		if len(metaData) > 22 {
			metaData = metaData[22:]
			// Base64 解码
			decoded, err := base64.StdEncoding.DecodeString(string(metaData))
			if err == nil {
				// AES 解密
				decrypted, err := aesECBDecrypt(decoded, metaKey)
				if err == nil && len(decrypted) > 6 {
					// 去除前 6 字节 ("music:")
					jsonData := decrypted[6:]
					var info struct {
						Format string `json:"format"`
					}
					if json.Unmarshal(jsonData, &info) == nil && info.Format != "" {
						format = strings.ToLower(info.Format)
					}
				}
			}
		}
	}
	offset += metaLen

	// 4. 跳过 CRC(4) + Gap(5)，提取封面图片
	if offset+9 > len(data) {
		return "", errors.New("文件格式错误: crc/gap")
	}
	offset += 9 // CRC(4) + gap(5)

	if offset+4 > len(data) {
		return "", errors.New("文件格式错误: image size")
	}
	imageSize := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4

	var imageData []byte
	if imageSize > 0 && offset+imageSize <= len(data) {
		imageData = make([]byte, imageSize)
		copy(imageData, data[offset:offset+imageSize])
	}
	offset += imageSize

	// 5. 解密音乐数据
	if offset >= len(data) {
		return "", errors.New("文件格式错误: 无音乐数据")
	}
	musicData := make([]byte, len(data)-offset)
	copy(musicData, data[offset:])

	// RC4-KSA 生成 S 盒
	sBox := rc4KSA(rc4Key)

	// 自定义解密（非标准 RC4-PRGA）
	rc4Decrypt(sBox, musicData)

	// 6. 嵌入封面并写出文件
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outPath := filepath.Join(outputDir, baseName+"."+format)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	if len(imageData) > 0 {
		switch format {
		case "mp3":
			musicData = embedCoverMP3(musicData, imageData)
		case "flac":
			musicData = embedCoverFLAC(musicData, imageData)
		}
	}

	if err := os.WriteFile(outPath, musicData, 0644); err != nil {
		return "", fmt.Errorf("写出文件失败: %w", err)
	}
	return outPath, nil
}

// ─── AES-128-ECB 解密 ───────────────────────────────────────────

func aesECBDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("数据长度不是块大小的整数倍")
	}
	out := make([]byte, len(data))
	for i := 0; i < len(data); i += bs {
		block.Decrypt(out[i:i+bs], data[i:i+bs])
	}
	// 去除 PKCS7 填充
	return pkcs7Unpad(out)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("空数据")
	}
	pad := int(data[len(data)-1])
	if pad == 0 || pad > aes.BlockSize {
		return data, nil // 无填充或非标准，直接返回
	}
	return data[:len(data)-pad], nil
}

// ─── RC4 ────────────────────────────────────────────────────────

func rc4KSA(key []byte) [256]byte {
	var s [256]byte
	for i := range s {
		s[i] = byte(i)
	}
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) % 256
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// NCM 使用的自定义解密（非标准 PRGA）
func rc4Decrypt(s [256]byte, data []byte) {
	for idx := range data {
		i := (idx + 1) % 256
		j := (i + int(s[i])) % 256
		k := (int(s[i]) + int(s[j])) % 256
		data[idx] ^= s[k]
	}
}

// ─── 封面嵌入 ────────────────────────────────────────────────────

// embedCoverMP3 将封面图片嵌入 MP3 的 ID3v2 tag（APIC 帧）
// 策略：若已有 ID3v2 tag，在最后一个有效帧之后插入 APIC 帧；否则在文件头部新建 ID3v2 tag
func embedCoverMP3(mp3Data, imageData []byte) []byte {
	mimeType := detectImageMIME(imageData)

	// 构造 APIC 帧内容: encoding(1) + mime + 0x00 + picType(1) + desc + 0x00 + imageData
	var apicContent []byte
	apicContent = append(apicContent, 0x00)             // encoding: ISO-8859-1
	apicContent = append(apicContent, []byte(mimeType)...)
	apicContent = append(apicContent, 0x00)             // mime 结束符
	apicContent = append(apicContent, 0x03)             // picture type: Cover (front)
	apicContent = append(apicContent, 0x00)             // description: empty
	apicContent = append(apicContent, imageData...)

	apicFrame := buildID3Frame("APIC", apicContent)

	// 检查是否已有 ID3v2 tag
	if len(mp3Data) >= 10 && string(mp3Data[:3]) == "ID3" {
		tagSize := id3SyncsafeToInt(mp3Data[6:10])
		headerSize := 10
		tagEnd := headerSize + tagSize

		// 找到最后一个有效帧的结束位置（跳过 padding）
		lastFrameEnd := headerSize
		offset := headerSize
		for offset+10 <= tagEnd {
			fsize := int(mp3Data[offset+4])<<24 | int(mp3Data[offset+5])<<16 |
				int(mp3Data[offset+6])<<8 | int(mp3Data[offset+7])
			fid := string(mp3Data[offset : offset+4])

			// 遇到空帧或全零说明进入 padding 区域
			if fsize == 0 || mp3Data[offset] == 0x00 {
				break
			}
			// 已有 APIC 帧，直接返回
			if fid == "APIC" {
				return mp3Data
			}
			lastFrameEnd = offset + 10 + fsize
			offset = lastFrameEnd
		}

		// 在最后一个有效帧之后插入 APIC，更新 tag size
		newTagSize := tagSize + len(apicFrame)
		var out []byte
		out = append(out, mp3Data[:3]...)                        // "ID3"
		out = append(out, mp3Data[3:6]...)                       // version + flags
		out = append(out, intToID3Syncsafe(newTagSize)...)       // 新 tag 大小
		out = append(out, mp3Data[headerSize:lastFrameEnd]...)   // 原有有效帧
		out = append(out, apicFrame...)                          // 新 APIC 帧
		// 保留剩余 padding + 音频数据
		out = append(out, mp3Data[lastFrameEnd:]...)
		return out
	}

	// 没有 ID3v2 tag，新建一个
	newTagSize := len(apicFrame)
	var tag []byte
	tag = append(tag, []byte("ID3")...)
	tag = append(tag, 0x03, 0x00) // version 2.3.0
	tag = append(tag, 0x00)       // flags
	tag = append(tag, intToID3Syncsafe(newTagSize)...)
	tag = append(tag, apicFrame...)

	return append(tag, mp3Data...)
}

// embedCoverFLAC 将封面图片嵌入 FLAC 的 PICTURE metadata block
func embedCoverFLAC(flacData, imageData []byte) []byte {
	if len(flacData) < 4 || string(flacData[:4]) != "fLaC" {
		return flacData
	}

	mimeType := detectImageMIME(imageData)

	// 构造 PICTURE block 数据
	var pic []byte
	pic = appendUint32BE(pic, 3)                          // picture type: Cover (front)
	pic = appendUint32BE(pic, uint32(len(mimeType)))
	pic = append(pic, []byte(mimeType)...)
	pic = appendUint32BE(pic, 0)                          // description length: 0
	pic = appendUint32BE(pic, 0)                          // width (unknown)
	pic = appendUint32BE(pic, 0)                          // height (unknown)
	pic = appendUint32BE(pic, 0)                          // color depth (unknown)
	pic = appendUint32BE(pic, 0)                          // color count (unknown)
	pic = appendUint32BE(pic, uint32(len(imageData)))
	pic = append(pic, imageData...)

	// 找到第一个 STREAMINFO 之后的位置，在 STREAMINFO 后插入 PICTURE block
	// FLAC 结构: "fLaC" + metadata blocks + audio frames
	// 每个 metadata block: 1字节(last-flag|type) + 3字节(length) + data
	offset := 4
	insertAt := -1
	for offset+4 <= len(flacData) {
		blockHeader := flacData[offset]
		isLast := (blockHeader & 0x80) != 0
		blockType := blockHeader & 0x7F
		blockLen := int(flacData[offset+1])<<16 | int(flacData[offset+2])<<8 | int(flacData[offset+3])

		// 如果已有 PICTURE block (type=6)，跳过嵌入
		if blockType == 6 {
			return flacData
		}

		if isLast {
			insertAt = offset
			break
		}
		offset += 4 + blockLen
	}

	if insertAt < 0 {
		return flacData
	}

	// 将原来的 last block 的 last-flag 清除，插入 PICTURE block（设为 last）
	var out []byte
	out = append(out, flacData[:insertAt]...)

	// 修改原 last block 的 last-flag → 不再是 last
	origHeader := flacData[insertAt] & 0x7F // 清除 last-flag
	out = append(out, origHeader)
	origBlockLen := int(flacData[insertAt+1])<<16 | int(flacData[insertAt+2])<<8 | int(flacData[insertAt+3])
	out = append(out, flacData[insertAt+1:insertAt+4+origBlockLen]...)

	// 追加 PICTURE block（设为 last）
	out = append(out, 0x80|0x06) // last=1, type=6 (PICTURE)
	picLen := len(pic)
	out = append(out, byte(picLen>>16), byte(picLen>>8), byte(picLen))
	out = append(out, pic...)

	// 追加剩余音频数据
	out = append(out, flacData[insertAt+4+origBlockLen:]...)
	return out
}

// ─── 工具函数 ────────────────────────────────────────────────────

func detectImageMIME(data []byte) string {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}
	return "image/jpeg" // 默认
}

func buildID3Frame(id string, content []byte) []byte {
	var frame []byte
	frame = append(frame, []byte(id)...)
	// 帧大小（4字节大端，ID3v2.3 不用 syncsafe）
	size := len(content)
	frame = append(frame, byte(size>>24), byte(size>>16), byte(size>>8), byte(size))
	frame = append(frame, 0x00, 0x00) // flags
	frame = append(frame, content...)
	return frame
}

func hasID3Frame(tagData []byte, frameID string) bool {
	offset := 0
	for offset+10 <= len(tagData) {
		id := string(tagData[offset : offset+4])
		size := int(tagData[offset+4])<<24 | int(tagData[offset+5])<<16 | int(tagData[offset+6])<<8 | int(tagData[offset+7])
		if id == frameID {
			return true
		}
		if size == 0 {
			break
		}
		offset += 10 + size
	}
	return false
}

func id3SyncsafeToInt(b []byte) int {
	return int(b[0])<<21 | int(b[1])<<14 | int(b[2])<<7 | int(b[3])
}

func intToID3Syncsafe(n int) []byte {
	return []byte{
		byte((n >> 21) & 0x7F),
		byte((n >> 14) & 0x7F),
		byte((n >> 7) & 0x7F),
		byte(n & 0x7F),
	}
}

func appendUint32BE(b []byte, v uint32) []byte {
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
