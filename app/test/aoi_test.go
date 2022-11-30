package test

import (
	"fmt"
	"server/app/core"
	"testing"
)

func TestSurrounding(t *testing.T) {
	aoi := core.NewAoiManager(0, 0, 250, 5, 0, 250, 5)
	for gid, _ := range aoi.Grids {
		grids := aoi.GetSurroundingByGId(gid)
		fmt.Printf("%d %d\n", gid%5, gid/5)
		for _, g := range grids {
			fmt.Printf("%d %d\n", g.GId%5, g.GId/5)
		}
		fmt.Println()
	}
}
