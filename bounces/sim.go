package bounces

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const maxBounce = 1000

func debug(f string, a ...interface{}) {
	return
	fmt.Printf(f+"\n", a...)
}

func (inp Input) simulate(r *rand.Rand) (float64, float64, bool) {
	debug("\n\nstart")
	x, y := inp.Start[0], inp.Start[1]
	vx, vy := inp.randInitial(r)

	//v0 := math.Sqrt(vx*vx + vy*vy)
	i := 0
	tdist := 0.0
	//defer func() {
	//	println(v0, i, tdist)
	//}()
	for ; i < maxBounce; i++ {
		debug("pos         %+.3f %+.3f %+.3f %+.3f", x, y, vx, vy)
		if inp.stopped(vx, vy) {
			return x, y, true
		}

		ox, oy := x, y
		dist, obstacle := inp.closesObstacle(x, y, vx, vy)
		debug(" - closest: %+.3f", dist)

		x, y, vx, vy = inp.advanceBall(dist, x, y, vx, vy)
		tdist += math.Sqrt((ox-x)*(ox-x) + (oy-y)*(oy-y))
		debug(" - advance: %+.3f %+.3f %+.3f %+.3f", x, y, vx, vy)
		if inp.stopped(vx, vy) {
			return x, y, true
		}

		vx, vy = obstacle.Bounce(x, y, vx, vy, inp.Ball, inp.Elasticity)
		debug(" - bounce:  %+.3f %+.3f %+.3f %+.3f", x, y, vx, vy)
	}

	inp.Error(errors.New("maxBounce reached - did you have a bad config?"))
	return 0, 0, false
}

func (inp Input) stopped(vx, vy float64) bool {
	return math.Sqrt(vx*vx+vy*vy) < inp.Terminal
}

func (inp Input) randInitial(r *rand.Rand) (float64, float64) {
	return inp.Velocity(r)
}

func (inp Input) closesObstacle(x0, y0, vx, vy float64) (float64, Obstacle) {
	closest, closestID := math.MaxFloat64, -1
	for i, o := range inp.Obstacles {
		d, ok := o.DistToColl(x0, y0, vx, vy, inp.Ball)
		if !ok {
			continue
		}
		if d >= 0 && d < closest {
			closest = d
			closestID = i
		}
	}
	if closestID < 0 {
		inp.Error(errors.New("did not collie with any obstacle"))
	}
	return closest * math.Sqrt(vx*vx+vy*vy), inp.Obstacles[closestID]
}

func (inp Input) advanceBall(dist, x, y, vx, vy float64) (float64, float64, float64, float64) {
	// Advance the ball accounting for fricition, either stoping before the obstacle or retaining some velocity
	v := math.Sqrt(vx*vx + vy*vy)
	n0, n1 := vx/v, vy/v
	distToTerminal := (v + inp.Terminal) / 2 * (inp.Terminal - v) / inp.Friction
	if distToTerminal > 0 && distToTerminal <= dist {
		return x + n0*distToTerminal, y + n1*distToTerminal, 0, 0
	}
	v = math.Sqrt(v*v + 2*dist*inp.Friction)
	return x + dist*n0, y + dist*n1, n0 * v, n1 * v
}
