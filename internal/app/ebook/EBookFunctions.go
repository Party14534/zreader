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

func ReadEBook(ebook ebooktype.EBook, page int) (string, error) {
    switch ebook.Type {
        case ebooktype.EPUB:
            return epub.ReadPage(ebook, page)            
    }

    return "", fmt.Errorf("Ebook type not supported\n")
}

func ReadEBookChunks(ebook ebooktype.EBook, page int) ([]string, []int, error) {
    switch ebook.Type {
        case ebooktype.EPUB:
            return epub.ReadPageChunks(ebook, page)            
    }

    return nil, nil, fmt.Errorf("Ebook type not supported\n")
}
