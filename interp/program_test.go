package interp

import (
	"reflect"
	"testing"
)

func TestREPLMultiValue(t *testing.T) {

	tests := []struct {
		src      string
		expected []reflect.Value
	}{
		{
			// Set three variables
			src: `a, b, c := func() (int, string, bool) { return 42, "foo", true }()`,
			expected: []reflect.Value{
				reflect.ValueOf("a"),
				reflect.ValueOf("b"),
				reflect.ValueOf("c"),
			},
		},
		// {
		// 	// Return three values from named function.
		// 	src: `strings.Cut("Gopher", "ph")`,
		// 	expected: []reflect.Value{
		// 		reflect.ValueOf("Go),
		// 		reflect.ValueOf("er"),
		// 		reflect.ValueOf(true),
		// 	},
		// },
		// {
		// 	// Return three values from anonymous function.
		// 	src: `func() (int, string, bool) { return 42, "foo", true }()`,
		// 	expected: []reflect.Value{
		// 		reflect.ValueOf(42),
		// 		reflect.ValueOf("foo"),
		// 		reflect.ValueOf(true),
		// 	},
		// },
		{
			// Set two variables.
			src: `a, b := func() (int, string) { return 42, "foo" }()`,
			expected: []reflect.Value{
				reflect.ValueOf("a"),
				reflect.ValueOf("b"),
			},
		},
		// {
		// 	// Return two values from anonymous function.
		// 	src: `func() (int, string) { return 42, "foo" }()`,
		// 	expected: []reflect.Value{
		// 		reflect.ValueOf(42),
		// 		reflect.ValueOf("foo"),
		// 	},
		// },
		{
			// Set one variable.
			src: `a := func() (int) { return 42 }()`,
			expected: []reflect.Value{
				reflect.ValueOf("a"),
			},
		},
		// {
		// 	// Return one value from anonymous function.
		// 	src: `func() (int) { return 42 }()`,
		// 	expected: []reflect.Value{
		// 		reflect.ValueOf(42),
		// 	},
		// },
		// {
		// 	// Return zero values.
		// 	src:      `func() () { return }()`,
		// 	expected: nil,
		// },
		{
			// Define named function.
			src: `func foo() (int) { return 42 }`,
			expected: []reflect.Value{
				reflect.ValueOf("foo"),
			},
		},
		// {
		// 	// Define anonymous function.
		// 	src: `func () (int) { return 42 }`,
		// 	expected: []reflect.Value{
		// 		reflect.ValueOf("foo"), // what for this anonymous fn?
		// 	},
		// },
		{
			// Set two variables to two values.
			src: `a, b := 7, 42`,
			expected: []reflect.Value{
				reflect.ValueOf("a"),
				reflect.ValueOf("b"),
			},
		},
		{
			// Set one variable to one value.
			src: `a := 7`,
			expected: []reflect.Value{
				reflect.ValueOf("a"),
			},
		},
	}

	WrapTopValues = true

	for j, tt := range tests {
		debugFn("\nTest: %d\n", j)
		i := New(Options{})

		// First compile the source
		p, err := i.Compile(tt.src)
		if err != nil {
			t.Errorf("Compile error for %q: %v", tt.src, err)
			continue
		}

		// Execute with multi-value return
		res, err := i.Execute(p)
		if err != nil {
			t.Errorf("Execute error for %q: %v", tt.src, err)
			continue
		}

		// Unwrap the result
		values, err := UnwrapReflectValueSlice(res)
		if err != nil {
			t.Errorf("Unwrap error for %q: %v", tt.src, err)
			continue
		}

		for i, got := range values {
			if !got.IsValid() {
				t.Errorf("For %q value[%d]: got invalid value", tt.src, i)
				continue
			}
			// Unwrap the reflect.Value if it represents a slice of values
			if IsReflectValueSlice(got) {
				sliceValues, err := UnwrapReflectValueSlice(got)
				if err != nil {
					t.Errorf("Unwrap error for %q: %v", tt.src, err)
					continue
				}
				// Compare each value in the slice
				if len(sliceValues) != len(tt.expected) {
					t.Errorf("For %q: got %d values, want %d", tt.src, len(sliceValues), len(tt.expected))
					continue
				}
				for j, sv := range sliceValues {
					if j == len(tt.expected) {
						t.Errorf("For %q value: got length of %d, want length of %d", tt.src, len(sliceValues), len(tt.expected))
						break
					}
					want := tt.expected[j]
					if !reflect.DeepEqual(sv.Interface(), want.Interface()) {
						t.Errorf("For %q value[%d]: got %v, want %v", tt.src, j, sv.Interface(), want.Interface())
					}
				}
			} else {
				want := tt.expected[i]
				if !reflect.DeepEqual(got.Interface(), want.Interface()) {
					t.Errorf("For %q value[%d]: got %v, want %v", tt.src, i, got.Interface(), want.Interface())
				}
			}
		}
		debugFn("test passed.\n")
	}

	WrapTopValues = false
}
