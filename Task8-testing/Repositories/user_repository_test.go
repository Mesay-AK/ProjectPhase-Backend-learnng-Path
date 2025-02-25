package repositories

import (
	"context"
	"testing"
	"time"

	domain "github.com/Mesay-AK/ProjectPhase-Backend-learnng-Path/Task8-testing/Domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var id = primitive.NewObjectID()
var firstName = "John"
var lastName = "Doe"
var userName = "johndoe"
var email = "johndoe@example.com"
var password = "password123"
var userRole = "USER"

type UserRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection string
	repo       domain.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Connect to the MongoDB test instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	// Ping the MongoDB instance to ensure connection is established
	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	suite.client = client
	suite.db = client.Database("testdb")
	suite.collection = "users"
	suite.repo = NewUserRepository(*suite.db, suite.collection)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {

	user := &domain.User{
		ID:         id,
		First_Name: &firstName,
		Last_Name:  &lastName,
		User_Name:  &userName,
		Email:      &email,
		Password:   &password,
		User_Role:  &userRole,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	err := suite.repo.CreateUser(context.Background(), user)
	suite.Require().NoError(err)

	// Verify user was inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.User

	err = collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), user.Email, result.Email)
	suite.NoError(err, "no error when creating and finding the same user")
}

func (suite *UserRepositoryTestSuite) TestFetch() {
	fetchedUsers, err := suite.repo.Fetch(context.Background())
	suite.Require().NoError(err)

	// Verify that the correct number of users were fetched
	assert.Equal(suite.T(), 1, len(fetchedUsers))
	for _, user := range fetchedUsers {
		assert.Equal(suite.T(), email, *user.Email, "email same with the previous inserted")
		assert.Equal(suite.T(), firstName, *user.First_Name, "first name same with the previous inserted")
		assert.Equal(suite.T(), lastName, *user.Last_Name, "last name same with the previous inserted")
		assert.Equal(suite.T(), userName, *user.User_Name, "user name same with the previous inserted")
		assert.Equal(suite.T(), userRole, *user.User_Role, "user role same with the previous inserted")
		// Ensure the password field is not returned due to the projection
		assert.Empty(suite.T(), user.Password, "password is not visible due to projection")
	}

}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	userFound, err := suite.repo.GetByID(context.Background(), id.Hex())
	suite.Require().NoError(err)

	// Verify the correct user is being fetched.
	assert.Equal(suite.T(), email, *userFound.Email, "same email with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same first name with the previous inserted")
	assert.Equal(suite.T(), lastName, *userFound.Last_Name, "same last name with the previous inserted")
	assert.Equal(suite.T(), userName, *userFound.User_Name, "same user name with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same  user role with the previous inserted")
}

func (suite *UserRepositoryTestSuite) TestGetByEmail() {
	userFound, err := suite.repo.GetByEmail(context.Background(), email)
	suite.Require().NoError(err)

	// Verify the correct user is being fetched.
	assert.Equal(suite.T(), email, *userFound.Email, "same email with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same first name with the previous inserted")
	assert.Equal(suite.T(), lastName, *userFound.Last_Name, "same last name with the previous inserted")
	assert.Equal(suite.T(), userName, *userFound.User_Name, "same user name with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same  user role with the previous inserted")
}

func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	userFound, err := suite.repo.GetByUsername(context.Background(), userName)
	suite.Require().NoError(err)

	// Verify the correct user is being fetched.
	assert.Equal(suite.T(), email, *userFound.Email, "same email with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same first name with the previous inserted")
	assert.Equal(suite.T(), lastName, *userFound.Last_Name, "same last name with the previous inserted")
	assert.Equal(suite.T(), userName, *userFound.User_Name, "same user name with the previous inserted")
	assert.Equal(suite.T(), firstName, *userFound.First_Name, "same  user role with the previous inserted")
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	newUserName := "Abebe"
	newUser := domain.UserUpdate{
		First_Name: &firstName,
		Last_Name:  &lastName,
		User_Name:  &newUserName,
		Email:      &email,
		Password:   &password,
		User_Role:  &userRole,
	}
	updatedUser, err := suite.repo.UpdateUser(context.Background(), id.Hex(), newUser)
	suite.Require().NoError(err)

	// find out if the user is updated
	collection := suite.db.Collection(suite.collection)
	var result domain.User

	err = collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), email, *result.Email)
	// check if user name is updated
	assert.Equal(suite.T(), newUserName, *result.User_Name)
	assert.Equal(suite.T(), id.Hex(), updatedUser.ID.Hex())
	suite.NoError(err, "no error when creating and finding the same user")
}

func (suite *UserRepositoryTestSuite) TestPromote() {
	promotedUser, err := suite.repo.Promote(context.Background(), id.Hex())
	suite.Require().NoError(err)

	// Verify the correct user is being fetched.
	assert.Equal(suite.T(), email, *promotedUser.Email, "same email with the previous inserted")

	// Verify if the user is promoted
	assert.Equal(suite.T(), "ADMIN", *promotedUser.User_Role, "new user role is ADMIN")

}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
