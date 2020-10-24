package activity

import (
	"context"

	"github.com/calebamiles/example-figlet-service/figlet"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

// FigletizeText applies a Figlet transform to provided text
func FigletizeText(ctx context.Context, inputText string) (string, error) {
	t := figlet.NewTransformer()

	transformedTxt, err := t.Figletize(inputText)
	if err != nil {
		activity.GetLogger(ctx).Error("transforming text", zap.Error(err))
		return "", nil
	}

	activity.GetLogger(ctx).Info("FigletizeText called", zap.String("Text", transformedTxt))
	return transformedTxt, nil
}
