package main

import (
	"fmt"
	"os"

	"github.com/clbanning/mxj"
	"github.com/kr/pretty"

	. "github.com/dave/jennifer/jen"
)

var structNames map[string]map[string]string = map[string]map[string]string{
	"//OPNsense/Firewall/Filter": map[string]string{
		"rule":      "FilterRule",
		"snatrules": "FilterSnatRule",
	},
}

func main() {

	// Open our xmlFile
	xmlFile, err := os.Open("models/Filter.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.

	type Map map[string]interface{}

	mv, err := mxj.NewMapXmlReader(xmlFile)
	if err != nil {
		fmt.Println(err)
	}

	pretty.Println(mv)

	model := mv["model"].(map[string]interface{})
	statements := createStructsFromModel(model)

	f := NewFile("main")
	f.Add(statements[1])
	fmt.Printf("%#v", f)
}

func createStructsFromModel(model map[string]interface{}) []*Statement {
	items := model["items"].(map[string]interface{})
	mount := model["mount"].(string)

	statements := make([]*Statement, 0)

	for key, value := range items {
		item := value.(map[string]interface{})
		name := structNames[mount][key]

		for itemKey, itemValue := range item {
			itemContent := itemValue.(map[string]interface{})
			createStructFromItem(name+itemKey, itemContent, statements, nil)
		}
	}

	return statements
}

func createStructFromItem(name string, field map[string]interface{}, structs []*Statement, parent *Statement) {
	fmt.Println(name, field)
	t := field["-type"].(string)

	switch t {
	case "IntegerField":
		jso := createJSONTag(name, field["Required"])

		parent.Add(
			Id(name).Int().Tag(map[string]string{"json": jso}),
		)
	case "BooleanField":
		jso := createJSONTag(name, field["Required"])

		parent.Add(
			Id(name).Bool().Tag(map[string]string{"json": jso}),
		)
	case "OptionField":
		jso := createJSONTag(name, field["Required"])

		parent.Add(
			Id(name).String().Tag(map[string]string{"json": jso}),
		)
	case ".\\FilterRuleField":
		// struc.Add(
		// 	Id(fieldName).Id("i").Tag(map[string]string{"json": jso}),
		// )
		struc := Type().Id(name).Struct()

		createStructFromItem(name, field, structs, struc)
		fmt.Printf("Struct: %#v\n", struc)

		fmt.Printf("Structs: %#v\n", structs)

		structs = append(structs, struc)

	default:
		fmt.Printf("Unsupported field type: %s\n", t)
	}

}

func createJSONTag(name string, required interface{}) string {
	if required != nil && required.(string) == "Y" {
		return name
	} else {
		return name + ",omitempty"
	}

}
