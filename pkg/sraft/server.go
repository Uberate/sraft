package sraft

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/plugins"
	"github.io/uberate/sraft/pkg/plugins/point"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"strconv"
	"time"
)

const (
	DefaultInnerLogPath = "_inner"
	DefaultDataLogPath  = "_data"
	DefaultLogLogPath   = "_log"
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
	// raft value
	currentTerm uint64
	logs        Logs
	votedFor    string

	// server value
	Id           string
	logger       *logrus.Logger
	status       Status
	server       point.Server
	serverConfig plugins.AnyConfig
	serverKind   string

	// storage value
	dataStorage       storage.Storage
	DataPath          string
	DataStorageKind   string
	DataStorageConfig plugins.AnyConfig
	logsStorage       storage.Storage
	LogsPath          string
	LogsStorageKind   string
	LogsStorageConfig plugins.AnyConfig

	// cluster value
	ClusterConfigPath string // CLusterConfigPath was storage in data storage.
	ClusterConfig     ClusterConfig
	ClusterClients    map[string]point.Client
}

func (s *Server) Init(config ServerConfig) error {

	// ==================================== init logger
	logger := config.LogConfig.ToLogger()
	if logger == nil {
		return fmt.Errorf("Log init error, stop init ")
	}

	if configStr, err := json.Marshal(config); err == nil {
		logger.Debugf("Init config: \n %s", string(configStr))
	}

	// ==================================== init id
	s.Id = config.Id
	if len(s.Id) == 0 {
		logger.Warnf("Server.Id is nil, use now nano second as id.")
		s.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// ==================================== init storages
	s.DataStorageKind = config.DataStorageKind
	s.DataPath = config.DataStoragePath
	s.DataStorageConfig = config.DataStorageConfig
	if len(s.DataPath) == 0 {
		logger.Warnf("Datas log-path is nil, use default value: [%s]", DefaultDataLogPath)
		s.DataPath = DefaultDataLogPath
	}
	var ok bool
	s.dataStorage, ok = storage.GetStorageEngine(s.DataStorageKind)
	if !ok {
		err := fmt.Errorf("Init data stroage fatal, can't find specify implement of storage: [%s],"+
			" please check yout config [DataStorageKind] ", s.DataStorageKind)
		logger.Fatal(err)
		return err
	}
	if err := s.dataStorage.SetConfig(s.DataStorageConfig); err != nil {
		err = fmt.Errorf("Init data storage with [DataStorageConfig] error: %s ", err.Error())
		logger.Fatal(err)
		return err
	}

	s.LogsStorageKind = config.LogsStorageKind
	s.LogsPath = config.LogsStoragePath
	s.LogsStorageConfig = config.LogsStorageConfig
	if len(s.LogsPath) == 0 {
		logger.Warnf("Logs log-path is nil, use default value: [%s]", DefaultLogLogPath)
		s.LogsPath = DefaultLogLogPath
	}
	s.logsStorage, ok = storage.GetStorageEngine(s.LogsStorageKind)
	if !ok {
		err := fmt.Errorf("Init log stroage fatal, can't find specify implement of storage: [%s],"+
			" please check yout config [LogsStorageKind] ", s.LogsStorageKind)
		logger.Fatal(err)
		return err
	}
	if err := s.logsStorage.SetConfig(s.LogsStorageConfig); err != nil {
		err = fmt.Errorf("Init log storage with [LogsStorageConfig] error: %s ", err.Error())
		logger.Fatal(err)
		return err
	}

	// ==================================== init point
	s.serverKind = config.PointKind
	s.serverConfig = config.PointConfig

	if pointInstance, ok := point.GetPoint(s.serverKind); !ok {
		err := fmt.Errorf("Get server point fatal, can't find specify implement of point: [%s],"+
			" please check your config [PointKind] ", s.serverKind)
		logger.Fatal(err)
		return err
	} else {
		var err error
		s.server, err = pointInstance.Server(s.Id, s.serverConfig, s.logger)
		if err != nil {
			err = fmt.Errorf("Init point with [PonitConfig] error: %s ", err.Error())
			logger.Fatal(err)
			return err
		}
	}
	// ==================================== init cluster info

	logger.Info("Server init done")
	return nil
}
