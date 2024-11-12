package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/traefik/yaegi/extract"
)

func extractCmd(arg []string) error {
	profPath := "/tmp/pprof/extract.prof"
	if err := os.MkdirAll(filepath.Dir(profPath), 0770); err != nil {
		return err
	}
	f, err := os.Create(profPath)
	if err != nil {
		return err
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var licensePath string
	var dest string
	var filename string
	var name string
	var exclude string
	var include string
	var tag string

	eflag := flag.NewFlagSet("run", flag.ContinueOnError)
	eflag.StringVar(&licensePath, "license", "", "path to a LICENSE file")
	eflag.StringVar(&dest, "dest", "", "the package")
	eflag.StringVar(&filename, "filename", "", "the filename for the extracted symbols")
	eflag.StringVar(&name, "name", "", "the namespace for the extracted symbols")
	eflag.StringVar(&exclude, "exclude", "", "comma separated list of regexp matching symbols to exclude")
	eflag.StringVar(&include, "include", "", "comma separated list of regexp matching symbols to include")
	eflag.StringVar(&tag, "tag", "", "comma separated list of build tags to be added to the created package")
	eflag.Usage = func() {
		fmt.Println("Usage: yaegi extract [options] packages...")
		fmt.Println("Options:")
		eflag.PrintDefaults()
	}

	if err := eflag.Parse(arg); err != nil {
		return err
	}

	args := eflag.Args()
	if len(args) == 0 {
		return errors.New("missing package")
	}

	license, err := genLicense(licensePath)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if name == "" {
		name = filepath.Base(wd)
	}
	if dest == "" {
		dest = name
	}
	ext := extract.Extractor{
		Dest:    dest,
		License: license,
	}
	if tag != "" {
		ext.Tag = strings.Split(tag, ",")
	}

	if exclude != "" {
		ext.Exclude = strings.Split(exclude, ",")
	}
	if include != "" {
		ext.Include = strings.Split(include, ",")
	}

	r := strings.NewReplacer("/", "-", ".", "_", "~", "_")

	for _, pkgIdent := range args {
		fmt.Printf("Extracting from package: %q...\n", pkgIdent)
		var buf bytes.Buffer
		importPath, err := ext.Extract(pkgIdent, name, &buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if filename != "" {
			importPath = filename
		}

		fmt.Printf("Extracted import path: %q...\n", importPath)
		oFile := r.Replace(importPath) + ".go"
		f, err := os.Create(oFile)
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, &buf); err != nil {
			_ = f.Close()
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

// genLicense generates the correct LICENSE header text from the provided
// path to a LICENSE file.
func genLicense(fname string) (string, error) {
	if fname == "" {
		return "", nil
	}

	f, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("could not open LICENSE file: %w", err)
	}
	defer func() { _ = f.Close() }()

	license := new(strings.Builder)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		txt := sc.Text()
		if txt != "" {
			txt = " " + txt
		}
		license.WriteString("//" + txt + "\n")
	}
	if sc.Err() != nil {
		return "", fmt.Errorf("could not scan LICENSE file: %w", err)
	}

	return license.String(), nil
}
