package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type request struct {
	Address string `json:"address"`
}

type server struct {
	faucet *faucet
}

func NewServer(faucet *faucet) *server {
	return &server{faucet: faucet}
}

func (s server) Run(port int) {
	r := mux.NewRouter()
	r.Methods("POST").Path("/api/faucet").HandlerFunc(s.faucetHandler)
	r.Methods("GET").Path("/api/info").HandlerFunc(s.infoHandler)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.Use(negroni.NewStatic(http.Dir("web/public")))
	n.UseHandler(r)

	log.Infof("Starting http server %d", port)
	http.ListenAndServe(":"+strconv.Itoa(port), n)
}

func (s server) faucetHandler(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if len(req.Address) == 0 || !re.MatchString(req.Address) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}

	if !s.faucet.isEmptyQueue() {
		if s.faucet.tryEnqueue(req.Address) {
			log.WithFields(log.Fields{
				"address": req.Address,
			}).Info("Added to queue successfully")
			fmt.Fprintf(w, "Added %s to the queue", req.Address)
		} else {
			log.Warn("Max queue capacity reached")
			http.Error(w, "Faucet queue is too long, please try again later.", http.StatusServiceUnavailable)
		}
		return
	}

	txHash, err := s.faucet.Transfer(context.Background(), req.Address, s.faucet.GetPayoutWei())
	if err != nil {
		log.WithError(err).Error("Could not send transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"txHash":  txHash,
		"address": req.Address,
	}).Info("Funded directly successfully")
	fmt.Fprintf(w, txHash)
}

func (s server) infoHandler(w http.ResponseWriter, r *http.Request) {
	type Info struct {
		Account string `json:"account"`
		Payout  string `json:"payout"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Info{
		Account: s.faucet.FromAddress().String(),
		Payout:  s.faucet.GetPayoutWei().String(),
	})
}
