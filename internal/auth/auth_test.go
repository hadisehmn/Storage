package auth

import (
	models "go-practice/STORAGE/internal/model"
	"testing"
)

func TestSignUp(t *testing.T) {

	tests := []struct {
		name      string
		req       models.SignUpRequest
		wantError string
	}{
		{
			name: "name required",
			req: models.SignUpRequest{
				Name:     "",
				Email:    "ali@test.com",
				Password: "123456",
			},
			wantError: "name is required",
		},
		{
			name: "email required",
			req: models.SignUpRequest{
				Name:     "Ali",
				Email:    "",
				Password: "123456",
			},
			wantError: "email is required",
		},
		{
			name: "password required",
			req: models.SignUpRequest{
				Name:     "Ali",
				Email:    "ali@test.com",
				Password: "",
			},
			wantError: "password is required",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := &AuthService{}

			_, err := service.SignUp(test.req)

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != test.wantError {
				t.Errorf(
					"expected %q, got %q",
					test.wantError,
					err.Error(),
				)
			}
		})
	}
}

func TestSignIn(t *testing.T) {

	tests := []struct {
		name      string
		req       models.SignInRequest
		wantError string
	}{
		{
			name: "email required",
			req: models.SignInRequest{
				Email:    "",
				Password: "123456",
			},
			wantError: "email is required",
		},
		{
			name: "password required",
			req: models.SignInRequest{
				Email:    "ali@test.com",
				Password: "",
			},
			wantError: "password is required",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			service := &AuthService{}

			_, err := service.SignIn(test.req)

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != test.wantError {
				t.Errorf(
					"expected %q, got %q",
					test.wantError,
					err.Error(),
				)
			}
		})
	}
}
