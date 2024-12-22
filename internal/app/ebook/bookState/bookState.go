package bookstate

type BookState struct {
    ChapterNumber int `json:"ChapterNummber"`
    PageNumber int `json:"PageNumber"`
    FontScale float64 `json:"FontScale"`
    DarkMode bool `json:"DarkMode"`
}

