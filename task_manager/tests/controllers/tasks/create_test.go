package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCreateTask is used to test CreateTask controller
func (s *SuiteTaskUsecase) TestCreateTask() {

	// prepare tests
	tests := []TaskListTestCase{
		{
			Name: "passing test",
			Payload: TaskRequest{
				Title:       "new task",
				Description: "not mandatory",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "pending",
			},
			Expected: http.StatusCreated,
			MockSetup: func() {
				s.mockUsecase.On("Create", mock.Anything, mock.MatchedBy(func(t *domain.Task) bool {
					return t.Title == "new task"
				})).Return(nil).Once()
			},
		},
		{
			Name: "missing title",
			Payload: TaskRequest{
				Title:       "",
				Description: "not mandatory",
				DueDate:     time.Now().Add(24 * time.Hour),
				Status:      "pending",
			},
			Expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			// prepare request body
			body, _ := json.Marshal(tt.Payload)

			// create a request and response
			req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte(body)))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// run the request
			s.router.ServeHTTP(resp, req)

			require.Equal(s.T(), tt.Expected, resp.Code)
		})
		s.mockUsecase.AssertExpectations(s.T())
	}
}
