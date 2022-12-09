package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parse(file *os.File) ([][]uint8, error) {
	var grid [][]uint8

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		cells := strings.Split(line, "")
		size := len(cells)

		row := make([]uint8, size, size)

		for i, char := range cells {
			value, err := strconv.ParseUint(char, 10, 8)

			if err != nil {
				return nil, err
			}

			row[i] = uint8(value)
		}

		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return grid, nil
}

func visibility(grid [][]uint8) [][]uint8 {
	var vismap [][]uint8

	height := len(grid)
	width := len(grid[0]) // assume at least 1 row and rows of equal width

	for y := 0; y < height; y++ {
		vrow := make([]uint8, width, width)
		vismap = append(vismap, vrow)

		vrow[0] = 1
		vrow[width-1] = 1

		if y == 0 || y == height-1 {
			for x := 1; x < width-1; x++ {
				vrow[x] = 1
			}
		}
	}

	for y := 1; y < height-1; y++ {
		row := grid[y]
		vrow := vismap[y]

		left := row[0]
		right := row[width-1]

		for x, xx := 1, width-2; x < width-1 && xx > 0; x, xx = x+1, xx-1 {
			if row[x] > left {
				vrow[x] = 1
				left = row[x]
			}

			if row[xx] > right {
				vrow[xx] = 1
				right = row[xx]
			}

			if x >= xx && left == right {
				break
			}
		}
	}

	for x := 1; x < width-1; x++ {
		top := grid[0][x]
		bottom := grid[height-1][x]

		for y, yy := 1, height-2; y < height-1 && yy > 0; y, yy = y+1, yy-1 {
			if row, vrow := grid[y], vismap[y]; row[x] > top {
				vrow[x] = 1
				top = row[x]
			}

			if row, vrow := grid[yy], vismap[yy]; row[x] > bottom {
				vrow[x] = 1
				bottom = row[x]
			}

			if y >= yy && top == bottom {
				break
			}
		}
	}

	return vismap
}

func sum(grid [][]uint8) uint {
	var total uint

	for _, row := range grid {
		for _, cell := range row {
			total += uint(cell)
		}
	}

	return total
}

func main() {
	grid, err := parse(os.Stdin)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
	}

	vismap := visibility(grid)
	count := sum(vismap)

	fmt.Println(count)
}
