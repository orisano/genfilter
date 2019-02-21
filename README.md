# genfilter
genfilter is a tool for collect to generated files and filter to them in stream.

## Installation
```
go get github.com/orisano/genfilter
```

## How to use
```
$ genfilter
genfilter: subcommand is required:
Available SubCommands:
 - apply
 - build
```

```
$ genfilter build -h
Usage of build:
  -d string
    	root directory (default ".")
  -o string
    	output filter binary path (default "filter.gob")
```

```
$ genfilter apply -h
Usage of apply:
  -f string
    	filter binary path (default "filter.gob")
  -i string
    	input file (default "-")
```

## Author
Nao YONASHIRO (@orisano)

## License
MIT
