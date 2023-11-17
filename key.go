package main

import "github.com/veandco/go-sdl2/sdl"

type Key struct {
	KeyboardEvent sdl.Keycode
	MappedKey     byte
	IsPressed     bool
}

func (k *Key) SetMappedKey() {
	k.MappedKey = MapKey(byte(k.KeyboardEvent))
}

/*
Key Mapping:
1 -> 1
2 -> 2
3 -> 3
4 -> C

q -> 4
w -> 5
e -> 6
r -> D

a -> 7
s -> 8
d -> 9
f -> E

z -> A
x -> 0
c -> B
v -> F
*/
func MapKey(c byte) byte {
	switch c {
	case byte('1'):
		return byte('1')
	case byte('2'):
		return byte('2')
	case byte('3'):
		return byte('3')
	case byte('4'):
		return byte('C')

	case byte('q'):
		return byte('4')
	case byte('w'):
		return byte('5')
	case byte('e'):
		return byte('6')
	case byte('r'):
		return byte('D')

	case byte('a'):
		return byte('7')
	case byte('s'):
		return byte('8')
	case byte('d'):
		return byte('9')
	case byte('f'):
		return byte('E')

	case byte('z'):
		return byte('A')
	case byte('x'):
		return byte('0')
	case byte('c'):
		return byte('B')
	case byte('v'):
		return byte('F')
	default:
		return c
	}
}
