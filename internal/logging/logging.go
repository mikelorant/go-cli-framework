package logging

import (
	"context"
	"log"
	"os"
)

type Logger struct {
	log.Logger
}

func New(ctx context.Context) *log.Logger {
	return log.New(os.Stdout, "", 0)
}
