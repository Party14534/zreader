package ereader

import (
	"image"
	"log"
	"os"
	"path/filepath"

	"gioui.org/unit"
	"github.com/Party14534/zReader/internal/app/ebook"
	bookstate "github.com/Party14534/zReader/internal/app/ebook/bookState"
	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
	"github.com/Party14534/zReader/internal/pkg"
)

func loadImage(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}
	return img
}

func clearChapterLengths() {
    for i := range chapterLengths {
        chapterLengths[i] = 0
    }
}

func loadEbookHistory() {
    state, err := ebook.GetEBookHistory(currentBook)
    if err != nil { 
        return 
    }

    chapterNumber = state.ChapterNumber
    pageNumber = state.PageNumber
    fontScale = unit.Sp(state.FontScale)

    isDarkMode = state.DarkMode
    if isDarkMode {
        textColor = darkModeTextColor
        backgroundColor = darkModeBackgroundColor
    } else {
        textColor = lightModeTextColor
        backgroundColor = lightModeBackgroundColor
    }
}

func quitEReader() {
    state := bookstate.BookState{
        ChapterNumber: chapterNumber,
        PageNumber: pageNumber,
        FontScale: float64(fontScale),
        DarkMode: isDarkMode,
    }
    err := ebook.SaveEBookHistory(currentBook, state)
    if err != nil { panic(err) }
}

func initializeEReader(book ebooktype.EBook) {
    needToBuildPages = true
    justStarted = true

    currentBook = book
    numberOfChapters = len(book.Chapters)
    chapterChunks = make([][]string, len(book.Chapters))
    chunkTypes = make([][]int, len(book.Chapters))
    chapterLengths = make([]unit.Dp, len(book.Chapters))
    pageLabelStyles = make([][]pageStyleIndices, len(book.Chapters))

    loadEbookHistory()
    readChapter()
}

func initializeMenu() {
    err := getEBooks()
    if err != nil {
        log.Print(err)
    }
}

func getEBooks() error {
    // Get the location of the ebooks
    var err error
    if basePath == "" {
        basePath, err = pkg.GetAppDataDir("zreader")    
        if err != nil { return err }
    }

    bookPaths := filepath.Join(basePath, ".ebookfiles")

    // Get every folder in bookPaths
    entries, err := os.ReadDir(bookPaths)
    if err != nil { return err }

    for _, entry := range entries {
        if entry.IsDir() {
            book, err := getEBookMenuData(filepath.Join(bookPaths, entry.Name()))
            if err != nil { continue }

            menuBooks = append(menuBooks, book)
        }
    }

    return nil
}

func getEBookMenuData(path string) (book ebooktype.EBook, err error) {
    return ebook.GetEBookMetaData(path)
}

/*
func chunkString(input string) (chunks []string) {
    start := 0
    alreadyChunked := false
    for i := 1; i < len(input); i++ {
          if input[i] == '\n' && !alreadyChunked {
              chunks = append(chunks, input[start:i])
              start = i+1
              alreadyChunked = true
          } else { alreadyChunked = false }
    }

    chunks = append(chunks, input[start:])

	return chunks
} 
*/

