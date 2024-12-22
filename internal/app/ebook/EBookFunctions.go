package ebook

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	bookstate "github.com/Party14534/zreader/internal/app/ebook/bookState"
	ebooktype "github.com/Party14534/zreader/internal/app/ebook/ebookType"
	"github.com/Party14534/zreader/internal/app/ebook/epub"
)

func LoadFile(ebookPath, dest string) (ebooktype.EBook, error) {
    extension := path.Ext(ebookPath)

    var ebook ebooktype.EBook
    var err error

    /* Don't do this so users can rebuild the ebook if they need to
    Check if the ebook has been read before
    pathPieces := strings.Split(ebookPath, string(os.PathSeparator))
    bookPath := filepath.Join(dest, pathPieces[len(pathPieces) - 1])

    ebook, err = GetEBookMetaData(bookPath)
    if err == nil {
        return ebook, nil
    }
    */

    switch extension {
        case ".epub":
            err = epub.LoadEpubBook(ebookPath, dest, &ebook)            
        default:
            err = fmt.Errorf("Ebook type not supported\n")
    }

    // Save the metadata to not do it again
    if err == nil {
        err = SaveEBookMetaData(ebook)
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

func GetEBookHistory(ebook ebooktype.EBook) (bookstate.BookState, error) {
    var state bookstate.BookState
    stateJson, err := os.ReadFile(filepath.Join(ebook.Dest, "state.json"))
    if err != nil {
         return state, err
    }

    err = json.Unmarshal(stateJson, &state)
    return state, err
}

func SaveEBookHistory(ebook ebooktype.EBook, state bookstate.BookState) error {
    jsonState, err := json.Marshal(state)
    if err != nil {
        return err
    }

    err = os.WriteFile(filepath.Join(ebook.Dest, "state.json"), jsonState, 0644)
    return err
}

func GetEBookMetaData(bookPath string) (book ebooktype.EBook, err error) {
    metaJson, err := os.ReadFile(filepath.Join(bookPath, "metadata.json"))
    if err != nil { return book, err }

    err = json.Unmarshal(metaJson, &book)
    return book, err
}

func SaveEBookMetaData(ebook ebooktype.EBook) error {
    metaJson, err := json.Marshal(ebook)
    if err != nil { return err }

    err = os.WriteFile(filepath.Join(ebook.Dest, "metadata.json"), metaJson, 0644)
    return err
}
