# deplist

deplist prints a detailed list of package imports for a package.
The output contains the path to the source on the GOPATH for those imported packages.

## Installation

```
$ go get github.com/broady/deplist
```

## Usage

```
usage: deplist [-tags] [-goroot] <dirs...>
  -goroot
      include imports in GOROOT
  -tags string
      comma-separated list of build tags to apply
  -tsv
      use only a single tab between columns
```

## Support

This is not an official Google product.

## License

deplist is available under the Apache 2 license. See [LICENSE](LICENSE).
