package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const PROGRAM_START_ADDRESS = 0x200
const FONT_START_ADDRESS = 0x050
const FONT_END_ADDRESS = 0x09F

// initiate PC
var PC uint16 = PROGRAM_START_ADDRESS
var instruction uint16
var font [80]byte = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

var V [16]byte = [16]byte{
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
}

var memory [4096]byte
var stack Stack[uint16]
var index uint16
var display Display
var k Key

// ignore the memory clock atm
func main() {
	loadFileToMemory(&memory)
	go chip()
	display.InitDisplay(&k)
}

func chip() {
	time.Sleep(2 * time.Second)
	initMemory(&memory)
	for {

		// fetch
		instruction = (uint16(0|memory[PC]) << 8) | uint16(0|(memory[PC+1]))
		PC += 2

		// decode & execute
		firstNibble := (instruction & uint16(61440)) >> 12
		switch firstNibble {
		case 0:
			switch twelveBitAddress := instruction & (uint16(0xFFFF) >> 4); twelveBitAddress {
			case uint16(0x00E0):
				display.ClearScreen()
			case uint16(0x00EE):
				address, _ := stack.Pop()
				PC = address.(uint16)
			default:
				fmt.Printf("unsupported %04X\n", twelveBitAddress)
			}
		case 1:
			PC = instruction & (uint16(0xFFFF) >> 4)
		case 2:
			stack.Push(PC)
			PC = instruction & (uint16(0xFFFF) >> 4)
		case 3:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			if V[secondNibble] == byte(instruction&uint16(0x00FF)) {
				PC += 1
			}
		case 4:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			if V[secondNibble] != byte(instruction&uint16(0x00FF)) {
				PC += 1
			}
		case 5:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)

			if V[secondNibble] == V[thirdNibble] {
				PC += 1
			}
		case 6:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			V[secondNibble] = byte(instruction & 0x00FF)
		case 7:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			V[secondNibble] += uint8((instruction & (uint16(0x00FF))))
		case 8:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)
			VX := &V[secondNibble]
			VY := &V[thirdNibble]

			switch instruction & uint16(0x000F) {
			case 0:
				*VX = *VY
			case 1:
				*VX = *VX | *VY
			case 2:
				*VX = *VX & *VY
			case 3:
				*VX = *VX ^ *VY
			case 4:
				result := *VX + *VY

				if result < *VX || result < *VY {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 5:
				*VX = *VX - *VY

				if *VX > *VY {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 6:
				*VX = *VY
				if (*VX & 0x01) == 1 {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}

				*VX = *VX >> 1
			case 7:
				*VX = *VY - *VX

				if *VY > *VX {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 0xE:
				*VX = *VY

				if (*VX & 0x80) == 1 {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}

				*VX = *VX << 1
			}
		case 9:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)

			if V[secondNibble] != V[thirdNibble] {
				PC += 1
			}
		case 0xA:
			index = instruction & (uint16(0xFFFF) >> 4)
		case 0xB:
			twelveBitAddress := instruction & (uint16(0xFFFF) >> 4)
			PC = twelveBitAddress + uint16(V[0])
		case 0xC:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			lastByte := byte(instruction & (0x00FF))
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			V[secondNibble] = byte(r.Intn(int(lastByte))) & (lastByte)
		case 0xD:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)
			N := (instruction & (uint16(0x000F)))

			y := int32(V[thirdNibble] & 31)

			V[0xF] = 0x0

			pointer := index
			for N > 0 {
				spriteAddress := memory[pointer]
				x := int32(V[secondNibble] & 63)

				for i := 7; i >= 0; i-- {
					bit := (spriteAddress >> i) & 1

					if bit == 1 {
						if display.IsPixelOn(x, y) {
							display.TurnOffPixel(x, y)

							V[0xF] = 1
						} else {
							display.TurnOnPixel(x, y)
						}
					}
					x++
				}
				pointer++

				y++
				N--
			}
		case 0xE:
			lastByte := byte(instruction & (0x00FF))
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			switch lastByte {
			case 0x9E:
				if byte(secondNibble) == key.MappedKey {
					PC += 2
				}
			case 0xA1:
				if byte(secondNibble) != key.MappedKey {
					PC += 2
				}
			}

		default:
			fmt.Printf("unsupported instruction: %X\n", firstNibble)
		}
	}
}

func initMemory(mem *[4096]byte) {
	j := 0
	for i := FONT_START_ADDRESS; i <= FONT_END_ADDRESS; i++ {
		mem[i] = font[j]
		j++
	}
}

func loadFileToMemory(mem *[4096]byte) {
	pointer := PROGRAM_START_ADDRESS
	// dat, err := os.ReadFile("test-roms/chip8-test-rom/test_opcode.ch8")
	dat, err := os.ReadFile("test-roms/ibm-logo.ch8")

	if err != nil {
		panic(err)
	}

	for _, b := range dat {
		mem[pointer] = b
		pointer++
	}
}
