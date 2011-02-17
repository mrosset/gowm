/* 
 *  Copyright (C) 2005 Bodo Giannone
 *  
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *  
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *  
 *  You should have received a copy of the GNU General Public License
 *  along with this program; if not, write to the Free Software
 *  Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.
 */

#include <stdio.h>
#include <stdlib.h>
#include <X11/Xlib.h>


int xmax, ymax;
static int screen;
static XSetWindowAttributes attr; 

Window root;
Display *display;

int main(void)
{
    display    = XOpenDisplay(0); 

    if (!display) {
        printf("error: could not open display\n");
        exit(1);
    }
    XWindowChanges wc;
    wc.border_width = 1;
    screen = DefaultScreen(display);
    xmax   = DisplayWidth (display, screen);
    ymax   = DisplayHeight(display, screen);
    root   = DefaultRootWindow(display);
    Colormap cmap = DefaultColormap(display, screen);
    XColor color;
    XAllocNamedColor(display, cmap, "turquoise", &color, &color);
    attr.event_mask = SubstructureNotifyMask;
    XChangeWindowAttributes(display, root, CWEventMask, &attr);

    while ( 1 )
    {
        static XWindowAttributes client;
        static XEvent event;
        XNextEvent(display, &event);
        if(event.type == MapNotify)
        {
            printf("got a notify event");
            fflush(stderr);
            XGetWindowAttributes(display, event.xmap.window, &client);
            //XMoveWindow(display, event.xmap.window, xmax/2 - client.width/2, ymax/2 - client.height/2);
            XResizeWindow(display, event.xmap.window, xmax/2, ymax);
            XSetWindowBorder(display, event.xmap.window, color.pixel);
            XConfigureWindow(display, event.xmap.window, CWBorderWidth, &wc);
        }
    }
}
