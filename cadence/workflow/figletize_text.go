package workflow

import (
	"time"

	"github.com/calebamiles/example-figlet-service/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	// TaskList is the queue to use for GetFortune execution
	TaskList = "figletizeTextTaskList"
)

// FigletizeText applies a Figlet transformation to provided text
func FigletizeText(ctx workflow.Context, inputText string) (string, error) {
	ao := workflow.ActivityOptions{
		TaskList:               TaskList,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, activity.FigletizeText, inputText)
	var figletedText string
	if err := future.Get(ctx, &figletedText); err != nil {

		workflow.GetLogger(ctx).Error("Executing FigletizeText activity", zap.Error(err))
		return "", err
	}

	workflow.GetLogger(ctx).Info("FigletizeText workflow done")
	return figletedText, nil
}
