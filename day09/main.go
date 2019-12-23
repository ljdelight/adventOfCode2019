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

const (
	// Parameters can be of two types:
	//   - Position mode:  the arg is the memory address of the value to use ('10' is an address and results in the lookup memory['10'])
	//   - Immediate mode: the arg should be interpreted as a literal ('15' is the value, 15)
	// NOTE: Parameters that an instruction writes to will never be in immediate mode.
	ADD               = 1
	MUL               = 2
	INPUT             = 3
	OUTPUT            = 4
	JMP_IF_TRUE       = 5
	JMP_IF_FALSE      = 6
	LESS_THAN         = 7
	EQUALS            = 8
	ADJ_RELATIVE_BASE = 9
	HALT              = 99
)

const (
	POSITION_MODE  int = 0
	IMMEDIATE_MODE int = 1
	RELATIVE_MODE  int = 2
)

func MakeComputer(memory []int, input <-chan int, output chan<- int) *Computer {
	if input == nil {
		input = make(chan int, 3000)
	}
	if output == nil {
		output = make(chan int, 3000)
	}
	mem := make([]int, 3000)
	copy(mem, memory)
	c := Computer{memory: mem, input: input, output: output}
	return &c
}

type Computer struct {
	relativeBase int
	ip           int
	memory       []int
	input        <-chan int
	output       chan<- int
}

// E the instruction, return false to HALT execution
func (c *Computer) E() bool {
	stop := false
	instruction := (c.memory[c.ip]/10)%10*10 + (c.memory[c.ip] % 10)
	logSugar.Debug("Processing instruction ", c.ip, c.memory[c.ip])
	switch instruction {
	case ADD:
		c.Add()
	case MUL:
		c.Multiply()
	case INPUT:
		c.Input()
	case OUTPUT:
		c.Output()
	case JMP_IF_TRUE:
		c.JumpIfTrue()
	case JMP_IF_FALSE:
		c.JumpIfFalse()
	case LESS_THAN:
		c.LessThan()
	case EQUALS:
		c.OpEquals()
	case ADJ_RELATIVE_BASE:
		c.OpAdjustRelativeBase()
	case HALT:
		stop = true
	default:
		log.Fatal("instruction does not exist", zap.Int("instruction", instruction))
	}
	return stop
}

func (c *Computer) arg(pos int) *int {
	mode := c.memory[c.ip] / 100
	for i := 0; i < pos; i++ {
		mode = mode / 10
	}
	mode = mode % 10

	switch mode {
	case IMMEDIATE_MODE:
		return &c.memory[c.ip+1+pos]
	case POSITION_MODE:
		return &c.memory[c.memory[c.ip+1+pos]]
	case RELATIVE_MODE:
		return &c.memory[c.relativeBase+c.memory[c.ip+1+pos]]
	default:
		panic("unknonwn addressing mode")
	}
}

// Add (opcode=1) the first two arguments and store into the third. The first two argument addressing modes support POSITION and IMMEDIATE.
func (c *Computer) Add() {
	log.Debug("Add instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	arg2 := c.arg(1)
	out := c.arg(2)
	*out = *arg1 + *arg2
	c.ip += 4
}

// Multiply (opcode=2) the first two arguments and store into the third. The first two argument addressing modes support POSITION and IMMEDIATE.
func (c *Computer) Multiply() {
	log.Debug("Multiply instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	arg2 := c.arg(1)
	out := c.arg(2)
	*out = (*arg1) * (*arg2)
	c.ip += 4
}

// Input (opcode=3) takes a single integer from input and saves it to the position given by its (only) argument.
func (c *Computer) Input() {
	log.Debug("Input instruction", zap.Int("instructionPointer", c.ip))
	out := c.arg(0)
	*out = <-c.input
	c.ip += 2
}

// Output (opcode=4) gets its argument and writes it to the output. The argument supports addressing modes POSITION and IMMEDIATE.
func (c *Computer) Output() {
	log.Debug("Output instruction", zap.Int("instructionPointer", c.ip))
	arg := c.arg(0)
	log.Debug("output value", zap.Int("output", *arg))
	c.output <- *arg
	c.ip += 2
}

// JumpIfTrue (opcode=5): if the first argument is non-zero, then set the instruction pointer to the value from the second argument. Otherwise do nothing. The arguments support addressing modes POSITION and IMMEDIATE.
func (c *Computer) JumpIfTrue() {
	log.Debug("JumpIfTrue instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	if *arg1 != 0 {
		arg2 := c.arg(1)
		c.ip = *arg2
	} else {
		c.ip += 3
	}
}

// JumpIfFalse (opcode=6): if the first argument is non-zero, then set the instruction pointer to the value from the second argument. Otherwise do nothing. The arguments support addressing modes POSITION and IMMEDIATE.
func (c *Computer) JumpIfFalse() {
	log.Debug("JumpIfFalse instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	if *arg1 == 0 {
		arg2 := c.arg(1)
		c.ip = *arg2
	} else {
		c.ip += 3
	}
}

// LessThan (opcode=7) takes two arguments and if arg1 is less than arg2 write 1 into the third location of the third argument, otherwise write 0. The first two argument addressing modes support POSITION and IMMEDIATE.
func (c *Computer) LessThan() {
	log.Debug("LessThan instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	arg2 := c.arg(1)
	out := c.arg(2)
	if *arg1 < *arg2 {
		*out = 1
	} else {
		*out = 0
	}

	c.ip += 4
}

// OpEquals (opcode=8) takes two arguments and if arg1 equals arg2 write 1 into the third location of the third argument, otherwise write 0. The first two argument addressing modes support POSITION and IMMEDIATE.
func (c *Computer) OpEquals() {
	log.Debug("OpEquals instruction", zap.Int("instructionPointer", c.ip))

	arg1 := c.arg(0)
	arg2 := c.arg(1)
	out := c.arg(2)
	if *arg1 == *arg2 {
		*out = 1
	} else {
		*out = 0
	}

	c.ip += 4
}

// OpAdjustRelativeBase (opcode=9) takes an adjustment to the relative base. The first two argument addressing modes support POSITION and IMMEDIATE.
func (c *Computer) OpAdjustRelativeBase() {
	log.Debug("OpAdjustRelativeBase instruction", zap.Int("instructionPointer", c.ip))

	arg := c.arg(0)
	c.relativeBase += *arg
	c.ip += 2
}

func main() {
	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}

	memory := readProgram(strings.Join(lines, ""))
	input := make(chan int, 3000)
	output := make(chan int, 3000)

	input <- 2
	solve(memory, input, output)
	close(input)
	close(output)
	for data := range output {
		fmt.Printf("Part1: %d\n", data)
	}
}

func solve(memory []int, input <-chan int, output chan<- int) {
	c := MakeComputer(memory, input, output)
	stop := false
	for !stop {
		stop = c.E()
	}
	//log.Debug("Memory", zap.Ints("memory", c.memory))
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
