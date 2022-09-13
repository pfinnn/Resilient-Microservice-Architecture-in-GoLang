package circuitbreaker

import "errors"

type ICircuitBreaker interface {
	Close()
	Open()
	isOpen() bool
	isClosed() bool
	GetThreshold() int
	GetErrorCount() int
	ResetErrorCounter()
	IncrementErrorCounter()
}

// private CircuitBreaker struct that can be initialized with the NewCircuitBreaker function call
// it will be closed on default
type CircuitBreaker struct {
	closed    bool
	threshold int
	counter   int
	ErrInvalidParameter error
	ErrCircuitBreakerOpen error
}



func NewCircuitBreaker (errorThreshold int) (CircuitBreaker, error){

var (
	// Constructor Errors
	ErrInvalidParameter = errors.New("invalid parameters")
	// Util Errors
	ErrCircuitBreakerOpen = errors.New("circuit breaker open, denying requests to save resources")
)

// Throw error if parameters invalid
if errorThreshold < 1 {
	panic(ErrInvalidParameter)
}

cb := CircuitBreaker {
	// default values, all fields need to be initialized
	true,
	errorThreshold,
	0,
	ErrInvalidParameter, // error into its own struct?
	ErrCircuitBreakerOpen,
}

return cb, nil
}

// GETTER

func (cb *CircuitBreaker) isOpen() bool{
	return !cb.closed
}

func (cb *CircuitBreaker) IsClosed() bool{
	return cb.closed
}

func (cb *CircuitBreaker) GetThreshold() int {
	return cb.threshold
}

func (cb *CircuitBreaker) GetErrorCount() int {
	return cb.counter
}

// SETTER

func (cb *CircuitBreaker) Close(){
	cb.closed = true
}

func (cb *CircuitBreaker) Open() {
	cb.closed = false
}

func (cb *CircuitBreaker) ResetErrorCounter(){
	cb.counter = 0
}

func (cb *CircuitBreaker) IncrementErrorCounter() {
	cb.counter += 1
}

// TODO: Add Timer that closes the CB after interval without error



