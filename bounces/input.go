package bounces

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/vron/bounces/line"
	"github.com/vron/bounces/shape"
)

type VelocitySampler func(r *rand.Rand) (float64, float64)

type Obstacle interface {
	Bounds() (x, y, w, h float64)
	DistToColl(x, y, vx, vy, r float64) (float64, bool)
	Bounce(x, y, vx, vy, r, el float64) (float64, float64)
}

type Input struct {
	Ball       float64
	Start      [2]float64
	Velocity   VelocitySampler
	Friction   float64
	Terminal   float64
	Elasticity float64
	Measures   []Rect

	ImageRes  int
	ImagePath string
	Error     func(e error, a ...interface{})

	Obstacles []Obstacle
}

type Rect struct {
	X, Y, W, H float64
	Name       string
}

type Parser interface {
	Handle([]string, *Input) bool
}

func ParseInput(r io.Reader, fatal func(e error, a ...interface{})) Input {
	buf, err := ioutil.ReadAll(r)
	fatal(err, "error reading input:")

	inp := Input{Error: fatal}
	lines := bytes.Split(buf, []byte("\n"))
	parsers := []Parser{
		parseBall{},
		parseLine{},
		parseCircle{},
		parseStart{},
		parseVelocity{},
		parseFriction{},
		parseTerminal{},
		parseElasticity{},
		parseMeasure{},
	}
Lines:
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		parts := strings.Split(sanitize(string(l)), " ")
		if len(parts) < 1 {
			continue
		}
		for _, p := range parsers {
			if p.Handle(parts, &inp) {
				continue Lines
			}
		}
		fatal(errors.New("unknown command: '" + string(l) + "'"))
	}

	if inp.Measures == nil || len(inp.Measures) < 1 {
		fatal(errors.New("must have a measure for convergence"))
	}

	return inp
}

func sanitize(s string) string {
	r := regexp.MustCompile(`\s+`)
	sn := r.ReplaceAllString(s, " ")
	r = regexp.MustCompile(` ?//.*`)
	sn = r.ReplaceAllString(sn, "")
	// fmt.Printf("'%v' -> '%v'\n", s, sn)
	return sn
}

type parseStart struct{}

func (p parseStart) Handle(d []string, i *Input) bool {
	if d[0] != "start" {
		return false
	}
	if len(d) != 3 {
		i.Error(errors.New(""), "start expects 2 arguments")
	}
	i.Start[0] = i.num(d[1])
	i.Start[1] = i.num(d[2])
	return true
}

type parseVelocity struct{}

func (p parseVelocity) Handle(d []string, i *Input) bool {
	if d[0] != "velocity" {
		return false
	}
	if len(d) < 4 {
		i.Error(errors.New(""), "velocity expects 3 or 4 arguments")
	}
	switch d[1] {
	case "uniform":
		l, h := i.num(d[2]), i.num(d[3])
		i.Velocity = func(r *rand.Rand) (float64, float64) {
			theta := r.Float64() * math.Pi * 2
			v := r.Float64()*(h-l) + l
			return math.Cos(theta) * v, math.Sin(theta) * v
		}
	case "lognormal":
		mu, sig, lim := i.num(d[2]), i.num(d[3]), i.num(d[4])
		i.Velocity = func(r *rand.Rand) (float64, float64) {
			theta := r.Float64() * math.Pi * 2
			v := 1e99
			for v > lim {
				v = math.Exp(r.NormFloat64()*sig + mu)
			}
			return math.Cos(theta) * v, math.Sin(theta) * v
		}
	default:
		i.Error(errors.New(""), "velocity did not recognize the distr.")
	}
	return true
}

type parseBall struct{}

func (p parseBall) Handle(d []string, i *Input) bool {
	if d[0] != "ball" {
		return false
	}
	if len(d) != 2 {
		i.Error(errors.New(""), "ball expects 1 argument")
	}
	i.Ball = i.num(d[1])
	return true
}

type parseLine struct{}

func (p parseLine) Handle(d []string, i *Input) bool {
	if d[0] != "line" {
		return false
	}
	if len(d) != 5 {
		i.Error(errors.New(""), "line expects 4 arguments")
	}
	i.Obstacles = append(i.Obstacles, line.SegmentFromPoints(i.num(d[1]), i.num(d[2]), i.num(d[3]), i.num(d[4])))
	i.Obstacles = append(i.Obstacles, shape.NewCircle(i.num(d[3]), i.num(d[4]), 0))
	i.Obstacles = append(i.Obstacles, shape.NewCircle(i.num(d[1]), i.num(d[2]), 0))
	return true
}

type parseCircle struct{}

func (p parseCircle) Handle(d []string, i *Input) bool {
	if d[0] != "circle" {
		return false
	}
	if len(d) != 4 {
		i.Error(errors.New(""), "circle expects 3 arguments")
	}
	i.Obstacles = append(i.Obstacles, shape.NewCircle(i.num(d[1]), i.num(d[2]), i.num(d[3])))
	return true
}

type parseFriction struct{}

func (p parseFriction) Handle(d []string, i *Input) bool {
	if d[0] != "friction" {
		return false
	}
	if len(d) != 2 {
		i.Error(errors.New(""), "friction expects 1 arguments")
	}
	i.Friction = i.num(d[1])
	return true
}

type parseTerminal struct{}

func (p parseTerminal) Handle(d []string, i *Input) bool {
	if d[0] != "terminal" {
		return false
	}
	if len(d) != 2 {
		i.Error(errors.New(""), "terminal expects 1 arguments")
	}
	i.Terminal = i.num(d[1])
	return true
}

type parseElasticity struct{}

func (p parseElasticity) Handle(d []string, i *Input) bool {
	if d[0] != "elasticity" {
		return false
	}
	if len(d) != 2 {
		i.Error(errors.New(""), "elasticity expects 1 arguments")
	}
	i.Elasticity = i.num(d[1])
	return true
}

type parseMeasure struct{}

func (p parseMeasure) Handle(d []string, i *Input) bool {
	if d[0] != "measure" {
		return false
	}
	if len(d) != 6 {
		i.Error(errors.New(""), "meassure expects 5 arguments")
	}
	i.Measures = append(i.Measures, Rect{i.num(d[2]), i.num(d[3]), i.num(d[4]) - i.num(d[2]), i.num(d[5]) - i.num(d[3]), d[1]})
	return true
}

func (i *Input) num(b string) float64 {
	f, err := strconv.ParseFloat(b, 64)
	i.Error(err, "parsing i.number")
	return f
}
