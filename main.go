package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	right = iota + 65514
	left
	down
	up
)

type TwentyFortyEight struct {
	n, m       int
	box        [][]int
	direction  int
	length     int
	gameover   bool
	shifts     map[int]bool
	shiftState bool
	Player
}
type Player struct {
	name    string
	score   int
	level   int
	sizeBox string
	rating  int
}

func (tfe *TwentyFortyEight) Init() {
	tfe.length = 5
	tfe.shifts = make(map[int]bool)
}

func (tfe *TwentyFortyEight) genRandValue() (res int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	res = 2
	genNumber := rnd.Intn(10)
	if genNumber == 4 {
		res = 4
	}
	return
}

func (tfe *TwentyFortyEight) generate() {
	ni, mi := 0, 0
	save := make(map[int]int)
	for {
		for save[mi] == ni {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			ni = rnd.Intn(tfe.n)
			if tfe.direction == right || tfe.direction == down {
				mi = rnd.Intn(tfe.m/2 + tfe.m%2)
			} else {
				mi = tfe.m/2 + rnd.Intn(tfe.m/2+tfe.m%2)
			}
		}
		save[mi] = ni

		if tfe.box[ni][mi] == 0 {
			tfe.box[ni][mi] = tfe.genRandValue()
			break
		}
	}
}

func (tfe *TwentyFortyEight) reverse(ind int) {
	switch tfe.direction {
	case right, down:
		for i := tfe.m/2 - 1; i >= 0; i-- {
			opp := tfe.m - 1 - i
			tfe.box[ind][i], tfe.box[ind][opp] = tfe.box[ind][opp], tfe.box[ind][i]
		}
	}
}

func (tfe *TwentyFortyEight) shiftLine(ind int) {
	var flag bool
	temp := make([]int, tfe.m)
	tfe.reverse(ind)
	for j, i := 0, 0; i < tfe.m; i++ {
		if tfe.box[ind][i] != temp[j] && tfe.box[ind][i] != 0 && temp[j] != 0 || flag {
			j++
			flag = false
		} else if tfe.box[ind][i] != 0 && temp[j] != 0 {
			flag = true
			tfe.score += tfe.box[ind][i] * 2
		}
		temp[j] += tfe.box[ind][i]
		if tfe.box[ind][j] != temp[j] && !tfe.shiftState {
			tfe.shiftState = true
		}
	}
	tfe.box[ind] = temp
	tfe.reverse(ind)
}

func (tfe *TwentyFortyEight) transposition() {
	switch tfe.direction {
	case down, up:
		tempbox := make([][]int, tfe.m)

		for i := 0; i < tfe.m; i++ {
			tempbox[i] = make([]int, tfe.n)
			for j := 0; j < tfe.n; j++ {
				tempbox[i][j] = tfe.box[j][i]
			}
		}
		tfe.box = tempbox
	case right, left:
		return
	}
}

func (tfe *TwentyFortyEight) shiftLines() {
	tfe.shiftState = false
	tfe.transposition()
	for i := range tfe.box {
		tfe.shiftLine(i)
	}
	if tfe.shiftState {
		tfe.generate()
		tfe.shifts = make(map[int]bool)
	} else {
		tfe.shifts[tfe.direction] = tfe.shiftState
	}
	tfe.transposition()
	if len(tfe.shifts) == 4 {
		tfe.gameover = true
	} else {
		tfe.gameover = false
	}
}

func (tfe *TwentyFortyEight) showBox() {
	fmt.Print("\033[H\033[2J")
	border := fmt.Sprint("+" + repeat("-", tfe.length))
	for i := 0; i < tfe.n; i++ {
		fmt.Println(repeat(border, tfe.m) + "+")
		for j := 0; j < tfe.m; j++ {
			if tfe.box[i][j] == 0 {
				fmt.Print("|" + center("", tfe.length))
			} else {
				number := fmt.Sprintf("%d", tfe.box[i][j])
				fmt.Printf("|" + center(number, tfe.length))
			}
		}
		fmt.Print("|\n")
	}
	fmt.Println(repeat(border, tfe.m) + "*")

	numScore := fmt.Sprintf("Score: %d", tfe.score)
	fmt.Println("|" + center(numScore, tfe.m*(tfe.length+1)-1) + "|")
}

func (tfe *TwentyFortyEight) initBox() {
	tfe.box = make([][]int, tfe.n)
	for i := range tfe.box {
		tfe.box[i] = make([]int, tfe.m)
	}
}

func (tfe *TwentyFortyEight) registration() {
	fmt.Print("Введие ваше имя: ")
	_, err := fmt.Scanf("%s", &tfe.name)
	if err != nil {
		fmt.Println("Возникла ошибка")
		return
	}

	fmt.Print("Введите размер поля в формате Х У:")
	_, err = fmt.Scan(&tfe.n, &tfe.m)
	if err != nil {
		fmt.Println("Возникла ошибка")
		return
	}
}

func (tfe *TwentyFortyEight) play2048() {
	fmt.Println("\033[H\033[2J")

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	tfe.initBox()

	tfe.generate()
	tfe.showBox()

	for event := range keysEvents {
		if event.Err != nil {
			panic(event.Err)
		}
		switch key := int(event.Key); key {
		case right, left, down, up:
			tfe.direction = key
			tfe.shiftLines()
			tfe.showBox()
		case int(keyboard.KeyEsc):
			return
		}
		if tfe.gameover {
			fmt.Println("|" + center("-= Game Over =-", tfe.m*(tfe.length+1)-1) + "|")
			return
		}
	}
}

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func repeat(s string, n int) (res string) {
	switch n {
	case 0:
		return ""
	case 1:
		return s
	}

	for i := 0; i < n; i++ {
		res += s
	}
	return
}

func main() {
	tfe := TwentyFortyEight{}
	tfe.Init()
	tfe.play2048()

}
