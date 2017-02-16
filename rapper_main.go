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
	fmt.Fprintf(os.Stderr, "  %s [FLAGS] path(s)\n", app)
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

	if flag.NArg() <= 0 {
		fmt.Fprintf(os.Stderr, "%s: no paths specified, nothing to wrap\n", app)
		fmt.Fprint(os.Stderr, "\n")
		usage()
		os.Exit(1)
	}

	if cfg.Exclude && cfg.Include {
		fmt.Fprintf(os.Stderr, "%s: include and exclude flags were both set to true; these are mutually exclusive flags\n", app)
		os.Exit(1)
	}

	if cfg.Comment != "" {
		cfg.commentStyle = linewrap.StringAsCommentStyle(cfg.Comment)
		if cfg.commentStyle == linewrap.NoComment {
			fmt.Fprintf(os.Stderr, "%s: invalid comment style: %q\n", app, cfg.Comment)
			os.Exit(1)
		}
	}

	if cfg.Verbose && cfg.commentStyle != linewrap.NoComment {
		fmt.Printf("%s: the files will be wrapped as %s\n", app, cfg.commentStyle)
	}

	if cfg.LogFile != "" && cfg.LogFile != "stderr" { // open the logfile if one is specified
		cfg.f, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: open logfile: %s\n", app, err)
			os.Exit(1)
		}
	}
}

func wrapperMain(paths []string) int {
	if cfg.f != nil {
		defer cfg.f.Close() // make sure the logfile is closed if there is one
	}

	var processed, updated int
	// process the path
	for _, p := range paths {
		p, u, err := dir(p)
		processed += p
		updated += u
		if err != nil {
			doneMessage(processed, updated)
			log.Printf("wrap %s: error: %s", p, err)
			return 1
		}

	}

	doneMessage(processed, updated)
	return 0
}

func dir(p string) (processed, updated int, err error) {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return processed, updated, err
	}
	// wrap the bytes
	w := linewrap.New()
	w.CommentStyle = cfg.commentStyle

	for _, f := range files {
		w.Reset()
		processed++
		ext := filepath.Ext(f.Name())
		if cfg.Include {
			if !cfg.Ext.Exists(ext) {
				if cfg.Verbose {
					fmt.Printf("file skipped:    %s\n", f.Name())
				}
				continue
			}
		}
		if cfg.Exclude {
			if cfg.Ext.Exists(ext) {
				if cfg.Verbose {
					fmt.Printf("file skipped:    %s\n", f.Name())
				}
				continue
			}
		}

		b, err := ioutil.ReadFile(filepath.Join(p, f.Name()))
		if err != nil {
			return processed, updated, fmt.Errorf("read %s: %s", filepath.Join(p, f.Name()), err)
		}

		b, err = w.Bytes(b)
		if err != nil {
			return processed, updated, fmt.Errorf("wrap %s: %s", filepath.Join(p, f.Name()), err)
		}

		// write the bytes
		err = ioutil.WriteFile(filepath.Join(p, f.Name()), b, f.Mode())
		if err != nil {
			return processed, updated, fmt.Errorf("write %s: %s", filepath.Join(p, f.Name()), err)
		}
		updated++
		if cfg.Verbose {
			fmt.Printf("file processed:  %s\n", f.Name())
		}
	}

	return processed, updated, nil
}

func doneMessage(p, u int) {
	s := p - u

	dash := make([]byte, 0, len(app)+21)
	for i := 0; i < cap(dash); i++ {
		dash = append(dash, '-')
	}
	if cfg.Verbose {
		fmt.Print("\n")
	}
	fmt.Printf("%s: processing complete\n", app)
	fmt.Printf("%s\n", string(dash))
	fmt.Printf("%d files processed\n", p)
	fmt.Printf("%d were skipped\n", s)
	fmt.Printf("%d were updated\n", u)
}
