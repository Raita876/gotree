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
	Root      string
	DirNum    int
	FileNum   int
	IsLastDir bool
}

type Row struct {
	Name      string
	Level     int
	IsEnd     bool
	WithBlank bool
}

func (row *Row) Str() string {
	c := CONNECTOR_CROSS
	if row.IsEnd {
		c = CONNECTOR_RIGHT_ANGLE
	}

	if row.WithBlank {
		c = strings.Repeat(CONNECTOR_BLANK, row.Level-1) + c
	} else {
		c = strings.Repeat(CONNECTOR_LINE, row.Level-1) + c
	}

	str := c + row.Name

	return str
}

func (w *Walker) Walk(dir string, level int) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for i, file := range files {
		isEnd := false
		if i == len(files)-1 {
			isEnd = true
			if dir == w.Root {
				w.IsLastDir = true
			}
		}

		row := Row{
			Name:      file.Name(),
			Level:     level,
			IsEnd:     isEnd,
			WithBlank: w.IsLastDir,
		}

		fmt.Println(row.Str())

		if file.IsDir() {
			path := filepath.Join(dir, file.Name())
			err := w.Walk(path, level+1)
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
		Root:      root,
		DirNum:    0,
		FileNum:   0,
		IsLastDir: false,
	}

	fmt.Println(root)

	err := w.Walk(root, 1)
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
