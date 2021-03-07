package xsd

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Union struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema union"`
	// ElementList []Element `xml:"element"`
	// Choices     []Choice  `xml:"choice"`
	MemberTypes reference `xml:"memberTypes,attr"`
	allElements []Element `xml:"-"`
}

func (s *Union) Elements() []Element {
	return s.allElements
}

func (s *Union) compile(sch *Schema, parentElement *Element) {
	types := strings.Split(string(s.MemberTypes), " ")

	for _, typeName := range types {
		var el Element
		el.typ = sch.findReferencedType(reference(typeName))
		if el.typ == nil {
			panic("Cannot resolve type reference: " + string(el.Type))
		}

		el.compile(sch, parentElement)

		fmt.Println(el)
	}

	// s.allElements = s.ElementList
}
