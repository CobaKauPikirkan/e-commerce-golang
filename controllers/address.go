package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CobaKauPikirkan/e-commerce-golang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		address , err :=primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var addresses models.Address

		addresses.AdressId =primitive.NewObjectID()
		err = c.BindJSON(&addresses)
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$adress_id"}, {Key: "count",Value: bson.D{primitive.E{Key: "$sum",Value: 1 }}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var addressinfo []bson.M
		err = pointcursor.All(ctx, &addressinfo)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		var size int32

		for _, adress_no := range addressinfo{
		count := adress_no["count"]
		size = count.(int32)	
		}
		if size < 2{
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update :=bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_,err :=UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		}else{
			c.IndentedJSON(400, "not allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}

		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil{
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set",Value: bson.D{primitive.E{Key: "address.0.house_name",Value: editAddress.House},{Key: "address.0.street_name",Value: editAddress.Street}, {Key: "address.0.city_name",Value: editAddress.City}, {Key: "address.0.postal_code", Value: editAddress.PostalCode}}}}
		_, err = UserCollection.UpdateOne(ctx,filter, update)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "succesfuly updated the home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var editAddress models.Address
		
		if err := c.BindJSON(&editAddress); err != nil{
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
	
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name",Value: editAddress.House},{Key: "address.1.street_name",Value: editAddress.Street}, {Key: "address.1.city_name",Value: editAddress.City}, {Key: "address.1.postal_code", Value: editAddress.PostalCode}}}}
		
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "successfuly updated the Work Address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id :=c.Query("id")
		if user_id =="" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "internal server error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}

		update := bson.D{{Key : "$set", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}
		_ , err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "not found")
			return
		}
		defer cancel()

		ctx.Done()

		c.IndentedJSON(200, "succesfully deleted")
	}
}