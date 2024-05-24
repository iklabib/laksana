package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"codeberg.org/iklabib/markisa/model"
	"codeberg.org/iklabib/markisa/toolchains"
)

func main() {
	resp := model.RunResult{}
	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "FAILED TO READ FROM STDIN"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
		return
	}

	var req model.Submission
	if err := json.Unmarshal(buff, &req); err != nil {
		resp.ExitCode = -1
		resp.Status = "UNMARSHAL FAILED"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
		return
	}

	switch (req.Type) {
	case "python":
		python := toolchains.NewPython()
		dir, err := python.Prep(req.Src, req.SrcTest)
		if err != nil {
			resp.ExitCode = -1
			resp.Status = err.Error()
			jsonified, _ := json.Marshal(resp)
			fmt.Print(string(jsonified))
			return
		}

		resp = python.Eval(dir)
	}

	marshaled, err := json.Marshal(resp)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "MARSHAL FAILED"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
		return
	}

	fmt.Print(string(marshaled))
}
