package ereader

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Party14534/zReader/internal/app/ebook"
	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
	"github.com/Party14534/zReader/internal/app/parser"
)

type C = layout.Context
type D = layout.Dimensions

var chapterNumber int
var currentBook ebooktype.EBook
var numberOfChapters int

var chapterProgress []unit.Dp
var chapterChunks [][]string
var chunkTypes [][]int
var chapterLengths []unit.Dp

var textWidth unit.Dp = 550
var marginWidth unit.Dp
var fontScale unit.Sp = 1.0
var smallScrollStepSize unit.Dp = 50
var largeScrollStepSize unit.Dp = 50
var labelStyles []material.LabelStyle
var atBottom bool = false

var textColor uint8 = 255
var backgroundColor uint8 = 0
var darkModeTextColor uint8 = 255
var darkModeBackgroundColor uint8 = 0
var lightModeTextColor uint8 = 0
var lightModeBackgroundColor uint8 = 255
var isDarkMode bool = true
var ereaderFont string = "RobotoMono Nerd Font, Times New Roman"

func StartReader(book ebooktype.EBook, chapter int) {
    chapterNumber = chapter
    numberOfChapters = len(book.Chapters)
    chapterProgress = make([]unit.Dp, len(book.Chapters))
    chapterChunks = make([][]string, len(book.Chapters))
    chunkTypes = make([][]int, len(book.Chapters))
    chapterLengths = make([]unit.Dp, len(book.Chapters))

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

    smallScrollStepSize = 32

    // Read first chapter
    readChapter(theme)

    if isDarkMode {
        textColor = darkModeTextColor
        backgroundColor = darkModeBackgroundColor
    } else {
        textColor = lightModeTextColor
        backgroundColor = lightModeBackgroundColor
    }

    for {
        switch e := window.Event().(type) {
        case app.DestroyEvent:
            return e.Err

        case app.FrameEvent:
            // This graphics context is used for managing the rendering state
            gtx := app.NewContext(&ops, e)

            largeScrollStepSize = unit.Dp(float32(gtx.Constraints.Max.Y) * 0.95)

            // Handle key events
            handleKeyEvents(&gtx, theme)

            flexCol := layout.Flex {
                Axis: layout.Vertical,
                Spacing: layout.SpaceStart,
            }

            // Before drawing get chapter length so we can reset the ops after
            if chapterLengths[chapterNumber] == 0 {
                chapterLengths[chapterNumber] = getChapterLength(gtx, &ops)

                gtx.Reset()
            } 

            /*
                Prevent overscroll here instead of in the the event handler
                to prevent going over 100% chapter progress when making the
                font scale smaller then going back to a previous chapter whose
                progress was at the previous chapters old length
            */
            if chapterLengths[chapterNumber] > 0 &&
                chapterProgress[chapterNumber] > chapterLengths[chapterNumber] {
                chapterProgress[chapterNumber] = chapterLengths[chapterNumber]
            }

            // Drawing to screen
            paint.Fill(&ops, color.NRGBA{R: backgroundColor,
                        G: backgroundColor, B: backgroundColor, A: 255})
            
            layoutList(gtx, &ops)

            // Chapter number
            flexCol.Layout(gtx,
                layout.Rigid(
                    func(gtx C) D{
                        chapterNumber := material.Body2(theme, strconv.Itoa(chapterNumber) + " ")
                        chapterNumber.Font.Typeface = font.Typeface(ereaderFont)

                        chapterNumber.TextSize *= fontScale
                        chapterNumber.Alignment = text.End
                        chapterNumber.Color = color.NRGBA{R: textColor,
                                    G: textColor, B: textColor, A: 255}
                        return chapterNumber.Layout(gtx)
                    },
                ),
            )

            // Chapter completion percent
            flexCol.Layout(gtx,
                layout.Rigid(
                    func(gtx C) D{
                        percentage := 0.0
                        if chapterLengths[chapterNumber] <= 0 {
                            percentage = 1.0
                        } else {
                            percentage = float64(chapterProgress[chapterNumber] /
                                chapterLengths[chapterNumber])
                        }
                        completion := " " + fmt.Sprintf("%.0f", percentage * 100) + 
                            "%"
                        chapterNumber := material.Body2(theme, completion)
                        chapterNumber.Font.Typeface = font.Typeface(ereaderFont)

                        chapterNumber.TextSize *= fontScale
                        chapterNumber.Alignment = text.Start
                        chapterNumber.Color = color.NRGBA{R: textColor,
                                    G: textColor, B: textColor, A: 255}
                        return chapterNumber.Layout(gtx)
                    },
                ),           
            )

            // Pass the drawing operations to the GPU
            e.Frame(gtx.Ops)
        }
    }
}

func handleKeyEvents(gtx *layout.Context, theme *material.Theme) {
    // Handle key events
    for {
        keyEvent, ok := gtx.Event(
            key.Filter {
                Name: "L",
            },
            key.Filter {
                Name: "H",
            },
            key.Filter {
                Name: "J",
            },
            key.Filter {
                Name: "K",
            },
            key.Filter {
                Name: "D",
                Required: key.ModCtrl,
            },
            key.Filter {
                Name: "U",
                Required: key.ModCtrl,
            },
            key.Filter {
                Name: "[",
            },
            key.Filter {
                Name: "]",
            },
            key.Filter{
                Name: key.NameSpace,
            },
        )
        if !ok { break }

        ev, ok := keyEvent.(key.Event)
        if !ok { break }

        switch ev.Name {
        case key.Name("L"):
            if ev.State == key.Release { 
                chapterNumber++ 
                if chapterNumber >= numberOfChapters { 
                    chapterNumber = numberOfChapters - 1
                } else {
                    readChapter(theme)
                }
            }

        case key.Name("H"):
            if ev.State == key.Release { 
                chapterNumber--
                if chapterNumber < 0 { 
                    chapterNumber = 0 
                } else {
                    readChapter(theme)
                }
            }
            
        case key.Name("J"):
            if ev.State == key.Release { continue }
            if !atBottom { 
                chapterProgress[chapterNumber] += smallScrollStepSize 

                // Prevent overscroll
                if chapterLengths[chapterNumber] > 0 && chapterProgress[chapterNumber] > chapterLengths[chapterNumber] {
                    chapterProgress[chapterNumber] = chapterLengths[chapterNumber]
                }
            }

        case key.Name("K"):
            if ev.State == key.Release { continue }
            chapterProgress[chapterNumber] -= smallScrollStepSize 
            if chapterProgress[chapterNumber] < 0 { chapterProgress[chapterNumber] = 0 }

        case key.Name("D"):
            if ev.State == key.Release && !atBottom { 
                chapterProgress[chapterNumber] += largeScrollStepSize 

                // Prevent overscroll
                if chapterLengths[chapterNumber] > 0 && chapterProgress[chapterNumber] > chapterLengths[chapterNumber] {
                    chapterProgress[chapterNumber] = chapterLengths[chapterNumber]
                }
            }

        case key.Name("U"):
            if ev.State == key.Release { 
                chapterProgress[chapterNumber] -= largeScrollStepSize 
                if chapterProgress[chapterNumber] < 0 { chapterProgress[chapterNumber] = 0 }
            }

        case key.Name("["):
            if ev.State == key.Release {
                fontScale -= 0.05
                if fontScale < 0.05 { fontScale = 0.05 }
                buildPageLayout(theme)
                resetScrollsAfterScaleChange(fontScale + 0.05)
                clearChapterLengths()
            }

        case key.Name("]"):
            if ev.State == key.Release {
                fontScale += 0.05
                buildPageLayout(theme)
                resetScrollsAfterScaleChange(fontScale - 0.05)
                clearChapterLengths()
            }

        case key.NameSpace:
            if ev.State == key.Release {
                isDarkMode = !isDarkMode

                if isDarkMode {
                    textColor = darkModeTextColor
                    backgroundColor = darkModeBackgroundColor
                } else {
                    textColor = lightModeTextColor
                    backgroundColor = lightModeBackgroundColor
                }

                buildPageLayout(theme)
            }

        }
    }
}

func readChapter(theme *material.Theme) {
    var err error
    
    if chapterChunks[chapterNumber] == nil {
        chapterChunks[chapterNumber], chunkTypes[chapterNumber], err =
            ebook.ReadEBookChunks(currentBook, chapterNumber)
        if err != nil { panic(err) }
    }

    buildPageLayout(theme)

    // Set to previous scroll
    chapterProgress[chapterNumber] = unit.Dp(chapterProgress[chapterNumber])
}

func buildPageLayout(theme *material.Theme) {
    labelStyles = labelStyles[:0]
    for i, chunk := range chapterChunks[chapterNumber] {
        var label material.LabelStyle
        switch chunkTypes[chapterNumber][i] {
        case parser.H1:
            label = material.H1(theme, chunk)
        case parser.H2:
            label = material.H2(theme, chunk)
        case parser.H3:
            label = material.H3(theme, chunk)
        case parser.H4:
            label = material.H4(theme, chunk)
        case parser.H5:
            label = material.H5(theme, chunk)
        case parser.H6:
            label = material.H6(theme, chunk)
        case parser.Img:
            // Separating in case I need to make changes to label
            label = material.Body1(theme, chunk)
        default:
            label = material.Body1(theme, chunk)
        }

        label.Font.Typeface = font.Typeface(ereaderFont)
        label.TextSize *= fontScale
        label.LineHeight *= fontScale // Idk if this does anything but it feels nice to have
        label.Alignment = text.Middle

        label.Color = color.NRGBA{R: textColor, G: textColor, B: textColor, A: 255}

        labelStyles = append(labelStyles, label)
    }
}

// layoutList handles the layout of the list
func layoutList(gtx layout.Context, ops *op.Ops) {
    textWidth = unit.Dp(gtx.Constraints.Max.X) * 0.95
    marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 2
    pageMargins := layout.Inset {
        Left:   marginWidth,
        Right:  marginWidth,
        Top: unit.Dp(12),
        Bottom: unit.Dp(12),
    }

    var visList = layout.List {
        Axis: layout.Vertical,
        Position: layout.Position {
            Offset: int(chapterProgress[chapterNumber]),
        },
    }

    visList.Layout(gtx, len(labelStyles), func(gtx C, i int) D {
        // Render each item in the list
        return pageMargins.Layout(gtx, func(gtx C) D{
            if chunkTypes[chapterNumber][i] == parser.Img {
                // Draw the image in the window
                return layout.Center.Layout(gtx, func(gtx C) D {
                    // Build image 
                    img := loadImage(labelStyles[i].Text)
                    imgOp := paint.NewImageOp(img)
                    imgOp.Filter = paint.FilterNearest
                    imgOp.Add(ops)

                    scale := 2
                    fScale := float32(scale)
                    imgSize := img.Bounds().Size()
                    imgSize.X *= scale
                    imgSize.Y *= scale

                    op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), 
                        f32.Pt(fScale, fScale))).Add(ops)
                    paint.PaintOp{}.Add(gtx.Ops)

                    return layout.Dimensions{Size: imgSize}
                })
            } else {
                return labelStyles[i].Layout(gtx)
            }
        },)
    },)

    // To prevent overscroll
    atBottom = !visList.Position.BeforeEnd
}

func getChapterLength(gtx C, ops *op.Ops) unit.Dp {
    pageMargins := layout.Inset {
        Left:   marginWidth,
        Right:  marginWidth,
        Top: unit.Dp(12),
        Bottom: unit.Dp(12),
    }

    var visList = layout.List {
        Axis: layout.Vertical,
        Position: layout.Position {
            Offset: int(chapterProgress[chapterNumber]),
        },
    }

    var emptyList = layout.List {
        Axis: layout.Vertical,
    }

    emptyList.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
        return visList.Layout(gtx, len(labelStyles), func(gtx C, i int) D {
            return pageMargins.Layout(gtx, func(gtx C) D {
                if chunkTypes[chapterNumber][i] == parser.Img {
                    // Draw the image in the window
                    return layout.Center.Layout(gtx, func(gtx C) D {
                        // Build image 
                        img := loadImage(labelStyles[i].Text)
                        imgOp := paint.NewImageOp(img)
                        imgOp.Filter = paint.FilterNearest
                        imgOp.Add(ops)

                        scale := 2
                        fScale := float32(scale)
                        imgSize := img.Bounds().Size()
                        imgSize.X *= scale
                        imgSize.Y *= scale

                        op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), 
                            f32.Pt(fScale, fScale))).Add(ops)
                        paint.PaintOp{}.Add(gtx.Ops)

                        return layout.Dimensions{Size: imgSize}
                    })
                } else {
                    return labelStyles[i].Layout(gtx)
                }
            })
        })
    })

    return unit.Dp(emptyList.Position.OffsetLast * -1)
}

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

func resetScrollsAfterScaleChange(previousScale unit.Sp) {
    ratio := fontScale / previousScale
    for i := range chapterProgress {
        chapterProgress[i] *= unit.Dp(ratio)
    }
}

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

