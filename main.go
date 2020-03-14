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

type Walker struct {
	DirNum  int
	FileNum int
}

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

func (w *Walker) Walk(root string, level int) error {

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
			err := w.Walk(path, level+1)
			if err != nil {
				return err
			}

			w.DirNum += 1
		} else {
			w.FileNum += 1
		}

	}

	return nil
}

func Tree(root string) error {
	fmt.Println(root)

	w := Walker{
		DirNum:  0,
		FileNum: 0,
	}

	err := w.Walk(root, 1)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%d directories, %d files\n", w.DirNum, w.FileNum)

	return nil
}

func main() {
	root := "sample"
	err := Tree(root)
	if err != nil {
		log.Fatal(err)
	}
}
