package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

var mp map[int]bool = map[int]bool{}

func genRandValue() int {
	res := 2
	rnd := rand.Intn(10)
	if rnd == 4 {
		res = 4
	} else {
		res = 2
	}
	return res
}

func generate(box *[][]int, dir int) {
	n, m, ni, mi := len(*box), len((*box)[0]), 0, 0
	save := make(map[int]int)
	for {
		for save[mi] == ni {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			ni = rnd.Intn(n)
			if dir == 65514 || dir == 65516 {
				mi = rnd.Intn(m/2 + m%2)
			} else {
				mi = m/2 + rnd.Intn(m/2+m%2)
			}
		}
		save[mi] = ni

		if (*box)[ni][mi] == 0 {
			(*box)[ni][mi] = genRandValue()
			break
		}
	}
}

func reverse(line *[]int, dir int) {
	if dir == 65514 || dir == 65516 {
		for i := len(*line)/2 - 1; i >= 0; i-- {
			opp := len(*line) - 1 - i
			(*line)[i], (*line)[opp] = (*line)[opp], (*line)[i]
		}
	}
}

func shiftLine(line *[]int, dir int) bool {
	var state, flag bool
	temp := make([]int, len(*line))
	reverse(line, dir)
	for j, i := 0, 0; i < len(*line); i++ {
		if (*line)[i] != temp[j] && (*line)[i] != 0 && temp[j] != 0 || flag {
			j++
			flag = false
		} else if (*line)[i] != 0 && temp[j] != 0 {
			flag = true
		}
		temp[j] += (*line)[i]
		if (*line)[j] != temp[j] && !state {
			state = true
		}
	}
	*line = temp
	reverse(line, dir)
	return state
}

func shiftLines(box *[][]int, dir int) (bool, bool) {
	var state, sl bool
	transposition(box, dir)
	for i := range *box {
		sl = shiftLine(&(*box)[i], dir)
		state = state || sl
	}
	if state {
		generate(box, dir)
		mp = make(map[int]bool)
	} else {
		mp[dir] = state
	}
	transposition(box, dir)
	if len(mp) == 4 {
		return state, false
	} else {
		return state, true
	}

}

func transposition(box *[][]int, dir int) {
	if dir < 65516 {
		return
	}
	n, m := len(*box), len((*box)[0])
	tempbox := make([][]int, m)

	for i := 0; i < m; i++ {
		tempbox[i] = make([]int, n)
		for j := 0; j < n; j++ {
			tempbox[i][j] = (*box)[j][i]
		}
	}
	*box = tempbox
}

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func showBox(box *[][]int) {
	max := 0
	ln := len(*box)
	lm := len((*box)[0])
	fmt.Print("\033[H\033[2J")
	for i := 0; i < ln; i++ {
		fmt.Println(strings.Repeat("+------", lm) + "+")
		for j := 0; j < lm; j++ {
			if (*box)[i][j] == 0 {
				fmt.Printf("|%6s", "")
			} else {
				if (*box)[i][j] > max {
					max = (*box)[i][j]
				}
				fmt.Print("|" + center(fmt.Sprintf("%d", (*box)[i][j]), 6))
			}
		}
		fmt.Printf("|\n")
	}
	fmt.Println(strings.Repeat("+------", lm) + "*")
	fmt.Println("|" + center("Score: "+strconv.Itoa(max), lm*7-1) + "|")
}

func initBox(n, m int) *[][]int {
	box := make([][]int, n)
	for i := range box {
		box[i] = make([]int, m)
	}
	return &box
}

func play2048() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Введите размер поля в формате: Х У, например, 4 4:")
	var n, m int
	_, err := fmt.Scan(&n, &m)
	if err != nil {
		fmt.Println("Возникла ошибка")
		return
	}
	var gameover bool
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	box := initBox(n, m)
	generate(box, 0)
	showBox(box)

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		_, gameover = shiftLines(box, int(event.Key))
		showBox(box)

		if !gameover {
			lm := len((*box)[0])
			fmt.Println("|" + center("-= Game Over =-", lm*7-1) + "|")
			return
		}

		if event.Key == keyboard.KeyEsc {
			return
		}
	}
}

func main() {
	play2048()
}
