package twophase

func getBasicMoveCube() []cubie {
	basicMoveCube := []cubie{newCubieDefault(), newCubieDefault(), newCubieDefault(), newCubieDefault(), newCubieDefault(), newCubieDefault()}

	cpU := [8]int{coUBR, coURF, coUFL, coULB, coDFR, coDLF, coDBL, coDRB}
	coU := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	epU := [12]int{edUB, edUR, edUF, edUL, edDR, edDF, edDL, edDB, edFR, edFL, edBL, edBR}
	eoU := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	cpR := [8]int{coDFR, coUFL, coULB, coURF, coDRB, coDLF, coDBL, coUBR}
	coR := [8]int{2, 0, 0, 1, 1, 0, 0, 2}
	epR := [12]int{edFR, edUF, edUL, edUB, edBR, edDF, edDL, edDB, edDR, edFL, edBL, edUR}
	eoR := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	cpF := [8]int{coUFL, coDLF, coULB, coUBR, coURF, coDFR, coDBL, coDRB}
	coF := [8]int{1, 2, 0, 0, 2, 1, 0, 0}
	epF := [12]int{edUR, edFL, edUL, edUB, edDR, edFR, edDL, edDB, edUF, edDF, edBL, edBR}
	eoF := [12]int{0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0}

	cpD := [8]int{coURF, coUFL, coULB, coUBR, coDLF, coDBL, coDRB, coDFR}
	coD := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	epD := [12]int{edUR, edUF, edUL, edUB, edDF, edDL, edDB, edDR, edFR, edFL, edBL, edBR}
	eoD := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	cpL := [8]int{coURF, coULB, coDBL, coUBR, coDFR, coUFL, coDLF, coDRB}
	coL := [8]int{0, 1, 2, 0, 0, 2, 1, 0}
	epL := [12]int{edUR, edUF, edBL, edUB, edDR, edDF, edFL, edDB, edFR, edUL, edDL, edBR}
	eoL := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	cpB := [8]int{coURF, coUFL, coUBR, coDRB, coDFR, coDLF, coULB, coDBL}
	coB := [8]int{0, 0, 1, 2, 0, 0, 2, 1}
	epB := [12]int{edUR, edUF, edUL, edBR, edDR, edDF, edDL, edBL, edFR, edFL, edUB, edDB}
	eoB := [12]int{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1}

	basicMoveCube[cU] = newCubie(cpU, coU, epU, eoU)
	basicMoveCube[cR] = newCubie(cpR, coR, epR, eoR)
	basicMoveCube[cF] = newCubie(cpF, coF, epF, eoF)
	basicMoveCube[cD] = newCubie(cpD, coD, epD, eoD)
	basicMoveCube[cL] = newCubie(cpL, coL, epL, eoL)
	basicMoveCube[cB] = newCubie(cpB, coB, epB, eoB)

	return basicMoveCube
}
