package filter_test

import (
	"fmt"
	"regexp"

	"github.com/Icenium/go-personal-data-filter/filter"
)

func ExampleNewBuilder() {
	f, err := filter.NewBuilder().
		WithMask(`¯\_(:|)_/¯`).
		Build()
	if err != nil {
		panic(err)
	}

	fmt.Println(f.RemovePersonalData("some@mail.com"))
	// Output:
	// ¯\_(:|)_/¯
}

func ExamplePersonalDataFilterBuilder_WithPersonalDataProperties() {
	f, err := filter.NewBuilder().
		WithPersonalDataProperties(`myprop`).
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

	fmt.Printf("%#v\n", f.RemovePersonalData(input))
	// Output:
	// struct { Email string; MyProp string }{Email:"not-personal", MyProp:""}
}

func ExamplePersonalDataFilterBuilder_WithAdditionalPersonalDataProperties() {
	f, err := filter.NewBuilder().
		WithAdditionalPersonalDataProperties(`myprop`).
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
	fmt.Printf("%#v\n", f.RemovePersonalData(input))
	// Output:
	// struct { Email string; MyProp string }{Email:"", MyProp:""}
}

func ExamplePersonalDataFilterBuilder_WithRegExp() {
	f, err := filter.NewBuilder().
		WithRegExp(regexp.MustCompile(`\-.*`)). // override all default regular expressions.
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
	fmt.Printf("%#v\n", f.RemovePersonalData(input))
	// Output:
	// struct { Personal string; MyProp string }{Personal:"some@mail.com", MyProp:"not"}
}

func ExamplePersonalDataFilterBuilder_WithAdditionalRegularExpressions() {
	f, err := filter.NewBuilder().
		WithAdditionalRegularExpressions(`\-.*`).
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
	fmt.Printf("%#v\n", f.RemovePersonalData(input))
	// Output:
	// struct { Personal string; MyProp string }{Personal:"", MyProp:"not"}
}

func ExamplePersonalDataFilterBuilder_WithDefaultMatchReplacer() {
	f, err := filter.NewBuilder().
		WithDefaultMatchReplacer().
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
	// Output:
	// struct { Personal string; MyProp string }{Personal:"3d13579f08e876d2d2d94da15ea657fb39795dd2a59e3378c9e58c4f4b0d053b", MyProp:"not-personal"}
}

func ExamplePersonalDataFilterBuilder_WithMatchReplacer() {
	f, err := filter.NewBuilder().
		WithMatchReplacer(func(input string) string { return input + "-replaced" }).
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
	// Output:
	// struct { Personal string; MyProp string }{Personal:"email@mail.com-replaced", MyProp:"not-personal"}
}
