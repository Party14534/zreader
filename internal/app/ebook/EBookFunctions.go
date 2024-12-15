package ebook

import (
	"fmt"
	"path"

	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
	"github.com/Party14534/zReader/internal/app/ebook/epub"
)

func LoadFile(ebookPath, dest string) (ebooktype.EBook, error) {
    extension := path.Ext(ebookPath)

    var ebook ebooktype.EBook
    var err error

    switch extension {
        case ".epub":
            err = epub.LoadEpubBook(ebookPath, dest, &ebook)            
        default:
            err = fmt.Errorf("Ebook type not supported\n")
    }

    return ebook, err
}

func ReadEBook(ebook ebooktype.EBook, chapter int) (string, error) {
    switch ebook.Type {
        case ebooktype.EPUB:
            return epub.ReadChapter(ebook, chapter)            
    }

    return "", fmt.Errorf("Ebook type not supported\n")
}

func ReadEBookChunks(ebook ebooktype.EBook, chapter int) ([]string, []int, error) {
    switch ebook.Type {
        case ebooktype.EPUB:
            return epub.ReadChapterChunks(ebook, chapter)            
    }

    return nil, nil, fmt.Errorf("Ebook type not supported\n")
}
