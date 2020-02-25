package bounces

import (
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/gonum/stat"
)

func Run(inp Input, prec float64, min, max int64) Results {

	// run a set of simulations, calculating both the measures ans the distribution image.
	// we always want 5 independent simulations for the measures such that
	// we can calculate the standard errors. The distribution image we simply
	// accumulate all to smooth as much as possible

	nos := 5
	var blockSize = int(10000)
	if int(min) < blockSize {
		blockSize = int(min)
	}

	r := rand.New(rand.NewSource(0))
	res := allocateResults(nos, inp)
	queue := make(chan bool, runtime.NumCPU())
	wg := sync.WaitGroup{}

	for targetSamples := int64(blockSize); true; targetSamples *= 2 {

		for _, sample := range res {
			notr := targetSamples - sample.no
			notr = notr / int64(blockSize)
			for i := int64(0); i < notr; i++ {

				// start workers in blocks of blockSize until we have sufficient ones
				queue <- true
				wg.Add(1)
				go func(res *result, r *rand.Rand) {
					for i := 0; i < blockSize; i++ {
						x, y, _ := inp.simulate(r)
						res.addPos(x, y)
					}
					<-queue
					wg.Done()
				}(sample, rand.New(rand.NewSource(r.Int63())))
			}
		}
		wg.Wait()

		// now each of the results should have targetSamples number of samples, check if the total is big enoguh
		// or the precision is good enough so we can quit.
		if int64(nos)*targetSamples < min {
			continue
		}

		if int64(nos)*targetSamples >= max {
			break
		}
		_, maxErr := calculateActualResults(res)
		if maxErr <= prec {
			break
		}
	}

	rr, _ := calculateActualResults(res)
	return rr
}

type Results struct {
	Measures      []float64
	MeasureErrors []float64
	Image         []int32
}

func calculateActualResults(results []*result) (Results, float64) {
	r := Results{
		Image: results[0].image.image,
	}

	maxErr := 0.0
	for mi := range results[0].measures {
		vals := make([]float64, 0, len(results))
		for _, r := range results {
			vals = append(vals, float64(r.measures[mi])/float64(r.no))
		}
		r.Measures = append(r.Measures, stat.Mean(vals, nil))
		v := stat.StdErr(stat.StdDev(vals, nil), float64(len(results)))
		r.MeasureErrors = append(r.MeasureErrors, v)
		if v > maxErr {
			maxErr = v
		}
	}
	return r, maxErr
}

func allocateResults(nos int, inp Input) []*result {
	img := &image{
		image: make([]int32, inp.ImageRes*inp.ImageRes),
		res:   inp.ImageRes,
	}
	img.bounds[0], img.bounds[1], img.bounds[2], img.bounds[3] = inp.bounds()
	results := make([]*result, 0, nos)
	for i := 0; i < nos; i++ {
		res := &result{
			input: &inp,
			image: img,
		}

		for range inp.Measures {
			res.measures = append(res.measures, 0)
		}
		results = append(results, res)
	}
	return results
}

func (inp Input) bounds() (x, y, w, h float64) {
	xm, xM := math.MaxFloat64, -math.MaxFloat64
	ym, yM := math.MaxFloat64, -math.MaxFloat64
	for _, o := range inp.Obstacles {
		x, y, w, h := o.Bounds()
		if x < xm {
			xm = x
		}
		if y < ym {
			ym = y
		}
		if x+w > xM {
			xM = x + w
		}
		if y+h > yM {
			yM = y + h
		}
	}
	return xm, ym, xM - xm, yM - ym
}

type image struct {
	sync.Mutex
	bounds [4]float64
	image  []int32
	res    int
}

type result struct {
	input *Input

	sync.Mutex
	measures []int64
	image    *image
	no       int64
}

func (r *result) addPos(x, y float64) {
	r.image.addPos(x, y)
	r.Lock()
	defer r.Unlock()

	for i, in := range r.input.Measures {
		if x < in.X || x > in.X+in.W || y < in.Y || y > in.Y+in.H {
			continue
		}
		r.measures[i]++
	}
	r.no++
}

func (r *image) addPos(x, y float64) {
	// find the grid position where this should be drawn
	res := r.res
	x -= r.bounds[0]
	y -= r.bounds[1]
	m := r.bounds[2]
	if r.bounds[3] > r.bounds[2] {
		m = r.bounds[3]
	}
	x /= m
	y /= m
	xi := int(x * float64(res))
	yi := int(y * float64(res))
	if xi >= res {
		xi = res - 1
	}
	if yi >= res {
		yi = res - 1
	}
	r.Lock()
	r.image[xi+res*yi]++
	r.Unlock()
}
