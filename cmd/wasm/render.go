package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)


type Renderer struct {
	canvas js.Value
	ctx    js.Value
	width  float64
	height float64
}

func NewRenderer(canvas js.Value) *Renderer {
	ctx := canvas.Call("getContext", "2d")
	// Turn off image smoothing for authentic crisp pixel art
	ctx.Set("imageSmoothingEnabled", false)

	return &Renderer{
		canvas: canvas,
		ctx:    ctx,
		width:  256,
		height: 224,
	}
}

func (r *Renderer) Render(game *Game) {
	ctx := r.ctx
	stage := game.Stage
	camX := stage.CameraX

	// Clear canvas
	ctx.Set("fillStyle", "#0a0a1a")
	ctx.Call("fillRect", 0, 0, r.width, r.height)

	switch game.State {
	case StateTitle:
		r.drawBackground(stage)
		r.drawTitleScreen(game)
		return
	case StateGameOver:
		r.drawBackground(stage)
		r.drawGameOverScreen(game)
		return
	case StateVictory:
		r.drawBackground(stage)
		r.drawVictoryScreen(game)
		return
	}

	ctx.Call("save")

	if game.ScreenShake > 0 {
		sx := (rand.Float64() - 0.5) * game.ShakeIntensity
		sy := (rand.Float64() - 0.5) * game.ShakeIntensity
		ctx.Call("translate", sx, sy)
	}

	// 1. Draw Parallax Background
	r.drawBackground(stage)

	// 2. Draw Platforms & Terrain
	r.drawPlatforms(stage)

	// 3. Draw Drop Items
	r.drawDropItems(game.DropItems, camX)

	// 4. Draw Enemies
	r.drawEnemies(game.EnemyMgr.Enemies, camX)

	// 5. Draw Boss
	if game.Boss.Active || game.Boss.Defeated {
		r.drawBoss(game.Boss, camX)
	}

	// 6. Draw Bullets
	r.drawBullets(game.Bullets, camX)

	// 7. Draw Player
	if game.Player.Lives >= 0 {
		r.drawPlayer(game.Player, camX)
	}

	// 8. Draw Particles & Floating Texts
	r.drawParticles(game.Particles, camX)
	r.drawFloatingTexts(game.FloatingTexts, camX)

	// 9. Draw HUD
	r.drawHUD(game)

	ctx.Call("restore")
}


func (r *Renderer) drawBackground(stage *Stage) {
	ctx := r.ctx
	camX := stage.CameraX

	// Sky gradient
	grad := ctx.Call("createLinearGradient", 0, 0, 0, r.height)
	grad.Call("addColorStop", 0.0, "#080c1c")
	grad.Call("addColorStop", 0.6, "#182848")
	grad.Call("addColorStop", 1.0, "#0b3b4f")
	ctx.Set("fillStyle", grad)
	ctx.Call("fillRect", 0, 0, r.width, r.height)

	// Stars
	ctx.Set("fillStyle", "#ffffff")
	for i := 0; i < 25; i++ {
		sx := math.Mod(float64(i*37)-camX*0.05, r.width)
		if sx < 0 {
			sx += r.width
		}
		sy := math.Mod(float64(i*29), 90)
		ctx.Call("fillRect", sx, sy, 1.5, 1.5)
	}

	// Distant Mountains (parallax 0.15x)
	ctx.Set("fillStyle", "#102238")
	ctx.Call("beginPath")
	ctx.Call("moveTo", 0, 150)
	for x := 0.0; x <= r.width; x += 32 {
		worldX := x + camX*0.15
		my := 115 + math.Sin(worldX*0.015)*25 + math.Cos(worldX*0.04)*10
		ctx.Call("lineTo", x, my)
	}
	ctx.Call("lineTo", r.width, 180)
	ctx.Call("lineTo", 0, 180)
	ctx.Call("fill")

	// Mid-ground Jungle Canopy (parallax 0.4x)
	ctx.Set("fillStyle", "#0a382c")
	ctx.Call("beginPath")
	ctx.Call("moveTo", 0, 170)
	for x := 0.0; x <= r.width; x += 20 {
		worldX := x + camX*0.4
		jy := 140 + math.Sin(worldX*0.03)*16
		ctx.Call("lineTo", x, jy)
	}
	ctx.Call("lineTo", r.width, 180)
	ctx.Call("lineTo", 0, 180)
	ctx.Call("fill")
}

func (r *Renderer) drawPlatforms(stage *Stage) {
	ctx := r.ctx
	camX := stage.CameraX

	for _, plat := range stage.Platforms {
		if plat.Destroyed {
			continue
		}

		screenX := plat.X - camX
		if screenX+plat.W < 0 || screenX > r.width {
			continue
		}

		if plat.IsWater {
			// Animated water stream
			ctx.Set("fillStyle", "#0d4b68")
			ctx.Call("fillRect", screenX, plat.Y, plat.W, plat.H)

			// Water wave crest
			ctx.Set("fillStyle", "#40a9d4")
			ctx.Call("beginPath")
			ctx.Call("moveTo", screenX, plat.Y)
			for wx := 0.0; wx <= plat.W; wx += 8 {
				waveY := plat.Y + stage.GetWaterWave(plat.X+wx)
				ctx.Call("lineTo", screenX+wx, waveY)
			}
			ctx.Call("lineTo", screenX+plat.W, plat.Y+6)
			ctx.Call("lineTo", screenX, plat.Y+6)
			ctx.Call("fill")
			continue
		}

		if plat.IsDropThru {
			// Wooden suspension bridge or cliff ledge
			bridgeColor := "#8b5a2b"
			if plat.Exploding {
				if (plat.ExplodeTimer/2)%2 == 0 {
					bridgeColor = "#ff3300"
				} else {
					bridgeColor = "#ffcc00"
				}
			}
			ctx.Set("fillStyle", bridgeColor)
			ctx.Call("fillRect", screenX, plat.Y, plat.W, plat.H)
			// Plank lines
			ctx.Set("fillStyle", "#5c3a1e")
			for bx := 0.0; bx < plat.W; bx += 10 {
				ctx.Call("fillRect", screenX+bx, plat.Y, 2, plat.H)
			}
			ctx.Set("fillStyle", "#a66f38")

			ctx.Call("fillRect", screenX, plat.Y, plat.W, 2)
		} else {
			// Solid earth jungle ground
			// Earth body
			ctx.Set("fillStyle", "#382513")
			ctx.Call("fillRect", screenX, plat.Y, plat.W, plat.H)

			// Rock texture
			ctx.Set("fillStyle", "#2b1c0e")
			for rx := 0.0; rx < plat.W; rx += 16 {
				ctx.Call("fillRect", screenX+rx+4, plat.Y+8, 8, 8)
			}

			// Top lush grass / moss
			ctx.Set("fillStyle", "#228b22")
			ctx.Call("fillRect", screenX, plat.Y, plat.W, 4)
			ctx.Set("fillStyle", "#32cd32")
			for gx := 0.0; gx < plat.W; gx += 4 {
				gh := 2.0
				if int(gx)%8 == 0 {
					gh = 4.0
				}
				ctx.Call("fillRect", screenX+gx, plat.Y-gh+2, 2, gh)
			}
		}
	}
}

func (r *Renderer) drawPlayer(p *Player, camX float64) {
	ctx := r.ctx
	sx := p.X - camX
	sy := p.Y

	// Flashing transparency when invincible after respawn
	if p.InvincibleTimer > 0 && (p.InvincibleTimer/4)%2 == 0 {
		return
	}

	ctx.Call("save")
	ctx.Call("translate", sx, sy)

	// Barrier Shield Effect
	if p.BarrierTimer > 0 {
		ctx.Set("strokeStyle", "#00ffff")
		ctx.Set("lineWidth", 2)
		ctx.Call("beginPath")
		ctx.Call("arc", 0, 0, 24, 0, math.Pi*2)
		ctx.Call("stroke")

		// Orbiting orbs
		orbAngle := float64(p.BarrierTimer) * 0.15
		for o := 0; o < 3; o++ {
			ang := orbAngle + float64(o)*(math.Pi*2/3)
			ox := math.Cos(ang) * 24
			oy := math.Sin(ang) * 24
			ctx.Set("fillStyle", "#ffff00")
			ctx.Call("beginPath")
			ctx.Call("arc", ox, oy, 4, 0, math.Pi*2)
			ctx.Call("fill")
		}
	}

	pantsColor := "#1b3f8b"
	headbandColor := "#e60000"
	hairColor := "#6b4423"
	if p.IsLanceBean {
		pantsColor = "#cc2222"
		headbandColor = "#1166ee"
		hairColor = "#111111"
	}

	// Somersault Jump (Spinning ball)
	if p.State == PlayerJumping {
		ctx.Call("rotate", p.FlipAngle)
		// Body ball
		ctx.Set("fillStyle", pantsColor)
		ctx.Call("beginPath")
		ctx.Call("arc", 0, 0, 10, 0, math.Pi*2)
		ctx.Call("fill")
		// Torso & head
		ctx.Set("fillStyle", "#d2906b")
		ctx.Call("beginPath")
		ctx.Call("arc", 0, -3, 7, 0, math.Pi*2)
		ctx.Call("fill")
		// Headband
		ctx.Set("fillStyle", headbandColor)
		ctx.Call("fillRect", -5, -8, 10, 3)
		ctx.Call("restore")
		return
	}

	// Lying Prone / Crouching
	if p.State == PlayerCrouching {
		dir := 1.0
		if !p.FacingRight {
			dir = -1.0
		}
		// Torso lying flat
		ctx.Set("fillStyle", "#d2906b")
		ctx.Call("fillRect", -12*dir, 2, 20*dir, 8)
		// Pants
		ctx.Set("fillStyle", pantsColor)
		ctx.Call("fillRect", -18*dir, 4, 10*dir, 6)
		// Head & headband
		ctx.Set("fillStyle", "#d2906b")
		ctx.Call("fillRect", 6*dir, 0, 8*dir, 8)
		ctx.Set("fillStyle", headbandColor)
		ctx.Call("fillRect", 6*dir, 0, 8*dir, 3)
		// Gun forward
		ctx.Set("fillStyle", "#444444")
		ctx.Call("fillRect", 12*dir, 4, 14*dir, 3)
		ctx.Call("restore")
		return
	}

	// Standard Upright / Running pose
	dir := 1.0
	if !p.FacingRight {
		dir = -1.0
	}

	// Legs
	ctx.Set("fillStyle", pantsColor)
	if p.State == PlayerRunning {
		// Running frame cycle
		switch p.RunFrame {
		case 0:
			ctx.Call("fillRect", -6*dir, 6, 6*dir, 12)
			ctx.Call("fillRect", 2*dir, 6, 6*dir, 12)
		case 1:
			ctx.Call("fillRect", -10*dir, 6, 6*dir, 10)
			ctx.Call("fillRect", 4*dir, 6, 8*dir, 12)
		case 2:
			ctx.Call("fillRect", -4*dir, 6, 8*dir, 12)
			ctx.Call("fillRect", 0, 6, 6*dir, 12)
		case 3:
			ctx.Call("fillRect", -8*dir, 6, 8*dir, 12)
			ctx.Call("fillRect", 6*dir, 6, 6*dir, 10)
		}
	} else {
		// Idle standing legs
		ctx.Call("fillRect", -6*dir, 6, 5*dir, 12)
		ctx.Call("fillRect", 2*dir, 6, 5*dir, 12)
	}

	// Bare Torso (muscle tone)
	ctx.Set("fillStyle", "#d2906b")
	ctx.Call("fillRect", -5*dir, -8, 10*dir, 14)

	// Head & Headband
	ctx.Set("fillStyle", "#d2906b")
	ctx.Call("fillRect", -4*dir, -18, 8*dir, 9)
	// Hair
	ctx.Set("fillStyle", hairColor)
	ctx.Call("fillRect", -4*dir, -18, 8*dir, 3)
	// Headband
	ctx.Set("fillStyle", headbandColor)
	ctx.Call("fillRect", -5*dir, -16, 9*dir, 3)
	// Flowing headband ribbon
	ctx.Call("fillRect", -10*dir, -16, 5*dir, 2)


	// Gun & Arms based on Aim Direction
	ctx.Set("fillStyle", "#444444")
	switch p.AimDir {
	case DirUp:
		ctx.Call("fillRect", 0, -28, 3*dir, 18) // Gun pointing straight up
	case DirUpRight, DirUpLeft:
		// Diagonal 45°
		ctx.Call("save")
		ctx.Call("rotate", -math.Pi/4*dir)
		ctx.Call("fillRect", 2*dir, -4, 16*dir, 4)
		ctx.Call("restore")
	case DirDownRight, DirDownLeft:
		// Diagonal down 45°
		ctx.Call("save")
		ctx.Call("rotate", math.Pi/4*dir)
		ctx.Call("fillRect", 2*dir, -4, 16*dir, 4)
		ctx.Call("restore")
	default:
		// Straight forward
		ctx.Call("fillRect", 2*dir, -4, 16*dir, 4)
	}

	ctx.Call("restore")
}

func (r *Renderer) drawBullets(bullets []*Bullet, camX float64) {
	ctx := r.ctx
	for _, b := range bullets {
		sx := b.X - camX
		sy := b.Y

		ctx.Set("fillStyle", b.Color)
		ctx.Call("beginPath")
		if b.Type == WeaponLaser {
			// Elongated laser bolt
			ctx.Call("save")
			ctx.Call("translate", sx, sy)
			angle := math.Atan2(b.VY, b.VX)
			ctx.Call("rotate", angle)
			ctx.Call("fillRect", -12, -2, 24, 4)
			ctx.Call("restore")
		} else {
			ctx.Call("arc", sx, sy, b.Radius, 0, math.Pi*2)
			ctx.Call("fill")
		}
	}
}

func (r *Renderer) drawEnemies(enemies []*Enemy, camX float64) {
	ctx := r.ctx
	for _, e := range enemies {
		if !e.Active {
			continue
		}
		sx := e.X - camX
		sy := e.Y

		ctx.Call("save")
		ctx.Call("translate", sx, sy)

		dir := 1.0
		if e.FacingLeft {
			dir = -1.0
		}

		switch e.Type {
		case EnemyTypeSoldier:
			// Red uniform soldier
			ctx.Set("fillStyle", "#cc2222") // Red jacket
			ctx.Call("fillRect", -4*dir, -10, 8*dir, 12)
			// Pants
			ctx.Set("fillStyle", "#223366")
			if e.AnimFrame%2 == 0 {
				ctx.Call("fillRect", -6*dir, 2, 5*dir, 10)
				ctx.Call("fillRect", 1*dir, 2, 5*dir, 10)
			} else {
				ctx.Call("fillRect", -8*dir, 2, 6*dir, 8)
				ctx.Call("fillRect", 3*dir, 2, 6*dir, 10)
			}
			// Head & helmet
			ctx.Set("fillStyle", "#d2906b")
			ctx.Call("fillRect", -3*dir, -16, 6*dir, 6)
			ctx.Set("fillStyle", "#445566") // Helmet
			ctx.Call("fillRect", -4*dir, -18, 8*dir, 4)
			// Rifle
			ctx.Set("fillStyle", "#222222")
			ctx.Call("fillRect", 2*dir, -4, 10*dir, 3)

		case EnemyTypeSniper:
			// Prone sniper
			ctx.Set("fillStyle", "#2d5a27") // Camouflage green
			ctx.Call("fillRect", -10*dir, 0, 18*dir, 8)
			ctx.Set("fillStyle", "#1b3f18")
			ctx.Call("fillRect", -14*dir, 2, 8*dir, 6)
			// Head
			ctx.Set("fillStyle", "#d2906b")
			ctx.Call("fillRect", 6*dir, -2, 6*dir, 6)
			// Long rifle barrel
			ctx.Set("fillStyle", "#111111")
			ctx.Call("fillRect", 10*dir, 0, 16*dir, 3)

		case EnemyTypeTurret:
			// Armored Pillbox Turret
			ctx.Set("fillStyle", "#555566")
			ctx.Call("beginPath")
			ctx.Call("arc", 0, 4, 12, math.Pi, 0)
			ctx.Call("fill")
			// Rotating Cannon Barrel
			ctx.Call("save")
			ctx.Call("rotate", e.AimAngle)
			ctx.Set("fillStyle", "#222233")
			ctx.Call("fillRect", 0, -3, 16, 6)
			ctx.Call("restore")

		case EnemyTypeFalconCapsule:
			// Winged power capsule
			ctx.Set("fillStyle", "#ee2222")
			ctx.Call("beginPath")
			ctx.Call("arc", 0, 0, 8, 0, math.Pi*2)
			ctx.Call("fill")
			// Falcon wings
			ctx.Set("fillStyle", "#ffffff")
			wingFlap := math.Sin(float64(e.StateTimer)*0.2) * 5
			ctx.Call("beginPath")
			ctx.Call("moveTo", -12, wingFlap)
			ctx.Call("lineTo", 0, 0)
			ctx.Call("lineTo", 12, wingFlap)
			ctx.Call("lineTo", 0, -4)
			ctx.Call("fill")
		}

		ctx.Call("restore")
	}
}

func (r *Renderer) drawBoss(b *Boss, camX float64) {
	ctx := r.ctx
	sx := b.X - camX
	sy := b.Y

	ctx.Call("save")
	ctx.Call("translate", sx, sy)

	// Metal fortress wall base
	ctx.Set("fillStyle", "#3a3d4d")
	ctx.Call("fillRect", 0, -80, b.Width, b.Height)

	// Steel plates & rivets
	ctx.Set("fillStyle", "#2a2c38")
	ctx.Call("fillRect", 4, -76, b.Width-8, 4)
	ctx.Call("fillRect", 4, 0, b.Width-8, 4)
	ctx.Call("fillRect", 4, 60, b.Width-8, 4)

	// Top Cannon
	if !b.TopTurret.Destroyed {
		if b.TopTurret.Flashing > 0 {
			ctx.Set("fillStyle", "#ffffff")
		} else {
			ctx.Set("fillStyle", "#606880")
		}
		ctx.Call("fillRect", b.TopTurret.OffsetX-12, b.TopTurret.OffsetY-12, b.TopTurret.Width, b.TopTurret.Height)
		// Cannon muzzle
		ctx.Set("fillStyle", "#1a1b24")
		ctx.Call("fillRect", b.TopTurret.OffsetX-18, b.TopTurret.OffsetY-4, 8, 8)
	}

	// Left & Right Turrets
	drawTurret := func(p *BossPart) {
		if p.Destroyed {
			ctx.Set("fillStyle", "#222222")
			ctx.Call("fillRect", p.OffsetX-8, p.OffsetY-8, 16, 16)
			return
		}
		if p.Flashing > 0 {
			ctx.Set("fillStyle", "#ffffff")
		} else {
			ctx.Set("fillStyle", "#707890")
		}
		ctx.Call("beginPath")
		ctx.Call("arc", p.OffsetX, p.OffsetY, 10, 0, math.Pi*2)
		ctx.Call("fill")
		ctx.Set("fillStyle", "#222222")
		ctx.Call("fillRect", p.OffsetX-14, p.OffsetY-2, 10, 4)
	}
	drawTurret(&b.LeftTurret)
	drawTurret(&b.RightTurret)

	// Center Sensor Core (Red Falcon Core)
	coreX := 50.0
	coreY := -10.0
	if b.Flashing > 0 {
		ctx.Set("fillStyle", "#ffffff")
	} else if b.Defeated {
		ctx.Set("fillStyle", "#331111")
	} else {
		// Pulsing red heart
		pulse := (math.Sin(float64(b.SoldierTimer)*0.1) + 1.0) * 0.5
		if pulse > 0.5 {
			ctx.Set("fillStyle", "#ff1122")
		} else {
			ctx.Set("fillStyle", "#aa0011")
		}
	}
	ctx.Call("beginPath")
	ctx.Call("arc", coreX, coreY, 14, 0, math.Pi*2)
	ctx.Call("fill")

	// Lower Soldier Hatch Door
	ctx.Set("fillStyle", "#1b1d24")
	ctx.Call("fillRect", 20, 30, 40, 34)

	ctx.Call("restore")
}

func (r *Renderer) drawDropItems(items []*DropItem, camX float64) {
	ctx := r.ctx
	for _, it := range items {
		if !it.Active {
			continue
		}
		sx := it.X - camX
		sy := it.Y

		// Badge background
		ctx.Set("fillStyle", it.Color)
		ctx.Call("beginPath")
		ctx.Call("arc", sx, sy, 8, 0, math.Pi*2)
		ctx.Call("fill")

		// Letter text (S, M, L, B)
		ctx.Set("fillStyle", "#ffffff")
		ctx.Set("font", "bold 10px monospace")
		ctx.Set("textAlign", "center")
		ctx.Set("textBaseline", "middle")
		ctx.Call("fillText", it.Letter, sx, sy)
	}
}

func (r *Renderer) drawParticles(particles []*Particle, camX float64) {
	ctx := r.ctx
	for _, p := range particles {
		sx := p.X - camX
		sy := p.Y
		ctx.Set("fillStyle", p.Color)
		ctx.Call("fillRect", sx-p.Size/2, sy-p.Size/2, p.Size, p.Size)
	}
}

func (r *Renderer) drawFloatingTexts(texts []*FloatingText, camX float64) {
	ctx := r.ctx
	for _, t := range texts {
		sx := t.X - camX
		sy := t.Y
		ctx.Set("fillStyle", t.Color)
		ctx.Set("font", "bold 10px monospace")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", t.Text, sx, sy)
	}
}

func (r *Renderer) drawHUD(game *Game) {
	ctx := r.ctx

	ctx.Set("font", "bold 9px monospace")
	ctx.Set("textBaseline", "top")

	// 1P SCORE
	ctx.Set("fillStyle", "#ff3333")
	ctx.Set("textAlign", "left")
	ctx.Call("fillText", "1P", 8, 6)
	ctx.Set("fillStyle", "#ffffff")
	ctx.Call("fillText", fmt.Sprintf("%06d", game.Player.Score), 24, 6)

	// HI-SCORE
	ctx.Set("fillStyle", "#ff3333")
	ctx.Set("textAlign", "center")
	ctx.Call("fillText", "HI-SCORE", r.width/2-24, 6)
	ctx.Set("fillStyle", "#ffffff")
	ctx.Call("fillText", fmt.Sprintf("%06d", game.HiScore), r.width/2+32, 6)

	// REST LIVES (mini player icons)
	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("textAlign", "left")
	ctx.Call("fillText", "REST", 8, 18)
	for i := 0; i < game.Player.Lives && i < 10; i++ {
		lx := 36.0 + float64(i)*8
		ctx.Set("fillStyle", "#1b3f8b")
		ctx.Call("fillRect", lx, 18, 4, 6)
		ctx.Set("fillStyle", "#e60000")
		ctx.Call("fillRect", lx, 16, 4, 2)
	}
	if game.Player.Lives > 10 {
		ctx.Set("fillStyle", "#ffff00")
		ctx.Call("fillText", fmt.Sprintf("x%d", game.Player.Lives), 36+80, 18)
	}

	// CURRENT WEAPON ICON
	wName := "NORMAL"
	wColor := "#ffffff"
	switch game.Player.Weapon {
	case WeaponSpread:
		wName = "[S] SPREAD"
		wColor = "#ff3333"
	case WeaponLaser:
		wName = "[L] LASER"
		wColor = "#00ffff"
	case WeaponMachine:
		wName = "[M] MACHINE"
		wColor = "#ffdd00"
	case WeaponBarrier:
		wName = "[B] BARRIER"
		wColor = "#00ff66"
	}
	ctx.Set("fillStyle", wColor)
	ctx.Set("textAlign", "right")
	ctx.Call("fillText", wName, r.width-8, 6)

	// BOSS HP BAR (when active)
	if game.Boss.Active && !game.Boss.Defeated {
		ctx.Set("fillStyle", "#ff0000")
		ctx.Set("textAlign", "center")
		ctx.Call("fillText", "BOSS CORE", r.width/2, 20)

		barW := 80.0
		barH := 4.0
		bx := r.width/2 - barW/2
		by := 30.0

		ctx.Set("fillStyle", "#333333")
		ctx.Call("fillRect", bx, by, barW, barH)

		ratio := float64(game.Boss.CoreHP) / float64(game.Boss.MaxCoreHP)
		if ratio < 0 {
			ratio = 0
		}
		ctx.Set("fillStyle", "#00ff00")
		ctx.Call("fillRect", bx, by, barW*ratio, barH)
	}
}

func (r *Renderer) drawTitleScreen(game *Game) {
	ctx := r.ctx

	ctx.Set("textAlign", "center")

	// Retro Logo Box
	ctx.Set("fillStyle", "#e60000")
	ctx.Set("font", "900 28px sans-serif")
	ctx.Call("fillText", "CONTRA", r.width/2, 45)

	ctx.Set("fillStyle", "#ffd700")
	ctx.Set("font", "bold 16px sans-serif")
	ctx.Call("fillText", "魂 斗 罗", r.width/2, 75)

	ctx.Set("fillStyle", "#00ffff")
	ctx.Set("font", "10px monospace")
	ctx.Call("fillText", "STAGE 1: JUNGLE ISLAND", r.width/2, 98)

	// Konami code indicator if triggered
	if game.Player.Lives == 30 {
		ctx.Set("fillStyle", "#00ff00")
		ctx.Set("font", "bold 9px monospace")
		ctx.Call("fillText", "★ 30 LIVES CODE ACTIVATED! ★", r.width/2, 116)
	} else {
		ctx.Set("fillStyle", "#888888")
		ctx.Set("font", "8px monospace")
		ctx.Call("fillText", "CHEAT: UP UP DOWN DOWN LEFT RIGHT LEFT RIGHT B A", r.width/2, 116)
	}

	// Character Select indicator
	commandoName := "1P: BILL RIZER [比尔]"
	if game.Player.IsLanceBean {
		commandoName = "1P: LANCE BEAN [兰斯]"
	}
	ctx.Set("fillStyle", "#ffd700")
	ctx.Set("font", "bold 9px monospace")
	ctx.Call("fillText", commandoName+" (按 C 切换角色)", r.width/2, 130)

	// Start blinking prompt
	if (game.FrameCount/30)%2 == 0 {
		ctx.Set("fillStyle", "#ffffff")
		ctx.Set("font", "bold 11px monospace")
		ctx.Call("fillText", "PRESS ENTER / SPACE / TAP TO START", r.width/2, 148)
	}


	// Instructions
	ctx.Set("fillStyle", "#aaaaaa")
	ctx.Set("font", "8px monospace")
	ctx.Call("fillText", "WASD / ARROWS: MOVE & AIM", r.width/2, 172)
	ctx.Call("fillText", "J / Z / SPACE: SHOOT  |  K / X: JUMP", r.width/2, 186)
	ctx.Call("fillText", "DOWN + JUMP: DROP DOWN PLATFORMS", r.width/2, 200)
}

func (r *Renderer) drawGameOverScreen(game *Game) {
	ctx := r.ctx
	ctx.Set("textAlign", "center")

	ctx.Set("fillStyle", "#e60000")
	ctx.Set("font", "bold 22px monospace")
	ctx.Call("fillText", "GAME OVER", r.width/2, 60)

	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("font", "11px monospace")
	ctx.Call("fillText", fmt.Sprintf("FINAL SCORE: %06d", game.Player.Score), r.width/2, 95)
	ctx.Call("fillText", fmt.Sprintf("TOTAL KILLS: %d", game.Player.TotalKills), r.width/2, 115)

	if (game.FrameCount/30)%2 == 0 {
		ctx.Set("fillStyle", "#ffd700")
		ctx.Call("fillText", "PRESS ENTER / JUMP TO RETRY", r.width/2, 155)
	}
}

func (r *Renderer) drawVictoryScreen(game *Game) {
	ctx := r.ctx
	ctx.Set("textAlign", "center")

	ctx.Set("fillStyle", "#ffd700")
	ctx.Set("font", "bold 20px monospace")
	ctx.Call("fillText", "STAGE 1 CLEAR!", r.width/2, 55)

	ctx.Set("fillStyle", "#00ffff")
	ctx.Set("font", "12px monospace")
	ctx.Call("fillText", "BASE DESTROYED!", r.width/2, 82)

	ctx.Set("fillStyle", "#ffffff")
	ctx.Set("font", "10px monospace")
	ctx.Call("fillText", fmt.Sprintf("FINAL SCORE: %06d", game.Player.Score), r.width/2, 115)
	ctx.Call("fillText", fmt.Sprintf("TOTAL KILLS: %d", game.Player.TotalKills), r.width/2, 132)

	if (game.FrameCount/30)%2 == 0 {
		ctx.Set("fillStyle", "#00ff00")
		ctx.Call("fillText", "PRESS ENTER TO PLAY AGAIN", r.width/2, 168)
	}
}
