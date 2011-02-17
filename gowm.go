package main

import (
	"os"
	"log"
	"x-go-binding.googlecode.com/hg/xgb"
)

var (
	bcolor     = "turquoise"
	conn       *xgb.Conn
	screen     *xgb.ScreenInfo
	root       xgb.Id
	logger     = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	envdisplay = os.Getenv("DISPLAY")
)


func main() {
	lprintf("starting")
	bcolor = "red"
	connectToX()
	screen = conn.DefaultScreen()
	root = screen.Root
	registerEvents()
	checkWM()
	setupScreen()
	run()
	shutdown()
}

func setupScreen() {
	qtr, err := conn.QueryTree(root)
	atom_desktop, _ := conn.InternAtom(true, "_NET_WM_DESKTOP")
	if err != nil {
		lfatalf(err.String())
	}
	for _, child := range qtr.Children {
		lprintf("found window %v", child)
		attr, err := conn.GetWindowAttributes(child)
		if err != nil {
			lprintf("couldnt get attribute for %v", child)
			continue
		}
		if !attr.OverrideRedirect || attr.MapState == xgb.MapStateViewable {
			conn.ChangeProperty(xgb.PropModeReplace, child, atom_desktop.Atom, xgb.AtomString, 8, []uint8{0})
		}

	}
}

func getWmDestop() {
}

func connectToX() {
	var e os.Error
	lprintf("connect to %v", envdisplay)
	conn, e = xgb.Dial(envdisplay)
	if e != nil {
		lfatalf("connecting to %v", envdisplay)
	}
	lprintf("connected to %v", envdisplay)
}

func setupscreen() {
}


func rootFlags() []uint32 {
	return []uint32{
		xgb.EventMaskSubstructureRedirect |
			xgb.EventMaskSubstructureNotify |
			xgb.EventMaskStructureNotify |
			xgb.EventMaskLeaveWindow |
			xgb.EventMaskEnterWindow |
			xgb.EventMaskPropertyChange,
	}
}

func dRootFlags() []uint32 {
	return []uint32{
		xgb.EventMaskSubstructureRedirect |
			xgb.EventMaskSubstructureNotify |
			xgb.EventMaskStructureNotify |
			xgb.EventMaskLeaveWindow |
			xgb.EventMaskEnterWindow |
			xgb.EventMaskPropertyChange,
	}
}


func registerEvents() {
	conn.ChangeWindowAttributes(root, xgb.CWEventMask, dRootFlags())

}

func checkWM() {
	_, err := conn.WaitForEvent()
	if err != nil {
		lfatalf("Is a another window manager running?")
	}
}


func shutdown() {
	lprintf("closing %v", envdisplay)
	conn.Close()
	os.Exit(0)
}

func run() {
	conn.MapWindow(root)
	for {
		reply, err := conn.WaitForEvent()
		if err != nil {
			lfatalf("error : %v", err)
		}
		lprintf("event %T", reply)
		switch event := reply.(type) {
		case xgb.ExposeEvent:
		case xgb.MapRequestEvent:
			lprintf("%T from %v", reply, event.Window)
			conn.MapWindow(event.Window)
		}
	}
}

func lprintf(format string, i ...interface{}) {
	logger.Printf("gowm: "+format, i...)
}

func lfatalf(format string, i ...interface{}) {
	logger.Fatalf("gowm: error: "+format, i...)
}
