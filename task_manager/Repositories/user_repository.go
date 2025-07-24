package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	domain "github.com/A2SVTask7/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// userRepository implements the domain.UserRepository interface
type userRepository struct {
	database   mongo.Database // MongoDB database instance
	collection string         // Name of the users collection
}

// NewUserRepository returns a new userRepository instance
func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

// Custom error variables for user-related failures
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUserID     = errors.New("invalid user id")
)

// FetchByUsername retrieves a user by their username
// Returns ErrUserNotFound if no user is found
func (ur *userRepository) FetchByUsername(ctx context.Context, username string) (domain.User, error) {
	users := ur.database.Collection(ur.collection)
	filter := bson.D{
		{Key: "username", Value: username},
	}

	var user domain.User
	result := users.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, result.Err()
	}

	err := result.Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// Create inserts a new user into the collection
// Assigns the generated ObjectID to the user
func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	users := ur.database.Collection(ur.collection)

	// Insert user
	result, err := users.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	user.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// PromoteByUserID sets the IsAdmin field to true for a specific user
// Returns the number of matched documents
func (ur *userRepository) PromoteByUserID(ctx context.Context, userID string) (int, error) {
	// check for valid id
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, ErrInvalidUserID
	}

	users := ur.database.Collection(ur.collection)
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_admin", Value: true},
		}},
	}

	result, err := users.UpdateByID(ctx, objID, update)
	if err != nil {
		return 0, err
	}

	return int(result.MatchedCount), err
}

// FetchAllUsers retrieves all users from the collection
func (ur *userRepository) FetchAllUsers(ctx context.Context) ([]domain.User, error) {
	users := ur.database.Collection(ur.collection)

	var results []domain.User
	cursor, err := users.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	for cursor.TryNext(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Failed to decode users in GetAllUsers")
			continue
		}
		results = append(results, user)
	}

	return results, nil
}

// FetchByUserID retrieves a user by their unique ID
// Returns ErrUserNotFound if no user is found
func (ur *userRepository) FetchByUserID(ctx context.Context, userID string) (domain.User, error) {
	// check for valid id
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, ErrInvalidUserID
	}

	users := ur.database.Collection(ur.collection)
	// prepare filter
	filter := bson.D{{Key: "_id", Value: objID}}

	// fetch the result
	result := users.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}

	// decode the user
	var user domain.User
	err = result.Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
