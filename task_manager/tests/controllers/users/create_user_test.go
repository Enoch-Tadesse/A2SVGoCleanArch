package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// TestCreateUser is use to test CreateUser controller
func (s *SuiteUserUsecase) TestCreateUser() {
	tests := []UserListTestCase{
		{
			Name: "Valid Request",
			Payload: UserRequest{
				Username: "username",
				Password: "password",
			},
			Expected: http.StatusCreated,
			MockSetup: func() {
				s.mockUsecase.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return true
				})).Return(nil).Once()
			},
		},

		{
			Name: "User Already Exists",
			Payload: UserRequest{
				Username: "username",
				Password: "password",
			},
			Expected: http.StatusBadRequest,
			MockSetup: func() {
				s.mockUsecase.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return true
				})).Return(domain.ErrUserAlreadyExists).Once()
			},
		},

		{
			Name: "Missing Username",
			Payload: UserRequest{
				Username: "",
				Password: "password",
			},
			Expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {

			s.PrepareTest(tt)

			// create the request
			body, _ := json.Marshal(tt.Payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(body)))
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)

			s.Equal(tt.Expected, resp.Code)

		})
	}
	s.mockUsecase.AssertExpectations(s.T())
}
