package epub

import (
	"os"

	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
)

func ReadPage(ebook ebooktype.EBook, page int) (string, error) {
    pageName := ebook.Pages[page]
    text, err := os.ReadFile(ebook.Dest + string(os.PathSeparator) + pageName)
    
    return string(text), err
}
