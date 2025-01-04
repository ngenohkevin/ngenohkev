package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ngenohkevin/ngenohkev/internals/components"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//create a new router
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := components.Home("NgenohKev Blog").Render(r.Context(), w)
		if err != nil {
			logger.Error("An error occurred while rendering the home page", err)
		}
	})

	router.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {

	})

	//start the server
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Info("Server is running on port 8080")

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error("An error occurred while starting the server", err)
	}
}
