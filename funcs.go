package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"
)

// This functions are required in template FuncMap
// You have to use them for Template builders (or mix with your own)
func Funcs() template.FuncMap {
	return template.FuncMap{
		"meta": func(p Page) template.HTML {
			builder := ""
			meta := Meta{}
			if p, ok := p.(ImplementsMeta); ok {
				meta = p.Meta()
			}
			if meta.Title != "" {
				builder += "<title>" + meta.Title + "</title>\n"
			}
			if meta.Canonical != "" {
				builder += `<link rel="canonical" href="` + meta.Canonical + `">` + "\n"
			}
			if len(meta.Hreflangs) != 0 {
				for _, hreflang := range meta.Hreflangs {
					builder += `<link rel="alternate" hreflang="` + hreflang.Lang + `" href="` + hreflang.Href + `">` + "\n"
				}
			}
			return template.HTML(builder)
		},
		"dynamics": func() template.HTML {
			return template.HTML(ssaclient)
		},
		"json": func(data interface{}) string {
			d, _ := json.Marshal(data)
			return string(d)
		},
		"componentattrs": func(c Component) template.HTMLAttr {
			// Extract component data
			name := reflect.ValueOf(c).Elem().Type().Name()
			statebytes, err := json.Marshal(c)
			if err != nil {
				panic(err)
			}
			state := string(statebytes)
			state = url.QueryEscape(state)
			// Build attributes
			builder := fmt.Sprintf(`name='%s' state='%s'`, name, state)
			return template.HTMLAttr(builder)
		},
		"action": func(action string, args ...interface{}) template.JS {
			formattedargs := []string{}
			for _, arg := range args {
				b, _ := json.Marshal(arg)
				formattedargs = append(formattedargs, string(b))
			}

			return template.JS(fmt.Sprintf("Action(this, '%s', %s)", action, strings.Join(formattedargs, ", ")))
		},
		"bind": func(field string) template.JS {
			return template.JS(fmt.Sprintf("Bind(this, '%s')", field))
		},
		"formsubmit": func() template.JS {
			return template.JS("FormSubmit(this, event)")
		},
	}
}