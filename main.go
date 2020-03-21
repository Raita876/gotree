package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	version string
	name    string
)

const (
	CONNECTOR_CROSS       = "├── "
	CONNECTOR_RIGHT_ANGLE = "└── "
	CONNECTOR_LINE        = "│   "
	CONNECTOR_BLANK       = "    "
)

type Walker struct {
	dirNum   int
	fileNum  int
	isEndDir []bool
	colored  bool
}

type Row struct {
	name         string
	level        int
	onRightAngle bool
	isBlank      []bool
	isDir        bool
	colored      bool
}

func (row *Row) Name(colored bool) string {
	name := row.name

	if colored {
		if row.isDir {
			name = fmt.Sprintf("\x1b[34m%s\x1b[0m", row.name)
		} else {
			name = fmt.Sprintf("\x1b[31m%s\x1b[0m", row.name)
		}

	}

	return name
}

func (row *Row) Str() string {
	var str string
	for i := 0; i < row.level-1; i++ {
		if row.isBlank[i] {
			str += CONNECTOR_BLANK
		} else {
			str += CONNECTOR_LINE
		}
	}

	if row.onRightAngle {
		str += CONNECTOR_RIGHT_ANGLE + row.Name(row.colored)
	} else {
		str += CONNECTOR_CROSS + row.Name(row.colored)
	}

	// debug
	// str += fmt.Sprintf("(level=%d, isblank=%v)", row.level, row.isBlank)

	return str
}

func (w *Walker) PrintRoot(root string) {
	fmt.Println(root)
}

func (w *Walker) PrintRow(row Row) {
	fmt.Println(row.Str())
}

func (w *Walker) PrintResult() {
	fmt.Printf("\n%d directories, %d files\n", w.dirNum, w.fileNum)
}

func (w *Walker) Walk(dir string, level int) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for i, file := range files {

		if level-len(w.isEndDir) == 1 {
			w.isEndDir = append(w.isEndDir, false)
		}

		if level < len(w.isEndDir) {
			w.isEndDir = w.isEndDir[:level]
		}

		var onRightAngle bool
		if i == len(files)-1 {
			onRightAngle = true
			w.isEndDir[level-1] = true
		}

		row := Row{
			name:         file.Name(),
			level:        level,
			onRightAngle: onRightAngle,
			isBlank:      w.isEndDir,
			isDir:        file.IsDir(),
			colored:      w.colored,
		}

		w.PrintRow(row)

		if file.IsDir() {
			path := filepath.Join(dir, file.Name())
			err := w.Walk(path, level+1)
			if err != nil {
				return err
			}

			w.dirNum++
		} else {
			w.fileNum++
		}

	}

	return nil
}

func Tree(root string) error {
	w := Walker{
		dirNum:   0,
		fileNum:  0,
		isEndDir: []bool{},
		colored:  false,
	}

	w.PrintRoot(root)

	err := w.Walk(root, 1)
	if err != nil {
		return err
	}

	w.PrintResult()

	return nil
}

func main() {
	app := &cli.App{
		Version: version,
		Name:    name,
		Usage:   "Golang tree command.",
		Action: func(c *cli.Context) error {
			root := c.Args().Get(0)
			err := Tree(root)
			if err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
