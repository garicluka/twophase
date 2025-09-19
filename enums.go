package twophase

const (
	cU = iota
	cR
	cF
	cD
	cL
	cB
)

const (
	coURF = iota
	coUFL
	coULB
	coUBR
	coDFR
	coDLF
	coDBL
	coDRB
)

const (
	edUR = iota
	edUF
	edUL
	edUB
	edDR
	edDF
	edDL
	edDB
	edFR
	edFL
	edBL
	edBR
)

// ..5 are commented out because they are never used
const (
	fU1 = 0
	fU2 = 1
	fU3 = 2
	fU4 = 3
	// fU5 = 4
	fU6 = 5
	fU7 = 6
	fU8 = 7
	fU9 = 8
	fR1 = 9
	fR2 = 10
	fR3 = 11
	fR4 = 12
	// fR5 = 13
	fR6 = 14
	fR7 = 15
	fR8 = 16
	fR9 = 17
	fF1 = 18
	fF2 = 19
	fF3 = 20
	fF4 = 21
	// fF5 = 22
	fF6 = 23
	fF7 = 24
	fF8 = 25
	fF9 = 26
	fD1 = 27
	fD2 = 28
	fD3 = 29
	fD4 = 30
	// fD5 = 31
	fD6 = 32
	fD7 = 33
	fD8 = 34
	fD9 = 35
	fL1 = 36
	fL2 = 37
	fL3 = 38
	fL4 = 39
	// fL5 = 40
	fL6 = 41
	fL7 = 42
	fL8 = 43
	fL9 = 44
	fB1 = 45
	fB2 = 46
	fB3 = 47
	fB4 = 48
	// fB5 = 49
	fB6 = 50
	fB7 = 51
	fB8 = 52
	fB9 = 53
)

const (
	mU1 = 0
	mU2 = 1
	mU3 = 2
	mR1 = 3
	mR2 = 4
	mR3 = 5
	mF1 = 6
	mF2 = 7
	mF3 = 8
	mD1 = 9
	mD2 = 10
	mD3 = 11
	mL1 = 12
	mL2 = 13
	mL3 = 14
	mB1 = 15
	mB2 = 16
	mB3 = 17
)

const (
	bs_ROT_URF3 = 0
	bs_ROT_F2   = 1
	bs_ROT_U4   = 2
	bs_MIRR_LR2 = 3
)
