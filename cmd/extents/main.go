package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	maxEasting := -math.MaxFloat64
	minEasting := math.MaxFloat64
	maxNorthing := -math.MaxFloat64
	minNorthing := math.MaxFloat64
	maxHeight := -math.MaxFloat64
	minHeight := math.MaxFloat64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var easting float64
		var northing float64
		var height float64
		if _, err := fmt.Sscanf(scanner.Text(), "%f %f %f", &easting, &northing, &height); err != nil {
			log.Fatal(err)
		}

		maxEasting = math.Max(maxEasting, easting)
		minEasting = math.Min(minEasting, easting)
		maxNorthing = math.Max(maxNorthing, northing)
		minNorthing = math.Min(minNorthing, northing)
		maxHeight = math.Max(maxHeight, height)
		minHeight = math.Min(minHeight, height)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%10.10f,%10.10f,%10.10f\n", minEasting, minNorthing, minHeight)
	fmt.Printf("%10.10f,%10.10f,%10.10f\n", maxEasting, maxNorthing, maxHeight)
}
