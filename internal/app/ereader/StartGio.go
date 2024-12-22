package ereader

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Party14534/zreader/internal/app/ebook"
	ebooktype "github.com/Party14534/zreader/internal/app/ebook/ebookType"
	"github.com/Party14534/zreader/internal/app/parser"
)

var readingBook bool = false
var switched bool = false

func StartReader(book ebooktype.EBook) {
    theme = material.NewTheme()

    initializeEReader(book)

    go func() {
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

func StartMenu() {
    theme = material.NewTheme()

    initializeMenu()

    go func() {
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
    var ops op.Ops

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
            if readingBook { quitEReader() }
            return e.Err

        case app.ConfigEvent: 
            needToBuildPages = true 

        case app.FrameEvent:
            // This graphics context is used for managing the rendering state
            gtx := app.NewContext(&ops, e)

            // If the user is reading a book draw the ereader screen, else
            // they are on the main menu
            if readingBook {
                drawEReaderScreen(&gtx, &ops, theme)
            } else {
                drawMenuScreen(&gtx, &ops, theme)
            }

            // If we switched screens we need to redraw directly after
            if switched {
                window.Invalidate()
            }

            // Pass the drawing operations to the GPU
            e.Frame(gtx.Ops)
        }
    }
}

func readChapter() {
    var err error

    scrollY = 0
    
    if chapterChunks[chapterNumber] == nil {
        chapterChunks[chapterNumber], chunkTypes[chapterNumber], err =
            ebook.ReadEBookChunks(currentBook, chapterNumber)
        if err != nil { panic(err) }
    }

    buildPageLayout()
}

func buildPageLayout() {
    labelStyles = labelStyles[:0]
    for i, chunk := range chapterChunks[chapterNumber] {
        var label material.LabelStyle
        switch chunkTypes[chapterNumber][i] {
        case parser.H1:
            label = material.H1(theme, chunk)
            label.Alignment = text.Middle
        case parser.H2:
            label = material.H2(theme, chunk)
            label.Alignment = text.Middle
        case parser.H3:
            label = material.H3(theme, chunk)
            label.Alignment = text.Middle
        case parser.H4:
            label = material.H4(theme, chunk)
            label.Alignment = text.Middle
        case parser.H5:
            label = material.H5(theme, chunk)
            label.Alignment = text.Middle
        case parser.H6:
            label = material.H6(theme, chunk)
            label.Alignment = text.Middle
        case parser.Img:
            // Separating in case I need to make changes to label
            label = material.Body1(theme, chunk)
        default:
            label = material.Body1(theme, chunk)
            label.Alignment = text.Middle
        }

        label.Font.Typeface = font.Typeface(ereaderFont)
        label.TextSize *= fontScale
        label.LineHeight *= fontScale // Idk if this does anything but it feels nice to have

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
        Alignment: layout.Start,
        Position: layout.Position{
            Offset: scrollY,
        },
    }

    indices := pageLabelStyles[chapterNumber][pageNumber]
    page := labelStyles[indices.start:indices.end]
    
    pageMargins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return visList.Layout(gtx, len(page), func(gtx C, index int) D {
            // Render each item in the list
            i := indices.start + index 
            if chunkTypes[chapterNumber][i] == parser.Img {
                // Draw the image in the window
                return layout.Center.Layout(gtx, func(gtx C) D {
                    // Build image 
                    img, err := loadImage(labelStyles[i].Text)
                    if err != nil {
                        return labelStyles[i].Layout(gtx)
                    }

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
    })

    beforeEnd = visList.Position.BeforeEnd
}

func buildChapterPages(gtx C, ops *op.Ops) unit.Dp {
    // Clear page label styles
    pageLabelStyles[chapterNumber] = pageLabelStyles[chapterNumber][:0]

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
            Offset: int(0),
        },
    }

    var emptyList = layout.List {
        Axis: layout.Vertical,
    }

    // Prevents text from being cut off
    maxPageHeight := gtx.Constraints.Max.Y - int(gtx.Metric.PxPerDp * 24)

    height := 0
    startIndex := 0
    emptyList.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
        return pageMargins.Layout(gtx, func(gtx C) D {
            return visList.Layout(gtx, len(labelStyles), func(gtx C, i int) D {
                isLast := i + 1 == len(labelStyles)
                if chunkTypes[chapterNumber][i] == parser.Img {
                    // Draw the image in the window
                    return layout.Center.Layout(gtx, func(gtx C) D {
                        // Build image 
                        img, err := loadImage(labelStyles[i].Text)
                        if err != nil {
                            return labelStyles[i].Layout(gtx)
                        }
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

                        if height + imgSize.Y > maxPageHeight {
                            if height == 0 { 
                                pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                    start: startIndex,
                                    end: i + 1,
                                })

                                startIndex = i + 1
                            } else {
                                pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                    start: startIndex,
                                    end: i,
                                })

                                if isLast {
                                    pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                        start: i,
                                        end: i + 1,
                                    })
                                }

                                startIndex = i
                                height = imgSize.Y
                            }
                        } else {
                            height += imgSize.Y
                            if isLast {
                                pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                    start: startIndex,
                                    end: i + 1,
                                })
                            }
                        }

                        return layout.Dimensions{Size: imgSize}
                    })
                } else {
                    dim := labelStyles[i].Layout(gtx)

                    if height + dim.Size.Y > maxPageHeight {
                        if height == 0 { 
                            pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                start: startIndex,
                                end: i + 1,
                            })

                            startIndex = i + 1
                        } else {
                            pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                start: startIndex,
                                end: i,
                            })

                            if isLast {
                                pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                    start: i,
                                    end: i + 1,
                                })
                            }

                            startIndex = i
                            height = dim.Size.Y
                        }
                    } else {
                        height += dim.Size.Y
                        if isLast {
                            pageLabelStyles[chapterNumber] = append(pageLabelStyles[chapterNumber], pageStyleIndices{
                                start: startIndex,
                                end: i + 1,
                            })
                        }
                    }

                    return dim
                }
            })
        })
    })

    return unit.Dp(emptyList.Position.OffsetLast * -1)
}

