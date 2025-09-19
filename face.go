package twophase

import "errors"

type face struct {
	f []int
}

func newFace() face {
	face := face{}
	for range 9 {
		face.f = append(face.f, cU)
	}
	for range 9 {
		face.f = append(face.f, cR)
	}
	for range 9 {
		face.f = append(face.f, cF)
	}
	for range 9 {
		face.f = append(face.f, cD)
	}
	for range 9 {
		face.f = append(face.f, cL)
	}
	for range 9 {
		face.f = append(face.f, cB)
	}

	return face
}

func (face *face) fromString(s string) error {
	if len(s) < 54 {
		return errors.New("less than 54 facelts")
	} else if len(s) > 54 {
		return errors.New("more than 54 facelts")
	}

	cnt := [6]int{0, 0, 0, 0, 0, 0}

	for i := range 54 {
		switch s[i] {
		case 'U':
			face.f[i] = cU
			cnt[cU] += 1
		case 'R':
			face.f[i] = cR
			cnt[cR] += 1
		case 'F':
			face.f[i] = cF
			cnt[cF] += 1
		case 'D':
			face.f[i] = cD
			cnt[cD] += 1
		case 'L':
			face.f[i] = cL
			cnt[cL] += 1
		case 'B':
			face.f[i] = cB
			cnt[cB] += 1
		}
	}

	for _, x := range cnt {
		if x != 9 {
			return errors.New("not exactly 9 facelts of each color")
		}
	}

	return nil
}

func (face face) toString() string {
	s := ""
	for i := range 54 {
		switch face.f[i] {
		case cU:
			s += "U"
		case cR:
			s += "R"
		case cF:
			s += "F"
		case cD:
			s += "D"
		case cL:
			s += "L"
		case cB:
			s += "B"
		}
	}
	return s
}

func (face *face) to_cubie_cube() cubie {
	cc := newCubieDefault()
	cc.cp = [8]int{-1, -1, -1, -1, -1, -1, -1, -1}
	cc.ep = [12]int{-1, -1, -1, -1, -1, -1, -1, -1 - 1, -1, -1, -1, -1}

	for i := range 8 {
		fac := cornerFaceLet[i]
		ori := 0
		for ori = range 3 {
			if face.f[fac[ori]] == cU || face.f[fac[ori]] == cD {
				break
			}
		}
		col1 := face.f[fac[(ori+1)%3]]
		col2 := face.f[fac[(ori+2)%3]]
		for j := range 8 {
			col := cornerColor[j]
			if col1 == col[1] && col2 == col[2] {
				cc.cp[i] = j
				cc.co[i] = ori
				break
			}
		}
	}

	for i := range 12 {
		for j := range 12 {
			if face.f[edgeFaceLet[i][0]] == edgeColor[j][0] && face.f[edgeFaceLet[i][1]] == edgeColor[j][1] {
				cc.ep[i] = j
				cc.eo[i] = 0
				break
			}
			if face.f[edgeFaceLet[i][0]] == edgeColor[j][1] && face.f[edgeFaceLet[i][1]] == edgeColor[j][0] {
				cc.ep[i] = j
				cc.eo[i] = 1
				break
			}
		}
	}

	return cc
}
