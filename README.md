[![Go Reference](https://pkg.go.dev/badge/github.com/juanbrotenelle/go_salamoonder.svg)](https://pkg.go.dev/github.com/juanbrotenelle/go_salamoonder)

# go_salamoonder

The go_salamoonder package provides a complete set of methods for interacting with the [Salamoonder API](https://salamoonder.com/).

## Installation

```bash
go get -u github.com/juanbrotenelle/go_salamoonder
```

## Getting started

### Simple examples

#### Get account balance

```go
package main

import (
	"fmt"
	"log"
	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	api := salamoonder.New("sr-YOUR-API-KEY")

	balance, err := api.Balance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)
	// 499.81070
}
```

#### Polling Kasada Solve

```go
package main

import (
	"fmt"
	"time"
	"log"
	"errors"
	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	api := salamoonder.New("sr-YOUR-API-KEY")

	taskId, err := api.KasadaCreate("...862e0f06eea3/p.js", true)
	if err != nil {
		fmt.Println(err)
	}

	var result salamoonder.KasadaSolution
	for {
		result, err = api.Kasada(taskId)
		if err == nil && result.UserAgent != "" {
			fmt.Println(res)
			return
		}
		if errors.Is(err, salamoonder.ErrTaskNotReady) {
			time.Sleep(2 * time.Second)
			continue
		}
		log.Fatal(err)
	}
	// {xxx xxx xxx xx xxx xxx xxx 1704673650160}
}
```

#### Extract p.js from site

```go
package main

import (
	"fmt"
	"log"
	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	pjs, err := salamoonder.FindPJS("https://nike.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pjs)
	// https://www.nike.com/149e9513-01fa-4fb0-aad4-566afd725d1b/2d206a39-8ed7-437e-a3be-862e0f06eea3/p.js
}
```
