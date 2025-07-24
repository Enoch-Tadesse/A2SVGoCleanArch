package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DTO used only inside repository
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	DueDate     time.Time          `bson:"due_date"`
	Status      string             `bson:"status"`
}

// Convert domain.Task → repositories.Task
func fromDomainToTask(t *domain.Task) (Task, error) {
	objID, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:          objID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Status:      t.Status,
	}, nil
}

// Convert repositories.Task → domain.Task
func (t *Task) toDomain() domain.Task {
	return domain.Task{
		ID:          t.ID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Status:      t.Status,
	}
}

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

// UpdateByTaskID updates a task in the collection using its ID
// Returns the number of matched and modified documents
func (tr *taskRepository) UpdateByTaskID(ctx context.Context, task *domain.Task) (int, int, error) {

	taskEntity, err := fromDomainToTask(task)
	if err != nil {
		return 0, 0, err
	}

	tasks := tr.database.Collection(tr.collection)
	// prepare filter and update
	filter := bson.D{{Key: "_id", Value: taskEntity.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: taskEntity.Title},
			{Key: "description", Value: taskEntity.Description},
			{Key: "due_date", Value: taskEntity.DueDate},
			{Key: "status", Value: taskEntity.Status},
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
		return 0, domain.ErrInvalidTaskID
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
	taskEntity, err := fromDomainToTask(task)
	if err != nil {
		return err
	}

	tasks := tr.database.Collection(tr.collection)

	result, err := tasks.InsertOne(ctx, taskEntity)
	if err != nil {
		return err
	}
	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unexpected InsertedID type: %T", result.InsertedID)
	}
	task.ID = objID.Hex()
	return nil
}

// FetchByTaskID retrieves a task by its ID
// Returns ErrTaskNotFound if no document is found
func (tr *taskRepository) FetchByTaskID(ctx context.Context, taskID string) (domain.Task, error) {
	// check for valid ID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.Task{}, domain.ErrInvalidUserID
	}

	tasks := tr.database.Collection(tr.collection)

	// prepare filter
	filter := bson.D{{Key: "_id", Value: objID}}

	// fetch the result
	result := tasks.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, result.Err()
	}

	// decode the task
	var task Task
	err = result.Decode(&task)
	if err != nil {
		return domain.Task{}, err
	}
	return task.toDomain(), nil
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
		var task Task
		if err := cursor.Decode(&task); err != nil {
			log.Println("Failed to decode tasks in GetAllTasks")
			continue
		}
		results = append(results, task.toDomain())
	}

	return results, nil
}
