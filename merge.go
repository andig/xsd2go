package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/gocomply/xsd2go/pkg/xsd"
)

var summary *xsd.Schema

func merge(schema *xsd.Schema) {
	if summary == nil {
		summary = schema
		return
	}

	summary.AttributeGroups = append(summary.AttributeGroups, schema.AttributeGroups...)
	summary.Attributes = append(summary.Attributes, schema.Attributes...)
	summary.ComplexTypes = append(summary.ComplexTypes, schema.ComplexTypes...)
	summary.Elements = append(summary.Elements, schema.Elements...)
	summary.Imports = append(summary.Imports, schema.Imports...)
	summary.SimpleTypes = append(summary.SimpleTypes, schema.SimpleTypes...)
}

func main() {
	files := os.Args[1:]
	for _, f := range files {
		// fmt.Println(f)
		r, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		schema, err := xsd.ParseSchema(r)
		if err != nil {
			log.Fatal(err)
		}

		// schema := xsd.Schema{}
		// data, err := os.ReadFile(f)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if err := xml.Unmarshal(data, &schema); err != nil {
		// 	log.Fatal(err)
		// }
		// merge(&schema)

		// fmt.Println(schema)
		merge(schema)
	}

	// fmt.Println(summary)
	io.WriteString(os.Stdout, xml.Header)
	e := xml.NewEncoder(os.Stdout)
	e.Indent("", "    ")
	if err := e.Encode(&summary); err != nil {
		log.Fatal(err)
	}

	// schema := xsd.Schema{}
	// data, err := os.ReadFile(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := xml.Unmarshal(data, &schema); err != nil {
	// 	log.Fatal(err)
	// }

}
