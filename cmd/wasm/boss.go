package main

import (
	"math"
	"math/rand"
)

func NewBoss(x, y float64) *Boss {
	return &Boss{
		X:         x,
		Y:         y,
		Width:     110,
		Height:    160,
		Active:    false,
		Defeated:  false,
		CoreHP:    80,
		MaxCoreHP: 80,
		TopTurret: BossPart{
			Name:    "TopCannon",
			OffsetX: 50,
			OffsetY: -50,
			Width:   26,
			Height:  26,
			HP:      25,
			MaxHP:   25,
		},
		LeftTurret: BossPart{
			Name:    "LeftTurret",
			OffsetX: 25,
			OffsetY: 10,
			Width:   22,
			Height:  22,
			HP:      20,
			MaxHP:   20,
		},
		RightTurret: BossPart{
			Name:    "RightTurret",
			OffsetX: 75,
			OffsetY: 10,
			Width:   22,
			Height:  22,
			HP:      20,
			MaxHP:   20,
		},
	}
}

func (b *Boss) Reset(x, y float64) {
	b.X = x
	b.Y = y
	b.Active = false
	b.Defeated = false
	b.CoreHP = 80
	b.DeathTimer = 0
	b.DefeatAnim = 0
	b.TopTurret.Destroyed = false
	b.TopTurret.HP = b.TopTurret.MaxHP
	b.LeftTurret.Destroyed = false
	b.LeftTurret.HP = b.LeftTurret.MaxHP
	b.RightTurret.Destroyed = false
	b.RightTurret.HP = b.RightTurret.MaxHP
}

func (b *Boss) Update(player *Player, playerBullets []*Bullet, particles *[]*Particle, texts *[]*FloatingText, enemies *[]*Enemy, em *EnemyManager) []*Bullet {
	if !b.Active {
		return nil
	}

	var bossBullets []*Bullet

	if b.Defeated {
		b.DeathTimer++
		// Fireworks explosion sequence
		if b.DeathTimer%8 == 0 && b.DeathTimer < 180 {
			ex := b.X + rand.Float64()*b.Width - 10
			ey := b.Y - 60 + rand.Float64()*b.Height
			if globalAudio != nil {
				globalAudio.PlayBossExplosion()
			}
			for i := 0; i < 10; i++ {
				*particles = append(*particles, &Particle{
					X:       ex,
					Y:       ey,
					VX:      (rand.Float64() - 0.5) * 8,
					VY:      (rand.Float64() - 0.5) * 8,
					Life:    30,
					MaxLife: 30,
					Color:   "#FFAA00",
					Size:    5,
				})
			}
		}
		return nil
	}

	if b.Flashing > 0 {
		b.Flashing--
	}
	if b.TopTurret.Flashing > 0 {
		b.TopTurret.Flashing--
	}
	if b.LeftTurret.Flashing > 0 {
		b.LeftTurret.Flashing--
	}
	if b.RightTurret.Flashing > 0 {
		b.RightTurret.Flashing--
	}

	// 1. Soldier Hatch Logic (spawns soldiers from the bunker floor)
	b.SoldierTimer++
	if b.SoldierTimer > 150 {
		b.SoldierTimer = 0
		*enemies = append(*enemies, em.createEnemy(EnemyTypeSoldier, b.X+30, b.Y+40, false, WeaponNormal))
	}

	// 2. Top Cannon behavior
	if !b.TopTurret.Destroyed {
		b.TopTurret.ShootTimer++
		if b.TopTurret.ShootTimer > 120 {
			b.TopTurret.ShootTimer = 0
			// Fire parabolic shell or aimed high-power shot
			cx := b.X + b.TopTurret.OffsetX
			cy := b.Y + b.TopTurret.OffsetY
			angle := math.Atan2(player.Y-cy, player.X-cx)
			speed := 4.2
			bossBullets = append(bossBullets, &Bullet{
				X:       cx,
				Y:       cy,
				VX:      math.Cos(angle) * speed,
				VY:      math.Sin(angle) * speed,
				Damage:  1,
				IsEnemy: true,
				Life:    180,
				Radius:  5.0,
				Color:   "#FF3300",
			})
		}
	}

	// 3. Left Turret behavior
	if !b.LeftTurret.Destroyed {
		b.LeftTurret.ShootTimer++
		if b.LeftTurret.ShootTimer > 100 {
			b.LeftTurret.ShootTimer = 0
			lx := b.X + b.LeftTurret.OffsetX
			ly := b.Y + b.LeftTurret.OffsetY
			angle := math.Atan2(player.Y-ly, player.X-lx)
			speed := 3.6
			bossBullets = append(bossBullets, &Bullet{
				X:       lx,
				Y:       ly,
				VX:      math.Cos(angle) * speed,
				VY:      math.Sin(angle) * speed,
				Damage:  1,
				IsEnemy: true,
				Life:    180,
				Radius:  4.0,
				Color:   "#FFAA22",
			})
		}
	}

	// 4. Right Turret behavior
	if !b.RightTurret.Destroyed {
		b.RightTurret.ShootTimer++
		if b.RightTurret.ShootTimer > 100 {
			b.RightTurret.ShootTimer = 0
			rx := b.X + b.RightTurret.OffsetX
			ry := b.Y + b.RightTurret.OffsetY
			angle := math.Atan2(player.Y-ry, player.X-rx)
			speed := 3.6
			bossBullets = append(bossBullets, &Bullet{
				X:       rx,
				Y:       ry,
				VX:      math.Cos(angle) * speed,
				VY:      math.Sin(angle) * speed,
				Damage:  1,
				IsEnemy: true,
				Life:    180,
				Radius:  4.0,
				Color:   "#FFAA22",
			})
		}
	}

	// 5. Player Bullets vs Boss
	for _, bBullet := range playerBullets {
		if bBullet.Life <= 0 || bBullet.IsEnemy {
			continue
		}

		// Check Top Turret
		if !b.TopTurret.Destroyed {
			tx := b.X + b.TopTurret.OffsetX
			ty := b.Y + b.TopTurret.OffsetY
			if math.Abs(bBullet.X-tx) < (b.TopTurret.Width/2+bBullet.Radius) && math.Abs(bBullet.Y-ty) < (b.TopTurret.Height/2+bBullet.Radius) {
				bBullet.Life = 0
				b.TopTurret.HP -= bBullet.Damage
				b.TopTurret.Flashing = 5
				if b.TopTurret.HP <= 0 {
					b.TopTurret.Destroyed = true
					player.Score += 1000
					b.createExplosionAt(tx, ty, particles)
				}
				continue
			}
		}

		// Check Left Turret
		if !b.LeftTurret.Destroyed {
			lx := b.X + b.LeftTurret.OffsetX
			ly := b.Y + b.LeftTurret.OffsetY
			if math.Abs(bBullet.X-lx) < (b.LeftTurret.Width/2+bBullet.Radius) && math.Abs(bBullet.Y-ly) < (b.LeftTurret.Height/2+bBullet.Radius) {
				bBullet.Life = 0
				b.LeftTurret.HP -= bBullet.Damage
				b.LeftTurret.Flashing = 5
				if b.LeftTurret.HP <= 0 {
					b.LeftTurret.Destroyed = true
					player.Score += 800
					b.createExplosionAt(lx, ly, particles)
				}
				continue
			}
		}

		// Check Right Turret
		if !b.RightTurret.Destroyed {
			rx := b.X + b.RightTurret.OffsetX
			ry := b.Y + b.RightTurret.OffsetY
			if math.Abs(bBullet.X-rx) < (b.RightTurret.Width/2+bBullet.Radius) && math.Abs(bBullet.Y-ry) < (b.RightTurret.Height/2+bBullet.Radius) {
				bBullet.Life = 0
				b.RightTurret.HP -= bBullet.Damage
				b.RightTurret.Flashing = 5
				if b.RightTurret.HP <= 0 {
					b.RightTurret.Destroyed = true
					player.Score += 800
					b.createExplosionAt(rx, ry, particles)
				}
				continue
			}
		}

		// Check Central Core Sensor (Heart)
		coreX := b.X + 50
		coreY := b.Y - 10
		if math.Abs(bBullet.X-coreX) < 22 && math.Abs(bBullet.Y-coreY) < 22 {
			bBullet.Life = 0
			b.CoreHP -= bBullet.Damage
			b.Flashing = 6
			if b.CoreHP <= 0 {
				b.Defeated = true
				b.CoreHP = 0
				player.Score += 10000
				if globalAudio != nil {
					globalAudio.PlayBossExplosion()
				}
				*texts = append(*texts, &FloatingText{
					Text:  "STAGE 1 CLEAR!",
					X:     coreX,
					Y:     coreY - 30,
					Life:  180,
					Color: "#FFFF00",
				})
				break
			}
		}
	}

	return bossBullets
}

func (b *Boss) createExplosionAt(x, y float64, particles *[]*Particle) {
	if globalAudio != nil {
		globalAudio.PlayExplosion()
	}
	for i := 0; i < 12; i++ {
		*particles = append(*particles, &Particle{
			X:       x,
			Y:       y,
			VX:      (rand.Float64() - 0.5) * 6,
			VY:      (rand.Float64() - 0.5) * 6,
			Life:    25,
			MaxLife: 25,
			Color:   "#FF6600",
			Size:    4,
		})
	}
}
