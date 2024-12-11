package epub

import (
	"os"

	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
	"github.com/Party14534/zReader/internal/app/parser"
)

func ReadPage(ebook ebooktype.EBook, page int) (string, error) {
    // Get the name of the file and load the file in
    pageName := ebook.Pages[page]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + pageName)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))

    parsedText := ElementsToText(htmlElements)
    
    return parsedText, err
}

func ElementsToText(elements []parser.HTMLElement) (parsedText string) {
    for _, element := range elements {
        //if element.TagCode == parser.Undefined { continue }

        parsedText = parsedText + element.Content
    }

    return parsedText
}
