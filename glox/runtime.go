package main

// LoxRuntime interface defines the methods needed for Lox runtime operations.
// This interface is implemented by both GLox (for normal execution) and TestGLox (for testing).
// It provides a common contract for error reporting and runtime state management.
type LoxRuntime interface {
	// error reports a general error at the specified line
	error(line int, message string)

	// parseError reports a parsing error at the specified token
	parseError(token Token, message string)

	// runtimeError reports a runtime error that occurred during execution
	runtimeError(err error)
}
