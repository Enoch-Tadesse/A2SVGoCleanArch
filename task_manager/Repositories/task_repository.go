package repositories

import (
	"context"
	"errors"
	"log"

	domain "github.com/A2SVTask7/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// taskRepository implements the domain.TaskRepository interface
type taskRepository struct {
	database   mongo.Database // MongoDB database instance
	collection string         // Name of the collection to operate on
}

// NewTaskRepository returns a new taskRepository instance
func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

// ErrTaskNotFound is returned when a task with the given ID does not exist
var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidTaskID = errors.New("invalid task id")
)

// UpdateByTaskID updates a task in the collection using its ID
// Returns the number of matched and modified documents
func (tr *taskRepository) UpdateByTaskID(ctx context.Context, id string, task *domain.Task) (int, int, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, 0, ErrInvalidTaskID
	}
	// set the task id
	task.ID = objID

	tasks := tr.database.Collection(tr.collection)
	// prepare filter and update
	filter := bson.D{{Key: "_id", Value: task.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: task.Title},
			{Key: "description", Value: task.Description},
			{Key: "due_date", Value: task.DueDate},
			{Key: "status", Value: task.Status},
		}},
	}

	// execute update command
	result, err := tasks.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, 0, err
	}
	return int(result.MatchedCount), int(result.ModifiedCount), nil
}

// DeleteByTaskID deletes a task by its ID
// Returns the number of documents deleted
func (tr *taskRepository) DeleteByTaskID(ctx context.Context, taskID string) (int, error) {
	// check for valid ID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return 0, ErrInvalidUserID
	}

	tasks := tr.database.Collection(tr.collection)

	filter := bson.D{{Key: "_id", Value: objID}}

	result, err := tasks.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}

// Create inserts a new task into the collection
// Assigns the generated ObjectID back to the task
func (tr *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	tasks := tr.database.Collection(tr.collection)

	result, err := tasks.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FetchByTaskID retrieves a task by its ID
// Returns ErrTaskNotFound if no document is found
func (tr *taskRepository) FetchByTaskID(ctx context.Context, taskID string) (domain.Task, error) {
	// check for valid ID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.Task{}, ErrInvalidUserID
	}

	tasks := tr.database.Collection(tr.collection)

	// prepare filter
	filter := bson.D{{Key: "_id", Value: objID}}

	// fetch the result
	result := tasks.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return domain.Task{}, ErrTaskNotFound
		}
		return domain.Task{}, result.Err()
	}

	// decode the task
	var task domain.Task
	err = result.Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// FetchAllTasks retrieves all tasks from the collection
func (tr *taskRepository) FetchAllTasks(ctx context.Context) ([]domain.Task, error) {
	tasks := tr.database.Collection(tr.collection)

	var results []domain.Task
	cursor, err := tasks.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}

	for cursor.TryNext(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Failed to decode tasks in GetAllTasks")
			continue
		}
		results = append(results, task)
	}

	return results, nil
}
