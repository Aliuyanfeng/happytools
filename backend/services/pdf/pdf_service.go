/*
 * @Author: LiuYanFeng
 * @Date: 2026-07-21
 * @Description: PDF 解析服务 - 元数据读取与修改
 */
package pdf

import (
	"fmt"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// PDFMetadata PDF 元数据结构
type PDFMetadata struct {
	Title            string            `json:"title"`
	Author           string            `json:"author"`
	Subject          string            `json:"subject"`
	Keywords         []string          `json:"keywords"`
	Creator          string            `json:"creator"`
	Producer         string            `json:"producer"`
	CreationDate     string            `json:"creationDate"`
	ModificationDate string            `json:"modificationDate"`
	PageCount        int               `json:"pageCount"`
	Version          string            `json:"version"`
	Properties       map[string]string `json:"properties"`
	Tagged           bool              `json:"tagged"`
	Encrypted        bool              `json:"encrypted"`
}

// PDFService PDF 解析服务
type PDFService struct {
	app *application.App
}

// NewPDFService 创建服务实例
func NewPDFService(app *application.App) *PDFService {
	return &PDFService{app: app}
}

// OpenFileDialog 打开文件选择对话框，返回选中的 PDF 文件路径
func (s *PDFService) OpenFileDialog() (string, error) {
	result, err := s.app.Dialog.OpenFile().
		SetTitle("选择 PDF 文件").
		AddFilter("PDF 文件", "*.pdf").
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return result, nil
}

// SaveFileDialog 打开保存文件对话框，返回输出路径
func (s *PDFService) SaveFileDialog(defaultName string) (string, error) {
	result, err := s.app.Dialog.SaveFile().SetMessage("保存修改后的 PDF 文件").
		AddFilter("PDF 文件", "*.pdf").
		SetFilename(defaultName).
		PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetMetadata 读取 PDF 文件的元数据信息
func (s *PDFService) GetMetadata(filePath string) (*PDFMetadata, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	info, err := api.PDFInfo(f, filePath, nil, false, conf)
	if err != nil {
		return nil, fmt.Errorf("读取 PDF 信息失败: %v", err)
	}

	metadata := &PDFMetadata{
		Title:            info.Title,
		Author:           info.Author,
		Subject:          info.Subject,
		Keywords:         info.Keywords,
		Creator:          info.Creator,
		Producer:         info.Producer,
		CreationDate:     info.CreationDate,
		ModificationDate: info.ModificationDate,
		PageCount:        info.PageCount,
		Version:          info.Version,
		Properties:       info.Properties,
		Tagged:           info.Tagged,
		Encrypted:        info.Encrypted,
	}

	return metadata, nil
}

// SetStandardMetadata 修改 PDF 标准元数据（Title, Author, Subject）
// 由于 pdfcpu 没有提供直接修改标准元数据的公共 API，
// 此方法通过读取 PDF 上下文，直接操作 infoDict 来实现
func (s *PDFService) SetStandardMetadata(inFile, outFile, title, author, subject string) error {
	f, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	ctx, err := api.ReadValidateAndOptimize(f, conf)
	if err != nil {
		return fmt.Errorf("读取 PDF 失败: %v", err)
	}

	// 获取或创建 infoDict
	if err := ensureInfoDictForWrite(ctx); err != nil {
		return fmt.Errorf("确保 infoDict 失败: %v", err)
	}

	// 解引用 infoDict
	d, err := ctx.DereferenceDict(*ctx.Info)
	if err != nil || d == nil {
		return fmt.Errorf("读取 infoDict 失败: %v", err)
	}

	// 更新标准元数据字段
	if title != "" {
		d.Update("Title", types.StringLiteral(title))
	} else {
		d.Delete("Title")
	}

	if author != "" {
		d.Update("Author", types.StringLiteral(author))
	} else {
		d.Delete("Author")
	}

	if subject != "" {
		d.Update("Subject", types.StringLiteral(subject))
	} else {
		d.Delete("Subject")
	}

	// 写入输出文件
	tmpFile := inFile + ".tmp"
	if outFile != "" && inFile != outFile {
		tmpFile = outFile
	}

	f2, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer func() {
		f2.Close()
		if outFile == "" || inFile == outFile {
			os.Rename(tmpFile, inFile)
		}
	}()

	if err := api.Write(ctx, f2, conf); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("写入 PDF 失败: %v", err)
	}

	return nil
}

// AddKeywords 添加关键词到 PDF 文件
func (s *PDFService) AddKeywords(inFile, outFile string, keywords []string) error {
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	return api.AddKeywordsFile(inFile, outFile, keywords, conf)
}

// RemoveKeywords 从 PDF 文件中移除关键词
func (s *PDFService) RemoveKeywords(inFile, outFile string, keywords []string) error {
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	return api.RemoveKeywordsFile(inFile, outFile, keywords, conf)
}

// AddProperties 添加自定义属性到 PDF 文件
func (s *PDFService) AddProperties(inFile, outFile string, properties map[string]string) error {
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	return api.AddPropertiesFile(inFile, outFile, properties, conf)
}

// RemoveProperties 从 PDF 文件中移除自定义属性
func (s *PDFService) RemoveProperties(inFile, outFile string, properties []string) error {
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed
	return api.RemovePropertiesFile(inFile, outFile, properties, conf)
}

// ensureInfoDictForWrite 确保上下文中存在 infoDict，如果不存在则创建
func ensureInfoDictForWrite(ctx *model.Context) error {
	if ctx.Info != nil {
		return nil
	}

	// 创建新的 infoDict
	d := types.NewDict()
	d.InsertString("Producer", "pdfcpu "+model.VersionStr)

	ir, err := ctx.IndRefForNewObject(d)
	if err != nil {
		return err
	}

	ctx.Info = ir
	return nil
}

// GetPageInfo 获取 PDF 页面详细信息（尺寸等）
func (s *PDFService) GetPageInfo(filePath string) (*PDFMetadata, error) {
	return s.GetMetadata(filePath)
}

// ReadPDFContent 读取 PDF 文件用于流式处理（内部辅助方法）
func readPDFForWrite(inFile string) (*os.File, *model.Context, *model.Configuration, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("打开文件失败: %v", err)
	}

	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	ctx, err := api.ReadValidateAndOptimize(f, conf)
	if err != nil {
		f.Close()
		return nil, nil, nil, fmt.Errorf("读取 PDF 失败: %v", err)
	}

	return f, ctx, conf, nil
}

// writePDFContext 将修改后的 PDF 上下文写入文件
func writePDFContext(ctx *model.Context, conf *model.Configuration, inFile, outFile string) error {
	tmpFile := inFile + ".tmp"
	if outFile != "" && inFile != outFile {
		tmpFile = outFile
	}

	f2, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}

	ok := false
	defer func() {
		f2.Close()
		if !ok {
			os.Remove(tmpFile)
		}
		if ok && (outFile == "" || inFile == outFile) {
			os.Rename(tmpFile, inFile)
		}
	}()

	if err := api.Write(ctx, f2, conf); err != nil {
		return fmt.Errorf("写入 PDF 失败: %v", err)
	}

	ok = true
	return nil
}

