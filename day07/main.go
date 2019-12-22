package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	log, _ = zap.NewDevelopment()
	//log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

const (
	// Parameters can be of two types:
	//   - Position mode:  the arg is the memory address of the value to use ('10' is an address and results in the lookup memory['10'])
	//   - Immediate mode: the arg should be interpreted as a literal ('15' is the value, 15)
	// NOTE: Parameters that an instruction writes to will never be in immediate mode.
	ADD          = 1
	MUL          = 2
	INPUT        = 3
	OUTPUT       = 4
	JMP_IF_TRUE  = 5
	JMP_IF_FALSE = 6
	LESS_THAN    = 7
	EQUALS       = 8
	HALT         = 99
)

const (
	POSITION_MODE  int = 0
	IMMEDIATE_MODE int = 1
)

// type OpCode interface {
// 	// execute the opcode given the instruction pointer and memory, returning the next instruction index (skipping the OpCode's arguments)
// 	func Exec(ip int, memory []int) (int)
// }

// Add (opcode=1) the first two arguments and store into the third. The first two argument addressing modes support POSITION and IMMEDIATE.
func Add(ip int, memory []int) int {
	log.Debug("Add intruction", zap.Int("instructionPointer", ip))
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}

	arg2 := memory[ip+2]
	if (instruction / 1000 % 10) == POSITION_MODE {
		arg2 = memory[arg2]
	}

	memory[memory[ip+3]] = arg1 + arg2
	return ip + 4
}

// Multiply (opcode=2) the first two arguments and store into the third. The first two argument addressing modes support POSITION and IMMEDIATE.
func Multiply(ip int, memory []int) int {
	log.Debug("Multiply intruction", zap.Int("instructionPointer", ip))
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}

	arg2 := memory[ip+2]
	if (instruction / 1000 % 10) == POSITION_MODE {
		arg2 = memory[arg2]
	}

	memory[memory[ip+3]] = arg1 * arg2
	return ip + 4
}

// Input (opcode=3) takes a single integer from input and saves it to the position given by its (only) argument.
func Input(ip int, memory []int, input <-chan int) int {
	log.Debug("Input intruction", zap.Int("instructionPointer", ip))
	memory[memory[ip+1]] = <-input
	return ip + 2
}

// Output (opcode=4) gets its argument and writes it to the output. The argument supports addressing modes POSITION and IMMEDIATE.
func Output(ip int, memory []int, output chan<- int) int {
	log.Debug("Output intruction", zap.Int("instructionPointer", ip))
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}

	output <- arg1
	return ip + 2
}

// JumpIfTrue (opcode=5): if the first argument is non-zero, then set the instruction pointer to the value from the second argument. Otherwise do nothing. The arguments support addressing modes POSITION and IMMEDIATE.
func JumpIfTrue(ip int, memory []int) int {
	log.Debug("JumpIfTrue intruction", zap.Int("instructionPointer", ip))
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}
	if arg1 != 0 {
		arg2 := memory[ip+2]
		if (instruction / 1000 % 10) == POSITION_MODE {
			arg2 = memory[arg2]
		}
		return arg2
	} else {
		return ip + 3
	}
}

// JumpIfFalse (opcode=6): if the first argument is non-zero, then set the instruction pointer to the value from the second argument. Otherwise do nothing. The arguments support addressing modes POSITION and IMMEDIATE.
func JumpIfFalse(ip int, memory []int) int {
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}
	if arg1 == 0 {
		arg2 := memory[ip+2]
		if (instruction / 1000 % 10) == POSITION_MODE {
			arg2 = memory[arg2]
		}
		return arg2
	} else {
		return ip + 3
	}
}

// LessThan (opcode=7) takes two arguments and if arg1 is less than arg2 write 1 into the third location of the third argument, otherwise write 0. The first two argument addressing modes support POSITION and IMMEDIATE.
func LessThan(ip int, memory []int) int {
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}

	arg2 := memory[ip+2]
	if (instruction / 1000 % 10) == POSITION_MODE {
		arg2 = memory[arg2]
	}

	if arg1 < arg2 {
		memory[memory[ip+3]] = 1
	} else {
		memory[memory[ip+3]] = 0
	}

	return ip + 4
}

// OpEquals (opcode=8) takes two arguments and if arg1 equals arg2 write 1 into the third location of the third argument, otherwise write 0. The first two argument addressing modes support POSITION and IMMEDIATE.
func OpEquals(ip int, memory []int) int {
	instruction := memory[ip]
	arg1 := memory[ip+1]
	if (instruction / 100 % 10) == POSITION_MODE {
		arg1 = memory[arg1]
	}

	arg2 := memory[ip+2]
	if (instruction / 1000 % 10) == POSITION_MODE {
		arg2 = memory[arg2]
	}

	if arg1 == arg2 {
		memory[memory[ip+3]] = 1
	} else {
		memory[memory[ip+3]] = 0
	}

	return ip + 4
}

func main() {
	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}
	memory := readProgram(strings.Join(lines, ""))
	p1(memory)
	p2(memory)
}

func p1(memory []int) {
	permutations := perms([]int{0, 1, 2, 3, 4}, 0, 4)

	var bestPhaseSettings []int
	max := math.MinInt32
	log.Debug("Memory", zap.Ints("mem", memory))
	for _, p := range permutations {
		phaseSettings := p
		res := solveWithPhaseSettings(memory, phaseSettings, 4)
		//fmt.Printf("%v\n", phaseSettings)
		//log.Info("Phase setting result", zap.Ints("phaseSettings", phaseSettings), zap.Int("thruster", res))
		if res > max {
			bestPhaseSettings = phaseSettings
			max = res
		}

	}
	fmt.Printf("Part1: %d %v\n", max, bestPhaseSettings)
}

func p2(memory []int) {
	permutations := perms([]int{5, 6, 7, 8, 9}, 0, 4)

	var bestPhaseSettings []int
	max := math.MinInt32
	for _, p := range permutations {
		phaseSettings := p
		res := solveWithPhaseSettings(memory, phaseSettings, 0)
		//fmt.Printf("%v\n", phaseSettings)
		//log.Info("Phase setting result", zap.Ints("phaseSettings", phaseSettings), zap.Int("thruster", res))
		if res > max {
			bestPhaseSettings = phaseSettings
			max = res
		}
	}
	fmt.Printf("Part2: %d %v\n", max, bestPhaseSettings)
}

func perms(arr []int, i int, n int) (res [][]int) {
	if i == n {
		t := make([]int, len(arr))
		copy(t, arr)
		return append(res, t)
	}
	for j := 0; j <= n; j++ {
		arr[i], arr[j] = arr[j], arr[i]
		res = append(res, perms(arr, i+1, n)...)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return res
}
func solveWithPhaseSettings(memoryOriginal []int, phaseSettings []int, r int) int {
	var channels [5]chan int
	for i := range channels {
		channels[i] = make(chan int, 3)
	}

	var allMemory [][]int
	for i := 0; i < 5; i++ {
		temp := append([]int(nil), memoryOriginal...)
		allMemory = append(allMemory, temp)
	}

	var ips [5]int

	// first write the phase settings
	for idx, phase := range phaseSettings {
		channels[idx] <- phase
	}
	// the input to the first amp is 0
	channels[0] <- 0

	// the amps all tie together
	for allMemory[4][ips[4]]%10 != 9 {
		go func() {
			logSugar.Debugf("Amp %d", 0)
			ips[0] = solve(ips[0], allMemory[0], channels[0], channels[1])
		}()

		go func() {
			logSugar.Debugf("Amp %d", 1)
			ips[1] = solve(ips[1], allMemory[1], channels[1], channels[2])
		}()

		go func() {
			logSugar.Debugf("Amp %d", 2)
			ips[2] = solve(ips[2], allMemory[2], channels[2], channels[3])
		}()

		go func() {
			logSugar.Debugf("Amp %d", 3)
			ips[3] = solve(ips[3], allMemory[3], channels[3], channels[4])
		}()

		logSugar.Debugf("Amp %d", 4)
		ips[4] = solve(ips[4], allMemory[4], channels[4], channels[r])
	}
	for i := range channels {
		close(channels[i])
	}

	// the output of the last amp goes to the thrusters
	return <-channels[r]
}

func solve(ip int, memory []int, input <-chan int, output chan<- int) int {
	log.Debug("solving")

	stop := false
	for !stop {
		instruction := (memory[ip]/10)%10*10 + (memory[ip] % 10)
		logSugar.Debug("Processing ip ", ip, memory[ip])
		switch instruction {
		case ADD:
			ip = Add(ip, memory)
		case MUL:
			ip = Multiply(ip, memory)
		case INPUT:
			ip = Input(ip, memory, input)
		case OUTPUT:
			ip = Output(ip, memory, output)
		case JMP_IF_TRUE:
			ip = JumpIfTrue(ip, memory)
		case JMP_IF_FALSE:
			ip = JumpIfFalse(ip, memory)
		case LESS_THAN:
			ip = LessThan(ip, memory)
		case EQUALS:
			ip = OpEquals(ip, memory)
		case HALT:
			stop = true
		default:
			log.Fatal("Failed")
		}
	}
	log.Debug("Memory", zap.Ints("memory", memory))
	return ip
}

func readProgram(str string) (res []int) {
	for _, opcode := range strings.Split(str, ",") {
		res = append(res, MakeAtoi(opcode))
	}
	return
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

func test() {
	addTests := []struct {
		ip         int
		memory     []int
		wantIp     int
		wantMemory []int
	}{
		{0, []int{1, 4, 3, 4, 33}, 4, []int{1, 4, 3, 4, 37}},       // 4 position mode is 33 + 3 position mode is 4 == 37
		{0, []int{101, 4, 3, 4, 33}, 4, []int{101, 4, 3, 4, 8}},    // 4 immediate mode is 4 + 3 position mode is 4 == 8
		{0, []int{1001, 4, 3, 4, 33}, 4, []int{1001, 4, 3, 4, 36}}, // 4 position mode is 33 + 3 immediate mode is 3 == 36
		{0, []int{1101, 4, 3, 4, 33}, 4, []int{1101, 4, 3, 4, 7}},  // 4 immediate mode is 4 + 3 immediate mode is 3 == 7
	}

	for _, test := range addTests {
		ip := Add(test.ip, test.memory)
		if ip != test.wantIp || !reflect.DeepEqual(test.memory, test.wantMemory) {
			log.Error("add failed", zap.Int("ip", ip), zap.Int("wantIp", test.wantIp), zap.Ints("memory", test.memory), zap.Ints("wantMemory", test.wantMemory))
		}
	}

	multiplyTests := []struct {
		ip         int
		memory     []int
		wantIp     int
		wantMemory []int
	}{
		{0, []int{2, 4, 3, 4, 33}, 4, []int{2, 4, 3, 4, 132}},      // 4 position mode is 33 * 3 position mode is 4 == 132
		{0, []int{102, 4, 3, 4, 33}, 4, []int{102, 4, 3, 4, 16}},   // 4 immediate mode is 4 * 3 position mode is 4 == 16
		{0, []int{1002, 4, 3, 4, 33}, 4, []int{1002, 4, 3, 4, 99}}, // 4 position mode is 33 * 3 immediate mode is 3 == 99
		{0, []int{1102, 4, 3, 4, 33}, 4, []int{1102, 4, 3, 4, 12}}, // 4 immediate mode is 4 * 3 immediate mode is 3 == 12
	}

	for _, test := range multiplyTests {
		ip := Multiply(test.ip, test.memory)
		if ip != test.wantIp || !reflect.DeepEqual(test.memory, test.wantMemory) {
			log.Error("multiply failed", zap.Int("ip", ip), zap.Int("wantIp", test.wantIp), zap.Ints("memory", test.memory), zap.Ints("wantMemory", test.wantMemory))
		}
	}
}
