package lucene

import (
	"encoding/json"
	"strings"
)

const (
	DataTypeString  = "string"
	DataTypeUuid    = "uuid"
	DataTypeInteger = "integer"
	DataTypeDate    = "date"
	DataTypeText    = "text"
	DataTypeBoolean = "boolean"

	OptionPattern       = "pattern"
	OptionCaseSensitive = "case_sensitive"
)

type IndexSchema []IndexSchemaField

type IndexSchemaField struct {
	Name     string
	DataType string
	Options  FieldOptions
}

type FieldOptions map[string]interface{}

func (s *IndexSchema) WithFieldOptions(name string, dataType string, options FieldOptions) *IndexSchema {
	field := IndexSchemaField{
		Name:     name,
		DataType: dataType,
		Options:  options,
	}
	*s = append(*s, field)
	return s
}

func (s *IndexSchema) WithField(name string, dataType string) *IndexSchema {
	return s.WithFieldOptions(name, dataType, FieldOptions{})
}

func (s *IndexSchema) String() string {
	sb := new(strings.Builder)
	sb.WriteString(`{"fields":{`)
	for i, field := range *s {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('"')
		sb.WriteString(field.Name)
		sb.WriteString(`":{"type":"`)
		sb.WriteString(field.DataType)
		sb.WriteRune('"')
		for optionKey, optionValue := range field.Options {
			sb.WriteString(`,"`)
			sb.WriteString(optionKey)
			sb.WriteString(`":`)
			bytes, _ := json.Marshal(optionValue)
			sb.Write(bytes)
		}
		sb.WriteRune('}')
	}
	sb.WriteString("}}")

	return sb.String()
}
