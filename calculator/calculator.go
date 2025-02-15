// Package calculator provides basic arithmetic operations such as addition, subtraction,
// multiplication, division, and square root.
package calculator

import (
	"errors"
	"math"
)

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
