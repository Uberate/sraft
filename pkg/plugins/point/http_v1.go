package point

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/sraft"
	"io/ioutil"
	"net/http"
)

const HttpV1EngineKind = "http_v1"

// ==================================== Engine define

type HttpV1Engine struct {
}

func (h HttpV1Engine) Client(id string, config sraft.AnyConfig, logger *logrus.Logger) (sraft.Client, error) {
	client := HttpV1Client{
		Id:     id,
		Logger: logger,
	}

	if err := config.ToAny(&client); err != nil {
		return nil, err
	}

	return &client, nil
}

func (h HttpV1Engine) Server(id string, config sraft.AnyConfig, logger *logrus.Logger) (sraft.Server, error) {
	server := HttpV1Server{}

	if err := config.ToAny(&server); err != nil {
		return nil, err
	}
	return &server, nil
}

// ==================================== Client define

type HttpV1Client struct {
	Id     string         `mapstructure:"-"`
	Logger *logrus.Logger `mapstructure:"-"`

	// custom config
	TargetPoint string
}

func (h *HttpV1Client) Name() string {
	return HttpV1EngineKind + ":" + h.Id
}

func (h *HttpV1Client) SendMessage(path string, message sraft.SendMessage) (*sraft.ReceiveMessage, error) {
	requestJson, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	receiveMessage := &sraft.ReceiveMessage{}
	if err = json.Unmarshal(resBody, receiveMessage); err != nil {
		return nil, err
	}

	return receiveMessage, nil

}

func (h *HttpV1Client) SendAny(path string, value any) (*sraft.ReceiveMessage, error) {
	message, err := sraft.SendMessageFromAny(value)
	if err != nil {
		return nil, err
	}

	return h.SendMessage(path, message)
}

// ==================================== Server define

type HttpV1Server struct {
	// private values

	stopChan    chan bool
	alreadyStop chan bool

	handlers map[string]sraft.Handler

	// custom config
	Point string
}

func (h *HttpV1Server) Name() string {
	return HttpV1EngineKind
}

func (h *HttpV1Server) Handler(path string, hand sraft.Handler) {
	h.handlers[path] = hand
}

func (h *HttpV1Server) Run() error {

	if err := http.ListenAndServe(h.Point, h); err != nil {
		return err
	}

	select {
	case <-h.stopChan:
	}

	h.alreadyStop <- true
	return nil
}

func (h *HttpV1Server) Stop() error {
	h.stopChan <- true

	select {
	case <-h.alreadyStop:
		break
	}

	return nil
}

func (h *HttpV1Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
