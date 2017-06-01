package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	messagebird "github.com/messagebird/go-rest-api"

	"errors"
	"strconv"

	"github.com/kisPocok/sms-gateway-api/logger"
	Q "github.com/kisPocok/sms-gateway-api/queue"
	"github.com/kisPocok/sms-gateway-api/split"
)

//const apiKey = "0E8kldYVc5JwvaeXF2j0ew0ty" // account 1
const apiKey = "HXxIc0tw0xmCNa6TnQjHuZvXs" // acount 2
const usage = "usage: msgbird --port={port}"
const title = "SMS Gateway API\n"

func main() {
	serv(getPort())
}

// Server store application-level values
type Server struct {
	port    int
	version string
	console logger.Logging
	msgbird msgbirdClient
}

type msgbirdClient interface {
	NewMessage(string, []string, string, *messagebird.MessageParams) (*messagebird.Message, error)
}

func createServer(port int) *Server {
	return &Server{
		port:    port,
		version: "v1.0",
		console: logger.Logger{},
		msgbird: messagebird.New(apiKey),
	}
}

func (srv *Server) getListenerAddr() string {
	return ":" + strconv.Itoa(srv.port)
}

func (srv *Server) rootHandler(w http.ResponseWriter, request *http.Request) {
	srv.console.Log("Root called")
	io.WriteString(w, title)
}

func (srv *Server) messageHandler(w http.ResponseWriter, request *http.Request) {
	srv.console.Log("API Request has arrived")
	w.Header().Set("Content-Type", "application/json")

	if request.Method != http.MethodPost {
		errorNotSupportedMethod(w)
		return
	}

	request.ParseForm()
	msg := &SMS{
		Recipient:  request.Form.Get("recipient"),
		Originator: request.Form.Get("originator"),
		Message:    request.Form.Get("message"),
	}

	if err := msg.Validate(); err != nil {
		errorWriter(w, http.StatusBadRequest, err.Error(), "")
		return
	}

	var apiResponse ResponseRoot
	for _, msg := range srv.sendMessages(w, msg) {
		if len(msg.Errors) > 0 {
			responseError(w, msg)
			return
		}
		apiResponse.Data = append(apiResponse.Data, formatMessageToResponse(msg))
	}
	json.NewEncoder(w).Encode(apiResponse)
}

type messageSent struct {
	Message messagebird.Message
	Error   error
}

func (srv *Server) sendMessages(w http.ResponseWriter, msg *SMS) []messagebird.Message {
	var responses []messagebird.Message
	recipients := append(make([]string, 0), msg.Recipient)

	for _, m := range split.Splitter(msg.Message) {
		// Transaction handling would be nice here.
		// Limiter is missing here. Sorry.
		msg, _ := srv.msgbird.NewMessage(msg.Originator, recipients, m, nil)
		responses = append(responses, *msg)
	}
	return responses
}

func responseError(w http.ResponseWriter, msg messagebird.Message) {
	// It's require more error handling.
	switch msg.Errors[0].Code {
	case 9:
		errorWriter(w, http.StatusBadRequest,
			"Wrong param "+msg.Errors[0].Parameter, msg.Errors[0].Description)
	default:
		internalError(w)
	}
}

func formatMessageToResponse(msg messagebird.Message) MessageResponse {
	return MessageResponse{
		Type:    "message",
		ID:      msg.Id,
		HRef:    msg.HRef,
		Body:    msg.Body,
		Created: msg.CreatedDatetime.Unix(),
	}
}

// SMS example from the docs: {
//  "recipient":31612345678,
//  "originator":"MessageBird",
//  "message":"This is a test message."}
type SMS struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

// Validate their own properties
// (Here comes some example validation, I don't know is it reallistic or not.)
func (sms *SMS) Validate() error {
	if sms.Recipient == "" {
		return errors.New("Invalid recipient number")
	}

	if sms.Originator == "" {
		return errors.New("Originator is missing")
	}

	if sms.Message == "" {
		return errors.New("Message is empty")
	}

	return nil
}

func serv(port int) {
	srv := createServer(port)
	middlewares := Q.Create(rateLimitOnePerSec)

	http.HandleFunc("/", middlewares.Then(srv.rootHandler))
	http.HandleFunc("/api/v1/message", middlewares.Then(srv.messageHandler))

	fmt.Println("Server listening on port nr.", port)
	log.Fatal(http.ListenAndServe(srv.getListenerAddr(), nil))
}

func rateLimitOnePerSec(next http.HandlerFunc) http.HandlerFunc {
	limiter := time.Tick(time.Second * 1)
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		<-limiter
		next.ServeHTTP(w, request)
	})
}

func internalError(w http.ResponseWriter) {
	errorWriter(w, http.StatusInternalServerError, "Internal server error", "")
}

func errorNotSupportedMethod(w http.ResponseWriter) {
	errorWriter(w, http.StatusMethodNotAllowed, "Not supported method", "")
}

func errorWriter(w http.ResponseWriter, status int, title, details string) {
	errors := make([]ErrorResponse, 1)
	errors[0] = ErrorResponse{
		Status: status,
		Code:   strings.Replace(strings.ToUpper(title), " ", "_", -1),
		Title:  title,
		Detail: details,
	}
	e := ResponseRoot{Errors: errors}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(e)
}

func getPort() int {
	var port int
	flag.IntVar(&port, "port", 8080, usage)
	flag.Parse()
	return port
}
