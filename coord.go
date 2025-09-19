package twophase

import "slices"

type coord struct {
	twist        int
	flip         int
	slice_sorted int
	u_edges      int
	d_edges      int
	corners      int
	ud_edges     int

	flipslice_classidx uint16
	flipslice_sym      uint8
	flipslice_rep      uint32
	corner_classidx    uint16
	corner_sym         uint8
	corner_rep         uint16
}

func newCoord(cc *cubie, tables tables) coord {
	var co coord
	if cc == nil {
		co.u_edges = 1656
	} else {
		co.twist = cc.getTwist()
		co.flip = cc.getFlip()
		co.slice_sorted = cc.getSliceSorted()
		co.u_edges = cc.getUEdges()
		co.d_edges = cc.getDEdges()
		co.corners = cc.getCorners()
		if co.slice_sorted < 24 {
			co.ud_edges = cc.getUDEdges()
		} else {
			co.ud_edges = -1
		}

		co.flipslice_classidx = tables.fs_classidx[2048*(co.slice_sorted/24)+co.flip]
		co.flipslice_sym = tables.fs_sym[2048*(co.slice_sorted/24)+co.flip]
		co.flipslice_rep = tables.fs_rep[co.flipslice_classidx]

		co.corner_classidx = tables.corner_classidx[co.corners]
		co.corner_sym = tables.corner_sym[co.corners]
		co.corner_rep = tables.corner_rep[co.corner_classidx]
	}

	return co
}

func (co *coord) getDepthPhase2(corners uint32, ud_edges uint16, tables tables) uint32 {
	classidx := tables.corner_classidx[corners]
	sym := tables.corner_sym[corners]
	depth_mod3 := getCornersUDEdgesDepth3(40320*uint32(classidx)+uint32(tables.conj_ud_edges[(uint32(ud_edges)<<4)+uint32(sym)]), tables)
	if depth_mod3 == 3 {
		return 11
	}
	depth := uint32(0)
	for corners != 0 || ud_edges != 0 {
		if depth_mod3 == 0 {
			depth_mod3 = 3
		}
		for _, m := range []int{0, 1, 2, 4, 7, 9, 10, 11, 13, 16} {
			corners1 := uint32(tables.move_corners[18*(corners)+uint32(m)])
			ud_edges1 := tables.move_ud_edges[18*uint32(ud_edges)+uint32(m)]
			classidx1 := tables.corner_classidx[corners1]
			sym = tables.corner_sym[corners1]
			if getCornersUDEdgesDepth3(
				40320*uint32(classidx1)+uint32(tables.conj_ud_edges[(uint32(ud_edges1)<<4)+uint32(sym)]), tables) == depth_mod3-1 {
				depth += 1
				corners = corners1
				ud_edges = ud_edges1
				depth_mod3 -= 1
				break
			}
		}
	}
	return depth
}

func (co *coord) getDepthPhase1(tables tables) uint32 {
	slice_ := uint32(co.slice_sorted / 24)
	flip := uint32(co.flip)
	twist := uint32(co.twist)
	flipslice := 2048*slice_ + flip
	classidx := tables.fs_classidx[flipslice]
	sym := tables.fs_sym[flipslice]
	depth_mod3 := getFlipsliceTwistDepth3(2187*uint32(classidx)+uint32(tables.conj_twist[(twist<<4)+uint32(sym)]), tables)

	var depth uint32 = 0
	for flip != 0 || slice_ != 0 || twist != 0 {
		if depth_mod3 == 0 {
			depth_mod3 = 3
		}

		for m := range uint32(18) {
			twist1 := tables.move_twist[18*twist+m]
			flip1 := tables.move_flip[18*flip+m]
			slice1 := tables.move_slice_sorted[18*slice_*24+m] / 24
			flipslice1 := 2048*uint32(slice1) + uint32(flip1)
			classidx1 := tables.fs_classidx[flipslice1]
			sym := tables.fs_sym[flipslice1]
			if getFlipsliceTwistDepth3(
				2187*uint32(classidx1)+uint32(tables.conj_twist[(uint32(twist1)<<4)+uint32(sym)]), tables) == depth_mod3-1 {
				depth += 1
				twist = uint32(twist1)
				flip = uint32(flip1)
				slice_ = uint32(slice1)
				depth_mod3 -= 1
				break
			}
		}
	}
	return depth
}

func contains(s []int, e int) bool {
	return slices.Contains(s, e)
}
