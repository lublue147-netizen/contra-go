package main

import (
	"syscall/js"
)

var (
	game     *Game
	renderer *Renderer
	audio    *AudioSystem
)

func main() {
	c := make(chan struct{}, 0)

	window := js.Global()
	document := window.Get("document")

	canvas := document.Call("getElementById", "gameCanvas")
	if !canvas.Truthy() {
		println("Error: gameCanvas element not found")
		return
	}

	audio = InitAudio()
	game = NewGame(audio)
	renderer = NewRenderer(canvas)

	// Setup Keyboard Listeners
	setupKeyboardInput(window)

	// Setup Touch / Mobile API
	setupVirtualInput(window)

	// Setup Game loop via requestAnimationFrame
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Poll Gamepad if connected
		pollGamepad(window)

		// Update game logic
		game.Update()

		// Render frame
		renderer.Render(game)

		// Request next frame
		window.Call("requestAnimationFrame", renderFrame)
		return nil
	})
	window.Call("requestAnimationFrame", renderFrame)

	// Inform web UI that Wasm is loaded and ready
	onReady := window.Get("onContraWasmReady")
	if onReady.Truthy() {
		onReady.Invoke()
	}

	<-c
}

func setupKeyboardInput(window js.Value) {
	window.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		key := e.Get("key").String()
		code := e.Get("code").String()

		// Audio context unlock on first key
		audio.ensureContext()

		// Konami code check
		game.Player.CheckKonami(key)

		switch key {
		case "ArrowLeft", "a", "A":
			game.Input.Left = true
		case "ArrowRight", "d", "D":
			game.Input.Right = true
		case "ArrowUp", "w", "W":
			game.Input.Up = true
		case "ArrowDown", "s", "S":
			game.Input.Down = true
		case "j", "J", "z", "Z", " ":
			game.Input.Shoot = true
		case "k", "K", "x", "X":
			game.Input.Jump = true
		case "Enter":
			game.Input.Start = true
		case "m", "M":
			audio.ToggleMute()
		case "c", "C":
			game.Player.ToggleCharacter()
		}

		// Prevent browser scrolling with arrow keys or space
		if key == "ArrowUp" || key == "ArrowDown" || key == "ArrowLeft" || key == "ArrowRight" || key == " " {
			e.Call("preventDefault")
		}
		_ = code
		return nil
	}))

	window.Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		key := e.Get("key").String()

		switch key {
		case "ArrowLeft", "a", "A":
			game.Input.Left = false
		case "ArrowRight", "d", "D":
			game.Input.Right = false
		case "ArrowUp", "w", "W":
			game.Input.Up = false
		case "ArrowDown", "s", "S":
			game.Input.Down = false
		case "j", "J", "z", "Z", " ":
			game.Input.Shoot = false
		case "k", "K", "x", "X":
			game.Input.Jump = false
		case "Enter":
			game.Input.Start = false
		}
		return nil
	}))
}

func setupVirtualInput(window js.Value) {
	// Expose setVirtualInput for mobile on-screen buttons
	window.Set("setVirtualInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return nil
		}
		btn := args[0].String()
		pressed := args[1].Bool()

		audio.ensureContext()

		switch btn {
		case "left":
			game.Input.Left = pressed
		case "right":
			game.Input.Right = pressed
		case "up":
			game.Input.Up = pressed
		case "down":
			game.Input.Down = pressed
		case "shoot":
			game.Input.Shoot = pressed
		case "jump":
			game.Input.Jump = pressed
		case "start":
			game.Input.Start = pressed
		case "konami":
			game.Player.Lives = 30
			game.Player.Weapon = WeaponSpread
			if globalAudio != nil {
				globalAudio.PlayKonami()
			}
		case "character":
			game.Player.ToggleCharacter()
		}
		return nil
	}))

	// Expose sound toggle
	window.Set("toggleMute", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return audio.ToggleMute()
	}))

	// Expose character toggle
	window.Set("toggleCharacter", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return game.Player.ToggleCharacter()
	}))

}

func pollGamepad(window js.Value) {
	navigator := window.Get("navigator")
	if !navigator.Truthy() {
		return
	}
	getGamepads := navigator.Get("getGamepads")
	if !getGamepads.Truthy() {
		return
	}

	gamepads := navigator.Call("getGamepads")
	if !gamepads.Truthy() || gamepads.Length() == 0 {
		return
	}

	pad := gamepads.Index(0)
	if !pad.Truthy() {
		return
	}

	axes := pad.Get("axes")
	buttons := pad.Get("buttons")

	if axes.Truthy() && axes.Length() >= 2 {
		axisX := axes.Index(0).Float()
		axisY := axes.Index(1).Float()

		if axisX < -0.3 {
			game.Input.Left = true
			game.Input.Right = false
		} else if axisX > 0.3 {
			game.Input.Right = true
			game.Input.Left = false
		}

		if axisY < -0.3 {
			game.Input.Up = true
			game.Input.Down = false
		} else if axisY > 0.3 {
			game.Input.Down = true
			game.Input.Up = false
		}
	}

	if buttons.Truthy() && buttons.Length() >= 10 {
		// Button 0 (A / Cross) -> Jump
		if buttons.Index(0).Get("pressed").Bool() {
			game.Input.Jump = true
		}
		// Button 2 (X / Square) or Button 1 (B / Circle) -> Shoot
		if buttons.Index(1).Get("pressed").Bool() || buttons.Index(2).Get("pressed").Bool() {
			game.Input.Shoot = true
		}
		// Button 9 (Start) -> Start
		if buttons.Index(9).Get("pressed").Bool() {
			game.Input.Start = true
		}
	}
}
