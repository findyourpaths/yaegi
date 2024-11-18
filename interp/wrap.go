package interp

import (
	"fmt"
	"reflect"
)

// WrapReflectValueSlice takes a slice of reflect.Value and returns a single reflect.Value
// that represents the slice. This is useful when working with reflection and you need
// to handle a slice of reflect.Value as a single value.
func WrapReflectValueSlice(values []reflect.Value) reflect.Value {
	// Create a new slice with the same length as the input
	sliceType := reflect.SliceOf(reflect.TypeOf((*reflect.Value)(nil)).Elem())
	slice := reflect.MakeSlice(sliceType, len(values), len(values))

	// Copy each value into the new slice
	for i, v := range values {
		slice.Index(i).Set(reflect.ValueOf(v))
	}

	return slice
}

// IsReflectValueSlice checks if the given reflect.Value represents a []reflect.Value
func IsReflectValueSlice(v reflect.Value) bool {
	// Check if it's a slice first
	if v.Kind() != reflect.Slice {
		return false
	}

	// Get the type of the slice elements
	elemType := v.Type().Elem()

	// Compare with reflect.Value's type
	return elemType == reflect.TypeOf((*reflect.Value)(nil)).Elem()
}

// UnwrapReflectValueSlice takes a reflect.Value that contains a slice of reflect.Value
// and returns it as a []reflect.Value. This is the inverse operation of WrapValueSlice.
func UnwrapReflectValueSlice(wrapped reflect.Value) ([]reflect.Value, error) {
	if wrapped.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %v", wrapped.Kind())
	}

	length := wrapped.Len()
	result := make([]reflect.Value, length)

	for i := 0; i < length; i++ {
		elem := wrapped.Index(i).Interface()
		val, ok := elem.(reflect.Value)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a reflect.Value", i)
		}
		result[i] = val
	}

	return result, nil
}
