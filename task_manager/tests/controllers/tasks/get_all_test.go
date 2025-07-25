package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestGetTask is used to test GetAllTask controller
func (s *SuiteTaskUsecase) TestGetTask() {
	// expected returning value
	expectedTasks := sampleDatas
	// a list of tests
	tests := []TaskListTestCase{

		{
			Name: "Successfully fetch tasks",
			MockSetup: func() {
				s.mockUsecase.On("FetchAllTasks", mock.Anything).Return(expectedTasks, nil)
			},
			Expected: http.StatusOK,
			ValidateSlice: func(body []domain.Task) {
				require.Equal(s.T(), len(expectedTasks), len(body))
				require.Equal(s.T(), expectedTasks[0].Title, body[0].Title)
				require.Equal(s.T(), expectedTasks[1].Status, body[1].Status)
			},
		},
		{
			Name: "internal server error",
			MockSetup: func() {
				s.mockUsecase.On("FetchAllTasks", mock.Anything).Return([]domain.Task{}, fmt.Errorf("db error")).Once()
			},
			Expected: http.StatusInternalServerError,
			ValidateSlice: func(body []domain.Task) {
				require.Equal(s.T(), 0, len(body))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)
			// create request and response
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			resp := httptest.NewRecorder()

			s.router.ServeHTTP(resp, req)

			// check for valid responses
			require.Equal(s.T(), tt.Expected, resp.Code)

			var responseBody []domain.Task
			if resp.Code == http.StatusOK {
				err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
				require.NoError(s.T(), err, "Failed to unmarshal response body")
			}
			tt.ValidateSlice(responseBody)

		})
	}
	s.mockUsecase.AssertExpectations(s.T())

}
