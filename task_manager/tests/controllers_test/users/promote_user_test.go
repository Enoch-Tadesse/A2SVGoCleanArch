package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
)

// TestPromoteUser is used to test Promote controller
func (s *SuiteUserUsecase) TestPromoteUser() {
	tests := []UserListTestCase{
		{
			Name: "Missing ID",
			Payload: UserRequest{
				ID: "",
			},
			Expected: http.StatusNotFound,
		},
		{
			Name: "Successful Promotion",
			Payload: UserRequest{
				ID: "user1234",
			},
			MockSetup: func() {
				s.mockUsecase.On("PromoteByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return true
				})).Return(nil).Once()
			},
			Expected: http.StatusOK,
		},
		{
			Name: "User Not Found",
			Payload: UserRequest{
				ID: "user1234",
			},
			MockSetup: func() {
				s.mockUsecase.On("PromoteByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return true
				})).Return(domain.ErrUserNotFound).Once()
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Random Error",
			Payload: UserRequest{
				ID: "user1234",
			},
			MockSetup: func() {
				s.mockUsecase.On("PromoteByUserID", mock.Anything, mock.MatchedBy(func(id string) bool {
					return true
				})).Return(fmt.Errorf("some random error")).Once()
			},
			Expected: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			req, _ := http.NewRequest(http.MethodPatch, "/promote/"+tt.Payload.ID, nil)
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)
			s.Equal(tt.Expected, resp.Code)
		})
	}
	s.mockUsecase.AssertExpectations(s.T())
}
