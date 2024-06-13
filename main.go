package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var movies []Movie

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env failed")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	movies = append(movies, Movie{ID: 1, Title: "Inception", Director: "Christopher Nolan", Year: 2010})
	movies = append(movies, Movie{ID: 2, Title: "BlacKkKlansman", Director: "Spike Lee", Year: 2018})

	app.Post("/login", login)
	app.Get("/config", getEnv)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Use(logging)

	app.Use(checkMiddleware)

	app.Get("/movies", getMovies)
	app.Get("/movies/:id", getMovie)
	app.Post("/movies", createMovie)
	app.Put("/movies/:id", updateMovie)
	app.Delete("/movies/:id", deleteMovie)
	app.Post("/poster", uploadFile)
	app.Get("/test-html", testHTML)

	app.Listen(":8080")
}

func getMovies(c *fiber.Ctx) error {
	return c.JSON(movies)
}

func getMovie(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, movie := range movies {
		if movie.ID == movieId {
			return c.JSON(movie)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func createMovie(c *fiber.Ctx) error {
	movie := new(Movie)

	if err := c.BodyParser(movie); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	movies = append(movies, *movie)
	return c.JSON(movie)
}

func updateMovie(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	movieUpdate := new(Movie)

	if err := c.BodyParser(movieUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for index, movie := range movies {
		if movie.ID == movieId {
			movies[index].Director = movieUpdate.Director
			movies[index].Title = movieUpdate.Title
			movies[index].Year = movieUpdate.Year
			return c.JSON(movies[index])
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func deleteMovie(c *fiber.Ctx) error {
	movieId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, movie := range movies {
		if movie.ID == movieId {
			movies = append(movies[:i], movies[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("poster")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/"+file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File Upload Complete!!")
}

func testHTML(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":  "Hello, World!!",
		"Auther": "Book",
	})
}

func getEnv(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"JWT_SECRET": os.Getenv("JWT_SECRET"),
	})
}

func logging(c *fiber.Ctx) error {
	start := time.Now()

	fmt.Printf("URL = %s, METHOD = %s, TIME = %s\n", c.OriginalURL(), c.Method(), start)

	return c.Next()
}

var member = User{
	Email:    "hello@example.com",
	Password: "P@ssw0rd",
}

func login(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != member.Email || user.Password != member.Password {
		return fiber.ErrUnauthorized
	}

	// Using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Fetching JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("JWT secret not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Login Successfully !!",
		"token":   t,
		"err":     err,
	})
}

func checkMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "member" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}
