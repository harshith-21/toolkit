# toolkit
GoLang module made be ✨me✨

> this module makes a random string of variable length given

eg:
```go
package main

import (
	"fmt"
	"github.com/harshith-21/toolkit"
)

func main() {
	var tools toolkit.Tools

	s := tools.RandomString(5)
	fmt.Println(s)
}
```

A simple example of how to create a reusable Go module with commonly used tools.

The included tools are:

- [X] Read JSON
- [X] Write JSON
- [X] Produce a JSON encoded error response
- [X] Upload a file to a specified directory
- [X] Download a static file
- [X] Get a random string of length n
- [X] Post JSON to a remote service 
- [X] Create a directory, including all parent directories, if it does not already exist
- [X] Create a URL safe slug from a string

## Installation

`go get -u github.com/harshith-21/toolbox`