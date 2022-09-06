package shell

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"goa.design/goa/v3/http"
	stdhttp "net/http"
)

type EncoderWithLogs struct {
	encoder http.Encoder
	logger  *zap.Logger
	kind    string
}

func ResponseEncoderWithLogs(ctx context.Context, r stdhttp.ResponseWriter) http.Encoder {
	return &EncoderWithLogs{
		encoder: http.ResponseEncoder(ctx, r),
		logger:  Logger(ctx),
		kind:    "response",
	}
}

func (e EncoderWithLogs) Encode(v interface{}) error {
	err := e.encoder.Encode(v)
	if err != nil {
		e.logger.Error(fmt.Sprintf("failed to decode %s", e.kind), zap.Error(err))
	}

	return err
}
