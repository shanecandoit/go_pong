package main

import (
	"image/color"
	"math/rand"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 10
	paddleHeight = 70
	ballRadius   = 5
)

type object struct {
	x      float32
	y      float32
	dx     float32
	dy     float32
	width  float32
	height float32
}

type game struct {
	ball                    object
	leftPaddle, rightPaddle object
	scoreLeft, scoreRight   int
}

func initObject(o *object, x, y, width, height float32) {
	o.x = x
	o.y = y
	o.width = width
	o.height = height
}

func (g *game) Update() error {
	// Update ball position
	g.ball.x += g.ball.dx
	g.ball.y += g.ball.dy

	// Check for wall collisions
	if g.ball.y < 0 || g.ball.y+ballRadius > screenHeight {
		g.ball.dy *= -1
	}

	// Check for paddle collisions
	if g.ball.x < paddleWidth+g.leftPaddle.width &&
		g.ball.y > g.leftPaddle.y &&
		g.ball.y+ballRadius < g.leftPaddle.y+paddleHeight {
		g.ball.dx *= -1
	} else if g.ball.x > screenWidth-paddleWidth-g.rightPaddle.width &&
		g.ball.y > g.rightPaddle.y &&
		g.ball.y+ballRadius < g.rightPaddle.y+paddleHeight {
		g.ball.dx *= -1
	}

	// Check for scoring
	if g.ball.x < 0 {
		//g.ball.initPosition()
		// reset the ball position
		g.ball.x = screenWidth / 2
		g.ball.y = screenHeight / 2
		// random velocity
		g.ball.dx = 2
		g.ball.dy = 2
		fiftyFifty := rand.Intn(2)
		if fiftyFifty == 0 {
			g.ball.dx *= -1
		}
		fiftyFifty = rand.Intn(2)
		if fiftyFifty == 0 {
			g.ball.dy *= -1
		}

		g.scoreRight++
	} else if g.ball.x > screenWidth {
		// reset the ball position
		g.ball.x = screenWidth / 2
		g.ball.y = screenHeight / 2
		// random velocity
		g.ball.dx = 2
		g.ball.dy = 2
		fiftyFifty := rand.Intn(2)
		if fiftyFifty == 0 {
			g.ball.dx *= -1
		}
		fiftyFifty = rand.Intn(2)
		if fiftyFifty == 0 {
			g.ball.dy *= -1
		}
		g.scoreLeft++
	}

	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.leftPaddle.y > 0 {
			g.leftPaddle.y -= 4
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.leftPaddle.y < screenHeight-paddleHeight {
			g.leftPaddle.y += 4
		}
	}

	// AI movement
	if g.ball.y < g.rightPaddle.y+paddleHeight/2 {
		if g.rightPaddle.y > 0 {
			g.rightPaddle.y -= 2
		}
	}
	if g.ball.y > g.rightPaddle.y+paddleHeight/2 {
		if g.rightPaddle.y < screenHeight-paddleHeight {
			g.rightPaddle.y += 2
		}
	}

	// update the ball
	g.ball.x += g.ball.dx
	g.ball.y += g.ball.dy

	// press escape to quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	return nil
}

func (g *game) initPosition() {
	g.ball.x = screenWidth / 2
	g.ball.y = screenHeight / 2
	g.ball.dx = 2
	g.ball.dy = 2
	g.leftPaddle.y = screenHeight/2 - paddleHeight/2
	g.rightPaddle.y = screenHeight/2 - paddleHeight/2
	g.scoreLeft, g.scoreRight = 0, 0
}

func (g *game) Draw(screen *ebiten.Image) {
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	// Draw background
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, black, false)

	// Draw ball
	// vector.DrawFilledCircle(screen, g.ball.x, g.ball.y, ballRadius, black, false)
	vector.DrawFilledCircle(screen, float32(g.ball.x), float32(g.ball.y), ballRadius, white, false)

	// Draw paddles
	// vector.DrawFilledRect(screen, &ebiten.GeoM{Rect: ebiten.Rect(0, int(g.leftPaddle.y), paddleWidth, paddleHeight)}, black)
	vector.DrawFilledRect(screen, g.leftPaddle.x, g.leftPaddle.y, paddleWidth, paddleHeight, white, false)
	// vector.DrawFilledRect(screen, &ebiten.GeoM{Rect: ebiten.Rect(screenWidth-paddleWidth, int(g.rightPaddle.y), paddleWidth, paddleHeight)}, black)
	vector.DrawFilledRect(screen, g.rightPaddle.x, g.rightPaddle.y, paddleWidth, paddleHeight, white, false)

	// Draw score
	// text.Draw(screen, "Score Left: "+string(g.scoreLeft), 10, 20, font, black)
	// text.Draw(screen, "Score Right: "+string(g.scoreRight), screenWidth-100, 20, font, black)
	t := "Score Left: " + strconv.Itoa(g.scoreLeft) + " Score Right: " + strconv.Itoa(g.scoreRight)
	ebitenutil.DebugPrintAt(screen, t, 230, 20)
}

func main() {
	g := &game{}
	g.scoreLeft, g.scoreRight = 0, 0
	initObject(&g.ball, screenWidth/2, screenHeight/2, ballRadius, ballRadius)
	g.ball.dx = rand.Float32()*2 + 1
	g.ball.dy = rand.Float32()*2 + 1
	initObject(&g.leftPaddle, 0, 0, paddleWidth, paddleHeight)
	initObject(&g.rightPaddle, screenWidth-paddleWidth, 0, paddleWidth, paddleHeight)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func (g *game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}
