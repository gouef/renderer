<img align=right width="168" src="docs/gouef_logo.png">

# gouef/mode
Mode of project

[![GoDoc](https://pkg.go.dev/badge/github.com/gouef/mode.svg)](https://pkg.go.dev/github.com/gouef/mode)
[![GitHub stars](https://img.shields.io/github/stars/gouef/mode?style=social)](https://github.com/gouef/mode/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouef/mode)](https://goreportcard.com/report/github.com/gouef/mode)
[![codecov](https://codecov.io/github/gouef/mode/branch/main/graph/badge.svg?token=YUG8EMH6Q8)](https://codecov.io/github/gouef/mode)


## Vesions
![Stable Version](https://img.shields.io/github/v/release/gouef/mode?label=Stable&labelColor=green)
![GitHub Release](https://img.shields.io/github/v/release/gouef/mode?label=RC&include_prereleases&filter=*rc*&logoSize=diago)
![GitHub Release](https://img.shields.io/github/v/release/gouef/mode?label=Beta&include_prereleases&filter=*beta*&logoSize=diago)

## Introduction

Mode of project, like Release, Debug, Testing

## Examples

### Basic
```go
package main
import "github.com/gouef/mode"

func main()  {
    m, err := mode.NewBasicMode()
	
	if err != nil {
		// do something
    }
    
    // some code
    if r, _ := m.IsRelease(); r {
        // some code
    }
}
```


### Additional modes
```go
package main
import "github.com/gouef/mode"

func main()  {
	modes := []string{"staging"}
	m, err := mode.NewMode(modes)

	if err != nil {
		// do something
	}

	// some code
	if r, _ := m.IsRelease(); r {
		m.EnableMode("staging")
	}
	
	if sm, _ := m.IsMode("staging"); sm {
		// some code
	}
}
```


