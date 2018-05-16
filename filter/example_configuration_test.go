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
