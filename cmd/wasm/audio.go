package main

import (
	"syscall/js"
)

type AudioSystem struct {
	ctx       js.Value
	enabled   bool
	bgmNode   js.Value
	bgmGain   js.Value
	bgmTimer  js.Value
	muted     bool
}

var globalAudio *AudioSystem

func InitAudio() *AudioSystem {
	a := &AudioSystem{
		enabled: false,
		muted:   false,
	}
	globalAudio = a
	return a
}

func (a *AudioSystem) ensureContext() bool {
	if a.muted {
		return false
	}
	if !a.ctx.Truthy() {
		window := js.Global()
		audioCtxClass := window.Get("AudioContext")
		if !audioCtxClass.Truthy() {
			audioCtxClass = window.Get("webkitAudioContext")
		}
		if audioCtxClass.Truthy() {
			a.ctx = audioCtxClass.New()
			a.enabled = true
		}
	}
	if a.ctx.Truthy() && a.ctx.Get("state").String() == "suspended" {
		a.ctx.Call("resume")
	}
	return a.enabled
}

func (a *AudioSystem) ToggleMute() bool {
	a.muted = !a.muted
	if a.muted && a.bgmGain.Truthy() {
		a.bgmGain.Get("gain").Set("value", 0)
	}
	return a.muted
}

func (a *AudioSystem) PlayShoot() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	osc := a.ctx.Call("createOscillator")
	gain := a.ctx.Call("createGain")

	osc.Set("type", "square")
	osc.Get("frequency").Call("setValueAtTime", 900, now)
	osc.Get("frequency").Call("exponentialRampToValueAtTime", 200, now+0.08)

	gain.Get("gain").Call("setValueAtTime", 0.15, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.08)

	osc.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	osc.Call("start", now)
	osc.Call("stop", now+0.08)
}

func (a *AudioSystem) PlaySpread() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	for _, freq := range []float64{700, 520, 360} {
		osc := a.ctx.Call("createOscillator")
		gain := a.ctx.Call("createGain")

		osc.Set("type", "square")
		osc.Get("frequency").Call("setValueAtTime", freq, now)
		osc.Get("frequency").Call("exponentialRampToValueAtTime", 120, now+0.12)

		gain.Get("gain").Call("setValueAtTime", 0.1, now)
		gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.12)

		osc.Call("connect", gain)
		gain.Call("connect", a.ctx.Get("destination"))

		osc.Call("start", now)
		osc.Call("stop", now+0.12)
	}
}

func (a *AudioSystem) PlayLaser() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	osc := a.ctx.Call("createOscillator")
	gain := a.ctx.Call("createGain")

	osc.Set("type", "sawtooth")
	osc.Get("frequency").Call("setValueAtTime", 1400, now)
	osc.Get("frequency").Call("exponentialRampToValueAtTime", 100, now+0.16)

	gain.Get("gain").Call("setValueAtTime", 0.18, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.16)

	osc.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	osc.Call("start", now)
	osc.Call("stop", now+0.16)
}

func (a *AudioSystem) PlayJump() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	osc := a.ctx.Call("createOscillator")
	gain := a.ctx.Call("createGain")

	osc.Set("type", "square")
	osc.Get("frequency").Call("setValueAtTime", 160, now)
	osc.Get("frequency").Call("exponentialRampToValueAtTime", 500, now+0.14)

	gain.Get("gain").Call("setValueAtTime", 0.14, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.14)

	osc.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	osc.Call("start", now)
	osc.Call("stop", now+0.14)
}

func (a *AudioSystem) PlayExplosion() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	// Noise buffer simulation via random sound wave
	sampleRate := a.ctx.Get("sampleRate").Int()
	bufferSize := sampleRate / 6 // ~160ms
	buffer := a.ctx.Call("createBuffer", 1, bufferSize, sampleRate)
	data := buffer.Call("getChannelData", 0)

	js.Global().Call("eval", `(function(data, len) {
		for (var i = 0; i < len; i++) {
			data[i] = (Math.random() * 2 - 1) * Math.exp(-i / (len * 0.3));
		}
	})`).Invoke(data, bufferSize)

	noise := a.ctx.Call("createBufferSource")
	noise.Set("buffer", buffer)

	filter := a.ctx.Call("createBiquadFilter")
	filter.Set("type", "lowpass")
	filter.Get("frequency").Call("setValueAtTime", 800, now)
	filter.Get("frequency").Call("exponentialRampToValueAtTime", 60, now+0.18)

	gain := a.ctx.Call("createGain")
	gain.Get("gain").Call("setValueAtTime", 0.3, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.18)

	noise.Call("connect", filter)
	filter.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	noise.Call("start", now)
}

func (a *AudioSystem) PlayBossExplosion() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	sampleRate := a.ctx.Get("sampleRate").Int()
	bufferSize := sampleRate / 2 // 500ms
	buffer := a.ctx.Call("createBuffer", 1, bufferSize, sampleRate)
	data := buffer.Call("getChannelData", 0)

	js.Global().Call("eval", `(function(data, len) {
		for (var i = 0; i < len; i++) {
			data[i] = (Math.random() * 2 - 1) * Math.exp(-i / (len * 0.5));
		}
	})`).Invoke(data, bufferSize)

	noise := a.ctx.Call("createBufferSource")
	noise.Set("buffer", buffer)

	filter := a.ctx.Call("createBiquadFilter")
	filter.Set("type", "lowpass")
	filter.Get("frequency").Call("setValueAtTime", 400, now)
	filter.Get("frequency").Call("exponentialRampToValueAtTime", 40, now+0.5)

	gain := a.ctx.Call("createGain")
	gain.Get("gain").Call("setValueAtTime", 0.45, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.5)

	noise.Call("connect", filter)
	filter.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	noise.Call("start", now)
}

func (a *AudioSystem) PlayPowerup() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	notes := []float64{330, 440, 554.37, 659.25}
	for i, f := range notes {
		t := now + float64(i)*0.06
		osc := a.ctx.Call("createOscillator")
		gain := a.ctx.Call("createGain")

		osc.Set("type", "triangle")
		osc.Get("frequency").Call("setValueAtTime", f, t)

		gain.Get("gain").Call("setValueAtTime", 0.2, t)
		gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, t+0.07)

		osc.Call("connect", gain)
		gain.Call("connect", a.ctx.Get("destination"))

		osc.Call("start", t)
		osc.Call("stop", t+0.07)
	}
}

func (a *AudioSystem) PlayDeath() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	osc := a.ctx.Call("createOscillator")
	gain := a.ctx.Call("createGain")

	osc.Set("type", "sawtooth")
	osc.Get("frequency").Call("setValueAtTime", 440, now)
	osc.Get("frequency").Call("linearRampToValueAtTime", 80, now+0.35)

	gain.Get("gain").Call("setValueAtTime", 0.25, now)
	gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, now+0.35)

	osc.Call("connect", gain)
	gain.Call("connect", a.ctx.Get("destination"))

	osc.Call("start", now)
	osc.Call("stop", now+0.35)
}

func (a *AudioSystem) PlayKonami() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	fanfare := []float64{261.6, 329.6, 392.0, 523.2, 659.2, 784.0, 1046.5}
	for i, f := range fanfare {
		t := now + float64(i)*0.07
		osc := a.ctx.Call("createOscillator")
		gain := a.ctx.Call("createGain")

		osc.Set("type", "square")
		osc.Get("frequency").Call("setValueAtTime", f, t)

		gain.Get("gain").Call("setValueAtTime", 0.2, t)
		gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, t+0.12)

		osc.Call("connect", gain)
		gain.Call("connect", a.ctx.Get("destination"))

		osc.Call("start", t)
		osc.Call("stop", t+0.12)
	}
}

func (a *AudioSystem) PlayVictory() {
	if !a.ensureContext() {
		return
	}
	now := a.ctx.Get("currentTime").Float()

	victoryNotes := []struct {
		f float64
		d float64
	}{
		{523.25, 0.12}, // C5
		{523.25, 0.12},
		{523.25, 0.12},
		{659.25, 0.36}, // E5
		{587.33, 0.24}, // D5
		{659.25, 0.48}, // E5
		{783.99, 0.60}, // G5
	}

	t := now
	for _, n := range victoryNotes {
		osc := a.ctx.Call("createOscillator")
		gain := a.ctx.Call("createGain")

		osc.Set("type", "square")
		osc.Get("frequency").Call("setValueAtTime", n.f, t)

		gain.Get("gain").Call("setValueAtTime", 0.25, t)
		gain.Get("gain").Call("exponentialRampToValueAtTime", 0.01, t+n.d)

		osc.Call("connect", gain)
		gain.Call("connect", a.ctx.Get("destination"))

		osc.Call("start", t)
		osc.Call("stop", t+n.d)

		t += n.d * 1.05
	}
}
