package main

import (
	"math"
	"math/rand"
)

type EnemySpawner struct {
	TriggerX  float64
	Type      EnemyType
	X, Y      float64
	DropsItem WeaponType
	HasDrop   bool
	Spawned   bool
}

type EnemyManager struct {
	Enemies   []*Enemy
	Spawners  []EnemySpawner
	nextID    int
	lastSpawn int
}

func NewEnemyManager() *EnemyManager {
	em := &EnemyManager{
		nextID: 1,
	}
	em.InitStage1Spawners()
	return em
}

func (em *EnemyManager) InitStage1Spawners() {
	// Populate Stage 1 jungle spawners along stage progression
	spawners := []EnemySpawner{
		// First section: intro soldiers & first falcon capsule
		{TriggerX: 120, Type: EnemyTypeSoldier, X: 450, Y: 140},
		{TriggerX: 200, Type: EnemyTypeSoldier, X: 520, Y: 140},
		{TriggerX: 280, Type: EnemyTypeFalconCapsule, X: 600, Y: 60, HasDrop: true, DropsItem: WeaponSpread}, // Give 'S' early!
		{TriggerX: 350, Type: EnemyTypeSniper, X: 680, Y: 130},
		{TriggerX: 420, Type: EnemyTypeSoldier, X: 750, Y: 140},

		// First bridge & island: Turrets and more soldiers
		{TriggerX: 550, Type: EnemyTypeTurret, X: 820, Y: 140},
		{TriggerX: 620, Type: EnemyTypeSoldier, X: 900, Y: 140},
		{TriggerX: 700, Type: EnemyTypeFalconCapsule, X: 1000, Y: 50, HasDrop: true, DropsItem: WeaponMachine},
		{TriggerX: 780, Type: EnemyTypeSniper, X: 1100, Y: 80},
		{TriggerX: 850, Type: EnemyTypeSoldier, X: 1180, Y: 140},

		// Mountain pass & waterfall section
		{TriggerX: 950, Type: EnemyTypeTurret, X: 1280, Y: 130},
		{TriggerX: 1020, Type: EnemyTypeFalconCapsule, X: 1350, Y: 70, HasDrop: true, DropsItem: WeaponLaser},
		{TriggerX: 1100, Type: EnemyTypeSoldier, X: 1420, Y: 90},
		{TriggerX: 1200, Type: EnemyTypeTurret, X: 1550, Y: 130},
		{TriggerX: 1300, Type: EnemyTypeSniper, X: 1650, Y: 70},
		{TriggerX: 1380, Type: EnemyTypeFalconCapsule, X: 1720, Y: 60, HasDrop: true, DropsItem: WeaponBarrier},
		{TriggerX: 1450, Type: EnemyTypeSoldier, X: 1800, Y: 140},

		// Final approach before boss fortress
		{TriggerX: 1600, Type: EnemyTypeTurret, X: 1950, Y: 130},
		{TriggerX: 1700, Type: EnemyTypeSniper, X: 2050, Y: 70},
		{TriggerX: 1800, Type: EnemyTypeFalconCapsule, X: 2150, Y: 50, HasDrop: true, DropsItem: WeaponSpread},
		{TriggerX: 1900, Type: EnemyTypeSoldier, X: 2280, Y: 140},
		{TriggerX: 2000, Type: EnemyTypeSoldier, X: 2380, Y: 140},
	}
	em.Spawners = spawners
}

func (em *EnemyManager) Reset() {
	em.Enemies = nil
	em.InitStage1Spawners()
	em.lastSpawn = 0
}

func (em *EnemyManager) Update(stage *Stage, player *Player, playerBullets []*Bullet, particles *[]*Particle, items *[]*DropItem, texts *[]*FloatingText) ([]*Bullet, []*Bullet) {
	var enemyBullets []*Bullet

	// 1. Check spawners
	for i := range em.Spawners {
		s := &em.Spawners[i]
		if !s.Spawned && stage.CameraX >= s.TriggerX {
			s.Spawned = true
			e := em.createEnemy(s.Type, s.X, s.Y, s.HasDrop, s.DropsItem)
			em.Enemies = append(em.Enemies, e)
		}
	}

	// Dynamic continuous soldier spawning to keep the Contra intensity high
	if stage.CameraX < stage.BossArenaX-100 {
		em.lastSpawn++
		if em.lastSpawn > 180 { // Every ~3 seconds, spawn running soldier from ahead
			em.lastSpawn = 0
			spawnX := stage.CameraX + stage.ViewWidth + 30
			if spawnX < stage.Length-100 {
				em.Enemies = append(em.Enemies, em.createEnemy(EnemyTypeSoldier, spawnX, 120, false, WeaponNormal))
			}
		}
	}

	// 2. Update active enemies
	var activeEnemies []*Enemy
	for _, e := range em.Enemies {
		if !e.Active {
			continue
		}

		e.StateTimer++
		e.AnimTimer++
		if e.AnimTimer%8 == 0 {
			e.AnimFrame = (e.AnimFrame + 1) % 4
		}

		// Enemy behavior by type
		switch e.Type {
		case EnemyTypeSoldier:
			// Run towards player
			if player.X < e.X {
				e.FacingLeft = true
				e.VX = -1.8
			} else {
				e.FacingLeft = false
				e.VX = 1.8
			}

			// Gravity & platform collision
			e.VY += 0.4
			if e.VY > 8.0 {
				e.VY = 8.0
			}
			e.X += e.VX
			e.Y += e.VY

			// Collide with platforms
			e.OnGround = false
			for _, plat := range stage.Platforms {
				if !plat.IsWater && e.X >= plat.X-8 && e.X <= plat.X+plat.W+8 {
					if e.Y >= plat.Y-e.Height/2-4 && e.Y <= plat.Y-e.Height/2+8 && e.VY > 0 {
						e.Y = plat.Y - e.Height/2
						e.VY = 0
						e.OnGround = true
						break
					}
				}
			}

		case EnemyTypeSniper:
			// Face player
			e.FacingLeft = player.X < e.X
			e.ShootTimer++
			if e.ShootTimer > 110 {
				e.ShootTimer = 0
				// Shoot bullet aimed at player
				angle := math.Atan2(player.Y-e.Y, player.X-e.X)
				speed := 3.2
				enemyBullets = append(enemyBullets, &Bullet{
					X:       e.X,
					Y:       e.Y - 2,
					VX:      math.Cos(angle) * speed,
					VY:      math.Sin(angle) * speed,
					Damage:  1,
					IsEnemy: true,
					Life:    180,
					Radius:  3.5,
					Color:   "#FF2222",
				})
			}

		case EnemyTypeTurret:
			// Rotating bunker turret
			e.AimAngle = math.Atan2(player.Y-e.Y, player.X-e.X)
			e.ShootTimer++
			if e.ShootTimer > 130 {
				e.ShootTimer = 0
				speed := 3.5
				enemyBullets = append(enemyBullets, &Bullet{
					X:       e.X + math.Cos(e.AimAngle)*12,
					Y:       e.Y + math.Sin(e.AimAngle)*12,
					VX:      math.Cos(e.AimAngle) * speed,
					VY:      math.Sin(e.AimAngle) * speed,
					Damage:  1,
					IsEnemy: true,
					Life:    180,
					Radius:  4.0,
					Color:   "#FFAA00",
				})
			}

		case EnemyTypeFalconCapsule:
			// Smooth sine wave flying across screen
			e.X -= 2.2
			e.Y += math.Sin(float64(e.StateTimer)*0.08) * 1.5

		case EnemyTypeSensor:
			// Stationary wall sensor
		}

		// Despawn if fallen off screen or way behind camera
		if e.X < stage.CameraX-100 || e.X > stage.CameraX+stage.ViewWidth+150 || e.Y > stage.Height+50 {
			e.Active = false
			continue
		}

		// 3. Collision with Player Bullets
		for _, b := range playerBullets {
			if b.Life <= 0 || b.IsEnemy {
				continue
			}
			dx := b.X - e.X
			dy := b.Y - e.Y
			if math.Abs(dx) < (e.Width/2+b.Radius) && math.Abs(dy) < (e.Height/2+b.Radius) {
				// Hit!
				e.HP -= b.Damage
				if b.PierceCount > 0 {
					b.PierceCount--
				} else {
					b.Life = 0 // Bullet consumed
				}

				// Spark particles
				*particles = append(*particles, &Particle{
					X:       b.X,
					Y:       b.Y,
					VX:      (rand.Float64() - 0.5) * 4,
					VY:      (rand.Float64() - 0.5) * 4,
					Life:    15,
					MaxLife: 15,
					Color:   "#FFFF55",
					Size:    3,
					IsSpark: true,
				})

				if e.HP <= 0 {
					e.Active = false
					player.Score += e.ScoreValue
					player.TotalKills++

					if globalAudio != nil {
						globalAudio.PlayExplosion()
					}

					// Explosion particles
					for pIdx := 0; pIdx < 8; pIdx++ {
						*particles = append(*particles, &Particle{
							X:       e.X,
							Y:       e.Y,
							VX:      (rand.Float64() - 0.5) * 6,
							VY:      (rand.Float64() - 0.5) * 6,
							Life:    25,
							MaxLife: 25,
							Color:   "#FF5500",
							Size:    4,
						})
					}

					// Item drop
					if e.HasDrop {
						letter := "S"
						color := "#FF3333"
						switch e.DropsItem {
						case WeaponSpread:
							letter = "S"
							color = "#FF3333"
						case WeaponLaser:
							letter = "L"
							color = "#00FFFF"
						case WeaponMachine:
							letter = "M"
							color = "#FFDD00"
						case WeaponBarrier:
							letter = "B"
							color = "#00FF66"
						}
						*items = append(*items, &DropItem{
							X:      e.X,
							Y:      e.Y,
							VY:     -3.0,
							Type:   e.DropsItem,
							Letter: letter,
							Color:  color,
							Active: true,
							Life:   600, // 10 seconds before despawn
						})
					}
					break
				}
			}
		}

		// 4. Collision with Player Body
		if e.Active && player.State != PlayerDying {
			pdx := player.X - e.X
			pdy := player.Y - e.Y
			if math.Abs(pdx) < (e.Width/2+player.Width/2-4) && math.Abs(pdy) < (e.Height/2+player.Height/2-4) {
				if player.BarrierTimer > 0 {
					// Barrier kills enemy instantly!
					e.Active = false
					player.Score += e.ScoreValue
					if globalAudio != nil {
						globalAudio.PlayExplosion()
					}
				} else {
					player.Die()
				}
			}
		}

		if e.Active {
			activeEnemies = append(activeEnemies, e)
		}
	}
	em.Enemies = activeEnemies

	// Return any new enemy bullets and updated player bullets
	return enemyBullets, playerBullets
}

func (em *EnemyManager) createEnemy(t EnemyType, x, y float64, hasDrop bool, dropType WeaponType) *Enemy {
	id := em.nextID
	em.nextID++

	switch t {
	case EnemyTypeSoldier:
		return &Enemy{
			ID:         id,
			Type:       t,
			X:          x,
			Y:          y,
			Width:      20,
			Height:     34,
			HP:         1,
			MaxHP:      1,
			ScoreValue: 100,
			FacingLeft: true,
			Active:     true,
			HasDrop:    hasDrop,
			DropsItem:  dropType,
		}
	case EnemyTypeSniper:
		return &Enemy{
			ID:         id,
			Type:       t,
			X:          x,
			Y:          y,
			Width:      22,
			Height:     32,
			HP:         2,
			MaxHP:      2,
			ScoreValue: 200,
			FacingLeft: true,
			Active:     true,
			HasDrop:    hasDrop,
			DropsItem:  dropType,
		}
	case EnemyTypeTurret:
		return &Enemy{
			ID:         id,
			Type:       t,
			X:          x,
			Y:          y,
			Width:      28,
			Height:     26,
			HP:         4,
			MaxHP:      4,
			ScoreValue: 300,
			Active:     true,
			HasDrop:    hasDrop,
			DropsItem:  dropType,
		}
	case EnemyTypeFalconCapsule:
		return &Enemy{
			ID:         id,
			Type:       t,
			X:          x,
			Y:          y,
			Width:      24,
			Height:     18,
			HP:         1,
			MaxHP:      1,
			ScoreValue: 200,
			Active:     true,
			HasDrop:    hasDrop,
			DropsItem:  dropType,
		}
	default:
		return &Enemy{
			ID:         id,
			Type:       t,
			X:          x,
			Y:          y,
			Width:      20,
			Height:     20,
			HP:         1,
			MaxHP:      1,
			ScoreValue: 100,
			Active:     true,
		}
	}
}
