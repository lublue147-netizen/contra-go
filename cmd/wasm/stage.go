package main

import (
	"math"
)

type Stage struct {
	Length      float64
	Height      float64
	ViewWidth   float64
	ViewHeight  float64
	CameraX     float64
	BossArenaX  float64
	Platforms   []Platform
	WaterLevel  float64
	WaveOffset  float64
}

func NewStage1() *Stage {
	st := &Stage{
		Length:     2560, // Total stage length
		Height:     224,  // NES native vertical resolution
		ViewWidth:  256,  // NES native horizontal resolution
		ViewHeight: 224,
		CameraX:    0,
		BossArenaX: 2200, // Where the camera locks for the boss
		WaterLevel: 210,
	}
	st.InitPlatforms()
	return st
}

func (st *Stage) InitPlatforms() {
	platforms := []Platform{
		// 1. Initial beach & jungle
		{X: 0, Y: 180, W: 450, H: 44, IsWater: false, IsDropThru: false},
		{X: 150, Y: 120, W: 140, H: 12, IsWater: false, IsDropThru: true},
		{X: 260, Y: 70, W: 120, H: 12, IsWater: false, IsDropThru: true},

		// 2. First bridge over water
		{X: 450, Y: 180, W: 150, H: 12, IsWater: false, IsDropThru: true},
		{X: 450, Y: 206, W: 150, H: 18, IsWater: true, IsDropThru: false}, // Water stream

		// 3. Central Jungle Ridge
		{X: 600, Y: 180, W: 380, H: 44, IsWater: false, IsDropThru: false},
		{X: 660, Y: 125, W: 150, H: 12, IsWater: false, IsDropThru: true},
		{X: 780, Y: 80, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 4. Waterfall & stepped cliffs
		{X: 980, Y: 160, W: 200, H: 64, IsWater: false, IsDropThru: false},
		{X: 1040, Y: 105, W: 180, H: 12, IsWater: false, IsDropThru: true},
		{X: 1180, Y: 150, W: 240, H: 74, IsWater: false, IsDropThru: false},
		{X: 1240, Y: 90, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 5. Second bridge & deep lagoon
		{X: 1420, Y: 180, W: 180, H: 12, IsWater: false, IsDropThru: true},
		{X: 1420, Y: 206, W: 180, H: 18, IsWater: true, IsDropThru: false},

		// 6. Mountain approach leading to Fortress
		{X: 1600, Y: 180, W: 350, H: 44, IsWater: false, IsDropThru: false},
		{X: 1680, Y: 120, W: 180, H: 12, IsWater: false, IsDropThru: true},
		{X: 1820, Y: 75, W: 160, H: 12, IsWater: false, IsDropThru: true},
		{X: 1950, Y: 180, W: 250, H: 44, IsWater: false, IsDropThru: false},
		{X: 2000, Y: 125, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 7. Boss Arena Wall
		{X: 2200, Y: 180, W: 360, H: 44, IsWater: false, IsDropThru: false},
		{X: 2240, Y: 130, W: 120, H: 12, IsWater: false, IsDropThru: true},
		{X: 2320, Y: 85, W: 80, H: 12, IsWater: false, IsDropThru: true},
	}
	st.Platforms = platforms
}

func (st *Stage) Update(player *Player, boss *Boss) {
	st.WaveOffset += 0.05

	// Camera moves right when player passes halfway screen
	targetCamX := player.X - st.ViewWidth*0.4
	if targetCamX > st.CameraX {
		st.CameraX = targetCamX
	}

	// Lock camera at Boss Arena
	if st.CameraX >= st.BossArenaX {
		st.CameraX = st.BossArenaX
		if !boss.Active && !boss.Defeated {
			boss.Active = true
		}
	}

	// Stage end boundary
	maxCamX := st.Length - st.ViewWidth
	if st.CameraX > maxCamX {
		st.CameraX = maxCamX
	}
}

func (st *Stage) Reset() {
	st.CameraX = 0
	st.WaveOffset = 0
}

func (st *Stage) GetWaterWave(x float64) float64 {
	return math.Sin(x*0.06+st.WaveOffset) * 2.5
}
