// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"math"
	"io"
	"os"
	"net/http"
	"log"
	"strconv"
)

const (
	cells = 100                 // number of grid cells
	xyrange = 30.0                // axis ranges (-xyrange..+xyrange)
	angle = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	configs := make(map[string]string)
	if len(os.Args) > 1 && os.Args[1] == "web" {
		// Test with http://localhost:8000/?w=350&h=250&s=red&f=blue
		//!+http
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Print(err)
			}
			if f := r.Form["w"]; f != nil {
				configs["width"] = f[0]
			}
			if f := r.Form["h"]; f != nil {
				configs["height"] = f[0]
			}
			if f := r.Form["s"]; f != nil {
				configs["scolor"] = f[0]
			}
			if f := r.Form["f"]; f != nil {
				configs["fcolor"] = f[0]
			}

			w.Header().Set("Content-Type", "image/svg+xml")
			surface(w, configs)
		}
		http.HandleFunc("/", handler)
		//!-http
		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}
	surface(os.Stdout, configs)
}

func surface(out io.Writer, configs map[string]string) {
	width := 600
	height := 320
	scolor := "grey"
	fcolor := "white"
	if w, ok := configs["width"]; ok {
		if wt, err := strconv.Atoi(w); err == nil {
			width = wt
		}
	}
	if h, ok := configs["height"]; ok {
		if ht, err := strconv.Atoi(h); err == nil {
			height = ht
		}
	}
	if sc, ok := configs["scolor"]; ok {
		scolor = sc
	}
	if fc, ok := configs["fcolor"]; ok {
		fcolor = fc
	}

	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' " +
		"style='stroke: %s; fill: %s; stroke-width: 0.7' " +
		"width='%d' height='%d'>", scolor, fcolor, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i + 1, j, width, height)
			bx, by := corner(i, j, width, height)
			cx, cy := corner(i, j + 1, width, height)
			dx, dy := corner(i + 1, j + 1, width, height)
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(out, "</svg>")
}

func corner(i, j, width, height int) (float64, float64) {
	xyscale := float64(width) / 2 / xyrange // pixels per x or y unit
	zscale := float64(height) * 0.4        // pixels per z unit
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width) / 2 + (x - y) * cos30 * xyscale
	sy := float64(height) / 2 + (x + y) * sin30 * xyscale - z * zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
