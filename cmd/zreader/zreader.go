package main

import (
	"os"
	"path/filepath"

	"github.com/Party14534/zReader/internal/app/ebook"
	"github.com/Party14534/zReader/internal/app/ereader"
	"github.com/Party14534/zReader/internal/pkg"
)

func main() {
    if len(os.Args) == 2 {
        epubPath := os.Args[1]

        // Get the base path of the ebook metadata
        basePath, err := pkg.GetAppDataDir("zreader")
        if err != nil {
            panic(err)
        }

        book ,err := ebook.LoadFile(epubPath, filepath.Join(basePath, ".ebookfiles"))
        if err != nil {
            panic(err)
        }

        ereader.StartReader(book)
    } else {
        ereader.StartMenu()
    }

}
