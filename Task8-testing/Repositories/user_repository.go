package repositories

import (
	"fmt"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

// NewUserRepository creates a new instance of UserRepository.
// It takes a mongo.Database and a collection name as parameters.
// It returns a domain.UserRepository interface.
// The returned UserRepository instance uses the provided database and collection for data operations.
func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

// CreateUser creates a new user in the database.
// It takes a context.Context and a *domain.User as parameters.
// It returns an error if there was a problem inserting the user into the database.
func (ur *userRepository) CreateUser(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)
	_, err := collection.InsertOne(c, user)

	return err
}

// Fetch retrieves all users from the database.
// It returns a slice of domain.User and an error, if any.
func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User
	err = cursor.All(c, &users)
	if err != nil {
		return []domain.User{}, err
	}

	return users, err
}

// GetByID retrieves a user from the database based on the provided ID.
// It takes a context.Context and the ID of the user as parameters.
// It returns the retrieved user and an error if any occurred.
func (ur *userRepository) GetByID(c context.Context, id string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user domain.User
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, fmt.Errorf("user not found")
	}

	err = collection.FindOne(c, bson.D{{Key: "_id", Value: idHex}}).Decode(&user)
	return user, err
}

// GetByEmail retrieves a user from the repository based on the provided email.
// It takes a context.Context and an email string as parameters.
// It returns a domain.User and an error.
// The domain.User represents the user entity with all its properties.
// The error indicates any error that occurred during the retrieval process.
func (ur *userRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	var user domain.User
	err := collection.FindOne(c, bson.D{{Key: "email", Value: email}}).Decode(&user)
	return user, err

}

// GetByUsername retrieves a user from the repository based on the given username.
// It takes a context.Context and a string representing the username as parameters.
// It returns a domain.User and an error. The domain.User represents the retrieved user,
// and the error indicates any errors that occurred during the retrieval process.
func (ur *userRepository) GetByUsername(c context.Context, userName string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	var user domain.User
	err := collection.FindOne(c, bson.D{{Key: "user_name", Value: userName}}).Decode(&user)
	return user, err
}

// UpdateUser updates a user's information in the database.
// It takes a context, an ID string, and a UserUpdate struct as parameters.
// The function returns the updated User object and an error, if any.
// The UserUpdate struct contains optional fields for updating the user's first name, last name, username, email, password, and user role.
// If the ID is not found in the database, the function returns an error.
// If the update is successful, the function returns the updated User object.
func (ur *userRepository) UpdateUser(c context.Context, id string, user domain.UserUpdate) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	updateFields := make(bson.M)
	if user.First_Name != nil {
		updateFields["first_name"] = user.First_Name
	}
	if user.Last_Name != nil {
		updateFields["last_name"] = user.Last_Name
	}
	if user.User_Name != nil {
		updateFields["user_name"] = user.User_Name
	}
	if user.Email != nil {
		updateFields["email"] = user.Email
	}
	if user.Password != nil {
		updateFields["password"] = user.Password
	}
	if user.User_Role != nil {
		updateFields["user_role"] = user.User_Role
	}

	// Add the updated_at field with the current timestamp
	updateFields["updated_at"] = time.Now()

	var updatedUser domain.User
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return updatedUser, fmt.Errorf("user not found")
	}

	update := bson.D{{Key: "$set", Value: updateFields}}
	result, err := collection.UpdateOne(
		c,
		bson.D{{Key: "_id", Value: idHex}},
		update,
	)

	if err != nil {
		return updatedUser, err
	}

	if result.MatchedCount == 0 {
		return domain.User{}, fmt.Errorf("no user found with the given ID")
	}

	err = collection.FindOne(c, bson.D{{Key: "_id", Value: idHex}}).Decode(&updatedUser)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil

}

// Promote promotes a user to an admin role.
// It takes a context and the user ID as parameters.
// It returns the updated user and an error if any
func (ur *userRepository) Promote(c context.Context, id string) (domain.User, error) {
	admin := "ADMIN"
	newUser := domain.UserUpdate{
		User_Role: &admin,
	}

	updatedUser, err := ur.UpdateUser(c, id, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("unable to promote user")
	}

	return updatedUser, nil
}
