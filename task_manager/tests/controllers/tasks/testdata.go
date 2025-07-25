package tasks
// this file contains shared data, and struct within the tasks test

import (
	"time"

	domain "github.com/A2SVTask7/Domain"
)

// request struct
type TaskRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

type TaskListTestCase struct {
	Name          string              // name of the test
	MockSetup     func()              // mock strup
	Expected      int                 // expected status
	ValidateSlice func([]domain.Task) // validation datas of a slice
	Validate      func(domain.Task)   // validation of a single data
	Payload       TaskRequest         // payload if it is a post request
}

// valid sample tasks
var sampleDatas = []domain.Task{
	{
		ID:          "task1",
		Title:       "First Task",
		Description: "Do something",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	},
	{
		ID:          "task2",
		Title:       "Second Task",
		Description: "Do something else",
		DueDate:     time.Now().Add(48 * time.Hour),
		Status:      "done",
	},
}
