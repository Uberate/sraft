package sraft

import (
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/plugins"
	"github.io/uberate/sraft/pkg/plugins/point"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"time"
)

const (
	DefaultLogLevel = "INFO"
)

type ServerConfig struct {
	ID string

	LogConfig LogConfig

	ServerPointKind    string
	ServerPointConf    plugins.AnyConfig
	ClusterConfig      ClusterConfig
	StorageKind        string
	StorageConf        plugins.AnyConfig
	DataStoragePath    string
	LogStoragePath     string
	ClusterStoragePath string

	TimeoutBase  int64
	TimeoutRange int64
}

type LogConfig struct {
	TimestampFormat  string
	DisableTimestamp bool
	DisableColor     bool
	FullTimestamp    bool
	Level            string
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		DisableColor:     false,
		FullTimestamp:    true,
		Level:            DefaultLogLevel,
	}
}

// ToLogger return a logrus.Logger instance, use LogConfig init.
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
		logger.Errorf("Parse level error: LevelValue: [%s], Err: [%s]; use defualt level: [%s]",
			lc.Level, err.Error(), DefaultLogLevel)
		level, _ = logrus.ParseLevel(DefaultLogLevel)
	} else {
		logger.SetLevel(level)
	}

	logger.Info("Log init done")

	return logger
}

type Server struct {

	// ==================================== Server base param

	// ID is the server identify, in the cluster, the ID is name of one server. The cluster can't have two server use
	// same ID value.
	ID             string
	stopChan       chan bool
	heartbeatTimer *time.Ticker

	// ==================================== Behavior params
	serverImpl point.Server
	clients    map[string]point.Client

	dataStorage       storage.Storage
	logEntryPath      string
	dataPath          string
	clusterConfigPath string

	randomTimeout int64

	// ==================================== Server status info

	// Status represent the server status: FollowerStatus CandidateStatus LeaderStatus, any other status is error
	// status. On server start, the Status should be FollowerStatus.
	status Status
	// VoteFor the leader id, and the may be candidate. On server start, the VoteFor should empty.
	voteFor     string
	currentTerm uint64

	// ==================================== State machine value he

	// logs save all the log in the array. In one cluster, any server will get same log by same log-index in different
	// server.
	//
	// The log not care which log kind: DataLogEntry or ClusterLogEntry. And not all the logs are committed. Once the
	// log already be committed by majority server, the logs is committed.
	logs []Log
	// nextAppendsAt storage the server next append index. The keys are Server.ID, and keys is controller by
	// clusterConfig. And the keys should stable in one server.
	nextAppendsAt map[string]uint64
	// nextCommittedAt the next committed log's index. On server start, the nextCommittedAt is zero.
	nextCommittedAt uint64

	clusterConfig ClusterConfig
}
