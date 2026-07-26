package handlers

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v5"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func initTemplates() *Template {
	return &Template{
		templates: parseIndexTemplate(),
	}
}
