package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	children []*node
	parent   *node
	name     string
}

func (n *node) add_child(c *node) {
	n.children = append(n.children, c)
	c.parent = n
}

func (n *node) String() string {
	if len(n.children) == 0 {
		return n.name
	}
	subtrees := make([]string, len(n.children))
	for i, c := range n.children {
		subtrees[i] = fmt.Sprint(c)
	}
	return "(" + strings.Join(subtrees, ",") + ")" + n.name
}

func parse(filename string) (m map[string]*node) {
	f, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Could not open file %s", filename))
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	m = make(map[string]*node)
	for s.Scan() {
		parseLine(s.Text(), m)
	}
	if err := s.Err(); err != nil {
		fmt.Println(os.Stderr, "reading file:", err)
	}
	return
}

func parseLine(line string, m map[string]*node) {
	elements := strings.Split(line, "\t")
	for i := 0; i < (len(elements) - 1); i++ {
		parent := elements[i]
		child := elements[i+1]
		if _, ok := m[parent]; !ok {
			m[parent] = &node{name: parent}
		}
		if _, ok := m[child]; !ok {
			m[child] = &node{name: child}
		}
		m[parent].add_child(m[child])
	}
}

func main() {
	m := parse("example.txt")
	// l := "NGo1	ASk1	VBo1	ASk2"
	// m := make(map[string]*node)
	// parseLine(l, m)
	fmt.Println(m)
	n := m["NGo1"]
	fmt.Println(n.children)
}
