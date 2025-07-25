package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// TestGetAllUsers is used to test GetAllUsers controller
func (s *SuiteUserUsecase) TestGetAllUsers() {
	expected := sampleUsers

	tests := []UserListTestCase{
		{
			Name: "Retrives Successfully",
			ValidateSlice: func(u []domain.User) {
				s.Equal(len(expected), len(u))
				s.Equal(expected[0].ID, u[0].ID)
				s.Equal(expected[1].Username, u[1].Username)
			},
			Expected: http.StatusOK,
			MockSetup: func() {
				s.mockUsecase.On("FetchAllUsers", mock.Anything).Return(expected, nil)
			},
		},
		{
			Name:     "Random Error",
			Expected: http.StatusInternalServerError,
			MockSetup: func() {
				s.mockUsecase.On("FetchAllUsers", mock.Anything).Return(expected, errors.New("random error"))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)

			var responseBody []domain.User
			if resp.Code == http.StatusOK {
				err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
				s.NoError(err)
			}

			s.Equal(tt.Expected, resp.Code)
			if tt.ValidateSlice != nil {
				tt.ValidateSlice(responseBody)
			}
		})
	}
}
