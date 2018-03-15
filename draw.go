package main

import (
	"errors"
	// "fmt"
	"math"
)

func (image Image) DrawLines(edges *Matrix, c Color) {
	m := edges.mat
	for i := 0; i < edges.cols-1; i += 2 {
		image.DrawLine(c, int(m[0][i]), int(m[1][i]), int(m[0][i+1]), int(m[1][i+1]))
	}
}

func (image Image) DrawLine(c Color, x0, y0, x1, y1 int) error {
	if x0 < 0 || y0 < 0 || x1 > image.width || y1 > image.height {
		return errors.New("Error: Coordinates out of bounds")
	}
	if x0 > x1 {
		x1, x0 = x0, x1
		y1, y0 = y0, y1
	}

	deltaX := x1 - x0
	deltaY := y1 - y0
	if deltaY >= 0 {
		if math.Abs(float64(deltaY)) <= math.Abs(float64(deltaX)) {
			image.drawLineOctant1(c, deltaY, deltaX*-1, x0, y0, x1, y1)
		} else {
			image.drawLineOctant2(c, deltaY, deltaX*-1, x0, y0, x1, y1)
		}
	} else {
		if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
			image.drawLineOctant7(c, deltaY, deltaX*-1, x0, y0, x1, y1)
		} else {
			image.drawLineOctant8(c, deltaY, deltaX*-1, x0, y0, x1, y1)
		}
	}
	return nil
}

func (image Image) drawLineOctant1(c Color, lA, lB, x0, y0, x1, y1 int) error {
	y := y0
	lD := 2*lA + lB
	for x := x0; x < x1; x++ {
		err := image.plot(c, x, y)
		if err != nil {
			return err
		}
		if lD > 0 {
			y++
			lD += 2 * lB
		}
		lD += 2 * lA
	}
	return nil
}

func (image Image) drawLineOctant2(c Color, lA, lB, x0, y0, x1, y1 int) error {
	x := x0
	lD := lA + 2*lB
	for y := y0; y < y1; y++ {
		err := image.plot(c, x, y)
		if err != nil {
			return err
		}
		if lD < 0 {
			x++
			lD += 2 * lA
		}
		lD += 2 * lB
	}
	return nil
}

func (image Image) drawLineOctant7(c Color, lA, lB, x0, y0, x1, y1 int) error {
	x := x0
	lD := lA - 2*lB
	for y := y0; y > y1; y-- {
		err := image.plot(c, x, y)
		if err != nil {
			return err
		}
		if lD > 0 {
			x++
			lD += 2 * lA
		}
		lD -= 2 * lB
	}
	return nil
}

func (image Image) drawLineOctant8(c Color, lA, lB, x0, y0, x1, y1 int) error {
	y := y0
	lD := 2*lA - lB
	for x := x0; x < x1; x++ {
		err := image.plot(c, x, y)
		if err != nil {
			return err
		}
		if lD < 0 {
			y--
			lD -= 2 * lB
		}
		lD += 2 * lA
	}
	return nil
}

func (m *Matrix) AddCircle(cx, cy, cz, r float64) {
	var oldX, oldY float64 = -1, -1
	// TODO: No magic numbers wow i have so much to fix
	for i := 0; i <= 100; i++ {
		var t float64 = float64(i) / float64(100)
		x := r*math.Cos(2*math.Pi*t) + cx
		y := r*math.Sin(2*math.Pi*t) + cy
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, cz, x, y, cz)
		oldX = x
		oldY = y
	}
}

func makeHermiteCoefs(p0, p1, rp0, rp1 float64) (*Matrix, error) {
	h := MakeMatrix(4, 0)
	h.AddCol([]float64{2, -3, 0, 1})
	h.AddCol([]float64{-2, 3, 0, 0})
	h.AddCol([]float64{1, -2, 1, 0})
	h.AddCol([]float64{1, -1, 0, 0})

	mat := MakeMatrix(4, 0)
	mat.AddCol([]float64{p0, p1, rp0, rp1})

	return mat.Mult(h)
}

func (m *Matrix) AddHermite(x0, y0, x1, y1, rx0, ry0, rx1, ry1, stepSize float64) error {
	xC, err := makeHermiteCoefs(x0, x1, rx0, rx1)
	if err != nil {
		return err
	}
	yC, err := makeHermiteCoefs(y0, y1, ry0, ry1)
	if err != nil {
		return err
	}
	// TODO: Figure out a better way to do this
	var oldX, oldY float64 = -1, -1
	var steps int = int(1 / stepSize)
	for i := 0; i <= steps; i++ {
		var t float64 = float64(i) / float64(steps)
		x := xC.mat[0][0]*math.Pow(t, 3.0) + xC.mat[1][0]*math.Pow(t, 2.0) + xC.mat[2][0]*t + xC.mat[3][0]
		y := yC.mat[0][0]*math.Pow(t, 3.0) + yC.mat[1][0]*math.Pow(t, 2.0) + yC.mat[2][0]*t + yC.mat[3][0]
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, 0.0, x, y, 0.0)
		oldX = x
		oldY = y
	}
	return nil
}

func makeBezierCoefs(p0, p1, p2, p3 float64) (*Matrix, error) {
	h := MakeMatrix(4, 0)
	h.AddCol([]float64{-1, 3, -3, 1})
	h.AddCol([]float64{3, -6, 3, 0})
	h.AddCol([]float64{-3, 3, 0, 0})
	h.AddCol([]float64{1, 0, 0, 0})

	mat := MakeMatrix(4, 0)
	mat.AddCol([]float64{p0, p1, p2, p3})

	return mat.Mult(h)
}

// TODO: maybe combine with hermite bc a lot of duplicate code?
// Or make a separate parametric fxn
func (m *Matrix) AddBezier(x0, y0, x1, y1, x2, y2, x3, y3, stepSize float64) error {
	xC, err := makeBezierCoefs(x0, x1, x2, x3)
	if err != nil {
		return err
	}
	yC, err := makeBezierCoefs(y0, y1, y2, y3)
	if err != nil {
		return err
	}
	// TODO: Figure out a better way to do this
	var oldX, oldY float64 = -1, -1
	var steps int = int(1 / stepSize)
	for i := 0; i <= steps; i++ {
		var t float64 = float64(i) / float64(steps)
		x := xC.mat[0][0]*math.Pow(t, 3.0) + xC.mat[1][0]*math.Pow(t, 2.0) + xC.mat[2][0]*t + xC.mat[3][0]
		y := yC.mat[0][0]*math.Pow(t, 3.0) + yC.mat[1][0]*math.Pow(t, 2.0) + yC.mat[2][0]*t + yC.mat[3][0]
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, 0.0, x, y, 0.0)
		oldX = x
		oldY = y
	}
	return nil
}
