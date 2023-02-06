package myers

import (
	"bytes"
	"fmt"
	"sort"
)

type point struct {
	X int
	Y int
}

type step struct {
	from point
	to   point
}

func (s step) String() string {
	return fmt.Sprintf("(%d, %d) -> (%d, %d)", s.from.X, s.from.Y, s.to.X, s.to.Y)
}

type editList []int

func (e editList) get(idx int) int {
	if idx < 0 {
		return e[len(e)+idx]
	}
	return e[idx]
}

type operation int

const (
	ins operation = iota
	del
	eql
)

type lineDiff struct {
	op operation
	a  string
}

func (l lineDiff) String() string {
	buf := bytes.Buffer{}
	switch l.op {
	case del:
		buf.WriteString("\033[31m")
		buf.WriteString(l.a)
		buf.WriteString("\033[0m")
	case ins:
		buf.WriteString("\033[32m")
		buf.WriteString(l.a)
		buf.WriteString("\033[0m")
	case eql:
		buf.WriteString(l.a)
	}

	return buf.String()
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Diff between two strings (first, second) using Myer's Algorithm.
// Implemented based on the excellent blog at https://blog.jcoglan.com/2017/02/12/the-myers-diff-algorithm-part-1/
// And the original paper "An O(ND) Difference Algorithm and Its Variations" by Eugene W. Myer
// See: https://link.springer.com/article/10.1007/BF01840446
func Diff(first, second string) (string, bool) {
	buf := bytes.Buffer{}

	lhs := []byte(first)
	rhs := []byte(second)
	steps, err := backtrack(lhs, rhs)
	if err != nil {
		return "", false
	}

	unequal := false

	sort.SliceStable(steps, func(i, j int) bool {
		return true
	})
	prevA := 0
	for _, s := range steps {
		// fmt.Printf("%s (prevA=%d) - ", s.String(), prevA)
		if s.to.X == s.from.X {
			unequal = true
			// inserted, rhs
			buf.WriteString(lineDiff{eql, string(rhs[min(prevA, s.to.Y):s.to.Y])}.String())
			inserted := string(rhs[s.from.Y:s.to.Y])
			buf.WriteString(lineDiff{ins, inserted}.String())
			prevA = s.to.Y + 1
			// fmt.Printf("INSERT: char %d (rhs)[%s]\n", s.to.Y, buf.String())
		} else if s.to.Y == s.from.Y {
			unequal = true
			// deleted, rhs
			buf.WriteString(lineDiff{eql, string(rhs[min(prevA, s.from.X):s.from.X])}.String())
			remove := string(lhs[s.from.X:s.to.X])
			buf.WriteString(lineDiff{del, remove}.String())
			prevA = s.from.X + 1
			// fmt.Printf("DELETE: char %d (lhs)[%s]\n", s.from.X, buf.String())
		} else {
			// x and y both change
			buf.WriteString(lineDiff{eql, string(rhs[s.from.Y:s.to.Y])}.String())
			prevA = s.to.Y + 1
			// fmt.Printf("EQUAL: from %d to %d (lhs)[%s]\n", s.from.X, s.to.X, buf.String())
		}
	}

	if unequal && prevA < len(rhs) {
		buf.WriteString(string(rhs[prevA:]))
	}

	if unequal {
		return buf.String(), true
	}

	return "", false
}

func backtrack(lhs, rhs []byte) ([]step, error) {
	edits, err := ses(lhs, rhs)
	if err != nil {
		return nil, err
	}

	x := len(lhs)
	y := len(rhs)
	steps := make([]step, 0)
	// traverse in reverse
	for d := len(edits) - 1; d >= 0; d-- {
		v := edits[d]
		k := x - y

		var prevK int
		if k == -d || (k != d && (v.get(k-1) < v.get(k+1))) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		prevX := v.get(prevK)
		prevY := prevX - prevK

		for x > prevX && y > prevY {
			steps = append(steps, step{
				from: point{x - 1, y - 1},
				to:   point{x, y},
			})
			x--
			y--
		}
		if d > 0 {
			steps = append(steps, step{
				from: point{prevX, prevY},
				to:   point{x, y},
			})
		}
		x = prevX
		y = prevY
	}
	return steps, nil
}

// ses (Shorted Edit Search) is a graph search
func ses(lhs, rhs []byte) ([]editList, error) {
	var x int
	n := len(lhs)
	m := len(rhs)
	maxLen := n + m

	var v editList = make([]int, 2*maxLen+1)
	trace := make([]editList, 0)
	for d := 0; d <= maxLen; d++ {

		// s.g. trace << v.clone from blog post
		current := make([]int, len(v))
		copy(current, v)
		trace = append(trace, current)

		for k := -d; k <= d; k += 2 {
			preK := v.get(k - 1)
			postK := v.get(k + 1)
			if k == -d || (k != d && preK < postK) {
				x = postK
			} else {
				x = preK + 1
			}
			y := x - k

			for x < n && y < m && lhs[x] == rhs[y] {
				x++
				y++
			}

			if k < 0 {
				v[len(v)+k] = x
			} else {
				v[k] = x
			}

			if x >= n && y >= m {
				return trace, nil
			}
		}
	}

	return trace, nil
}
