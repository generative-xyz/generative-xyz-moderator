package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/schema"
	"github.com/valyala/bytebufferpool"
)

const (
	queryTagAtStruct = "query"
)

var (
	// decoderPoolMap helps to improve BodyParser's, QueryParser's and ReqHeaderParser's performance
	decoderPoolMap = map[string]*sync.Pool{}
)

// ParserConfig form decoder config for SetParserDecoder
type ParserConfig struct {
	IgnoreUnknownKeys bool
	SetAliasTag       string
	ParserType        []ParserType
	ZeroEmpty         bool
}

// ParserType require two element, type and converter for register.
// Use ParserType with BodyParser for parsing custom type in form data.
type ParserType struct {
	Customtype interface{}
	Converter  func(string) reflect.Value
}

func init() {
	decoderPoolMap[queryTagAtStruct] = &sync.Pool{New: func() interface{} {
		return decoderBuilder(ParserConfig{
			IgnoreUnknownKeys: true,
			ZeroEmpty:         true,
		})
	}}
}

func QueryParser(r *http.Request, out interface{}) error {
	data := make(map[string][]string)
	var err error
	urlValues := r.URL.Query()
	for k := range urlValues {
		for _, v := range urlValues[k] {
			if strings.Contains(k, "[") {
				k, err = parseParamSquareBrackets(k)
			}
			if strings.Contains(v, ",") && equalFieldType(out, reflect.Slice, k) {
				values := strings.Split(v, ",")
				for i := 0; i < len(values); i++ {
					data[k] = append(data[k], values[i])
				}
			} else {
				data[k] = append(data[k], v)
			}
		}
	}
	if err != nil {
		return err
	}

	return parseToStruct(queryTagAtStruct, out, data)
}

func decoderBuilder(parserConfig ParserConfig) interface{} {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(parserConfig.IgnoreUnknownKeys)
	if parserConfig.SetAliasTag != "" {
		decoder.SetAliasTag(parserConfig.SetAliasTag)
	}
	for _, v := range parserConfig.ParserType {
		decoder.RegisterConverter(reflect.ValueOf(v.Customtype).Interface(), v.Converter)
	}
	decoder.ZeroEmpty(parserConfig.ZeroEmpty)
	return decoder
}

func equalFieldType(out interface{}, kind reflect.Kind, key string) bool {
	// Get type of interface
	outTyp := reflect.TypeOf(out).Elem()
	key = strings.ToLower(key)
	// Must be a struct to match a field
	if outTyp.Kind() != reflect.Struct {
		return false
	}
	// Copy interface to an value to be used
	outVal := reflect.ValueOf(out).Elem()
	// Loop over each field
	for i := 0; i < outTyp.NumField(); i++ {
		// Get field value data
		structField := outVal.Field(i)
		// Can this field be changed?
		if !structField.CanSet() {
			continue
		}
		// Get field key data
		typeField := outTyp.Field(i)
		// Get type of field key
		structFieldKind := structField.Kind()
		// Does the field type equals input?
		if structFieldKind != kind {
			continue
		}
		// Get tag from field if exist
		inputFieldName := typeField.Tag.Get(queryTagAtStruct)
		if inputFieldName == "" {
			inputFieldName = typeField.Name
		} else {
			inputFieldName = strings.Split(inputFieldName, ",")[0]
		}
		if strings.ToLower(inputFieldName) == key {
			return true
		}
	}
	return false
}

func parseParamSquareBrackets(k string) (string, error) {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	kbytes := []byte(k)

	for i, b := range kbytes {
		if b == '[' && kbytes[i+1] != ']' {
			if err := bb.WriteByte('.'); err != nil {
				return "", fmt.Errorf("failed to write: %w", err)
			}
		}

		if b == '[' || b == ']' {
			continue
		}

		if err := bb.WriteByte(b); err != nil {
			return "", fmt.Errorf("failed to write: %w", err)
		}
	}

	return bb.String(), nil
}

func parseToStruct(aliasTag string, out interface{}, data map[string][]string) error {
	// Get decoder from pool
	schemaDecoder, ok := decoderPoolMap[aliasTag].Get().(*schema.Decoder)
	if !ok {
		panic(fmt.Errorf("failed to type-assert to *schema.Decoder"))
	}
	defer decoderPoolMap[aliasTag].Put(schemaDecoder)

	// Set alias tag
	schemaDecoder.SetAliasTag(aliasTag)

	if err := schemaDecoder.Decode(out, data); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}
