package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// EqualSlice --
func EqualSlice(a interface{}, b interface{}) bool {
	return strings.EqualFold(fmt.Sprintf("%b", a), fmt.Sprintf("%b", b))
}

// StructToMap --
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}

// StructToPbStruct --
func StructToPbStruct(obj interface{}) (newStruct *structpb.Struct, err error) {
	newMap, err := StructToMap(obj)
	if err != nil {
		return
	}
	newStruct, err = structpb.NewStruct(newMap)
	return
}

// InterfaceToStruct Convert an interface to a specify struct
func InterfaceToStruct(source interface{}, result interface{}) error {
	sourceJson, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(sourceJson, result)
}

func StructToString(v interface{}) string {
	ops := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	if protoMessage, ok := v.(proto.Message); ok {
		if b, err := ops.Marshal(protoMessage); err == nil {
			return string(b)
		}
	}
	if b, err := json.Marshal(v); err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

type jsonObjectMarshaler struct {
	obj any
}

func (j *jsonObjectMarshaler) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(j.obj)
	if err != nil {
		return nil, fmt.Errorf("json marshaling failed: %w", err)
	}
	return bytes, nil
}

func ZapJson(key string, obj any) zap.Field {
	return zap.Reflect(key, &jsonObjectMarshaler{obj: obj})
}
