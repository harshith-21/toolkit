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
