package c

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gitlab.com/iklabib/markisa/container"
	"gitlab.com/iklabib/markisa/model"
	"gitlab.com/iklabib/markisa/util"
)

func Run(encoded string) model.RunResult {
	decoded := util.DecodeAscii85([]byte(encoded))
	return container.RunContainer(decoded, "markisa:common")
}

func Build(source string) model.BuildResult {
	var resp model.BuildResult
	payload := strings.NewReader(source)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080", payload)
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
	json.Unmarshal(body, &resp)

	return resp
}
