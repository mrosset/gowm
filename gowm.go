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
	createTestWindow()
	checkWM()
	setupScreen()
	run()
	shutdown()
}

func createTestWindow() {
	win := conn.NewId()
	gc := conn.NewId()
	conn.CreateWindow(0, win, root, 150, 150, 200, 200, 0, 0, 0, 0, nil)
	conn.ChangeWindowAttributes(win, xgb.CWBackPixel|xgb.CWEventMask,
		[]uint32{
			screen.WhitePixel,
			xgb.EventMaskExposure | xgb.EventMaskKeyRelease | xgb.EventMaskEnterWindow,
		})
	conn.CreateGC(gc, win, xgb.GCForeground, []uint32{screen.WhitePixel})
	conn.MapWindow(win)
	setupWindow(win)
}

func setupScreen() {
	qtr, err := conn.QueryTree(root)
	//atom_desktop, _ := conn.InternAtom(true, "_NET_WM_DESKTOP")
	if err != nil {
		lfatalf(err.String())
	}
	for _, child := range qtr.Children {
		lprintf("found window %v", child)
		setupWindow(child)
	}
}

func setupWindow(win xgb.Id) {
	attr, err := conn.GetWindowAttributes(win)
	if err != nil {
		lprintf("couldnt get attribute for %v", win)
		return
	}
	if !attr.OverrideRedirect || attr.MapState == xgb.MapStateViewable {
		setBorderColor(win, "pink")
		setBorderWidth(win, 2)
		conn.MapWindow(win)

	}
}

func getWmDestop() {
}

func connectToX() {
	lprintf("connecting to %v", envdisplay)
	var err os.Error
	conn, err = xgb.Dial(envdisplay)
	if err != nil {
		lfatalf("%v", err)
	}
	lprintf("connected to %v", envdisplay)
}


func registerEvents() {
	conn.ChangeWindowAttributes(root, xgb.CWEventMask, []uint32{
		xgb.EventMaskSubstructureRedirect |
			xgb.EventMaskSubstructureNotify |
			xgb.EventMaskStructureNotify |
			xgb.EventMaskLeaveWindow |
			xgb.EventMaskEnterWindow |
			xgb.EventMaskPropertyChange |
			xgb.EventMaskKeyPress,
	})

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

func setupwin(win *xgb.Id) {
	//cond.ConfigureWindow(win, xdb.ConfigureWindow:
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
		case xgb.CreateNotifyEvent:
			lprintf("create event for %v", event.Window)
			setupWindow(event.Window)
		case xgb.ExposeEvent:
		case xgb.EnterNotifyEvent:
			lprintf("root window is %v", root)
			lprintf("setting focus to %v", event.Event)
			conn.SetInputFocus(byte(0), event.Event, event.Time)
		case xgb.MapRequestEvent:
			lprintf("%T from %v", reply, event.Window)
			setupWindow(event.Window)
		}
	}
}

func setBorderWidth(win xgb.Id, width uint32) {
	lprintf("setting %v border width to %v", win, width)
	conn.ConfigureWindow(win, 16, []uint32{width})
}

func setBorderColor(win xgb.Id, color string) {
	lprintf("setting %v border color to %v", win, color)
	conn.ChangeWindowAttributes(win, xgb.CWBorderPixel|xgb.CWEventMask, []uint32{getColorByName(color), xgb.EventMaskEnterWindow})
}

func lprintf(format string, i ...interface{}) {
	logger.Printf("gowm: "+format, i...)
}

func lfatalf(format string, i ...interface{}) {
	logger.Fatalf("gowm: error: "+format, i...)
}

func getColorByName(color string) uint32 {
	cr, _ := conn.AllocNamedColor(screen.DefaultColormap, color)
	return cr.Pixel
}
