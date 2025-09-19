package twophase

func getConjMove() []uint16 {
	symCube := getSymCube()
	inv_idx := getInvIdx()
	conj_move := []uint16{}
	moveCube := getMoveCube()
	for range 48 * 48 {
		conj_move = append(conj_move, 0)
	}

	for s := range 48 {
		for m := range 18 {
			ss := newCubie(symCube[s].cp, symCube[s].co, symCube[s].ep, symCube[s].eo)
			ss.multiply(moveCube[m])
			ss.multiply(symCube[inv_idx[s]])
			for m2 := range 18 {
				if ss == moveCube[m2] {
					conj_move[18*s+m] = uint16(m2)
				}
			}
		}
	}
	return conj_move
}

func getBasicSymCube() []cubie {
	cpROT_URF3 := [8]int{coURF, coDFR, coDLF, coUFL, coUBR, coDRB, coDBL, coULB}
	coROT_URF3 := [8]int{1, 2, 1, 2, 2, 1, 2, 1}
	epROT_URF3 := [12]int{edUF, edFR, edDF, edFL, edUB, edBR, edDB, edBL, edUR, edDR, edDL, edUL}
	eoROT_URF3 := [12]int{1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1}

	cpROT_F2 := [8]int{coDLF, coDFR, coDRB, coDBL, coUFL, coURF, coUBR, coULB}
	coROT_F2 := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	epROT_F2 := [12]int{edDL, edDF, edDR, edDB, edUL, edUF, edUR, edUB, edFL, edFR, edBR, edBL}
	eoROT_F2 := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	cpROT_U4 := [8]int{coUBR, coURF, coUFL, coULB, coDRB, coDFR, coDLF, coDBL}
	coROT_U4 := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	epROT_U4 := [12]int{edUB, edUR, edUF, edUL, edDB, edDR, edDF, edDL, edBR, edFR, edFL, edBL}
	eoROT_U4 := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1}

	cpMIRR_LR2 := [8]int{coUFL, coURF, coUBR, coULB, coDLF, coDFR, coDRB, coDBL}
	coMIRR_LR2 := [8]int{3, 3, 3, 3, 3, 3, 3, 3}
	epMIRR_LR2 := [12]int{edUL, edUF, edUR, edUB, edDL, edDF, edDR, edDB, edFL, edFR, edBR, edBL}
	eoMIRR_LR2 := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	basicSymCube := []cubie{newCubieDefault(), newCubieDefault(), newCubieDefault(), newCubieDefault()}
	basicSymCube[bs_ROT_URF3] = newCubie(cpROT_URF3, coROT_URF3, epROT_URF3, eoROT_URF3)
	basicSymCube[bs_ROT_F2] = newCubie(cpROT_F2, coROT_F2, epROT_F2, eoROT_F2)
	basicSymCube[bs_ROT_U4] = newCubie(cpROT_U4, coROT_U4, epROT_U4, eoROT_U4)
	basicSymCube[bs_MIRR_LR2] = newCubie(cpMIRR_LR2, coMIRR_LR2, epMIRR_LR2, eoMIRR_LR2)

	return basicSymCube
}

func getSymCube() []cubie {
	basicSymCube := getBasicSymCube()

	symCube := []cubie{}
	cc := newCubieDefault()
	idx := 0
	for range 3 {
		for range 2 {
			for range 4 {
				for range 2 {
					symCube = append(symCube, newCubie(cc.cp, cc.co, cc.ep, cc.eo))
					idx += 1
					cc.multiply(basicSymCube[bs_MIRR_LR2])
				}
				cc.multiply(basicSymCube[bs_ROT_U4])
			}
			cc.multiply(basicSymCube[bs_ROT_F2])
		}
		cc.multiply(basicSymCube[bs_ROT_URF3])
	}

	return symCube
}

func getInvIdx() [48]uint8 {
	symCube := getSymCube()
	inv_idx := [48]uint8{}

	for j := range 48 {
		for i := range 48 {
			cc := newCubie(symCube[j].cp, symCube[j].co, symCube[j].ep, symCube[j].eo)
			cc.cornerMultiply(symCube[i])
			if cc.cp[coURF] == coURF && cc.cp[coUFL] == coUFL && cc.cp[coULB] == coULB {
				inv_idx[j] = uint8(i)
				break
			}
		}
	}

	return inv_idx
}
