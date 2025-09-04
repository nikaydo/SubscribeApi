package handlers

import (
	"encoding/json"
	"main/internal/config"
	"main/internal/database"
	"main/internal/models"
	"net/http"

	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return Handler
}

func (h *Handlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
		return
	case http.MethodDelete:
		h.Delete(w, r)
		return
	case http.MethodPost:
		h.Post(w, r)
		return
	case http.MethodPut:
		h.Put(w, r)
		return
	default:
		WriteResponse(http.StatusMethodNotAllowed, Error{Error: "Method not allowed"}, w)
	}
}

// Get godoc
// @Summary      Get subscriptions
// @Description  Retrieves subscriptions based on filters
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query  string  false  "User ID"
// @Param        service_name query  string  false  "Service name"
// @Param        start_date   query  string  false  "Start date"
// @Param        end_date     query  string  false  "End date"
// @Param        summary      query  bool    false  "Include summary"
// @Success      200     {array}   models.Subscribe
// @Success      200     {object}  models.Summary  "If summary=true"
// @Failure      400     {object}  Error
// @Failure      500     {object}  Error
// @Router       /api/subscribe [get]
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {

	filter := models.Filters{
		ServiceName: r.FormValue("service_name"),
		StartDate:   r.FormValue("start_date"),
		EndDate:     r.FormValue("end_date"),
		Summary:     r.URL.Query().Get("summary") == "false",
	}
	if err := uuid.Validate(r.FormValue("user_id")); r.FormValue("user_id") != "" && err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	if r.FormValue("user_id") != "" {
		filter.UserId = uuid.MustParse(r.FormValue("user_id"))
	}
	filter, err := filter.ParseFilters()
	if err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	subs, err := h.Database.GetSubscribe(filter)
	if err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	b, err := models.ParseJson(subs)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, Error{Error: err.Error()}, w)
		return
	}
	if filter.Summary {
		var sum models.Summary
		for _, i := range subs {
			sum.Symmary += i.Price
		}
		WriteResponse(http.StatusOK, sum, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Delete godoc
// @Summary      Delete a subscription
// @Description  Deletes a subscription by its ID.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body  body      models.Id  true  "Subscription ID"
// @Success      200   {string}  string                      "OK"
// @Failure      400   {object}  Error                       "Bad Request"
// @Router       /api/subscribe [delete]
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	var subscribe models.Id
	if err := json.NewDecoder(r.Body).Decode(&subscribe); err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	if err := h.Database.DeleteSubscribe(subscribe.Id); err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Post handles HTTP POST requests to create a new subscription.
// @Summary      Create subscription
// @Description  Parses the request body to create a new subscription, validates the user ID, and stores the subscription in the database.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.Subscribe  true  "Subscription data"
// @Success      201  {object}  models.Id
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/subscribe [post]
func (h *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	sub, err := models.ParseResponse(r.Body)
	if err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	if err := uuid.Validate(sub.UserId.String()); err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	id, err := h.Database.Subscribe(sub)
	if err != nil {
		WriteResponse(http.StatusInternalServerError, Error{Error: err.Error()}, w)
		return
	}
	WriteResponse(http.StatusCreated, models.Id{Id: uint(id)}, w)
}

// Put handles HTTP PUT requests to update a subscription.
// @Summary      Update subscription
// @Description  Updates an existing subscription with the provided data.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body  body      models.Subscribe  true  "Subscription update payload"
// @Success      200   {object}  nil
// @Failure      400   {object}  Error  "Invalid input or date format"
// @Failure      500   {object}  Error  "Internal server error"
// @Router       /api/subscribe [put]
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	sub, err := models.ParseResponse(r.Body)
	if err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	if sub.Id == 0 {
		WriteResponse(http.StatusInternalServerError, Error{Error: "id cant be empty"}, w)
		return
	}
	if err := uuid.Validate(sub.UserId.String()); err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: err.Error()}, w)
		return
	}
	if err := h.Database.UpdateSubscribe(sub); err != nil {
		WriteResponse(http.StatusInternalServerError, Error{Error: err.Error()}, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func WriteResponse(code int, j any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	byteJson, err := json.Marshal(j)
	if err != nil {
		WriteResponse(http.StatusBadRequest, Error{Error: "cant marshal response"}, w)
		return
	}
	w.WriteHeader(code)
	w.Write(byteJson)
}
