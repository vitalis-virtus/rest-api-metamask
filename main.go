package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/vitalis-virtus/rest-api-metamask/model"
	"github.com/vitalis-virtus/rest-api-metamask/storage"

	"github.com/go-chi/chi"
)

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
var ErrInvalidAddress = errors.New("invalid address")

type RegisterPayload struct {
	Address string `json:"address"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.Address) {
		return ErrInvalidAddress
	}
	return nil
}

func registerHandler(s *storage.MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p RegisterPayload
		if err := bindReqBody(r, &p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := p.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := model.User{
			Address: strings.ToLower(p.Address),
		}

		if err := s.CreateIfNotExists(u); err != nil {
			switch errors.Is(err, storage.ErrUserExists) {
			case true:
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func userNonceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func signInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func welcomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func bindReqBody(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func run() error {
	r := chi.NewRouter()

	st := storage.NewMemoryStorage()

	//	registering handlers
	r.Post("/register", registerHandler(st))
	r.Get("/users/{address:^0x[a-fA-F0-9]{40}$}/nonce", userNonceHandler())
	r.Post("/signin", signInHandler())
	r.Get("/welcome", welcomeHandler())

	err := http.ListenAndServe("localhost:8002", r)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
