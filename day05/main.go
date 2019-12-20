package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
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
	memory[memory[ip+1]] = <-input
	return ip + 2
}

// Output (opcode=4) gets its argument and writes it to the output. The argument supports addressing modes POSITION and IMMEDIATE.
func Output(ip int, memory []int, output chan<- int) int {
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
	//test()

	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}

	memory := readProgram(strings.Join(lines, ""))
	input := make(chan int, 3000)
	output := make(chan int, 3000)

	input <- 5
	solve(memory, input, output)

	for data := range output {
		fmt.Printf("Part1: %d\n", data)
	}
}

func solve(memoryOriginal []int, input <-chan int, output chan<- int) {
	log.Debug("solving")
	memory := append([]int(nil), memoryOriginal...)

	ip := 0
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
