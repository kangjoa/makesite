package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
	Title        string
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
		title := strings.TrimSpace(lines[0])  // First line, trimmed of whitespace

		// Generate HTML filename by replacing .txt with .html
		htmlFileName := strings.Replace(fileName, ".txt", ".html", 1)
	
	// Create a Page struct with the content
	page := Page{
		TextFilePath: fileName,
		TextFileName: fileName,
		HTMLPagePath: htmlFileName,
		Content:      string(fileContents),
		Title:        title,
	}

	// Create a new template in memory named "template.tmpl"
	t := template.Must(template.New("template.tmpl").ParseFiles("template.tmpl"))

	// Print the rendered template to stdout
	fmt.Println("=== Rendered HTML ===")
	err = t.Execute(os.Stdout, page)
	if err != nil {
		panic(err)
	}

	// Create a new HTML file
	newFile, err := os.Create(htmlFileName)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Execute the template and save to file
	err = t.Execute(newFile, page)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== HTML written to first-post.html ===", htmlFileName)
	
}

func processDirectory(dirName string) {
	// Find all .txt files in the directory
	files, err := os.ReadDir(dirName)
	if err != nil {
		panic(err)
	}

		// Print found .txt files to stdout
		fmt.Println("\n=== Found .txt files: ===")
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".txt") {
				fmt.Println(file.Name())
			}
		}

	// Process each file
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			processSingleFile(filepath.Join(dirName, file.Name()))
		}
	}
}