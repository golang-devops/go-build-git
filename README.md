# go-build-git
Simple wrapper around `go build` to allow getting the git hash and injecting it into a variable


## Quick Start

### Install this package binary with

```
go get -u github.com/golang-devops/go-build-git
```

### Define your `main.go` file:

```
package main

var (
    GitHash = "NO GIT HASH"
)

func main() {
    fmt.Println(fmt.Sprintf("Running git hash '%s'", GitHash))
}
```

### Now inside the main file directory run:

```
go-build-git -out "./tmpbinary" -injectvar "main.GitHash"
./tmpbinary
```