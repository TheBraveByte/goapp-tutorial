### Handlers for user requests
To process every api requests make by the user to each endpoint, we need to compute the logic that will carry out the user request successfully. And for us to accompanied that ,we will create methods to help process each request.

To meet up with what we need. We will set up a repository pattern in the handlers package which is similar to how we put together our database queries methods in the query package.

Let's go straight to setting a proper structure in creating our handlers


```go
package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/go-app/modules/config"
	"github.com/yusuf/go-app/modules/database"
	"github.com/yusuf/go-app/modules/database/query"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoApp struct {
	App *config.GoAppTools
	DB  database.DBRepo
}

func NewGoApp(app *config.GoAppTools, db *mongo.Client) *GoApp {
	return &GoApp{
		App: app,
		DB:  query.NewGoAppDB(app, db),
	}
}

func (ga *GoApp) Home () gin.HandlerFunc{
	return func (ctx *gin.Context){
		fmt.Sprintln("Creating a scalable web application with Gin")
	}
}
```

We created a GoApp struct with field App of type config.GoAppTools which is a struct type that we added in the config package and a field DB of type database.DBRepo which is an interface that implements any type defined but here the database.DBRepo implements all the database query methods that will be use in each of the computed handlers' logic.

Now let's define a function of NewGoApp that takes in app of a pointer to the type config.GoAppTools and db of pointer to type mongo.Client as parameters then returns a type of pointer to the GoApp struct. This function defined will come in handy when called in the main package(main.go) and the value of its parameters will be passed.

All that done properly, we created a handler method using the GoApp struct that was added earlier to handle the user request to the application homepage, this handler return a type of gin.HandlerFunc. In the method, our return value is an anonymous function with a parameter of pointer to gin.Context. Gin context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request with the appropriate HTTP status code (200)and render a JSON response.


With the new approach toward building a stable, scalable and easy to maintain application , we need to modify our previous code in main.go

```go
// connecting to the database
	client := driver.Connection(uri)
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			app.ErrorLogger.Fatal(err)
			return
		}
	}()
    appRouter := gin.New()
    appRouter.GET("/", func(ctx *gin.Context) {
        app.InfoLogger.Println("Creating a scalable web application with Gin")
})
```
to this

```go

	// connecting to the database
	client := driver.Connection(uri)
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			app.ErrorLogger.Fatal(err)
			return
		}
	}()
    appRouter := gin.New()

    goApp := handlers.NewGoApp(&app, client)
    Routes(appRouter, goApp)

```
Now, what as really changed in our current code compare to the previous code?
After the Gin instance was created, we called on the NewGoApp function defined in the handlers package and pass the value of it parameters i.e (app struct and MongoDB client) whose value is a pointer to the struct created in the package. Since the struct was used to implement methods, this allows us to have access to all the methods that implements the struct type from the value returned from the function.

Routes for Endpoints

To map up each route's endpoints to their respective handlers which handle the user requests and also to add default security protocols to secure every request made by the user.

To get this done , We will define a function called Routes which takes in two parameter r of type pointer to gin.Engine, the Gin Engine is the Gin framework's instance, it contains the muxer, default middleware and configuration settings, and we've created an instance of Engine, by using New() in main.go and g of type pointer to the GoApp function that was previously define in the handlers package.
In addition to that , we will need to add some user data in session as cookies for easy usability among handlers methods and to do that we will need to install a gin session package tool. let's install that as shown below.

```go

  go get github.com/gin-contrib/sessions
  go get github.com/gin-contrib/sessions/cookie

```

That done, let go ahead and add up some code to set up the Routes function as mentioned.

```go
package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yusuf/go-app/handlers"
)

func Routes(r *gin.Engine, g *handlers.GoApp) {
	router := r.Use(gin.Logger(), gin.Recovery())

	router.GET("/", g.Home())

	// set up for storing details as cookies
	cookieData := cookie.NewStore([]byte("go-app"))
	router.Use(sessions.Sessions("session", cookieData))
}
```
From the implemented code, We make use of the Use method (an instance of Gin instance) attaches a global middleware to the router. i.e. the middleware attached through Use() will be included in the handlers chain for every single request. The middleware use are the Gin Logger instance to write out logs and the Recovery function that help recovers from any panics if there was one while the server is running.
Next we added the HTTP GET request for the Homepage endpoints for the API and went further to set up the cookies store that will be use in session and have it added to the Use method as well.