package main

import "github.com/nsf/termbox-go"
import "time"
import "fmt"

const shipColour = termbox.ColorRed
const spaceShip = 'Ÿ'
const levelDurationInSeconds = time.Millisecond * 1000 * 10

func drawMap(g *Game) {
	for y := 0; y < len(g.gameMap); y++ {
		for x := 0; x < g.screenWidth; x++ {
			if x < g.gameMap[y].left || x > g.gameMap[y].right {
				termbox.SetCell(x, y, ' ', termbox.ColorBlue, termbox.ColorBlue)
			}
		}
	}
}

func drawShip(g *Game) {
	termbox.SetCell(g.playerX, g.playerY, spaceShip, shipColour, termbox.ColorDefault)
}

func drawMessage(g *Game, message string, colour termbox.Attribute) {
	tbprint(g.screenWidth/2-len(message)/2, g.screenHeight/2-1, colour, termbox.ColorDefault, message)
}

func drawWidget(g *Game) {
	tbprint(0, 0, termbox.ColorYellow, termbox.ColorDefault, fmt.Sprintf("Level: %d", g.level))
	tbprint(0, 1, termbox.ColorYellow, termbox.ColorDefault, fmt.Sprintf("Depth: %d meters", g.getDepth()))
}

func render(g *Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawMap(g)
	drawShip(g)
	drawWidget(g)

	if !g.playing {
		drawMessage(g, "GAME OVER!", termbox.ColorRed)
	} else if g.paused {
		drawMessage(g, "PAUSED", termbox.ColorYellow)
	}
	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	w, h := termbox.Size()

	g := NewGame(w, h)
	ticker := time.NewTicker(levelDurationInSeconds)

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch {
				case ev.Key == termbox.KeyArrowLeft:
					g.moveLeft()
				case ev.Key == termbox.KeyArrowRight:
					g.moveRight()
				case ev.Key == termbox.KeyArrowUp:
					g.moveUp()
				case ev.Key == termbox.KeyArrowDown:
					g.moveDown()
				case ev.Ch == 'p':
					g.togglePause()
				case ev.Ch == 'q':
					return
				}
			}
		case <-ticker.C:
			g.incrementLevel()
		default:
			render(g)
			g.animate()

			time.Sleep(g.getAnimationLoopDelay())
		}
	}
}
