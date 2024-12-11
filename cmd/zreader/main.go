package main

import (
	"fmt"
	"os"

	"github.com/Party14534/zReader/internal/app/ebook"
	"github.com/Party14534/zReader/internal/app/ereader"
)

func main() {
    if len(os.Args) != 2 {
        panic(fmt.Errorf("Invalid number of arguments"))
    }

    epubPath := os.Args[1]
    book ,err := ebook.LoadFile(epubPath, ".ebookfiles")
    if err != nil {
        panic(err)
    }

    // Ebook is loaded in and we can now read the first page
    text, err := ebook.ReadEBook(book, 4)
    if err != nil {
        panic(err)
    }

    fmt.Println(text)

    ereader.CurrentPageText = text
}
