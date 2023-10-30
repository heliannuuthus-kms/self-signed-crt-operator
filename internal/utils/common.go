package utils

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	Logger logr.Logger
)

func init() {
	Logger = zap.New(zap.UseFlagOptions(&zap.Options{Development: false}))
}

type Entry[L any, R any] struct {
	Left  L
	Right R
}
