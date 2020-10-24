package service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/calebamiles/example-figlet-service/figlet"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HandleFigletizeTextDirect applies a Figlet transformation to provided text without using Cadence
func HandleFigletizeTextDirect(w http.ResponseWriter, req *http.Request) {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		logger.Error("Failed to setup logger", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	if len(bodyBytes) == 0 {
		logger.Error("Read empty request body from client")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := figlet.NewTransformer()
	figletedTxt, err := t.Figletize(string(bodyBytes))
	if err != nil {
		logger.Error("Getting Figlet transformation", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	n, err := w.Write([]byte(figletedTxt))
	if err != nil {
		logger.Error("Writing HandleFigletizeTextCadence response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n != len(figletedTxt) {
		logger.Error(fmt.Sprintf("Expected to write %d bytes, but only wrote %d", len(figletedTxt), n))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
