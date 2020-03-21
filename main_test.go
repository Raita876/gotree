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
		if err := os.MkdirAll(dir, 0777); err != nil {
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
	tests := []struct {
		name    string
		want    string
		colored bool
		level   int
	}{
		{
			name: "gotree --disable-color <dir>",
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
			colored: false,
			level:   math.MaxInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpStdout := os.Stdout

			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Tree(TMP_DIR, tt.colored, tt.level)
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
				t.Fatalf("Stdout missmatch (-got +want):\n%s", diff)
			}
		})

	}

}
