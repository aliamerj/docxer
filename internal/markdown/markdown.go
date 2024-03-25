package markdown

import (
	"archive/zip"
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func applyStyle(fileContent, markdownText string) string {
	content := fileContent
	var titleSection, bodySections strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(markdownText))
	for scanner.Scan() {
		line := escapeXMLChars(scanner.Text())
		style, placeholder := determineStyle(line)
		cleanText := strings.Replace(line, placeholder, "", -1)

		if style == "Title" {
			titleSection.Reset()
			titleSection.WriteString(cleanText)
		} else {
			bodySections.WriteString(newSection(style, cleanText))
		}
	}

	content = strings.Replace(content, "{{TITLE}}", titleSection.String(), 1)
	content = strings.Replace(content, "{{SECTION}}", bodySections.String(), 1)

	return content
}
func newSection(style, body string) string {
	formattedBody := processTextFormatting(body)
	if !strings.Contains(formattedBody, "<w:r>") && !strings.Contains(formattedBody, "<w:rPr>") {
		formattedBody = fmt.Sprintf("<w:r><w:t>%s</w:t></w:r>", formattedBody)
	}
	return fmt.Sprintf(`<w:p><w:pPr><w:pStyle w:val="%s"/></w:pPr>%s</w:p>`, style, formattedBody)
}

func processTextFormatting(text string) string {
	boldItalicPattern := regexp.MustCompile(`\*\*\*(.+?)\*\*\*`)
	boldPattern := regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicPattern := regexp.MustCompile(`\*(.+?)\*`)

	if strings.ContainsAny(text, "*") {
		// Split the text into words
		words := strings.Split(text, " ")
		var processedWords []string

		for _, w := range words {
			// Check for bold and italic
			if boldItalicPattern.MatchString(w) {
				processedWord := boldItalicPattern.ReplaceAllString(w, `<w:r><w:rPr><w:b/><w:i/></w:rPr><w:t xml:space="preserve">$1 </w:t></w:r>`)
				processedWords = append(processedWords, processedWord)
			} else if boldPattern.MatchString(w) { // Check for bold
				processedWord := boldPattern.ReplaceAllString(w, `<w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">$1 </w:t></w:r>`)
				processedWords = append(processedWords, processedWord)
			} else if italicPattern.MatchString(w) { // Check for italic
				processedWord := italicPattern.ReplaceAllString(w, `<w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">$1 </w:t></w:r>`)
				processedWords = append(processedWords, processedWord)
			} else { // No styling
				processedWords = append(processedWords, `<w:r><w:t xml:space="preserve">`+w+` </w:t></w:r>`)
			}
		}

		return strings.Join(processedWords, " ")
	}

	// If no markdown styling markers are found, return the text wrapped in WordprocessingML tags for plain text
	return `<w:r><w:t xml:space="preserve">` + text + `</w:t></w:r>`
}

func determineStyle(line string) (style string, placeholder string) {
	switch {
	case strings.HasPrefix(line, "#$# "):
		return "Title", "#$# "
	case strings.HasPrefix(line, "# "):
		return "Heading1", "# "
	case strings.HasPrefix(line, "## "):
		return "Heading2", "## "
	case strings.HasPrefix(line, "### "):
		return "Heading3", "### "
	case strings.HasPrefix(line, "#### "):
		return "Heading4", "#### "
	case strings.HasPrefix(line, "##### "):
		return "Heading5", "##### "
	case strings.HasPrefix(line, "###### "):
		return "Heading6", "###### "
	default:
		return "Normal", ""
	}
}
func escapeXMLChars(text string) string {
	// Replace special XML characters with their escape sequences
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")
	return text
}
