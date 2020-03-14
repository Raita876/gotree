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
	CONNECTOR_BLANK       = "    "
)

type Walker struct {
	DirNum  int
	FileNum int
}

type Row struct {
	Name  string
	Level int
	Done  int
	Blank int
	IsEnd bool
}

func (row *Row) Str() string {
	c := CONNECTOR_CROSS
	if row.IsEnd {
		c = CONNECTOR_RIGHT_ANGLE
	}

	c = strings.Repeat(CONNECTOR_BLANK, row.Done) + strings.Repeat(CONNECTOR_LINE, row.Level-1-row.Done-row.Blank) + strings.Repeat(CONNECTOR_BLANK, row.Blank) + c

	str := c + row.Name

	// debug
	str += fmt.Sprintf("(level=%d, done=%d, blank=%d)", row.Level, row.Done, row.Blank)

	return str
}

func (w *Walker) Walk(root string, dir string, level int, done int, blank int) error {

	files, err := ioutil.ReadDir(dir)
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
			Done:  done,
			Blank: blank,
			IsEnd: isEnd,
		}

		fmt.Println(row.Str())

		path := filepath.Join(dir, file.Name())

		if i == len(files)-1 {
			if dir == root {
				done++
				root = path
			} else {
				blank++
			}
		}

		if file.IsDir() {
			err := w.Walk(root, path, level+1, done, blank)
			if err != nil {
				return err
			}

			w.DirNum++
		} else {
			w.FileNum++
		}

	}

	return nil
}

func Tree(root string) error {
	w := Walker{
		DirNum:  0,
		FileNum: 0,
	}

	fmt.Println(root)

	err := w.Walk(root, root, 1, 0, 0)
	if err != nil {
		return err
	}

	fmt.Printf("\n%d directories, %d files\n", w.DirNum, w.FileNum)

	return nil
}

func main() {
	root := "sample"
	err := Tree(root)
	if err != nil {
		log.Fatal(err)
	}
}
