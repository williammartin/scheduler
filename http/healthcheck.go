package http

import "net/http"

func (s Server) HandleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
