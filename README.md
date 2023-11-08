# Hachigo - CHIP-8 Simulator

Hachigo is a simple CHIP-8 simulator.

## Features

### Emulated Instructions

Currently, Hachigo is limited to these instructions, with the only purpose to run `ibm-logo.ch8`, displaying the IBM logo.

[X] `00E0`: Clear Screen - Clears the display.
[X] `1NNN`: Jump - Jumps to a specified address.
[X] `6XNN`: Set Register VX - Sets a register VX to a specified value.
[X] `7XNN`: Add Value to Register VX - Adds a value to a register VX.
[X] `ANNN`: Set Index Register I - Sets the index register I to a specified value.
[X] `DXYN`: Display/Draw - Draws sprites on the screen.

## Usage

To run:
```bash
go run .
