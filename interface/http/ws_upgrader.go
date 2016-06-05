/*Package ifaceHttp provides HTTP interfaces
to the app's border entities.*/
package ifaceHttp

import (
	"github.com/gorilla/websocket"
	"github.com/utrack/gallery/client/ws"
	"github.com/utrack/gallery/hub"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWs handles websocket requests from the peer.
func ServeWs(h hub.ConnectionAcceptor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := clientWs.NewClient(ws)
		err = h.Accept(client)
		if err != nil {
			log.Printf("Error when accepting websocket client: %v", err)
			client.Disconnect()
		}
	}
}
