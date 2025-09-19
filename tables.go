package twophase

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
)

type tables struct {
	corner_classidx []uint16
	corner_sym      []uint8
	corner_rep      []uint16

	fs_classidx []uint16
	fs_sym      []uint8
	fs_rep      []uint32

	conj_twist []uint16

	move_twist        []uint16
	move_flip         []uint16
	move_slice_sorted []uint16

	fs_twist_depth3 []uint32

	conj_ud_edges []uint16

	move_corners []uint16

	move_d_edges  []uint16
	move_u_edges  []uint16
	move_ud_edges []uint16

	corners_ud_edges_depth3 []uint32

	cornslice_depth []int8

	u_edges_plus_d_edges_to_ud_edges []uint16
}

// Generates if tables are not already in dirPath.
//
// This is slow, especially when generating, tables are in total almost 22 milliion elements in length
func GetAndGenerateTables(dirPath string) (tables, error) {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return tables{}, err
	}

	allTables := tables{}

	allTables.corner_classidx, allTables.corner_sym, allTables.corner_rep, err = getCornerTables(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.fs_classidx, allTables.fs_sym, allTables.fs_rep, err = getFlipSliceTables(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.conj_twist, err = getConjTwistTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_twist, err = getMoveTwistTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_flip, err = getMoveFlipTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_slice_sorted, err = getMoveSliceSortedTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.fs_twist_depth3, err = getPhase1PrunTable(dirPath, allTables)
	if err != nil {
		return tables{}, err
	}

	allTables.conj_ud_edges, err = getConjUDEdgesTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_corners, err = getMoveCornersTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_d_edges, err = getMoveDEdgesTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_u_edges, err = getMoveUEdgesTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.move_ud_edges, err = getMoveUDEdgesTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	allTables.corners_ud_edges_depth3, err = getPhase2PrunTable(dirPath, allTables)
	if err != nil {
		return tables{}, err
	}

	allTables.cornslice_depth, err = getPhase2CornspliceprunTable(dirPath, allTables)
	if err != nil {
		return tables{}, err
	}

	allTables.u_edges_plus_d_edges_to_ud_edges, err = getPhase2EdgemergeTable(dirPath)
	if err != nil {
		return tables{}, err
	}

	return allTables, nil
}

func getCornerTables(folderName string) ([]uint16, []uint8, []uint16, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()

	coClassidxExists, err := pathExists(filepath.Join(folderName, "co_classidx"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint16{}, err
	}

	coSymExists, err := pathExists(filepath.Join(folderName, "co_sym"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint16{}, err
	}

	coRepExists, err := pathExists(filepath.Join(folderName, "co_rep"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint16{}, err
	}

	corner_classidx := make([]uint16, 40_320)
	corner_sym := make([]uint8, 40_320)
	corner_rep := make([]uint16, 2768)

	if !(coClassidxExists && coSymExists && coRepExists) {
		for i := range 40_320 {
			corner_classidx[i] = 65535
		}

		var classidx uint16 = 0
		cc := newCubieDefault()

		for cp := range uint16(40320) {
			cc.setCorners(int(cp))

			if corner_classidx[cp] == 65535 {
				corner_classidx[cp] = classidx
				corner_sym[cp] = 0
				corner_rep[classidx] = cp
			} else {
				continue
			}

			for s := range uint8(16) {
				ss := newCubie(symCube[inv_idx[s]].cp, symCube[inv_idx[s]].co, symCube[inv_idx[s]].ep,
					symCube[inv_idx[s]].eo)
				ss.cornerMultiply(cc)
				ss.cornerMultiply(symCube[s])
				cp_new := ss.getCorners()
				if corner_classidx[cp_new] == 65535 {
					corner_classidx[cp_new] = classidx
					corner_sym[cp_new] = s
				}
			}
			classidx += 1
		}

		err = writeToFile(filepath.Join(folderName, "co_classidx"), corner_classidx)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
		err = writeToFile(filepath.Join(folderName, "co_sym"), corner_sym)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
		err = writeToFile(filepath.Join(folderName, "co_rep"), corner_rep)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "co_classidx"), corner_classidx)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
		err = readFromFile(filepath.Join(folderName, "co_sym"), corner_sym)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
		err = readFromFile(filepath.Join(folderName, "co_rep"), corner_rep)
		if err != nil {
			return []uint16{}, []uint8{}, []uint16{}, err
		}
	}
	return corner_classidx,
		corner_sym,
		corner_rep,
		nil
}

func getFlipSliceTables(folderName string) ([]uint16, []uint8, []uint32, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()

	fsClassidxExists, err := pathExists(filepath.Join(folderName, "fs_classidx"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint32{}, err
	}

	fsSymExists, err := pathExists(filepath.Join(folderName, "fs_sym"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint32{}, err
	}

	fsRepExists, err := pathExists(filepath.Join(folderName, "fs_rep"))
	if err != nil {
		return []uint16{}, []uint8{}, []uint32{}, err
	}

	flipslice_classidx := make([]uint16, 1_013_760)
	flipslice_sym := make([]uint8, 1_013_760)
	flipslice_rep := make([]uint32, 64_430)

	if !(fsClassidxExists && fsSymExists && fsRepExists) {
		for i := range 1_013_760 {
			flipslice_classidx[i] = 65535
		}

		var classidx uint16 = 0
		cc := newCubieDefault()

		for slc := range 495 {
			cc.setSlice(slc)
			for flip := range 2048 {
				cc.setFlip(flip)
				var idx uint32 = 2048*uint32(slc) + uint32(flip)
				if flipslice_classidx[idx] == 65535 {
					flipslice_classidx[idx] = classidx
					flipslice_sym[idx] = 0
					flipslice_rep[classidx] = idx
				} else {
					continue
				}
				for s := range uint8(16) {
					ss := newCubie(symCube[inv_idx[s]].cp, symCube[inv_idx[s]].co, symCube[inv_idx[s]].ep,
						symCube[inv_idx[s]].eo)
					ss.edgeMultiply(cc)
					ss.edgeMultiply(symCube[s])
					idx_new := 2048*ss.getSlice() + ss.getFlip()
					if flipslice_classidx[idx_new] == 65535 {
						flipslice_classidx[idx_new] = classidx
						flipslice_sym[idx_new] = s
					}
				}
				classidx += 1
			}
		}

		err = writeToFile(filepath.Join(folderName, "fs_classidx"), flipslice_classidx)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
		err = writeToFile(filepath.Join(folderName, "fs_sym"), flipslice_sym)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
		err = writeToFile(filepath.Join(folderName, "fs_rep"), flipslice_rep)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "fs_classidx"), flipslice_classidx)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
		err = readFromFile(filepath.Join(folderName, "fs_sym"), flipslice_sym)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
		err = readFromFile(filepath.Join(folderName, "fs_rep"), flipslice_rep)
		if err != nil {
			return []uint16{}, []uint8{}, []uint32{}, err
		}
	}

	return flipslice_classidx,
		flipslice_sym,
		flipslice_rep, nil
}

func getConjTwistTable(folderName string) ([]uint16, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()

	conjTwistExists, err := pathExists(filepath.Join(folderName, "conj_twist"))
	if err != nil {
		return []uint16{}, err
	}

	twist_conj := make([]uint16, 34992)

	if !conjTwistExists {
		for t := range 2187 {
			cc := newCubieDefault()
			cc.setTwist(t)
			for s := range 16 {
				ss := newCubie(symCube[s].cp, symCube[s].co, symCube[s].ep, symCube[s].eo)
				ss.cornerMultiply(cc)
				ss.cornerMultiply(symCube[inv_idx[s]])
				twist_conj[16*t+s] = uint16(ss.getTwist())
			}
		}

		err = writeToFile(filepath.Join(folderName, "conj_twist"), twist_conj)
		if err != nil {
			return []uint16{}, err
		}
	} else {

		err = readFromFile(filepath.Join(folderName, "conj_twist"), twist_conj)
		if err != nil {
			return []uint16{}, err
		}
	}

	return twist_conj, nil
}

func getMoveTwistTable(folderName string) ([]uint16, error) {
	var moveTwistExists bool
	basicMoveCube := getBasicMoveCube()

	moveTwistExists, err := pathExists(filepath.Join(folderName, "move_twist"))
	if err != nil {
		return []uint16{}, err
	}

	twist_move := make([]uint16, 39366)

	if !moveTwistExists {
		for i := range 2187 {
			a := newCubieDefault()
			a.setTwist(i)
			for j := range 6 {
				for k := range 3 {
					a.cornerMultiply(basicMoveCube[j])
					twist_move[18*i+3*j+k] = uint16(a.getTwist())
				}
				a.cornerMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_twist"), twist_move)
		if err != nil {
			return []uint16{}, err
		}
	} else {

		err = readFromFile(filepath.Join(folderName, "move_twist"), twist_move)
		if err != nil {
			return []uint16{}, err
		}

	}

	return twist_move, nil
}

func getMoveFlipTable(folderName string) ([]uint16, error) {
	var moveFlipExists bool
	basicMoveCube := getBasicMoveCube()

	moveFlipExists, err := pathExists(filepath.Join(folderName, "move_flip"))
	if err != nil {
		return []uint16{}, err
	}

	flip_move := make([]uint16, 36864)

	if !moveFlipExists {
		for i := range 2048 {
			a := newCubieDefault()
			a.setFlip(i)
			for j := range 6 {
				for k := range 3 {
					a.edgeMultiply(basicMoveCube[j])
					flip_move[18*i+3*j+k] = uint16(a.getFlip())
				}
				a.edgeMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_flip"), flip_move)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_flip"), flip_move)
		if err != nil {
			return []uint16{}, err
		}
	}

	return flip_move, nil
}

func getMoveSliceSortedTable(folderName string) ([]uint16, error) {
	basicMoveCube := getBasicMoveCube()

	moveSliceSortedExists, err := pathExists(filepath.Join(folderName, "move_slice_sorted"))
	if err != nil {
		return []uint16{}, err
	}

	slice_sorted_move := make([]uint16, 213_840)

	if !moveSliceSortedExists {
		for i := range 11880 {
			a := newCubieDefault()
			a.setSliceSorted(i)
			for j := range 6 {
				for k := range 3 {
					a.edgeMultiply(basicMoveCube[j])
					slice_sorted_move[18*i+3*j+k] = uint16(a.getSliceSorted())
				}
				a.edgeMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_slice_sorted"), slice_sorted_move)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_slice_sorted"), slice_sorted_move)
		if err != nil {
			return []uint16{}, err
		}
	}

	return slice_sorted_move, nil
}

func getPhase1PrunTable(folderName string, tables tables) ([]uint32, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()
	total := 140_908_410

	phase1PrunExists, err := pathExists(filepath.Join(folderName, "phase1_prun"))
	if err != nil {
		return []uint32{}, err
	}

	fs_twist_depth3 := make([]uint32, total/16+1)

	if !phase1PrunExists {
		for i := range total/16 + 1 {
			fs_twist_depth3[i] = 0xffffffff
		}

		cc := newCubieDefault()

		var fs_sym []uint16 = []uint16{}
		for range 64430 {
			fs_sym = append(fs_sym, 0)
		}
		for i := range 64430 {
			rep := tables.fs_rep[i]

			cc.setSlice(int(rep) / 2048)
			cc.setFlip(int(rep) % 2048)

			for s := range 16 {
				ss := newCubie(symCube[s].cp, symCube[s].co, symCube[s].ep,
					symCube[s].eo)
				ss.edgeMultiply(cc)
				ss.edgeMultiply(symCube[inv_idx[s]])
				if ss.getSlice() == int(rep)/2048 && ss.getFlip() == int(rep)%2048 {
					fs_sym[i] |= 1 << s
				}
			}
		}

		fs_classidx := 0
		twist := 0
		set_flipslice_twist_depth3(2187*uint32(fs_classidx)+uint32(twist), 0, &fs_twist_depth3)

		done := 1
		var depth uint32 = 0
		backsearch := false

		for done != total {
			depth3 := depth % 3
			if depth == 9 {
				backsearch = true
			}

			var idx uint32 = 0

			for fs_classidx := range 64430 {

				twist = 0
				for twist < 2187 {
					if !backsearch && idx%16 == 0 && fs_twist_depth3[idx/16] == 0xffffffff && twist < 2187-16 {
						twist += 16
						idx += 16
						continue
					}
					var match bool
					if backsearch {
						match = get_flipslice_twist_depth3(idx, &fs_twist_depth3) == 3
					} else {
						match = get_flipslice_twist_depth3(idx, &fs_twist_depth3) == depth3
					}
					if match {
						flipslice := tables.fs_rep[fs_classidx]
						flip := flipslice % 2048
						slice_ := flipslice >> 11
						for m := range 18 {
							twist1 := tables.move_twist[18*twist+m]
							flip1 := tables.move_flip[18*int(flip)+m]
							slice1 := tables.move_slice_sorted[432*int(slice_)+m] / 24

							flipslice1 := uint32((uint32(slice1) << 11) + uint32(flip1))

							fs1_classidx := tables.fs_classidx[flipslice1]
							fs1_sym := tables.fs_sym[flipslice1]
							twist1 = tables.conj_twist[(twist1<<4)+uint16(fs1_sym)]
							idx1 := 2187*uint32(fs1_classidx) + uint32(twist1)
							if !backsearch {
								if get_flipslice_twist_depth3(idx1, &fs_twist_depth3) == 3 {
									set_flipslice_twist_depth3(idx1, (depth+1)%3, &fs_twist_depth3)
									done += 1
									sym := fs_sym[fs1_classidx]
									if sym != 1 {
										for k := 1; k < 16; k++ {
											sym >>= 1
											if sym%2 == 1 {
												twist2 := tables.conj_twist[int(twist1<<4)+k]
												idx2 := 2187*uint32(fs1_classidx) + uint32(twist2)
												if get_flipslice_twist_depth3(idx2, &fs_twist_depth3) == 3 {
													set_flipslice_twist_depth3(idx2, (depth+1)%3, &fs_twist_depth3)
													done += 1
												}
											}
										}
									}
								}

							} else {
								if get_flipslice_twist_depth3(idx1, &fs_twist_depth3) == depth3 {
									set_flipslice_twist_depth3(idx, (depth+1)%3, &fs_twist_depth3)
									done += 1
									break
								}
							}
						}
					}
					twist += 1
					idx += 1
				}
			}
			depth += 1
		}

		err = writeToFile(filepath.Join(folderName, "phase1_prun"), fs_twist_depth3)
		if err != nil {
			return []uint32{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "phase1_prun"), fs_twist_depth3)
		if err != nil {
			return []uint32{}, err
		}
	}

	return fs_twist_depth3, nil
}

func getConjUDEdgesTable(folderName string) ([]uint16, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()

	conjUdEdgesExists, err := pathExists(filepath.Join(folderName, "conj_ud_edges"))
	if err != nil {
		return []uint16{}, err
	}

	table := make([]uint16, 645_120)

	if !conjUdEdgesExists {
		for t := range 40320 {
			cc := newCubieDefault()
			cc.setUDEdges(t)
			for s := range 16 {
				ss := newCubie(symCube[s].cp, symCube[s].co, symCube[s].ep, symCube[s].eo)
				ss.edgeMultiply(cc)
				ss.edgeMultiply(symCube[inv_idx[s]])
				table[16*t+s] = uint16(ss.getUDEdges())
			}
		}

		err = writeToFile(filepath.Join(folderName, "conj_ud_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "conj_ud_edges"), table)
		if err != nil {
			return []uint16{}, err
		}

	}

	return table, nil
}

func getMoveCornersTable(folderName string) ([]uint16, error) {
	basicMoveCube := getBasicMoveCube()

	moveCornersExists, err := pathExists(filepath.Join(folderName, "move_corners"))
	if err != nil {
		return []uint16{}, err
	}

	table := make([]uint16, 725_760)

	if !moveCornersExists {
		for i := range 40_320 {
			a := newCubieDefault()
			a.setCorners(i)
			for j := range 6 {
				for k := range 3 {
					a.cornerMultiply(basicMoveCube[j])
					table[18*i+3*j+k] = uint16(a.getCorners())
				}
				a.cornerMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_corners"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_corners"), table)
		if err != nil {
			return []uint16{}, err
		}
	}

	return table, nil
}

func getMoveDEdgesTable(folderName string) ([]uint16, error) {
	basicMoveCube := getBasicMoveCube()

	moveDEdgesExists, err := pathExists(filepath.Join(folderName, "move_d_edges"))
	if err != nil {
		return []uint16{}, err
	}

	table := make([]uint16, 213_840)

	if !moveDEdgesExists {
		for i := range 11880 {
			a := newCubieDefault()
			a.setDEdges(i)
			for j := range 6 {
				for k := range 3 {
					a.edgeMultiply(basicMoveCube[j])
					table[18*i+3*j+k] = uint16(a.getDEdges())
				}
				a.edgeMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_d_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_d_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	}

	return table, nil
}

func getMoveUEdgesTable(folderName string) ([]uint16, error) {
	basicMoveCube := getBasicMoveCube()

	moveUEdgesExists, err := pathExists(filepath.Join(folderName, "move_u_edges"))
	if err != nil {
		return []uint16{}, err
	}

	table := make([]uint16, 213_840)

	if !moveUEdgesExists {
		for i := range 11880 {
			a := newCubieDefault()
			a.setUEdges(i)
			for j := range 6 {
				for k := range 3 {
					a.edgeMultiply(basicMoveCube[j])
					table[18*i+3*j+k] = uint16(a.getUEdges())
				}
				a.edgeMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_u_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_u_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	}

	return table, nil
}

func getMoveUDEdgesTable(folderName string) ([]uint16, error) {
	basicMoveCube := getBasicMoveCube()

	moveUdEdgesExists, err := pathExists(filepath.Join(folderName, "move_ud_edges"))
	if err != nil {
		return []uint16{}, err
	}

	table := make([]uint16, 725_760)

	if !moveUdEdgesExists {
		for i := range 40320 {
			a := newCubieDefault()
			a.setUDEdges(i)
			for j := range 6 {
				for k := range 3 {
					a.edgeMultiply(basicMoveCube[j])
					if (j == 1 || j == 2 || j == 4 || j == 5) && k != 1 {
						continue
					}
					table[18*i+3*j+k] = uint16(a.getUDEdges())
				}
				a.edgeMultiply(basicMoveCube[j])
			}
		}

		err = writeToFile(filepath.Join(folderName, "move_ud_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "move_ud_edges"), table)
		if err != nil {
			return []uint16{}, err
		}
	}

	return table, nil
}

func getPhase2PrunTable(folderName string, tables tables) ([]uint32, error) {
	symCube := getSymCube()
	inv_idx := getInvIdx()
	total := 111_605_760

	phase2PrunExists, err := pathExists(filepath.Join(folderName, "phase2_prun"))
	if err != nil {
		return []uint32{}, err
	}

	corners_ud_edges_depth3 := make([]uint32, total/16)

	if !phase2PrunExists {
		for i := range total / 16 {
			corners_ud_edges_depth3[i] = 0xffffffff
		}
		cc := newCubieDefault()

		c_sym := []uint16{}
		for range 2768 {
			c_sym = append(c_sym, 0)
		}
		for i := range 2768 {
			rep := tables.corner_rep[i]
			cc.setCorners(int(rep))
			for s := range uint16(16) {
				ss := newCubie(symCube[s].cp, symCube[s].co, symCube[s].ep,
					symCube[s].eo)
				ss.cornerMultiply(cc)
				ss.cornerMultiply(symCube[inv_idx[s]])
				if ss.getCorners() == int(rep) {
					c_sym[i] |= 1 << s
				}
			}
		}

		c_classidx := 0
		ud_edge := 0
		set_corners_ud_edges_depth3(uint32(40320*c_classidx+ud_edge), 0, &corners_ud_edges_depth3)
		done := 1
		depth := uint32(0)
		for depth < 10 {
			depth3 := depth % 3
			idx := uint32(0)
			for c_classidx := range 2768 {
				ud_edge = 0

				for ud_edge < 40320 {
					if idx%16 == 0 && corners_ud_edges_depth3[idx/16] == 0xffffffff && ud_edge < 40320-16 {
						ud_edge += 16
						idx += 16
						continue
					}

					if get_corners_ud_edges_depth3(uint32(idx), &corners_ud_edges_depth3) == uint32(depth3) {
						corner := tables.corner_rep[c_classidx]

						for _, m := range []int{0, 1, 2, 4, 7, 9, 10, 11, 13, 16} {
							ud_edge1 := uint32(tables.move_ud_edges[18*ud_edge+m])
							corner1 := uint32(tables.move_corners[18*int(corner)+m])
							c1_classidx := uint32(tables.corner_classidx[corner1])
							c1_sym := uint32(tables.corner_sym[corner1])
							ud_edge1 = uint32(tables.conj_ud_edges[(uint32(ud_edge1)<<4)+uint32(c1_sym)])
							idx1 := 40320*uint32(c1_classidx) + ud_edge1

							if get_corners_ud_edges_depth3(uint32(idx1), &corners_ud_edges_depth3) == 3 {
								set_corners_ud_edges_depth3(uint32(idx1), uint32(depth+1)%3, &corners_ud_edges_depth3)
								done += 1
								sym := uint32(c_sym[c1_classidx])
								if sym != 1 {
									for k := 1; k < 16; k++ {
										sym >>= 1
										if sym%2 == 1 {
											ud_edge2 := tables.conj_ud_edges[uint32(uint32(ud_edge1)<<4)+uint32(k)]
											idx2 := 40320*c1_classidx + uint32(ud_edge2)
											if get_corners_ud_edges_depth3(uint32(idx2), &corners_ud_edges_depth3) == 3 {
												set_corners_ud_edges_depth3(uint32(idx2), uint32(depth+1)%3, &corners_ud_edges_depth3)
												done += 1
											}
										}
									}
								}
							}
						}
					}
					ud_edge += 1
					idx += 1
				}
			}
			depth++
		}

		err = writeToFile(filepath.Join(folderName, "phase2_prun"), corners_ud_edges_depth3)
		if err != nil {
			return []uint32{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "phase2_prun"), corners_ud_edges_depth3)
		if err != nil {
			return []uint32{}, err
		}
	}

	return corners_ud_edges_depth3, nil
}

func getPhase2CornspliceprunTable(folderName string, tables tables) ([]int8, error) {

	phase2Cornsliceprun, err := pathExists(filepath.Join(folderName, "phase2_cornsliceprun"))
	if err != nil {
		return []int8{}, err
	}

	cornslice_depth := make([]int8, 967_680)

	if !phase2Cornsliceprun {
		for i := range 967_680 {
			cornslice_depth[i] = -1
		}
		corners := 0
		slice_ := 0
		cornslice_depth[24*corners+slice_] = 0
		done := 1
		depth := 0
		for done != 967_680 {
			for corners := range 40320 {
				for slice_ := range 24 {
					if int(cornslice_depth[24*corners+slice_]) == depth {
						for _, m := range []int{0, 1, 2, 4, 7, 9, 10, 11, 13, 16} {
							corners1 := uint32(tables.move_corners[18*corners+m])
							slice_1 := uint32(tables.move_slice_sorted[18*slice_+m])
							idx1 := 24*corners1 + slice_1
							if cornslice_depth[idx1] == -1 {
								cornslice_depth[idx1] = int8(depth) + 1
								done += 1
							}
						}
					}
				}
			}
			depth += 1
		}

		err = writeToFile(filepath.Join(folderName, "phase2_cornsliceprun"), cornslice_depth)
		if err != nil {
			return []int8{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "phase2_cornsliceprun"), cornslice_depth)
		if err != nil {
			return []int8{}, err
		}

	}

	return cornslice_depth, nil
}

func getPhase2EdgemergeTable(folderName string) ([]uint16, error) {
	phase2EdgemergeExists, err := pathExists(filepath.Join(folderName, "phase2_edgemerge"))
	if err != nil {
		return []uint16{}, err
	}

	c_u := newCubieDefault()
	c_d := newCubieDefault()
	c_ud := newCubieDefault()

	edge_u := []int{edUR, edUF, edUL, edUB}
	edge_d := []int{edDR, edDF, edDL, edDB}
	edge_ud := []int{edUR, edUF, edUL, edUB, edDR, edDF, edDL, edDB}

	table := make([]uint16, 40320)

	if !phase2EdgemergeExists {
		for i := range 1680 {
			c_u.setUEdges(i)
			for j := range 70 {
				c_d.setDEdges(j * 24)
				invalid := false
				for _, e := range edge_ud {
					c_ud.ep[e] = -1
					if contains(edge_u, c_u.ep[e]) {
						c_ud.ep[e] = c_u.ep[e]
					}
					if contains(edge_d, c_d.ep[e]) {
						c_ud.ep[e] = c_d.ep[e]
					}
					if c_ud.ep[e] == -1 {
						invalid = true
						break
					}
				}
				if !invalid {
					for k := range 24 {
						c_d.setDEdges(j*24 + k)
						for _, e := range edge_ud {
							if contains(edge_u, c_u.ep[e]) {
								c_ud.ep[e] = c_u.ep[e]
							}
							if contains(edge_d, c_d.ep[e]) {
								c_ud.ep[e] = c_d.ep[e]
							}
						}
						table[24*i+k] = uint16(c_ud.getUDEdges())
					}
				}
			}
		}

		err = writeToFile(filepath.Join(folderName, "phase2_edgemerge"), table)
		if err != nil {
			return []uint16{}, err
		}
	} else {
		err = readFromFile(filepath.Join(folderName, "phase2_edgemerge"), table)
		if err != nil {
			return []uint16{}, err
		}
	}

	return table, nil
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func writeToFile(path string, data any) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0666)
}

func readFromFile(path string, dist any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	buf1 := bytes.NewBuffer(data)

	return binary.Read(buf1, binary.LittleEndian, dist)
}
