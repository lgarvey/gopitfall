package main

import "math/rand"
import "time"

const animationSpeed = 10 * time.Millisecond
const minStep = 1
const maxStep = 6
const maxCaveWidth = 50
const minCaveWidth = 50
const maxLevel = 10
const playerStartY = 10

type GameRow struct {
	left  int
	right int
}

type Game struct {
	screenWidth  int
	screenHeight int
	playing      bool
	paused       bool
	level        int
	depth        uint
	playerX      int
	playerY      int

	gameMap []GameRow
}

func NewGame(screenWidth, screenHeight int) *Game {
	g := new(Game)
	g.init(screenWidth, screenHeight)
	return g
}

func (g *Game) init(screenWidth, screenHeight int) {
	g.screenWidth = screenWidth
	g.screenHeight = screenHeight
	g.playing = true
	g.paused = false
	g.level = 1
	g.depth = 0
	g.initMap()
	g.initPlayer()
}

func (g *Game) initMap() {
	g.gameMap = []GameRow{}

	startRow := GameRow{
		left:  g.screenWidth/2 - maxCaveWidth/2,
		right: g.screenWidth/2 + maxCaveWidth/2}

	g.gameMap = append(g.gameMap, startRow)

	for y := 1; y < g.screenHeight; y++ {
		g.gameMap = append(g.gameMap, g.genNextRow(g.gameMap[y-1]))
	}
}

func (g *Game) initPlayer() {
	g.playerY = playerStartY
	row := g.gameMap[g.playerY]
	g.playerX = row.left + ((row.right - row.left) / 2)
}

func (g *Game) genNextRow(prevRow GameRow) GameRow {

	var newLeft, newRight int

	direction := rand.Intn(3)
	step := minStep + rand.Intn(maxStep-minStep+1)

	switch direction {
	case 0: // no moveement
		newLeft = prevRow.left
		newRight = prevRow.right
		break

	case 1: // left
		newLeft = prevRow.left - step
		newRight = prevRow.right - step
		break

	case 2: // right
		newLeft = prevRow.left + step
		newRight = prevRow.right + step
	}

	if newLeft < 0 {
		newLeft = 0
	}

	if newRight < minCaveWidth {
		newRight = minCaveWidth
	}

	if newRight > g.screenWidth {
		newRight = g.screenWidth
	}

	if newLeft > g.screenWidth-minCaveWidth {
		newLeft = g.screenWidth - minCaveWidth
	}

	return GameRow{newLeft, newRight}
}

func (g *Game) scrollMap() {
	g.gameMap = g.gameMap[1:]
	newRow := g.genNextRow(g.gameMap[len(g.gameMap)-1])
	g.gameMap = append(g.gameMap, newRow)
}

func (g *Game) isCollision() bool {
	row := g.gameMap[g.playerY]

	if g.playerX > row.left && g.playerX < row.right {
		return false
	} else {
		return true
	}
}

func (g *Game) animate() {
	if g.isPlaying() {
		if g.isCollision() {
			g.playing = false
		} else {
			g.scrollMap()
		}
	}
}

func (g *Game) moveUp() {
	if g.isPlaying() {
		g.playerY -= 2

		if g.playerY < 0 {
			g.playerY = 0
		}
	}
}

func (g *Game) moveDown() {
	if g.isPlaying() {
		g.playerY += 2

		if g.playerY > g.screenHeight {
			g.playerY = g.screenHeight
		}
	}
}

func (g *Game) moveLeft() {
	if g.isPlaying() {
		g.playerX -= 4

		if g.playerX < 0 {
			g.playerX = 0
		}
	}
}

func (g *Game) moveRight() {
	if g.isPlaying() {
		g.playerX += 4

		if g.playerX > g.screenWidth {
			g.playerX = g.screenWidth
		}
	}
}

func (g *Game) togglePause() {
	if g.playing {
		g.paused = !g.paused
	}
}

func (g *Game) isPlaying() bool {
	return g.playing && !g.paused
}

func (g *Game) incrementLevel() {
	if g.level < maxLevel {
		g.level++
	}
}

func (g *Game) getAnimationLoopDelay() time.Duration {
	return time.Duration(20+(maxLevel-g.level)*10) * time.Millisecond
}
