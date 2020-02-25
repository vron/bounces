package line

import (
	"math"
)

// TODO: update this to use the offsets based on ball radius

const tol = 1e-8
const R = 0.001

// A Segment represents a line segment as (x0, y0, dx, dy)
type Segment [4]float64

func SegmentFromPoints(x0, y0, x1, y1 float64) Segment {
	return Segment{x0, y0, x1 - x0, y1 - y0}
}

func (s Segment) Bounds() (x, y, w, h float64) {
	if s[2] < 0 {
		x, w = s[0]+s[2], -s[2]
	} else {
		x, w = s[0], s[2]
	}
	if s[3] < 0 {
		y, h = s[1]+s[3], -s[3]
	} else {
		y, h = s[1], s[3]
	}
	return
}

func (s Segment) DistToColl(x, y, vx, vy, r float64) (float64, bool) {
	// we check this by checking to lintes offset (the one in correct dir), and using two circles at the ends.
	d, n := s.distAndItem(x, y, vx, vy, r)
	if n < 0 {
		return -1, false
	}
	return d, true
}

func (s Segment) distAndItem(x, y, vx, vy, r float64) (float64, int) {
	l := math.Sqrt(s[2]*s[2] + s[3]*s[3])
	nx, ny := -s[3]/l, s[2]/l
	if nx*vx+ny*vy > 0 {
		nx, ny = -nx, -ny
	}
	d, ok := s.distToColl(x, y, vx, vy, nx*r, ny*r)
	if ok {
		return d, 0
	}
	return -1, -1

}

func (s Segment) distToColl(x, y, vx, vy, ox, oy float64) (float64, bool) {
	// see also http://geomalgorithms.com/a05-_intersect-1.html
	u1, u2 := s[2], s[3]
	v1, v2 := vx, vy
	w1, w2 := ox+s[0]-x, oy+s[1]-y

	d := (v1*u2 - v2*u1)
	if math.Abs(d) < tol {
		return -1, false
	}

	si := (v2*w1 - v1*w2) / d
	if si < 0 || si > 1 {
		return -1, false
	}

	ti := (u2*w1 - u1*w2) / d
	return ti, true
}

func (s Segment) Bounce(x, y, vx, vy, r, el float64) (float64, float64) {
	// so we bounce against the line
	// split velocity into (t, n) components, t velocity must be
	// unchanged thanks to conservation of momentum
	l := math.Sqrt(s[2]*s[2] + s[3]*s[3])
	nx, ny := -s[3]/l, s[2]/l
	tx, ty := s[2]/l, s[3]/l

	n, t := nx*vx+ny*vy, tx*vx+ty*vy
	E := (n * n) * el
	nn := math.Sqrt(E)
	if n > 0 {
		nn = -nn
	}

	//print(math.Sqrt(vx*vx + vy*vy))
	vx = nx*nn + tx*t
	vy = ny*nn + ty*t
	//println(" ", math.Sqrt(vx*vx+vy*vy))

	return vx, vy

}
