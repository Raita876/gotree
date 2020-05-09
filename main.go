package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

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
	uid        bool
	gid        bool
	size       bool
	includeDot bool
}

type Row struct {
	fileInfo     os.FileInfo
	level        uint
	onRightAngle bool
	isBlank      []bool
	colored      bool
	permission   bool
	uid          bool
	gid          bool
	size         bool
}

func (row *Row) Status() string {
	status := ""

	if row.permission {
		status += row.Mode() + " "
	}

	if row.uid {
		status += row.User() + " "
	}

	if row.gid {
		status += row.Group() + " "
	}

	if row.size {
		status += row.Size() + " "
	}

	if status != "" {
		return fmt.Sprintf("[%s]  ", strings.TrimSpace(status))
	}

	return status
}

func (row *Row) Size() string {
	if row.fileInfo.IsDir() {
		return "-"
	}

	size := row.fileInfo.Size()
	fs := FormatSize(size)

	if row.colored {
		fs = fmt.Sprintf(PRINT_COLOR_GREEN, fs)
	}

	return fs
}

func FormatSize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d", size)
	}

	prefix := "kMGTP"
	for i := 0; i < len(prefix); i++ {
		s := int(size) / 1000 * (i + 1)
		if s < 1000 {
			return fmt.Sprintf("%d%s", s, string(prefix[i]))
		}
	}

	return "?????"
}

func (row *Row) User() string {
	var userName string
	var uid string

	if stat, ok := row.fileInfo.Sys().(*syscall.Stat_t); ok {
		uid = fmt.Sprintf("%d", stat.Uid)
	} else {
		uid = fmt.Sprintf("%d", os.Getuid())
	}

	u, err := user.LookupId(uid)
	if err != nil {
		userName = uid
	} else {
		userName = u.Username
	}

	if row.colored {
		userName = fmt.Sprintf(PRINT_COLOR_YELLOW, userName)
	}

	return userName
}

func (row *Row) Group() string {
	var group string
	var gid string

	if stat, ok := row.fileInfo.Sys().(*syscall.Stat_t); ok {
		gid = fmt.Sprintf("%d", stat.Gid)
	} else {
		gid = fmt.Sprintf("%d", os.Getgid())
	}

	g, err := user.LookupGroupId(gid)
	if err != nil {
		group = gid
	} else {
		group = g.Name
	}

	if row.colored {
		group = fmt.Sprintf(PRINT_COLOR_YELLOW, group)
	}

	return group
}

func (row *Row) Name() string {
	name := row.fileInfo.Name()

	if row.colored {
		if row.fileInfo.IsDir() {
			name = fmt.Sprintf(PRINT_COLOR_BLUE, name)
		} else {
			if row.isExec() {
				name = fmt.Sprintf(PRINT_COLOR_GREEN, name)
			}
		}
	}

	return name
}

func (row *Row) File() string {
	return fmt.Sprintf("%s%s", row.Status(), row.Name())
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
		str += CONNECTOR_RIGHT_ANGLE + row.File()
	} else {
		str += CONNECTOR_CROSS + row.File()
	}

	// debug
	// str += fmt.Sprintf("(level=%d, isblank=%v)", row.level, row.isBlank)

	return str
}

func (row *Row) Mode() string {
	var m uint32
	m = uint32(row.fileInfo.Mode())
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
	m = uint32(row.fileInfo.Mode())

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
			fileInfo:     file,
			level:        level,
			onRightAngle: onRightAngle,
			isBlank:      w.isEndDir,
			colored:      w.colored,
			permission:   w.permission,
			uid:          w.uid,
			gid:          w.gid,
			size:         w.size,
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

func Tree(root string, colored bool, level uint, permission bool, uid bool, gid bool, size bool, includeDot bool) error {
	w := Walker{
		dirNum:     0,
		fileNum:    0,
		isEndDir:   []bool{},
		colored:    colored,
		level:      level,
		permission: permission,
		uid:        uid,
		gid:        gid,
		size:       size,
		includeDot: includeDot,
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
			&cli.BoolFlag{
				Name:    "uid",
				Aliases: []string{"u"},
				Usage:   "Print file owner or UID number.",
			},
			&cli.BoolFlag{
				Name:    "gid",
				Aliases: []string{"g"},
				Usage:   "Print file group or GID number.",
			},
			&cli.BoolFlag{
				Name:    "size",
				Aliases: []string{"s"},
				Usage:   "Print the size.",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "All files are listed.",
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

			uid := c.Bool("uid")
			gid := c.Bool("gid")

			size := c.Bool("size")

			includeDot := c.Bool("all")

			err := Tree(root, colored, level, permission, uid, gid, size, includeDot)
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
