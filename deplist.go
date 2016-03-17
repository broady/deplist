// Command deplist prints a detailed list of package imports for a package.
// The output contains the path to the source on the GOPATH for those imported packages.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	tags   = flag.String("tags", "", "comma-separated list of build tags to apply")
	goroot = flag.Bool("goroot", false, "include imports in GOROOT")
	usage  = flag.Bool("h", false, "print usage")
)

func main() {
	flag.Parse()

	if *usage {
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: pkglist [-tags] [-h] <dirs...>")
		os.Exit(1)
	}

	buildctx := build.Default
	buildctx.BuildTags = strings.Split(*tags, ",")

	visited := make(map[imp]bool)
	var imports []imp

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
			i := imp{path: importPath, from: abs}
			imports = append(imports, i)
		}
	}

	for len(imports) != 0 {
		i := imports[0]
		imports = imports[1:] // shift

		if _, ok := visited[i]; ok || i.path == "C" {
			continue
		}
		visited[i] = true

		pkg, err := buildctx.Import(i.path, i.from, 0)
		if err != nil {
			log.Fatalf("could not get package %q, imported from %q: %v", i.path, i.from, err)
		}

		if !*goroot && pkg.Goroot {
			continue
		}

		fmt.Printf("%s\t%s\t%s\n", i.from, i.path, pkg.Dir)

		for _, i := range pkg.Imports {
			imports = append(imports, imp{path: i, from: pkg.Dir})
		}
	}
}

type imp struct {
	path, from string
}
