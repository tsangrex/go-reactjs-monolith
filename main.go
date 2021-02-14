package main

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type User struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}
func main() {
	// Create new Fiber instance
	engine := html.New("./resources/views", ".html")
	engine.Reload(true)
	engine.Debug(true)
	engine.Layout("embed")
	engine.Delims("{{", "}}")
	store := session.New()
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(requestid.New())
	app.Use(compress.New())

	// Create new GET route on path "/"
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("login", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/login", func(ctx *fiber.Ctx) error {
		user := new(User)

		if err := ctx.BodyParser(user); err != nil {
			return err
		}

		log.Println(user.Username)

		sess, err := store.Get(ctx)
		if err != nil {
			panic(err)
		}
		sess.Set("Username", user.Username)
		if err := sess.Save(); err != nil {
			panic(err)
		}
		return ctx.Redirect("home", 302)
	})

	app.Get("/home", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			panic(err)
		}

		if sess.Fresh() == true {
			return ctx.Redirect("/", 302)
		}

		username := sess.Get("Username")
		return ctx.Render("home", fiber.Map{
			"Username": username,
		})
	})

	app.Get("/logout", func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			panic(err)
		}
		// Destry session
		if err := sess.Destroy(); err != nil {
			panic(err)
		}
		return ctx.Redirect("/", 302)
	})

	app.Get("/download", func(ctx *fiber.Ctx) error {
		return ctx.Download("./resources/files/test.png")
	})

	// Start server on http://localhost:3000
	log.Fatal(app.Listen(":3000"))
}
