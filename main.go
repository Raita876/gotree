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
	PRINT_COLOR_BLACK         = "\x1b[30m%s\x1b[0m"
	PRINT_COLOR_RED           = "\x1b[31m%s\x1b[0m"
	PRINT_COLOR_GREEN         = "\x1b[32m%s\x1b[0m"
	PRINT_COLOR_YELLOW        = "\x1b[33m%s\x1b[0m"
	PRINT_COLOR_BLUE          = "\x1b[34m%s\x1b[0m"
	PRINT_COLOR_PURPLE        = "\x1b[35m%s\x1b[0m"
	PRINT_COLOR_CYAN          = "\x1b[36m%s\x1b[0m"
	PRINT_COLOR_LIGHT_GRAY    = "\x1b[37m%s\x1b[0m"
	PRINT_COLOR_DARK_GRAY     = "\x1b[90m%s\x1b[0m"
	PRINT_COLOR_LIGHT_RED     = "\x1b[91m%s\x1b[0m"
	PRINT_COLOR_LIGHT_GREEN   = "\x1b[92m%s\x1b[0m"
	PRINT_COLOR_LIGHT_YELLOW  = "\x1b[93m%s\x1b[0m"
	PRINT_COLOR_LIGHT_BLUE    = "\x1b[94m%s\x1b[0m"
	PRINT_COLOR_LIGHT_MAGENTA = "\x1b[95m%s\x1b[0m"
	PRINT_COLOR_LIGHT_CYAN    = "\x1b[96m%s\x1b[0m"
	PRINT_COLOR_LIGHT_WHITE   = "\x1b[97m%s\x1b[0m"
)

func ColorRed(s string) string {
	return fmt.Sprintf(PRINT_COLOR_RED, s)
}

func ColorGreen(s string) string {
	return fmt.Sprintf(PRINT_COLOR_GREEN, s)
}

func ColorYellow(s string) string {
	return fmt.Sprintf(PRINT_COLOR_YELLOW, s)
}

func ColorBlue(s string) string {
	return fmt.Sprintf(PRINT_COLOR_BLUE, s)
}

func ColorPurple(s string) string {
	return fmt.Sprintf(PRINT_COLOR_PURPLE, s)
}

func ColorCyan(s string) string {
	return fmt.Sprintf(PRINT_COLOR_CYAN, s)
}

func ColorLightGray(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_GRAY, s)
}

func ColorDarkGray(s string) string {
	return fmt.Sprintf(PRINT_COLOR_DARK_GRAY, s)
}

func ColorLightRed(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_RED, s)
}

func ColorLightGreen(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_GREEN, s)
}

func ColorLightYellow(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_YELLOW, s)
}

func ColorLightBlue(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_BLUE, s)
}

func ColorLightMagenta(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_MAGENTA, s)
}

func ColorLightCyan(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_CYAN, s)
}

func ColorLightWhite(s string) string {
	return fmt.Sprintf(PRINT_COLOR_LIGHT_WHITE, s)
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

func contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}

	return false
}

func ext(fileName string) string {
	ext := filepath.Ext(fileName)

	if ext == "" {
		return ext
	} else {
		return ext[1:]
	}
}

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
	datetime   bool
}

type Option interface {
	apply(*Walker)
}

type coloredOption bool

func (c coloredOption) apply(w *Walker) {
	w.colored = bool(c)
}

type levelOption uint

func (l levelOption) apply(w *Walker) {
	w.level = uint(l)
}

type permissionOption bool

func (p permissionOption) apply(w *Walker) {
	w.permission = bool(p)
}

type uidOption bool

func (u uidOption) apply(w *Walker) {
	w.uid = bool(u)
}

type gidOption bool

func (g gidOption) apply(w *Walker) {
	w.gid = bool(g)
}

type sizeOption bool

func (s sizeOption) apply(w *Walker) {
	w.size = bool(s)
}

type includeDotOption bool

func (i includeDotOption) apply(w *Walker) {
	w.includeDot = bool(i)
}

type datetimeOption bool

func (dt datetimeOption) apply(w *Walker) {
	w.datetime = bool(dt)
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
	datetime     bool
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

	if row.datetime {
		status += row.Datetime() + " "
	}

	if status != "" {
		return fmt.Sprintf("[%s]  ", strings.TrimSpace(status))
	}

	return status
}

func (row *Row) Datetime() string {
	mt := row.fileInfo.ModTime().Format("2006-01-02 15:04")

	if row.colored {
		mt = ColorBlue(mt)
	}

	return mt
}

func (row *Row) Size() string {
	if row.fileInfo.IsDir() {
		return "-"
	}

	size := row.fileInfo.Size()
	fs := FormatSize(size)

	if row.colored {
		fs = ColorGreen(fs)
	}

	return fs
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
		userName = ColorYellow(userName)
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
		group = ColorYellow(group)
	}

	return group
}

func (row *Row) Name() string {
	name := row.fileInfo.Name()

	if row.colored {
		if row.isDir() {
			return ColorBlue(name)
		}

		if row.isExec() {
			return ColorLightGreen(name)
		}

		if row.isImmediate() {
			return ColorLightYellow(name)
		}

		if row.isImage() {
			return ColorLightMagenta(name)
		}

		if row.isVideo() {
			return ColorPurple(name)
		}

		if row.isMusic() {
			return ColorPurple(name)
		}

		if row.isCrypto() {
			return ColorLightCyan(name)
		}

		if row.isDocument() {
			return ColorGreen(name)
		}

		if row.isCompressed() {
			return ColorRed(name)
		}

		if row.isTemp() {
			return ColorDarkGray(name)
		}

		if row.isCompiled() {
			return ColorYellow(name)
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
				modeStr[0] = ColorBlue(string(c))
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
					modeStr[w] = ColorYellow(string(c))
				case "w":
					modeStr[w] = ColorRed(string(c))
				case "x":
					modeStr[w] = ColorGreen(string(c))
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

func (row *Row) isDir() bool {
	return row.fileInfo.IsDir()
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

func (row *Row) isImmediate() bool {
	name := row.fileInfo.Name()

	nameWithoutExt := name[:len(name)-len(filepath.Ext(name))]
	if strings.ToLower(nameWithoutExt) == "readme" {
		return true
	}

	immediateFiles := []string{
		"Makefile", "Cargo.toml", "SConstruct", "CMakeLists.txt",
		"build.gradle", "pom.xml", "Rakefile", "package.json", "Gruntfile.js",
		"Gruntfile.coffee", "BUILD", "BUILD.bazel", "WORKSPACE", "build.xml",
		"webpack.config.js", "meson.build",
	}

	if contains(immediateFiles, name) {
		return true
	}

	return false
}

func (row *Row) isImage() bool {
	imageExts := []string{
		"png", "jpeg", "jpg", "gif", "bmp", "tiff", "tif",
		"ppm", "pgm", "pbm", "pnm", "webp", "raw", "arw",
		"svg", "stl", "eps", "dvi", "ps", "cbr", "jpf",
		"cbz", "xpm", "ico", "cr2", "orf", "nef",
	}

	ext := ext(row.fileInfo.Name())

	if contains(imageExts, ext) {
		return true
	}

	return false
}

func (row *Row) isVideo() bool {
	videoExts := []string{
		"avi", "flv", "m2v", "m4v", "mkv", "mov", "mp4", "mpeg",
		"mpg", "ogm", "ogv", "vob", "wmv", "webm", "m2ts",
	}

	ext := ext(row.fileInfo.Name())

	if contains(videoExts, ext) {
		return true
	}

	return false
}

func (row *Row) isMusic() bool {
	musicExts := []string{
		"aac", "m4a", "mp3", "ogg", "wma", "mka", "opus",
		"alac", "ape", "flac", "wav",
	}

	ext := ext(row.fileInfo.Name())

	if contains(musicExts, ext) {
		return true
	}

	return false
}

func (row *Row) isCrypto() bool {
	cryptoExts := []string{
		"asc", "enc", "gpg", "pgp", "sig", "signature", "pfx", "p12",
	}

	ext := ext(row.fileInfo.Name())

	if contains(cryptoExts, ext) {
		return true
	}

	return false
}

func (row *Row) isDocument() bool {
	documentExts := []string{
		"djvu", "doc", "docx", "dvi", "eml", "eps", "fotd",
		"odp", "odt", "pdf", "ppt", "pptx", "rtf",
		"xls", "xlsx",
	}

	ext := ext(row.fileInfo.Name())

	if contains(documentExts, ext) {
		return true
	}

	return false
}

func (row *Row) isCompressed() bool {
	compressedExts := []string{
		"zip", "tar", "Z", "z", "gz", "bz2", "a", "ar", "7z",
		"iso", "dmg", "tc", "rar", "par", "tgz", "xz", "txz",
		"lz", "tlz", "lzma", "deb", "rpm", "zst",
	}

	ext := ext(row.fileInfo.Name())

	if contains(compressedExts, ext) {
		return true
	}

	return false
}

func (row *Row) isTemp() bool {
	name := row.fileInfo.Name()

	// XXXX~ or #XXXX#
	if name[len(name)-1:] == "~" || (name[:1] == "#" && name[len(name)-1:] == "#") {
		return true
	}

	tempExts := []string{"tmp", "swp", "swo", "swn", "bak", "bk"}

	ext := ext(name)

	if contains(tempExts, ext) {
		return true
	}

	return false
}

func (row *Row) isCompiled() bool {
	compiledExts := []string{"class", "elc", "hi", "o", "pyc", "zwc"}

	ext := ext(row.fileInfo.Name())

	if contains(compiledExts, ext) {
		return true
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
			datetime:     w.datetime,
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

func Tree(root string, opts ...Option) error {
	w := &Walker{
		dirNum:     0,
		fileNum:    0,
		isEndDir:   []bool{},
		colored:    true,
		level:      math.MaxUint64,
		permission: false,
		uid:        false,
		gid:        false,
		size:       false,
		includeDot: false,
		datetime:   false,
	}

	for _, o := range opts {
		o.apply(w)
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
				Name:    "datetime",
				Aliases: []string{"D"},
				Usage:   "Print file datetime.",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "All files are listed.",
			},
		},
		Action: func(c *cli.Context) error {
			root := c.Args().Get(0) // TODO: 引数の数をチェックする。

			level := levelOption(c.Uint("level"))
			colored := coloredOption(!c.Bool("disable-color"))
			permission := permissionOption(c.Bool("permission"))
			uid := uidOption(c.Bool("uid"))
			gid := gidOption(c.Bool("gid"))
			size := sizeOption(c.Bool("size"))
			includeDot := includeDotOption(c.Bool("all"))
			datetime := datetimeOption(c.Bool("datetime"))

			err := Tree(root, colored, level, permission, uid, gid, size, includeDot, datetime)
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
