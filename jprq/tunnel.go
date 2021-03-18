package jprq

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/gosimple/slug"
	"github.com/labstack/gommon/log"
)

type Tunnel struct {
	host           string
	port           int
	conn           *websocket.Conn
	token          string
	requests       map[uuid.UUID]RequestMessage
	requestChan    chan RequestMessage
	responseChan   chan ResponseMessage
	numOfReqServed int
}

func (j Cfc) GetTunnelByHost(host string) (*Tunnel, error) {
	t, ok := j.tunnels[host]
	if !ok {
		return t, errors.New("Tunnel doesn't exist")
	}

	return t, nil
}

func (j *Cfc) AddTunnel(username string, port int, conn *websocket.Conn) *Tunnel {
	username = slug.Make(username)
	host := fmt.Sprintf("%s.%s", username, j.baseHost)

	_, err := j.GetTunnelByHost(host)
	if err == nil {
		adj := getRandomAdj()
		host = fmt.Sprintf("%s-%s", adj, host)
	}

	token := generateToken()
	requests := make(map[uuid.UUID]RequestMessage)
	requestChan, responseChan := make(chan RequestMessage), make(chan ResponseMessage)
	tunnel := Tunnel{
		host:         host,
		port:         port,
		conn:         conn,
		token:        token,
		requests:     requests,
		requestChan:  requestChan,
		responseChan: responseChan,
	}

	log.Info("New Tunnel: ", host)
	j.tunnels[host] = &tunnel
	return &tunnel
}

func (j *Cfc) DeleteTunnel(host string) {
	tunnel, ok := j.tunnels[host]
	if !ok {
		return
	}
	log.Infof("Deleted Tunnel: %s, Number Of Requests Served: %d", host, tunnel.numOfReqServed)
	close(tunnel.requestChan)
	close(tunnel.responseChan)
	delete(j.tunnels, host)
}

func (tunnel *Tunnel) DispatchRequests() {
	for {
		select {
		case requestMessage, more := <-tunnel.requestChan:
			if !more {
				return
			}
			messageContent, _ := json.Marshal(requestMessage)
			tunnel.requests[requestMessage.ID] = requestMessage
			tunnel.conn.WriteMessage(websocket.TextMessage, messageContent)
		}
	}
}

func (tunnel *Tunnel) DispatchResponses() {
	for {
		select {
		case responseMessage, more := <-tunnel.responseChan:
			if !more {
				return
			}
			requestMessage, ok := tunnel.requests[responseMessage.RequestId]
			if !ok {
				log.Error("Request Not Found", responseMessage.RequestId)
				continue
			}

			requestMessage.ResponseChan <- responseMessage
			delete(tunnel.requests, requestMessage.ID)
			tunnel.numOfReqServed++
		}
	}
}
