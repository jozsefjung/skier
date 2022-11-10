package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Config struct {
	windowWidth  int32
	windowHeight int32
	windowTitle  string
}

type Slope struct {
	points     []rl.Vector2
	pointCount int32
	active     bool
	config     *Config
}

type Slps struct {
	slopes []Slope
}

type ParallaxBackground struct {
	textures         []rl.Texture2D
	speeds           []float32
	positions        []rl.Vector2
	initialPositions []rl.Vector2
	speedModifier    float32
	config           *Config
}

func (b *ParallaxBackground) init() {
	b.textures = make([]rl.Texture2D, 0)
	b.speeds = make([]float32, 0)
	b.positions = make([]rl.Vector2, 0)
	b.initialPositions = make([]rl.Vector2, 0)
	b.speedModifier = 1.0
}

func (b *ParallaxBackground) add(fileName string, scrollSpeed float32, initialPosition rl.Vector2) {
	img := rl.LoadImage(fileName)
	if img.Width != b.config.windowWidth || img.Height != b.config.windowHeight {
		rl.ImageResize(img, b.config.windowWidth, b.config.windowHeight)
	}

	b.textures = append(b.textures, rl.LoadTextureFromImage(img))
	b.speeds = append(b.speeds, scrollSpeed)
	b.positions = append(b.positions, initialPosition)
	b.initialPositions = append(b.initialPositions, initialPosition)
}

func (b *ParallaxBackground) update() {
	for i := 0; i < len(b.textures); i++ {
		b.positions[i].X -= b.speeds[i] * b.speedModifier
		if b.positions[i].X <= -float32(b.textures[i].Width) {
			b.positions[i].X = b.initialPositions[i].X
		}
	}
}

func (b *ParallaxBackground) draw() {
	for i, texture := range b.textures {
		rl.DrawTextureEx(texture, b.positions[i], 0, 1, rl.White)
		rl.DrawTexture(texture, int32(b.positions[i].X)+texture.Width, int32(b.positions[i].Y), rl.White)
	}
}

func NewParallaxBackground(config *Config) (b ParallaxBackground) {
	b.init()
	b.config = config
	return b
}

func (slope *Slope) init() {
	slope.pointCount = 0
	slope.points = make([]rl.Vector2, 0)
	slope.active = false
}

func (slope *Slope) add() {
	slope.pointCount++
	slope.points = append(slope.points, rl.GetMousePosition())
}

func (slope *Slope) draw() {
	if !slope.active && slope.lastPoint().X < -10 {
		return
	}
	if slope.pointCount > 0 {
		for i := 1; i < len(slope.points); i++ {

			pointA := slope.points[i-1]
			pointB := slope.points[i]
			pointNum := pointB.X - pointA.X
			diffX := pointB.X - pointA.X
			diffY := pointB.Y - pointA.Y
			intervalX := diffX / (pointNum + 1)
			intervalY := diffY / (pointNum + 1)
			rl.DrawLineEx(slope.points[i-1], slope.points[i], 5, rl.RayWhite)
			for j := 0; j < int(pointNum); j++ {
				x1, y1 := pointA.X+intervalX*float32(j), pointA.Y+intervalY*float32(j)
				x2, y2 := pointA.X+intervalX*float32(j), float32(slope.config.windowHeight)+pointA.Y+intervalY*float32(j)
				rl.DrawLineEx(rl.NewVector2(x1, y1+2.5), rl.NewVector2(x2, y2), 4, rl.LightGray)
				rl.DrawLineEx(rl.NewVector2(x1, y1+30), rl.NewVector2(x2, y2), 4, rl.Gray)
				rl.DrawLineEx(rl.NewVector2(x1, y1+60), rl.NewVector2(x2, y2), 4, rl.DarkGray)
			}

		}

	}
}

func (slope *Slope) scroll(speed float32) {
	for i := 0; i < int(slope.pointCount); i++ {
		slope.points[i].X -= speed
	}
}

func (slope *Slope) lastPoint() rl.Vector2 {
	return slope.points[len(slope.points)-1]
}

func NewSlope(config *Config) (slope Slope) {
	slope.init()
	slope.config = config
	return slope
}

func main() {
	cfg := Config{
		1920, 1080, "Skier",
	}
	rl.InitWindow(cfg.windowWidth, cfg.windowHeight, cfg.windowTitle)
	rl.SetTargetFPS(60)
	rl.HideCursor()

	slopes := make([]Slope, 0)
	bkg := NewParallaxBackground(&cfg)
	bkg.add("assets/landscape_0004_5_clouds.png", 1, rl.NewVector2(0, -100))
	bkg.add("assets/landscape_0003_4_mountain.png", 2, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0002_3_trees.png", 3, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0001_2_trees.png", 4, rl.NewVector2(0, 0))
	bkg.add("assets/landscape_0000_1_trees.png", 5, rl.NewVector2(0, 0))
	bkgGrad := rl.LoadTextureFromImage(rl.GenImageGradientV(int(cfg.windowWidth), int(0.65*float32(cfg.windowHeight)), rl.SkyBlue, rl.Beige))
	for !rl.WindowShouldClose() {

		mousepos := rl.GetMousePosition()
		for i := 0; i < len(slopes); i++ {
			slopes[i].scroll(6)
		}

		if (mousepos.X >= 0 && mousepos.X <= float32(cfg.windowWidth)) && (mousepos.Y >= 0 && mousepos.Y <= float32(cfg.windowHeight)) {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				slopes = append(slopes, NewSlope(&cfg))
				slopes[len(slopes)-1].add()
				slopes[len(slopes)-1].active = true
			}
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				if mousepos.X > slopes[len(slopes)-1].lastPoint().X {
					slopes[len(slopes)-1].add()
				}
			}
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				slopes[len(slopes)-1].active = false
			}
		}
		bkg.update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(235, 239, 242, 255))
		rl.DrawTexture(bkgGrad, 0, 0, rl.White)

		bkg.draw()
		for i := 0; i < len(slopes); i++ {
			slopes[i].draw()
		}
		rl.DrawCircleV(mousepos, 10, rl.RayWhite)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}