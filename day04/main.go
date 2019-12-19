package main

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	//log, _ = zap.NewDevelopment()
	log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

func isValidP1(in int) bool {
	orig := in
	hasAdjacentDigits := false
	last := -1
	for in > 0 {
		div := in % 10
		in = in / 10
		if last != -1 && last == div {
			hasAdjacentDigits = true
		}
		if last != -1 && div > last {
			return false
		}
		last = div
	}
	if hasAdjacentDigits {
		logSugar.Debugf("Valid: %d", orig)
	}

	return hasAdjacentDigits
}

func isValidP2(in int) bool {
	logSugar.Debugf("Checking '%d'", in)
	dd := 0
	prev := 10
	countOfThisDigit := 0
	for in > 0 {
		div := in % 10
		in = in / 10
		if div > prev {
			// digits are out of order, fail
			return false
		} else if div == prev {
			countOfThisDigit++
			if countOfThisDigit == 1 {
				dd++
			} else if countOfThisDigit == 2 {
				dd--
			}
		} else {
			countOfThisDigit = 0
		}
		prev = div
	}
	return dd > 0
}

func main() {
	solve()
}

func solve() {
	lhs := 168630
	rhs := 718098

	count := 0
	for i := lhs; i <= rhs; i++ {
		if isValidP1(i) {
			count++
		}
	}
	fmt.Printf("Part1: %d\n", count)

	count = 0
	for i := lhs; i <= rhs; i++ {
		if isValidP2(i) {
			count++
		}
	}
	fmt.Printf("Part2: %d\n", count)
}
