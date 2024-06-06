package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) postMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	m := &Message{
		Content: input.Content,
	}

	if err := app.Insert(m); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/messages?uuid=%s", m.UUID.String()))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"message": m}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) retrieveAndRemoveMessageHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	m := &Message{UUID: uuid}

	if err := app.RetrieveAndRemove(m); err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": m}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type metricsResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}

func (app *application) metrics(next http.Handler) http.Handler {
	var totalRequestsReceived = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_requests_received",
			Help: "number of http requests received",
		},
	)

	var totalResponsesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_responses_sent",
			Help: "number of http responses sent",
		},
	)

	var totalResponsesSentByStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_responses_sent_by_status",
			Help: "http responses sent by status",
		},
		[]string{
			"response_code",
		},
	)

	var requestersNumOfRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requesters_number_of_request",
			Help: "number of requests for each requester",
		},
		[]string{
			"host",
		},
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Inc()

		mw := newMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		totalResponsesSent.Inc()

		totalResponsesSentByStatus.With(prometheus.Labels{
			"response_code": strconv.Itoa(mw.statusCode),
		}).Inc()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "parse error"
		}

		requestersNumOfRequests.With(prometheus.Labels{
			"host": ip,
		}).Inc()
	})
}
