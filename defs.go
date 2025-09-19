package twophase

var cornerFaceLet = [8][3]int{
	{fU9, fR1, fF3},
	{fU7, fF1, fL3},
	{fU1, fL1, fB3},
	{fU3, fB1, fR3},
	{fD3, fF9, fR7},
	{fD1, fL9, fF7},
	{fD7, fB9, fL7},
	{fD9, fR9, fB7},
}

var cornerColor = [8][3]int{
	{cU, cR, cF},
	{cU, cF, cL},
	{cU, cL, cB},
	{cU, cB, cR},
	{cD, cF, cR},
	{cD, cL, cF},
	{cD, cB, cL},
	{cD, cR, cB},
}

var edgeFaceLet = [12][2]int{
	{fU6, fR2},
	{fU8, fF2},
	{fU4, fL2},
	{fU2, fB2},
	{fD6, fR8},
	{fD2, fF8},
	{fD4, fL8},
	{fD8, fB8},
	{fF6, fR4},
	{fF4, fL6},
	{fB6, fL4},
	{fB4, fR6},
}

var edgeColor = [12][2]int{
	{cU, cR},
	{cU, cF},
	{cU, cL},
	{cU, cB},
	{cD, cR},
	{cD, cF},
	{cD, cL},
	{cD, cB},
	{cF, cR},
	{cF, cL},
	{cB, cL},
	{cB, cR},
}
