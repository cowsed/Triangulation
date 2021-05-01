package main

import "github.com/go-gl/mathgl/mgl64"

func ccw(A, B, C mgl64.Vec2) bool {
	return (C.Y()-A.Y())*(B.X()-A.X()) > (B.Y()-A.Y())*(C.X()-A.X())
}

//Return true if line segments AB and CD intersect
func intersect(A, B, C, D mgl64.Vec2) bool {
	return ccw(A, C, D) != ccw(B, C, D) && ccw(A, B, C) != ccw(A, B, D)
}
