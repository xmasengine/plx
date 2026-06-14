package pir

/*

PIR -> z80 translator conventions:

We use the hardware stack as the virtual stack.

The register pairs HL, DE are used as the top (TOS) of stack and next
of stack (NXT), and over if these are used we spill to the stack.

The calling convetion for funtions and operator is the same,
they use the stack, but HL and DE are used for the first two parameters,
because the are TOS and NXT.

The A and F registers will be used by all operators so are always clobbered.
The BC register is free for other uses.

*/
