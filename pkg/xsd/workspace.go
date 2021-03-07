package xsd

import (
	"fmt"
	"os"
	"path/filepath"
)

type Workspace struct {
	Cache         map[string]*Schema
	GoModulesPath string
}

func NewWorkspace(goModulesPath, xsdPath string) (*Workspace, error) {
	ws := Workspace{
		Cache:         map[string]*Schema{},
		GoModulesPath: goModulesPath,
	}

	_, err := ws.loadXsd(xsdPath)

	cbd := ws.Cache["combined.xsd"]
	fmt.Println(cbd)

	res := Schema{
		Xmlns:           cbd.Xmlns,
		XMLName:         cbd.XMLName,
		TargetNamespace: cbd.TargetNamespace,
		ModulesPath:     cbd.ModulesPath,
		filePath:        cbd.filePath,
	}

	for k, s := range ws.Cache {
		res.SimpleTypes = append(res.SimpleTypes, s.SimpleTypes...)
		res.ComplexTypes = append(res.ComplexTypes, s.ComplexTypes...)
		res.Attributes = append(res.Attributes, s.Attributes...)
		res.Elements = append(res.Elements, s.Elements...)
		res.AttributeGroups = append(res.AttributeGroups, s.AttributeGroups...)

		delete(ws.Cache, k)
	}

	ws.Cache["base"] = &res
	res.compile()

	return &ws, err
}

func (ws *Workspace) loadXsd(xsdPath string) (*Schema, error) {
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	f, err := os.Open(xsdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	schema, err := parseSchema(f)
	if err != nil {
		return nil, err
	}

	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	ws.Cache[xsdPath] = schema

	dir := filepath.Dir(xsdPath)
	for idx, _ := range schema.Imports {
		if err := schema.Imports[idx].load(ws, dir); err != nil {
			return nil, err
		}
	}

	// schema.compile()
	return schema, nil
}
