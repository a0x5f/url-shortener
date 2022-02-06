package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	netUrl "net/url"
	"strings"
	"url-shortener/configs"
	"url-shortener/internal/domain/inmemory"
	"url-shortener/internal/domain/postgres"
	"url-shortener/internal/usecases/translator"
	"url-shortener/pkg/errors"
)

type App struct {
	config     *configs.Configuration
	translator *translator.UrlTranslator
}

type Response struct {
	Url string `json:"url"`
}

func New(c *configs.Configuration) *App {
	if c.Server.UsePostgres {
		return newPostgresApp(c)
	} else {
		return newInmemoryApp(c)
	}
}

func newInmemoryApp(c *configs.Configuration) *App {
	db := make(inmemory.InmemoryDb)

	dbService := inmemory.New(&db)
	urlTranslator := translator.New(dbService)

	return &App{
		config:     c,
		translator: urlTranslator,
	}
}

func newPostgresApp(c *configs.Configuration) *App {
	pgConnString := fmt.Sprintf(
		"port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Port,
		c.Postgres.Host,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Dbname)

	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		//log.Fatal("Error opening database")
		log.Fatal(errors.Wrap(err, "can't open database"))
		return nil
	}

	if err = db.Ping(); err != nil {
		db.Close()
		log.Fatal(errors.Wrap(err, "Error opening database"))
		return nil
	}

	dbService := postgres.New(db)
	urlTranslator := translator.New(dbService)

	return &App{
		config:     c,
		translator: urlTranslator,
	}
}

func (h *App) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/short", h.getShortUrl)
	r.Get("/full", h.getFullUrl)

	log.Printf("Listening port %v...", h.config.Server.Port)

	addr := fmt.Sprintf(":%d", h.config.Server.Port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err.Error())
	}
}

func (h *App) getShortUrl(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.FormValue("url")

	if _, err := netUrl.ParseRequestURI(requestUrl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Wrap(err, "Bad request").Error()))
		return
	}

	url, err := h.translator.ShortenUrl(r.Context(), requestUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrap(err, "Internal server error").Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{url})
}

func (h *App) getFullUrl(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Query().Get("url")

	if !strings.HasPrefix(requestUrl, translator.ShortLinkDomain) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	url, err := h.translator.ExtendUrl(r.Context(), requestUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.Wrap(err, "Internal server error").Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{url})
}
