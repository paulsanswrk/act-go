package bngx

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"testing"
)

func TestBuildUrl(t *testing.T) {
	params := map[string]interface{}{"b": 1, "a": "qqq"}
	res := build_and_sign_url(params)
	fmt.Println(res)
}

func TestMapStructure(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n\n", result)

	input = map[string]interface{}{}
	err = mapstructure.Decode(result, &input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", input)
}
