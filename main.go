package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

var movies []Movie

func main() {
	app := fiber.New()

	movies = append(movies, Movie{ID: 1, Title: "Inception", Director: "Christopher Nolan", Year: 2010})
	movies = append(movies, Movie{ID: 2, Title: "BlacKkKlansman", Director: "Spike Lee", Year: 2018})

	app.Get("/movies", getMovies)
	app.Get("/movies/:id", getMovie)
	app.Post("/movies", createMovie)
	app.Put("/movies/:id", updateMovie)
	app.Delete("/movies/:id", deleteMovie)

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
