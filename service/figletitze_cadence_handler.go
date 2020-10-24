package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/calebamiles/example-figlet-service/cadence/workflow"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// hostPort is the location of the cadence frontend
	hostPort = "127.0.0.1:7933"

	// clientName is the name of the client
	clientName = "figletize-text-http-handler"

	// clientService is the Cadence service to connect to
	clientService = "cadence-frontend"
)

// HandleFigletizeTextCadence applies a Figlet transformation through executation of a Cadence workflow
func HandleFigletizeTextCadence(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cadenceOpts := &client.Options{
		Identity: clientName,
	}

	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		logger.Error("Failed to setup logger", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.Method != http.MethodPost {
		logger.Error("non POST request: %s cannot be handled", zap.String("HTTP method", req.Method))
		w.WriteHeader(http.StatusBadRequest)
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

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(clientName))
	if err != nil {
		logger.Error("Failed to setup tchannel", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: clientName,
		Outbounds: yarpc.Outbounds{
			clientService: {Unary: ch.NewSingleOutbound(hostPort)},
		},
	})

	if err := dispatcher.Start(); err != nil {
		logger.Error("Failed to start dispatcher", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thriftService := workflowserviceclient.New(dispatcher.ClientConfig(clientService))
	cadence := client.NewClient(thriftService, "hcp", cadenceOpts)

	startOpts := client.StartWorkflowOptions{
		TaskList:                        workflow.TaskList,
		ExecutionStartToCloseTimeout:    60 * time.Second,
		DecisionTaskStartToCloseTimeout: 20 * time.Second,
		WorkflowIDReusePolicy:           client.WorkflowIDReusePolicyAllowDuplicateFailedOnly,
		Memo:                            map[string]interface{}{"workflow-type": "local-development"},
	}

	future, err := cadence.ExecuteWorkflow(ctx, startOpts, workflow.FigletizeText, string(bodyBytes))
	if err != nil {
		logger.Error("Executing FigletizeText workflow", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var figletedTxt string
	err = future.Get(ctx, &figletedTxt)
	if err != nil {
		logger.Error("Getting FigletizeText workflow result", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
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
