package ereader

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)
var flexCol layout.Flex = layout.Flex {
    Axis: layout.Vertical,
    Spacing: layout.SpaceStart,
}

var menuFlexCol layout.Flex = layout.Flex {
    Axis: layout.Vertical,
    Spacing: layout.SpaceBetween,
    Alignment: layout.Middle,
}

var bookMenuMargins layout.Inset = layout.Inset {
    Left: unit.Dp(12),
    Right: unit.Dp(12),
}

var loadEBookButton widget.Clickable
var inFileMenu bool = false


func drawMenuScreen(gtx *layout.Context, ops *op.Ops, theme *material.Theme) {
    handleMenuEvents(gtx)

    // Drawing to screen
    paint.Fill(ops, color.NRGBA{R: backgroundColor,
                G: backgroundColor, B: backgroundColor, A: 255})

    // Handle button press
    if loadEBookButton.Clicked(*gtx) {
        err := openFileViewer()
        if err != nil {
            log.Println(err)
        }
    }

    // Main layout
    menuFlexCol.Layout(*gtx,
        // Title
        layout.Rigid( func(gtx C) D {
            title := material.H1(theme, "zreader")
            title.Font.Typeface = font.Typeface(ereaderFont)

            title.TextSize *= fontScale
            title.Alignment = text.Middle
            title.Color = color.NRGBA{R: textColor,
                        G: textColor, B: textColor, A: 255}
            return title.Layout(gtx)
        },),

        // Books
        layout.Rigid( func(gtx C) D {
            return bookMenuMargins.Layout(gtx, func(gtx C) D {
                if menuBookIndex < len(menuBooks) {
                    return layout.Center.Layout(gtx, func(gtx C) D {
                        // Build image 
                        coverPath := filepath.Join(menuBooks[menuBookIndex].Dest,
                            menuBooks[menuBookIndex].Cover)
                        img, err := loadImage(coverPath)
                        if err != nil {
                            message := material.H3(theme, menuBooks[menuBookIndex].Title)
                            message.Font.Typeface = font.Typeface(ereaderFont)
                            message.TextSize *= fontScale
                            message.Alignment = text.Middle
                            message.Color = color.NRGBA{R: textColor,
                                        G: textColor, B: textColor, A: 255}
                            return message.Layout(gtx)
                        }

                        imgOp := paint.NewImageOp(img)
                        imgOp.Filter = paint.FilterNearest
                        imgOp.Add(ops)

                        scale := 1
                        fScale := float32(scale)
                        imgSize := img.Bounds().Size()
                        imgSize.X *= scale
                        imgSize.Y *= scale

                        op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), 
                            f32.Pt(fScale, fScale))).Add(ops)
                        paint.PaintOp{}.Add(gtx.Ops)

                        return layout.Dimensions{Size: imgSize}
                    }) 
                } 

                return layout.Inset {
                    Left: unit.Dp(500),
                    Right: unit.Dp(500),
                }.Layout(gtx, func(gtx C) D {
                    loadBtn := material.Button(theme, &loadEBookButton, "Load Book")
                    loadBtn.Font.Typeface = font.Typeface(ereaderFont)
                    return loadBtn.Layout(gtx)
                })
            })
        }),

        // Book Info
        layout.Rigid( func(gtx C) D {
            if menuBookIndex < len(menuBooks) {
                book := menuBooks[menuBookIndex]
                infoText := book.Title + "\n" + book.Creator
                info := material.Body1(theme, infoText)
                info.Font.Typeface = font.Typeface(ereaderFont)

                info.TextSize *= fontScale
                info.Alignment = text.Middle
                info.Color = color.NRGBA{R: textColor,
                            G: textColor, B: textColor, A: 255}
                return info.Layout(gtx)
            }
            return layout.Spacer{Height: 1}.Layout(gtx)
        },),

        // Spacer
        layout.Rigid( func(gtx C) D {
            return layout.Spacer{Height: 1}.Layout(gtx)
        }),
    )

}

func drawEReaderScreen(gtx *layout.Context, ops *op.Ops, theme *material.Theme) {
    // Handle key events
    handleEReaderEvents(gtx)

    // Before drawing get chapter length so we can reset the ops after
    if chapterLengths[chapterNumber] == 0 || needToBuildPages {
        chapterLengths[chapterNumber] = buildChapterPages(*gtx, ops)
        pageNumber = min(pageNumber, len(pageLabelStyles[chapterNumber]) - 1)

        gtx.Reset()
        needToBuildPages = false
    }

    // Drawing to screen
    paint.Fill(ops, color.NRGBA{R: backgroundColor,
                G: backgroundColor, B: backgroundColor, A: 255})
    
    layoutList(*gtx, ops)

    // Chapter number
    flexCol.Layout(*gtx,
        layout.Rigid(
            func(gtx C) D{
                chapterNumber := material.Body2(theme, currentBook.Chapters[chapterNumber].ID + " ")
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
    flexCol.Layout(*gtx,
        layout.Rigid(
            func(gtx C) D{
                percentage := 0.0
                if chapterLengths[chapterNumber] <= 0 {
                    percentage = 1.0
                } else {
                    percentage = float64(pageNumber) /
                        float64(len(pageLabelStyles[chapterNumber]) - 1)
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
}
