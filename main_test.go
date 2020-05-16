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
├── corge
├── [34mfoo[0m
│   ├── [34mbar[0m
│   │   └── baz
│   ├── quux
│   └── qux
├── [34mgrault[0m
│   ├── [34mgarply[0m
│   │   ├── fred
│   │   └── [34mwaldo[0m
│   │       ├── wibble
│   │       └── wobble
│   └── plugh
└── [34mxyzzy[0m
    └── [34mthud[0m
        ├── flob
        └── wubble

7 directories, 10 files`,
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
├── corge
├── foo
│   ├── bar
│   │   └── baz
│   ├── quux
│   └── qux
├── grault
│   ├── garply
│   │   ├── fred
│   │   └── waldo
│   │       ├── wibble
│   │       └── wobble
│   └── plugh
└── xyzzy
    └── thud
        ├── flob
        └── wubble

7 directories, 10 files`,
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
├── corge
├── [34mfoo[0m
│   ├── [34mbar[0m
│   ├── quux
│   └── qux
├── [34mgrault[0m
│   ├── [34mgarply[0m
│   └── plugh
└── [34mxyzzy[0m
    └── [34mthud[0m

6 directories, 4 files`,
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
├── corge
├── foo
│   ├── bar
│   ├── quux
│   └── qux
├── grault
│   ├── garply
│   └── plugh
└── xyzzy
    └── thud

6 directories, 4 files`,
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
├── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  corge
├── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mfoo[0m
│   ├── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mbar[0m
│   │   └── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  baz
│   ├── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  quux
│   └── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  qux
├── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mgrault[0m
│   ├── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mgarply[0m
│   │   ├── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  fred
│   │   └── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mwaldo[0m
│   │       ├── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wibble
│   │       └── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wobble
│   └── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  plugh
└── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mxyzzy[0m
    └── [[34md[0m[33mr[0m[31mw[0m[32mx[0m[33mr[0m-[32mx[0m[33mr[0m-[32mx[0m]  [34mthud[0m
        ├── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  flob
        └── [.[33mr[0m[31mw[0m-[33mr[0m--[33mr[0m--]  wubble

7 directories, 10 files`,
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
├── [.rw-r--r--]  corge
├── [drwxr-xr-x]  foo
│   ├── [drwxr-xr-x]  bar
│   │   └── [.rw-r--r--]  baz
│   ├── [.rw-r--r--]  quux
│   └── [.rw-r--r--]  qux
├── [drwxr-xr-x]  grault
│   ├── [drwxr-xr-x]  garply
│   │   ├── [.rw-r--r--]  fred
│   │   └── [drwxr-xr-x]  waldo
│   │       ├── [.rw-r--r--]  wibble
│   │       └── [.rw-r--r--]  wobble
│   └── [.rw-r--r--]  plugh
└── [drwxr-xr-x]  xyzzy
    └── [drwxr-xr-x]  thud
        ├── [.rw-r--r--]  flob
        └── [.rw-r--r--]  wubble

7 directories, 10 files`,
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
├── .aaa
├── [34m.bbb[0m
│   └── .ccc
├── corge
├── [34mfoo[0m
│   ├── [34mbar[0m
│   │   └── baz
│   ├── quux
│   └── qux
├── [34mgrault[0m
│   ├── [34mgarply[0m
│   │   ├── fred
│   │   └── [34mwaldo[0m
│   │       ├── wibble
│   │       └── wobble
│   └── plugh
└── [34mxyzzy[0m
    └── [34mthud[0m
        ├── flob
        └── wubble

8 directories, 12 files`,
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
├── .aaa
├── .bbb
│   └── .ccc
├── corge
├── foo
│   ├── bar
│   │   └── baz
│   ├── quux
│   └── qux
├── grault
│   ├── garply
│   │   ├── fred
│   │   └── waldo
│   │       ├── wibble
│   │       └── wobble
│   └── plugh
└── xyzzy
    └── thud
        ├── flob
        └── wubble

8 directories, 12 files`,
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
├── [[33mrunner[0m [33mdocker[0m]  corge
├── [[33mrunner[0m [33mdocker[0m]  [34mfoo[0m
│   ├── [[33mrunner[0m [33mdocker[0m]  [34mbar[0m
│   │   └── [[33mrunner[0m [33mdocker[0m]  baz
│   ├── [[33mrunner[0m [33mdocker[0m]  quux
│   └── [[33mrunner[0m [33mdocker[0m]  qux
├── [[33mrunner[0m [33mdocker[0m]  [34mgrault[0m
│   ├── [[33mrunner[0m [33mdocker[0m]  [34mgarply[0m
│   │   ├── [[33mrunner[0m [33mdocker[0m]  fred
│   │   └── [[33mrunner[0m [33mdocker[0m]  [34mwaldo[0m
│   │       ├── [[33mrunner[0m [33mdocker[0m]  wibble
│   │       └── [[33mrunner[0m [33mdocker[0m]  wobble
│   └── [[33mrunner[0m [33mdocker[0m]  plugh
└── [[33mrunner[0m [33mdocker[0m]  [34mxyzzy[0m
    └── [[33mrunner[0m [33mdocker[0m]  [34mthud[0m
        ├── [[33mrunner[0m [33mdocker[0m]  flob
        └── [[33mrunner[0m [33mdocker[0m]  wubble

7 directories, 10 files`,
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
├── [runner docker]  corge
├── [runner docker]  foo
│   ├── [runner docker]  bar
│   │   └── [runner docker]  baz
│   ├── [runner docker]  quux
│   └── [runner docker]  qux
├── [runner docker]  grault
│   ├── [runner docker]  garply
│   │   ├── [runner docker]  fred
│   │   └── [runner docker]  waldo
│   │       ├── [runner docker]  wibble
│   │       └── [runner docker]  wobble
│   └── [runner docker]  plugh
└── [runner docker]  xyzzy
    └── [runner docker]  thud
        ├── [runner docker]  flob
        └── [runner docker]  wubble

7 directories, 10 files`,
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
├── [[32m0[0m]  corge
├── [-]  [34mfoo[0m
│   ├── [-]  [34mbar[0m
│   │   └── [[32m0[0m]  baz
│   ├── [[32m0[0m]  quux
│   └── [[32m0[0m]  qux
├── [-]  [34mgrault[0m
│   ├── [-]  [34mgarply[0m
│   │   ├── [[32m0[0m]  fred
│   │   └── [-]  [34mwaldo[0m
│   │       ├── [[32m0[0m]  wibble
│   │       └── [[32m0[0m]  wobble
│   └── [[32m0[0m]  plugh
└── [-]  [34mxyzzy[0m
    └── [-]  [34mthud[0m
        ├── [[32m0[0m]  flob
        └── [[32m0[0m]  wubble

7 directories, 10 files`,
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
├── [0]  corge
├── [-]  foo
│   ├── [-]  bar
│   │   └── [0]  baz
│   ├── [0]  quux
│   └── [0]  qux
├── [-]  grault
│   ├── [-]  garply
│   │   ├── [0]  fred
│   │   └── [-]  waldo
│   │       ├── [0]  wibble
│   │       └── [0]  wobble
│   └── [0]  plugh
└── [-]  xyzzy
    └── [-]  thud
        ├── [0]  flob
        └── [0]  wubble

7 directories, 10 files`,
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
├── [[34m__DATETIME__[0m]  corge
├── [[34m__DATETIME__[0m]  [34mfoo[0m
│   ├── [[34m__DATETIME__[0m]  [34mbar[0m
│   │   └── [[34m__DATETIME__[0m]  baz
│   ├── [[34m__DATETIME__[0m]  quux
│   └── [[34m__DATETIME__[0m]  qux
├── [[34m__DATETIME__[0m]  [34mgrault[0m
│   ├── [[34m__DATETIME__[0m]  [34mgarply[0m
│   │   ├── [[34m__DATETIME__[0m]  fred
│   │   └── [[34m__DATETIME__[0m]  [34mwaldo[0m
│   │       ├── [[34m__DATETIME__[0m]  wibble
│   │       └── [[34m__DATETIME__[0m]  wobble
│   └── [[34m__DATETIME__[0m]  plugh
└── [[34m__DATETIME__[0m]  [34mxyzzy[0m
    └── [[34m__DATETIME__[0m]  [34mthud[0m
        ├── [[34m__DATETIME__[0m]  flob
        └── [[34m__DATETIME__[0m]  wubble

7 directories, 10 files`

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
