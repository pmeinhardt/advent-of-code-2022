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

func score(grid [][]uint8, px, py int) uint {
	top, right, bottom, left := uint(0), uint(0), uint(0), uint(0)

	height, width := len(grid), len(grid[0])

	value := grid[py][px]

	for x := px - 1; x >= 0; x-- {
		left += 1

		if grid[py][x] >= value {
			break
		}
	}

	for x := px + 1; x < width; x++ {
		right += 1

		if grid[py][x] >= value {
			break
		}
	}

	for y := py - 1; y >= 0; y-- {
		top += 1

		if grid[y][px] >= value {
			break
		}
	}

	for y := py + 1; y < height; y++ {
		bottom += 1

		if grid[y][px] >= value {
			break
		}
	}

	return top * right * bottom * left
}

func survey(grid [][]uint8) [][]uint {
	var scores [][]uint

	height, width := len(grid), len(grid[0])

	for _, row := range grid {
		scores = append(scores, make([]uint, len(row)))
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			scores[y][x] = score(grid, x, y)
		}
	}

	return scores
}

func max(grid [][]uint) uint {
	result := uint(0) // no negative values allowed

	for _, row := range grid {
		for _, cell := range row {
			if result < cell {
				result = cell
			}
		}
	}

	return result
}

func main() {
	grid, err := parse(os.Stdin)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
	}

	scores := survey(grid)
	value := max(scores)

	fmt.Println(value)
}
