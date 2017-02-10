package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mohae/linewrap"
)

//const text = `MIT License
//Copyright (c) <year> <copyright holders>
//
const text = `Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`

// there will be trailing spaces until linewrap implements a lexer based wrapper or
// some other method of eliding trailing spaces is added.
//const expected = `MIT License
//Copyright (c) <year> <copyright holders>
const expected = `Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in 
the Software without restriction, including without limitation the rights to 
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies 
of the Software, and to permit persons to whom the Software is furnished to do 
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all 
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR 
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, 
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE 
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER 
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, 
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE 
SOFTWARE.
`

func TestWrapDir(t *testing.T) {
	linewrap.LineComment(true, "")
	d, err := ioutil.TempDir("", app)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		err = ioutil.WriteFile(filepath.Join(d, fmt.Sprintf("mit%d", i)), []byte(text), 0664)
		if err != nil {
			os.RemoveAll(d)
			t.Fatal(err)
		}
	}

	err = dir(d)
	if err != nil {
		os.RemoveAll(d)
		t.Fatal(err)
	}

	files, err := ioutil.ReadDir(d)
	if err != nil {
		os.RemoveAll(d)
		t.Fatal(err)
	}

	want := strings.Split(expected, "\n")
	for _, f := range files {
		if f.IsDir() { // skip directories
			continue
		}
		b, err := ioutil.ReadFile(filepath.Join(d, f.Name()))
		if err != nil {
			os.RemoveAll(d)
			t.Fatal(err)
		}
		got := strings.Split(string(b), "\n")
		if len(got) != len(want) {
			t.Errorf("got %d lines; want %d", len(got), len(want))
		}

		for i, line := range got { // if the above errored; this will probably panic but it's fine
			if line != want[i] {
				t.Errorf("got\t%q\nwant\t%q", line, want[i])
			}
		}
	}
	os.RemoveAll(d)
}
