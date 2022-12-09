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

		row := make([]uint8, len(cells))

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
	var vismap [][]uint8 // 0 = not visible, 1 = visible

	height, width := len(grid), len(grid[0])

	for y := 0; y < height; y++ {
		vrow := make([]uint8, width)
		vismap = append(vismap, vrow)

		vrow[0], vrow[width-1] = 1, 1 // left and right edge

		if y > 0 && y < height-1 {
			continue
		}

		for x := 1; x < width-1; x++ { // top and bottom edge
			vrow[x] = 1
		}
	}

	for y := 1; y < height-1; y++ {
		row := grid[y]
		vrow := vismap[y]

		left, right := row[0], row[width-1]

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
		top, bottom := grid[0][x], grid[height-1][x]

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
	total := uint(0)

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
