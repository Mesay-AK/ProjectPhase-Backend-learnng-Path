package repository

import (
	"context"
	"fmt"
	domain "task_manager_clean/Domain"
	infrastructure "task_manager_clean/Infrastructure"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection("users")
	return &UserRepository{UserCollection: collection}
}

func (repo *UserRepository) Register(ctx context.Context, user *domain.User) error {
	_, err := repo.UserCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Print(err)
		return fmt.Errorf("failed to register user: %v", err)
	}
	return nil
}

func (repo *UserRepository) LogIn(ctx context.Context, email, password string) (domain.User, string, string, error) {
	var foundUser domain.User
	err := repo.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("invalid Email or Password")
	}

	passwordIsValid, msg := infrastructure.VerifyPassword(password, *foundUser.Password)
	if !passwordIsValid {
		return foundUser, "", "", fmt.Errorf(msg)
	}

	token, refreshToken, err := infrastructure.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.UserType, *foundUser.UserID)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("error generating tokens: %v", err)
	}

	err = infrastructure.UpdateAllTokens(repo.UserCollection, token, refreshToken, *foundUser.UserID)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("error updating tokens: %v", err)
	}

	return foundUser, token, refreshToken, nil
}

func (repo *UserRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	cursor, err := repo.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	filter := bson.M{"_id": objectID}
	err = repo.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) PromoteUser(ctx context.Context, id string, promoted *domain.User) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"user_type": "ADMIN"}}
	result := repo.UserCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if err := result.Decode(&promoted); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return promoted, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID format")
	}

	filter := bson.M{"_id": objectID}
	result, err := repo.UserCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (repo *UserRepository) CheckExistingUser(email string) (int64, error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second) 
	defer cancel()

	count, err := repo.UserCollection.CountDocuments(c, bson.M{"email": email})
	if err != nil {
		return 0, fmt.Errorf("error checking for email: %v", err)
	}
	if count > 0 {
		return count, fmt.Errorf("this email already exists")
	}
	return count, nil
}
