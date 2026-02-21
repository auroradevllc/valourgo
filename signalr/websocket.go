package signalr

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WSConn struct {
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
	once sync.Once
}

func DialWS(ctx context.Context, url string, h http.Header) (*WSConn, error) {
	ws, _, err := websocket.DefaultDialer.DialContext(ctx, url, h)

	if err != nil {
		return nil, err
	}

	w := &WSConn{
		conn: ws,
		send: make(chan []byte, 128),
		done: make(chan struct{}),
	}

	go w.sendLoop()

	return w, nil
}

func (w *WSConn) sendLoop() {
	for {
		b, ok := <-w.send

		if !ok {
			return
		}

		err := w.conn.WriteMessage(websocket.TextMessage, b)

		if err != nil {
			log.WithError(err).Error("Unable to send data to websocket")
		}
	}
}

func (w *WSConn) Send(b []byte) error {
	w.send <- append(b, 0x1e)
	return nil
}

func (w *WSConn) Read() ([]byte, error) {
	_, msg, err := w.conn.ReadMessage()
	return msg, err
}

func (w *WSConn) Close() error {
	var err error

	w.once.Do(func() {
		close(w.done)
		err = w.conn.Close()
	})

	return err
}
