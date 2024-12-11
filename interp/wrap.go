package interp

import (
	"fmt"
	"reflect"
)

type Variable struct {
	Name string
}

var VariableType = reflect.TypeOf(Variable{})

var WrapTopValues = true

// var WrapTopValues = false

// var doDebug = true

// var debugFn = fmt.Printf

var doDebug = false

var debugFn = fmt.Sprintf

func genValueTop(n *node) func(*frame) reflect.Value {
	debugFn("interp.genValueTop(n.kind: %s), WrapTopValues: %t\n", kinds[n.kind], WrapTopValues)
	if doDebug {
		printNodeTree(n, 0)
	}

	if !WrapTopValues {
		return genValue(n)
	}

	switch n.kind {
	case blockStmt:
		debugFn("in interp.genValueTop(), for blockStmt, delegating to last child node\n")
		return genValueTop(n.child[len(n.child)-1])

	// case callExpr:
	// 	fmt.Printf("in interp.genValueTop(), for callExpr, fetching function return values\n")
	// 	return func(f *frame) reflect.Value {
	// 		if doDebug {
	// 			printFrameData(f.data)
	// 		}

	// 		// // Get the function to execute
	// 		// fn := n.exec
	// 		// // Prepare a place to store the return values
	// 		// prevLen := len(f.data)
	// 		// // Execute the function
	// 		// fn(f)
	// 		// // Get the new values added to the frame's data
	// 		// newData := f.data[prevLen:]
	// 		// // Wrap the results in a slice of reflect.Value
	// 		// return WrapReflectValueSlice(newData)
	// 		return reflect.ValueOf(nil)
	// 	}

	case defineStmt:
		debugFn("in interp.genValueTop(), for defineStmt, collecting child identExpr nodes\n")
		return func(f *frame) reflect.Value { return reflect.ValueOf(nodeChildIdents(n)) }

	case defineXStmt:
		debugFn("in interp.genValueTop(), for defineXStmt, collecting child identExpr nodes\n")
		return func(f *frame) reflect.Value { return reflect.ValueOf(nodeChildIdents(n)) }

	case exprStmt:
		debugFn("in interp.genValueTop(), for exprStmt, checking\n")
		// Return empty slice if this evaluates to nothing (e.g. calling func () {})
		if n.findex == -1 {
			idents := []reflect.Value{}
			return func(f *frame) reflect.Value { return reflect.ValueOf(idents) }
		}
		debugFn("in interp.genValueTop(), for exprStmt, delegating to genValue\n")
		return genValue(n)

	case fileStmt:
		debugFn("in interp.genValueTop(), for fileStmt, delegating to last child node\n")
		return genValueTop(n.child[len(n.child)-1])

	case funcDecl:
		debugFn("in interp.genValueTop(), for funcDecl, collecting child identExpr nodes\n")
		idents := nodeChildIdents(n)
		if len(idents) != 1 {
			fmt.Printf("warn: expecting one identExpr, got %d\n", len(idents))
		}
		if idents[0].String() == "main" {
			return genValueTop(n.child[len(n.child)-1])
		}
		return func(f *frame) reflect.Value { return reflect.ValueOf(idents) }

	default:
		return genValue(n)
	}
}

func nodeChildIdents(n *node) []reflect.Value {
	rs := []reflect.Value{}
	for _, c := range n.child {
		if c.kind == identExpr {
			rs = append(rs, reflect.ValueOf(Variable{Name: c.ident}))
		}
	}
	return rs
}

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
	debugFn("UnwrapReflectValueSlice(wrapped.Kind(): %q)\n", wrapped.Kind().String())
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

	debugFn("in UnwrapReflectValueSlice(), returning (%d)\n", len(result))
	return result, nil
}
