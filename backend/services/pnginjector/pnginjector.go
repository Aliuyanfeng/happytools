/*
 * @Author: LiuYanFeng
 * @Date: 2026-03-16
 * @Description: PNG Chunk 注入服务
 */
package pnginjector

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var pngSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

// PNGChunk 表示一个 PNG chunk 的信息
type PNGChunk struct {
	Index    int    `json:"index"`
	Type     string `json:"type"`
	Length   uint32 `json:"length"`
	DataHex  string `json:"dataHex"`
	CRC      string `json:"crc"`
	Offset   int64  `json:"offset"`
	Critical bool   `json:"critical"`
}

// PNGInjectorService PNG chunk 注入服务
type PNGInjectorService struct {
	app *application.App
}

// NewPNGInjectorService 创建服务实例
func NewPNGInjectorService(app *application.App) *PNGInjectorService {
	return &PNGInjectorService{app: app}
}


// OpenFileDialog 打开文件选择对话框，返回选中的 PNG 文件路径
func (s *PNGInjectorService) OpenFileDialog() (string, error) {
	result, err := s.app.Dialog.OpenFile().
		SetTitle("选择 PNG 文件").
		AddFilter("PNG 图片", "*.png").
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return result, nil
}

// SaveFileDialog 打开保存文件对话框，返回输出路径
func (s *PNGInjectorService) SaveFileDialog(defaultName string) (string, error) {
	result, err := s.app.Dialog.SaveFile().SetMessage("保存注入后的 PNG 文件").
		AddFilter("PNG 图片", "*.png").
		SetFilename(defaultName).
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return result, nil
}

// ParsePNG 解析 PNG 文件，返回所有 chunk 信息
func (s *PNGInjectorService) ParsePNG(filePath string) ([]PNGChunk, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], pngSignature) {
		return nil, fmt.Errorf("不是有效的 PNG 文件")
	}

	var chunks []PNGChunk
	offset := int64(8)
	index := 0

	for offset+12 <= int64(len(data)) {
		length := binary.BigEndian.Uint32(data[offset : offset+4])
		chunkType := string(data[offset+4 : offset+8])

		dataEnd := offset + 8 + int64(length)
		if dataEnd+4 > int64(len(data)) {
			break
		}

		crcVal := binary.BigEndian.Uint32(data[dataEnd : dataEnd+4])

		// 数据预览：最多取前 64 字节
		previewLen := int64(64)
		if int64(length) < previewLen {
			previewLen = int64(length)
		}
		dataHex := hex.EncodeToString(data[offset+8 : offset+8+previewLen])
		if int64(length) > previewLen {
			dataHex += "…"
		}

		critical := chunkType[0] >= 'A' && chunkType[0] <= 'Z'

		chunks = append(chunks, PNGChunk{
			Index:    index,
			Type:     chunkType,
			Length:   length,
			DataHex:  dataHex,
			CRC:      fmt.Sprintf("%08X", crcVal),
			Offset:   offset,
			Critical: critical,
		})

		offset = dataEnd + 4
		index++
	}

	return chunks, nil
}

// InjectChunk 在指定位置注入 chunk
// position: "before" 或 "after"，targetIndex: 目标 chunk 的索引
func (s *PNGInjectorService) InjectChunk(inputPath, outputPath, chunkType, payload, position string, targetIndex int) error {
	if len(chunkType) != 4 {
		return fmt.Errorf("chunk 类型必须是 4 个字符，当前为 %d 个字符", len(chunkType))
	}
	if payload == "" {
		return fmt.Errorf("payload 不能为空")
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}
	if len(data) < 8 || !bytes.Equal(data[:8], pngSignature) {
		return fmt.Errorf("不是有效的 PNG 文件")
	}

	chunks, err := s.ParsePNG(inputPath)
	if err != nil {
		return err
	}
	if targetIndex < 0 || targetIndex >= len(chunks) {
		return fmt.Errorf("无效的 chunk 索引: %d（共 %d 个 chunk）", targetIndex, len(chunks))
	}

	target := chunks[targetIndex]
	var injectAt int64
	if position == "before" {
		injectAt = target.Offset
	} else {
		injectAt = target.Offset + 12 + int64(target.Length)
	}

	newChunk := buildChunk([]byte(chunkType), []byte(payload))

	var out bytes.Buffer
	out.Write(data[:injectAt])
	out.Write(newChunk)
	out.Write(data[injectAt:])

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}
	return os.WriteFile(outputPath, out.Bytes(), 0644)
}

// buildChunk 按 PNG 规范构建 chunk: length(4) + type(4) + data + crc(4)
func buildChunk(chunkType, data []byte) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint32(len(data)))
	buf.Write(chunkType)
	buf.Write(data)
	crc := crc32.ChecksumIEEE(append(chunkType, data...))
	binary.Write(&buf, binary.BigEndian, crc)
	return buf.Bytes()
}
