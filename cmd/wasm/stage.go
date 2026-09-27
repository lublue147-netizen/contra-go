package main

import (
	"math"
	"math/rand"
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
		Length:     2560,
		Height:     224,
		ViewWidth:  256,
		ViewHeight: 224,
		CameraX:    0,
		BossArenaX: 2200,
		WaterLevel: 210,
	}
	st.InitPlatforms()
	return st
}

func (st *Stage) InitPlatforms() {
	platforms := []Platform{
		// 1. Initial beach & jungle
		{ID: 1, X: 0, Y: 180, W: 450, H: 44, IsWater: false, IsDropThru: false},
		{ID: 2, X: 150, Y: 120, W: 140, H: 12, IsWater: false, IsDropThru: true},
		{ID: 3, X: 260, Y: 70, W: 120, H: 12, IsWater: false, IsDropThru: true},

		// 2. First bridge over water (3 explosive segments)
		{ID: 10, X: 450, Y: 180, W: 50, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 11, X: 500, Y: 180, W: 50, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 12, X: 550, Y: 180, W: 50, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 13, X: 450, Y: 206, W: 150, H: 18, IsWater: true, IsDropThru: false}, // Water stream

		// 3. Central Jungle Ridge
		{ID: 20, X: 600, Y: 180, W: 380, H: 44, IsWater: false, IsDropThru: false},
		{ID: 21, X: 660, Y: 125, W: 150, H: 12, IsWater: false, IsDropThru: true},
		{ID: 22, X: 780, Y: 80, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 4. Waterfall & stepped cliffs
		{ID: 30, X: 980, Y: 160, W: 200, H: 64, IsWater: false, IsDropThru: false},
		{ID: 31, X: 1040, Y: 105, W: 180, H: 12, IsWater: false, IsDropThru: true},
		{ID: 32, X: 1180, Y: 150, W: 240, H: 74, IsWater: false, IsDropThru: false},
		{ID: 33, X: 1240, Y: 90, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 5. Second bridge & deep lagoon (3 explosive segments)
		{ID: 40, X: 1420, Y: 180, W: 60, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 41, X: 1480, Y: 180, W: 60, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 42, X: 1540, Y: 180, W: 60, H: 12, IsWater: false, IsDropThru: true, IsBridge: true, ExplodeTimer: 24},
		{ID: 43, X: 1420, Y: 206, W: 180, H: 18, IsWater: true, IsDropThru: false},

		// 6. Mountain approach leading to Fortress
		{ID: 50, X: 1600, Y: 180, W: 350, H: 44, IsWater: false, IsDropThru: false},
		{ID: 51, X: 1680, Y: 120, W: 180, H: 12, IsWater: false, IsDropThru: true},
		{ID: 52, X: 1820, Y: 75, W: 160, H: 12, IsWater: false, IsDropThru: true},
		{ID: 53, X: 1950, Y: 180, W: 250, H: 44, IsWater: false, IsDropThru: false},
		{ID: 54, X: 2000, Y: 125, W: 160, H: 12, IsWater: false, IsDropThru: true},

		// 7. Boss Arena Wall
		{ID: 60, X: 2200, Y: 180, W: 360, H: 44, IsWater: false, IsDropThru: false},
		{ID: 61, X: 2240, Y: 130, W: 120, H: 12, IsWater: false, IsDropThru: true},
		{ID: 62, X: 2320, Y: 85, W: 80, H: 12, IsWater: false, IsDropThru: true},
	}
	st.Platforms = platforms
}

func (st *Stage) Update(player *Player, boss *Boss, particles *[]*Particle, triggerShake func(frames int, intensity float64)) {
	st.WaveOffset += 0.05

	// 1. Camera moves right when player passes halfway screen
	targetCamX := player.X - st.ViewWidth*0.4
	if targetCamX > st.CameraX {
		st.CameraX = targetCamX
	}

	// 2. Lock camera at Boss Arena
	if st.CameraX >= st.BossArenaX {
		st.CameraX = st.BossArenaX
		if !boss.Active && !boss.Defeated {
			boss.Active = true
		}
	}

	// 3. Stage end boundary
	maxCamX := st.Length - st.ViewWidth
	if st.CameraX > maxCamX {
		st.CameraX = maxCamX
	}

	// 4. Bridge Explosion Mechanics (iconic Contra Stage 1 bridge sequential collapse)
	for i := range st.Platforms {
		p := &st.Platforms[i]
		if !p.IsBridge || p.Destroyed {
			continue
		}

		// Trigger explosion sequence when player crosses onto or near this bridge section
		if !p.Exploding && player.X >= p.X+10 {
			p.Exploding = true
		}

		if p.Exploding {
			p.ExplodeTimer--
			if p.ExplodeTimer <= 0 {
				p.Destroyed = true
				if globalAudio != nil {
					globalAudio.PlayExplosion()
				}
				if triggerShake != nil {
					triggerShake(8, 3.5)
				}
				// Spawn fireball particles
				for k := 0; k < 10; k++ {
					*particles = append(*particles, &Particle{
						X:       p.X + float64(k)*(p.W/10),
						Y:       p.Y,
						VX:      (rand.Float64() - 0.5) * 5,
						VY:      -rand.Float64() * 4,
						Life:    25,
						MaxLife: 25,
						Color:   "#FF5500",
						Size:    4,
					})
				}
			}
		}
	}
}

func (st *Stage) Reset() {
	st.CameraX = 0
	st.WaveOffset = 0
	st.InitPlatforms()
}

func (st *Stage) GetWaterWave(x float64) float64 {
	return math.Sin(x*0.06+st.WaveOffset) * 2.5
}
