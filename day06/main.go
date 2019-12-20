package main

import (
	"bufio"
	"container/list"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	//log, _ = zap.NewDevelopment()
	log, _   = zap.NewProduction()
	logSugar = log.Sugar()
)

type Node struct {
	name  string
	depth int
}

func main() {
	lines, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("failed", zap.Error(err))
	}

	graph := make(map[string][]string)

	for _, relation := range lines {
		if len(relation) == 0 {
			log.Debug("Line is empty, skipping parsing")
			continue
		}
		idx := strings.Index(relation, ")")
		if idx == -1 {
			logSugar.Debugf("Skipping line '%'", relation)
			continue
		}
		src := relation[:idx]
		dst := relation[idx+1:]

		graph[src] = append(graph[src], dst)
		graph[dst] = append(graph[dst], src)
		logSugar.Debugf("Found: %s is orbited by %s", src, dst)
	}

	visited := make(map[string]struct{})
	totalOrbits := 0
	q := list.New()
	q.PushBack(Node{name: "COM", depth: 0})
	for q.Len() > 0 {
		node := q.Front().Value.(Node)
		q.Remove(q.Front())
		visited[node.name] = struct{}{}
		//log.Debug("Node has orbits", zap.String("node", node.name), zap.Strings("children", graph[node.name]))
		for _, child := range graph[node.name] {
			if _, ok := visited[child]; !ok {
				totalOrbits += node.depth + 1
				q.PushBack(Node{name: child, depth: node.depth + 1})
			}
		}
	}

	p1(graph)
	p2(graph, "YOU", "SAN")
}

func p1(graph map[string][]string) {
	visited := make(map[string]struct{})
	totalOrbits := 0
	q := list.New()
	q.PushBack(Node{name: "COM", depth: 0})
	for q.Len() > 0 {
		node := q.Front().Value.(Node)
		q.Remove(q.Front())
		visited[node.name] = struct{}{}
		//log.Debug("Node has orbits", zap.String("node", node.name), zap.Strings("children", graph[node.name]))
		for _, child := range graph[node.name] {
			if _, ok := visited[child]; !ok {
				totalOrbits += node.depth + 1
				q.PushBack(Node{name: child, depth: node.depth + 1})
			}
		}
	}
	fmt.Printf("Part1: total orbits (direct and indirect) %d\n", totalOrbits)

}

func p2(graph map[string][]string, src string, target string) {
	visited := make(map[string]struct{})
	q := list.New()
	q.PushBack(Node{name: src, depth: 0})
	for q.Len() > 0 {
		node := q.Front().Value.(Node)
		q.Remove(q.Front())
		visited[node.name] = struct{}{}
		//log.Debug("Node has orbits", zap.String("node", node.name), zap.Strings("children", graph[node.name]))
		for _, child := range graph[node.name] {
			if child == target {
				fmt.Printf("Part2: path length %d\n", node.depth-1)
				return
			}
			if _, ok := visited[child]; !ok {
				q.PushBack(Node{name: child, depth: node.depth + 1})
			}
		}
	}

}

// read and trim each line from the given filename
func readInputFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Warn("failed to close", zap.Error(err))
		}
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}
