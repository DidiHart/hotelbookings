package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DidiHart/hotelbookings/internal/config"
	"github.com/DidiHart/hotelbookings/internal/driver"
	"github.com/DidiHart/hotelbookings/internal/handlers"
	"github.com/DidiHart/hotelbookings/internal/helpers"
	"github.com/DidiHart/hotelbookings/internal/models"
	"github.com/DidiHart/hotelbookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var logInfo *log.Logger
var errorInfo *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("starting mail listener...")

	listenForMail()

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("starting application on port %s", portNumber)

	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change to true in production
	app.Inproduction = false

	logInfo = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.LogInfo = logInfo

	errorInfo = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorInfo = errorInfo

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Inproduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("")

	if err != nil {
		log.Fatal("Cannot connect to database! Disconnected...")
	}

	log.Println("Connected to database!")

	tc, err := render.CacheTemplate()

	if err != nil {
		log.Fatal("cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false //dev mode with UseCache is false

	repo := handlers.NewRepo(&app, db)
	handlers.Newhandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
