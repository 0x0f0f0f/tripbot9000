package trip

import (
	// "fmt"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
)

// Polar coordinates
type Polar struct {
	Dist vg.Length
	Rad  float64
}

// Convert polar coordiniates to cartesian, relative to an origin point
func (p Polar) ToXY(origin vg.Point) vg.Point {
	return vg.Point{
		X: origin.X + vg.Length(math.Cos(p.Rad))*p.Dist,
		Y: origin.Y + vg.Length(math.Sin(p.Rad))*p.Dist,
	}
}

// Divide a circle in n segments and return the length
func segmentLength(n int) float64 {
	return (math.Pi * 2.0) / float64(n)
}

func RegularPolygon(center vg.Point, radius vg.Length, ns int, centerlines bool) vg.Path {
	path := vg.Path{}
	path.Move(center)
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
		if centerlines {
			path.Line(outer)
		}
		path.Move(outer)
		path.Line(new_outer)
		path.Move(center)
		outer = new_outer
	}

	return path
}

func SierpinskyTriangle(center vg.Point, radius vg.Length, maxdepth int, path vg.Path, flipY bool) vg.Path {
	yScale := vg.Length(1)
	if flipY {
		yScale = vg.Length(-1)
	}

	top := vg.Point{center.X, yScale*center.Y + radius}

	botleft := vg.Point{
		X: center.X + radius*vg.Length(math.Cos(2*math.Pi/3+math.Pi/2)),
		Y: yScale*center.Y + radius*vg.Length(math.Sin(2*math.Pi/3+math.Pi/2)),
	}
	botright := vg.Point{
		X: center.X + radius*vg.Length(math.Cos(math.Pi/3.0-(math.Pi/2))),
		Y: yScale*center.Y + radius*vg.Length(math.Sin(math.Pi/3.0-(math.Pi/2))),
	}

	// fmt.Printf("%+v\n%+v\n%+v\n", top, botleft, botright)

	// Calculate subtriangle radius and center
	base := botright.X - botleft.X
	// height := top.Y - botright.Y
	// Draw the outer triangle
	path.Move(top)
	path.Line(botleft)
	path.Line(botright)
	path.Line(top)
	path.Move(center)

	if maxdepth > 0 {
		subradius := radius / 2
		top_center := vg.Point{center.X, center.Y + subradius}
		left_center := vg.Point{center.X - base/4, center.Y - subradius/2}
		right_center := vg.Point{center.X + base/4, center.Y - subradius/2}

		path = SierpinskyTriangle(top_center, subradius, maxdepth-1, path, flipY)
		path = SierpinskyTriangle(left_center, subradius, maxdepth-1, path, flipY)
		path = SierpinskyTriangle(right_center, subradius, maxdepth-1, path, flipY)
	}
	return path
}

// Generate a Random gemstone, given the center position, the radius, the minimum distance of gemstone points
// from the radius and the number of segments of the starting circle.
func Gemstone(center vg.Point, radius vg.Length, minpointradius vg.Length, ns int) vg.Path {
	// How many Randomly chosen points fit in a segment
	npps := RandInt(3, 12)
	sl := segmentLength(ns)

	// Random points inside a segment
	Randpolar := make([]Polar, npps)
	// Connections between points inside a segment
	connections := make([][]bool, npps)
	for i := range connections {
		connections[i] = make([]bool, npps)
	}
	// if a point should connect to itself in the next segment
	connect_to_next := make([]bool, npps)

	for i := range Randpolar {
		// Generate Random polar coordinates for points in the segment
		p := Polar{
			Dist: vg.Length(RandFloat(float64(minpointradius), float64(radius))),
			Rad:  RandFloat(0.0, sl),
		}
		Randpolar[i] = p

		// Points have 1/npps Chance of being connected
		for j := range connections {
			if RandInt(1, npps) == 1 {
				connections[i][j] = true
			}
		}

		if RandInt(0, npps/2) == 1 {
			connect_to_next[i] = true
		}

	}

	// fmt.Println(connections)
	xyps := make([][]vg.Point, ns)
	for i := range xyps {
		xyps[i] = make([]vg.Point, npps)

		// Convert to XY
		for j, p := range Randpolar {
			np := p
			// If the segment number is even, then invert axis
			if i%2 == 1 && (ns%2 == 0 && i < ns) {
				np.Rad = (np.Rad * -1) + sl
			}
			np.Rad += sl*float64(i) + math.Pi/2
			xyps[i][j] = np.ToXY(center)
		}

	}

	path := vg.Path{}
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

// Generate a Random concentric geometric pattern and draw it on a canvas
func DrawRandomGeom(
	c vg.Canvas,
	center vg.Point,
	outer_radius,
	radius_step vg.Length,
	palette []color.Color) vg.Canvas {
	for r := outer_radius; r > radius_step; r -= radius_step {
		c.SetColor(palette[RandInt(0, len(palette))])
		// Generate the number of segments
		ns := RandInt(3, 12)
		// 1/3 Chance of making a polygon
		n := RandInt(0, 7)
		switch {
		case n == 0:
			triangle := SierpinskyTriangle(center, r, RandInt(3, 10), vg.Path{}, Chance(2))
			c.Stroke(triangle)
		case n < 4:
			poly := RegularPolygon(center, r, ns, Chance(2))
			c.Stroke(poly)
		case n >= 4:
			if Chance(2) {
				poly := RegularPolygon(center, r, ns, Chance(2))
				c.Stroke(poly)
			}
			gem := Gemstone(center, r, r-(radius_step*2), ns)
			c.Stroke(gem)

		}
	}
	return c
}
