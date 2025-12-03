package main

import "errors"

// =============================================================================
// Test that closures are analyzed in their own scope
// =============================================================================

// Before the fix: The linter would check the closure's return statement
// against the OUTER function's named returns (result, err), causing false positives.
//
// After the fix: The linter checks the closure independently - it should only
// complain that the closure itself doesn't have named returns (if that's the rule),
// NOT that the closure isn't using "result" and "err" from the outer function.

func outerWithClosure() (result string, err error) {
	// This closure has named returns and uses them - should be fine
	helper := func() (helperResult string, helperErr error) {
		helperResult = "helper"
		helperErr = errors.New("helper error")
		return helperResult, helperErr
	}

	result, err = helper()
	return result, err
}

// Closure without named returns - the linter will flag the closure itself,
// but should NOT flag the closure's return statement for not using the outer function's
// named returns
func outerWithUnnamedClosure() (result string, err error) {
	helper := func() (string, error) { // want `unnamed return with type "string" found - named returns are required` `unnamed return with type "error" found - named returns are required`
		// The linter should NOT complain that this return doesn't use "result" and "err"
		// because those are from the outer function, not this closure
		return "helper", errors.New("error")
	}

	result, err = helper()
	return result, err
}
