package filter_test

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

func Example() {
	f, err := filter.NewBuilder().
		WithMask("*****").
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
	fmt.Printf("%#v\n", res)
	// Output:
	// filter_test.someData{FilterMe:"*****", DontFilterMe:"some-data", NextLevel:map[string]string{"filterMe":"*****", "dontFilterMe":"some-data", "email":"*****"}, Email:"*****", Items:[]string{"some-data", "*****", "some-data", "*****"}, ID:"1fec999a-7e81-4bce-8b32-1b6ddd144f1b"}
}
