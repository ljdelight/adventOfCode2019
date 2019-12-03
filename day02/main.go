package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

var (
	//log, _ = zap.NewDevelopment()
	log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

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
	input.Scan()
	line := input.Text()
	logSugar.Debugf("the input %s", line)

	program := readProgram(line)
	program[1] = 12
	program[2] = 2
	//program := readProgram("1,0,0,0,99")
	//program := readProgram("2,3,0,3,99")
	//program := readProgram("2,4,4,5,99,0")
	//program := readProgram("1,1,1,4,99,5,6,0,99")

	logSugar.Info("Program", program)

	result := run(program)
	logSugar.Info("Resulting", result)

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			program[1] = noun
			program[2] = verb
			result = run(program)
			if result[0] == 19690720 {
				fmt.Printf("Solution found noun=%d verb=%d %d\n", noun, verb, 100*noun+verb)
				return
			}
		}
	}

}

func run(programOriginal []int) []int {
	program := append([]int(nil), programOriginal...)
	ADD := 1
	MUL := 2
	HALT := 99

	pc := 0
	stop := false
	for !stop {
		switch program[pc] {
		case ADD:
			program[program[pc+3]] = program[program[pc+1]] + program[program[pc+2]]
		case MUL:
			program[program[pc+3]] = program[program[pc+1]] * program[program[pc+2]]
		case HALT:
			stop = true
		}
		pc += 4
	}

	return program
}

func readProgram(str string) []int {
	var res []int
	for _, opcode := range strings.Split(str, ",") {
		res = append(res, MakeAtoi(opcode))
	}
	return res
}

func MakeAtoi(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}
