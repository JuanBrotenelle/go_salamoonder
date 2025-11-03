<div align="left">

[![Go Reference](https://pkg.go.dev/badge/github.com/juanbrotenelle/go_salamoonder.svg)](https://pkg.go.dev/github.com/juanbrotenelle/go_salamoonder)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

# go-salamoonder

Go library for working with [Salamoonder API](https://salamoonder.com).

## Installation

```bash
go get -u github.com/juanbrotenelle/go_salamoonder
```

## Usage Examples

### Get account balance

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	client, err := salamoonder.New("sr-YOUR-API-KEY", nil)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := client.Balance(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance.Wallet) // 499.81070
}
```

### Extract p.js from website

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

### Create task and get result

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	client, err := salamoonder.New("sr-YOUR-API-KEY", nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.CreateTask(context.Background(), salamoonder.KasadaOptions{
		Pjs:    "https://example.com/p.js",
		CdOnly: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	taskResult, err := client.Task(context.Background(), result.TaskId)
	if err != nil {
		log.Fatal(err)
	}

	var solution salamoonder.KasadaSolution
	if err := json.Unmarshal(taskResult.Solution, &solution); err != nil {
		log.Fatal(err)
	}

	fmt.Println(solution.UserAgent) // Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
}
```

### Polling task result

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	client, err := salamoonder.New("sr-YOUR-API-KEY", nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.CreateTask(context.Background(), salamoonder.KasadaOptions{
		Pjs:    "https://example.com/p.js",
		CdOnly: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		taskResult, err := client.Task(context.Background(), result.TaskId)
		if err != nil {
			log.Fatal(err)
		}

		if taskResult.Status == "ready" {
			var solution salamoonder.KasadaSolution
			if err := json.Unmarshal(taskResult.Solution, &solution); err != nil {
				log.Fatal(err)
			}
			fmt.Println(solution) // {xxx xxx xxx xx xxx xxx xxx 1704673650160}
			break
		}

		time.Sleep(2 * time.Second)
	}
}
```

### Get result with generics

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/juanbrotenelle/go_salamoonder"
)

func main() {
	client, err := salamoonder.New("sr-YOUR-API-KEY", nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.CreateTask(context.Background(), salamoonder.KasadaOptions{
		Pjs:    "https://example.com/p.js",
		CdOnly: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	taskResult, err := salamoonder.GetTaskResult[salamoonder.KasadaSolution](
		client,
		context.Background(),
		result.TaskId,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(taskResult.Solution.UserAgent) // Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36
}
```
