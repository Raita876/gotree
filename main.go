package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

const (
	CONNECTOR_CROSS       = "├── "
	CONNECTOR_RIGHT_ANGLE = "└── "
	CONNECTOR_LINE        = "│   "
)

type Row struct {
	Name  string
	Level int
	IsEnd bool
}

func Walk(root string, level int) ([]Row, error) {
	rows := []Row{}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return rows, err
	}

	for i, file := range files {
		isEnd := false
		if i == len(files)-1 {
			isEnd = true
		}

		rows = append(rows, Row{
			Name:  file.Name(),
			Level: level,
			IsEnd: isEnd,
		})

		if file.IsDir() {
			path := filepath.Join(root, file.Name())
			r, err := Walk(path, level+1)
			if err != nil {
				return rows, err
			}
			rows = append(rows, r...)
		}

	}

	return rows, nil
}

func Tree(rows []Row, root string) {
	fmt.Println(root)

	for _, row := range rows {
		fmt.Println(row.Str())
	}
}

func (row *Row) Str() string {
	c := CONNECTOR_CROSS
	if row.IsEnd {
		c = CONNECTOR_RIGHT_ANGLE
	}

	c = strings.Repeat(CONNECTOR_LINE, row.Level-1) + c

	str := c + row.Name

	return str
}

func main() {
	root := "sample"
	rows, err := Walk(root, 1)
	if err != nil {
		log.Fatal(err)
	}

	Tree(rows, root)
}
