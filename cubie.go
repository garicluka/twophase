package twophase

import (
	"errors"
	"math/rand/v2"
	"reflect"
)

type cubie struct {
	cp [8]int
	co [8]int
	ep [12]int
	eo [12]int
}

func newCubieDefault() cubie {
	cp := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	co := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	ep := [12]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	eo := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	c := cubie{
		cp: cp,
		co: co,
		ep: ep,
		eo: eo,
	}

	return c
}

func newCubie(cp [8]int, co [8]int, ep [12]int, eo [12]int) cubie {
	c := cubie{
		cp: cp,
		co: co,
		ep: ep,
		eo: eo,
	}

	return c
}

func (cc *cubie) set_edges(idx int) {
	ep := [12]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	cc.ep = ep
	for j := range 12 {
		k := idx % (j + 1)
		idx /= j + 1
		for k > 0 {
			rotate_right_12(&cc.ep, 0, j)
			k -= 1
		}
	}
}

func newCubieRandom() cubie {
	cc := newCubieDefault()
	cc.set_edges(rand.IntN(479001600))
	p := cc.edgeParity()
	for {
		cc.setCorners(rand.IntN(40320))
		if p == cc.cornerParity() {
			break
		}
	}
	cc.setFlip(rand.IntN(2048))
	cc.setTwist(rand.IntN(2187))

	return cc
}

func (cubie *cubie) toFaceCube() face {
	fc := newFace()
	for i := range 8 {
		j := cubie.cp[i]
		ori := cubie.co[i]
		for k := range 3 {
			fc.f[cornerFaceLet[i][(k+ori)%3]] = cornerColor[j][k]
		}

	}
	for i := range 12 {
		j := cubie.ep[i]
		ori := cubie.eo[i]
		for k := range 2 {
			fc.f[edgeFaceLet[i][(k+ori)%2]] = edgeColor[j][k]
		}
	}

	return fc
}

func (cubie *cubie) multiply(b cubie) {
	cubie.cornerMultiply(b)
	cubie.edgeMultiply(b)
}

func getMoveCube() []cubie {
	basicMoveCube := getBasicMoveCube()
	moveCube := []cubie{}
	for range 18 {
		moveCube = append(moveCube, newCubieDefault())
	}
	for c1 := range 6 {
		cc := newCubieDefault()
		for k1 := range 3 {
			cc.multiply(basicMoveCube[c1])
			moveCube[3*c1+k1] = newCubie(cc.cp, cc.co, cc.ep, cc.eo)
		}
	}
	return moveCube
}

func (cubie *cubie) cornerMultiply(b cubie) {
	c_perm := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	c_ori := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	ori := 0

	for c := range 8 {
		c_perm[c] = cubie.cp[b.cp[c]]
		ori_a := cubie.co[b.cp[c]]
		ori_b := b.co[c]
		if ori_a < 3 && ori_b < 3 {
			ori = ori_a + ori_b
			if ori >= 3 {
				ori -= 3
			}
		} else if ori_a < 3 && 3 <= ori_b {
			ori = ori_a + ori_b
			if ori >= 6 {
				ori -= 3
			}
		} else if ori_a >= 3 && 3 > ori_b {
			ori = ori_a - ori_b
			if ori < 3 {
				ori += 3
			}

		} else if ori_a >= 3 && ori_b >= 3 {
			ori = ori_a - ori_b
			if ori < 0 {
				ori += 3
			}
		}
		c_ori[c] = ori
	}

	for c := range 8 {
		cubie.cp[c] = c_perm[c]
		cubie.co[c] = c_ori[c]
	}
}

func (cubie *cubie) edgeMultiply(b cubie) {
	e_perm := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	e_ori := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for e := range 12 {
		e_perm[e] = cubie.ep[b.ep[e]]
		e_ori[e] = (b.eo[e] + cubie.eo[b.ep[e]]) % 2
	}

	for e := range 12 {
		cubie.ep[e] = e_perm[e]
		cubie.eo[e] = e_ori[e]
	}
}

func (cubie cubie) invCubieCube(d *cubie) {
	for e := range 12 {
		d.ep[cubie.ep[e]] = e

	}
	for e := range 12 {
		d.eo[e] = cubie.eo[d.ep[e]]
	}

	for c := range 8 {
		d.cp[cubie.cp[c]] = c
	}
	for c := range 8 {
		ori := cubie.co[d.cp[c]]
		if ori >= 3 {
			d.co[c] = ori
		} else {
			d.co[c] = -ori
			if d.co[c] < 0 {
				d.co[c] += 3
			}
		}
	}
}

func (cubie *cubie) getTwist() int {
	ret := 0
	for i := range coDRB {
		ret = 3*ret + cubie.co[i]
	}
	return ret
}

func (cubie *cubie) setTwist(twist int) {
	twistparity := 0

	for i := coDRB - 1; i > coURF-1; i-- {
		cubie.co[i] = twist % 3
		twistparity += cubie.co[i]
		twist /= 3
	}

	cubie.co[coDRB] = ((3 - twistparity%3) % 3)
}

func (cubie *cubie) getFlip() int {
	ret := 0
	for i := range edBR {
		ret = 2*ret + cubie.eo[i]
	}
	return ret

}

func (cubie *cubie) getSliceSorted() int {
	a := 0
	x := 0
	edge4 := []int{0, 0, 0, 0}

	for j := edBR; j > edUR-1; j-- {
		if edFR <= cubie.ep[j] && cubie.ep[j] <= edBR {
			a += c_nk(11-j, x+1)
			edge4[3-x] = cubie.ep[j]
			x += 1
		}
	}
	b := 0
	for j := 3; j > 0; j-- {
		k := 0
		for edge4[j] != j+8 {
			rotate_left(edge4, 0, j)
			k += 1
		}
		b = (j+1)*b + k
	}

	return 24*a + b
}

func (cubie *cubie) getUEdges() int {
	a := 0
	x := 0
	edge4 := []int{0, 0, 0, 0}
	ep_mod := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(ep_mod, cubie.ep[:])

	for range 4 {
		rotate_right(ep_mod, 0, 11)
	}
	for j := edBR; j > edUR-1; j-- {
		if edUR <= ep_mod[j] && ep_mod[j] <= edUB {
			a += c_nk(11-j, x+1)
			edge4[3-x] = ep_mod[j]
			x += 1
		}
	}
	b := 0
	for j := 3; j > 0; j-- {
		k := 0
		for edge4[j] != j {
			rotate_left(edge4, 0, j)
			k += 1
		}
		b = (j+1)*b + k
	}
	return 24*a + b
}

func (cubie *cubie) getDEdges() int {
	a := 0
	x := 0
	edge4 := []int{0, 0, 0, 0}
	ep_mod := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	copy(ep_mod, cubie.ep[:])
	for range 4 {
		rotate_right(ep_mod, 0, 11)
	}
	for j := edBR; j > edUR-1; j-- {
		if edDR <= ep_mod[j] && ep_mod[j] <= edDB {
			a += c_nk(11-j, x+1)
			edge4[3-x] = ep_mod[j]
			x += 1
		}
	}
	b := 0

	for j := 3; j > 0; j-- {
		k := 0
		for edge4[j] != j+4 {
			rotate_left(edge4, 0, j)
			k += 1
		}
		b = (j+1)*b + k
	}
	return 24*a + b
}

func (cubie *cubie) getCorners() int {
	var perm []int = []int{0, 0, 0, 0, 0, 0, 0, 0}
	copy(perm, cubie.cp[:])

	b := 0

	for j := coDRB; j > coURF; j-- {
		k := 0
		for perm[j] != j {
			rotate_left(perm, 0, j)
			k += 1
		}
		b = (j+1)*b + k
	}

	return b
}
func (cubie *cubie) setUEdges(idx int) {
	slice_edge := []uint16{edUR, edUF, edUL, edUB}
	other_edge := []uint16{edDR, edDF, edDL, edDB, edFR, edFL, edBL, edBR}
	b := idx % 24
	a := idx / 24
	for e := range 12 {
		cubie.ep[e] = -1
	}

	j := 1
	for j < 4 {
		k := b % (j + 1)
		b /= j + 1
		for k > 0 {
			rotate_right(slice_edge, 0, j)
			k -= 1
		}
		j += 1

	}

	x := 4
	for j := range 12 {
		if a-c_nk(11-j, x) >= 0 {
			cubie.ep[j] = int(slice_edge[4-x])
			a -= c_nk(11-j, x)
			x -= 1
		}
	}

	x = 0
	for j := range 12 {
		if cubie.ep[j] == -1 {
			cubie.ep[j] = int(other_edge[x])
			x += 1
		}
	}
	for range 4 {
		rotate_left_12(&cubie.ep, 0, 11)
	}
}

func (cubie *cubie) setDEdges(idx int) {
	slice_edge := []uint16{edDR, edDF, edDL, edDB}
	other_edge := []uint16{edFR, edFL, edBL, edBR, edUR, edUF, edUL, edUB}
	b := idx % 24
	a := idx / 24
	for e := range 12 {
		cubie.ep[e] = -1
	}
	j := 1
	for j < 4 {
		k := b % (j + 1)
		b /= j + 1
		for k > 0 {
			rotate_right(slice_edge, 0, j)
			k -= 1
		}
		j += 1

	}

	x := 4
	for j := range 12 {
		if a-c_nk(11-j, x) >= 0 {
			cubie.ep[j] = int(slice_edge[4-x])
			a -= c_nk(11-j, x)
			x -= 1
		}
	}

	x = 0
	for j := range 12 {
		if cubie.ep[j] == -1 {
			cubie.ep[j] = int(other_edge[x])
			x += 1
		}
	}
	for range 4 {
		rotate_left_12(&cubie.ep, 0, 11)
	}
}

func (cubie *cubie) setCorners(idx int) {
	cubie.cp = [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	for j := range 8 {
		k := idx % (j + 1)
		idx /= j + 1
		for k > 0 {
			rotate_right_8(&cubie.cp, 0, j)
			k -= 1
		}
	}
}

func (cubie *cubie) setUDEdges(idx int) {
	for i := range 8 {
		cubie.ep[i] = i
	}

	for j := range 8 {
		k := idx % (j + 1)
		idx /= j + 1
		for k > 0 {
			rotate_right_12(&cubie.ep, 0, j)
			k -= 1
		}
	}
}

func (cubie *cubie) getUDEdges() int {
	var perm []int = []int{0, 0, 0, 0, 0, 0, 0, 0}
	copy(perm, cubie.ep[0:8])

	b := 0

	for j := edDB; j > edUR; j-- {
		k := 0
		for perm[j] != j {
			rotate_left(perm, 0, j)
			k += 1
		}
		b = (j+1)*b + k
	}

	return b

}

func (cubie *cubie) setSliceSorted(idx int) {
	slice_edge := []int{edFR, edFL, edBL, edBR}
	other_edge := []int{edUR, edUF, edUL, edUB, edDR, edDF, edDL, edDB}
	b := idx % 24
	a := idx / 24
	for e := range 12 {
		cubie.ep[e] = -1
	}
	j := 1
	for j < 4 {
		k := b % (j + 1)
		b /= j + 1
		for k > 0 {
			rotate_right(slice_edge, 0, j)
			k -= 1
		}
		j += 1
	}
	x := 4
	for j := range 12 {
		if a-c_nk(11-j, x) >= 0 {
			cubie.ep[j] = slice_edge[4-x]
			a -= c_nk(11-j, x)
			x -= 1
		}
	}
	x = 0
	for j := range 12 {
		if cubie.ep[j] == -1 {
			cubie.ep[j] = other_edge[x]
			x += 1
		}
	}
}

func (cubie *cubie) setSlice(idx int) {
	slice_edge := []int{edFR, edFL, edBL, edBR}
	other_edge := []int{edUR, edUF, edUL, edUB, edDR, edDF, edDL, edDB}
	a := idx
	for e := range 12 {
		cubie.ep[e] = -1
	}
	x := 4
	for j := range 12 {
		if a-c_nk(11-j, x) >= 0 {
			cubie.ep[j] = slice_edge[4-x]
			a -= c_nk(11-j, x)
			x -= 1
		}
	}

	x = 0
	for j := range 12 {
		if cubie.ep[j] == -1 {
			cubie.ep[j] = other_edge[x]
			x += 1
		}
	}
}

func (cubie *cubie) setFlip(flip int) {
	flipparity := 0
	for i := edBR - 1; i > edUR-1; i-- {
		cubie.eo[i] = flip % 2
		flipparity += cubie.eo[i]
		flip /= 2
	}
	cubie.eo[edBR] = ((2 - flipparity%2) % 2)
}

func (cubie *cubie) getSlice() int {
	a := 0
	x := 0
	for j := edBR; j > edUR-1; j-- {
		if edFR <= cubie.ep[j] && cubie.ep[j] <= edBR {
			a += c_nk(11-j, x+1)
			x += 1
		}
	}
	return a
}

func (cubie *cubie) edgeParity() int {
	s := 0

	for i := edBR; i > edUR; i-- {
		for j := i - 1; j > edUR-1; j-- {
			if cubie.ep[j] > cubie.ep[i] {
				s += 1
			}
		}
	}

	return s % 2
}

func (cubie *cubie) cornerParity() int {
	s := 0

	for i := coDRB; i > coURF; i-- {
		for j := i - 1; j > coURF-1; j-- {
			if cubie.cp[j] > cubie.cp[i] {
				s += 1
			}
		}
	}

	return s % 2
}

func (cubie *cubie) verify() error {
	for i := range 8 {
		if cubie.co[i] == -1 {
			return errors.New("Error: Some corner orientations are invalid.")
		}
	}
	for i := range 8 {
		if cubie.cp[i] == -1 {
			return errors.New("Error: Some corner permutations are invalid.")
		}
	}
	for i := range 12 {
		if cubie.eo[i] == -1 {
			return errors.New("Error: Some edge orientations are invalid.")
		}
	}
	for i := range 12 {
		if cubie.ep[i] == -1 {
			return errors.New("Error: Some edge permutations are invalid.")
		}
	}

	edge_count := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := range 12 {
		edge_count[cubie.ep[i]] += 1
	}
	for i := range 12 {
		if edge_count[i] != 1 {
			return errors.New("Error: Some edges are undefined.")
		}
	}

	s := 0
	for i := range 12 {
		s += cubie.eo[i]

	}
	if s%2 != 0 {
		return errors.New("Error: Total edge flip is wrong.")
	}

	corner_count := [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	for i := range 8 {
		corner_count[cubie.cp[i]] += 1
	}
	for i := range 8 {
		if corner_count[i] != 1 {
			return errors.New("Error: Some corners are undefined.")
		}
	}

	s2 := 0
	for i := range 8 {
		s2 += cubie.co[i]
	}
	if s2%3 != 0 {
		return errors.New("Error: Total corner twist is wrong.")
	}

	if cubie.edgeParity() != cubie.cornerParity() {
		return errors.New("Error: Wrong edge and corner parity")
	}

	return nil
}

func (cubie *cubie) symmetries() []uint8 {
	symCube := getSymCube()
	inv_idx := getInvIdx()

	s := []uint8{}
	d := newCubieDefault()
	for j := range 48 {
		c := newCubie(symCube[j].cp, symCube[j].co, symCube[j].ep, symCube[j].eo)
		c.multiply(*cubie)
		c.multiply(symCube[inv_idx[j]])
		if reflect.DeepEqual(*cubie, c) {
			s = append(s, uint8(j))
		}
		c.invCubieCube(&d)
		if reflect.DeepEqual(*cubie, d) {
			s = append(s, uint8(j+48))
		}
	}

	return s
}
