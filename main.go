package main

import (
	"REST-API-using-GO-Fiber/models"
	"REST-API-using-GO-Fiber/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Blog struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBlogs(context *fiber.Ctx) error {
	blog := Blog{}

	err := context.BodyParser(&blog)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&blog).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create a blog"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "blog created successfully"})

	return nil
}

func (r *Repository) GetBlogs(context *fiber.Ctx) error {
	blogModels := &[]models.Blog{}

	err := r.DB.Find(blogModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the blog"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "blog fetched successfully", "data": blogModels})
	return nil
}

func (r *Repository) GetBlogsByID(context *fiber.Ctx) error {
	blogModel := &models.Blog{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be found",
		})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(blogModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the blog"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "blog fetched successfully", "data": blogModel})
	return nil
}

func (r *Repository) DeleteBlogs(context *fiber.Ctx) error {
	blogModel := models.Blog{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be found",
		})
		return nil
	}

	err := r.DB.Delete(&blogModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete the blog"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "blog deleted successfully"})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_blogs", r.CreateBlogs)
	api.Delete("/delete_blogs/:id", r.DeleteBlogs)
	api.Get("/get_blogs/:id", r.GetBlogsByID)
	api.Get("/blogs", r.GetBlogs)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}

	db, err := storage.NewConnection(&config)

	if err != nil {
		log.Fatal("could not connect to the db")
	}

	err = models.MigrateBlogs(db)
	if err != nil {
		log.Fatal("could not create a db")
	}
	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
