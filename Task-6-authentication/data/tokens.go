package data

import (
	"context"
	"fmt"
	"time"
	"os"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	// "task_manager_with_jwt/models"
)


var SECRET_KEY string = os.Getenv("SECRET_KEY")

type signedDetails struct {
	Email        string
	Name           string
	Uid          string
	User_type    string
	jwt.StandardClaims
}


func (us *UserService) CheckExistingUser(email string) (int64, error) {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := us.UserCollection.CountDocuments(c, bson.M{"email": email})
	if err != nil {
		return 0, fmt.Errorf("error checking for email: %v", err)
	}
	if count > 0 {
		return count, fmt.Errorf("this email already exists")
	}
	return count, nil
}


func (us *UserService) GenerateAllTokens(email, fName,  userType, uid string) (signedToken, signedRefreshToken string, err error) {


	claims := &signedDetails{
		Email:      email,
		Name: 		fName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(), 
		},
	}

	refreshClaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(), 
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err = refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}



func (us *UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (us *UserService) VerifyPassword(userPassword, foundPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(userPassword))
	check := true
	mssg := ""

	if err != nil {
		mssg = "Invalid email or password"
		check = false
	}

	return check, mssg
}

func (us *UserService) UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	filter := bson.M{"user_id": userId}
	opt := options.Update().SetUpsert(true)

	_, err := us.UserCollection.UpdateOne(
		c,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		opt,
	)

	if err != nil {
		return fmt.Errorf("error updating tokens: %v", err)
	}

	return nil
}

func VerifyPassword(userPassword, foundPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(userPassword) )
	check := true
	mssg := ""

	if err != nil {
		mssg = "Invalid email or password"
		check = false
	}

	return check, mssg


}


func ValidateToken(signedToken string) (claims *signedDetails, msg string) {
    // Parse the token and extract claims
    token, err := jwt.ParseWithClaims(
        signedToken,
        &signedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )
    
    // Check if there was an error parsing the token
    if err != nil {
        msg = err.Error() // Set error message if token parsing fails
        return nil, msg
    }

    // Assert the token claims to our defined struct
    claims, ok := token.Claims.(*signedDetails)
    if !ok {
        msg = "Invalid token" // Set message if claims assertion fails
        return nil, msg
    }

    // Check if the token is expired
    if claims.ExpiresAt < time.Now().Local().Unix() {
        msg = "Expired Token" // Set message if token is expired
        return nil, msg
    }

    // Return claims and an empty message if everything is valid
    return claims, ""
}

