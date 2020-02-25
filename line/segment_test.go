package line

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	distance(t, 1, true, 0,
		0, 1, 1, 1,
		0, 0, 0, 1)
	distance(t, 0.5, true, 0,
		0, 1, 1, 0,
		0, 0, 1, 1)
	distance(t, -1, false, 0,
		0, 1, 1, 1,
		0, 0, 1, -0.1)
	distance(t, 0.9, true, 0.1,
		0, 1, 1, 1,
		0, 0, 0, 1)
	distance(t, math.Sqrt2/2-0.1, true, 0.1,
		0, 1, 1, 0,
		0, 0, 1/math.Sqrt2, 1/math.Sqrt2)
	distance(t, -1, false, 0.1,
		0, 1, 1, 1,
		0, 0, 1, -0.1)
}

func TestBounce(t *testing.T) {
	bounce(t, 1, 0.1,
		0, 0, 0, 1,
		1, 0, -1, 1,
		1, 1)
	bounce(t, 1, 0.1,
		0, 0, 0, 1,
		-1, 0, 1, 1,
		-1, 1)
	bounce(t, 1, 0.1,
		0, 1, 0, 0,
		-1, 0, 1, 1,
		-1, 1)
	bounce(t, 1, 0.1,
		0, 1, 0, 0,
		1, 0, -1, 1,
		1, 1)
}

func distance(t *testing.T, d float64, flag bool, r, px, py, qx, qy, x, y, vx, vy float64) {
	s := SegmentFromPoints(px, py, qx, qy)

	a, b := s.DistToColl(x, y, vx, vy, r)

	t.Log("expected:", flag, d, "got: ", b, a)
	if b != flag || math.Abs(a-d) > tol {
		t.Error()
	}
}

func bounce(t *testing.T, e, r float64, px, py, qx, qy, x, y, vx, vy, bx, by float64) {
	s := SegmentFromPoints(px, py, qx, qy)

	a, b := s.Bounce(x, y, vx, vy, e, r)

	t.Log("expected:", bx, by, "got: ", a, b)
	if math.Abs(a-bx) > tol || math.Abs(b-by) > tol {
		t.Error()
	}
}
