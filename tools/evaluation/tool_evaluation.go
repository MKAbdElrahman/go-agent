package evaluation

import (
	"errors"
	"fmt"
	"go-agent/metadata"
	"reflect"
)

var (
	ErrNotAFunction     = errors.New("entry is not a function")
	ErrArgumentMismatch = errors.New("argument count mismatch")
	ErrArgumentType     = errors.New("argument type mismatch")
	ErrFunctionPanic    = errors.New("function execution panicked")
)

// Tool represents a function along with its metadata and documentation.
type Tool struct {
	Metadata metadata.FunctionMetaData `json:"metadata"`
	Function interface{}               `json:"function"`
}

func (t Tool) Evaluate(args []interface{}) ([]interface{}, error) {
	functionValue := reflect.ValueOf(t.Function)
	if functionValue.Kind() != reflect.Func {
		return nil, ErrNotAFunction
	}

	functionType := functionValue.Type()
	isVariadic := functionType.IsVariadic()
	numIn := functionType.NumIn()

	if isVariadic {
		if len(args) < numIn-1 {
			return nil, fmt.Errorf("%w: expected at least %d arguments, got %d", ErrArgumentMismatch, numIn-1, len(args))
		}
	} else {
		if len(args) != numIn {
			return nil, fmt.Errorf("%w: expected %d arguments, got %d", ErrArgumentMismatch, numIn, len(args))
		}
	}

	argValues, err := convertArguments(args, functionType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArgumentType, err)
	}

	results, err := callFunction(functionValue, argValues)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFunctionPanic, err)
	}

	return extractResults(results)
}

// convertArguments converts and validates the provided arguments against the function's expected types.
func convertArguments(args []interface{}, functionType reflect.Type) ([]reflect.Value, error) {
	numIn := functionType.NumIn()
	isVariadic := functionType.IsVariadic()
	argValues := make([]reflect.Value, 0, len(args))

	for i, arg := range args {
		var expectedType reflect.Type

		if isVariadic && i >= numIn-1 {
			// For variadic functions, the last argument type is the element type of the slice.
			expectedType = functionType.In(numIn - 1).Elem()
		} else {
			expectedType = functionType.In(i)
		}

		argValue := reflect.ValueOf(arg)

		// Handle slices for variadic functions
		if isVariadic && i >= numIn-1 && argValue.Kind() == reflect.Slice {
			// Unpack the slice into individual arguments
			for j := 0; j < argValue.Len(); j++ {
				elem := argValue.Index(j)
				if !elem.Type().AssignableTo(expectedType) {
					return nil, fmt.Errorf("argument %d (element %d): expected %s, got %s", i+1, j+1, expectedType, elem.Type())
				}
				argValues = append(argValues, elem)
			}
			continue
		}

		// Handle non-slice arguments
		if !argValue.Type().AssignableTo(expectedType) {
			if argValue.CanConvert(expectedType) {
				argValue = argValue.Convert(expectedType)
			} else {
				return nil, fmt.Errorf("argument %d: expected %s, got %s", i+1, expectedType, argValue.Type())
			}
		}

		argValues = append(argValues, argValue)
	}

	return argValues, nil
}

func callFunction(functionValue reflect.Value, argValues []reflect.Value) (results []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("function panicked with argument(s) %v: %v", argValues, r)
		}
	}()

	results = functionValue.Call(argValues)
	return results, nil
}

func extractResults(results []reflect.Value) ([]interface{}, error) {
	if len(results) == 0 {
		return nil, nil
	}

	lastResult := results[len(results)-1]
	if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !lastResult.IsNil() {
			return convertResults(results[:len(results)-1]), lastResult.Interface().(error)
		}
		return convertResults(results[:len(results)-1]), nil
	}

	return convertResults(results), nil
}

func convertResults(results []reflect.Value) []interface{} {
	converted := make([]interface{}, len(results))
	for i, res := range results {
		converted[i] = res.Interface()
	}
	return converted
}
