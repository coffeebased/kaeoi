package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/coffeebased/kaeoi/internal/gameserver"
)

type Handler struct {
	ctx            context.Context
	store          gameserver.Store
	monitor        *gameserver.Monitor
	wsWriteTimeout time.Duration
	mux            *http.ServeMux
}

func New(ctx context.Context, logger *slog.Logger, store gameserver.Store, monitor *gameserver.Monitor, wsWriteTimeout time.Duration) http.Handler {
	if ctx == nil {
		panic("nil context")
	}

	if logger == nil {
		panic("nil logger")
	}

	if store == nil {
		panic("nil store")
	}

	if monitor == nil {
		panic("nil monitor")
	}

	if wsWriteTimeout <= 0 {
		panic("non-positive websocket write timeout")
	}

	handler := &Handler{
		ctx:            ctx,
		store:          store,
		monitor:        monitor,
		wsWriteTimeout: wsWriteTimeout,
		mux:            http.NewServeMux(),
	}

	handler.mux.HandleFunc("GET /ws", handler.subscribeServers)

	return logRequests(logger, handler)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) subscribeServers(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer func() {
		_ = conn.CloseNow()
	}()

	ctx := conn.CloseRead(h.ctx)

	for server := range h.monitor.Subscribe(ctx) {
		writeCtx, writeCancel := context.WithTimeout(ctx, h.wsWriteTimeout)
		err := wsjson.Write(writeCtx, conn, newGameServerResponse(server))
		writeCancel()

		if err != nil {
			return
		}
	}
}
