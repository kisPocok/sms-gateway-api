package queue

import "net/http"

type middleware func(http.HandlerFunc) http.HandlerFunc

// Q is a http middleware queue
// Yeah, I know... This worls like Alice, however it's written by me!
type Q struct {
	queue []middleware
}

// Create queue for middlewares
func Create(fn ...middleware) *Q {
	q := Q{}
	for _, f := range fn {
		q.queue = append(q.queue, f)
	}
	return &q
}

// Then add the final http.HandlerFunc and execute the queue
func (q *Q) Then(final http.HandlerFunc) http.HandlerFunc {
	for i := len(q.queue) - 1; i >= 0; i-- {
		final = q.queue[i](final)
	}
	return final
}
