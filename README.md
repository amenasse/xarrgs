# Xarrgs

A toy implementation of Xargs written in Golang. It's purpose was
for me to practice some Go. Only a tiny subset of Xargs features are implemented.

## Installation

    go get github.com/amenasse/xarrgs

## Usage

```console
$ xarrgs -h

Usage of xarrgs:
  -max-chars int
    	Use at most max-chars per command line (default 2048)
  -max-procs int
    	Maximum number of processes to use (default 1)
  -null
    	items are seperated by a null not whitespace
```

## Examples


```console

# Perform work over 4 processes
$ cat input.txt | xarrgs  --max-procs 4 work

# Use at most max chars per command line
$ echo "An old silent pond..." | xarrgs --max-chars 10 echo
An old
silent
pond...


# Handle null terminated input
$ find . -type d -print0 | ./xarrgs -null du -s
```


## Lessons Learnt


1. A slice references an underlying array.  

    when passing a slice to a function or channel the same underlying array is still referenced by the receiver.

    From https://blog.golang.org/go-slices-usage-and-internals:

    > "A slice is a descriptor of an array segment. It consists of a pointer to the
    > array, the length of the segment, and its capacity (the maximum length of the
    > segment)."

2. GNU Xargs and quotes

    GNU Xargs treats quotes specially by default. They are removed from the
    output. Quotes must be matched. Quoted text is seen as a single argument.
