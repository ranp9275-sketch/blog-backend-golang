package handlers

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// extractPDFText 从PDF文件中提取文本内容
func extractPDFText(pdfPath string) (string, error) {
	// 打开PDF文件
	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var textBuilder strings.Builder
	totalPages := r.NumPage()

	// 遍历所有页面
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		p := r.Page(pageNum)
		if p.V.IsNull() {
			continue
		}

		// 提取页面文本
		text, err := p.GetPlainText(nil)
		if err != nil {
			// 如果某一页提取失败，继续处理下一页
			continue
		}

		// 添加页面文本
		if len(text) > 0 {
			textBuilder.WriteString(text)
			textBuilder.WriteString("\n\n")
		}
	}

	extractedText := strings.TrimSpace(textBuilder.String())

	// 如果提取的文本太少，可能提取失败
	if len(extractedText) < 50 {
		return "", fmt.Errorf("extracted text too short")
	}

	return extractedText, nil
}
