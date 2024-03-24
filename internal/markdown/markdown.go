package markdown

import (
	"archive/zip"
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliamerj/docxer/internal/template"
	"github.com/aliamerj/docxer/internal/utils"
)

//go:embed template/*
var documentXml embed.FS

func CreateMarkdownDocx(path string, markdown string) (string, error) {
	outputFilePath := filepath.Join(path, "docx_markdown.docx")
	file, err := os.Create(outputFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	if err := template.CreateDocxTemplate(zipFile); err != nil {
		return "", err
	}
	docxer := docxWriter(markdown)

	if err := utils.CreateDocx(documentXml, zipFile, docxer); err != nil {
		return "", err
	}
	return outputFilePath, err
}

func docxWriter(markdownText string) utils.DocxWriter {

	return func(fileContent string) string {
		return applyStyle(fileContent, markdownText)
	}
}

func applyStyle(fileContent string, markdownText string) string {
	content := fileContent

	scanner := bufio.NewScanner(strings.NewReader(markdownText))
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmedLine, "#$# "):
			content = strings.Replace(content, "{{TITLE}}", strings.Replace(line, "#$# ", "", -1), -1)
		case strings.HasPrefix(trimmedLine, "# "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading1", strings.Replace(line, "# ", "", -1)), -1)
		case strings.HasPrefix(trimmedLine, "## "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading2", strings.Replace(line, "## ", "", -1)), -1)
		case strings.HasPrefix(trimmedLine, "### "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading3", strings.Replace(line, "### ", "", -1)), -1)
		case strings.HasPrefix(trimmedLine, "#### "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading4", strings.Replace(line, "#### ", "", -1)), -1)
		case strings.HasPrefix(trimmedLine, "##### "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading5", strings.Replace(line, "##### ", "", -1)), -1)
		case strings.HasPrefix(trimmedLine, "###### "):
			content = strings.Replace(content, "{{SECTION}}", newSection("Heading6", strings.Replace(line, "###### ", "", -1)), -1)
		default:
			content = strings.Replace(content, "{{SECTION}}", newSection("Normal", line), -1)
		}
	}
	content = strings.ReplaceAll(content, "{{TITLE}}", "")
	content = strings.ReplaceAll(content, "{{SECTION}}", "")
	return content

}

func newSection(style string, body string) string {
	return fmt.Sprintf(`<w:p>
      <w:pPr>
        <w:pStyle w:val="%s" />
        <w:bidi w:val="0" />
        <w:spacing w:before="0" w:after="140" />
        <w:jc w:val="left" />
        <w:rPr></w:rPr>
      </w:pPr>
      <w:r>
        <w:rPr></w:rPr>
        <w:t xml:space="preserve">%s</w:t>
      </w:r>
    </w:p>
    {{SECTION}}
    `, style, body)
}
