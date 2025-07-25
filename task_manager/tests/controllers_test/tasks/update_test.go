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

// TestUpdateTask is used to test UpdateTask controller
func (s *SuiteTaskUsecase) TestUpdateTask() {

	original := TaskRequest{
		ID:          sampleDatas[0].ID,
		Title:       sampleDatas[0].Title,
		Description: sampleDatas[0].Description,
		Status:      sampleDatas[0].Status,
		DueDate:     sampleDatas[0].DueDate,
	}
	modified := original
	modified.Status = "completed"

	tests := []TaskListTestCase{
		{
			Name:     "Status is modified",
			Payload:  original,
			Expected: http.StatusOK,
			Validate: func(t domain.Task) {
				require.Equal(s.T(), modified.ID, t.ID)
				require.Equal(s.T(), modified.Title, t.Title)
				require.Equal(s.T(), modified.Description, t.Description)
				require.Equal(s.T(), modified.Status, t.Status)
				require.WithinDuration(s.T(), modified.DueDate, t.DueDate, 1*time.Second)
			},
			MockSetup: func() {
				s.mockUsecase.On("UpdateByTaskID", mock.Anything, mock.MatchedBy(func(t *domain.Task) bool {
					t.Status = "completed"
					return true
				})).Return(nil).Once()
			},
		},
		{
			Name:     "Invalid task ID",
			Payload:  original,
			Expected: http.StatusBadRequest,
			Validate: func(t domain.Task) {},
			MockSetup: func() {
				s.mockUsecase.On("UpdateByTaskID", mock.Anything, mock.Anything).
					Return(domain.ErrInvalidTaskID).Once()
			},
		},
		{
			Name: "Due date in the past",
			Payload: TaskRequest{
				ID:          original.ID,
				Title:       original.Title,
				Description: original.Description,
				Status:      original.Status,
				DueDate:     time.Now().Add(-24 * time.Hour), // past date
			},
			Expected: http.StatusBadRequest,
			Validate: func(t domain.Task) {},
			MockSetup: func() {
				s.mockUsecase.On("UpdateByTaskID", mock.Anything, mock.Anything).
					Return(domain.ErrInvalidDueDate).Once()
			},
		},
		{
			Name:     "Task not found",
			Payload:  original,
			Expected: http.StatusBadRequest,
			Validate: func(t domain.Task) {},
			MockSetup: func() {
				s.mockUsecase.On("UpdateByTaskID", mock.Anything, mock.Anything).
					Return(domain.ErrTaskNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		s.PrepareTest(tt)

		body, _ := json.Marshal(tt.Payload)
		req, _ := http.NewRequest(http.MethodPut, "/tasks/"+tt.Payload.ID, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()

		s.router.ServeHTTP(resp, req)

		var responseBody domain.Task

		if resp.Code == http.StatusOK {
			err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
			require.NoError(s.T(), err)
			tt.Validate(responseBody)
		}

		require.Equal(s.T(), tt.Expected, resp.Code)
	}
	s.mockUsecase.AssertExpectations(s.T())

}
