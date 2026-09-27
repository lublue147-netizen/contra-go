package main

import (
	"math"
)

type Player struct {
	X, Y            float64
	VX, VY          float64
	Width, Height   float64
	State           PlayerState
	FacingRight     bool
	AimDir          Direction
	OnGround        bool
	InWater         bool
	CanDropThru     bool
	DropThruTimer   int
	Lives           int
	Score           int
	Weapon          WeaponType
	ShootCooldown   int
	InvincibleTimer int
	BarrierTimer    int
	FlipAngle       float64
	RunFrame        int
	RunTimer        int
	RespawnTimer    int
	KonamiProgress  int
	TotalKills      int
}

var konamiSequence = []string{
	"ArrowUp", "ArrowUp",
	"ArrowDown", "ArrowDown",
	"ArrowLeft", "ArrowRight",
	"ArrowLeft", "ArrowRight",
	"b", "a",
}

func NewPlayer(startX, startY float64) *Player {
	return &Player{
		X:               startX,
		Y:               startY,
		Width:           24,
		Height:          38,
		State:           PlayerIdle,
		FacingRight:     true,
		AimDir:          DirRight,
		Lives:           3,
		Score:           0,
		Weapon:          WeaponNormal,
		InvincibleTimer: 120, // 2 seconds flashing at start
	}
}

func (p *Player) Reset(x, y float64) {
	p.X = x
	p.Y = y
	p.VX = 0
	p.VY = 0
	p.State = PlayerIdle
	p.FacingRight = true
	p.AimDir = DirRight
	p.OnGround = true
	p.InWater = false
	p.Weapon = WeaponNormal
	p.BarrierTimer = 0
	p.InvincibleTimer = 180
	p.RespawnTimer = 0
	p.FlipAngle = 0
}

func (p *Player) CheckKonami(key string) bool {
	if p.KonamiProgress < len(konamiSequence) {
		expected := konamiSequence[p.KonamiProgress]
		if key == expected || (key == "B" && expected == "b") || (key == "A" && expected == "a") {
			p.KonamiProgress++
			if p.KonamiProgress == len(konamiSequence) {
				p.Lives = 30
				p.Weapon = WeaponSpread
				p.KonamiProgress = 0
				if globalAudio != nil {
					globalAudio.PlayKonami()
				}
				return true
			}
			return false
		}
	}
	p.KonamiProgress = 0
	return false
}

func (p *Player) Update(input *InputState, stage *Stage) []*Bullet {
	if p.State == PlayerDying {
		p.RespawnTimer++
		p.VY += 0.35
		p.Y += p.VY
		p.X += p.VX
		p.FlipAngle += 0.2
		if p.RespawnTimer > 90 {
			if p.Lives > 0 {
				p.Lives--
				// Respawn slightly back and above ground
				spawnX := math.Max(stage.CameraX+60, p.X-100)
				p.Reset(spawnX, 80)
			}
		}
		return nil
	}

	if p.InvincibleTimer > 0 {
		p.InvincibleTimer--
	}
	if p.BarrierTimer > 0 {
		p.BarrierTimer--
	}
	if p.ShootCooldown > 0 {
		p.ShootCooldown--
	}
	if p.DropThruTimer > 0 {
		p.DropThruTimer--
	}

	// 1. Determine Aim Direction & Movement
	moveSpeed := 3.2
	if p.InWater {
		moveSpeed = 2.0
	}

	dx := 0.0
	if input.Left {
		dx -= moveSpeed
		p.FacingRight = false
	}
	if input.Right {
		dx += moveSpeed
		p.FacingRight = true
	}

	// Crouch logic
	p.CanDropThru = false
	if input.Down && p.OnGround && !p.InWater {
		if input.Left || input.Right {
			// Prone crawl or aim down while running
			p.State = PlayerCrouching
			p.Height = 20
		} else {
			// Pure prone
			p.State = PlayerCrouching
			p.Height = 18
			dx = 0
		}
		p.CanDropThru = true
	} else if !p.OnGround {
		p.State = PlayerJumping
		p.Height = 24
		p.FlipAngle += 0.25
	} else if dx != 0 {
		p.State = PlayerRunning
		p.Height = 38
		p.RunTimer++
		if p.RunTimer%6 == 0 {
			p.RunFrame = (p.RunFrame + 1) % 4
		}
	} else {
		p.State = PlayerIdle
		p.Height = 38
		p.FlipAngle = 0
	}

	// Aim Direction calculation
	if p.State == PlayerCrouching {
		if p.FacingRight {
			p.AimDir = DirRight
		} else {
			p.AimDir = DirLeft
		}
	} else if p.State == PlayerJumping {
		if input.Up {
			if input.Right {
				p.AimDir = DirUpRight
			} else if input.Left {
				p.AimDir = DirUpLeft
			} else {
				p.AimDir = DirUp
			}
		} else if input.Down {
			if input.Right {
				p.AimDir = DirDownRight
			} else if input.Left {
				p.AimDir = DirDownLeft
			} else {
				p.AimDir = DirDown
			}
		} else if p.FacingRight {
			p.AimDir = DirRight
		} else {
			p.AimDir = DirLeft
		}
	} else { // Idle or Running on ground
		if input.Up {
			if input.Right {
				p.AimDir = DirUpRight
			} else if input.Left {
				p.AimDir = DirUpLeft
			} else {
				p.AimDir = DirUp
			}
		} else if input.Down && (input.Right || input.Left) {
			if input.Right {
				p.AimDir = DirDownRight
			} else {
				p.AimDir = DirDownLeft
			}
		} else if p.FacingRight {
			p.AimDir = DirRight
		} else {
			p.AimDir = DirLeft
		}
	}

	// 2. Jump Handling
	if input.JumpJustPressed {
		if input.Down && p.CanDropThru {
			// Drop through platform
			p.DropThruTimer = 16
			p.OnGround = false
			p.VY = 2.0
		} else if p.OnGround || p.InWater {
			p.VY = -7.6
			p.OnGround = false
			p.InWater = false
			p.State = PlayerJumping
			p.FlipAngle = 0.1
			if globalAudio != nil {
				globalAudio.PlayJump()
			}
		}
	}

	// 3. Apply Horizontal Movement & Boundaries
	p.VX = dx
	p.X += p.VX

	// Camera left border collision (Contra classic: player cannot scroll backwards)
	minX := stage.CameraX + 12
	if p.X < minX {
		p.X = minX
	}
	// Boss right boundary
	if p.X > stage.Length-24 {
		p.X = stage.Length - 24
	}

	// 4. Gravity & Vertical Physics
	gravity := 0.42
	p.VY += gravity
	if p.VY > 9.0 {
		p.VY = 9.0
	}
	p.Y += p.VY

	// 5. Platform Collisions
	p.checkPlatformCollisions(stage)

	// Pit death
	if p.Y > stage.Height+40 {
		p.Die()
		return nil
	}

	// 6. Shooting
	var newBullets []*Bullet
	if input.Shoot && p.ShootCooldown == 0 {
		newBullets = p.fireWeapon()
	}

	return newBullets
}

func (p *Player) checkPlatformCollisions(stage *Stage) {
	footX := p.X
	footY := p.Y + p.Height/2
	prevFootY := footY - p.VY

	p.OnGround = false
	p.InWater = false

	for _, plat := range stage.Platforms {
		// Water detection
		if plat.IsWater {
			if p.X >= plat.X && p.X <= plat.X+plat.W && footY >= plat.Y && footY <= plat.Y+plat.H {
				p.InWater = true
				if p.VY > 0 {
					p.Y = plat.Y + 6 - p.Height/2
					p.VY = 0
					p.OnGround = true
				}
				continue
			}
		}

		// Drop thru platform check
		if plat.IsDropThru && p.DropThruTimer > 0 {
			continue
		}

		// Solid top landing check
		if footX >= plat.X-10 && footX <= plat.X+plat.W+10 {
			// Land on top edge if falling
			if prevFootY <= plat.Y+6 && footY >= plat.Y {
				p.Y = plat.Y - p.Height/2
				p.VY = 0
				p.OnGround = true
				break
			}
		}
	}
}

func (p *Player) fireWeapon() []*Bullet {
	var bullets []*Bullet
	fireX := p.X
	fireY := p.Y - 6

	if p.State == PlayerCrouching {
		fireY = p.Y + 4
	} else if p.AimDir == DirUp {
		fireY = p.Y - p.Height/2 - 4
	}

	angle := p.getAimAngle()

	switch p.Weapon {
	case WeaponNormal:
		p.ShootCooldown = 14
		speed := 8.5
		bullets = append(bullets, &Bullet{
			X:       fireX,
			Y:       fireY,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Damage:  1,
			Type:    WeaponNormal,
			IsEnemy: false,
			Life:    140,
			Radius:  3.5,
			Color:   "#FFFFFF",
		})
		if globalAudio != nil {
			globalAudio.PlayShoot()
		}

	case WeaponSpread:
		p.ShootCooldown = 22
		speed := 7.5
		spreadAngles := []float64{-0.35, -0.17, 0, 0.17, 0.35}
		for _, sa := range spreadAngles {
			a := angle + sa
			bullets = append(bullets, &Bullet{
				X:       fireX,
				Y:       fireY,
				VX:      math.Cos(a) * speed,
				VY:      math.Sin(a) * speed,
				Damage:  1,
				Type:    WeaponSpread,
				IsEnemy: false,
				Life:    120,
				Radius:  4.5,
				Color:   "#FF3333",
			})
		}
		if globalAudio != nil {
			globalAudio.PlaySpread()
		}

	case WeaponLaser:
		p.ShootCooldown = 18
		speed := 12.0
		bullets = append(bullets, &Bullet{
			X:           fireX,
			Y:           fireY,
			VX:          math.Cos(angle) * speed,
			VY:          math.Sin(angle) * speed,
			Damage:      3,
			Type:        WeaponLaser,
			IsEnemy:     false,
			Life:        100,
			Radius:      5.0,
			Color:       "#00FFFF",
			PierceCount: 4,
		})
		if globalAudio != nil {
			globalAudio.PlayLaser()
		}

	case WeaponMachine:
		p.ShootCooldown = 7 // Super fast auto-fire
		speed := 9.0
		bullets = append(bullets, &Bullet{
			X:       fireX,
			Y:       fireY,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Damage:  1,
			Type:    WeaponMachine,
			IsEnemy: false,
			Life:    130,
			Radius:  4.0,
			Color:   "#FFDD00",
		})
		if globalAudio != nil {
			globalAudio.PlayShoot()
		}

	case WeaponBarrier:
		p.ShootCooldown = 12
		speed := 8.5
		bullets = append(bullets, &Bullet{
			X:       fireX,
			Y:       fireY,
			VX:      math.Cos(angle) * speed,
			VY:      math.Sin(angle) * speed,
			Damage:  2,
			Type:    WeaponBarrier,
			IsEnemy: false,
			Life:    120,
			Radius:  4.0,
			Color:   "#00FF66",
		})
		if globalAudio != nil {
			globalAudio.PlayShoot()
		}
	}

	return bullets
}

func (p *Player) getAimAngle() float64 {
	switch p.AimDir {
	case DirRight:
		return 0
	case DirLeft:
		return math.Pi
	case DirUp:
		return -math.Pi / 2
	case DirDown:
		return math.Pi / 2
	case DirUpRight:
		return -math.Pi / 4
	case DirUpLeft:
		return -3 * math.Pi / 4
	case DirDownRight:
		return math.Pi / 4
	case DirDownLeft:
		return 3 * math.Pi / 4
	default:
		if p.FacingRight {
			return 0
		}
		return math.Pi
	}
}

func (p *Player) Die() {
	if p.State == PlayerDying || p.InvincibleTimer > 0 || p.BarrierTimer > 0 {
		return
	}
	p.State = PlayerDying
	p.VY = -5.0
	p.VX = -1.5
	if !p.FacingRight {
		p.VX = 1.5
	}
	p.RespawnTimer = 0
	p.Weapon = WeaponNormal
	if globalAudio != nil {
		globalAudio.PlayDeath()
	}
}
