package csharp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"gitlab.com/iklabib/markisa/container"
	"gitlab.com/iklabib/markisa/model"
)

func Run(bin []byte) model.RunResult {
	return container.RunContainer(bin, "markisa:common")
}

func Build(source string) model.BuildResult {
	var resp model.BuildResult
	payload := strings.NewReader(source)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081", payload)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		return resp
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		return resp
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		return resp
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		log.Println(err.Error())
		resp.ExitCode = -1
		resp.Status = "SERIALIZATION_ERROR"
		return resp
	}

	return resp
}
