package main

import (
	"fmt"
	"html/template"
	"os"
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
	// Read the contents of first-post.txt
	fileContents, err := os.ReadFile("first-post.txt")
	if err != nil {
		panic(err)
	}

	// Create a Page struct with the content
	page := Page{
		TextFilePath: "first-post.txt",
		TextFileName: "first-post",
		HTMLPagePath: "first-post.html",
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
	newFile, err := os.Create("first-post.html")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Execute the template and save to file
	err = t.Execute(newFile, page)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== HTML written to first-post.html ===")
}