package shape

import (
	"math"
	"testing"
)

const n = 1 / math.Sqrt2

func TestDistance(t *testing.T) {
	distance(t, math.Sqrt2-0.2, true,
		1, 1, 0.1,
		0.1, 0, 0, n, n)
	distance(t, -1, false,
		1, 1, 0.1,
		0.1, 0, 0, 0, n)
	distance(t, -1, false,
		1, 1, 0.1,
		0.1, 0, 0, -n, -n)
	distance(t, -1, false,
		1, 1, 0.1,
		0.1, 0.95, 0.95, n, n)
}

func TestBounce(t *testing.T) {
	bounce(t, 1,
		1, 1, 0.1, 0.1,
		0, 0, n, n,
		-n, -n)
	bounce(t, 1,
		1, 1, 0.1, 0.1,
		0, 0, 4*n, 4*n,
		-4*n, -4*n)
	bounce(t, 0,
		1, 1, 0.1, 0.1,
		0, 0, n, n,
		0, 0)
	bounce(t, 0.5,
		1, 0, 0.1, 0.1,
		0, 0, 1, 0,
		-1/math.Sqrt2, 0)
	bounce(t, 1,
		1, 1.2, 0.1, 0.1,
		0, 0, n, n,
		n, -n)
	bounce(t, 0.5,
		1, 1.2, 0.1, 0.1,
		0, 0, n, n,
		n, -0.5)
	bounce(t, 0.5,
		1, 1.2, 0.1, 0.1,
		0, 0, 4*n, 4*n,
		4*n, -4*0.5)
}

func distance(t *testing.T, d float64, flag bool, px, py, r, R, x, y, vx, vy float64) {
	s := NewCircle(px, py, r)

	a, b := s.DistToColl(x, y, vx, vy, R)

	t.Log("expected:", flag, d, "got: ", b, a)
	if b != flag || math.Abs(a-d) > tol {
		t.Error()
	}
}

func bounce(t *testing.T, e float64, px, py, r, R, x, y, vx, vy, bx, by float64) {
	s := NewCircle(px, py, r)

	a, b := s.Bounce(x, y, vx, vy, e, R)

	t.Log("expected:", bx, by, "got: ", a, b)
	if math.Abs(a-bx) > tol || math.Abs(b-by) > tol {
		t.Error()
	}
}
