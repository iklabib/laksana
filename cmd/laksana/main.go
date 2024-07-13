package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/iklabib/laksana/model"
	"codeberg.org/iklabib/laksana/toolchains"
)

func main() {
	BaseUrl := os.Getenv("BASE_URL")

	workdir := "/tmp/laksana"
	if err := os.Mkdir(workdir, 0o775); err != nil {
		if !os.IsExist(err) {
			log.Panicln("failed to create workdir")
		}
	}

	evaluator := toolchains.NewEvaluator(workdir)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /run", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var submission model.Submission
		if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		eval := evaluator.Eval(r.Context(), submission)

		result, err := json.Marshal(eval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
	})

	fmt.Printf("Serving at %s \n", BaseUrl)
	err := http.ListenAndServe(BaseUrl, mux)
	log.Fatal(err)
}
