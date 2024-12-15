package ebooktype

const (
    EPUB = iota + 1
    MOGI
)

type EBook struct {
    // File Info
    Dest string
    ContentFilePath string
    Type int
    
    // Book metadata
    Title string
    Creator string
    Language string

    // Chapters
    Chapters []string
}

