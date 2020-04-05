package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// User type
type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// or
	// return c.XML(http.StatusCreated, u)
}

//e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

// e.POST("/save", save)
func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}

func main() {
	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middleware
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "you are authorized user!")
	})

	// Route level middleware
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}

	e.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "/users")
	}, track)

	e.GET("/show", show)

	e.POST("/save", save)

	e.GET("/users/:id", getUser)
	e.POST("/users", saveUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
