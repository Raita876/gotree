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
	Name      string
	Level     int
	DoneLevel int
	IsEnd     bool
}

func (row *Row) Str() string {
	c := CONNECTOR_CROSS
	if row.IsEnd {
		c = CONNECTOR_RIGHT_ANGLE
	}

	c = strings.Repeat(CONNECTOR_BLANK, row.DoneLevel) + strings.Repeat(CONNECTOR_LINE, row.Level-1-row.DoneLevel) + c

	str := c + row.Name

	// debug
	str += fmt.Sprintf("(level=%d, done=%d)", row.Level, row.DoneLevel)

	return str
}

func (w *Walker) Walk(dir string, level int, doneLevel int) error {

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
			Name:      file.Name(),
			Level:     level,
			DoneLevel: doneLevel,
			IsEnd:     isEnd,
		}

		fmt.Println(row.Str(), dir)

		if i == len(files)-1 {
			doneLevel++
		}

		if file.IsDir() {
			path := filepath.Join(dir, file.Name())
			err := w.Walk(path, level+1, doneLevel)
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

	err := w.Walk(root, 1, 0)
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
