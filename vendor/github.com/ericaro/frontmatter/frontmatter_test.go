package frontmatter

import (
	"fmt"
)

type Example struct {
	Name     string   `yaml:"name"`
	Variants []string `yaml:"variants,flow"`
	Content  string   `fm:"content" yaml:"-"`
}

func ExampleUnmarshal() {
	data := `---
name: toto
variants: [titi, tutu]
---
Hello!`

	v := new(Example)
	err := Unmarshal(([]byte)(data), v)
	if err != nil {
		fmt.Printf("err! %s", err.Error())
	}

	if v.Name == "toto" && v.Variants[0] == "titi" && v.Variants[1] == "tutu" && v.Content == "Hello!" {
		fmt.Println("Ok")
	} else {
		fmt.Printf("%v, %v, %v, %v, %v\n", v, v.Name == "toto", v.Variants[0] == "titi", v.Variants[1] == "tutu", v.Content == "Hello!")
	}
	//Output: Ok
}

func ExampleMarshal() {

	v := Example{
		Name:     "toto",
		Variants: []string{"titi", "tutu"},
		Content:  "Hello!",
	}

	data, err := Marshal(v)
	if err != nil {
		fmt.Printf("err! %s", err.Error())
	}
	fmt.Println(string(data))
	//Output:
	//---
	//name: toto
	//variants: [titi, tutu]
	//---
	//Hello!
}
