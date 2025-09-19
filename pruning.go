package twophase

func getCornersUDEdgesDepth3(ix uint32, tables tables) uint32 {
	y := tables.corners_ud_edges_depth3[ix/16]
	y >>= (ix % 16) * 2
	return y & 3
}

func get_corners_ud_edges_depth3(ix uint32, fs_twist_depth3 *[]uint32) uint32 {
	y := (*fs_twist_depth3)[ix/16]
	y >>= (ix % 16) * 2
	return y & 3
}

func set_corners_ud_edges_depth3(ix uint32, value uint32, fs_twist_depth3 *[]uint32) {
	shift := (ix % 16) * 2
	base := ix >> 4

	(*fs_twist_depth3)[base] &= ^(3 << shift) & 0xffffffff
	(*fs_twist_depth3)[base] |= value << shift
}

func getFlipsliceTwistDepth3(ix uint32, tables tables) uint32 {
	y := tables.fs_twist_depth3[ix/16]
	y >>= (ix % 16) * 2
	return y & 3
}

func get_flipslice_twist_depth3(ix uint32, fs_twist_depth3 *[]uint32) uint32 {
	y := (*fs_twist_depth3)[ix/16]
	y >>= (ix % 16) * 2
	return y & 3
}

func set_flipslice_twist_depth3(ix uint32, value uint32, fs_twist_depth3 *[]uint32) {
	shift := (ix % 16) * 2
	base := ix >> 4

	(*fs_twist_depth3)[base] &= ^(3 << shift) & 0xffffffff
	(*fs_twist_depth3)[base] |= value << shift
}

func getDistance() []int8 {
	var distance []int8 = []int8{}
	for range 60 {
		distance = append(distance, 0)
	}
	for i := range 20 {
		for j := range 3 {
			distance[3*i+j] = int8((i/3)*3 + j)
			if i%3 == 2 && j == 0 {
				distance[3*i+j] += 3
			} else if i%3 == 0 && j == 2 {
				distance[3*i+j] -= 3
			}
		}
	}
	return distance
}
