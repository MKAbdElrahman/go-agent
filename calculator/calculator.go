// Package calculator provides advanced arithmetic and mathematical operations.
package calculator

import (
	"errors"
	"math"
)

func GetPublicFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Add":        Add,
		"Subtract":   Subtract,
		"Multiply":   Multiply,
		"Divide":     Divide,
		"SquareRoot": SquareRoot,
		"Power":      Power,
		"Factorial":  Factorial,
		"Modulus":    Modulus,
		"Sin":        Sin,
		"Cos":        Cos,
		"Tan":        Tan,
		"Log":        Log,
		"Log10":      Log10,
	}
}

// Add returns the sum of two numbers.
// @param a: The first number.
// @param b: The second number.
// @return float64: The sum of a and b.
// @example: Add(3, 4) // returns 7
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference between two numbers.
// @param a: The first number.
// @param b: The second number.
// @return float64: The difference between a and b.
// @example: Subtract(10, 4) // returns 6
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of two numbers.
// @param a: The first number.
// @param b: The second number.
// @return float64: The product of a and b.
// @example: Multiply(3, 4) // returns 12
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of two numbers.
// @param a: The dividend.
// @param b: The divisor.
// @return float64: The quotient of a divided by b.
// @constraint b != 0: b must not be zero.
// @example: Divide(10, 2) // returns 5
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

// SquareRoot calculates the square root of a number.
// @param x: The number to calculate the square root of.
// @return float64: The square root of the input number.
// @constraint x >= 0: x must be non-negative.
// @example: SquareRoot(4) // returns 2
func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, errors.New("square root of a negative number is not allowed")
	}
	return math.Sqrt(a), nil
}

// Power returns the result of raising a to the power of b.
// @param a: The base.
// @param b: The exponent.
// @return float64: The result of a raised to the power of b.
// @example: Power(2, 3) // returns 8
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// Factorial calculates the factorial of a non-negative integer.
// @param n: The number to calculate the factorial of.
// @return float64: The factorial of the input number.
// @constraint n >= 0: n must be non-negative.
// @example: Factorial(5) // returns 120
func Factorial(n int) (float64, error) {
	if n < 0 {
		return 0, errors.New("factorial of a negative number is not allowed")
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result, nil
}

// Modulus returns the remainder of a divided by b.
// @param a: The dividend.
// @param b: The divisor.
// @return float64: The remainder of a divided by b.
// @constraint b != 0: b must not be zero.
// @example: Modulus(10, 3) // returns 1
func Modulus(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return math.Mod(a, b), nil
}

// Sin calculates the sine of a number in radians.
// @param x: The angle in radians.
// @return float64: The sine of the input angle.
// @example: Sin(math.Pi / 2) // returns 1
func Sin(x float64) float64 {
	return math.Sin(x)
}

// Cos calculates the cosine of a number in radians.
// @param x: The angle in radians.
// @return float64: The cosine of the input angle.
// @example: Cos(0) // returns 1
func Cos(x float64) float64 {
	return math.Cos(x)
}

// Tan calculates the tangent of a number in radians.
// @param x: The angle in radians.
// @return float64: The tangent of the input angle.
// @example: Tan(math.Pi / 4) // returns 1
func Tan(x float64) float64 {
	return math.Tan(x)
}

// Log calculates the natural logarithm of a number.
// @param x: The number to calculate the logarithm of.
// @return float64: The natural logarithm of the input number.
// @constraint x > 0: x must be positive.
// @example: Log(2.71828) // returns 1
func Log(x float64) (float64, error) {
	if x <= 0 {
		return 0, errors.New("logarithm of a non-positive number is not allowed")
	}
	return math.Log(x), nil
}

// Log10 calculates the base-10 logarithm of a number.
// @param x: The number to calculate the logarithm of.
// @return float64: The base-10 logarithm of the input number.
// @constraint x > 0: x must be positive.
// @example: Log10(100) // returns 2
func Log10(x float64) (float64, error) {
	if x <= 0 {
		return 0, errors.New("logarithm of a non-positive number is not allowed")
	}
	return math.Log10(x), nil
}
