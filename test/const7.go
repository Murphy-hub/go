// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check that the compiler refuses excessively long constants.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// testProg creates a package called name, with path dir/name.go,
// which declares an untyped constant of the given length.
// testProg compiles this package and checks for the absence or
// presence of a constant literal error.
func testProg(dir, name string, length int, ok bool) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf,
		"package %s; const _ = %s // %d digits",
		name, strings.Repeat("9", length), length,
	)

	filename := filepath.Join(dir, fmt.Sprintf("%s.go", name))
	if err := os.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "tool", "compile", filename)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if ok {
		// no error expected
		if err != nil {
			log.Fatalf("%s: compile failed unexpectedly: %v", name, err)
		}
		return
	}

	// error expected
	if err == nil {
		log.Fatalf("%s: compile succeeded unexpectedly", name)
	}
	if !bytes.Contains(output, []byte("excessively long constant")) {
		log.Fatalf("%s: wrong compiler error message:\n%s\n", name, output)
	}
}

func main() {
	if runtime.GOOS == "js" || runtime.Compiler != "gc" {
		return
	}

	dir, err := ioutil.TempDir("", "const7_")
	if err != nil {
		log.Fatalf("creating temp dir: %v\n", err)
	}
	defer os.RemoveAll(dir)

	const limit = 10000 // compiler-internal constant length limit
	testProg(dir, "x1", limit, true)
	testProg(dir, "x2", limit+1, false)
}
