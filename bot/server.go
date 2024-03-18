package bot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int32  `yaml:"port"`
	Path        string `yaml:"path"`
	AccessToken string `yaml:"access_token"`
}

type Server struct {
	engine   *nbhttp.Engine
	upgrader *websocket.Upgrader
}

func NewServer(cfg ServerConfig) *Server {
	mux := http.NewServeMux()
	upgrader := newUpgrader()
	server := &Server{
		upgrader: upgrader,
	}
	mux.HandleFunc(cfg.Path, server.onWebsocket)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{addr},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
	})
	logging.SetLevel(logging.LevelNone) // disable nbio logging
	server.engine = engine
	return server
}

func (s *Server) ListenAndServe() error {
	if err := s.engine.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}

func (s *Server) onWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log().Errorf("[%s] Websocket upgrade failed: %v", conn.RemoteAddr().String(), err)
		return
	}
	Log().Infof("[%s] Websocket connection established", conn.RemoteAddr().String())
}

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.OnOpen(onOpen)
	u.OnMessage(onMessage)
	u.OnClose(onClose)
	return u
}

func onOpen(c *websocket.Conn) {
	Log().Debugf("[%s] Websocket connection opened", c.RemoteAddr().String())
}

func onMessage(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
	Log().Debugf("[%s] Message received: type=%s, data=%s", c.RemoteAddr().String(), messageTypeToString(messageType), string(data))
}

func messageTypeToString(t websocket.MessageType) string {
	switch t {
	case 0:
		return "FragmentMessage"
	case 1:
		return "TextMessage"
	case 2:
		return "BinaryMessage"
	case 8:
		return "CloseMessage"
	case 9:
		return "PingMessage"
	case 10:
		return "PongMessage"
	default:
		return "UndefinedMessage"
	}
}

func onClose(c *websocket.Conn, err error) {
	if err != nil {
		Log().Warnf("[%s] Websocket connection closed with error: %v", c.RemoteAddr().String(), err)
		return
	}
	Log().Debugf("[%s] Websocket connection closed cleanly", c.RemoteAddr().String())
}
