package ereader

import (
	"log"
	"math"

	"gioui.org/io/key"
)

func handleMenuEvents(gtx *C) {
    if inFileMenu { return }
    // Handle key events
    for {
        keyEvent, ok := gtx.Event(
            key.Filter {
                Name: "J",
            },
            key.Filter {
                Name: "K",
            },
            key.Filter {
                Name: key.NameReturn,
            },
            key.Filter{
                Name: key.NameSpace,
            },
        )
        if !ok { break }

        ev, ok := keyEvent.(key.Event)
        if !ok { break }

        switch ev.Name {
        case key.Name("J"):
            if ev.State == key.Release {
                menuBookIndex--
                menuBookIndex = max(menuBookIndex, 0)
            }

        case key.Name("K"):
            if ev.State == key.Release {
                menuBookIndex++
                menuBookIndex = min(menuBookIndex, len(menuBooks))
            }

        case key.NameReturn:
            if ev.State != key.Release {
                if menuBookIndex != len(menuBooks) {
                    initializeEReader(menuBooks[menuBookIndex])
                    switched = true
                } else {
                    err := openFileViewer()
                    if err != nil {
                        log.Println(err)
                    }
                }
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
            }
        }
    }

}

func handleEReaderEvents(gtx *C) {
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
                Name: "[",
            },
            key.Filter {
                Name: "]",
            },
            key.Filter {
                Name: key.NameUpArrow,
            },
            key.Filter {
                Name: key.NameDownArrow,
            },
            key.Filter {
                Name: key.NameEscape,
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
                    pageNumber = 0
                    readChapter()
                }
            }

        case key.Name("H"):
            if ev.State == key.Release { 
                chapterNumber--
                if chapterNumber < 0 { 
                    chapterNumber = 0 
                } else {
                    pageNumber = 0
                    readChapter()
                }
            }
            
        case key.Name("J"):
            if ev.State == key.Release { 
                pageNumber--
                if pageNumber < 0 {
                    chapterNumber--
                    if chapterNumber < 0 { 
                        chapterNumber = 0 
                        pageNumber = 0
                    } else {
                        if chapterLengths[chapterNumber] == 0 {
                            pageNumber = math.MaxInt32
                        } else {
                            pageNumber = len(pageLabelStyles[chapterNumber]) - 1
                        }
                        readChapter()
                    }
                } else { scrollY = 0 }
            }

        case key.Name("K"):
            if ev.State == key.Release { 
                pageNumber++
                if pageNumber >= len(pageLabelStyles[chapterNumber]) {
                    chapterNumber++
                    if chapterNumber >= numberOfChapters { 
                        chapterNumber = numberOfChapters - 1
                        pageNumber--
                    } else {
                        pageNumber = 0
                        readChapter()
                    }
                } else { scrollY = 0 }
            }

        case key.Name("["):
            if ev.State == key.Release {
                fontScale -= 0.05
                if fontScale < 0.05 { fontScale = 0.05 }
                buildPageLayout()
                needToBuildPages = true
                clearChapterLengths()
            }

        case key.Name("]"):
            if ev.State == key.Release {
                fontScale += 0.05
                buildPageLayout()
                needToBuildPages = true
                clearChapterLengths()
            }

        case key.NameUpArrow:
            scrollY = max(0, scrollY - scrollStep)

        case key.NameDownArrow:
            if beforeEnd { scrollY += scrollStep }

        case key.NameEscape:
            quitEReader()
            initializeMenu()
            switched = true

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

                buildPageLayout()
            }
        }
    }
}
