package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// TestLoginUser is used to test Login controller
func (s *SuiteUserUsecase) TestLoginUser() {
	tests := []UserListTestCase{

		{
			Name: "User Not Found",
			Payload: UserRequest{
				Username: "Username",
				Password: "Password",
			},
			MockSetup: func() {
				s.mockUsecase.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(domain.User{}, "", domain.ErrUserNotFound).Once()
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Username Missing",
			Payload: UserRequest{
				Username: "",
				Password: "Password",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Password missing",
			Payload: UserRequest{
				Username: "Username",
				Password: "",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Random Error",
			Payload: UserRequest{
				Username: "Username",
				Password: "Password",
			},
			MockSetup: func() {
				s.mockUsecase.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(domain.User{}, "", errors.New("random error")).Once()
			},
			Expected: http.StatusInternalServerError,
		},
		{
			Name: "Invalid Password",
			Payload: UserRequest{
				Username: "Username",
				Password: "Password",
			},
			MockSetup: func() {
				s.mockUsecase.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(domain.User{}, "", domain.ErrIncorrectPassword).Once()
			},
			Expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			body, _ := json.Marshal(tt.Payload)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(body)))
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)

			s.Equal(tt.Expected, resp.Code)
		})
	}
}
