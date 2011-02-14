// Copyright 2009 The XGB Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"x-go-binding.googlecode.com/hg/xgb"
)

func main() {
    log("init")
    otherwm()
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		fmt.Printf("cannot connect: %v\n", err)
		os.Exit(1)
	}

	for {
		reply, err := c.WaitForEvent()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("event %T\n", reply)
		switch event := reply.(type) {
		case xgb.ExposeEvent:
		case xgb.KeyReleaseEvent:
        case xgb.MappingNotifyEvent:
		    //fmt.Printf("event %s\n", event.Window)
		}
	}

	c.Close()
}

func otherwm( ) bool {
}

func log( msg string ) {
    println("dwm: "+msg)
}
