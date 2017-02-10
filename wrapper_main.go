package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mohae/linewrap"
)

// usage is the usage func for flag.Usage.
func usage() {
	fmt.Fprint(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [FLAGS] path\n", app)
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s takes a path and wraps its contents so that the lines are of a certain length.\n", app)
	fmt.Fprint(os.Stderr, "  If the path is a file, the contents of that file will be wrapped and saved.\n")
	fmt.Fprint(os.Stderr, "  If the path is a directory, all files within that directory will be wrapped and saved.\n")
	fmt.Fprint(os.Stderr, "  Directory operations are not recursive.\n")
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprint(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

// FlagParse handles flag parsing, validation, and any side affects of flag
// states. Errors or invalid states should result in printing a message to
// os.Stderr and an os.Exit() with a non-zero int.
func FlagParse() {
	var err error

	flag.Parse()

	if cfg.LogFile != "" && cfg.LogFile != "stdout" { // open the logfile if one is specified
		cfg.f, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: open logfile: %s", app, err)
			os.Exit(1)
		}
	}
}

func wrapperMain(paths []string) int {
	/// this combination ensures that wrapped lines have leading blanks elided
	linewrap.LineComment(true, "")
	if cfg.f != nil {
		defer cfg.f.Close() // make sure the logfile is closed if there is one
	}

	// process the path
	for _, p := range paths {
		err := dir(p)
		if err != nil {
			log.Printf("wrap %s: error: %s", p, err)
			return 1
		}
	}

	return 0
}

func dir(p string) error {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, f := range files {
		b, err := ioutil.ReadFile(filepath.Join(p, f.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %s", filepath.Join(p, f.Name()), err)
		}

		// wrap the bytes
		b, err = linewrap.Bytes(b)
		if err != nil {
			return fmt.Errorf("wrap %s: %s", filepath.Join(p, f.Name()), err)
		}

		// write the bytes
		err = ioutil.WriteFile(filepath.Join(p, f.Name()), b, f.Mode())
		if err != nil {
			return fmt.Errorf("write %s: %s", filepath.Join(p, f.Name()), err)
		}

	}

	return nil
}
