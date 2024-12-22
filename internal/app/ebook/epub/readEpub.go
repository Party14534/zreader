package epub

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ebooktype "github.com/Party14534/zreader/internal/app/ebook/ebookType"
	"github.com/Party14534/zreader/internal/app/parser"
)

func ReadChapter(ebook ebooktype.EBook, chapter int) (string, error) {
    // Get the name of the file and load the file in
    chapterName := ebook.Chapters[chapter]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + chapterName.Path)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))

    parsedText := ElementsToText(htmlElements)
    
    return parsedText, err
}

func ReadChapterChunks(ebook ebooktype.EBook, chapter int) (parsedChunks []string, imageIndices []int, err error) {
    // Get the name of the file and load the file in
    chapterName := ebook.Chapters[chapter]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + chapterName.Path)

    // Convert the html into html elements and parse it for the text
    htmlElements := parser.ParseHTML(string(text))
    
    fileDomainSlice := splitFilePath(chapterName.Path)
    fileDomain := ""
    for i, dir := range fileDomainSlice {
        if i == len(fileDomainSlice) - 1 { continue }

        fileDomain += dir + string(os.PathSeparator)
    }

    imageDomain := filepath.Join(ebook.Dest, fileDomain) + string(os.PathSeparator)

    chunks, types := ElementsToChunks(htmlElements, imageDomain)
    
    return chunks, types, err
}

func splitFilePath(name string) (slice []string) {
    startIndex := 0
    for i, r := range name {
        if r == '/' || r == '\\' {
            slice = append(slice, name[startIndex:i])
            startIndex = i + 1
        }
    } 
    
    if startIndex != len(name) - 1 {
        slice = append(slice, name[startIndex:len(name)])
    }

    return slice
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

func ElementsToChunks(elements []parser.HTMLElement, imageDomain string) (parsedChunks []string, chunkTypes []int) {
    var chunk string = ""
    var previousChunkType int
    for i := 0; i < len(elements); i++ {
        element := elements[i]

        // If element is an inline element or next element is an inline element do not chunk
        // Treat unknown elements as inline
        _, isInline := parser.HtmlInlineTagMap[element.Tag]
        
        if element.TagCode == parser.Undefined {
            isInline = true
        }

        nextIsInline := false 
        var nextCode int
        if i < len(elements) - 1 {
            nextCode, nextIsInline = parser.HtmlInlineTagMap[elements[i+1].Tag]

            // If element is an image it is not inline
            if nextCode == parser.Img {
                nextIsInline = false
            }
        }

        if (isInline || nextIsInline) && element.TagCode != parser.Img {
            chunk += element.Content
            previousChunkType = element.TagCode
        } else {
            if element.TagCode == parser.Img {
                // If we get to an image while inside an inline element 
                // create a new chunk
                if chunk != "" {
                    parsedChunks = append(parsedChunks, chunk)
                    chunkTypes = append(chunkTypes, previousChunkType)
                    chunk = ""
                }

                chunk += imageDomain
            }

            if element.Content == "" && chunk == "" {
                continue
            }

            chunk += element.Content

            // Handle double periods in image path
            if element.TagCode == parser.Img {
                chunk = removeDoublePeriodsInPath(chunk) 
            }

            parsedChunks = append(parsedChunks, chunk)
            chunkTypes = append(chunkTypes, element.TagCode)

            previousChunkType = element.TagCode
            chunk = ""
        }
    }

    return parsedChunks, chunkTypes
}

func removeDoublePeriodsInPath(path string) string {
    dirs := strings.Split(path, string(os.PathSeparator))
    newPath := ""

    for i := 0; i < len(dirs); i++ {
        if strings.Compare(dirs[i], "..") == 0 {
            dirs = append(dirs[:i-1], dirs[i+1:]...)
            continue
        } else if strings.Compare(dirs[i], ".") == 0 {
            fmt.Println(dirs)
            dirs = append(dirs[:i], dirs[i+1:]...)
            fmt.Println(dirs)
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

