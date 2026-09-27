package main

import (
	"math"
	"strconv"
	"syscall/js"
)

type InputState struct {
	Left           bool
	Right          bool
	Up             bool
	Down           bool
	Shoot          bool
	ShootJustPress bool
	Jump           bool
	JumpJustPressed bool
	Start          bool
}

type Game struct {
	State         GameState
	Player        *Player
	Stage         *Stage
	EnemyMgr      *EnemyManager
	Boss          *Boss
	Bullets       []*Bullet
	Particles     []*Particle
	DropItems     []*DropItem
	FloatingTexts []*FloatingText
	Input         InputState
	HiScore       int
	FrameCount    int
	audio         *AudioSystem
	prevJump      bool
	prevShoot     bool
}

func NewGame(audio *AudioSystem) *Game {
	st := NewStage1()
	pl := NewPlayer(60, 140)
	em := NewEnemyManager()
	boss := NewBoss(2400, 160)

	g := &Game{
		State:    StateTitle,
		Player:   pl,
		Stage:    st,
		EnemyMgr: em,
		Boss:     boss,
		HiScore:  20000,
		audio:    audio,
	}

	// Try reading HiScore from localStorage
	storage := js.Global().Get("localStorage")
	if storage.Truthy() {
		saved := storage.Call("getItem", "contra_hiscore")
		if saved.Truthy() {
			if val, err := strconv.Atoi(saved.String()); err == nil && val > g.HiScore {
				g.HiScore = val
			}
		}
	}

	return g
}

func (g *Game) StartGame() {
	g.State = StatePlaying
	g.Stage.Reset()
	g.Player.Reset(60, 140)
	g.EnemyMgr.Reset()
	g.Boss.Reset(2400, 160)
	g.Bullets = nil
	g.Particles = nil
	g.DropItems = nil
	g.FloatingTexts = nil
	g.FrameCount = 0
}

func (g *Game) Update() {
	g.FrameCount++

	// Just pressed flags
	g.Input.JumpJustPressed = g.Input.Jump && !g.prevJump
	g.Input.ShootJustPress = g.Input.Shoot && !g.prevShoot
	g.prevJump = g.Input.Jump
	g.prevShoot = g.Input.Shoot

	// State machine
	switch g.State {
	case StateTitle:
		if g.Input.Start || g.Input.JumpJustPressed || g.Input.ShootJustPress {
			g.StartGame()
		}
		return

	case StateGameOver:
		if g.Input.Start || g.Input.JumpJustPressed || g.Input.ShootJustPress {
			g.StartGame()
		}
		return

	case StateVictory:
		if g.Input.Start || g.Input.JumpJustPressed || g.Input.ShootJustPress {
			g.StartGame()
		}
		return

	case StatePlaying:
		g.updatePlaying()
	}
}

func (g *Game) updatePlaying() {
	// 1. Update Player
	newBullets := g.Player.Update(&g.Input, g.Stage)
	if len(newBullets) > 0 {
		g.Bullets = append(g.Bullets, newBullets...)
	}

	// Update HiScore
	if g.Player.Score > g.HiScore {
		g.HiScore = g.Player.Score
		storage := js.Global().Get("localStorage")
		if storage.Truthy() {
			storage.Call("setItem", "contra_hiscore", strconv.Itoa(g.HiScore))
		}
	}

	// Game Over check
	if g.Player.Lives < 0 {
		g.State = StateGameOver
		return
	}

	// 2. Update Stage & Camera
	g.Stage.Update(g.Player, g.Boss)

	// 3. Update Enemies
	var playerBullets []*Bullet
	for _, b := range g.Bullets {
		if !b.IsEnemy {
			playerBullets = append(playerBullets, b)
		}
	}

	enemyBullets, updatedPlayerBullets := g.EnemyMgr.Update(
		g.Stage,
		g.Player,
		playerBullets,
		&g.Particles,
		&g.DropItems,
		&g.FloatingTexts,
	)
	if len(enemyBullets) > 0 {
		g.Bullets = append(g.Bullets, enemyBullets...)
	}

	// 4. Update Boss
	bossBullets := g.Boss.Update(
		g.Player,
		updatedPlayerBullets,
		&g.Particles,
		&g.FloatingTexts,
		&g.EnemyMgr.Enemies,
		g.EnemyMgr,
	)
	if len(bossBullets) > 0 {
		g.Bullets = append(g.Bullets, bossBullets...)
	}

	// Check Boss Victory
	if g.Boss.Defeated && g.Boss.DeathTimer > 200 {
		g.State = StateVictory
		if globalAudio != nil {
			globalAudio.PlayVictory()
		}
		return
	}

	// 5. Update Bullets & Collide Enemy Bullets with Player
	g.Bullets = UpdateBullets(g.Bullets, g.Stage)

	if g.Player.State != PlayerDying && g.Player.InvincibleTimer == 0 && g.Player.BarrierTimer == 0 {
		for _, b := range g.Bullets {
			if !b.IsEnemy || b.Life <= 0 {
				continue
			}
			pdx := b.X - g.Player.X
			pdy := b.Y - g.Player.Y
			if math.Abs(pdx) < (g.Player.Width/2+b.Radius-2) && math.Abs(pdy) < (g.Player.Height/2+b.Radius-2) {
				b.Life = 0
				g.Player.Die()
				break
			}
		}
	}

	// 6. Update Drop Items
	g.DropItems = UpdateDropItems(g.DropItems, g.Stage, g.Player, &g.FloatingTexts)

	// 7. Update Particles
	var activeParticles []*Particle
	for _, p := range g.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life--
		if p.Life > 0 {
			activeParticles = append(activeParticles, p)
		}
	}
	g.Particles = activeParticles

	// 8. Update Floating Texts
	var activeTexts []*FloatingText
	for _, t := range g.FloatingTexts {
		t.Y -= 0.4
		t.Life--
		if t.Life > 0 {
			activeTexts = append(activeTexts, t)
		}
	}
	g.FloatingTexts = activeTexts
}
