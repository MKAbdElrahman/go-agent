package tools

import (
	"fmt"
	"go-agent/metadata"
	"reflect"
)

// Tool represents a function along with its metadata and documentation.
type Tool struct {
	Metadata metadata.FunctionMetaData `json:"metadata"`
	Doc      string                    `json:"doc"`
	Function interface{}               `json:"function"`
}

// Evaluate executes the function stored in the Tool with the provided arguments.
func (t Tool) Evaluate(args []interface{}) ([]interface{}, error) {
	functionValue := reflect.ValueOf(t.Function)
	if functionValue.Kind() != reflect.Func {
		return nil, ErrNotAFunction
	}

	functionType := functionValue.Type()
	if len(args) != functionType.NumIn() {
		return nil, fmt.Errorf("%w: expected %d arguments, got %d", ErrArgumentMismatch, functionType.NumIn(), len(args))
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
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		expectedType := functionType.In(i)

		if !argValue.Type().AssignableTo(expectedType) {
			if argValue.CanConvert(expectedType) {
				argValue = argValue.Convert(expectedType)
			} else {
				return nil, fmt.Errorf("argument %d: expected %s, got %s", i+1, expectedType, argValue.Type())
			}
		}
		argValues[i] = argValue
	}
	return argValues, nil
}

// callFunction calls the function with the provided arguments and handles panics.
func callFunction(functionValue reflect.Value, argValues []reflect.Value) (results []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	results = functionValue.Call(argValues)
	return results, nil
}

// extractResults processes the function's return values, separating results from errors.
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

// convertResults converts reflect.Value results to interface{}.
func convertResults(results []reflect.Value) []interface{} {
	converted := make([]interface{}, len(results))
	for i, res := range results {
		converted[i] = res.Interface()
	}
	return converted
}
