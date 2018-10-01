# GoDerive

Add derive comment above your type, and generate source code for the marked type.

## Install

```
go get -u -v github.com/nextzhou/goderive
```

## Usage

```
$ goderive help
GoDerive

Add derive comment above your type, and generate source code for the marked type.

Comment Format:
  // derive-<plugin>
  // derive-<plugin>: flag;!negative_flag;arg=single_value; arg2=val1,val2
  type YourType struct{/* ... */}

Usage:
  goderive [flags] [path ...]
  goderive help [plugin ...]

Flags:
  -d, --delete          delete existing generated file when no derived type (default true)
  -h, --help            help for goderive
  -o, --output string   output file name (default "derived.gen.go")

Plugins:
  set    set collection
```

## Plugins

```
$ goderive help set
Plugin: set

set collection

Flags:
  ByPoint        store elements by pointer

Args:
  Rename         single value                            assign set type name manually
  Order          single value    [Unstable Append Key]   keep order(default: Unstable)
```
