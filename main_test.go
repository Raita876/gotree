package main

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const TMP_DIR = "tmp"

func setup() error {
	files := []string{
		TMP_DIR + "/.aaa",
		TMP_DIR + "/.bbb/.ccc",
		TMP_DIR + "/01/README.md",
		TMP_DIR + "/01/compiled.o",
		TMP_DIR + "/01/compressed.zip",
		TMP_DIR + "/01/crypto.asc",
		TMP_DIR + "/01/document.xlsx",
		TMP_DIR + "/01/exec",
		TMP_DIR + "/01/image.png",
		TMP_DIR + "/01/music.mp3",
		TMP_DIR + "/01/tmp.bk",
		TMP_DIR + "/01/video.mp4",
		TMP_DIR + "/01/wav.wav",
		TMP_DIR + "/foo/bar/baz",
		TMP_DIR + "/foo/qux",
		TMP_DIR + "/foo/quux",
		TMP_DIR + "/corge",
		TMP_DIR + "/grault/garply/waldo/wibble",
		TMP_DIR + "/grault/garply/waldo/wobble",
		TMP_DIR + "/grault/garply/fred",
		TMP_DIR + "/grault/plugh",
		TMP_DIR + "/xyzzy/thud/wubble",
		TMP_DIR + "/xyzzy/thud/flob",
	}

	for _, f := range files {
		dir := filepath.Dir(f)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		_, err := os.Create(f)
		if err != nil {
			return err
		}

	}

	// TODO: 個別にパーミッションを付与しているので、struct で管理できるようにする。
	os.Chmod(TMP_DIR+"/01/exec", 0777)

	return nil
}

func reset() error {
	if fi, err := os.Stat(TMP_DIR); err == nil && fi.IsDir() {
		err := os.RemoveAll(TMP_DIR)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	setup()

	result := m.Run()

	reset()

	os.Exit(result)

}
func TestTree(t *testing.T) {
	testCaseWithDate, err := testCaseWithDate()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		want       string
		colored    coloredOption
		level      levelOption
		permission permissionOption
		uid        uidOption
		gid        gidOption
		size       sizeOption
		includeDot includeDotOption
		datetime   datetimeOption
	}{
		{
			name: "gotree <directory>",
			want: `tmp
[90m├── [0m[94m01[0m/
[90m│   [0m[90m├── [0m[4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[33mcompiled.o[0m
[90m│   [0m[90m├── [0m[31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[92mexec[0m*
[90m│   [0m[90m├── [0m[95mimage.png[0m
[90m│   [0m[90m├── [0m[35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[90mtmp.bk[0m
[90m│   [0m[90m├── [0m[35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[35mwav.wav[0m
[90m├── [0mcorge
[90m├── [0m[94mfoo[0m/
[90m│   [0m[90m├── [0m[94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0mbaz
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0m[94mgrault[0m/
[90m│   [0m[90m├── [0m[94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0mfred
[90m│   [0m[90m│   [0m[90m└── [0m[94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0mwibble
[90m│   [0m[90m│   [0m    [90m└── [0mwobble
[90m│   [0m[90m└── [0mplugh
[90m└── [0m[94mxyzzy[0m/
    [90m└── [0m[94mthud[0m/
        [90m├── [0mflob
        [90m└── [0mwubble

8 directories, 21 files`,
			colored:    true,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --disable-color <directory>",
			want: `tmp
[90m├── [0m01
[90m│   [0m[90m├── [0mREADME.md
[90m│   [0m[90m├── [0mcompiled.o
[90m│   [0m[90m├── [0mcompressed.zip
[90m│   [0m[90m├── [0mcrypto.asc
[90m│   [0m[90m├── [0mdocument.xlsx
[90m│   [0m[90m├── [0mexec
[90m│   [0m[90m├── [0mimage.png
[90m│   [0m[90m├── [0mmusic.mp3
[90m│   [0m[90m├── [0mtmp.bk
[90m│   [0m[90m├── [0mvideo.mp4
[90m│   [0m[90m└── [0mwav.wav
[90m├── [0mcorge
[90m├── [0mfoo
[90m│   [0m[90m├── [0mbar
[90m│   [0m[90m│   [0m[90m└── [0mbaz
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0mgrault
[90m│   [0m[90m├── [0mgarply
[90m│   [0m[90m│   [0m[90m├── [0mfred
[90m│   [0m[90m│   [0m[90m└── [0mwaldo
[90m│   [0m[90m│   [0m    [90m├── [0mwibble
[90m│   [0m[90m│   [0m    [90m└── [0mwobble
[90m│   [0m[90m└── [0mplugh
[90m└── [0mxyzzy
    [90m└── [0mthud
        [90m├── [0mflob
        [90m└── [0mwubble

8 directories, 21 files`,
			colored:    false,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree -L 2 <directory>",
			want: `tmp
[90m├── [0m[94m01[0m/
[90m│   [0m[90m├── [0m[4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[33mcompiled.o[0m
[90m│   [0m[90m├── [0m[31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[92mexec[0m*
[90m│   [0m[90m├── [0m[95mimage.png[0m
[90m│   [0m[90m├── [0m[35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[90mtmp.bk[0m
[90m│   [0m[90m├── [0m[35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[35mwav.wav[0m
[90m├── [0mcorge
[90m├── [0m[94mfoo[0m/
[90m│   [0m[90m├── [0m[94mbar[0m/
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0m[94mgrault[0m/
[90m│   [0m[90m├── [0m[94mgarply[0m/
[90m│   [0m[90m└── [0mplugh
[90m└── [0m[94mxyzzy[0m/
    [90m└── [0m[94mthud[0m/

7 directories, 15 files`,
			colored:    true,
			level:      2,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --disable-color -L 2 <directory>",
			want: `tmp
[90m├── [0m01
[90m│   [0m[90m├── [0mREADME.md
[90m│   [0m[90m├── [0mcompiled.o
[90m│   [0m[90m├── [0mcompressed.zip
[90m│   [0m[90m├── [0mcrypto.asc
[90m│   [0m[90m├── [0mdocument.xlsx
[90m│   [0m[90m├── [0mexec
[90m│   [0m[90m├── [0mimage.png
[90m│   [0m[90m├── [0mmusic.mp3
[90m│   [0m[90m├── [0mtmp.bk
[90m│   [0m[90m├── [0mvideo.mp4
[90m│   [0m[90m└── [0mwav.wav
[90m├── [0mcorge
[90m├── [0mfoo
[90m│   [0m[90m├── [0mbar
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0mgrault
[90m│   [0m[90m├── [0mgarply
[90m│   [0m[90m└── [0mplugh
[90m└── [0mxyzzy
    [90m└── [0mthud

7 directories, 15 files`,
			colored:    false,
			level:      2,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --permission <directory>",
			want: `tmp
[90m├── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94m01[0m/
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [33mcompiled.o[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m[32mx[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m[31mw[0m[32mx[0m]  [92mexec[0m*
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [95mimage.png[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [90mtmp.bk[0m
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  [35mwav.wav[0m
[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  corge
[90m├── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mfoo[0m/
[90m│   [0m[90m├── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  baz
[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  quux
[90m│   [0m[90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  qux
[90m├── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mgrault[0m/
[90m│   [0m[90m├── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wobble
[90m│   [0m[90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  plugh
[90m└── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mxyzzy[0m/
    [90m└── [0m[[94md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [94mthud[0m/
        [90m├── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  flob
        [90m└── [0m[.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wubble

8 directories, 21 files`,
			colored:    true,
			level:      math.MaxInt64,
			permission: true,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --disable-color --permission <directory>",
			want: `tmp
[90m├── [0m[drwxr-xr-x]  01
[90m│   [0m[90m├── [0m[.rw-r--r--]  README.md
[90m│   [0m[90m├── [0m[.rw-r--r--]  compiled.o
[90m│   [0m[90m├── [0m[.rw-r--r--]  compressed.zip
[90m│   [0m[90m├── [0m[.rw-r--r--]  crypto.asc
[90m│   [0m[90m├── [0m[.rw-r--r--]  document.xlsx
[90m│   [0m[90m├── [0m[.rwxrwxrwx]  exec
[90m│   [0m[90m├── [0m[.rw-r--r--]  image.png
[90m│   [0m[90m├── [0m[.rw-r--r--]  music.mp3
[90m│   [0m[90m├── [0m[.rw-r--r--]  tmp.bk
[90m│   [0m[90m├── [0m[.rw-r--r--]  video.mp4
[90m│   [0m[90m└── [0m[.rw-r--r--]  wav.wav
[90m├── [0m[.rw-r--r--]  corge
[90m├── [0m[drwxr-xr-x]  foo
[90m│   [0m[90m├── [0m[drwxr-xr-x]  bar
[90m│   [0m[90m│   [0m[90m└── [0m[.rw-r--r--]  baz
[90m│   [0m[90m├── [0m[.rw-r--r--]  quux
[90m│   [0m[90m└── [0m[.rw-r--r--]  qux
[90m├── [0m[drwxr-xr-x]  grault
[90m│   [0m[90m├── [0m[drwxr-xr-x]  garply
[90m│   [0m[90m│   [0m[90m├── [0m[.rw-r--r--]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[drwxr-xr-x]  waldo
[90m│   [0m[90m│   [0m    [90m├── [0m[.rw-r--r--]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[.rw-r--r--]  wobble
[90m│   [0m[90m└── [0m[.rw-r--r--]  plugh
[90m└── [0m[drwxr-xr-x]  xyzzy
    [90m└── [0m[drwxr-xr-x]  thud
        [90m├── [0m[.rw-r--r--]  flob
        [90m└── [0m[.rw-r--r--]  wubble

8 directories, 21 files`,
			colored:    false,
			level:      math.MaxInt64,
			permission: true,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree -a <directory>",
			want: `tmp
[90m├── [0m.aaa
[90m├── [0m[94m.bbb[0m/
[90m│   [0m[90m└── [0m.ccc
[90m├── [0m[94m01[0m/
[90m│   [0m[90m├── [0m[4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[33mcompiled.o[0m
[90m│   [0m[90m├── [0m[31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[92mexec[0m*
[90m│   [0m[90m├── [0m[95mimage.png[0m
[90m│   [0m[90m├── [0m[35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[90mtmp.bk[0m
[90m│   [0m[90m├── [0m[35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[35mwav.wav[0m
[90m├── [0mcorge
[90m├── [0m[94mfoo[0m/
[90m│   [0m[90m├── [0m[94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0mbaz
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0m[94mgrault[0m/
[90m│   [0m[90m├── [0m[94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0mfred
[90m│   [0m[90m│   [0m[90m└── [0m[94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0mwibble
[90m│   [0m[90m│   [0m    [90m└── [0mwobble
[90m│   [0m[90m└── [0mplugh
[90m└── [0m[94mxyzzy[0m/
    [90m└── [0m[94mthud[0m/
        [90m├── [0mflob
        [90m└── [0mwubble

9 directories, 23 files`,
			colored:    true,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: true,
			datetime:   false,
		},
		{
			name: "gotree --disable-color -a <directory>",
			want: `tmp
[90m├── [0m.aaa
[90m├── [0m.bbb
[90m│   [0m[90m└── [0m.ccc
[90m├── [0m01
[90m│   [0m[90m├── [0mREADME.md
[90m│   [0m[90m├── [0mcompiled.o
[90m│   [0m[90m├── [0mcompressed.zip
[90m│   [0m[90m├── [0mcrypto.asc
[90m│   [0m[90m├── [0mdocument.xlsx
[90m│   [0m[90m├── [0mexec
[90m│   [0m[90m├── [0mimage.png
[90m│   [0m[90m├── [0mmusic.mp3
[90m│   [0m[90m├── [0mtmp.bk
[90m│   [0m[90m├── [0mvideo.mp4
[90m│   [0m[90m└── [0mwav.wav
[90m├── [0mcorge
[90m├── [0mfoo
[90m│   [0m[90m├── [0mbar
[90m│   [0m[90m│   [0m[90m└── [0mbaz
[90m│   [0m[90m├── [0mquux
[90m│   [0m[90m└── [0mqux
[90m├── [0mgrault
[90m│   [0m[90m├── [0mgarply
[90m│   [0m[90m│   [0m[90m├── [0mfred
[90m│   [0m[90m│   [0m[90m└── [0mwaldo
[90m│   [0m[90m│   [0m    [90m├── [0mwibble
[90m│   [0m[90m│   [0m    [90m└── [0mwobble
[90m│   [0m[90m└── [0mplugh
[90m└── [0mxyzzy
    [90m└── [0mthud
        [90m├── [0mflob
        [90m└── [0mwubble

9 directories, 23 files`,
			colored:    false,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: true,
			datetime:   false,
		},
		{
			// This test case was created for "github actions". uid has a value according to it.
			// TODO: allow user group to be specified.
			name: "gotree --uid --gid <directory>",
			want: `tmp
[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [94m01[0m/
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [33mcompiled.o[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [92mexec[0m*
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [95mimage.png[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [90mtmp.bk[0m
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[[33mrunner[0m [33mdocker[0m]  [35mwav.wav[0m
[90m├── [0m[[33mrunner[0m [33mdocker[0m]  corge
[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [94mfoo[0m/
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0m[[33mrunner[0m [33mdocker[0m]  baz
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  quux
[90m│   [0m[90m└── [0m[[33mrunner[0m [33mdocker[0m]  qux
[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [94mgrault[0m/
[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  [94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0m[[33mrunner[0m [33mdocker[0m]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[[33mrunner[0m [33mdocker[0m]  [94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0m[[33mrunner[0m [33mdocker[0m]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[[33mrunner[0m [33mdocker[0m]  wobble
[90m│   [0m[90m└── [0m[[33mrunner[0m [33mdocker[0m]  plugh
[90m└── [0m[[33mrunner[0m [33mdocker[0m]  [94mxyzzy[0m/
    [90m└── [0m[[33mrunner[0m [33mdocker[0m]  [94mthud[0m/
        [90m├── [0m[[33mrunner[0m [33mdocker[0m]  flob
        [90m└── [0m[[33mrunner[0m [33mdocker[0m]  wubble

8 directories, 21 files`,
			colored:    true,
			level:      math.MaxInt64,
			permission: false,
			uid:        true,
			gid:        true,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			// This test case was created for "github actions". uid has a value according to it.
			// TODO: allow user group to be specified.
			name: "gotree --disable-color --uid --gid <directory>",
			want: `tmp
[90m├── [0m[runner docker]  01
[90m│   [0m[90m├── [0m[runner docker]  README.md
[90m│   [0m[90m├── [0m[runner docker]  compiled.o
[90m│   [0m[90m├── [0m[runner docker]  compressed.zip
[90m│   [0m[90m├── [0m[runner docker]  crypto.asc
[90m│   [0m[90m├── [0m[runner docker]  document.xlsx
[90m│   [0m[90m├── [0m[runner docker]  exec
[90m│   [0m[90m├── [0m[runner docker]  image.png
[90m│   [0m[90m├── [0m[runner docker]  music.mp3
[90m│   [0m[90m├── [0m[runner docker]  tmp.bk
[90m│   [0m[90m├── [0m[runner docker]  video.mp4
[90m│   [0m[90m└── [0m[runner docker]  wav.wav
[90m├── [0m[runner docker]  corge
[90m├── [0m[runner docker]  foo
[90m│   [0m[90m├── [0m[runner docker]  bar
[90m│   [0m[90m│   [0m[90m└── [0m[runner docker]  baz
[90m│   [0m[90m├── [0m[runner docker]  quux
[90m│   [0m[90m└── [0m[runner docker]  qux
[90m├── [0m[runner docker]  grault
[90m│   [0m[90m├── [0m[runner docker]  garply
[90m│   [0m[90m│   [0m[90m├── [0m[runner docker]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[runner docker]  waldo
[90m│   [0m[90m│   [0m    [90m├── [0m[runner docker]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[runner docker]  wobble
[90m│   [0m[90m└── [0m[runner docker]  plugh
[90m└── [0m[runner docker]  xyzzy
    [90m└── [0m[runner docker]  thud
        [90m├── [0m[runner docker]  flob
        [90m└── [0m[runner docker]  wubble

8 directories, 21 files`,
			colored:    false,
			level:      math.MaxInt64,
			permission: false,
			uid:        true,
			gid:        true,
			size:       false,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --size <directory>",
			want: `tmp
[90m├── [0m[-]  [94m01[0m/
[90m│   [0m[90m├── [0m[[32m0[0m]  [4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [33mcompiled.o[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [92mexec[0m*
[90m│   [0m[90m├── [0m[[32m0[0m]  [95mimage.png[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [90mtmp.bk[0m
[90m│   [0m[90m├── [0m[[32m0[0m]  [35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[[32m0[0m]  [35mwav.wav[0m
[90m├── [0m[[32m0[0m]  corge
[90m├── [0m[-]  [94mfoo[0m/
[90m│   [0m[90m├── [0m[-]  [94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0m[[32m0[0m]  baz
[90m│   [0m[90m├── [0m[[32m0[0m]  quux
[90m│   [0m[90m└── [0m[[32m0[0m]  qux
[90m├── [0m[-]  [94mgrault[0m/
[90m│   [0m[90m├── [0m[-]  [94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0m[[32m0[0m]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[-]  [94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0m[[32m0[0m]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[[32m0[0m]  wobble
[90m│   [0m[90m└── [0m[[32m0[0m]  plugh
[90m└── [0m[-]  [94mxyzzy[0m/
    [90m└── [0m[-]  [94mthud[0m/
        [90m├── [0m[[32m0[0m]  flob
        [90m└── [0m[[32m0[0m]  wubble

8 directories, 21 files`,
			colored:    true,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       true,
			includeDot: false,
			datetime:   false,
		},
		{
			name: "gotree --disable-color --size <directory>",
			want: `tmp
[90m├── [0m[-]  01
[90m│   [0m[90m├── [0m[0]  README.md
[90m│   [0m[90m├── [0m[0]  compiled.o
[90m│   [0m[90m├── [0m[0]  compressed.zip
[90m│   [0m[90m├── [0m[0]  crypto.asc
[90m│   [0m[90m├── [0m[0]  document.xlsx
[90m│   [0m[90m├── [0m[0]  exec
[90m│   [0m[90m├── [0m[0]  image.png
[90m│   [0m[90m├── [0m[0]  music.mp3
[90m│   [0m[90m├── [0m[0]  tmp.bk
[90m│   [0m[90m├── [0m[0]  video.mp4
[90m│   [0m[90m└── [0m[0]  wav.wav
[90m├── [0m[0]  corge
[90m├── [0m[-]  foo
[90m│   [0m[90m├── [0m[-]  bar
[90m│   [0m[90m│   [0m[90m└── [0m[0]  baz
[90m│   [0m[90m├── [0m[0]  quux
[90m│   [0m[90m└── [0m[0]  qux
[90m├── [0m[-]  grault
[90m│   [0m[90m├── [0m[-]  garply
[90m│   [0m[90m│   [0m[90m├── [0m[0]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[-]  waldo
[90m│   [0m[90m│   [0m    [90m├── [0m[0]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[0]  wobble
[90m│   [0m[90m└── [0m[0]  plugh
[90m└── [0m[-]  xyzzy
    [90m└── [0m[-]  thud
        [90m├── [0m[0]  flob
        [90m└── [0m[0]  wubble

8 directories, 21 files`,
			colored:    false,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       true,
			includeDot: false,
			datetime:   false,
		},
		{
			name:       "gotree -D <directory>",
			want:       testCaseWithDate,
			colored:    true,
			level:      math.MaxInt64,
			permission: false,
			uid:        false,
			gid:        false,
			size:       false,
			includeDot: false,
			datetime:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpStdout := os.Stdout

			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Tree(TMP_DIR, tt.colored, tt.level, tt.permission, tt.uid, tt.gid, tt.size, tt.includeDot, tt.datetime)
			if err != nil {
				t.Fatal(err)
			}
			w.Close()

			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatal(err)
			}

			got := strings.TrimRight(buf.String(), "\n")

			os.Stdout = tmpStdout

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Stdout missmatch (-got +want):\n%s", diff)
			}
		})

	}

}

func testCaseWithDate() (string, error) {
	testCase := `tmp
[90m├── [0m[[34m__DATETIME__[0m]  [94m01[0m/
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [4m[93mREADME.md[0m[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [33mcompiled.o[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [31mcompressed.zip[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [96mcrypto.asc[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [32mdocument.xlsx[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [92mexec[0m*
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [95mimage.png[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [35mmusic.mp3[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [90mtmp.bk[0m
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [35mvideo.mp4[0m
[90m│   [0m[90m└── [0m[[34m__DATETIME__[0m]  [35mwav.wav[0m
[90m├── [0m[[34m__DATETIME__[0m]  corge
[90m├── [0m[[34m__DATETIME__[0m]  [94mfoo[0m/
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [94mbar[0m/
[90m│   [0m[90m│   [0m[90m└── [0m[[34m__DATETIME__[0m]  baz
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  quux
[90m│   [0m[90m└── [0m[[34m__DATETIME__[0m]  qux
[90m├── [0m[[34m__DATETIME__[0m]  [94mgrault[0m/
[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  [94mgarply[0m/
[90m│   [0m[90m│   [0m[90m├── [0m[[34m__DATETIME__[0m]  fred
[90m│   [0m[90m│   [0m[90m└── [0m[[34m__DATETIME__[0m]  [94mwaldo[0m/
[90m│   [0m[90m│   [0m    [90m├── [0m[[34m__DATETIME__[0m]  wibble
[90m│   [0m[90m│   [0m    [90m└── [0m[[34m__DATETIME__[0m]  wobble
[90m│   [0m[90m└── [0m[[34m__DATETIME__[0m]  plugh
[90m└── [0m[[34m__DATETIME__[0m]  [94mxyzzy[0m/
    [90m└── [0m[[34m__DATETIME__[0m]  [94mthud[0m/
        [90m├── [0m[[34m__DATETIME__[0m]  flob
        [90m└── [0m[[34m__DATETIME__[0m]  wubble

8 directories, 21 files`

	modTime, err := modTime(TMP_DIR)
	if err != nil {
		return "", err
	}

	testCase = strings.Replace(testCase, "__DATETIME__", modTime, -1)

	return testCase, nil
}

func modTime(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	modTime := fi.ModTime().Format("2006-01-02 15:04")

	return modTime, nil
}
