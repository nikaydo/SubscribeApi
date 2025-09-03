package handlers

import (
	"encoding/json"
	"main/internal/config"
	"main/internal/database"
	"main/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Handlers struct {
	Database database.DatabaseInterface
	Env      config.Env
}

type Error struct {
	Error string `json:"error"`
}

func HandlersInit(db database.DatabaseInterface, env config.Env, mux *http.ServeMux) Handlers {
	Handler := Handlers{Database: db, Env: env}
	mux.HandleFunc("/api/subscribe", Handler.Subscribe)
	return Handler
}

func (h *Handlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	switch r.Method {
	case http.MethodGet:
		filter, err := models.ParseFilters(r.Body)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		subs, err := h.Database.GetSubscribe(filter)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		b, err := models.ParseJson(subs)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if filter.Summary {
			var sum models.Summary
			for _, i := range subs {
				sum.Symmary += i.Price
			}
			WriteResponse(400, sum, w)
			return
		}

		w.WriteHeader(200)
		w.Write(b)
		return
	case http.MethodDelete:
		idRaw := r.FormValue("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		if err := h.Database.DeleteSubscribe(id); err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		w.WriteHeader(200)
		return
	case http.MethodPost:
		sub, err := models.ParseResponse(r.Body)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		if err := uuid.Validate(sub.UserId.String()); err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		id, err := h.Database.Subscribe(sub)
		if err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		WriteResponse(201, id, w)
		return
	case http.MethodPut:
		var sub models.Subscribe
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		if err := uuid.Validate(sub.UserId.String()); err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		t, err := time.Parse("02.01.2006", "01."+sub.StartDate)
		if err != nil {
			WriteResponse(400, Error{Error: "use date format XX.XXXX"}, w)
			return
		}
		if t.Before(time.Now()) {
			WriteResponse(400, Error{Error: "date too late"}, w)
			return
		}
		if err := h.Database.UpdateSubscribe(sub); err != nil {
			WriteResponse(400, Error{Error: err.Error()}, w)
			return
		}
		w.WriteHeader(200)
	default:
		WriteResponse(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"}, w)
	}
}

func WriteResponse(code int, j any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	byteJson, err := json.Marshal(j)
	if err != nil {
		WriteResponse(500, Error{Error: "cant marshal response"}, w)
		return
	}
	w.WriteHeader(code)
	w.Write(byteJson)
}

//60601fee-2bf1-4721-ae6f-7636e79a0cba
