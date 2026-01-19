package custommarshaller

import (
	"encoding/json"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// CustomMarshaller is an auxiliary struct that helps to change the behavior of the default
// MarshalJSON method for specific types allowing custom marshalling of the fields
type CustomMarshaller struct {
	any
}

func New(i any) CustomMarshaller {
	return CustomMarshaller{i}
}

// MarshalJSON iterates over all the fields of the struct allowing custom marshalling of the fields
// - [32]byte will be marshalled as a common.Hash
// - [20]byte will be marshalled as a common.Address
// - []byte will be marshalled as a common.Hash
// - If the field is a struct or a pointer to a struct, it will be marshalled recursively
// - All the rest of the fields will be marshalled as is
//
// Example:
//
//	type MyStruct struct {
//		Field1 string
//		Field2 [32]byte
//		Field3 [20]byte
//		Field4 []byte
//	}
//	myStruct := MyStruct{
//		Field1: "value",
//		Field2: [32]byte{},
//		Field3: [20]byte{},
//		Field4: []byte{},
//	}
//
//	result, _ := json.Marshal(myStruct)
//	fmt.Println(string(result))
//
// Result returned when marshalling MyStruct directly:
//
//	{
//		"Field1": "value",
//		"Field2": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
//		"Field3": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
//		"Field4": [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]
//	}
//
//	result, _ = json.Marshal(custommarshaller.New(myStruct))
//	fmt.Println(string(result))
//
// Result returned when marshalling MyStruct using CustomMarshaller:
//
//	{
//		"Field1": "value",
//		"Field2": "0x0000000000000000000000000000000000000000000000000000000000000000",
//		"Field3": "0x0000000000000000000000000000000000000000",
//		"Field4": "0x0000000000000000000000000000000000000000000000000000000000000000"
//	}
func (c CustomMarshaller) MarshalJSON() ([]byte, error) {
	result := map[string]any{}
	instanceType := reflect.TypeOf(c.any)
	instanceValue := reflect.ValueOf(c.any)

	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
		instanceValue = instanceValue.Elem()
	}

	for i := range instanceType.NumField() {
		f := instanceType.Field(i)

		if !f.IsExported() {
			continue
		}

		fieldKind := f.Type.Kind()

		v := instanceValue.Field(i)
		if fieldKind == reflect.Array {
			var fieldInterfaceValue any
			if v.CanAddr() { // check if array is addressable
				v = v.Slice(0, f.Type.Len())
				fieldInterfaceValue = v.Interface()
				if f.Type.Len() == 20 {
					result[f.Name] = common.BytesToAddress(fieldInterfaceValue.([]byte))
				} else {
					result[f.Name] = common.BytesToHash(fieldInterfaceValue.([]byte))
				}
			} else {
				result[f.Name] = v.Interface()
			}
		} else if fieldKind == reflect.Slice {
			if f.Type.Elem().Kind() == reflect.Uint8 {
				result[f.Name] = common.BytesToHash(v.Bytes())
			} else {
				result[f.Name] = v.Interface()
			}
		} else if fieldKind == reflect.Struct || fieldKind == reflect.Ptr {
			result[f.Name] = CustomMarshaller{v.Interface()}
		} else {
			result[f.Name] = v.Interface()
		}

	}

	return json.Marshal(result)
}
