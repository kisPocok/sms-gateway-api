package main

import (
	"encoding/json"
	"errors"
	"msgbird/logger"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	testKey = "0E8kldYVc5JwvaeXF2j0ew0ty"

	port            = 9000
	testPath        = "/api/v1/message"
	NL              = "\n"
	_missingParam   = ""
	badRequestError = "bad request"
	longText        = "Bacon ipsum dolor amet kielbasa tenderloin brisket short ribs tri-tip venison turkey ground round sausage corned beef hamburger pork chop turducken jerky. Drumstick doner landjaeger, frankfurter flank meatball strip steak pig jowl sirloin cow ham hock tenderloin venison. Turkey capicola salami, bresaola biltong meatloaf andouille shankle cupim. Pastrami ribeye meatball, strip steak pig picanha shankle ham hock drumstick pancetta jowl kielbasa. Turducken andouille doner, salami kielbasa meatloaf strip steak biltong pork belly alcatra."

	sampleRecipient  = "3631203040"
	sampleOriginator = "David"
	sampleMessage    = "Hello World!"
)

func createTestServer() *Server {
	return &Server{
		port:    port,
		version: "v1.0-beta",
		console: logger.Null{},
		msgbird: MockClient{},
	}
}

func TestGetShouldRespond405(t *testing.T) {
	srv := createTestServer()
	r := callHTTPServer(t, http.MethodGet, nil, srv.messageHandler)
	expectNotAllowedMethod(t, r.Code)
	expectNoNewMessageCall(t)
}

func TestWithoutDataShouldRespondBadRequest(t *testing.T) {
	srv := createTestServer()
	r := callHTTPServer(t, http.MethodPost, nil, srv.messageHandler)
	expectMissingRecipient(t, r)
	expectBadRequest(t, r.Code)
	expectNoNewMessageCall(t)
}

func TestMissingOriginatorShouldBeWarn(t *testing.T) {
	srv := createTestServer()
	data := sampleData(sampleRecipient, _missingParam, sampleMessage)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expectMissingOriginator(t, r)
	expectBadRequest(t, r.Code)
	expectNoNewMessageCall(t)
}

func TestMissingMessageShouldBeWarn(t *testing.T) {
	srv := createTestServer()
	data := sampleData(sampleRecipient, sampleOriginator, _missingParam)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expectMissingMessage(t, r)
	expectBadRequest(t, r.Code)
	expectNoNewMessageCall(t)
}

func TestSendingEveryDataShouldFine(t *testing.T) {
	// Skipped because my account only support my phone number
	t.Skip()

	srv := createTestServer()
	data := sampleData(sampleRecipient, sampleOriginator, sampleMessage)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expect200OK(t, r.Code)
}

func TestUnsentMailErrorHandler(t *testing.T) {
	clearStage()
	srv := createTestServer()
	data := sampleData(sampleRecipient, sampleOriginator, sampleMessage)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expect200OK(t, r.Code)
	expectNewMessageCall(t, 1)
}
func TestMessagebirdApiResponseBadRequest(t *testing.T) {
	clearStage()
	srv := createTestServer()
	mockClientError = errors.New(badRequestError)
	data := sampleData(sampleRecipient, sampleOriginator, sampleMessage)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expectBadRequest(t, r.Code)
	expectNewMessageCall(t, 1)
}

func TestMessagebirdApiResponseUnknownError(t *testing.T) {
	clearStage()
	srv := createTestServer()
	mockClientError = errors.New("unknown error")
	data := sampleData(sampleRecipient, sampleOriginator, sampleMessage)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expectInternalError(t, r.Code)
	expectNewMessageCall(t, 1)
}

func TestLongMessagesShouldBeSentSeparately(t *testing.T) {
	clearStage()
	srv := createTestServer()
	data := sampleData(sampleRecipient, sampleOriginator, longText)
	r := callHTTPServer(t, http.MethodPost, data, srv.messageHandler)
	expect200OK(t, r.Code)
	expectNewMessageCall(t, 4)
}
func TestTheLastFunctionToRiseCoverage100Percent(t *testing.T) {
	p := getPort()
	if p != 8080 {
		t.Errorf("Expected port number is %v, got %v.", 8080, p)
	}
}

type MockClient struct{}

func (mock MockClient) NewMessage(originator string, recipients []string, body string, msgParams *messagebird.MessageParams) (*messagebird.Message, error) {
	// count occurence
	mockClientCall++

	var t time.Time
	t = time.Now()
	msg := &messagebird.Message{CreatedDatetime: &t}

	if mockClientError != nil {
		code := 0
		if mockClientError.Error() == badRequestError {
			code = 9
		}
		// if error exits add some custom error here too
		e := messagebird.Error{
			Code:        code,
			Description: "Dummy description",
			Parameter:   "Dummy param are missing",
		}
		msg.Errors = append(make([]messagebird.Error, 0), e)
	}

	return msg, mockClientError
}

// I hate global variables, but this was the easiest way
var mockClientCall = 0
var mockClientError error

func clearStage() {
	mockClientCall = 0
	mockClientError = nil
}

func expectNoNewMessageCall(t *testing.T) {
	expectNewMessageCall(t, 0)
}

func expectNewMessageCall(t *testing.T, expected int) {
	if expected != mockClientCall {
		t.Errorf("Expect NewMessage() call exactly %v times. Called %v times.", expected, mockClientCall)
	}
}

func expectMissingRecipient(t *testing.T, r *httptest.ResponseRecorder) {
	expectMissingParam(t, r, "Invalid recipient number")
}

func expectMissingOriginator(t *testing.T, r *httptest.ResponseRecorder) {
	expectMissingParam(t, r, "Originator is missing")
}

func expectMissingMessage(t *testing.T, r *httptest.ResponseRecorder) {
	expectMissingParam(t, r, "Message is empty")
}

func expectMissingParam(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	s := r.Body.String()
	var data ResponseRoot
	json.Unmarshal([]byte(s), &data)
	if actual := data.Errors[0].Title; actual != expected {
		t.Errorf("Wrong response. Expected `%s`, got `%s.`", expected, actual)
	}
}

func expect200OK(t *testing.T, actual int) {
	expectStatusCode(t, http.StatusOK, actual)
}

func expectInternalError(t *testing.T, actual int) {
	expectStatusCode(t, http.StatusInternalServerError, actual)
}

func expectBadRequest(t *testing.T, actual int) {
	expectStatusCode(t, http.StatusBadRequest, actual)
}

func expectNotAllowedMethod(t *testing.T, actual int) {
	expectStatusCode(t, http.StatusMethodNotAllowed, actual)
}

func expectStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Handler returned wrong status code: expected %v, got %v", expected, actual)
	}
}

func sampleData(recipient, originator, message string) url.Values {
	data := url.Values{}
	data.Add("recipient", recipient)
	data.Add("originator", originator)
	data.Add("message", message)
	return data
}

func callHTTPServer(t *testing.T, method string, params url.Values, fn http.HandlerFunc) *httptest.ResponseRecorder {
	request, err := http.NewRequest(method, testPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.PostForm = params // pass post data
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(response, request)
	return response
}

func TestRootShouldRespondDummyData(t *testing.T) {
	srv := createTestServer()
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.rootHandler)
	handler.ServeHTTP(response, request)
	if response.Body.String() != title {
		t.Error("Root response failed.")
	}
}

func TestHTTPThroughputIntegration(t *testing.T) {
	go serv(10011)

	// TODO http server
	rate := uint64(2) // per second
	duration := 2 * time.Second

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    "http://localhost:10011" + testPath,
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	if metrics.StatusCodes["405"] != 4 {
		t.Error("Throughput test failed. Expect 9 `bad request` calls.")
	}

	if metrics.Latencies.Max < duration {
		t.Error("Throughput test failed. Expect greater than", duration, "latency, got", metrics.Latencies.Max.Seconds(), "s.")
	}
}
