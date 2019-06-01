package kotlin

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/AntonRutkevich/astparser"
)

type Generator struct {
	fieldNameConverter fieldNameConverter
	outPackage         string
}

type Config struct {
	UseJsonTagNames bool
	OutPackage      string
}

func NewGenerator(cfg Config) *Generator {
	generator := &Generator{
		fieldNameConverter: defaultFieldNameConverter,
		outPackage:         cfg.OutPackage,
	}
	if cfg.UseJsonTagNames {
		generator.fieldNameConverter = jsonTagFieldNameConverter
	}
	return generator
}

func (g *Generator) Generate(sources map[string]astparser.ParsedFile) map[string][]byte {
	temp := template.New("tmpl").Funcs(isLastElemFn)
	t, err := temp.Parse(tmpl)
	if err != nil {
		panic(fmt.Sprintf("failed to parse template: %v", err))
	}

	result := make(map[string][]byte, len(sources))
	for name, file := range sources {
		outPackage := file.Package
		if g.outPackage != "" {
			outPackage = g.outPackage
		}
		f := File{
			Classes: make([]Class, 0, len(file.Structs)),
			Package: outPackage,
		}
	Loop:
		for _, structDef := range file.Structs {
			class := Class{
				Name:   structDef.Name,
				Fields: make([]Field, 0, len(structDef.Fields)),
			}
			for _, fieldDef := range structDef.Fields {
				if fieldDef.JsonName == "" {
					continue Loop
				}
				field := convertField(fieldDef, g.fieldNameConverter)
				class.addDataClass(fieldDef, &field)
				class.Fields = append(class.Fields, field)
			}
			f.Classes = append(f.Classes, class)
		}
		data := bytes.Buffer{}
		if err := t.Execute(&data, f); err != nil {
			log.Fatalf("failed to execute template for file %s: %v", name, err)
		}
		result[name] = data.Bytes()
	}
	return result
}

func (class *Class) addDataClass(fieldDef astparser.FieldDef, field *Field) {
	simpleTypeDef, ok := fieldDef.FieldType.(astparser.TypeSimple)
	if !ok {
		return
	}

	field.Type = fieldDef.FieldName

	simpleType := parseSimpleType(simpleTypeDef)
	class.DataClasses = append(class.DataClasses, DataClass{Name: fieldDef.FieldName, Type: simpleType})
}

func parseType(t astparser.Type) string {
	switch v := t.(type) {
	case astparser.TypeSimple:
		return parseSimpleType(v)
	case astparser.TypeArray:
		return fmt.Sprintf("List<%s>", parseType(v.InnerType))
	case astparser.TypeMap:
		return fmt.Sprintf("Map<%s,%s>", parseType(v.KeyType), parseType(v.ValueType))
	case astparser.TypePointer:
		return parseType(v.InnerType)
	case astparser.TypeCustom:
		//TODO handle dependency
		return v.Name
	default:
		panic(fmt.Sprintf("unknown type %+[1]v: %[1]T", t))
	}
}

func parseSimpleType(simple astparser.TypeSimple) string {
	switch simple.Name {
	case "string":
		return "String"
	case "int8":
		return "Byte"
	case "int16":
		return "Short"
	case "int", "int32":
		return "Int"
	case "int64":
		return "Long"
	case "float32":
		return "Float"
	case "float64":
		return "Double"
	case "bool":
		return "Boolean"
	default:
		panic(fmt.Sprintf("unknown go simple type %s", simple.Name))
	}
}
