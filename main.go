package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	version string
	name    string
)

const (
	// connector
	CONNECTOR_CROSS       = "├── "
	CONNECTOR_RIGHT_ANGLE = "└── "
	CONNECTOR_LINE        = "│   "
	CONNECTOR_BLANK       = "    "

	// print color
	PRINT_COLOR_RED    = "\x1b[31m%s\x1b[0m"
	PRINT_COLOR_GREEN  = "\x1b[32m%s\x1b[0m"
	PRINT_COLOR_YELLOW = "\x1b[33m%s\x1b[0m"
	PRINT_COLOR_BLUE   = "\x1b[34m%s\x1b[0m"
)

type Walker struct {
	dirNum     int
	fileNum    int
	isEndDir   []bool
	colored    bool
	level      uint
	permission bool
	includeDot bool
}

type Row struct {
	file         os.FileInfo
	level        uint
	onRightAngle bool
	isBlank      []bool
	colored      bool
	permission   bool
}

func (row *Row) Name() string {
	name := row.file.Name()

	if row.colored {
		if row.file.IsDir() {
			name = fmt.Sprintf(PRINT_COLOR_BLUE, name)
		} else {
			if row.isExec() {
				name = fmt.Sprintf(PRINT_COLOR_GREEN, name)
			}
		}
	}

	if row.permission {
		name = fmt.Sprintf("[%s]  %s", row.Mode(), name)
	}

	return name
}

func (row *Row) Str() string {
	var str string
	for i := 0; i < int(row.level-1); i++ {
		if row.isBlank[i] {
			str += CONNECTOR_BLANK
		} else {
			str += CONNECTOR_LINE
		}
	}

	if row.onRightAngle {
		str += CONNECTOR_RIGHT_ANGLE + row.Name()
	} else {
		str += CONNECTOR_CROSS + row.Name()
	}

	// debug
	// str += fmt.Sprintf("(level=%d, isblank=%v)", row.level, row.isBlank)

	return str
}

func (row *Row) Mode() string {
	var m uint32
	m = uint32(row.file.Mode())
	const str = "dalTLDpSugct?"
	var modeStr [10]string

	for i, c := range str {
		if m&(1<<uint(32-1-i)) != 0 {
			if row.colored {
				modeStr[0] = fmt.Sprintf(PRINT_COLOR_BLUE, string(c))
			} else {
				modeStr[0] = string(c)
			}
		}
	}

	if modeStr[0] == "" {
		modeStr[0] = "."
	}

	w := 1
	const rwx = "rwxrwxrwx"
	for i, c := range rwx {
		if m&(1<<uint(9-1-i)) != 0 {
			if row.colored {
				switch s := string(c); s {
				case "r":
					modeStr[w] = fmt.Sprintf(PRINT_COLOR_YELLOW, string(c))
				case "w":
					modeStr[w] = fmt.Sprintf(PRINT_COLOR_RED, string(c))
				case "x":
					modeStr[w] = fmt.Sprintf(PRINT_COLOR_GREEN, string(c))
				}
			} else {
				modeStr[w] = string(c)
			}
		} else {
			modeStr[w] = "-"
		}
		w++
	}

	return strings.Join(modeStr[:], "")
}

func (row *Row) isExec() bool {
	var m uint32
	m = uint32(row.file.Mode())

	const rwx = "rwxrwxrwx"
	for i := 0; i < 9; i++ {
		if m&(1<<uint(9-1-i)) != 0 && i%3 == 2 {
			return true
		}
	}

	return false
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

func (w *Walker) Walk(dir string, level uint) error {
	if level > w.level {
		return nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for i, file := range files {
		if !w.includeDot && file.Name()[:1] == "." && file.Name() != "." {
			continue
		}

		if int(level)-len(w.isEndDir) == 1 {
			w.isEndDir = append(w.isEndDir, false)
		}

		if int(level) < len(w.isEndDir) {
			w.isEndDir = w.isEndDir[:level]
		}

		var onRightAngle bool
		if i == len(files)-1 {
			onRightAngle = true
			w.isEndDir[level-1] = true
		}

		row := Row{
			file:         file,
			level:        level,
			onRightAngle: onRightAngle,
			isBlank:      w.isEndDir,
			colored:      w.colored,
			permission:   w.permission,
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

func Tree(root string, colored bool, level uint, permission bool) error {
	w := Walker{
		dirNum:     0,
		fileNum:    0,
		isEndDir:   []bool{},
		colored:    colored,
		level:      level,
		permission: permission,
		includeDot: false,
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
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:    "level",
				Aliases: []string{"L"},
				Value:   math.MaxUint64,
				Usage:   "Descend only level directories deep.",
			},
			&cli.BoolFlag{
				Name:    "disable-color",
				Aliases: []string{"d"},
				Usage:   "Disable color.",
			},
			&cli.BoolFlag{
				Name:    "permission",
				Aliases: []string{"p"},
				Usage:   "Print permission.",
			},
		},
		Action: func(c *cli.Context) error {
			root := c.Args().Get(0)

			colored := true
			if c.Bool("disable-color") {
				colored = false
			}

			level := c.Uint("level")

			permission := c.Bool("permission")

			err := Tree(root, colored, level, permission)
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
