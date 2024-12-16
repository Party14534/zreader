package ereader

import (
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Party14534/zReader/internal/app/ebook"
	ebooktype "github.com/Party14534/zReader/internal/app/ebook/ebookType"
)

type C = layout.Context
type D = layout.Dimensions

var chapterNumber int
var currentBook ebooktype.EBook
var chapterProgress []int
var numberOfChapters int

var chapterChunks []string
var imageIndices []int

var textWidth unit.Dp = 550
var marginWidth unit.Dp
var fontSize unit.Sp = 35
var smallScrollStepSize unit.Dp = 50
var largeScrollStepSize unit.Dp = 50
var scrollY unit.Dp = 0
var labelStyles []material.LabelStyle
var atBottom bool = false

var textColor uint8 = 255
var backgroundColor uint8 = 0
var darkModeTextColor uint8 = 255
var darkModeBackgroundColor uint8 = 0
var lightModeTextColor uint8 = 0
var lightModeBackgroundColor uint8 = 255
var isDarkMode bool = true


func StartReader(book ebooktype.EBook, chapter int) {
    chapterNumber = chapter
    numberOfChapters = len(book.Chapters)
    chapterProgress = make([]int, len(book.Chapters))

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

            // Drawing to screen
            paint.Fill(&ops, color.NRGBA{R: backgroundColor,
                        G: backgroundColor, B: backgroundColor, A: 255})

            flexCol := layout.Flex {
                Axis: layout.Vertical,
                Spacing: layout.SpaceStart,
            }

            flexCol.Layout(gtx,
                layout.Rigid(
                    func(gtx C) D{
                        numberFontSize := fontSize / 2
                        if numberFontSize < 0 { numberFontSize = 0 }
                        chapterNumber := material.Label(theme, numberFontSize, strconv.Itoa(chapterNumber) + " ")
                        chapterNumber.Font.Typeface = "RobotoMono Nerd Font"

                        chapterNumber.Alignment = text.End
                        chapterNumber.Color = color.NRGBA{R: textColor,
                                    G: textColor, B: textColor, A: 255}
                        return chapterNumber.Layout(gtx)
                    },
                ),
            )
            
            layoutList(gtx, &ops)

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
                Name: "-",
            },
            key.Filter {
                Name: "=",
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
                chapterProgress[chapterNumber] = int(scrollY)

                chapterNumber++ 
                if chapterNumber >= numberOfChapters { 
                    chapterNumber = numberOfChapters - 1
                } else {
                    scrollY = 0
                    readChapter(theme)
                }
            }

        case key.Name("H"):
            if ev.State == key.Release { 
                chapterProgress[chapterNumber] = int(scrollY)

                chapterNumber--
                if chapterNumber < 0 { 
                    chapterNumber = 0 
                } else {
                    scrollY = 0
                    readChapter(theme)
                }
            }
            
        case key.Name("J"):
            if !atBottom { scrollY += smallScrollStepSize }

        case key.Name("K"):
            scrollY -= smallScrollStepSize 
            if scrollY < 0 { scrollY = 0 }

        case key.Name("D"):
            if ev.State == key.Release && !atBottom { 
                scrollY += largeScrollStepSize 
            }

        case key.Name("U"):
            if ev.State == key.Release { 
                scrollY -= largeScrollStepSize 
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

    chapterChunks, imageIndices, err = ebook.ReadEBookChunks(currentBook, chapterNumber)
    if err != nil { panic(err) }

    buildPageLayout(theme)

    // Set to previous scroll
    scrollY = unit.Dp(chapterProgress[chapterNumber])
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

func buildPageLayout(theme *material.Theme) {
    
    labelStyles = labelStyles[:0]
    for _, chunk := range chapterChunks {
        label := material.Label(theme, unit.Sp(fontSize), chunk)

        label.Alignment = text.Middle
        label.Color = color.NRGBA{R: textColor, G: textColor, B: textColor, A: 255}
        label.Font.Typeface = "RobotoMono Nerd Font"

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
        Position: layout.Position{
            Offset: int(scrollY),
        },
    }

    imageIndex := 0
    visList.Layout(gtx, len(chapterChunks), func(gtx C, i int) D {
          // Render each item in the list
          return pageMargins.Layout(gtx, func(gtx C) D{
              if imageIndex < len(imageIndices) && i == imageIndices[imageIndex] {
                  imageIndex++
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

