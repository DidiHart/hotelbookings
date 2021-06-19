package config

import (
	"html/template"
	"log"

	"github.com/DidiHart/hotelbookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

//AppConfig stores the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	LogInfo       *log.Logger
	ErrorInfo     *log.Logger
	Inproduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
