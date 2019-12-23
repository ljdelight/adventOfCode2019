package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	log, _ = zap.NewDevelopment()
	//log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

func main() {
	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}
	in := strings.Join(lines, "")

	// cols := 3
	// rows := 2
	cols := 25 // wide
	rows := 6  // tall

	image, nLayers, _, _ := readImage(in, rows, cols)

	minZeroCount := math.MaxInt64
	layerWithMinimalZeros := 0
	for i := 0; i < nLayers; i++ {
		zeroCount := countDigits(image[i], 0)
		if zeroCount > 0 && zeroCount < minZeroCount {
			layerWithMinimalZeros = i
			minZeroCount = zeroCount
		}
	}

	logSugar.Debugf("Layer with minimal zeros layer=%d count=%d\n", layerWithMinimalZeros, minZeroCount)

	countOnes := countDigits(image[layerWithMinimalZeros], 1)
	countTwos := countDigits(image[layerWithMinimalZeros], 2)
	logSugar.Debugf("layer=%d ones=%d twos=%d", layerWithMinimalZeros, countOnes, countTwos)
	fmt.Printf("Part1: %d\n", countOnes*countTwos)

	res := analyzeLayers(image, rows, cols)

	fmt.Println("Part2:")
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[i]); j++ {
			if res[i][j] == colorBlack {
				fmt.Print(" ")
			} else if res[i][j] == colorWhite {
				fmt.Print("b")
			} else {
				fmt.Print("s")
			}
		}
		fmt.Print("\n")
	}
}

var (
	colorBlack       = 0
	colorWhite       = 1
	colorTransparent = 2
)

func analyzeLayers(image [][][]int, rows int, cols int) [][]int {
	finalImg := make([][]int, rows)
	for i := range finalImg {
		finalImg[i] = make([]int, cols)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			for layer := 0; layer < len(image); layer++ {
				if image[layer][r][c] != colorTransparent {
					finalImg[r][c] = image[layer][r][c]
					break
				}
			}
		}
	}
	return finalImg
}

func countDigits(layer [][]int, digit int) int {
	count := 0
	for j := 0; j < len(layer); j++ {
		for k := 0; k < len(layer[j]); k++ {
			if layer[j][k] == digit {
				count++
			}
		}
	}
	return count
}

func readImage(img string, rows int, cols int) ([][][]int, int, int, int) {
	if len(img)%cols*rows != 0 {
		logSugar.Debugf("len(img)=%d cols=%d rows=%d", len(img), cols, rows)
		log.Fatal("likely corruption since the img doesn't match the cols*rows layers")
	}
	cellsPerLayer := cols * rows
	nLayers := len(img) / (cols * rows)
	log.Debug("", zap.Int("nLayers", nLayers), zap.Int("cols", cols), zap.Int("rows", rows), zap.Int("cellsPerLayer", cellsPerLayer))
	layers := make([][][]int, nLayers)
	for i := range layers {
		layers[i] = make([][]int, rows)
		for j := range layers[i] {
			layers[i][j] = make([]int, cols)
		}
	}

	rIdx := 0
	cIdx := 0
	layerIdx := 0
	for _, c := range img {
		//logSugar.Debugf("- %d %d %d (%d %d %d)", layerIdx, rIdx, cIdx, nLayers, rows, cols)
		layers[layerIdx][rIdx][cIdx] = MakeAtoi(string(c))

		cIdx++
		if cIdx == cols {
			cIdx = 0
			rIdx++
			if rIdx == rows {
				rIdx = 0
				layerIdx++
			}

		}
		// log.Debug(string(c))
	}
	//log.Debug("", zap.Any("img", layers))
	return layers, nLayers, rows, cols

}

// read and trim each line from the given filename
func readInputFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Warn("failed to close", zap.Error(err))
		}
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

// MakeAtoi is equivalent to strconv.Atoi but will panic on failure
func MakeAtoi(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}
