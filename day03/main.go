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
	//log, _ = zap.NewDevelopment()
	log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

type direction rune

var (
	dirUp    direction = 'U'
	dirDown  direction = 'D'
	dirLeft  direction = 'L'
	dirRight direction = 'R'
)

type pair struct {
	a, b int
}
type Point struct {
	x, y int
}

func main() {
	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}

	var commands [][]string
	for _, line := range lines {
		commands = append(commands, strings.Split(line, ","))
	}
	log.Debug("Read the input", zap.Strings("commands", commands[0]))

	wire1Pts := buildSegments(commands[0])
	wire2Pts := buildSegments(commands[1])

	minManhattan := math.MaxInt32
	minWire := math.MaxInt32
	for k, v := range wire1Pts {
		// if a point is in both maps, it's a collision
		if val, ok := wire2Pts[k]; ok {
			if val.a < minManhattan {
				minManhattan = val.a
			}
			if v.b+val.b < minWire {
				minWire = v.b + val.b
			}
		}
	}
	fmt.Printf("minManhattan: %d\n", minManhattan)
	fmt.Printf("minWire: %d\n", minWire)
}

func buildSegments(commands []string) map[Point]pair {
	allpoints := make(map[Point]pair)
	curr := Point{x: 0, y: 0}
	ttlwalked := 0
	for _, cmd := range commands {
		dir := direction([]rune(cmd)[0])
		len := MakeAtoi(cmd[1:])

		t := curr
		for i := 1; i <= len; i++ {
			ttlwalked++
			t = walk(t, dir, 1)
			allpoints[t] = pair{a: int(math.Abs(float64(t.x)) + math.Abs(float64(t.y))), b: ttlwalked}
		}

		curr = t
	}
	return allpoints
}

func walk(p Point, dir direction, len int) Point {
	switch dir {
	case dirUp:
		p.y += len
	case dirDown:
		p.y -= len
	case dirLeft:
		p.x -= len
	case dirRight:
		p.x += len
	}
	return p
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
