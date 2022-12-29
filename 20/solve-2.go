package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"os"
	"strconv"
)

const key int = 811589153
const rounds int = 10

func parse(file *os.File) ([]int, error) {
	var nums []int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.ParseInt(line, 10, 32)

		if err != nil {
			return nil, err
		}

		nums = append(nums, int(num)*key)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nums, nil
}

func mix(r *ring.Ring) *ring.Ring {
	l := r.Len()
	elements := make([]*ring.Ring, l)

	for i := range elements {
		elements[i] = r
		r = r.Next()
	}

	for i := 0; i < rounds; i++ {
		for _, element := range elements {
			p := element.Prev()
			e := p.Unlink(1)

			value := (e.Value).(int)
			p = p.Move(value % (l - 1))

			p.Link(e)
		}
	}

	return r
}

func main() {
	nums, err := parse(os.Stdin)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		os.Exit(1)
	}

	l := len(nums)

	if l == 0 {
		fmt.Fprintln(os.Stderr, "No numbers read")
		os.Exit(1)
	}

	r := ring.New(l)

	// Copy values to ring
	for _, n := range nums {
		r.Value = n
		r = r.Next()
	}

	mix(r)

	// Rotate ring to 0-element
	for start := r; r.Value != 0 && r.Next() != start; r = r.Next() {
		// …
	}

	if r.Value != 0 {
		fmt.Fprintln(os.Stderr, "0 not found in mixed numbers")
		os.Exit(2)
	}

	// Find x, y and z by moving along the ring

	r = r.Move(1000)
	x := (r.Value).(int)

	r = r.Move(1000)
	y := (r.Value).(int)

	r = r.Move(1000)
	z := (r.Value).(int)

	fmt.Println(x + y + z)
}
