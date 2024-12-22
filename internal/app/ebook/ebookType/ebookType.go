package ebooktype

const (
    EPUB = iota + 1
    MOGI
)

type EBook struct {
    // File Info
    Dest string `json:"Dest"`
    ContentFilePath string `json:"ContentFilePath"`
    Type int `json:"Type"`
    
    // Book metadata
    Title string `json:"Title"`
    Creator string `json:"Creator"`
    Language string `json:"Language"`
    Cover string `json:"Cover"`

    // Chapters
    Chapters []Chapter `json:"Chapters"`
}

type Chapter struct {
    Path string `json:"Path"`
    ID string `json:"ID"`
}

