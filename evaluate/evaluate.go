package evaluate

import (
	"fmt"
	"go-agent/tools"
	"reflect"
)

type EvaluationResult struct {
	Result []interface{} `json:"result"` // The result of the function evaluation
	Error  error         `json:"error"`  // The error returned by the function, if any
}

func Evaluate(store tools.ToolStore, functionName string, args []interface{}) (result EvaluationResult) {
	// Step 1: Look up the function in the store
	entry, exists := store[functionName]
	if !exists {
		return EvaluationResult{
			Error: fmt.Errorf("function '%s' not found in store", functionName),
		}
	}

	fmt.Println("Tool Selected: ", functionName)

	// Step 2: Use reflection to get the function type
	functionValue := reflect.ValueOf(entry.Function)
	if functionValue.Kind() != reflect.Func {
		return EvaluationResult{
			Error: fmt.Errorf("entry for '%s' is not a function", functionName),
		}
	}

	// Step 3: Get the function's type information
	functionType := functionValue.Type()

	// Step 4: Validate the number of arguments
	if len(args) != functionType.NumIn() {
		return EvaluationResult{
			Error: fmt.Errorf("function '%s' expects %d arguments, but %d were provided",
				functionName, functionType.NumIn(), len(args)),
		}
	}

	// Step 5: Convert arguments to reflect.Value and validate types
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		expectedType := functionType.In(i)

		// Check if the argument type matches the expected type
		if !argValue.Type().AssignableTo(expectedType) {
			// Try to convert the argument to the expected type
			if argValue.CanConvert(expectedType) {
				argValue = argValue.Convert(expectedType)
			} else {
				return EvaluationResult{
					Error: fmt.Errorf("argument %d for function '%s' has type %s, but expected %s",
						i+1, functionName, argValue.Type(), expectedType),
				}
			}
		}
		argValues[i] = argValue
	}

	// Step 6: Call the function with the provided arguments
	defer func() {
		if r := recover(); r != nil {
			result = EvaluationResult{
				Error: fmt.Errorf("panic occurred while calling function '%s': %v", functionName, r),
			}
		}
	}()

	results := functionValue.Call(argValues)

	// Step 7: Split the results into result and error
	var resultInterfaces []interface{}
	var funcError error

	if len(results) > 0 {
		// Check if the last return value is an error
		lastResult := results[len(results)-1]
		if lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			// If the last result is an error, extract it
			if !lastResult.IsNil() {
				funcError = lastResult.Interface().(error)
			}
			// Exclude the error from the results
			resultInterfaces = make([]interface{}, len(results)-1)
			for i := 0; i < len(results)-1; i++ {
				resultInterfaces[i] = results[i].Interface()
			}
		} else {
			// If there's no error, include all results
			resultInterfaces = make([]interface{}, len(results))
			for i, result := range results {
				resultInterfaces[i] = result.Interface()
			}
		}
	}

	return EvaluationResult{
		Result: resultInterfaces,
		Error:  funcError,
	}
}
