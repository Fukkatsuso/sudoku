package sudoku

import (
	"fmt"
	"math/rand"
	"time"
)

// Sudoku9x9 is standard 9x9 sized sudoku
type Sudoku9x9 struct {
	Table [][]int
}

// NewSudoku9x9 create a new 9x9 sized puzzle with 'hole' holes
func NewSudoku9x9(hole int) *Sudoku9x9 {
	s := new(Sudoku9x9)
	s.Init()
	// 一度全部解いてから穴を空ける
	s.Solve()
	s.makeHole(hole)
	return s
}

// Init table
func (s *Sudoku9x9) Init() {
	s.Table = make([][]int, 9)
	for i := range s.Table {
		s.Table[i] = make([]int, 9)
		for j := range s.Table[i] {
			if i == 0 {
				s.Table[i][j] = j + 1
			} else {
				s.Table[i][j] = 0
			}
		}
	}

	shuffleSlice(s.Table[0])
}

// Solve 9x9 puzzle and return whether the puzzle is solved or not
func (s *Sudoku9x9) Solve() bool {
	return s.solve(0)
}

// Judge whether the answer is correct or not and return a truth value
func (s Sudoku9x9) Judge() bool {
	n := 9
	// existsX[i][j] : i番目のXに(j+1)が含まれるか
	existsRow := newBoolTable(n, n)
	existsCol := newBoolTable(n, n)
	existsBlock := newBoolTable(n, n)

	for index := 0; index < n*n; index++ {
		row, col := index/n, index%n
		bIndex := (row/3)*3 + col/3
		num := s.Table[row][col] - 1
		if num < 0 || n <= num {
			return false
		} else if existsRow[row][num] || existsCol[col][num] || existsBlock[bIndex][num] {
			return false
		}
		existsRow[row][num] = true
		existsCol[col][num] = true
		existsBlock[bIndex][num] = true
	}

	return true
}

// Print 9x9 table
func (s Sudoku9x9) Print() {
	for i := range s.Table {
		for j := range s.Table[i] {
			fmt.Printf(" %d", s.Table[i][j])
		}
		fmt.Printf("\n")
	}
}

// Solveのサブルーチン
func (s *Sudoku9x9) solve(index int) bool {
	// 埋め終わり
	if index == 9*9 {
		return true
	}

	row, col := index/9, index%9
	// すでに埋まっているなら次に進む
	if s.Table[row][col] > 0 {
		return s.solve(index + 1)
	}
	// 解候補は仮埋めして進める
	c := s.checkNumber(index)
	for i := 1; i <= 9; i++ {
		if !c[i-1] {
			continue
		}
		s.Table[row][col] = i
		if s.solve(index + 1) {
			return true
		}
	}

	// 1つも置けない場合は白紙に戻す
	s.Table[row][col] = 0
	return false
}

// 各数が置けるかどうか
func (s Sudoku9x9) checkNumber(index int) []bool {
	c := make([]bool, 9)
	for i := range c {
		c[i] = true
	}

	// ヨコ
	row := index / 9
	for i := 0; i < 9; i++ {
		if num := s.Table[row][i]; num > 0 {
			c[num-1] = false
		}
	}
	// タテ
	col := index % 9
	for i := 0; i < 9; i++ {
		if num := s.Table[i][col]; num > 0 {
			c[num-1] = false
		}
	}
	// ブロック
	bRow, bCol := row/3, col/3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			row := bRow*3 + i
			col := bCol*3 + j
			if num := s.Table[row][col]; num > 0 {
				c[num-1] = false
			}
		}
	}

	return c
}

// hole個の穴を空ける
func (s *Sudoku9x9) makeHole(hole int) {
	rand.Seed(time.Now().UnixNano())

	var row, col int
	for holeCnt, index := 0, rand.Intn(9*9); holeCnt < hole; {
		if holeCnt%2 == 0 {
			row, col = index/9, index%9
		} else {
			// 点対称の位置
			row, col = 9-1-index/9, 9-1-index%9
		}

		if s.Table[row][col] > 0 {
			s.Table[row][col] = 0
			holeCnt++
		} else {
			index = rand.Intn(9 * 9)
		}
	}
}

// row行col列の2次元bool配列を返す
func newBoolTable(row, col int) [][]bool {
	t := make([][]bool, row)
	for i := 0; i < row; i++ {
		t[i] = make([]bool, col)
	}
	return t
}

// Fisher-Yates shuffle
func shuffleSlice(s []int) {
	rand.Seed(time.Now().UnixNano())

	for i := len(s) - 1; i > 0; i-- {
		j := int(rand.Float64() * float64(i+1))
		s[i], s[j] = s[j], s[i]
	}
}
