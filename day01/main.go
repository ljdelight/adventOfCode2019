package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var (
	//log, _ = zap.NewDevelopment()
	log, _ = zap.NewProduction()
)

func fuelRequired(mass int) int {
	return mass/3 - 2
}
func fuelRequiredWithAdjustment(mass int) int {
	accFuel := fuelRequired(mass)
	delta := fuelRequired(accFuel)
	for delta > 0 {
		accFuel += delta
		delta = fuelRequired(delta)
	}
	return accFuel
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("failed to open file", zap.Error(err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Warn("failed to close", zap.Error(err))
		}
	}()

	input := bufio.NewScanner(f)
	var moduleMass []int
	for input.Scan() {
		x, err := strconv.ParseInt(input.Text(), 10, 64)
		if err != nil {
			log.Fatal("failed parsing input", zap.Error(err))
		}
		moduleMass = append(moduleMass, int(x))
	}

	accFuel := 0
	for _, mass := range moduleMass {
		accFuel += fuelRequired(mass)
		log.Sugar().Debugf("mass=%d fuel=%d", mass, fuelRequired(mass))
	}
	fmt.Printf("Part 1: %d\n", accFuel)

	accFuel = 0
	for _, mass := range moduleMass {
		accFuel += fuelRequiredWithAdjustment(mass)
	}
	fmt.Printf("Part 2: %d\n", accFuel)

}
