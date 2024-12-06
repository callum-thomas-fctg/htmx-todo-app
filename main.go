package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Todo struct {
	Completed bool
	Text      string
	Id        string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var todos []Todo

func findTodo(id string) (int, error) {
	todoId := -1
	for idx, t := range todos {
		if t.Id == id {
			todoId = idx
			break
		}
	}

	if todoId == -1 {
		return todoId, errors.New("failed to find todo")
	}

	return todoId, nil
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	if len(todos) == 0 {
		// Seed some todos.
		todos = append(todos, Todo{
			Completed: false,
			Text:      "Write a todo app",
			Id:        uuid.New().String(),
		})
		todos = append(todos, Todo{
			Completed: true,
			Text:      "Make some lunch",
			Id:        uuid.New().String(),
		})
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = t

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		type IndexData struct {
			Todos []Todo
		}

		templ, err := template.ParseFiles("templates/base.html", "templates/index.html")
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		data := IndexData{
			Todos: todos,
		}

		return templ.Execute(c.Response().Writer, data)
	})

	e.POST("/todo", func(c echo.Context) error {
		text := c.FormValue("text")
		newId := uuid.New()
		newTodo := Todo{
			Text:      text,
			Completed: false,
			Id:        newId.String(),
		}
		todos = append(todos, newTodo)

		templ, err := template.ParseFiles("templates/todo.html")
		if err != nil {
			c.Logger().Error("failed to parse template")
			return c.NoContent(http.StatusInternalServerError)
		}

		return templ.Execute(c.Response().Writer, newTodo)
	})

	e.DELETE("/todo/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		todoId, err := findTodo(id)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		todos = slices.Delete(todos, todoId, todoId)
		return c.HTML(http.StatusOK, "")
	})

	e.POST("/todo/:id/toggle", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			c.Logger().Error("failed to get id param")
			return c.NoContent(http.StatusNotFound)
		}

		idx, err := findTodo(id)
		if err != nil {
			c.Logger().Error("failed to find todo")
			return c.NoContent(http.StatusNotFound)
		}

		todo := &todos[idx]
		todo.Completed = !todo.Completed

		templ, err := template.ParseFiles("templates/todo.html")
		if err != nil {
			c.Logger().Error("failed to render todo template.")
			return c.NoContent(http.StatusInternalServerError)
		}

		return templ.Execute(c.Response().Writer, todo)
	})

	e.Start(":8080")
}
