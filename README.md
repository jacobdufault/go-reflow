# go-reflow

A simple utility to reformat comments. Accepts input on stdin, writes output
to stdout.

# Install

```sh
$ go install github.com/jacobdufault/go-reflow
$ go-reflow --help
```

# Usage

```
Usage of go-reflow:
  -tab-size int
        Size of a tab (\t) character (default 2)
  -width int
        Width to reflow text to (default 80)
```

# Example

```
$ echo " # foo bar baz a baaa barara" | ./go-reflow -width 20
 # foo bar baz a
 # baaa barara
```