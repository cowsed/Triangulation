package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/image/colornames"
)

/*
	Method:
	find closest point to hole on outer loop
	line from that point to closest point on loop
	go to other side and line from closest point on inside to out
	- splits main shape with hole into 2 polygons without holes
*/

func v2v(v mgl64.Vec2) pixel.Vec {
	return pixel.V(v.X(), v.Y())
}

func run() {
	fmt.Println("Intersect", intersect(mgl64.Vec2{-.5, -.5}, mgl64.Vec2{-.5, .5}, mgl64.Vec2{0, 0}, mgl64.Vec2{1, 0}))

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	numPts := 15

	points := circlePoints(numPts)
	activePt := 0

	edges, diags, ps := Triangulate(points)
	imd := imdraw.New(nil)
	redraw(activePt, ps, edges, diags, imd)

	s := 3.0
	fmt.Println("Beginning")
	for !win.Closed() {

		//Handle input
		{
			if win.JustPressed(pixelgl.KeyTab) {
				if win.Pressed(pixelgl.KeyLeftShift) {
					activePt--
					if activePt < 0 {
						activePt += len(points)
					}
				} else {
					activePt++
				}
				activePt = activePt % len(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.Pressed(pixelgl.KeyLeft) {
				points[activePt] = points[activePt].Add(mgl64.Vec2{-s, 0})
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.Pressed(pixelgl.KeyRight) {
				points[activePt] = points[activePt].Add(mgl64.Vec2{s, 0})
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.Pressed(pixelgl.KeyDown) {
				points[activePt] = points[activePt].Add(mgl64.Vec2{0, -s})
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.Pressed(pixelgl.KeyUp) {
				points[activePt] = points[activePt].Add(mgl64.Vec2{0, s})
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.JustPressed(pixelgl.KeyEqual) {
				numPts++
				activePt = activePt % numPts
				points = circlePoints(numPts)
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
			if win.JustPressed(pixelgl.KeyMinus) {
				numPts--
				activePt = activePt % numPts
				points = circlePoints(numPts)
				edges, diags, ps = Triangulate(points)
				redraw(activePt, ps, edges, diags, imd)
			}
		}
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}
func circlePoints(n int) []mgl64.Vec2 {
	points := make([]mgl64.Vec2, n)
	r := 200.0
	for i := range points {
		ang := (360.0 / float64(n)) * float64(i)
		x := math.Cos(mgl64.DegToRad(ang)) * r
		y := math.Sin(mgl64.DegToRad(ang)) * r

		points[i] = mgl64.Vec2{400 + x, 300 + y}

	}
	return points
}

func redraw(activePoint int, ps []mgl64.Vec2, edges, diags [][2]int, imd *imdraw.IMDraw) {
	//Draw Points
	imd.Clear()
	imd.Color = colornames.Red
	for _, e := range edges {
		imd.Push(v2v(ps[e[0]]))
		imd.Push(v2v(ps[e[1]]))
		imd.Line(4)
	}
	imd.Color = colornames.Blue
	for _, d := range diags {
		imd.Push(v2v(ps[d[0]]))
		imd.Push(v2v(ps[d[1]]))
		imd.Line(2)
	}
	imd.Color = colornames.Black

	for _, p := range ps {
		imd.Push(v2v(p))
		imd.Circle(8, 0)
	}
	imd.Color = colornames.Lightgreen

	imd.Push(v2v(ps[activePoint]))
	imd.Circle(8, 0)
}
func main() {
	pixelgl.Run(run)
}
