package sraft

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// ==================================== Handler define

type Handler func(message SendMessage) (*ReceiveMessage, error)

// ==================================== Define the message struct

// SendMessage is the point send to server package.
type SendMessage struct {
	Body map[string]interface{}
}

func SendMessageFromAny(obj any) (SendMessage, error) {
	value := map[string]interface{}{}
	err := mapstructure.Decode(obj, &value)
	if err != nil {
		return SendMessage{}, err
	}

	res := SendMessage{
		Body: value,
	}
	return res, nil
}

// ErrorMessage implement the error, and if contain the error info of server.
type ErrorMessage struct {
	code string
	Info string
}

// Code return the code of ErrorMessage. The Code follower the http.status code .
func (e ErrorMessage) Code() string {
	return e.code
}

func (e ErrorMessage) Error() string {
	return e.Info
}

// ReceiveMessage is the point server response of client.
type ReceiveMessage struct {
	Response map[string]interface{}
}

// ==================================== Define the point struct

// Point define the protocol of TransLayout of RPC.
type Point interface {
	Client(id string, config AnyConfig, logger *logrus.Logger) (Client, error)
	Server(id string, config AnyConfig, logger *logrus.Logger) (Server, error)
}

// Client is a protocol client, and Client only bind with one server, but server not care it.
type Client interface {
	Name() string
	SendMessage(path string, message SendMessage) (*ReceiveMessage, error)
	SendAny(path string, value any) (*ReceiveMessage, error)
}

type Server interface {
	Name() string

	Handler(path string, hand Handler)

	Run() error
	Stop() error
}
