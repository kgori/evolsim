package main

import (
	"bufio"
	"flag"
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
	return n.newick(true)
}

func (n *node) newick(inner bool) string {
	if len(n.children) == 0 {
		return n.name
	}
	subtrees := make([]string, len(n.children))
	for i, c := range n.children {
		subtrees[i] = c.newick(inner)
	}
	s := "(" + strings.Join(subtrees, ",") + ")"
	if inner {
		s += n.name
	}
	return s
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

func findRoot(n *node) *node {
	for n.parent != nil {
		n = n.parent
	}
	return n
}

func unroot(n *node) *node {
	c := n.children[0]
	for _, s := range n.children[1:len(n.children)] {
		c.add_child(s)
	}
	return c
}

var filename = flag.String("filename", "example.txt", "Tab-delimited lists of node names, one line per sheet")
var inner = flag.Bool("inner", true, "Print inner node labels, TRUE or false")
var unroot_tree = flag.Bool("unroot", false, "Unroot the tree, true or FALSE")

func main() {
	flag.Parse()
	m := parse(*filename)
	var root *node
	for _, nd := range m {
		root = findRoot(nd)
		break
	}
	if *unroot_tree {
		root = unroot(root)
	}
	fmt.Println(root.newick(*inner) + ";")
}
