package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"context"
    "fmt"
    "log"

	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

)

type User struct {
	ID    string `json:"id" form:"id" query:"id"`
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}
// 참조하도록 빼줌
var collection *mongo.Collection

func main(){

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("gofka").Collection("users")
	
	if(collection != nil){
		fmt.Println("Connected to MongoDB!")
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})


	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))


}
func saveUser(c echo.Context) error {
	ctx := c.Request().Context()

	u := new(User)
	if err := c.Bind(u); err!= nil{
		return err
	}

	insertResult, err := collection.InsertOne(ctx, u)
	if err != nil {
    	log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return c.JSON(http.StatusCreated, u)
}


func getUser(c echo.Context) error {
	ctx := c.Request().Context()
	// User ID from path `users/:id`
	id := c.Param("id")

	var user User
	if err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&user); err != nil{
		fmt.Println(id+" not exist: ", err)
		return c.String(http.StatusNotFound, id+" not exist")
	}

	return c.JSON(http.StatusOK, user)

}

func updateUser(c echo.Context) error {

	ctx := c.Request().Context()
	id := c.Param("id")

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	updatedResult, err := collection.ReplaceOne(ctx, bson.D{{"id", id}}, u)
	if err != nil || updatedResult == nil {
		return err
	}

	if updatedResult.MatchedCount == 0 {
		return c.String(http.StatusNotFound, id+" not exist")
	}
 
	fmt.Println("Updated a single documents: ", updatedResult.UpsertedID)
	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	deletedResult, err := collection.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		return err
	}

	fmt.Println("Deleted documents count: ", deletedResult.DeletedCount)
	return c.String(http.StatusNoContent, id+" is deleted")
}
