package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

const cellsize = 1.0
const quad_leaf_size = 10.0
const nodata_value = -9999.0

type point2 struct {
	e float64
	n float64
}

type point3 struct {
	e float64
	n float64
	h float64
}

type quad struct {
	//  lower left extent of the quad
	llp point2
	//  top right extent of the quad
	trp point2

	llq *quad
	lrq *quad
	tlq *quad
	trq *quad

	points []point3

	leaf bool
}

func newQuad(llp, trp point2) *quad {
	e_range := trp.e - llp.e
	n_range := trp.n - llp.n
	if e_range < quad_leaf_size && n_range < quad_leaf_size {
		return &quad{
			llp,
			trp,
			nil,
			nil,
			nil,
			nil,
			nil,
			true,
		}
	}

	e_range /= 2.0
	n_range /= 2.0

	return &quad{
		llp,
		trp,
		newQuad(point2{llp.e, llp.n}, point2{llp.e + e_range, llp.n + n_range}), //  llq
		newQuad(point2{llp.e + e_range, llp.n}, point2{trp.e, llp.n + n_range}), //  lrq
		newQuad(point2{llp.e, llp.n + n_range}, point2{llp.e + e_range, trp.n}), //  tlq
		newQuad(point2{llp.e + e_range, llp.n + n_range}, point2{trp.e, trp.n}), //  trq
		nil,
		false,
	}
}

func (q *quad) add(p point3) {
	if !q.covers(p) {
		return
	}

	if q.leaf {
		q.points = append(q.points, p)
	} else {
		if q.llq.covers(p) {
			q.llq.add(p)
		}
		if q.lrq.covers(p) {
			q.lrq.add(p)
		}
		if q.tlq.covers(p) {
			q.tlq.add(p)
		}
		if q.trq.covers(p) {
			q.trq.add(p)
		}
	}
}

func (q *quad) covers(p point3) bool {
	if p.e < q.llp.e || p.e >= q.trp.e {
		return false
	}
	if p.n < q.llp.n || p.n >= q.trp.n {
		return false
	}

	return true
}

func (q *quad) queryRange(ll, tr point2) []point3 {
	if tr.e <= q.llp.e || ll.e > q.trp.e {
		return nil
	}
	if tr.n <= q.llp.n || ll.n > q.trp.n {
		return nil
	}

	var points []point3
	if q.leaf {
		for _, p := range q.points {
			if p.e < ll.e || p.e >= tr.e {
				continue
			}
			if p.n < ll.n || p.n >= tr.n {
				continue
			}
			points = append(points, p)
		}
	} else {
		points = append(points, q.llq.queryRange(ll, tr)...)
		points = append(points, q.lrq.queryRange(ll, tr)...)
		points = append(points, q.tlq.queryRange(ll, tr)...)
		points = append(points, q.trq.queryRange(ll, tr)...)
	}
	return points
}

var llFlag = flag.String("ll", "", "lower left")
var trFlag = flag.String("tr", "", "top left")
var inFlag = flag.String("in", "", "input file")
var outFlag = flag.String("out", "", "output file")

func main() {
	flag.Parse()

	if *llFlag == "" || *trFlag == "" || *inFlag == "" || *outFlag == "" {
		return
	}

	ll := point2{}
	if _, err := fmt.Sscanf(*llFlag, "%f,%f", &ll.e, &ll.n); err != nil {
		log.Fatal(err)
	}

	tr := point2{}
	if _, err := fmt.Sscanf(*trFlag, "%f,%f", &tr.e, &tr.n); err != nil {
		log.Fatal(err)
	}

	in, err := os.Open(*inFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	make_start := time.Now()
	q := newQuad(ll, tr)
	log.Printf("time to make quad: %v\n", time.Now().Sub(make_start))

	populate_start := time.Now()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		var p point3
		if _, err := fmt.Sscanf(scanner.Text(), "%f %f %f", &p.e, &p.n, &p.h); err != nil {
			log.Fatal(err)
		}
		q.add(p)
	}
	log.Printf("time to populate quad: %v\n", time.Now().Sub(populate_start))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(*outFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	//  header
	e_range := tr.e - ll.e
	n_range := tr.n - ll.n
	fmt.Fprintf(w, "ncols %d\n", int(e_range / cellsize))
	fmt.Fprintf(w, "nrows %d\n", int(n_range / cellsize))
	fmt.Fprintf(w, "xllcorner %f\n", ll.e)
	fmt.Fprintf(w, "yllcorner %f\n", ll.n)
	fmt.Fprintf(w, "cellsize %f\n", cellsize)
	fmt.Fprintf(w, "nodata_value %.3f\n", nodata_value)

	output_start := time.Now()
	for n := n_range - cellsize; n >= 0.0; n -= cellsize {
		for e := 0.0; e < e_range; e += cellsize {
			ll_find := point2{ll.e + e, ll.n + n}
			tr_find := point2{ll.e + e + cellsize, ll.n + n + cellsize}
			points := q.queryRange(ll_find, tr_find)
			h := math.MaxFloat64
			if len(points) == 0 {
				h = nodata_value
			} else {
				for _, p := range points {
					h = math.Min(h, p.h)
				}
			}
			fmt.Fprintf(w, "%.3f ", h)
		}
		fmt.Fprintf(w, "\n")
	}
	log.Printf("time to output: %v\n", time.Now().Sub(output_start))

}
