package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	handler "filmsharing/user-management-service/handler"
	model "filmsharing/user-management-service/model"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func main() {
	engine := gin.Default()

	// Initialize the JWT middleware
	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Register middlewares
	engine.Use(handlerMiddleWare(authMiddleware))

	// Register routes
	registerRoute(engine, authMiddleware)

	engine.Run(":8080")
}

func registerRoute(r *gin.Engine, jwtHandler *jwt.GinJWTMiddleware) {
	r.NoRoute(jwtHandler.MiddlewareFunc(), handleNoRoute())

	r.POST("/register", handler.CreateUser)
	r.POST("/login", jwtHandler.LoginHandler)

	auth := r.Group("/auth", jwtHandler.MiddlewareFunc())
	auth.GET("/refresh_token", jwtHandler.RefreshHandler)
	auth.GET("/users", handler.GetUsers)
	auth.GET("/users/:id", handler.GetUserByID)
}

func handlerMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func initParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "filmsharing",
		Key:         []byte("thisisasecretkey"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		PayloadFunc: payloadFunc(),

		// IdentityHandler: identityHandler(),
		Authenticator: authenticator(),
		Authorizator:  authorizator(),
		Unauthorized:  unauthorized(),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*model.User); ok {
			return jwt.MapClaims{
				"username": v.Username,
				"email":    v.Email,
				"name":     v.Name,
			}
		}
		return jwt.MapClaims{}
	}
}

// func identityHandler() func(c *gin.Context) interface{} {
// 	return func(c *gin.Context) interface{} {
// 		claims := jwt.ExtractClaims(c)
// 		return &User{
// 			UserName: claims[identityKey].(string),
// 		}
// 	}
// }

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginUser model.User
		if err := c.ShouldBind(&loginUser); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		username := loginUser.Username
		password := loginUser.Password

		if checkUserExists(username, password) {
			return &model.User{}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func checkUserExists(username string, password string) bool {
	for _, user := range handler.UserMemDB {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		return true
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

// func helloHandler(c *gin.Context) {
// 	claims := jwt.ExtractClaims(c)
// 	user, _ := c.Get(identityKey)
// 	c.JSON(200, gin.H{
// 		"userID":   claims[identityKey],
// 		"userName": user.(*User).UserName,
// 		"text":     "Hello World.",
// 	})
// }
