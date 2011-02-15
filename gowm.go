package main

import (
	"fmt"
	"os"
	"x-go-binding.googlecode.com/hg/xgb"
)

func main() {
	log("main")
	c := connectX()
	screen := c.DefaultScreen()
	c.ChangeWindowAttributes(screen.Root, xgb.CWEventMask, []uint32{xgb.EventMaskSubstructureNotify|xgb.EventMaskEnterWindow})
	win := createwin(c)
	setupwin(win, c)
	c.MapWindow(win)
	qtr, _ := c.QueryTree(screen.Root)
	log(fmt.Sprintf("our win %s ", win))
	for _, value := range qtr.Children {
		log(fmt.Sprintf("%s ", value))
	}
	run(c)
	c.Close()
	os.Exit(0)
}


func run(con *xgb.Conn) {
	for {
		reply, err := con.WaitForEvent()
		log(fmt.Sprintf("%T", reply))
		if err != nil {
			log(fmt.Sprintf("error %v", err))
			os.Exit(1)
		}
		switch event := reply.(type) {
		case xgb.ExposeEvent:
		case xgb.ConfigureNotifyEvent:
			nc, _ := con.AllocNamedColor(con.DefaultScreen().DefaultColormap, "turquoise")
			con.ChangeWindowAttributes(event.Window, xgb.CWBorderPixel, []uint32{nc.Pixel})
			//setupwin(event.Window, con)
		}
	}
}

func createwin(con *xgb.Conn) xgb.Id {
	log("create win")
	win := con.NewId()
	gc := con.NewId()
	s := con.DefaultScreen()
	con.CreateWindow(0, win, s.Root, 1200, 150, 200, 200, xgb.WindowClassInputOutput, 0, 0, 0, nil)
	con.ChangeWindowAttributes(win, xgb.CWBackPixel|xgb.CWEventMask,
		[]uint32{
			s.BlackPixel,
			xgb.EventMaskExposure | xgb.EventMaskKeyRelease | xgb.EventMaskKeyPress | xgb.EventMaskEnterWindow,
		})
	con.CreateGC(gc, win, xgb.GCForeground, []uint32{s.WhitePixel})
	return win
}


func setupwin(win xgb.Id, con *xgb.Conn) {
	log("setupwin")
	nc, _ := con.AllocNamedColor(con.DefaultScreen().DefaultColormap, "turquoise")
	con.ChangeWindowAttributes(win, xgb.CWBorderPixel, []uint32{nc.Pixel})
}

func connectX() *xgb.Conn {
	log("connect")
	log("connecting to " + os.Getenv("DISPLAY"))
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		fmt.Printf("cannot connect: %v\n", err)
		os.Exit(1)
	}
	log("connected to " + os.Getenv("DISPLAY"))
	return c
}

func log(msg string) {
	println("gowm: " + msg)
}
