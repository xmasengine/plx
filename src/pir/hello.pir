// Set up VDP registers

LABL setup_vdp
PSHA VdpData
PSHW 4
COMP nz
JPIF setup_vdp

NAME VdpData
DATS `abcd`

