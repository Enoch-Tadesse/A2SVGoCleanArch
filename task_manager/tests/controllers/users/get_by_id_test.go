package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// TestGetUserByID is used to test GetUserByID controller
func (s *SuiteUserUsecase) TestGetUserByID() {
	expected := sampleUsers[0]

	tests := []UserListTestCase{
		{
			Name:     "Retrives Successfully",
			Expected: http.StatusOK,
			Payload: UserRequest{
				ID: expected.ID,
			},
			MockSetup: func() {
				s.mockUsecase.On("FetchByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return id == expected.ID
				})).Return(expected, nil)
			},
			Validate: func(u domain.User) {
				s.Equal(expected.ID, u.ID)
				s.Equal(expected.Username, u.Username)
				s.Equal(expected.IsAdmin, u.IsAdmin)
			},
		},
		{
			Name:     "Missing ID",
			Expected: http.StatusNotFound,
			Payload: UserRequest{
				ID: "",
			},
		},
		{
			Name: "Invalid User ID",
			Payload: UserRequest{
				ID: expected.ID,
			},
			Expected: http.StatusBadRequest,
			MockSetup: func() {
				s.mockUsecase.On("FetchByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return id == expected.ID
				})).Return(domain.User{}, domain.ErrInvalidUserID)
			},
		},
		{
			Name: "Random Error",
			Payload: UserRequest{
				ID: expected.ID,
			},
			Expected: http.StatusInternalServerError,
			MockSetup: func() {
				s.mockUsecase.On("FetchByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return id == expected.ID
				})).Return(domain.User{}, errors.New("Random Error"))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.Payload.ID, nil)
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)

			var responseBody domain.User
			if resp.Code == http.StatusOK {
				err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
				s.NoError(err)
			}

			s.Equal(tt.Expected, resp.Code)
			if tt.ValidateSlice != nil {
				tt.Validate(responseBody)
			}
		})
	}
}
