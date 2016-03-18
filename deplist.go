// Command deplist prints a detailed list of package imports for a package.
// The output contains the path to the source on the GOPATH for those imported packages.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var (
	tags   = flag.String("tags", "", "comma-separated list of build tags to apply")
	goroot = flag.Bool("goroot", false, "include imports in GOROOT")
	tsv    = flag.Bool("tsv", false, "use only a single tab between columns")
)

type flusher interface {
	Flush() error
}

type importFrom struct {
	path    string
	fromDir string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: deplist [-tags] [-goroot] <dirs...>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	buildctx := build.Default
	buildctx.BuildTags = strings.Split(*tags, ",")

	var w io.Writer = os.Stdout
	if !*tsv {
		w = tabwriter.NewWriter(os.Stdout, 8, 0, 4, ' ', 0)
	}

	visited := make(map[importFrom]bool)
	var imports []importFrom

	for _, dir := range flag.Args() {
		abs, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("could not get absolute path for dir %q: %v", dir, err)
		}

		pkg, err := buildctx.ImportDir(abs, 0)
		if err != nil {
			log.Fatalf("could not get package for dir %q: %v", dir, err)
		}

		for _, importPath := range pkg.Imports {
			imports = append(imports, importFrom{
				path:    importPath,
				fromDir: abs,
			})
		}
	}

	for len(imports) != 0 {
		i := imports[0]
		imports = imports[1:] // shift

		if _, ok := visited[i]; ok || i.path == "C" {
			continue
		}
		visited[i] = true

		pkg, err := buildctx.Import(i.path, i.fromDir, 0)
		if err != nil {
			log.Fatalf("could not get package %q, imported from %q: %v", i.path, i.fromDir, err)
		}

		if !*goroot && pkg.Goroot {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", i.fromDir, i.path, pkg.SrcRoot, pkg.ImportPath)

		for _, importPath := range pkg.Imports {
			imports = append(imports, importFrom{
				path:    importPath,
				fromDir: pkg.Dir,
			})
		}
	}

	if f, ok := w.(flusher); ok {
		f.Flush()
	}
}
