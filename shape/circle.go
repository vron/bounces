package shape

import "math"

const tol = 1e-8

// A Segment represents a line segment as (x0, y0, r)
type Circle struct {
	X, Y, R float64
}

func NewCircle(x, y, r float64) Circle {
	return Circle{X: x, Y: y, R: r}
}

func (s Circle) Bounds() (x, y, w, h float64) {
	return s.X - s.R, s.Y - s.R, 2 * s.R, 2 * s.R
}

func (s Circle) DistToColl(x, y, vx, vy, r float64) (float64, bool) {
	// https://stackoverflow.com/questions/1073336/circle-line-segment-collision-detection-algorithm
	a := vx*vx + vy*vy
	fx, fy := x-s.X, y-s.Y
	b := 2 * (fx*vx + fy*vy)
	c := fx*fx + fy*fy - (r+s.R)*(r+s.R)

	d := b*b - 4*a*c
	if d <= 0 {
		return -1, false // no intersection
	}
	d = math.Sqrt(d)
	t1 := (-b - d) / (2 * a)
	t2 := (-b + d) / (2 * a)
	_ = t2
	if t1 > 0 {
		return t1, true
	}
	return -1, false
}

func (s Circle) Bounce(x, y, vx, vy, r, el float64) (float64, float64) {
	px, py := x, y

	// split velocity into (t, n) components, t velocity must be
	// unchanged thanks to conservation of momentum
	nx, ny := (px-s.X)/(r+s.R), (py-s.Y)/(r+s.R)
	tx, ty := ny, -nx

	n, t := nx*vx+ny*vy, tx*vx+ty*vy

	E := (n * n) * el
	nn := math.Sqrt(E)
	if n > 0 {
		nn = -nn
	}

	return nx*nn + tx*t, ny*nn + ty*t
}
