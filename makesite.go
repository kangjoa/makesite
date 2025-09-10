package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
	Title        string
	IsMarkdown   bool
}

func main() {
	// Parse command line flags
	fileName := flag.String("file", "", "the path to the text file to read")
	dirName := flag.String("dir", "", "the directory to find all .txt files")
	flag.Parse()

	// Process the file or directory
	switch {
	case *fileName != "":
		processSingleFile(*fileName)
	case *dirName != "":
		processDirectory(*dirName)
	default:
		processSingleFile("first-post.txt")
	}
}

func processSingleFile(fileName string) {
	// Read the contents of the file
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// Extract first line for title
	lines := strings.Split(string(fileContents), "\n")
	firstLine := strings.TrimSpace(lines[0])

	// Remove Markdown header syntax
	title := strings.TrimLeft(firstLine, "# ")
	// Remove trailing punctuation
	title = strings.TrimRight(title, ".,!?;:")

	// Generate HTML filename by replacing .txt or .md with .html
	var baseName string
	if strings.HasSuffix(fileName, ".md") {
		baseName = strings.Replace(fileName, ".md", ".html", 1)
	} else {
		baseName = strings.Replace(fileName, ".txt", ".html", 1)
	}

	// Process content based on file type
	var processedContent string
	if strings.HasSuffix(fileName, ".md") {
		// Parse Markdown to HTML
		processedContent = parseMarkdown(string(fileContents))
	} else {
		processedContent = string(fileContents)
	}

	// Create a Page struct with the content
	page := Page{
		TextFilePath: fileName,
		TextFileName: fileName,
		HTMLPagePath: baseName,
		Content:      processedContent,
		Title:        title,
		IsMarkdown:   strings.HasSuffix(fileName, ".md"),
	}

	// Create a new template in memory named "template.tmpl"
	t := template.Must(template.New("template.tmpl").Funcs(template.FuncMap{
		// Allow html tags to be rendered
		"html": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprintf("%v", value))
		},
	}).ParseFiles("template.tmpl"))

	// Print the rendered template to stdout
	fmt.Println("=== Rendered HTML ===")
	err = t.Execute(os.Stdout, page)
	if err != nil {
		panic(err)
	}

	// Create a new HTML file
	newFile, err := os.Create(baseName)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Execute the template and save to file
	err = t.Execute(newFile, page)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== HTML written to first-post.html ===", baseName)

}

func parseMarkdown(content string) string {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		panic(err)
	}
	return buf.String()
}

const (
	Reset = "\033[0m"
	Green = "\033[32m"
	Bold  = "\033[1m"
)

func processDirectory(dirName string) {
	// Find all .txt files in the directory
	files, err := os.ReadDir(dirName)
	if err != nil {
		panic(err)
	}

	// Print found .txt files to stdout
	fmt.Println("\n=== Found .txt and .md files: ===")
	fileCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") || strings.HasSuffix(file.Name(), ".md") {
			fmt.Println(file.Name())
			fileCount++
		}
	}

	// Process each file
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") || strings.HasSuffix(file.Name(), ".md") {
			processSingleFile(filepath.Join(dirName, file.Name()))
		}
	}

	fmt.Println(Green + Bold + "Success!" + Reset + " Generated " + Bold + fmt.Sprintf("%d", fileCount) + Reset + " pages.")
}
