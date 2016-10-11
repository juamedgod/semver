[![Build Status](https://travis-ci.org/juamedgod/semver.svg?branch=master)](https://travis-ci.org/juamedgod/semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/juamedgod/semver)](https://goreportcard.com/report/github.com/juamedgod/semver)
[![GoDoc](https://godoc.org/github.com/juamedgod/semver?status.svg)](https://godoc.org/github.com/juamedgod/semver)


Go Semantic Versioning Parser
=============================

## Basic Usage

```
package main

import (
  "fmt"
  "os"

  "github.com/juamedgod/semver"
)

func main() {
  if result, err := semver.Satisfies("1.3.4", "1.2.1 - 1.4.0"); err != nil {
    fmt.Fprintf(os.Stderr, "Error evalutating semver expression: %v", err)
    os.Exit(1)
  } else {
    fmt.Println(result)
  }
}
```

## Versions

A `Version` object defines a [semver](http://semver.org/) version.

Versions can be created:

```go
// err1 should be nil, as the version is properly formated
v1, err1 := ParseVersion("1.3.5")

// err2 will contain an error
v2, err := ParseVersion("asdf")

// If you are confident your version string will validate
// you can use the "Must" version, which will panic on error but 
// will make the function more convenient

v3 := MustParseVersion("5.4")

// You can manually define them
v5 := &Version{Major: 6, Minor: 0, Patch: 4}
```

And compared

```go
// true
v3.Greater(v2)
 
// false
v2.Greater(v3) 

// true
v3.GreaterOrEqual(v3)

// true
v3.Less(v5)

// true
v5.LessOrEqual(v3)

// true
v5.Equal(v5)
```

## Ranges

A `Range` defines a range of versions. The syntax to create them is similar to the one used in [Versions](#versions)

```go

// err1 should be nil, as the range is properly formated
r1, err1 := ParseRange(">1.3.5")

// err2 will contain an error
r2, err := ParseRange("asdf")

// If you are confident your range string will validate
// you can use the "Must" version, which will panic on error but 
// will make the function more convenient

r3 := MustParseRange("^5.4")
```

And can be used to check if a certain `Version` fulfills it:

```go
v1 := MustParseVersion("1.0.4")
v2 := MustParseVersion("6.4.3")

r1 := MustParseRange(">1.3.5")
r2 := MustParseRange("< 4.0")

// false
r1.Contains(v1)

// true
r2.Contains(v1)

// true
r1.Contains(v2)

// false
r2.Contains(v2)
```

## Expressions

An `Expression` is a combination of ranges. Ranges separated by spaces act as an "AND" operation, and those serparated by `||` act as an "OR".

Creating expressions follow the usual procedure:

```go
// err1 should be nil, as the range is properly formated
e1, err1 := ParseExpr(">2.3.1 <3.4.1")

// err2 will contain an error
e2, err2 :=  ParseExpr("fooo")

// If you are confident your range string will validate
// you can use the "Must" version, which will panic on error but 
// will make the function more convenient
e3 := MustParseExpr("<1.5 || >5.4")
```

And can then be used to check if a version matches:

```go
v1 := MustParseVersion("3.1.0")
v2 := MustParseVersion("1.1.0")
v3 := MustParseVersion("6.1.0")

e1 := MustParseExpr(">2.3.1 <3.4.1")

// true
e1.Matches(v1)
// false
e1.Matches(v2)
// false
e1.Matches(v3)

e2 := MustParseExpr("<1.5 || >5.4")
// false
e2.Matches(v1)
// true
e2.Matches(v2)
// true
e2.Matches(v3)
```

