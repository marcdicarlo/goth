package handlers

import (
	"net/http"

	"goth/internal/components"


)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	component := components.Hello()
	component.Render(r.Context(), w)
}
