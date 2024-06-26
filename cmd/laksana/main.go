package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"

    "codeberg.org/iklabib/laksana/model"
    "codeberg.org/iklabib/laksana/toolchains"
    "codeberg.org/iklabib/laksana/util"
)

func main() {
    BaseUrl := os.Getenv("BASE_URL")

    workdir := "/tmp/laksana"
    if err := os.Mkdir(workdir, 0o775); err != nil {
        log.Panicln(fmt.Errorf("failed to create workdir"))
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != "POST" {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        var submission model.Submission
        if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        ctx := r.Context()
        evaluator := toolchains.NewEvaluator(workdir)
        resultChan := make(chan model.RunResult)

        go func() {
            defer close(resultChan)
            resultChan <- evaluator.Submission(ctx, submission)
        }()

        select {
        case <-ctx.Done():
            var res model.RunResult
            if err := ctx.Err(); !errors.Is(err, context.Canceled) {
                res = model.RunResult{
                    ExitCode: util.GetExitCode(&err),
                    Message:  err.Error(),
                }
            } else {
                res = model.RunResult{
                    ExitCode: util.GetExitCode(&err),
                    Message:  "canceled",
                }
            }

            result, err := json.Marshal(res)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Write(result)

        case res := <-resultChan:
            result, err := json.Marshal(res)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Write(result)
        }
    })

    fmt.Printf("Serving at %s \n", BaseUrl)
    err := http.ListenAndServe(BaseUrl, mux)
    log.Fatal(err)
}
