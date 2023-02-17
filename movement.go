package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

type Color struct {
	r int
	g int
	b int
}

func setColor(c Color) {
	fmt.Printf("\033[38;2;%v;%v;%vm", c.r, c.g, c.b)
}

func setPos(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x)
}

func randColor() Color {
	return Color{
		r: rand.Intn(256),
		g: rand.Intn(256),
		b: rand.Intn(256),
	}
}

func blendColors(c1, c2 Color, amount, totalAmount int) Color {
	rDif := int(math.Round(float64(c2.r-c1.r) * (float64(amount) / float64(totalAmount))))
	gDif := int(math.Round(float64(c2.g-c1.g) * (float64(amount) / float64(totalAmount))))
	bDif := int(math.Round(float64(c2.b-c1.b) * (float64(amount) / float64(totalAmount))))
	return Color{
		r: c1.r + rDif,
		g: c1.g + gDif,
		b: c1.b + bDif,
	}
}

func genCharFromDir(dir int) string {
	if dir%2 == 1 {
		return "│"
	}
	return "─"
}

func isNearEdge(x, y, dir int) bool {
	xsize, ysize, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Panic(err)
	}
	if dir == 1 && y < 3 {
		return true
	}
	if dir == 3 && y > (ysize-4) {
		return true
	}
	if dir == 2 && x < 3 {
		return true
	}
	if dir == 0 && x > (xsize-4) {
		return true
	}
	return false
}

func genDir(x, y, dir int) (int, string) {
	xsize, ysize, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Panic(err)
	}
	newDir := dir

	if isNearEdge(x, y, dir) || rand.Intn(20) == 0 {
		newDir = (dir + [2]int{1, 3}[rand.Intn(2)]) % 4
	} else if rand.Intn(20) == 0 {
		newDir = (dir + [2]int{1, 3}[rand.Intn(2)]) % 4
		if y < 3 && (dir == 0 || dir == 2) {
			newDir = 3
		}
		if y > ysize && (dir == 0 || dir == 2) {
			newDir = 3
		}
		if x < 3 && (dir == 1 || dir == 3) {
			newDir = 0
		}
		if x > xsize && (dir == 1 || dir == 3) {
			newDir = 2
		}
	}

	if dir == newDir {
		return newDir, genCharFromDir(dir)
	}
	// TODO: Make this cleaner
	if (dir == 2 && newDir == 1) || (dir == 3 && newDir == 0) {
		return newDir, "└"
	}
	if (dir == 0 && newDir == 1) || (dir == 3 && newDir == 2) {
		return newDir, "┘"
	}
	if (dir == 0 && newDir == 3) || (dir == 1 && newDir == 2) {
		return newDir, "┐"
	}
	return newDir, "┌"
}

func movingLine(startx, starty int, startColor Color, cloneProbability, deathProbability int) {
	x, y := startx, starty
	direction := rand.Intn(4)
	destinationColor := randColor()
	curProgress := 0
	maxProgess := rand.Intn(30) + 10
	char := genCharFromDir(direction)

	for {
		setPos(x, y)
		curColor := blendColors(startColor, destinationColor, curProgress, maxProgess)
		setColor(curColor)

		direction, char = genDir(x, y, direction)

		curProgress++
		if curProgress == maxProgess {
			startColor = destinationColor
			destinationColor = randColor()
			curProgress = 0
			maxProgess = rand.Intn(30) + 10
		}

		fmt.Print(char)
		switch direction {
		case 0:
			x += 1
		case 1:
			y -= 1
		case 2:
			x -= 1
		case 3:
			y += 1
		}
		time.Sleep(time.Second / 100)
	}
}

func main() {
	fmt.Print("\033[?25l")
	for {
		movingLine(50, 20, randColor(), 25, 60)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
