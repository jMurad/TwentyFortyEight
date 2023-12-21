package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	score      int
	length     int
	gameover   bool
	shifts     map[int]bool
	shiftState bool
}

func (tfe *TwentyFortyEight) Init() {
	tfe.length = 6
	tfe.shifts = make(map[int]bool)
}

func (tfe *TwentyFortyEight) genRandValue() int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := 2
	genNumber := rnd.Intn(10)
	if genNumber == 4 {
		res = 4
	}
	return res
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

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func (tfe *TwentyFortyEight) showBox() {
	fmt.Print("\033[H\033[2J")
	border := fmt.Sprint("+" + strings.Repeat("-", tfe.length))
	for i := 0; i < tfe.n; i++ {
		fmt.Println(strings.Repeat(border, tfe.m) + "+")
		for j := 0; j < tfe.m; j++ {
			if tfe.box[i][j] == 0 {
				fmt.Printf("|%6s", "")
			} else {
				if tfe.box[i][j] > tfe.score {
					tfe.score = tfe.box[i][j]
				}
				fmt.Print("|" + center(fmt.Sprintf("%d", tfe.box[i][j]), tfe.length))
			}
		}
		fmt.Printf("|\n")
	}
	fmt.Println(strings.Repeat(border, tfe.m) + "*")
	fmt.Println("|" + center("Score: "+strconv.Itoa(tfe.score), tfe.m*(tfe.length+1)-1) + "|")
}

func (tfe *TwentyFortyEight) initBox() {
	tfe.box = make([][]int, tfe.n)
	for i := range tfe.box {
		tfe.box[i] = make([]int, tfe.m)
	}
}

func (tfe *TwentyFortyEight) play2048() {
	fmt.Println("\033[H\033[2JВведите размер поля в формате: Х У, например, 4 4:")
	_, err := fmt.Scan(&tfe.n, &tfe.m)
	if err != nil {
		fmt.Println("Возникла ошибка")
		return
	}

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

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		tfe.direction = int(event.Key)
		tfe.shiftLines()
		tfe.showBox()

		if tfe.gameover {
			fmt.Println("|" + center("-= Game Over =-", tfe.m*(tfe.length+1)-1) + "|")
			return
		}

		if event.Key == keyboard.KeyEsc {
			return
		}
	}
}

func main() {
	tfe := TwentyFortyEight{}
	tfe.Init()
	tfe.play2048()
}
