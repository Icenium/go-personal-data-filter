# go-personal-data-filter
[![GoDoc](https://godoc.org/github.com/Icenium/go-personal-data-filter/filter?status.svg)](https://godoc.org/github.com/Icenium/go-personal-data-filter/filter)

## Contents
- [Installation](#installation)
- [What will be filtered](#what-will-be-filtered)
- [Example](#example)
- [Configuration](#configuration)

## Installation:
```shell
dep ensure -add github.com/Icenium/go-personal-data-filter/filter
```

## What will be filtered:
- Structs
	- recursive
	- properties with special names like `password`, `email` etc. will be filtered even if they don't contain personal data ([list of personal data properties](./filter/builder.go#L23))
	- properties with tag \``pdfilter:"nofilter"`\` will not be filtered
- Maps
	- recursive
	- the values with keys like `password`, `email` etc. will be filtered even if they don't contain personal data ([list of personal data properties](./filter/builder.go#L23))
- Arrays
	- each item will be checked
- Slices
	- each item will be checked
- Strings
	- Emails
	- GUIDs
	- IP v4
	- IP v6

## Example:
```Go
package main

import (
	"fmt"

	"github.com/Icenium/go-personal-data-filter/filter"
)

type someData struct {
	FilterMe     string
	DontFilterMe string
	NextLevel    map[string]string
	Email        string
	Items        []string
	ID           string `pdfilter:"nofilter"`
}

func main() {
	f, err := filter.NewBuilder().
		SetMask("*****").
		Build()
	if err != nil {
		panic(err)
	}

	input := someData{
		FilterMe:     "some@mail.com", // will be filtered and replaced with *****
		DontFilterMe: "some-data",
		NextLevel: map[string]string{
			"filterMe":     "1fec999a-7e81-4bce-8b32-1b6ddd144f1b", // will be filtered and replaced with *****
			"dontFilterMe": "some-data",
			"email":        "not-personal", // will be filtered and replaced with *****
		},
		Email: "some@mail.bg",                                                     // will be filtered and replaced with *****
		Items: []string{"some-data", "some@mail.bg", "some-data", "some@mail.bg"}, // will be filtered and the result will be ["some-data", "*****", "some-data", "*****"]
		ID:    "1fec999a-7e81-4bce-8b32-1b6ddd144f1b",                             // this field will not be filtered because of the nofilter setting in the pdfilter tag
	}
	res := f.RemovePersonalData(input)
	fmt.Println(res)
}
```

## Configuration:
- Mask:
```Go
package main

import (
	"fmt"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		WithMask(`¯\_(:|)_/¯`).
		Build()
	if err != nil {
		panic(err)
	}

	fmt.Println(f.RemovePersonalData("some@mail.com"))
}
```
- Personal data properties:
```Go
package main

import (
	"fmt"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		SetPersonalDataProperties(`myprop`). // override all default personal data properties.
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Email  string
		MyProp string
	}{
		Email:  "not-personal", // will not be filtered
		MyProp: "not-personal", // will be filtered
	}

	fmt.Println(f.RemovePersonalData(input))
}
```
```Go
package main

import (
	"fmt"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		AddPersonalDataProperties(`myprop`).
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Email  string
		MyProp string
	}{
		Email:  "not-personal", // will be filtered
		MyProp: "not-personal", // will be filtered
	}

	fmt.Println(f.RemovePersonalData(input))
}
```
- Regular expressions:
```Go
package main

import (
	"fmt"
	"regexp"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		SetRegExp(regexp.MustCompile(`\-.*`)). // override all default regular expressions.
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Personal string
		MyProp   string
	}{
		Personal: "some@mail.com", // will not be filtered
		MyProp:   "not-personal",  // will be filtered and te result will be "not"
	}

	fmt.Println(f.RemovePersonalData(input))
}
```
```Go
package main

import (
	"fmt"
	"regexp"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		AddRegularExpressions(`\-.*`).
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Personal string
		MyProp   string
	}{
		Personal: "some@mail.com", // will be filtered and te result will be ""
		MyProp:   "not-personal",  // will be filtered and te result will be "not"
	}

	fmt.Println(f.RemovePersonalData(input))
}
```
- Match filter function:
```Go
package main

import (
	"fmt"
	"regexp"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		UseDefaultMatchFilterFunc(). // use the default match filter function to the default one - sha256 sum.
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Personal string
		MyProp   string
	}{
		Personal: "email@mail.com", // will be replaced and the result will be sha256 hash
		MyProp:   "not-personal",   // will not be replaced
	}
	fmt.Printf("%#v\n", f.RemovePersonalData(input))
}
```
```Go
package main

import (
	"fmt"
	"regexp"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func main() {
	f, err := filter.NewBuilder().
		SetMatchFilterFunc(func(input string) string { return input + "-replaced" }).
		Build()
	if err != nil {
		panic(err)
	}

	input := struct {
		Personal string
		MyProp   string
	}{
		Personal: "email@mail.com", // will be replaced and the result will be email@mail.com-replaced
		MyProp:   "not-personal",   // will not be replaced
	}
	fmt.Printf("%#v\n", f.RemovePersonalData(input))
}
```

