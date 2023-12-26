package websockets

import (
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type WsClients map[*websocket.Conn]struct{}

type NotificationWS struct {
	log     *zap.SugaredLogger
	tracer  trace.Tracer
	mu      sync.Mutex
	clients WsClients
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewNotificationWS(log *zap.SugaredLogger, tracer trace.Tracer) *NotificationWS {
	return &NotificationWS{log: log, tracer: tracer, clients: make(WsClients)}
}

func (n *NotificationWS) HandleWS(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		n.log.Errorf("cannot upgrade http request: %v", err)
		return
	}
	defer conn.Close()

	n.mu.Lock()
	defer n.mu.Unlock()
	n.clients[conn] = struct{}{}

}
