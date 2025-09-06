package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"
)

// Page holds all the information we need to generate a new
// HTML page from a text file on the filesystem.
type Page struct {
	TextFilePath string
	TextFileName string
	HTMLPagePath string
	Content      string
}

func main() {
	// Parse command line flags
	fileName := flag.String("file", "first-post.txt", "the path to the text file to read")
	flag.Parse()

	// Read the contents of first-post.txt
	fileContents, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}

		// Generate HTML filename by replacing .txt with .html
		htmlFileName := strings.Replace(*fileName, ".txt", ".html", 1)
	
	// Create a Page struct with the content
	page := Page{
		TextFilePath: *fileName,
		TextFileName: *fileName,
		HTMLPagePath: htmlFileName,
		Content:      string(fileContents),
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