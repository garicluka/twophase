package twophase

import (
	"errors"
	"slices"
	"sync"
	"time"
)

type solverThread struct {
	cb_cube      cubie
	co_cube      *coord
	rot          uint32
	inv          uint32
	sofar_phase1 []uint32
	sofar_phase2 []uint32
	phase2_done  bool
	ret_length   int
	timeout      time.Duration
	start_time   time.Time

	cornersave uint32

	sharedState *sharedState
}

func newSolverThread(
	cb_cube cubie, rot uint32, inv uint32, ret_length int,
	timeout time.Duration, start_time time.Time, sharedState *sharedState) *solverThread {
	st := solverThread{
		cb_cube:      cb_cube,
		co_cube:      nil,
		rot:          rot,
		inv:          inv,
		sofar_phase1: nil,
		sofar_phase2: nil,
		phase2_done:  false,
		ret_length:   ret_length,
		timeout:      timeout,
		start_time:   start_time,

		cornersave: 0,

		sharedState: sharedState,
	}

	return &st
}

type sharedState struct {
	solutions       [][]uint32
	terminated      bool
	shortest_length []int
	mu              sync.Mutex
	count           int
}

func (st *solverThread) searchPhase2(corners uint32, ud_edges uint32, slice_sorted uint32, dist uint32, togo_phase2 uint32, tables tables) {
	st.sharedState.mu.Lock()
	terminated := st.sharedState.terminated
	st.sharedState.mu.Unlock()

	if terminated || st.phase2_done {
		return
	}

	conj_move := getConjMove()

	if togo_phase2 == 0 && slice_sorted == 0 {
		st.sharedState.mu.Lock()

		man := make([]uint32, len(st.sofar_phase1))
		copy(man, st.sofar_phase1)
		man = append(man, st.sofar_phase2...)

		if len(st.sharedState.solutions) == 0 || len((st.sharedState.solutions)[len(st.sharedState.solutions)-1]) > len(man) {
			if st.inv == 1 {
				slices.Reverse(man)
				for i := 0; i < len(man); i++ {
					m := man[i]
					man[i] = (m/3)*3 + (2 - m%3)
				}
			}
			for i := 0; i < len(man); i++ {
				m := man[i]
				man[i] = uint32(conj_move[18*16*st.rot+m])
			}
			st.sharedState.solutions = append(st.sharedState.solutions, man)
			(st.sharedState.shortest_length)[0] = len(man)
		}
		if (st.sharedState.shortest_length)[0] <= st.ret_length {
			st.sharedState.terminated = true
		}

		st.sharedState.mu.Unlock()
		st.phase2_done = true
	} else {
		distance := getDistance()
		for m := range uint32(18) {
			if m == mR1 || m == mR3 || m == mF1 || m == mF3 || m == mL1 || m == mL3 || m == mB1 || m == mB3 {
				continue
			}

			if len(st.sofar_phase2) > 0 {
				diff := (st.sofar_phase2)[len(st.sofar_phase2)-1]/3 - m/3
				if diff == 0 || diff == 3 {
					continue
				}
			} else {
				if len(st.sofar_phase1) > 0 {
					diff := (st.sofar_phase1)[len(st.sofar_phase1)-1]/3 - m/3
					if diff == 0 || diff == 3 {
						continue
					}
				}
			}

			corners_new := uint32(tables.move_corners[18*corners+m])
			ud_edges_new := uint32(tables.move_ud_edges[18*ud_edges+m])
			slice_sorted_new := uint32(tables.move_slice_sorted[18*slice_sorted+m])

			classidx := uint32(tables.corner_classidx[corners_new])
			sym := uint32(tables.corner_sym[corners_new])
			dist_new_mod3 := getCornersUDEdgesDepth3(
				40320*classidx+uint32(tables.conj_ud_edges[(ud_edges_new<<4)+sym]), tables)
			dist_new := uint32(distance[3*dist+dist_new_mod3])
			if max(dist_new, uint32(tables.cornslice_depth[24*corners_new+slice_sorted_new])) >= togo_phase2 {
				continue
			}

			st.sofar_phase2 = append(st.sofar_phase2, m)

			st.searchPhase2(corners_new, ud_edges_new, slice_sorted_new, dist_new, togo_phase2-1, tables)
			st.sofar_phase2 = st.sofar_phase2[:len(st.sofar_phase2)-1]
		}
	}
}

func (st *solverThread) search(flip uint32, twist uint32, slice_sorted uint32, dist uint32, togo_phase1 uint32, tables tables) {
	st.sharedState.mu.Lock()
	terminated := st.sharedState.terminated
	st.sharedState.mu.Unlock()

	if terminated {
		return
	}

	if togo_phase1 == 0 {
		st.sharedState.mu.Lock()
		if time.Since(st.start_time) > st.timeout && len(st.sharedState.solutions) > 0 {
			st.sharedState.terminated = true
		}
		st.sharedState.mu.Unlock()

		var m uint32
		if len(st.sofar_phase1) > 0 {
			m = (st.sofar_phase1)[len(st.sofar_phase1)-1]
		} else {
			m = fU1
		}

		var corners uint32
		if m == mR3 || m == mF3 || m == mL3 || m == mB3 {
			corners = uint32(tables.move_corners[18*st.cornersave+m-1])
		} else {
			corners = uint32(st.co_cube.corners)
			for _, m := range st.sofar_phase1 {
				corners = uint32(tables.move_corners[18*corners+m])
			}
			st.cornersave = corners
		}

		st.sharedState.mu.Lock()
		togo2_limit := min((st.sharedState.shortest_length)[0]-len(st.sofar_phase1), 11)
		st.sharedState.mu.Unlock()

		if int(tables.cornslice_depth[24*corners+slice_sorted]) >= togo2_limit {
			return
		}

		u_edges := uint32(st.co_cube.u_edges)
		d_edges := uint32(st.co_cube.d_edges)
		for _, m := range st.sofar_phase1 {
			u_edges = uint32(tables.move_u_edges[18*u_edges+m])
			d_edges = uint32(tables.move_d_edges[18*d_edges+m])
		}
		ud_edges := uint32(tables.u_edges_plus_d_edges_to_ud_edges[24*u_edges+d_edges%24])

		dist2 := st.co_cube.getDepthPhase2(corners, uint16(ud_edges), tables)

		for togo2 := int(dist2); togo2 < togo2_limit; togo2++ {
			st.sofar_phase2 = []uint32{}
			st.phase2_done = false

			st.searchPhase2(corners, ud_edges, slice_sorted, dist2, uint32(togo2), tables)
			if st.phase2_done {
				break
			}
		}
	} else {
		distance := getDistance()
		for m := range uint32(18) {
			if dist == 0 && togo_phase1 < 5 &&
				(m == mU1 || m == mU2 || m == mU3 || m == mR2 ||
					m == mF2 || m == mD1 || m == mD2 ||
					m == mD3 || m == mL2 || m == mB2) {
				continue
			}

			if len(st.sofar_phase1) > 0 {
				diff := st.sofar_phase1[len(st.sofar_phase1)-1]/3 - m/3
				if diff == 0 || diff == 3 {
					continue
				}
			}

			flip_new := uint32(tables.move_flip[18*flip+m])
			twist_new := uint32(tables.move_twist[18*twist+m])
			slice_sorted_new := uint32(tables.move_slice_sorted[18*slice_sorted+m])

			flipslice := 2048*(slice_sorted_new/24) + flip_new
			classidx := uint32(tables.fs_classidx[flipslice])
			sym := uint32(tables.fs_sym[flipslice])
			dist_new_mod3 := getFlipsliceTwistDepth3(2187*classidx+uint32(tables.conj_twist[(twist_new<<4)+sym]), tables)
			dist_new := uint32(distance[3*dist+dist_new_mod3])

			if dist_new >= togo_phase1 {
				continue
			}

			st.sofar_phase1 = append(st.sofar_phase1, uint32(m))

			st.search(flip_new, twist_new, slice_sorted_new, dist_new, togo_phase1-1, tables)
			st.sofar_phase1 = st.sofar_phase1[:len(st.sofar_phase1)-1]
		}
	}
}

func (st *solverThread) run(wg *sync.WaitGroup, tables tables) {
	defer wg.Done()
	symCube := getSymCube()
	var cb *cubie = nil
	switch st.rot {
	case 0:
		cb1 := newCubie(st.cb_cube.cp, st.cb_cube.co, st.cb_cube.ep, st.cb_cube.eo)
		cb = &cb1
	case 1:
		cb1 := newCubie(symCube[32].cp, symCube[32].co, symCube[32].ep, symCube[32].eo)
		cb = &cb1
		cb.multiply(st.cb_cube)
		cb.multiply(symCube[16])
	case 2:
		cb1 := newCubie(symCube[16].cp, symCube[16].co, symCube[16].ep, symCube[16].eo)
		cb = &cb1
		cb.multiply(st.cb_cube)
		cb.multiply(symCube[32])
	}

	if st.inv == 1 {
		tmp := newCubieDefault()
		cb.invCubieCube(&tmp)
		cb = &tmp
	}

	co1 := newCoord(cb, tables)
	st.co_cube = &co1

	dist := st.co_cube.getDepthPhase1(tables)

	for togo1 := dist; togo1 < 20; togo1++ {
		st.sofar_phase1 = []uint32{}

		st.search(uint32(st.co_cube.flip), uint32(st.co_cube.twist), uint32(st.co_cube.slice_sorted), dist, togo1, tables)
	}
}

func getUnique(a []uint8) []uint8 {
	var unique []uint8

	for _, v := range a {
		skip := slices.Contains(unique, v)
		if !skip {
			unique = append(unique, v)
		}
	}

	return unique
}

func getIntersection(a []uint8, b []uint8) []uint8 {
	set := make([]uint8, 0)
	hash := make(map[uint8]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

func filter(ss []uint32) (ret []uint32) {
	for _, s := range ss {
		if s < 3 {
			ret = append(ret, s)
		}
	}
	return
}

func numToMove(num uint32) (string, error) {
	switch num {
	case 0:
		return "U", nil
	case 1:
		return "U2", nil
	case 2:
		return "U'", nil
	case 3:
		return "R", nil
	case 4:
		return "R2", nil
	case 5:
		return "R'", nil
	case 6:
		return "F", nil
	case 7:
		return "F2", nil
	case 8:
		return "F'", nil
	case 9:
		return "D", nil
	case 10:
		return "D2", nil
	case 11:
		return "D'", nil
	case 12:
		return "L", nil
	case 13:
		return "L2", nil
	case 14:
		return "L'", nil
	case 15:
		return "B", nil
	case 16:
		return "B2", nil
	case 17:
		return "B'", nil
	default:
		return "", errors.New("Not a valid move")
	}
}

func newCubieFromFaceStr(faceStr string) (cubie, error) {
	f := newFace()
	err := f.fromString(faceStr)
	if err != nil {
		return cubie{}, err
	}
	c := f.to_cubie_cube()
	err = c.verify()

	if err != nil {
		return cubie{}, err
	}
	return c, nil
}

func solutionToStrMoves(solution []uint32) ([]string, error) {
	solutionStr := []string{}
	for _, m := range solution {
		move, err := numToMove(m)
		if err != nil {
			return []string{}, err
		}

		solutionStr = append(solutionStr, move)
	}

	return solutionStr, nil
}

func realSolve(cc cubie, max_length int, timeout time.Duration, tables tables) ([]uint32, error) {
	wg := sync.WaitGroup{}
	s_time := time.Now()
	syms := cc.symmetries()
	unique := getUnique(syms)
	intersection := getIntersection([]uint8{16, 20, 24, 28}, unique)

	var tr []uint32

	if len(intersection) > 0 {
		tr = []uint32{0, 3}
	} else {
		tr = []uint32{0, 1, 2, 3, 4, 5}
	}
	var set2 []uint8
	for i := 48; i < 96; i++ {
		set2 = append(set2, uint8(i))
	}
	intersection2 := getIntersection(set2, unique)
	if len(intersection2) > 0 {
		tr = filter(tr)
	}

	sharedState := sharedState{
		terminated:      false,
		solutions:       [][]uint32{},
		shortest_length: []int{999},
	}
	for _, i := range tr {
		th := newSolverThread(cc, i%3, i/3, max_length, timeout, s_time, &sharedState)
		wg.Add(1)
		go th.run(&wg, tables)
	}
	wg.Wait()

	if len(sharedState.solutions) == 0 {
		return []uint32{}, errors.New("len of solutions is 0")
	}

	return sharedState.solutions[len(sharedState.solutions)-1], nil
}

const SolvedCubeState string = "UUUUUUUUURRRRRRRRRFFFFFFFFFDDDDDDDDDLLLLLLLLLBBBBBBBBB"

// Verify that cubeState is possible, returns error if its not.
func VerifyCubeState(cubeState string) error {
	_, err := newCubieFromFaceStr(cubeState)

	return err
}

// Probability is same for all possible valid cubeStates.
func GetRandomCubeState() string {
	randomCubie := newCubieRandom()
	face := randomCubie.toFaceCube()

	return face.toString()
}

// Solve from fromCubeState to toCubeStats, return is []string solution.
//
// Max Length will stop solver early if solution of that or lower length is found.
//
// Timeout will stop solver after that duration, if no solution is found yet it will return when first one is found.
//
// CubeState format: "UUUUUUUUURRRRRRRRRFFFFFFFFFDDDDDDDDDLLLLLLLLLBBBBBBBBB".
//
// Solution is slice of moves, moves are 1 or 2 len strings that start with U, R, F, D, L or B, and optionally end with ' or 2.
//
// This function verifies both fromCubeState and toCubeState and returns error if they are not possible.
func Solve(fromCubeState string, toCubeState string, maxLength int, timeout time.Duration, tables tables) ([]string, error) {
	fromCubie, err := newCubieFromFaceStr(fromCubeState)
	if err != nil {
		return []string{}, err
	}
	if err := fromCubie.verify(); err != nil {
		return []string{}, err
	}

	toCubie, err := newCubieFromFaceStr(toCubeState)
	if err != nil {
		return []string{}, err
	}
	if err := toCubie.verify(); err != nil {
		return []string{}, err
	}

	cc := newCubieDefault()
	toCubie.invCubieCube(&cc)
	cc.multiply(fromCubie)
	err = cc.verify()
	if err != nil {
		return []string{}, err
	}

	solution, err := realSolve(cc, maxLength, timeout, tables)
	if err != nil {
		return []string{}, err
	}

	strSolution, err := solutionToStrMoves(solution)
	if err != nil {
		return []string{}, err
	}

	return strSolution, nil
}
