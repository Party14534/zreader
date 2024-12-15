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

func ReadPageChunks(ebook ebooktype.EBook, page int) (parsedChunks []string, imageIndices []int, err error) {
    // Get the name of the file and load the file in
    pageName := ebook.Pages[page]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + pageName)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))

    chunks, indices := ElementsToChunks(htmlElements, ebook.ContentFilePath)
    
    return chunks, indices, err
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

func ElementsToChunks(elements []parser.HTMLElement, imageDomain string) (parsedChunks []string, imageIndices []int) {
    var chunk string
    for i := 0; i < len(elements); i++ {
        element := elements[i]

        // If element is an inline element or next element is
        // an inline element do not chunk
        _, isInline := parser.HtmlInlineTagMap[element.Tag]
        nextIsInline := false 
        var nextCode int
        if i < len(elements) - 1 {
            nextCode, nextIsInline = parser.HtmlInlineTagMap[elements[i+1].Tag]
            if nextCode == parser.Img {
                nextIsInline = false
            }
        }

        if (isInline || nextIsInline) && element.TagCode != parser.Img {
            chunk += element.Content
        } else {
            if element.TagCode == parser.Img {
                // If we get to an image while inside an inline element 
                // create a new chunk
                if chunk != "" {
                    parsedChunks = append(parsedChunks, chunk)
                    chunk = ""
                }
                chunk += imageDomain
                imageIndices = append(imageIndices, len(parsedChunks))
            }

            chunk += element.Content
            parsedChunks = append(parsedChunks, chunk)
            chunk = ""
        }
    }

    return parsedChunks, imageIndices
}

