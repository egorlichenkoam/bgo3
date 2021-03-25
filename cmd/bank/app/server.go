package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const name = "server"

//bank server structure
type Server struct {
	name      string
	personSvc *person.Service
	cardSvc   *card.Service
	txSvc     *transaction.Service
	mux       *http.ServeMux
}

//server error
type ServerError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

//create and return server
func NewServer(personSvc *person.Service, cardSvc *card.Service, txSvc *transaction.Service, mux *http.ServeMux) *Server {
	log.Printf("%s - %s", name, "Create server")
	return &Server{
		name:      "Bank server",
		personSvc: personSvc,
		cardSvc:   cardSvc,
		txSvc:     txSvc,
		mux:       mux,
	}
}

//server http
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("%s - %s", name, "ServerHTTP")
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/addCard", s.addCard)
	s.mux.HandleFunc("/getCards", s.getCards)
}

func (s *Server) isPost(w http.ResponseWriter, method string) bool {
	log.Printf("%s - %s", s.name, "Start isPost")
	log.Printf("%s - method : %s", s.name, method)
	if method != "POST" {
		s.error(w, errors.New("support only post methods"), http.StatusMethodNotAllowed)
		return false
	}
	log.Printf("%s - %s", s.name, "End isPost")
	return true
}

func (s *Server) readFormValueInt64(key string, w http.ResponseWriter, r *http.Request) (int64, error) {
	log.Printf("%s - %s : %s", s.name, "Start readFormValueInt64", key)
	val, err := strconv.ParseInt(r.FormValue(key), 10, 64)
	if err != nil {
		s.error(w, errors.New(fmt.Sprintf("%s : %s", key, err.Error())), 500)
		return val, err
	}
	log.Printf("%s - %s : %s", s.name, "End readFormValueInt64", key)
	return val, nil
}

func (s *Server) readFormValueStr(key string, w http.ResponseWriter, r *http.Request) (string, error) {
	log.Printf("%s - %s : %s", s.name, "Start readFormValueStr", key)
	val := r.FormValue(key)
	if len(val) == 0 {
		err := errors.New(fmt.Sprintf("%s : %s", key, "can't be empty"))
		s.error(w, err, 517)
		return val, err
	}
	log.Printf("%s - %s : %s", s.name, "End readFormValueStr", key)
	return val, nil
}

//add card to particular person
func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", s.name, "Start addCard")
	if !s.isPost(w, r.Method) {
		return
	}
	personId, err := s.readFormValueInt64("personId", w, r)
	if err != nil {
		s.error(w, errors.New(fmt.Sprintf("%s", err.Error())), 500)
		return
	}
	cardTypeInt64, err := s.readFormValueInt64("cardType", w, r)
	if err != nil {
		s.error(w, errors.New(fmt.Sprintf("%s", err.Error())), 500)
		return
	}
	if cardTypeInt64 != 0 && cardTypeInt64 != 1 {
		s.error(w, errors.New("card type should be 0 for plastic or 1 for virtual"), 517)
		return
	}
	cardType := "PLASTIC"
	if cardTypeInt64 == 1 {
		cardType = "VIRTUAL"
	}
	issuer, err := s.readFormValueStr("issuer", w, r)
	if err != nil {
		s.error(w, errors.New(fmt.Sprintf("%s", err.Error())), 500)
		return
	}
	if !s.personSvc.Exist(personId) {
		s.error(w, errors.New("person does not exist"), 517)
		return
	}
	if len(s.cardSvc.ByPersonId(personId)) == 0 {
		s.error(w, errors.New("person have no cards, therefore person can't add card"), 517)
		return
	}
	newCard := s.cardSvc.Create(issuer, personId, 100_00, card.Rub, strconv.FormatInt(rand.Int63(), 10), cardType)
	responseBody, err := json.Marshal(newCard.DTO())
	if err != nil {
		s.error(w, err, 500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		s.error(w, err, 500)
	}
	log.Printf("%s - %s", s.name, "End addCard")
}

//return cards that belongs to particular person
func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s - %s", s.name, "Start getCards")
	if !s.isPost(w, r.Method) {
		return
	}
	personId, err := s.readFormValueInt64("personId", w, r)
	if err != nil {
		s.error(w, errors.New(fmt.Sprintf("%s", err.Error())), 500)
		return
	}
	if !s.personSvc.Exist(personId) {
		s.error(w, errors.New("person does not exist"), 517)
		return
	}
	cards := s.cardSvc.ByPersonId(personId)
	if len(cards) == 0 {
		s.error(w, errors.New(fmt.Sprintf("person %d have no cards", personId)), 517)
		return
	}
	responseBody, err := json.Marshal(cards)
	if err != nil {
		s.error(w, err, 517)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		s.error(w, err, 500)
	}
}

//return error from server
func (s *Server) error(w http.ResponseWriter, e error, code int) {
	log.Printf("%s - %s : %s %d", s.name, "send error", e, code)
	responseError, err := json.Marshal(ServerError{
		Error: e.Error(),
		Code:  code,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseError)))
	w.WriteHeader(code)
	w.Write(responseError)
}
