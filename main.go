package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/vron/bounces/bounces"
)

var (
	fTargetPrec float64
	fMaxIts     int64
	fMinIts     int64
	fOutput     string
	fRes        int
	fLog        bool
)

func init() {
	flag.Int64Var(&fMaxIts, "max", 1e9, "maximum number of its")
	flag.Int64Var(&fMinIts, "min", 1e4, "minimum number of its")
	flag.IntVar(&fRes, "res", 1024, "image resolution")
	flag.Float64Var(&fTargetPrec, "p", 1e-3, "approx target precision")
	flag.StringVar(&fOutput, "o", "./", "folder in which to save images")
	flag.BoolVar(&fLog, "log", true, "plot in logarithmic space")
}

func main() {
	flag.Parse()
	normalizeArgs()

	// run all the provided files
	for _, p := range flag.Args() {
		runFile(p)
	}

	// also run from stdin if no input
	if flag.NArg() == 0 {
		buf, err := ioutil.ReadAll(os.Stdin)
		fatal(err)
		if len(buf) > 0 {
			runInput("stdin", bytes.NewBuffer(buf))
		}
	}
}

func normalizeArgs() {
	if fMinIts < 1 {
		fMinIts = 1
	}
	if fMaxIts < fMinIts {
		fMaxIts = fMinIts * 2
	}
	if fTargetPrec < 0 {
		fTargetPrec = 0
	}
}

func runFile(p string) {
	file := filepath.Base(p)
	name := strings.TrimSuffix(file, filepath.Ext(file))

	f, err := os.Open(p)
	fatal(err, "error opening input file:")
	defer f.Close()

	runInput(name, f)
}

func runInput(name string, r io.Reader) {
	input := bounces.ParseInput(r, fatal)
	input.ImageRes = fRes
	input.ImagePath = filepath.Join(fOutput, name+".p")
	fmt.Printf("%20v ", name)
	res := bounces.Run(input, fTargetPrec, fMinIts, fMaxIts)
	fmt.Printf("%12v\t%.4g%%\t±%.4g\n", input.Measures[0].Name, 100*res.Measures[0], 100*res.MeasureErrors[0])
	for mi := range res.Measures[1:] {
		fmt.Printf("%20v %12v\t%.4g%%\t±%.4g\n", name, input.Measures[mi+1].Name, 100*res.Measures[mi+1], 100*res.MeasureErrors[mi+1])
	}

	// Write the image we might want to look at
	path := filepath.Join(fOutput, name+".png")
	f, err := os.Create(path)
	fatal(err)
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer func() {
		fatal(buf.Flush())
	}()
	pl := math.Log1p
	if !fLog {
		pl = func(a float64) float64 { return a }
	}
	fatal(res.Draw(buf, pl))
}

func fatal(e error, s ...interface{}) {
	if e != nil {
		str := fmt.Sprint(s...)
		log.Fatalln(str, e)
	}
}
