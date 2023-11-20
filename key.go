package main

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type Key struct {
	mu            sync.Mutex
	KeyboardEvent sdl.Keycode
	MappedKey     byte
	IsPressed     bool
}

func (k *Key) SetMappedKey() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.MappedKey = MapKey(byte(k.KeyboardEvent))
}

func (k *Key) GetMappedKey() byte {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.MappedKey
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
		return 0x1
	case byte('2'):
		return 0x2
	case byte('3'):
		return 0x3
	case byte('4'):
		return 0xC

	case byte('q'):
		return 0x4
	case byte('w'):
		return 0x5
	case byte('e'):
		return 0x6
	case byte('r'):
		return 0xD

	case byte('a'):
		return 0x7
	case byte('s'):
		return 0x8
	case byte('d'):
		return 0x9
	case byte('f'):
		return 0xE

	case byte('z'):
		return 0xA
	case byte('x'):
		return 0x0
	case byte('c'):
		return 0xB
	case byte('v'):
		return 0xF
	default:
		return c
	}
}
