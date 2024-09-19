package routes

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gemyago/top-k-system-go/pkg/di"
	"github.com/gemyago/top-k-system-go/pkg/diag"
	"go.uber.org/dig"
)

type router interface {
	Handle(pattern string, handler http.Handler)
}

type Group struct {
	dig.Out

	Mount MountFunc `group:"server"`
}

type MountFunc func(r router)

func WriteData(req *http.Request, log *slog.Logger, writer io.Writer, data []byte) {
	if _, err := writer.Write(data); err != nil {
		log.ErrorContext(req.Context(), "Failed to write response", diag.ErrAttr(err))
	}
}

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewHealthCheckRoutesGroup,
	)
}
