package main

// ResponseRoot describes jsonapi.org requirements
type ResponseRoot struct {
	Data   []MessageResponse `json:"data"`
	Errors []ErrorResponse   `json:"errors,omitempty"`
}

// MessageResponse store all the important (is it true?) return values.
type MessageResponse struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	HRef    string `json:"href"`
	Body    string `json:"body"`
	Created int64  `json:"created_at"`
}

// ErrorResponse for error handling
type ErrorResponse struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}
