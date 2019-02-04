# autocom

This is a simple auto-completion engine/server that is fast and efficient. It uses Ternary Search Tree (TST) as its main underlying data-structure which is storage efficient and allows for fast lookup.

## Getting Started
- [Install](#install)
- [Examples](#examples)
- [Engine Options](#engine-options)
- [Usage](#usage)
	- [Set](#set)
	- [Remove](#remove)
	- [TopN](#topn)
- [Run as a Server](#run-as-a-server)
- [Server endpoints](#server-endpoints)
- [Benchmarks](#benchmarks)


### Install
```bash
go get github.com/zeyadyasser/autocom
```

---
### Examples
Simple Example
```go
import (
	"fmt"
	"github.com/zeyadyasser/autocom/engine/skip"
)

func main() {
	opts := skip.Options{
		MaxLevels: 5,
		ToLower: true,
		SkipBegin: true,
	}
	E := skip.NewSkipEngine(opts, nil) // nil to use default backend (TST)
	E.Set("Mad Max: Fury Road", "value (any interface)")
	E.Set("Furious 7", nil)
	top, _ := E.TopN("Fu", 2)
	for K, V :=  range top {
		fmt.Printf("%v --> %v\n", K, V)
	}
	// Mad Max: Fury Road --> value (any interface)
	// Furious 7 --> nil
	top, _ = E.TopN("road", 2)
	for K, V :=  range top {
		fmt.Printf("%v --> %v\n", K, V)
	}
	// Mad Max: Fury Road --> value (any interface)
	E.Remove("Mad Max: Fury Road")
	top, _ = E.TopN("Fu", 2)
	for K, V :=  range top {
		fmt.Printf("%v --> %v\n", K, V)
	}
	// Furious 7 --> nil
}
```
Checkout the interactive [movies search sample](/movies_sample/main.go)

---
### Engine Options
Here the search is **case-sensitive** but **exact prefix** matches are required.
```go
opts := skip.Options{}
```
**Case-sensitive Example :**
```go
opts := skip.Options{}
E := skip.NewSkipEngine(opts, nil) // nil to use default backend (TST)
E.Set("Mad Max: Fury Road", nil)
E.Set("Furious 7", nil)
top, _ := E.TopN("mad", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// No results
top, _ = E.TopN("Mad", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// Mad Max: Fury Road --> nil
```
**Exact Prefix Example :**
```go
opts := skip.Options{}
E := skip.NewSkipEngine(opts, nil) // nil to use default backend (TST)
E.Set("Mad Max: Fury Road", nil)
E.Set("Furious 7", nil)
top, _ := E.TopN("Fury", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// No results, because "Fury is not the exact prefix of any key"
top, _ = E.TopN("Fur", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// Furious 7 --> nil
```
---
Here the search is **case insensitive** but **exact prefix** matches are required.
```go
opts := skip.Options{
	ToLower: true,
}
```
**Case-insensitive Example :**
```go
opts := skip.Options{
	ToLower: true,
}
E := skip.NewSkipEngine(opts, nil) // nil to use default backend (TST)
E.Set("Mad Max: Fury Road", nil)
E.Set("Furious 7", nil)
top, _ := E.TopN("mad", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// Mad Max: Fury Road --> nil
```
---
Here the search is **case-insensitive** and matches keys where any of the first N(**MaxLevels**) words of the key matches the query.
```go
opts := skip.Options{
	MaxLevels: 5,	// Here N is 5
	ToLower: true,  // Case-insensitive
	SkipBegin: true,// Match by skipping up to the first N words
}
```
**Skip Begin Example :**
```go
opts := skip.Options{
	MaxLevels: 3,	// Here N is 3
	ToLower: true,  // Case-insensitive
	SkipBegin: true,// Match by skipping up to the first N words
}
E := skip.NewSkipEngine(opts, nil) // nil to use default backend (TST)
E.Set("Mad Max: Fury Road", nil)
// The following keys are stored.
// 	"mad max: fury road"
// 	"max: fury road"
// 	"fury road"
// 	"road" is not stored, because N = 3 and this is the 4th word.
E.Set("Furious 7", nil)
// The following keys are stored.
// 	"furious 7"
// 	"7"
top, _ := E.TopN("fu", 2)
for K, V :=  range top {
	fmt.Printf("%v --> %v\n", K, V)
}
// Mad Max: Fury Road --> nil (Here Fury is matched)
// Furious 7 --> nil (Here Furious is matched)
```
---
### Usage
 The Engine interface offers three main methods.
 ### Set
 ```go
 Set(string, interface{}) error
```
example:
```go
E := skip.NewSkipEngine(opts, nil)
err := E.Set("Key", "any interface")
```
 ### Remove
 ```go
 Remove(string) error
```
example:
```go
E := skip.NewSkipEngine(opts, nil)
err := E.Remove("Key")
```
 ### TopN
 Maps Top N autocomplete matching keys to their values.
 ```go
 TopN(string, int) (complete.Map, error)
```
where ```complete.Map``` is ```map[string]interface{}```
example:
```go
E := skip.NewSkipEngine(opts, nil)
result, err := E.TopN("Key", 5)
// returns top 5 matching keys
```

---
### Run as a Server
- Clone the repo & cd to root of repository.
	```bash
	git clone https://github.com/zeyadyasser/autocom.git
	cd autocom
	```
- Create a new .env file using the .env.dist as a template
	-   Fill in the missing secrets suitable for your environment
	-   Never check in the .env file, it is already included in .gitignore, and never add actual secrets in .env.dist
- (Optional) To run on your local machine, run the following command
	```bash
	go run .
	```
- (Optional) To use Docker,  run the following command
	```bash
	docker build -t autocom .
	docker run --env-file=.env --rm -p 9696:9696 autocom
	```
- (Optional) To use docker-compose,  run the following command
	```bash
	docker-compose build
	docker-compose up
	```
---
### Server endpoints
**POST /set** 
You should send data in JSON format.
Example:
```bash
curl --header "Content-Type: application/json" \
  -X POST \
  --data '{"key":"xyz","value":"xyz"}' \
  http://user:password@localhost:9696/set
```

**DELETE /remove** 
You should send the key as a query parameter.
Example:
```bash
curl -X DELETE http://user:password@localhost:9696/remove?key=xyz
```

**GET /topn** 
You should send the key & number of results as query parameters.
If ```cnt``` is not a number ```400 Bad Request``` is returned.
Returns a JSON map from keys to values.
Example:
```bash
curl -X GET http://user:password@localhost:9696/topn?key=xyz&cnt=5
```

---
### Benchmarks
This was tested with Intel Core I5 8400 with 8GB of 2400hz DDR4 RAM.

**Disclaimer :** This not an exhaustive benchmark by any means.
```bash
cd path/to/autocom/engine/skip
go test -bench=.
goos: windows
goarch: amd64
pkg: github.com/zeyadyasser/autocom/engine/skip
BenchmarkTSTSet1ShortKey-6       1000000              1121 ns/op             620 B/op         15 allocs/op
BenchmarkTSTSet1LongKey-6        2000000              1086 ns/op             618 B/op         15 allocs/op
BenchmarkTSTTopN10ShortKey-6     1000000              1835 ns/op            1420 B/op         14 allocs/op
PASS
ok      github.com/zeyadyasser/autocom/engine/skip      6.629s
```
