package main

import (
	"fmt"
	"gonum.org/v1/plot/vg"
	// "gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
	// "gonum.org/v1/plot/vg/vgimg"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
	// "strings"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func randInt16(min int, max int) uint16 {
	return uint16(min + rand.Intn(max-min))
}
func randFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RegularPolygon(center vg.Point, radius vg.Length, ns int) vg.Path {
	path := vg.Path{}
	path.Move(center)
	// path.Arc(center, radius, 0, math.Pi*2)
	sl := segmentLength(ns)
	outer := vg.Point{
		X: center.X,
		Y: center.Y + radius,
	}
	for i := 0; i <= ns; i++ {
		new_outer := vg.Point{
			X: center.X + vg.Length(math.Cos(float64(i)*sl+math.Pi/2))*radius,
			Y: center.Y + vg.Length(math.Sin(float64(i)*sl+math.Pi/2))*radius,
		}

		path.Move(center)
		path.Line(outer)
		path.Move(outer)
		path.Line(new_outer)
		path.Move(center)
		outer = new_outer
	}

	return path
}

type Polar struct {
	Dist vg.Length
	Rad  float64
}

func (p Polar) ToXY(center vg.Point) vg.Point {
	return vg.Point{
		X: center.X + vg.Length(math.Cos(p.Rad))*p.Dist,
		Y: center.Y + vg.Length(math.Sin(p.Rad))*p.Dist,
	}
}

func Gemstone(center vg.Point, radius vg.Length, minpointradius vg.Length) vg.Path {
	// Generate the number of segments
	ns := randInt(3, 12)
	// How many randomly chosen points fit in a segment
	npps := randInt(3, 12)
	sl := segmentLength(ns)

	// Random points inside a segment
	randpolar := make([]Polar, npps)
	// Connections between points inside a segment
	connections := make([][]bool, npps)
	for i := range connections {
		connections[i] = make([]bool, npps)
	}
	// if a point should connect to itself in the next segment
	connect_to_next := make([]bool, npps)

	for i := range randpolar {
		// Generate random polar coordinates for points in the segment
		p := Polar{
			Dist: vg.Length(randFloat(float64(minpointradius), float64(radius))),
			Rad:  randFloat(0.0, sl),
		}
		randpolar[i] = p

		// Points have 1/npps chance of being connected
		for j := range connections {
			if randInt(1, npps) == 1 {
				connections[i][j] = true
			}
		}

		if randInt(0, npps/2) == 1 {
			connect_to_next[i] = true
		}

	}

	// fmt.Println(connections)
	xyps := make([][]vg.Point, ns)
	for i := range xyps {
		xyps[i] = make([]vg.Point, npps)

		// Convert to XY
		for j, p := range randpolar {
			np := p
			// If the segment number is even, then invert axis
			if i%2 == 1 && (ns%2 == 0 && i < ns) {
				np.Rad = (np.Rad * -1) + sl
			}
			np.Rad += sl*float64(i) + math.Pi/2
			xyps[i][j] = np.ToXY(center)
		}

	}

	path := RegularPolygon(center, radius, ns)
	path.Move(center)

	for i := range xyps {
		for j, point_conns := range connections {
			for k, isconn := range point_conns {
				if isconn {
					path.Move(xyps[i][j])
					path.Line(xyps[i][k])
					path.Move(center)
				}
			}
			if connect_to_next[j] {
				path.Move(xyps[(i+1)%ns][j])
				path.Line(xyps[i][j])
				path.Move(center)
			}
		}

		path.Move(center)
	}

	return path
}

// Divide a circle in n segments and return the length
func segmentLength(n int) float64 {
	return (math.Pi * 2.0) / float64(n)
}

func main() {
	var err error
	// TODO parse flags
	if len(os.Args) < 2 {
		log.Fatalf("need output filename")
	}
	file := os.Args[1]

	var seed int64
	if len(os.Args) < 3 {
		seed = time.Now().UnixNano()
	} else {
		seed, err = strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			log.Fatalf("error while parsing seed: %v", err)
		}

	}

	fmt.Println(seed)

	img_w := 20 * vg.Centimeter
	img_h := 20 * vg.Centimeter

	bg_rect := vg.Rectangle{Min: vg.Point{0, 0}, Max: vg.Point{img_w, img_h}}

	c := vgsvg.New(img_w, img_h)
	// d := draw.New(c)

	start_x := img_w / 2
	start_y := img_h / 2

	start_radius := img_w
	radius_step := 3 * vg.Centimeter

	start := vg.Point{start_x, start_y}
	rand.Seed(seed)
	c.SetLineWidth(5 * vg.Millimeter)
	c.Fill(bg_rect.Path())
	// for i := 0; i < 10; i++ {
	for r := start_radius; r > radius_step; r -= radius_step {
		c.SetColor(color.RGBA64{randInt16(0, 255), randInt16(0, 255), randInt16(0, 255), 255})
		poly := Gemstone(start, r, r-(radius_step*2))
		c.Stroke(poly)

	}

	// }

	// Save image
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not encode image: %+v", err)
	}
	return
}
