package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

const (
	CONNECTOR_CROSS       = "├── "
	CONNECTOR_RIGHT_ANGLE = "└── "
	CONNECTOR_LINE        = "│   "
	CONNECTOR_BLANK       = "    "
)

type Walker struct {
	DirNum   int
	FileNum  int
	IsEndDir []bool
}

type Row struct {
	Name         string
	Level        int
	OnRightAngle bool
	IsBlank      []bool
}

func (row *Row) Str() string {
	var str string
	for i := 0; i < row.Level-1; i++ {
		if row.IsBlank[i] {
			str += CONNECTOR_BLANK
		} else {
			str += CONNECTOR_LINE
		}
	}

	if row.OnRightAngle {
		str += CONNECTOR_RIGHT_ANGLE + row.Name
	} else {
		str += CONNECTOR_CROSS + row.Name
	}

	// debug
	// str += fmt.Sprintf("(level=%d, isblank=%v)", row.Level, row.IsBlank)

	return str
}

func (w *Walker) Walk(dir string, level int) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for i, file := range files {

		if level-len(w.IsEndDir) == 1 {
			w.IsEndDir = append(w.IsEndDir, false)
		}

		if level < len(w.IsEndDir) {
			w.IsEndDir = w.IsEndDir[:level]
		}

		var onRightAngle bool
		if i == len(files)-1 {
			onRightAngle = true
			w.IsEndDir[level-1] = true
		}

		row := Row{
			Name:         file.Name(),
			Level:        level,
			OnRightAngle: onRightAngle,
			IsBlank:      w.IsEndDir,
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
		DirNum:  0,
		FileNum: 0,
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
