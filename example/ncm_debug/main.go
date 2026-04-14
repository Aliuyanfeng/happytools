package main

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var coreKey = []byte{0x68, 0x7A, 0x48, 0x52, 0x41, 0x6D, 0x73, 0x6F, 0x35, 0x6B, 0x49, 0x6E, 0x62, 0x61, 0x78, 0x57}
var metaKey = []byte{0x23, 0x31, 0x34, 0x6C, 0x6A, 0x6B, 0x5F, 0x21, 0x5C, 0x5D, 0x26, 0x30, 0x55, 0x3C, 0x27, 0x28}
var ncmMagic = []byte{0x43, 0x54, 0x45, 0x4E, 0x46, 0x44, 0x41, 0x4D}

func main() {
	input := `D:\CloudMusic\VipSongsDownload\亚森 - 荒漠上行走.ncm`
	outDir := `D:\CloudMusic\VipSongsDownload\`

	outPath, err := convertNCM(input, outDir)
	if err != nil {
		fmt.Println("失败:", err)
		return
	}
	fmt.Println("✅ 成功:", outPath)

	// 验证 ID3 tag 里是否有 APIC 帧
	data, _ := os.ReadFile(outPath)
	if len(data) >= 10 && string(data[:3]) == "ID3" {
		tagSize := id3SyncsafeToInt(data[6:10])
		tagData := data[10 : 10+tagSize]
		if hasID3Frame(tagData, "APIC") {
			fmt.Println("✅ 封面已嵌入 (APIC 帧存在)")
		} else {
			fmt.Println("❌ 未找到 APIC 帧")
		}
	}
}

func convertNCM(inputPath, outputDir string) (string, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	offset := 0
	if len(data) < 10 || string(data[:8]) != string(ncmMagic) {
		return "", errors.New("不是有效的 NCM 文件")
	}
	offset += 10

	keyLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	keyData := make([]byte, keyLen)
	copy(keyData, data[offset:offset+keyLen])
	offset += keyLen
	for i := range keyData { keyData[i] ^= 0x64 }
	rc4Key, err := aesECBDecrypt(keyData, coreKey)
	if err != nil { return "", err }
	rc4Key = rc4Key[17:]

	metaLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	format := "mp3"
	if metaLen > 0 {
		metaData := make([]byte, metaLen)
		copy(metaData, data[offset:offset+metaLen])
		for i := range metaData { metaData[i] ^= 0x63 }
		if len(metaData) > 22 {
			decoded, err := base64.StdEncoding.DecodeString(string(metaData[22:]))
			if err == nil {
				decrypted, err := aesECBDecrypt(decoded, metaKey)
				if err == nil && len(decrypted) > 6 {
					var info struct{ Format string `json:"format"` }
					if json.Unmarshal(decrypted[6:], &info) == nil && info.Format != "" {
						format = strings.ToLower(info.Format)
					}
				}
			}
		}
	}
	offset += metaLen
	offset += 9 // CRC(4)+gap(5)

	imageSize := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	var imageData []byte
	if imageSize > 0 {
		imageData = make([]byte, imageSize)
		copy(imageData, data[offset:offset+imageSize])
		fmt.Printf("封面大小: %d bytes, MIME: %s\n", imageSize, detectImageMIME(imageData))
	}
	offset += imageSize

	musicData := make([]byte, len(data)-offset)
	copy(musicData, data[offset:])
	sBox := rc4KSA(rc4Key)
	rc4Decrypt(sBox, musicData)

	if len(imageData) > 0 && format == "mp3" {
		musicData = embedCoverMP3(musicData, imageData)
		fmt.Printf("嵌入封面后 MP3 大小: %d bytes\n", len(musicData))
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outPath := filepath.Join(outputDir, baseName+"."+format)
	os.MkdirAll(outputDir, 0755)
	os.WriteFile(outPath, musicData, 0644)
	return outPath, nil
}

func embedCoverMP3(mp3Data, imageData []byte) []byte {
	mimeType := detectImageMIME(imageData)
	var apicContent []byte
	apicContent = append(apicContent, 0x00)
	apicContent = append(apicContent, []byte(mimeType)...)
	apicContent = append(apicContent, 0x00)
	apicContent = append(apicContent, 0x03)
	apicContent = append(apicContent, 0x00)
	apicContent = append(apicContent, imageData...)
	apicFrame := buildID3Frame("APIC", apicContent)

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
			if fsize == 0 || mp3Data[offset] == 0x00 {
				break
			}
			if fid == "APIC" {
				return mp3Data
			}
			lastFrameEnd = offset + 10 + fsize
			offset = lastFrameEnd
		}
		fmt.Printf("lastFrameEnd=%d, tagEnd=%d, apicFrame len=%d\n", lastFrameEnd, tagEnd, len(apicFrame))

		newTagSize := tagSize + len(apicFrame)
		var out []byte
		out = append(out, mp3Data[:3]...)
		out = append(out, mp3Data[3:6]...)
		out = append(out, intToID3Syncsafe(newTagSize)...)
		out = append(out, mp3Data[headerSize:lastFrameEnd]...)
		out = append(out, apicFrame...)
		out = append(out, mp3Data[lastFrameEnd:]...)
		return out
	}
	newTagSize := len(apicFrame)
	var tag []byte
	tag = append(tag, []byte("ID3")...)
	tag = append(tag, 0x03, 0x00)
	tag = append(tag, 0x00)
	tag = append(tag, intToID3Syncsafe(newTagSize)...)
	tag = append(tag, apicFrame...)
	return append(tag, mp3Data...)
}

func detectImageMIME(data []byte) string {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 { return "image/jpeg" }
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" { return "image/png" }
	return "image/jpeg"
}

func buildID3Frame(id string, content []byte) []byte {
	var frame []byte
	frame = append(frame, []byte(id)...)
	size := len(content)
	frame = append(frame, byte(size>>24), byte(size>>16), byte(size>>8), byte(size))
	frame = append(frame, 0x00, 0x00)
	frame = append(frame, content...)
	return frame
}

func hasID3Frame(tagData []byte, frameID string) bool {
	offset := 0
	for offset+10 <= len(tagData) {
		id := string(tagData[offset : offset+4])
		size := int(tagData[offset+4])<<24 | int(tagData[offset+5])<<16 | int(tagData[offset+6])<<8 | int(tagData[offset+7])
		if id == frameID { return true }
		if size == 0 { break }
		offset += 10 + size
	}
	return false
}

func id3SyncsafeToInt(b []byte) int {
	return int(b[0])<<21 | int(b[1])<<14 | int(b[2])<<7 | int(b[3])
}

func intToID3Syncsafe(n int) []byte {
	return []byte{byte((n >> 21) & 0x7F), byte((n >> 14) & 0x7F), byte((n >> 7) & 0x7F), byte(n & 0x7F)}
}

func aesECBDecrypt(data, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	bs := block.BlockSize()
	if len(data)%bs != 0 { return nil, errors.New("长度错误") }
	out := make([]byte, len(data))
	for i := 0; i < len(data); i += bs { block.Decrypt(out[i:i+bs], data[i:i+bs]) }
	pad := int(out[len(out)-1])
	if pad > 0 && pad <= bs { return out[:len(out)-pad], nil }
	return out, nil
}

func rc4KSA(key []byte) [256]byte {
	var s [256]byte
	for i := range s { s[i] = byte(i) }
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) % 256
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func rc4Decrypt(s [256]byte, data []byte) {
	for idx := range data {
		i := (idx + 1) % 256
		j := (i + int(s[i])) % 256
		k := (int(s[i]) + int(s[j])) % 256
		data[idx] ^= s[k]
	}
}
