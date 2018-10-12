# GoDerive

Add derive comment above your type, and generate source code for the marked type.

## Install

### via go get

```
$ go get -u -v github.com/nextzhou/goderive
```

### via makefile

In this way you can get git version information

```
$ go get -u -v github.com/nextzhou/goderive
$ cd $GOPATH/src/github.com/nextzhou/goderive
$ make
$ goderive --verion # show git version information
Version: xxxxxxx
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
  goderive [flags] [path ...] # where a '/...' suffix includes all sub-directories
  goderive help [plugin ...]

Flags:
  -d, --delete                delete existing generated file when no derived type (default true)
  -D, --exclude-dir strings   exclude the given comma separated directories (default [vendor])
  -E, --exclude-ext strings   exclude the files having given file name ext (default [.gen.go,_test.go])
  -h, --help                  help for goderive
  -o, --output string         output file name (default "derived.gen.go")
  -v, --version               show version information

Plugins:
  set            set collection
  access         access fields for struct type
```

## Plugins

```
$ goderive help set
Plugin: set

set collection

Flags:
  Export         force the generated code to be exported/unexported

Args:
  Rename         single value                            assign set type name manually
  Order          single value    [Unstable Append Key]   keep order(default: Unstable)
```

```
$ goderive help access
Plugin: access

access fields for struct type

Args:
  Receiver       single value            receiver of methods
```
