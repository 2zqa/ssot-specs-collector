package aggregator

import (
	"fmt"
	"strings"
)

type Aggregator struct {
	Subject string
	Errors  []error
}

const initialErrorsCapacity = 10

func New(s string) *Aggregator {
	return &Aggregator{
		Subject: s,
		Errors:  make([]error, 0, initialErrorsCapacity)}
}

// CheckAndAdd appends the given error to the Aggregator's error list if it is not nil.
// Returns false if the error was added, true otherwise.
//
// This method is intended to be used with nil checks, e.g.:
//
//	if data, err := someFunc(); a.CheckAndAdd(err) {
//		// No error, do something with data
//	}
func (a *Aggregator) CheckAndAdd(err error) bool {
	if err != nil {
		a.Errors = append(a.Errors, err)
		return false
	}
	return true
}

// CheckPredicateAndAdd evaluates the given predicate and, if it is false,
// appends the given error message to the Aggregator's error list.
// Returns the value of the predicate.
//
// Example:
//
//	if a.CheckPredicateAndAdd(data != "", "no data found") {
//		// Data is not empty, do something with it
//	}
func (a *Aggregator) CheckPredicateAndAdd(predicate bool, errorMessage string) bool {
	if !predicate {
		a.Errors = append(a.Errors, fmt.Errorf(errorMessage))
	}
	return predicate
}

// Error returns a string with all errors concatenated, joined by a semicolon.
func (a *Aggregator) Error() string {
	errorCount := len(a.Errors)
	if errorCount == 0 {
		return ""
	}
	errorStrings := make([]string, 0, errorCount)

	for _, err := range a.Errors {
		errorStrings = append(errorStrings, err.Error())
	}
	joinedErrorStrings := strings.Join(errorStrings, "; ")
	return fmt.Sprintf("errors occured while fetching %s information: %s", a.Subject, joinedErrorStrings)
}

// ErrorOrNil returns any errors reported, or nil if none.
func (a *Aggregator) ErrorOrNil() error {
	if len(a.Errors) == 0 {
		return nil
	}
	return a
}
