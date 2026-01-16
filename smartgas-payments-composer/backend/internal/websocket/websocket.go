package websocket

import (
	"github.com/gorilla/websocket"
)

type Clients map[*websocket.Conn]bool

type Channel map[string]Clients

func (cs Clients) BroadcastJson(message any) {
	for client := range cs {
		client.WriteJSON(message)
	}
}

func (cs Clients) CloseConnections() {
	for client := range cs {
		cs.CloseClientConnection(client)
	}
}

func (cs Clients) AddClient(ws *websocket.Conn) {
	cs[ws] = true
}

func (cs Clients) CloseClientConnection(ws *websocket.Conn) {
	ws.WriteMessage(websocket.CloseNormalClosure, []byte("Normal closed"))
	ws.Close()
	delete(cs, ws)
}

func (ch Channel) GetChannel(channelId string) Clients {
	if ch[channelId] == nil {
		ch[channelId] = make(Clients)
	}

	return ch[channelId]
}

func (ch Channel) DeleteChannel(channelId string) {
	channel := ch[channelId]
	channel.CloseConnections()
	delete(ch, channelId)
}

func InitChannels() Channel {
	return make(Channel)
}
