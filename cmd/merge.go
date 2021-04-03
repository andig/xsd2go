package cmd

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/andig/xsd2go/pkg/xsd"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli"
)

func init() {
	merge := cli.Command{
		Name:  "merge",
		Usage: "merge XSD files into single file",
		Before: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.NewExitError("Need at least 1 file", 1)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			// xsdFile, goModule, outputDir := c.Args()[0], c.Args()[1], c.Args()[2]
			err := mergeFiles(c.Args())
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}

	app.Commands = append(app.Commands, merge)
}

var summary *xsd.Schema

func merge(schema *xsd.Schema) {
	if summary == nil {
		summary = schema
		return
	}

NAMESPACES:
	for _, ns := range schema.Xmlns {
		for _, summaryns := range summary.Xmlns {
			if summaryns.Prefix == ns.Prefix {
				continue NAMESPACES
			}
		}

		summary.Xmlns = append(summary.Xmlns, ns)
	}

	summary.Imports = append(summary.Imports, schema.Imports...)
	summary.AttributeGroups = append(summary.AttributeGroups, schema.AttributeGroups...)
	summary.Attributes = append(summary.Attributes, schema.Attributes...)
	summary.ComplexTypes = append(summary.ComplexTypes, schema.ComplexTypes...)
	summary.Elements = append(summary.Elements, schema.Elements...)
	summary.SimpleTypes = append(summary.SimpleTypes, schema.SimpleTypes...)
}

func check() {

	for _, e := range summary.AttributeGroups {
		if len(funk.Filter(summary.AttributeGroups, func(s xsd.AttributeGroup) bool {
			return e.Name == s.Name
		}).([]xsd.AttributeGroup)) > 1 {
			log.Fatal("duplicate AttributeGroup:", e.Name)
		}
	}

	for _, e := range summary.Attributes {
		if len(funk.Filter(summary.Attributes, func(s xsd.Attribute) bool {
			return e.Name == s.Name
		}).([]xsd.Attribute)) > 1 {
			log.Fatal("duplicate Attribute:", e.Name)
		}
	}

	for _, e := range summary.ComplexTypes {
		if len(funk.Filter(summary.ComplexTypes, func(s xsd.ComplexType) bool {
			return e.Name == s.Name
		}).([]xsd.ComplexType)) > 1 {
			log.Fatal("duplicate ComplexType:", e.Name)
		}
	}

	for _, e := range summary.Elements {
		if len(funk.Filter(summary.Elements, func(s xsd.Element) bool {
			return e.Name == s.Name
		}).([]xsd.Element)) > 1 {
			log.Fatal("duplicate Element:", e.Name)
		}
	}

	for _, e := range summary.SimpleTypes {
		if len(funk.Filter(summary.SimpleTypes, func(s xsd.SimpleType) bool {
			return e.Name == s.Name
		}).([]xsd.SimpleType)) > 1 {
			log.Fatal("duplicate SimpleType:", e.Name)
		}
	}
}

func namespaces() string {
	ns := funk.Map(summary.Xmlns, func(ns xsd.XmlnsAttr) string {
		return fmt.Sprintf("xmls:%s=%s", ns.Prefix, ns.Uri)
	}).([]string)

	return strings.Join(ns, " ")
}

func mergeFiles(files []string) error {
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

		// fmt.Println(schema)
		merge(schema)
	}

	check()

	io.WriteString(os.Stdout, fmt.Sprintf("<!-- %s -->\n", namespaces()))
	io.WriteString(os.Stdout, xml.Header)
	e := xml.NewEncoder(os.Stdout)
	e.Indent("", "    ")
	if err := e.Encode(&summary); err != nil {
		log.Fatal(err)
	}

	return nil
}
