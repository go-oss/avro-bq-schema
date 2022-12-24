package schema

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/actgardner/gogen-avro/v10/parser"
	"github.com/actgardner/gogen-avro/v10/resolver"
	"github.com/actgardner/gogen-avro/v10/schema"
)

const (
	logicalTypeAttr            = "logicalType"
	logicalTypeTimestampMicros = "timestamp-micros"
	logicalTypeDate            = "date"
	logicalTypeTimeMicros      = "time-micros"
)

func Convert(avroSchema []byte) (bigquery.Schema, error) {
	namespace := parser.NewNamespace(false)
	_, err := namespace.TypeForSchema(avroSchema)
	if err != nil {
		return nil, fmt.Errorf("namespace.TypeForSchema: %w", err)
	}

	for _, def := range namespace.Roots {
		if err := resolver.ResolveDefinition(def, namespace.Definitions); err != nil {
			return nil, fmt.Errorf("resolver.ResolveDefinition: %w", err)
		}
	}

	bqSchema := make([]*bigquery.FieldSchema, 0)
	for _, def := range namespace.Roots {
		if r, ok := def.(*schema.RecordDefinition); ok {
			for _, field := range r.Fields() {
				s := &bigquery.FieldSchema{
					Name:        field.Name(),
					Description: field.Doc(),
					Required:    true,
				}
				typeMapping(s, field.Type())
				bqSchema = append(bqSchema, s)
			}
		}
	}

	return bqSchema, nil
}

func typeMapping(s *bigquery.FieldSchema, avroType schema.AvroType) {
	switch t := avroType.(type) {
	case *schema.BoolField:
		s.Type = bigquery.BooleanFieldType
		return

	case *schema.IntField:
		s.Type = bigquery.IntegerFieldType
		if logicalType, ok := t.Attribute(logicalTypeAttr).(string); ok {
			if logicalType == logicalTypeDate {
				s.Type = bigquery.DateFieldType
			}
		}
		return

	case *schema.LongField:
		s.Type = bigquery.IntegerFieldType
		if logicalType, ok := t.Attribute(logicalTypeAttr).(string); ok {
			switch logicalType {
			case logicalTypeTimestampMicros:
				s.Type = bigquery.TimestampFieldType
			case logicalTypeTimeMicros:
				s.Type = bigquery.TimeFieldType
			}
		}
		return

	case *schema.FloatField, *schema.DoubleField:
		s.Type = bigquery.FloatFieldType
		return

	case *schema.BytesField:
		s.Type = bigquery.BytesFieldType
		return

	case *schema.StringField:
		s.Type = bigquery.StringFieldType
		return

	case *schema.ArrayField:
		s.Repeated = true
		typeMapping(s, t.ItemType())
		return

	case *schema.MapField:
		s.Repeated = true
		typeMapping(s, t.ItemType())
		return

	case *schema.UnionField:
		for _, typ := range t.AvroTypes() {
			if _, ok := typ.(*schema.NullField); ok {
				s.Required = false
				continue
			}
			typeMapping(s, typ)
		}
		return
	}
}
