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

var fontLocationAddress [16]byte

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
var sound Sound
var k Key

var timer Timer
var soundTimer Timer
var destroy bool = false

// ignore the memory clock atm
func main() {
	loadFileToMemory(&memory)
	go chip()
	k.SetIsPressed(false)
	sound.InitSound()
	display.InitDisplay()
}

func chip() {
	time.Sleep(5 * time.Second)
	initMemory(&memory)
	initFontLocationAddress()
	initTimers()
	timer.Decrement()
	soundTimer.Decrement()
	for {

		// fetch
		instruction = (uint16(0|memory[PC]) << 8) | uint16(0|(memory[PC+1]))
		PC += 2

		// decode & execute
		firstNibble := (instruction & uint16(61440)) >> 12
		switch firstNibble {
		case 0:
			switch twelveBitAddress := instruction & (uint16(0xFFFF) >> 4); twelveBitAddress {
			case 0x0E0:
				display.SendInstruction(Instruction{
					Fn: func(args ...any) {
						display.ClearScreen()
					},
					args: []any{},
				})
				for len(display.Instructions) > 0 {
				}
			case 0x00EE:
				address, _ := stack.Pop()
				PC = address.(uint16)
			default:
				fmt.Printf("unsupported %04X\n", twelveBitAddress)
			}
		case 1:
			PC = instruction & (uint16(0xFFFF) >> 4)
		case 2:
			stack.Push(PC)
			PC = instruction & 0x0FFF
		case 3:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			if V[secondNibble] == byte(instruction&uint16(0x00FF)) {
				PC += 2
			}
		case 4:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			if V[secondNibble] != byte(instruction&uint16(0x00FF)) {
				PC += 2
			}
		case 5:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)

			if V[secondNibble] == V[thirdNibble] {
				PC += 2
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
				*VX = result
				if result < *VX || result < *VY {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 5:
				minuend, subtrahend := *VX, *VY
				*VX = *VX - *VY

				if minuend >= subtrahend {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 6:
				// *VX = *VY

				flag := *VX & 0x01
				*VX >>= 1
				V[0xF] = flag
			case 7:
				*VX = *VY - *VX

				if *VY > *VX {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
			case 0xE:
				// *VX = *VY
				flag := (*VX >> 7) & 0x01
				*VX <<= 1
				V[0xF] = flag
			default:
				fmt.Printf("unsupported instruction: %X\n", instruction)
			}
		case 9:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			thirdNibble := ((instruction & (uint16(0x00F0))) >> 4)

			if V[secondNibble] != V[thirdNibble] {
				PC += 2
			}
		case 0xA:
			index = instruction & (uint16(0xFFFF) >> 4)
		case 0xB:
			// first logic
			twelveBitAddress := instruction & 0x0FFF
			PC = twelveBitAddress + uint16(V[0])

			// second logic, might need to adjust on other system
			// second := ((instruction & uint16(0x0F00)) >> 8)
			// eightBitAddress := instruction & (uint16(0x00FF))
			// PC = eightBitAddress + uint16(V[second])
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
			VF := &V[0xF]

			pointer := index
			for N > 0 {
				spriteAddress := memory[pointer]
				x := int32(V[secondNibble] & 63)

				for i := 7; i >= 0; i-- {
					bit := (spriteAddress >> i) & 1

					if bit == 1 {
						if display.IsPixelOn(x, y) {
							display.SendInstruction(Instruction{Fn: func(args ...any) {
								display.TurnOffPixel(args[0].(int32), args[1].(int32))
								*VF = 1
							}, args: []any{x, y}})
						} else {
							display.SendInstruction(Instruction{Fn: func(args ...any) {
								display.TurnOnPixel(args[0].(int32), args[1].(int32))
							}, args: []any{x, y}})

						}
					}
					x++
				}
				pointer++

				y++
				N--
			}
			display.SendInstruction(Instruction{Fn: func(args ...any) {
				display.Window.UpdateSurface()
			}, args: []any{}})
			for len(display.Instructions) > 0 {
			}
		case 0xE:
			lastByte := byte(instruction & (0x00FF))
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)

			switch lastByte {
			case 0x9E:
				ck := keys[V[secondNibble]]
				if ck != nil && ck.GetIsPressed() {
					PC += 2
				}
			case 0xA1:
				ck := keys[V[secondNibble]]
				if ck == nil || (ck != nil && !ck.GetIsPressed()) {
					PC += 2
				}
			default:
				fmt.Printf("unsupported instruction: %X\n", instruction)
			}
		case 0xF:
			secondNibble := ((instruction & (uint16(0x0F00))) >> 8)
			lastByte := byte(instruction & (0x00FF))

			switch lastByte {
			case 0x07:
				V[secondNibble] = timer.GetTime()
			case 0x15:
				timer.SetTime(V[secondNibble])
			case 0x18:
				soundTimer.SetTime(V[secondNibble])
			case 0x1E:
				index += uint16(V[secondNibble])

				if index > 0x0FFF {
					V[0xF] = 1
				}
			case 0x0A:
				if globalKey.GetIsPressed() {
					V[secondNibble] = globalKey.GetMappedKey()
				} else {
					PC -= 2
				}
			case 0x29:
				index = uint16(fontLocationAddress[(V[secondNibble])])
			case 0x33:
				ptr := index
				parseAndStore := func(n byte) {

					for i := 100; i > 0; i /= 10 {
						if n/byte(i) == 0 {
							memory[ptr] = 0
							ptr++
							continue
						}

						memory[ptr] = (n / byte(i))
						n = n % byte(i)
						ptr++
					}
				}

				parseAndStore(byte(V[secondNibble]))
			case 0x55:
				ptr := index
				for i := byte(0); i <= byte(secondNibble); i++ {
					memory[ptr] = V[i]
					ptr++
				}
			case 0x65:
				ptr := index
				for i := byte(0); i <= byte(secondNibble); i++ {
					V[i] = memory[ptr]
					ptr++
				}
			default:
				fmt.Printf("unsupported instruction: %X\n", instruction)
			}

		default:
			fmt.Printf("unsupported instruction: %X\n", instruction)
		}
		time.Sleep(time.Second / 2500)
	}
}

func initMemory(mem *[4096]byte) {
	j := 0
	for i := FONT_START_ADDRESS; i <= FONT_END_ADDRESS; i++ {
		mem[i] = font[j]
		j++
	}
}

func initFontLocationAddress() {
	inc := 0
	for i := 0; i <= 0xF; i++ {
		fontLocationAddress[i] = FONT_START_ADDRESS + byte(inc)
		inc += 5
	}
}

func loadFileToMemory(mem *[4096]byte) {
	pointer := PROGRAM_START_ADDRESS
	// dat, err := os.ReadFile("test-roms/chip8-test-rom/6-keypad.ch8")
	// dat, err := os.ReadFile("test-roms/chip8-test-rom/br8kout.ch8")
	dat, err := os.ReadFile("test-roms/chip8-test-rom/7-beep.ch8")
	// dat, err := os.ReadFile("test-roms/chip8-test-rom/chipquarium.ch8")
	// dat, err := os.ReadFile("test-roms/ibm-logo.ch8")

	if err != nil {
		panic(err)
	}

	for _, b := range dat {
		mem[pointer] = b
		pointer++
	}
}

func initTimers() {
	timer = Timer{Time: 0x0}
	soundTimer = Timer{Time: 0x0, TimerCallback: []func(){sound.Play}, TimerZeroCallback: []func(){sound.Stop}}
}

func printProgam() {
	count := 0
	dat, _ := os.ReadFile("test-roms/chip8-test-rom/delay_timer_test.ch8")
	for i, b := range dat {
		fmt.Printf("%02X", b)
		count++

		if count == 2 {
			count = 0
			fmt.Printf(" memory address: %04X\n", i+(0x200-1))
		}
	}
}
