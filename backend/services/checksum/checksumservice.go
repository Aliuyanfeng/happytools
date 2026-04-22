package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// FileChecksum 单个文件的校验和结果
type FileChecksum struct {
	Path   string `json:"path"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
	CRC32  string `json:"crc32"`
	Error  string `json:"error"`
}

// ChecksumService 文件校验和计算服务
type ChecksumService struct {
	app *application.App
}

// NewChecksumService 创建服务实例
func NewChecksumService(app *application.App) *ChecksumService {
	return &ChecksumService{app: app}
}

// Calculate 对每个文件路径计算 MD5/SHA1/SHA256/CRC32
func (s *ChecksumService) Calculate(paths []string) []FileChecksum {
	results := make([]FileChecksum, 0, len(paths))
	for _, p := range paths {
		results = append(results, calcOne(p))
	}
	return results
}

// SelectFiles 打开多选文件对话框，无文件类型过滤
func (s *ChecksumService) SelectFiles() ([]string, error) {
	return s.app.Dialog.OpenFile().
		SetTitle("选择文件").
		PromptForMultipleSelection()
}

func calcOne(path string) FileChecksum {
	fc := FileChecksum{
		Path: path,
		Name: filepath.Base(path),
	}

	f, err := os.Open(path)
	if err != nil {
		fc.Error = err.Error()
		return fc
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		fc.Error = err.Error()
		return fc
	}
	fc.Size = info.Size()

	hMD5 := md5.New()
	hSHA1 := sha1.New()
	hSHA256 := sha256.New()
	hCRC32 := crc32.NewIEEE()

	w := io.MultiWriter(hMD5, hSHA1, hSHA256, hCRC32)
	if _, err = io.Copy(w, f); err != nil {
		fc.Error = err.Error()
		return fc
	}

	fc.MD5 = hex.EncodeToString(hMD5.Sum(nil))
	fc.SHA1 = hex.EncodeToString(hSHA1.Sum(nil))
	fc.SHA256 = hex.EncodeToString(hSHA256.Sum(nil))
	fc.CRC32 = hex.EncodeToString(hCRC32.Sum(nil))
	return fc
}
