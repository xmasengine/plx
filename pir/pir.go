// package pir implement the PLX Intermediate Representation
package pir

/*

PIR is an intermediate representation that models an abstract machine
with a data stack, a return stack, variables in read/write locations,
4 16 bits registers and 4 16 bits temporaries, and data in read only locations.

Considerations on the PIR a virtual machine architecture.

On the stack: while mixed stack is easier to implement, a split call and data
stack is a lot safer, makes tasks easier, and on some 8/16 bits
architectures the hardware stack is small or only available just
for that purpose. A data stack is not strictly needed if there is no recursion,
but it is safer/easier for deep calls.

On variable use: it is necessary to allow loading and storing to
variables in RAM. Also access to ROM is definitely needed for data storage.

On virtual register size: Since we have to store adesses is is easiest
to simulate a 16 bits machine with 16bits registers that have H and L bytes.

On the virtual registers. It would be possible to use from 0 up to 16 virtual
16 bits registers. The less registers, the easier to implement on
register starved machines like the NES, but the more slow stack operations
or storage to memory will be needed to spill the register contents.
It is possible to use real registers for the top of stack or next of stack,
but in practice that is difficult to manage.
The more registers the better the performance on machines that have more real
registers, but the more registers must be emulated by using zero page or RAM
on machines with less real registers.

If I look at 0,1,2,4,8 or 16 16 bits registers, the balance looks like this:

0: Pure stack based. While, e.g, there are FORTHs that on this on z80,
	it is not that performant.
	While FORTH is nice in theory, in practice it is hard to balance the stack.
	I tried this first but it didn't work well.
1: Accumulator only. While this is nice for NES,
	for z80 it underuses the registers. It is also difficult to program for.
2: Accumulator and operand. Slightly easier than 1 register, but the result in
	the accumulator has to be spilled immediately.
4: Medium amount of spilling, not too many fake registers.
	Z80 and 8086 more or less have this, so ideal for those CPUs.
8: Less spilling but more fake registers.
	4 registers is still a bit few for complex operations, but on most
	 bits machines the registers will be fake.
16: Would only be feasible on an Atmel,
	too many fake registers on other machines.

Having tried 0 and 2 registers, I found thses hard to use, so I decided
on a middle ground betweem 4 and 8: the virtual PIR CPU will have 4 virtual
16bits registers R1..R4. But it also has 4 16 bits virtual temporary storage
locations, T1..T4. These can only be stored to and loaded from.
There may be further limitations on how the 4 registers may be used.

As for the intructions, a 1 operand instruction works well on a
stack machine or 1 register machine, but less so on a 4 register machine.
While I tried 1 operand before, it also was fragile at times, for example for
data, etc. Therefore PIR will use a 0 ,1 or 2 operand format.

The 2 operand instructions will somewhat unusually use the left hand operand
as the source or modifier and the right hand operand as as the target. E.g:

	* MOV R1, R2 // R2 = R1 as R2 is the target.
	* JMP Z, FOO ; conditional jump as Z the modifier and FOO is the target.

*/

// This is only used by the parser to skip empty lines.
const SKIP Operation = -1

// This is only used by the parser to indicate the end of file.
const DONE Operation = -2

// Instructions consist of an opcode and two mandatory operands.
type Instruction struct {
	Operation // Operation
	Operands  []Operand
}

// Program is a list of instructions
type Program []Instruction
