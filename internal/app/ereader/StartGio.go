package ereader

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var CurrentPageText string

func StartGio() {
   go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func run(window *app.Window) error {
    theme := material.NewTheme()
    var ops op.Ops

    for {
        switch e := window.Event().(type) {
            case app.DestroyEvent:
                return e.Err
            case app.FrameEvent:
                // Graphics context
                gtx := app.NewContext(&ops, e)

                title := material.H6(theme, CurrentPageText)
                title.Alignment = text.Middle
                
                title.Layout(gtx)
                e.Frame(gtx.Ops)
        }
    }
}
