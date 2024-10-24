package aggregator

import (
	"errors"
	"testing"
)

const aggregatorName string = "test"
const testErrorMessage string = "test error"

func TestAggregatorCheckPredicateWithTrue(t *testing.T) {
	a := New(aggregatorName)

	// Test adding a true predicate
	if a.CheckPredicateAndAdd(true, testErrorMessage) != true {
		t.Errorf("Expected CheckPredicate to return true, but got false")
	}

	if len(a.Errors) != 0 {
		t.Errorf("Expected Errors to have length 0, but got %d", len(a.Errors))
	}
}

func TestAggregatorCheckPredicateWithFalse(t *testing.T) {
	a := New(aggregatorName)

	// Test adding a true predicate
	if a.CheckPredicateAndAdd(false, testErrorMessage) != false {
		t.Errorf("Expected CheckPredicate to return true, but got false")
	}

	if len(a.Errors) != 1 {
		t.Errorf("Expected Errors to have length 1, but got %d", len(a.Errors))
	}
}

func TestAggregatorCheckWithError(t *testing.T) {
	a := New(aggregatorName)
	err := errors.New(testErrorMessage)

	// Test adding a non-nil error
	if a.CheckAndAdd(err) != false {
		t.Errorf("Expected Check to return false, but got true")
	}
	if len(a.Errors) != 1 {
		t.Errorf("Expected Errors to have length 1, but got %d", len(a.Errors))
	}
}

func TestAggregatorCheckWithNilError(t *testing.T) {
	a := New(aggregatorName)

	// Test adding a nil error
	if a.CheckAndAdd(nil) != true {
		t.Errorf("Expected Check to return true, but got false")
	}
	if len(a.Errors) != 0 {
		t.Errorf("Expected Errors to have length 0, but got %d", len(a.Errors))
	}
}

func TestAggregatorGetErrorString(t *testing.T) {
	a := New(aggregatorName)
	if a.Error() != "" {
		t.Errorf("Expected Error to return empty string, but got '%s'", a.Error())
	}

	a.CheckAndAdd(errors.New("test error 1"))
	a.CheckAndAdd(errors.New("test error 2"))

	want := "errors occured while fetching test information: test error 1; test error 2"
	got := a.Error()
	if got != want {
		t.Errorf("Expected Error to return '%s', but got '%s'", want, got)
	}
}

func TestAggregatorAggregateErrors(t *testing.T) {
	a := New(aggregatorName)
	if a.ErrorOrNil() != nil {
		t.Errorf("Expected ErrorOrNil to return nil, but got %v", a.ErrorOrNil())
	}

	a.CheckAndAdd(errors.New(testErrorMessage))
	if a.ErrorOrNil() == nil {
		t.Errorf("Expected ErrorOrNil to return non-nil error, but got nil")
	}
}
