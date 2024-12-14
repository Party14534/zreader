package ereader

import (
	"image/color"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Party14534/zReader/internal/app/ebook"
	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
)

type C = layout.Context
type D = layout.Dimensions

var CurrentPageText string

var pageNumber int
var currentBook ebooktype.EBook

var textWidth unit.Dp = 550
var marginWidth unit.Dp
var fontSize unit.Sp = 35
var scrollStepSize unit.Dp = 50
var scrollY unit.Dp = 0
var pageText string
var pageChunks []string
var labelStyles []material.LabelStyle


func StartReader(book ebooktype.EBook, page int) {
    pageNumber = page

    go func() {
        currentBook = book
        window := new(app.Window)
        window.Option(app.Title("ZReader"))

        err := run(window)
        if err != nil {
            log.Fatal(err)
        }

        os.Exit(0)
    }()

    app.Main()
}

func run(window *app.Window) error {
    theme := material.NewTheme()
    var ops op.Ops

    var nextPageButton widget.Clickable
    var previousPageButton widget.Clickable

    // Read first page
    readPage(theme)

    for {
        switch e := window.Event().(type) {
        case app.DestroyEvent:
            return e.Err

        case app.FrameEvent:
            // This graphics context is used for managing the rendering state
            gtx := app.NewContext(&ops, e)

            scrollStepSize = unit.Dp(float32(gtx.Constraints.Max.Y) * 0.95)

            // Handle key events
            handleKeyEvents(&gtx, theme)

            // Logic
            if nextPageButton.Clicked(gtx) {
                pageNumber++
            }
            if previousPageButton.Clicked(gtx) {
                if pageNumber > 0 {
                    pageNumber--
                }
            }

            paint.Fill(&ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

            // Drawing to screen
            flexCol := layout.Flex {
                Axis: layout.Vertical,
                Spacing: layout.SpaceStart,
            }

            /*var visList = layout.List{
                Axis: layout.Vertical,
                Position: layout.Position{
                    Offset: int(scrollY),
                },
            }

            textWidth = unit.Dp(gtx.Constraints.Max.X) * 0.95
            marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 2
            pageMargins := layout.Inset {
                Left:   marginWidth,
                Right:  marginWidth,
                Top:    unit.Dp(0),
                Bottom: unit.Dp(0),
            }

            
            flexCol.Layout(gtx,
                layout.Rigid(
                    func(gtx C) D {
                        return visList.Layout(gtx, 1, 
                            func(gtx C, index int) D {
                                return pageMargins.Layout(gtx, 
                                    func(gtx layout.Context) layout.Dimensions {
                                        page := material.Label(theme, fontSize, pageText)

                                        // Change the position of the label
                                        page.Alignment = text.Middle

                                        page.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

                                        // Draw the title to the context
                                        return page.Layout(gtx)
                                        return material.Label(theme, fontSize, "Hello").Layout(gtx)
                                    },
                                )
                            },
                        )
                    },
                ),
            )*/

            flexCol.Layout(gtx,
                layout.Rigid(
                    func(gtx C) D{
                        numberFontSize := fontSize / 2
                        if numberFontSize < 0 { numberFontSize = 0 }
                        chapterNumber := material.Label(theme, numberFontSize, strconv.Itoa(pageNumber) + " ")
                        chapterNumber.Font.Typeface = "monospace"

                        chapterNumber.Alignment = text.End
                        chapterNumber.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
                        return chapterNumber.Layout(gtx)
                    },
                ),
            )
            
            layoutList(gtx, theme)

            // Pass the drawing operations to the GPU
            e.Frame(gtx.Ops)
        }
    }
}

func handleKeyEvents(gtx *layout.Context, theme *material.Theme) {
    // Handle key events
    for {
        keyEvent, ok := gtx.Event(
            key.Filter{
                Name: "L",
            },
            key.Filter{
                Name: "H",
            },
            key.Filter{
                Name: "J",
            },
            key.Filter{
                Name: "K",
            },
            key.Filter{
                Name: "-",
            },
            key.Filter{
                Name: "=",
            },
        )
        if !ok { break }

        ev, ok := keyEvent.(key.Event)
        if !ok { break }

        switch ev.Name {
        case key.Name("L"):
            if ev.State == key.Release { 
                pageNumber++ 
                scrollY = 0
                readPage(theme)
            }

        case key.Name("H"):
            if ev.State == key.Release { 
                pageNumber--
                if pageNumber < 0 { 
                    pageNumber = 0 
                } else {
                    scrollY = 0
                    readPage(theme)
                }
            }
            
        case key.Name("J"):
            if ev.State == key.Release { scrollY += scrollStepSize }

        case key.Name("K"):
            if ev.State == key.Release { 
                scrollY -= scrollStepSize 
                if scrollY < 0 { scrollY = 0 }
            }

        case key.Name("-"):
            if ev.State == key.Release {
                fontSize -= 2
                if fontSize < 0 { fontSize = 0 }
                buildPageLayout(theme)
            }

        case key.Name("="):
            if ev.State == key.Release {
                fontSize += 2
                buildPageLayout(theme)
            }
        }
    }

}

func readPage(theme *material.Theme) {
    var err error
    pageText, err = ebook.ReadEBook(currentBook, pageNumber)
    if err != nil { panic(err) }
    pageChunks = chunkString(pageText)
    buildPageLayout(theme)
}

func chunkString(input string) (chunks []string) {
    start := 0
    alreadyChunked := false
	for i := 1; i < len(input); i++ {
        if input[i] == '\n' && !alreadyChunked {
            chunks = append(chunks, input[start:i])
            start = i
            alreadyChunked = true
        } else { alreadyChunked = false }
	}

    chunks = append(chunks, input[start:])

	return chunks
}

func buildPageLayout(theme *material.Theme) {
    
    labelStyles = labelStyles[:0]
    for _, chunk := range pageChunks {
        label := material.Label(theme, unit.Sp(fontSize), chunk)

        label.Alignment = text.Middle
        label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
        label.Font.Typeface = "monospace"

        labelStyles = append(labelStyles, label)
    }
}

// layoutList handles the layout of the list
func layoutList(gtx layout.Context, theme *material.Theme) {
    textWidth = unit.Dp(gtx.Constraints.Max.X) * 0.95
    marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 2
    pageMargins := layout.Inset {
        Left:   marginWidth,
        Right:  marginWidth,
        Top:    unit.Dp(12),
        Bottom: unit.Dp(0),
    }

    var visList = layout.List {
        Axis: layout.Vertical,
        Position: layout.Position{
            Offset: int(scrollY),
        },
    }

    visList.Layout(gtx, len(pageChunks), func(gtx layout.Context, i int) layout.Dimensions {
            // Render each item in the list
            return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return pageMargins.Layout(gtx, 
                    func(gtx C) D {
                        return labelStyles[i].Layout(gtx)
                    },  
                )
            },)
        },
    )

}
