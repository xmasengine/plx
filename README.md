# PLX

## Introduction

PLX is a PL/M and BASIC inspired compiler targeting 8/16 bits architectures,
written in Go. Pike PL/M, PLX is a high-level-syntax, low-level-effect
programming language.

While C is often called a low-level language, it abstracts the machine
and has no direct support for interrupts, chip-level I/O, ROM location,
or memory mapping. In C we rely on undefined behavior or assembly.

In PLX the compiler is not allowed to exhibit undefined behavior.
While there are 'dangerous' statements like DECLARE ... AT these have
defined semantics.  Statements like `INPUT`, `OUTPUT`,
and `INTERRUPT` allow true low-level programming without assembly.

The syntax uses easy-to-read full English UPPERCASE keywords, resembling
BASIC or PL/1 without being as verbose as COBOL or as obtuse as FORTRAN.

For easy portability and simplicity of compilation, PLX uses PIR as intermediate
representation. PIR is designed to be a FORTH like stack based language,
that then is compiled to assembly language, which then simulates a
stack machine on the target machine.

## Rationale

PLX is designed for the purpose of developing the Lord Of Christmas game
for the Sega Master System. While I started development in BASIC, I was not
satisfied with the result. And C is bad fit for small architectures. When I
found out about PL/M this inspired me to make PLX.

PIR implements a FORTH like stack machine with separate data and return stacks.
This simplifies the code generator and allows easier porting to other 8 bits
architectures like the NES.


## Credits

PLX uses [koron-go/z80](https://github.com/koron-go/z80) as the Z80 emulator
for testing, under the MIT License.

The Z80 assembler is based on [paulhankin/z80asm](https://github.com/paulhankin/z80asm),
forked and modified under the MIT License.

PLX uses the resource loader github.com/mrcook/smstilemap for loading
graphics, under the MIT license.

PLX uses github.com/user-none/go-chip-sn76489 as the sound chip emulator
for testing, under the MIT license.

PLX uses github.com/beevik/go6502 as the 6502 assembler, under the 2BSD license.

And Gary Kildall for his brilliant design of PL/M. We need such programming
languages to truly turn software development into software engineering.

