// Package handlers is generated automatically by Sunplate toolkit.
// Please, do not edit it manually.
package handlers

import (
	"net/http"

	contr "github.com/anonx/sunplate/controllers/templates"

	"github.com/anonx/sunplate/strconv"
)

// Templates is an insance of tTemplates that is automatically generated from Templates controller
// being found at "github.com/anonx/sunplate/controllers/templates/templates.go",
// and contains methods to be used as handler functions.
//
// Templates is a controller that provides support of HTML result
// rendering to your application.
// Use SetTemplatePaths to register templates and
// call c.RenderTemplate from your action to render some.
var Templates tTemplates

// tTemplates is a type with handler methods of Templates controller.
type tTemplates struct {
}

// New allocates (github.com/anonx/sunplate/controllers/templates).Templates controller,
// then returns it.
func (t tTemplates) New() *contr.Templates {
	c := &contr.Templates{}
	return c
}

// Before calls (github.com/anonx/sunplate/controllers/templates).Templates.Before.
func (t tTemplates) Before(c *contr.Templates, w http.ResponseWriter, r *http.Request) http.Handler {
	// Call magic Before action of (github.com/anonx/sunplate/controllers/templates).Templates.
	if res := c.Before( // "Binding" parameters.
	); res != nil {
		return res
	}
	return nil
}

// After is a dump method that always returns nil.
func (t tTemplates) After(c *contr.Templates, w http.ResponseWriter, r *http.Request) http.Handler {
	return nil
}

// Initially is a method that is started by every handler function at the very beginning
// of their execution phase.
func (t tTemplates) Initially(c *contr.Templates, w http.ResponseWriter, r *http.Request) (finish bool) {
	return
}

// Finally is a method that is started by every handler function at the very end
// of their execution phase no matter what.
func (t tTemplates) Finally(c *contr.Templates, w http.ResponseWriter, r *http.Request) (finish bool) {
	return
}

// RenderTemplate is a handler that was generated automatically.
// It calls Before, After, Finally methods, and RenderTemplate action found at
// github.com/anonx/sunplate/controllers/templates/templates.go
// in appropriate order.
//
// RenderTemplate initializes and returns HTML type that implements Result interface.
func (t tTemplates) RenderTemplate(w http.ResponseWriter, r *http.Request) {
	c := Templates.New()
	defer Templates.Finally(c, w, r)
	if finish := Templates.Initially(c, w, r); finish {
		return
	}
	if res := Templates.Before(c, w, r); res != nil {
		res.ServeHTTP(w, r)
		return
	}
	if res := c.RenderTemplate( // "Binding" parameters.
		strconv.String(r.Form, "templatePath"),
	); res != nil {
		res.ServeHTTP(w, r)
		return
	}
	if res := Templates.After(c, w, r); res != nil {
		res.ServeHTTP(w, r)
	}
}

func init() {
	_ = strconv.MeaningOfLife
}