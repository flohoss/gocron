package router

import (
	"errors"
	"html/template"
	"io"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return template.ExecuteTemplate(w, "layout.html", data)
}

func templateString(files []string) []string {
	templatesFolder := "web/templates/"
	baseTemplates := []string{templatesFolder + "layout.html", templatesFolder + "icons.html"}
	for i := 0; i < len(files); i++ {
		files[i] = templatesFolder + files[i]
	}
	combined := append(baseTemplates, files[:]...)
	return combined
}

func generateTemplate(files ...string) *template.Template {
	return template.Must(template.New("").Funcs(sprig.FuncMap()).ParseFiles(templateString(files)...))
}

func initTemplates() *Template {
	templates := make(map[string]*template.Template)

	templates["system"] = generateTemplate("pages/system.html")
	templates["logs"] = generateTemplate("pages/logs.html")
	templates["tools"] = generateTemplate("pages/tools.html")

	templates["jobs"] = generateTemplate("pages/jobs.html")
	templates["jobsForm"] = generateTemplate("forms.html", "pages/jobForm.html")

	templates["remotes"] = generateTemplate("pages/remotes.html")
	templates["remotesForm"] = generateTemplate("forms.html", "pages/remoteForm.html")

	return &Template{templates: templates}
}
