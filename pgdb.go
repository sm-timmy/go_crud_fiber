package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq" // add this
	"log"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

//var db *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)

	//connStr := "postgresql://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Sprint(db)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

type user struct {
	Id   int
	Name string
	Age  int
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var id int
	var name string
	var age int
	var users []user
	//var users = map[user]int{}
	rows, err := db.Query("SELECT * FROM users")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&id, &name, &age)
		users = append(users, user{id, name, age})
	}
	return c.JSON(users)
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newUser := new(user)
	if err := c.BodyParser(newUser); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newUser)
	if newUser.Name != "" || newUser.Age > 0 {
		_, err := db.Exec("INSERT into users VALUES (DEFAULT, $1, $2)", newUser.Name, newUser.Age)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}
	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	newUser := new(user)
	if err := c.BodyParser(newUser); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", newUser.Name, newUser.Age, newUser.Id)
	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	newUser := new(user)
	if err := c.BodyParser(newUser); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	db.Exec("DELETE from users WHERE id=$1", newUser.Id)
	return c.Redirect("/")
}
