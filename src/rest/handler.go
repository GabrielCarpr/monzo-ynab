package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"monzo-ynab/commands"
	"monzo-ynab/internal/config"
	"monzo-ynab/monzo"
	"net/http"
)

// NewHandler returns a HTTP hander.
func NewHandler(c config.Config, cmds *commands.Commands) *Handler {
	mux := http.NewServeMux()
	handler := &Handler{c, cmds, mux}

	mux.HandleFunc("/events/monzo/", handler.transaction)

	return handler
}

// Handler is a HTTP handler for the app.
type Handler struct {
	config   config.Config
	commands *commands.Commands
	mux      *http.ServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// Transaction is the endpoint for receiving Monzo transaction events.
func (h Handler) transaction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var event TransactionEvent
	err = json.Unmarshal(bod, &event)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if event.Type != "transaction.created" {
		log.Printf("Event was unknown type: %s", event.Type)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction := event.Data.Transaction()
	err = h.commands.Store.Execute(transaction)
	if err != nil {
		log.Printf("Failed storing: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TransactionEvent is a Monzo transaction event.
type TransactionEvent struct {
	Type string            `json:"type"`
	Data monzo.Transaction `json:"data"`
}
