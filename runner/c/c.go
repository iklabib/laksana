package c

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/iklabib/markisa/container"
	"github.com/iklabib/markisa/model"
)

func Run(bin []byte) model.RunResult {
	return container.RunContainer(bin, "markisa:common")
}

func Build(source string) model.BuildResult {
	payload := strings.NewReader(source)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080", payload)
	if err != nil {
		return internalError()
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return internalError()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return internalError()
	}

	// TODO: error logging
	var builderResp model.BuilderResponse
	json.Unmarshal(body, &builderResp)

	return builderResp.BuildResult
}

func internalError() model.BuildResult {
	return model.BuildResult{
		ExitCode: -1,
		Status:   "INTERNAL_ERROR",
	}
}
