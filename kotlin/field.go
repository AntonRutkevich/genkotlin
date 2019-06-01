package kotlin

import (
	"strings"
	"unicode"

	"github.com/AntonRutkevich/astparser"
)

type fieldNameConverter func(def astparser.FieldDef) string

func defaultFieldNameConverter(def astparser.FieldDef) string {
	return lowerCaseFirst(def.FieldName)
}

func jsonTagFieldNameConverter(def astparser.FieldDef) string {
	parts := strings.Split(def.JsonName, "_")
	multiWord := len(parts) > 1
	for i := 1; multiWord && i < len(parts); i++ {
		parts[i] = upperCaseFirst(parts[i])
	}
	parts[0] = lowerCaseFirst(parts[0])
	if !multiWord {
		return parts[0]
	}

	b := strings.Builder{}
	for _, s := range parts {
		b.WriteString(s)
	}
	return b.String()
}

func convertField(def astparser.FieldDef, name fieldNameConverter) Field {
	return Field{
		Name:     name(def),
		Type:     parseType(def.FieldType),
		Doc:      strings.Join(def.Comments, ", "),
		Optional: def.Nullable,
	}
}

func lowerCaseFirst(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func upperCaseFirst(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
