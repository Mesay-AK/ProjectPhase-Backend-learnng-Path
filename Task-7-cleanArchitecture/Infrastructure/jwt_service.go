package Infrastructure

import(
    "context"
    "fmt"
    "os"
    "time"
    jwt "github.com/dgrijalva/jwt-go"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

)

var secret_key = []byte(os.Getenv("SECRET_KEY"))

type signedDetails struct {
    Email       string
    Name        string
    Uid         string
    User_type   string
    jwt.StandardClaims
}

func GenerateAllTokens(email, fName, userType, uid string) (signedToken, signedRefreshToken string, err error) {
    // if secretKey == "" {
    //     fmt.Print("secret_key not set")
    //     return "", "", fmt.Errorf("SECRET_KEY environment variable is not set")
    // }

    claims := &signedDetails{
        Email:     email,
        Name:      fName,
        Uid:       uid,
        User_type: userType,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Add(time.Hour * 24).Unix(),
        },
    }

    refreshClaims := &signedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().UTC().Add(time.Hour * 168).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err = token.SignedString(secret_key)
    if err != nil {
        fmt.Println("tokens : ", token)
        return "", "", err
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedRefreshToken, err = refreshToken.SignedString(secret_key)
    if err != nil {
        fmt.Println("refreshToken: ", refreshToken)
        return "", "", err
    }
    fmt.Println("token and refreshtoken generation successful")
    return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userCollection *mongo.Collection, signedToken, signedRefreshToken, userId string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel()

    updateObj := bson.D{
        {Key: "token", Value: signedToken},
        {Key: "refresh_token", Value: signedRefreshToken},
        {Key: "updated_at", Value: time.Now().UTC()},
    }

    filter := bson.M{"user_id": userId}
    opt := options.Update().SetUpsert(true)

    _, err := userCollection.UpdateOne(
        ctx,
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

func ValidateToken(signedToken string) (claims *signedDetails, msg string) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &signedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return secret_key, nil
        },
    )

    if err != nil {
        return nil, err.Error()
    }

    claims, ok := token.Claims.(*signedDetails)
    if !ok {
        return nil, "Invalid token"
    }

    if claims.ExpiresAt < time.Now().UTC().Unix() {
        return nil, "Expired Token"
    }

    return claims, ""
}
