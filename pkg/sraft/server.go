package sraft

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/plugins"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"strconv"
	"time"
)

// ServerConfig save all the server config(not contain the cluster config).
type ServerConfig struct {
	Id string `json:"id" yaml:"id"`

	PointKind   string
	PointConfig plugins.AnyConfig

	LogConfig LogConfig

	DataStorageKind   string
	DataStorageConfig plugins.AnyConfig

	LogsStorageKind   string
	LogsStorageConfig plugins.AnyConfig

	LogsStoragePath    string
	DataStoragePath    string
	ClusterStoragePath string
}

type LogConfig struct {
	TimestampFormat  string
	DisableTimestamp bool
	DisableColor     bool
	FullTimestamp    bool
	Level            string
}

func (lc LogConfig) ToLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: lc.DisableTimestamp,
		TimestampFormat:  lc.TimestampFormat,
		DisableColors:    lc.DisableColor,
		FullTimestamp:    lc.FullTimestamp,
	})

	level, err := logrus.ParseLevel(lc.Level)
	if err != nil {
		logger.Errorf("Parse level error: [%s], %s", lc.Level, err.Error())
	} else {
		logger.SetLevel(level)
	}

	logger.Info("Log init done")

	return logger
}

type ClusterConfig struct {
	ClusterElement []ClusterElement
}

type ClusterElement struct {
	Id           string
	Kind         string
	TargetPoint  string
	ClientConfig plugins.AnyConfig
}

// Server is a node of sraft.
//
// Server abstract the network and storage layout. The network was abstract to point.Point. The storage was abstract to
// storage.Storage.
//
// Status:
// The serve has three status: FollowerStatus CandidateStatus LeaderStatus. For more status detail see Status.
type Server struct {
	Id string

	currentTerm uint64
	status      Status
	logger      *logrus.Logger

	logs     Logs
	votedFor string

	dataStorage       storage.Storage
	DataPath          string
	DataStorageKind   string
	DataStorageConfig plugins.AnyConfig

	ClusterConfigPath string // CLusterConfigPath was storage in data storage.
	ClusterConfig     ClusterConfig

	logsStorage storage.Storage // LogStorage can use same storage from data.
	LogsPath    string
	// If LogStorageKind is empty, use dataStorage with LogPath, but the LogPath can't same with DataPath.
	LogsStorageKind   string
	LogsStorageConfig plugins.AnyConfig
}

func (s *Server) Init(config ServerConfig) error {

	logger := config.LogConfig.ToLogger()
	if logger == nil {
		return fmt.Errorf("Log init error, stop init ")
	}

	if configStr, err := json.Marshal(config); err == nil {
		logger.Debugf("Init config: \n %s", string(configStr))
	}

	// init id
	s.Id = config.Id
	if len(s.Id) == 0 {
		logger.Warnf("Server.Id is nil, use now nano second as id.")
		s.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// init storages
	s.DataStorageKind = config.DataStorageKind

	return nil
}
