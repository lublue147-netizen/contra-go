package main

func UpdateBullets(bullets []*Bullet, stage *Stage) []*Bullet {
	var alive []*Bullet
	for _, b := range bullets {
		b.X += b.VX
		b.Y += b.VY
		b.Life--

		// Keep bullet if still within screen/stage bounds and has lifetime
		if b.Life > 0 && b.X >= stage.CameraX-50 && b.X <= stage.CameraX+stage.ViewWidth+50 && b.Y >= -50 && b.Y <= stage.Height+50 {
			alive = append(alive, b)
		}
	}
	return alive
}

func UpdateDropItems(items []*DropItem, stage *Stage, p *Player, texts *[]*FloatingText) []*DropItem {
	var alive []*DropItem
	for _, it := range items {
		if !it.Active {
			continue
		}
		it.Life--
		if it.Life <= 0 {
			continue
		}

		if !it.OnGround {
			it.VY += 0.25
			it.Y += it.VY

			// Check landing on platforms
			for _, plat := range stage.Platforms {
				if !plat.IsWater && it.X >= plat.X && it.X <= plat.X+plat.W {
					if it.Y >= plat.Y-8 && it.Y <= plat.Y+12 && it.VY > 0 {
						it.Y = plat.Y - 8
						it.VY = 0
						it.OnGround = true
						break
					}
				}
			}
		}

		// Check player pickup
		dx := p.X - it.X
		dy := p.Y - it.Y
		distSq := dx*dx + dy*dy
		if distSq < 26*26 && p.State != PlayerDying {
			// Collected!
			it.Active = false
			p.Weapon = it.Type
			if it.Type == WeaponBarrier {
				p.BarrierTimer = 600 // 10 seconds of invincibility
			}
			p.Score += 500
			if globalAudio != nil {
				globalAudio.PlayPowerup()
			}
			*texts = append(*texts, &FloatingText{
				Text:  "+" + it.Letter + " WEAPON",
				X:     it.X,
				Y:     it.Y - 10,
				Life:  60,
				Color: it.Color,
			})
			continue
		}

		alive = append(alive, it)
	}
	return alive
}
