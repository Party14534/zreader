package ereader

import (
	"fmt"
	"image"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"gioui.org/unit"
	"github.com/Party14534/zreader/internal/app/ebook"
	bookstate "github.com/Party14534/zreader/internal/app/ebook/bookState"
	ebooktype "github.com/Party14534/zreader/internal/app/ebook/ebookType"
	"github.com/Party14534/zreader/internal/pkg"
	"github.com/rymdport/portal/filechooser"
)

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
        return img, err
	}
	return img, nil
}

func openFileViewer() error {
    inFileMenu = true

    options := filechooser.OpenFileOptions{Multiple: false}
    file, err := filechooser.OpenFile("", "Select EBook", &options)
    if err == nil && len(file) != 0 {
        // Load the file into the reader
        err = loadBook(file[0][7:])
    }

    inFileMenu = false
    return err
}

func loadBook(filename string) (error) {
    // Format the filename
    fmt.Println(filename)
    filename, err := url.QueryUnescape(filename)
    if err != nil {
        return err
    }
    fmt.Println(filename)

    basePath, err := pkg.GetAppDataDir("zreader")
    if err != nil {
        return err
    }

    book, err := ebook.LoadFile(filename, filepath.Join(basePath, ".ebookfiles"))
    if err != nil {
        return err
    }

    switched = true
    initializeEReader(book)

    return nil
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
    readingBook = true

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
    readingBook = false

    err := getEBooks()
    if err != nil {
        log.Print(err)
    }
}

func getEBooks() error {
    // Clear old books
    menuBooks = menuBooks[:0]

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

