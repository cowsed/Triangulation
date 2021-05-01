package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
)

func sign(p1, p2, p3 mgl64.Vec2) float64 {
	return (p1.X()-p3.X())*(p2.Y()-p3.Y()) - (p2.X()-p3.X())*(p1.Y()-p3.Y())
}

func inTri(a, b, c, p mgl64.Vec2) bool {

	var has_neg, has_pos bool

	d1 := sign(p, a, b)
	d2 := sign(p, b, c)
	d3 := sign(p, c, a)

	has_neg = (d1 < 0) || (d2 < 0) || (d3 < 0)
	has_pos = (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(has_neg && has_pos)

}

func findPrevEnabled(removed []bool, start int) int {
	for i := start - 1; i != start; i-- {
		if i < 0 {
			i += len(removed)
		}
		if !removed[i] {
			return i
		}
	}
	return -1
}
func findNextEnabled(removed []bool, start int) int {
	for i := start + 1; i != start; i++ {
		i = i % len(removed)
		if !removed[i] {
			return i
		}
	}
	return -1
}

//Center is ear tip
func isEar(points []mgl64.Vec2, removed []bool, center int) bool {
	prev := findPrevEnabled(removed, center)
	next := findNextEnabled(removed, center)
	a, b, c := points[prev], points[center], points[next]
	//fmt.Printf("Checking inside A:%v@%d, B:%v@%d, C:%v@%d\n", a, prev, b, center, c, next)

	for i := 0; i < len(removed); i++ {
		//Dont Check self points or if point is gone
		if i == center || i == prev || i == next || removed[i] {
			continue
		}
		if inTri(a, b, c, points[i]) {
			//fmt.Printf("Point inside %v@%d\n", points[i], i)
			return false
		}

	}

	return true
}

func findEar(points []mgl64.Vec2, removed []bool) int {
	for i, r := range removed {
		if !r {
			if isEar(points, removed, i) {
				return i
			}
		}
	}
	return -1
}
func countPoints(b []bool) int {
	sum := 0
	for _, v := range b {
		if !v {
			sum++
		}
	}
	return sum
}

//points are just a list of points that are connected to adjacent points
func Triangulate(points []mgl64.Vec2) ([][2]int, [][2]int, []mgl64.Vec2) {

	//Duplicate points to avoid messing stuff up
	fmt.Println("Triangulating")
	L := points[:]

	//Removed to keep track of ears and such without reslicing
	removed := make([]bool, len(L))

	//Generate Border Edges
	edges := make([][2]int, len(points))
	for i := 0; i < len(points); i++ {
		a, b := i, i+1
		if b == len(points) {
			b = 0
		}
		edges = append(edges, [2]int{a, b})
	}
	maxDiags := len(points) - 3
	diags := make([][2]int, 0, maxDiags)

	for countPoints(removed) > 3 {
		//Find Ear
		earInd := findEar(L, removed)
		if earInd == -1 {
			fmt.Println("No ear found")
			break
		}
		//Get Side Point
		//a := (earInd - 1)
		//if a == -1 {
		//	a = len(L) - 1
		//}
		a := findPrevEnabled(removed, earInd)
		//Get Other Side Point
		c := findNextEnabled(removed, earInd)
		if !IntersectAll(a, c, points, edges) {
			//	fmt.Println("No intersect")
			diags = append(diags, [2]int{a, c})
			//Remove Vert from list
		}
		fmt.Println("removing", earInd)
		removed[earInd] = true

	}
	return edges, diags, L
}
func IntersectAll(a, c int, points []mgl64.Vec2, edges [][2]int) bool {
	for _, e := range edges {
		if intersect(points[a], points[c], points[e[0]], points[e[1]]) {
			return false
		}
	}
	return true
}
