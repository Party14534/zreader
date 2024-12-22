package ereader

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
)

type C = layout.Context
type D = layout.Dimensions

type pageStyleIndices struct {
    start int
    end int
}

var theme *material.Theme

// Menu values
var menuBooks []ebooktype.EBook
var basePath string = ""
var menuBookIndex int

// EBook metadata
var chapterNumber int
var currentBook ebooktype.EBook
var numberOfChapters int

// Current chapter data
var chapterChunks [][]string
var chunkTypes [][]int
var chapterLengths []unit.Dp
var labelStyles []material.LabelStyle
var pageLabelStyles [][]pageStyleIndices
// TODO: var chapterProgress []int
var pageNumber int = 0

// Page design
var textWidth unit.Dp = 550
var marginWidth unit.Dp
// var smallScrollStepSize unit.Dp = 50
// var largeScrollStepSize unit.Dp = 50
var fontScale unit.Sp = 1.0
var ereaderFont string = "RobotoMono Nerd Font, Times New Roman"
var textColor uint8 = 255
var backgroundColor uint8 = 0
var darkModeTextColor uint8 = 255
var darkModeBackgroundColor uint8 = 0
var lightModeTextColor uint8 = 0
var lightModeBackgroundColor uint8 = 255

// Booleans
var isDarkMode bool = true
// var atBottom bool = false
var needToBuildPages bool
var justStarted bool = true

