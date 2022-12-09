package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/CobaKauPikirkan/e-commerce-golang/database"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct{
	Email 		string 
	FirstName 	string
	LastName	string
	Uid			string
	jwt.StandardClaims
	// jwt.RegisteredClaims
}

var UserData  *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstname string,lastname string, uid string) (signedtoken string, signedrefreshtoken string, err error) {
	// claims := &SignedDetails{
	// 	Email: email,
	// 	FirstName: firstname,
	// 	LastName: lastname,
	// 	Uid: uid,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	// 	},
	// 	// RegisteredClaims: jwt.RegisteredClaims{
	// 	// 	ExpiresAt: &jwt.NumericDate{Time: time.Now().Local().Add(time.Hour * time.Duration(24))},
	// 	// },
	// }

	// refreshclaims := &SignedDetails{
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	// 	},
	// 	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 	// 	ExpiresAt: &jwt.NumericDate{Time: time.Now().Local().Add(time.Hour * time.Duration(24))},
	// 	// },
	// }

	// token, err := jwt.NewWithClaims(jwt.SigningMethodES256,claims).SignedString([]byte(SECRET_KEY))
	// if err != nil {
	// 	return "", "", err
	// }

	// refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshclaims).SignedString([]byte(SECRET_KEY)) 
	// if err != nil {
	// 	log.Panic(err)
	// 	return 
	// }

	// return token, refreshtoken, err
	claims := &SignedDetails{
		Email:      email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panicln(err)
		return
	}
	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claim *SignedDetails, msg string ) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The Token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}

func UpdateAllToken(signedtoken string, signedrefreshtoken string, userid string)  {
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()
	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})

	
	updated_at, _ := time.Parse(time.RFC822, time.Now().Format(time.RFC822))
	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: updated_at})

	upsert := true

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value:  updateobj}}, &opt)
	if err != nil {
		log.Panic(err)
		return
	}

	defer cancel()
}