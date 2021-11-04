package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// exports to all parts of our application, but doesn't import anything from anywhere else
// only uses packages already built into our standard library
// because it is a struct, we can put anything we need sitewide for our app, and it will be available to every package that imports this package
type AppConfig struct {
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}

// The config is just a way to share information among the parts of our application that are going to need it. For example, certain parts of the application
// need to know how to connect to the database, or whether or not we are in production, or whatever

// The config struct is used to "inject" this data into the parts of our application which need to have access to it. It's nothing more than a means of sharing
// data around to the various parts of the application
