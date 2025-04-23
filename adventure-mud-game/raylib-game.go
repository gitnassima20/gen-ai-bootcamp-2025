package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1600
	screenHeight = 700
)

var (
	running      = true
	bkgColor     = rl.NewColor(147, 211, 196, 255)
	grasseSprite rl.Texture2D
	playerSprite rl.Texture2D

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerRight, playerLeft bool
	playerFrame                                   int

	frameCount int

	playerSpeed float32 = 3

	bkgSoundPaused bool
	bkgSound       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	rl.DrawTexture(grasseSprite, 100, 50, rl.White)
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerDest.Y -= playerSpeed
		playerMoving = true
		playerDir = 1
		playerUp = true
	}

	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerDest.Y += playerSpeed
		playerMoving = true
		playerDir = 0
		playerDown = true
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerDest.X -= playerSpeed
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerDest.X += playerSpeed
		playerMoving = true
		playerDir = 3
		playerRight = true
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		bkgSoundPaused = !bkgSoundPaused
	}
}

func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = 0

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
		playerSrc.X = playerSrc.Width * float32(playerFrame)
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}

	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(bkgSound)

	if bkgSoundPaused {
		rl.PauseMusicStream(bkgSound)
	} else {
		rl.ResumeMusicStream(bkgSound)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2))

	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Gēmu")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grasseSprite = rl.LoadTexture("./res/tilesets/Grass.png")
	playerSprite = rl.LoadTexture("./res/characters/BasicCharakterSpritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 100, 100)

	rl.InitAudioDevice()
	bkgSound = rl.LoadMusicStream("./res/bkg-sound.mp3")
	bkgSoundPaused = false
	rl.PlayMusicStream(bkgSound)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2)), 0.0, 2.0)
}

func quit() {
	rl.UnloadTexture(grasseSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(bkgSound)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()
}
