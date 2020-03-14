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

func (row *Row) Str() string {
	c := CONNECTOR_CROSS
	if row.IsEnd {
		c = CONNECTOR_RIGHT_ANGLE
	}

	c = strings.Repeat(CONNECTOR_LINE, row.Level-1) + c

	str := c + row.Name

	return str
}

func Walk(root string, level int) error {

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for i, file := range files {
		isEnd := false
		if i == len(files)-1 {
			isEnd = true
		}

		row := Row{
			Name:  file.Name(),
			Level: level,
			IsEnd: isEnd,
		}

		fmt.Println(row.Str())

		if file.IsDir() {
			path := filepath.Join(root, file.Name())
			err := Walk(path, level+1)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func Tree(root string) error {
	fmt.Println(root)
	err := Walk(root, 1)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	root := "sample"
	err := Tree(root)
	if err != nil {
		log.Fatal(err)
	}
}
