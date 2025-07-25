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

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // Unique identifier for the user
	Username string             `bson:"username"`      // Username of the user
	Password string             `bson:"password"`      // Hashed password (excluded from JSON responses)
	IsAdmin  bool               `bson:"is_admin"`      // Flag indicating if the user is an admin
}

func (u *User) toDomain() domain.User {
	return domain.User{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func fromDomainToUser(u *domain.User) (User, error) {
	objID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:       objID,
		Username: u.Username,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}, nil
}

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

// FetchByUsername retrieves a user by their username
// Returns ErrUserNotFound if no user is found
func (ur *userRepository) FetchByUsername(ctx context.Context, username string) (domain.User, error) {
	users := ur.database.Collection(ur.collection)
	filter := bson.D{
		{Key: "username", Value: username},
	}

	var user User
	result := users.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, result.Err()
	}

	err := result.Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user.toDomain(), nil
}

// Create inserts a new user into the collection
// Assigns the generated ObjectID to the user
func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	userEntity := User{
		Username: user.Username,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}

	users := ur.database.Collection(ur.collection)

	// Insert user
	result, err := users.InsertOne(ctx, userEntity)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unexpected InsertedID type: %T", result.InsertedID)
	}
	user.ID = objID.Hex()

	return nil
}

// PromoteByUserID sets the IsAdmin field to true for a specific user
// Returns the number of matched documents
func (ur *userRepository) PromoteByUserID(ctx context.Context, userID string) (int, error) {
	// check for valid id
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, domain.ErrInvalidUserID
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
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Failed to decode users in GetAllUsers")
			continue
		}
		results = append(results, user.toDomain())
	}

	return results, nil
}

// CheckIfUsernameExists query the username and returns true with nil if it does
// or returns 0 if it does not. Might return errors in the way
func (ur *userRepository) CheckIfUsernameExists(c context.Context, username string) (bool, error) {
	users := ur.database.Collection(ur.collection)
	filter := bson.M{"username": username}
	count, err := users.CountDocuments(c, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountUsers counts the number of documents inside users
func (ur *userRepository) CountUsers(ctx context.Context) (int, error) {
	users := ur.database.Collection(ur.collection)
	count, err := users.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// FetchByUserID retrieves a user by their unique ID
// Returns ErrUserNotFound if no user is found
func (ur *userRepository) FetchByUserID(ctx context.Context, userID string) (domain.User, error) {
	// check for valid id
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, domain.ErrInvalidUserID
	}

	users := ur.database.Collection(ur.collection)
	// prepare filter
	filter := bson.D{{Key: "_id", Value: objID}}

	// fetch the result
	result := users.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	// decode the user
	var user User
	err = result.Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user.toDomain(), nil
}
