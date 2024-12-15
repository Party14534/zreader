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
    for i := 0; i < len(elements); i++ {
        element := elements[i]

        //if element.TagCode == parser.Undefined { continue }

        // If it is an inline element or next element is
        // an inline element do not add a line break
        lineBreak := "\n"
        _, isInline := parser.HtmlInlineTagMap[element.Tag]
        nextIsInline := false
        if i < len(elements) - 1 {
            _, nextIsInline = parser.HtmlInlineTagMap[elements[i+1].Tag]
        }
        if isInline || nextIsInline  {
            lineBreak = ""
        }

        parsedText = parsedText + element.Content + lineBreak
    }

    return parsedText
}

