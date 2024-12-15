package epub

import (
	"os"
	"strings"

	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
	"github.com/Party14534/zReader/internal/app/parser"
)

func ReadChapter(ebook ebooktype.EBook, chapter int) (string, error) {
    // Get the name of the file and load the file in
    chapterName := ebook.Chapters[chapter]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + chapterName)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))

    parsedText := ElementsToText(htmlElements)
    
    return parsedText, err
}

func ReadChapterChunks(ebook ebooktype.EBook, chapter int) (parsedChunks []string, imageIndices []int, err error) {
    // Get the name of the file and load the file in
    chapterName := ebook.Chapters[chapter]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + chapterName)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))
    
    fileDomainSlice := strings.Split(chapterName, string(os.PathSeparator))
    fileDomain := ""
    for i, dir := range fileDomainSlice {
        if i == len(fileDomainSlice) - 1 { continue }

        fileDomain += dir + string(os.PathSeparator)
    }

    imageDomain := ebook.Dest + string(os.PathSeparator) + fileDomain

    chunks, indices := ElementsToChunks(htmlElements, imageDomain)
    
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

            // Handle double periods in image path
            if element.TagCode == parser.Img {
                chunk = removeDoublePeriodsInPath(chunk) 
            }

            parsedChunks = append(parsedChunks, chunk)
            chunk = ""
        }
    }

    return parsedChunks, imageIndices
}

func removeDoublePeriodsInPath(path string) string {
    dirs := strings.Split(path, string(os.PathSeparator))
    newPath := ""

    for i := 0; i < len(dirs); i++ {
        if strings.Compare(dirs[i], "..") == 0 {
            dirs = append(dirs[:i-1], dirs[i+1:]...)
            continue
        } 
    }

    for i, dir := range dirs {
        separator := string(os.PathSeparator)
        if i == len(dirs) - 1 {
            separator = ""
        }

        newPath += dir + separator
    }

    return newPath
}

