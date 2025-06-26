package api

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(h Handler) Handler
