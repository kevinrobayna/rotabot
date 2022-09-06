package shell

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func ErrorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, rawError error) {
		Logger(ctx).Error(
			"unexpected_error",
			zap.Error(rawError),
		)
		http.Error(w, fmt.Sprintf("%s: Unexpected error occurred", "http"), http.StatusInternalServerError)
	}
}
