package tasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestGetByTaskID is used to test GetTaskByID
func (s *SuiteTaskUsecase) TestGetByTaskID() {
	expected := domain.Task{
		ID:          "1234",
		Title:       "Expected Task",
		Description: "Do something",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}
	tests := []TaskListTestCase{
		{
			Name: "All good",
			Payload: TaskRequest{
				ID: "1234",
			},
			Expected:  http.StatusOK,
			MockSetup: func() { s.mockUsecase.On("FetchByTaskID", mock.Anything, "1234").Return(expected, nil) },
			Validate: func(t domain.Task) {
				require.Equal(s.T(), expected.ID, t.ID)
				require.Equal(s.T(), expected.Title, t.Title)
				require.Equal(s.T(), expected.Description, t.Description)
				require.Equal(s.T(), expected.Status, t.Status)
				require.WithinDuration(s.T(), expected.DueDate, t.DueDate, time.Second) // allow 1s diff
			},
		},
		{
			Name: "No ID parameter",
			Payload: TaskRequest{
				ID: "",
			},
			Expected:  http.StatusNotFound,
			MockSetup: func() {},
			Validate: func(t domain.Task) {
				require.Empty(s.T(), t.ID)
				require.Empty(s.T(), t.Title)
				require.Empty(s.T(), t.Description)
			},
		},
		{
			Name: "Task not found",
			Payload: TaskRequest{
				ID: "9999",
			},
			Expected: http.StatusNotFound,
			MockSetup: func() {
				s.mockUsecase.On("FetchByTaskID", mock.Anything, "9999").
					Return(domain.Task{}, domain.ErrTaskNotFound).Once()
			},
			Validate: func(t domain.Task) {
				require.Empty(s.T(), t.ID)
				require.Empty(s.T(), t.Title)
				require.Empty(s.T(), t.Description)
			},
		},
		{
			Name: "Internal server error",
			Payload: TaskRequest{
				ID: "1234",
			},
			Expected: http.StatusInternalServerError,
			MockSetup: func() {
				s.mockUsecase.On("FetchByTaskID", mock.Anything, "1234").
					Return(domain.Task{}, errors.New("db error")).Once()
			},
			Validate: func(t domain.Task) {
				require.Empty(s.T(), t.ID)
				require.Empty(s.T(), t.Title)
				require.Empty(s.T(), t.Description)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.Name, func() {
			s.PrepareTest(tt)

			req, _ := http.NewRequest(http.MethodGet, "/tasks/"+tt.Payload.ID, nil)
			resp := httptest.NewRecorder()
			s.router.ServeHTTP(resp, req)

			require.Equal(s.T(), tt.Expected, resp.Code)

			var responseBody domain.Task
			if resp.Code == http.StatusOK {
				err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
				require.NoError(s.T(), err)
			}

			tt.Validate(responseBody)

		})
	}
	s.mockUsecase.AssertExpectations(s.T())
}
