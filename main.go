package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vitalis-virtus/rest-api-metamask/middleware"
	"github.com/vitalis-virtus/rest-api-metamask/model"
	"github.com/vitalis-virtus/rest-api-metamask/storage"
	"github.com/vitalis-virtus/rest-api-metamask/utils"
	"github.com/vitalis-virtus/rest-api-metamask/utils/fail"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

type RegisterPayload struct {
	Address string `json:"address"`
}

func (p RegisterPayload) Validate() error {
	if !hexRegex.MatchString(p.Address) {
		return fail.ErrInvalidAddress
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

		nonce, err := utils.GetNonce()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		u := model.User{
			Address: strings.ToLower(p.Address),
			Nonce:   nonce,
		}

		if err := s.CreateIfNotExists(u); err != nil {
			switch errors.Is(err, fail.ErrUserExists) {
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

func userNonceHandler(s *storage.MemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		if !hexRegex.MatchString(address) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := s.Get(strings.ToLower(address))
		if err != nil {
			switch errors.Is(err, fail.ErrUserNotExists) {
			case true:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		resp := struct {
			Nonce string
		}{
			Nonce: user.Nonce,
		}
		utils.RenderJson(r, w, http.StatusOK, resp)
	}
}

func signInHandler(s *storage.MemoryStorage, jwtProvider *utils.JwtHmacProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p model.SignInPayload
		if err := bindReqBody(r, &p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := utils.Validate(p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		address := strings.ToLower(p.Address)
		user, err := Authenticate(s, address, p.Nonce, p.Sig)
		switch err {
		case nil:
		case fail.ErrAuthError:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		signedToken, err := jwtProvider.CreateStandard(user.Address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := struct {
			AccessToken string `json:"access"`
		}{AccessToken: signedToken}
		utils.RenderJson(r, w, http.StatusOK, resp)
	}
}

func welcomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromReqContext(r)
		resp := struct {
			Msg string `json:"msg"`
		}{Msg: "Congrats " + user.Address + " you made it!"}
		utils.RenderJson(r, w, http.StatusOK, resp)
	}
}

func bindReqBody(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func getUserFromReqContext(r *http.Request) model.User {
	ctx := r.Context()
	key := ctx.Value("user").(model.User)
	return key
}

func Authenticate(s *storage.MemoryStorage, address string, nonce string, sigHex string) (model.User, error) {
	user, err := s.Get(address)
	if err != nil {
		return user, err
	}

	if user.Nonce != nonce {
		return user, fail.ErrAuthError
	}

	sig := hexutil.MustDecode(sigHex)
	sig[crypto.RecoveryIDOffset] -= 27
	msg := accounts.TextHash([]byte(nonce))
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return user, err
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	if user.Address != strings.ToLower(recoveredAddr.Hex()) {
		return user, fail.ErrAuthError
	}

	// update the nonce here so that the signature cannot be resused
	nonce, err = utils.GetNonce()
	if err != nil {
		return user, err
	}
	user.Nonce = nonce
	s.Update(user)

	return user, nil
}

func run() error {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)

	st := storage.NewMemoryStorage()
	jwtProv := utils.NewJwtHmacProvider("read something from env here maybe",
		"awesome-metamask-login",
		time.Minute*15)

	//	registering handlers
	r.Post("/register", registerHandler(st))
	r.Get("/users/{address:^0x[a-fA-F0-9]{40}$}/nonce", userNonceHandler(st))
	r.Post("/signin", signInHandler(st, jwtProv))
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(st, jwtProv))
		r.Get("/welcome", welcomeHandler())
	})

	err := http.ListenAndServe("localhost:8002", r)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
