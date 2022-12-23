package point

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/plugins"
	"net/http"
)

// ==================================== Handler define

type Handler func(path string, message SendMessage) *ReceiveMessage

// ==================================== Define the message struct

// SendMessage is the point send to server package.
type SendMessage struct {
	Body map[string]interface{}
}

func (sm *SendMessage) ToAny(obj any) error {
	return mapstructure.Decode(sm.Body, obj)
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
	Info string
}

func (e ErrorMessage) Error() string {
	return e.Info
}

// ReceiveMessage is the point server response of client.
type ReceiveMessage struct {
	Code         int
	Response     map[string]interface{}
	ErrorMessage ErrorMessage
}

func (rm ReceiveMessage) ToAny(obj any) error {
	return mapstructure.Decode(rm.Response, obj)
}

func ReceiveMessageFromAny(obj any) (ReceiveMessage, error) {
	value := map[string]interface{}{}
	err := mapstructure.Decode(obj, &value)
	if err != nil {
		return ReceiveMessage{}, err
	}

	res := ReceiveMessage{
		Code:     http.StatusOK,
		Response: value,
	}
	return res, nil
}

func QuickErrorReceiveMessage(code int, err error) *ReceiveMessage {
	errMessage := ErrorMessage{
		Info: err.Error(),
	}

	return &ReceiveMessage{
		Code:         code,
		ErrorMessage: errMessage,
	}
}

// ==================================== Define the point struct

// Point define the protocol of TransLayout of RPC.
type Point interface {
	Name() string
	Client(id string, config plugins.AnyConfig, logger *logrus.Logger) (Client, error)
	Server(id string, config plugins.AnyConfig, logger *logrus.Logger) (Server, error)
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

// ==================================== Points define

var points = map[string]Point{
	HttpV1Engine{}.Name(): HttpV1Engine{}, // HttpV1Engine
}

// RegisterPoint will append or cover a point instance to inner point set. GetPoint will get point from this set.
func RegisterPoint(instance Point) {
	points[instance.Name()] = instance
}

// GetPoint will return specify point instance from inner set by name. And not found specify point, return false.
func GetPoint(name string) (Point, bool) {
	res, ok := points[name]
	return res, ok
}
