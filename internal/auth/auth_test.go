package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test(t *testing.T) {
	type testCase struct {
		name           string
		userID         uuid.UUID
		secret         string
		validateSecret string
		expiresIn      time.Duration
		expectError    bool
	}

	runCases := []testCase{
		{
			name:           "valid token",
			userID:         uuid.New(),
			secret:         "secret",
			validateSecret: "secret",
			expiresIn:      time.Hour,
			expectError:    false,
		},
		{
			name:           "invalid secret",
			userID:         uuid.New(),
			secret:         "secret",
			validateSecret: "wrong-secret",
			expiresIn:      time.Hour,
			expectError:    true,
		},
		{
			name:           "expired token",
			userID:         uuid.New(),
			secret:         "secret",
			validateSecret: "secret",
			expiresIn:      -1 * time.Minute,
			expectError:    true,
		},
		{
			name:           "tampered token",
			userID:         uuid.New(),
			secret:         "secret",
			validateSecret: "secret",
			expiresIn:      time.Hour,
			expectError:    true,
		},
	}

	testCases := runCases

	passCount := 0
	failCount := 0

	for _, test := range testCases {
		token, err := MakeJWT(test.userID, test.secret, test.expiresIn)
		if err != nil {
			failCount++
			t.Errorf("MakeJWT failed for test %q: %v", test.name, err)
			continue
		}

		if test.name == "tampered token" {
			token += "abc"
		}

		uid, err := ValidateJWT(token, test.validateSecret)

		passed := true
		if test.expectError {
			if err == nil {
				passed = false
			}
		} else {
			if err != nil || uid != test.userID {
				passed = false
			}
		}

		if !passed {
			failCount++
			t.Errorf("---------------------------------\nTest Failed:\nName:\n- %s\nExpect error:\n- %v\nActual:\n- uid = %v\n- err = %v\nFail\n\n", test.name, test.expectError, uid, err)
		} else {
			passCount++
			fmt.Printf("---------------------------------\nTest Passed:\nName:\n- %s\nExpect error:\n- %v\nActual:\n- uid=%v\n- err=%v\nPass\n\n", test.name, test.expectError, uid, err)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n\n", passCount, failCount)
}
