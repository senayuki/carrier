package lego

import (
	"github.com/senayuki/carrier/types"
	"go.uber.org/zap"
)

type LegoCMD struct {
	C      *types.CertConfig
	path   string
	logger *zap.Logger
}
