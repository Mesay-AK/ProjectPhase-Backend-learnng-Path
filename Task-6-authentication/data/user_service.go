package data

import (
	"context"
	"errors"
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"task_manager_with_jwt/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserService struct {
	UserCollection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	collection := db.Collection("users")
	return &UserService{UserCollection: collection}
}

var ErrUserNotFound = errors.New("user not found")

func (service *UserService) Register(ctx context.Context, user *models.User) error {
	email := *user.Email
	name := *user.Name
	userType := *user.User_type

	// Check for existing email and get user count
	count, err := service.CheckExistingUser(email)
	if err != nil {
		return err
	}

	hashedPassword, err := service.HashPassword(*user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	user.Password = &hashedPassword

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.User_id = new(string)
	*user.User_id = user.ID.Hex()

	// If this is the first user, assign admin role
	if count == 0 {
		userType = "ADMIN"
	}else{
		userType = "USER"
	}


	token, refreshToken, err := service.GenerateAllTokens(email, name, userType, *user.User_id)
	if err != nil {
		return fmt.Errorf("failed to generate tokens: %v", err)
	}
	user.Token = &token
	user.Refresh_token = &refreshToken
	user.User_type = &userType

	_, err = service.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}

	return nil
}

func (service *UserService) LogIn(email, password string) (models.User, string, string, error) {
	var foundUser models.User
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()


	err := service.UserCollection.FindOne(c, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("invalid Email or Password")
	}

	passwordIsValid, msg := VerifyPassword(password, *foundUser.Password)
	if !passwordIsValid {
		return foundUser, "", "", fmt.Errorf(msg)
	}

	token, refreshToken, err := service.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.User_type, *foundUser.User_id)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("error generating tokens")
	}

	err = service.UpdateAllTokens(token, refreshToken, *foundUser.User_id)
	if err != nil {
		return foundUser, "", "", fmt.Errorf("error updating tokens")
	}

	return foundUser, token, refreshToken, nil
}


// These are requests for an authorized Admin only



func (service *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	collection := service.UserCollection

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	return users, nil
}


// user service reciever method to get user from the user collection
func (service *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	collection := service.UserCollection

	filter := bson.M{"_id": id}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}


// user service reciever method to promote a user on the user collection by it's ID
func (service *UserService) PromoteUser(id string, updateduser models.User) (*models.User, error) {
	collection := service.UserCollection

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{"user_type": "ADMIN"},
	}

	result := collection.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var user models.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (service *UserService) DeleteUser(id string) error {
	collection := service.UserCollection

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}