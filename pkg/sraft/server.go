package sraft

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.io/uberate/sraft/pkg/plugins"
	"github.io/uberate/sraft/pkg/plugins/point"
	"github.io/uberate/sraft/pkg/plugins/storage"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultLogLevel = "INFO"

	DefaultDataStoragePath    = "data"
	DefaultLogStoragePath     = "_inner/log"
	DefaultClusterStoragePath = "_inner/cluster"

	DefaultTimeout      = int64(150 * time.Millisecond)
	DefaultTimeoutRange = int64(100 * time.Millisecond)

	RequestVotePath   = "/sraft/request-vote"
	AppendEntriesPath = "/sraft/append-entries"
)

type ServerConfig struct {
	ID string

	LogConfig *LogConfig

	ServerPointKind    string
	ServerPointConf    plugins.AnyConfig
	ClusterConfig      *ClusterConfig
	StorageKind        string
	StorageConf        plugins.AnyConfig
	DataStoragePath    string
	LogStoragePath     string
	ClusterStoragePath string

	TimeoutBase  int64
	TimeoutRange int64
	RandomSeed   int64
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		ID: time.Now().Format(time.RFC3339),

		LogConfig: DefaultLogConfig(),

		ServerPointKind: "",
		ServerPointConf: plugins.AnyConfig{
			Config: map[string]string{},
		},
		ClusterConfig: &ClusterConfig{},
		StorageKind:   "",
		StorageConf: plugins.AnyConfig{
			Config: map[string]string{},
		},
		DataStoragePath:    DefaultDataStoragePath,
		LogStoragePath:     DefaultLogStoragePath,
		ClusterStoragePath: DefaultClusterStoragePath,

		TimeoutBase:  DefaultTimeout,
		TimeoutRange: DefaultTimeoutRange,
		RandomSeed:   int64(time.Now().Nanosecond()),
	}
}

func (s *Server) RandomTimeout() {
	// stop timer first
	s.heartbeatTicker.Stop()
	s.logger.Tracef("random timeout value: randomTimeout + randomRange")
	s.logger.Debugf("Server.randomTimeout: [%d]", s.randomTimeout)
	s.logger.Debugf("Server.randomTimeoutRange: [%d]", s.randomTimeoutRange)
	s.logger.Debugf("Server.randomSeed")
	randInstance := rand.New(rand.NewSource(s.randomSeed))
	s.randomTimeoutValue = s.randomTimeout + randInstance.Int63n(s.randomTimeoutRange)
	s.heartbeatTicker.Reset(time.Duration(s.randomTimeoutValue) * time.Millisecond)
	s.logger.Tracef("random timeout value success: new timer: [%dms]", s.randomTimeoutValue)
}

func (sc *ServerConfig) Server() (*Server, error) {
	s := &Server{
		stopChan: make(chan bool, 0),
	}

	// init logger
	if sc.LogConfig == nil {
		return nil, fmt.Errorf("log config is nil, stop server init")
	}

	s.logger = sc.LogConfig.ToLogger()
	s.logger.Trace("logger init done")

	if s.logger.Level <= logrus.DebugLevel {
		// show the config in json format
		configJsonStrByte, _ := json.Marshal(sc)
		s.logger.Debugf("server config: %s", string(configJsonStrByte))
	}

	// ==================================== Server base param
	s.logger.Trace("verify server ID")
	{
		s.ID = sc.ID
		// if s.ID length == nil, use time as ID
		if len(s.ID) == 0 {
			s.ID = time.Now().Format(time.RFC3339)
			s.logger.Warnf("server ID is nil, set now time str (with Format: RFC3339) as ID: [%s]", s.ID)
		}
		s.logger.Trace("verify server ID result: success")
	}
	s.currentTerm = 0
	s.logs = make([]Log, 0, 0)
	s.voteFor = ""
	s.status = FollowerStatus
	s.stopChan = make(chan bool, 0)
	if sc.ClusterConfig == nil {
		// is the single server mod
	} else {
		s.nextAppendsAt = make(map[string]uint64, len(sc.ClusterConfig.ClusterElement))
		s.nextCommittedAt = 0
	}
	s.singleLocker = &sync.Mutex{}

	// ==================================== Behavior params
	s.logger.Trace("init server point")
	{
		pointInstance, ok := point.GetPoint(sc.ServerPointKind)
		if !ok {
			err := fmt.Errorf("not found server pointInstance: [%s]", sc.ServerPointKind)
			s.logger.Errorf("init err: %v", err)
			return nil, err
		}
		if serverImpl, err := pointInstance.Server(s.ID, sc.ServerPointConf, s.logger); err != nil {
			s.logger.Errorf("init server impl error: %v", err)
			return nil, err
		} else {
			s.serverImpl = serverImpl
		}
		//TODO register proxy
		s.serverImpl.Handler(RequestVotePath, s.RequestVoteRPC)
		s.serverImpl.Handler(AppendEntriesPath, s.AppendEntryRPC)
		s.logger.Trace("init server point result: success")
	}

	s.logger.Tracef("init random timeout")
	{
		s.randomSeed = sc.RandomSeed
		s.randomTimeout = sc.TimeoutBase
		s.randomTimeoutRange = sc.TimeoutRange
		if s.randomSeed == 0 {
			// if random seed is zero, use now time unix time nanosecond as seed
			s.logger.Warnf("random seed is zero, use time.now().UnixNane() as seed")
			s.randomSeed = time.Now().UnixNano()
		}
		s.logger.Debugf("random seed: %v", s.randomSeed)
		s.heartbeatTicker = time.NewTicker(time.Duration(s.randomTimeout) * time.Millisecond)
		s.RandomTimeout()
		s.logger.Tracef("init random timeout result: success")
	}

	s.logger.Tracef("init storage")
	{
		s.dataPath = sc.DataStoragePath
		s.logEntryPath = sc.LogStoragePath
		s.clusterConfigPath = sc.ClusterStoragePath
		// those three params need different value
		if s.dataPath == s.logEntryPath ||
			s.dataPath == s.clusterConfigPath ||
			s.clusterConfigPath == s.logEntryPath {
			err := fmt.Errorf("the log path err, the value should different: dataPath != logPath != "+
				"clusterPath, but not: DataPath: [%s], LogPath: [%s], ClusterPath: [%s]", s.dataPath, s.logEntryPath,
				sc.ClusterStoragePath)
			s.logger.Errorf("Init storageImpl err: %v", err)
			return nil, err
		}
		if storageImpl, ok := storage.GetStorageEngine(sc.StorageKind); ok {
			s.dataStorage = storageImpl
			if err := s.dataStorage.SetConfig(sc.StorageConf); err != nil {
				s.logger.Errorf("init storage err: %v", err)
				return nil, err
			}
		} else {
			err := fmt.Errorf("not found data storage kind: [%s]", sc.StorageKind)
			s.logger.Errorf("init storage err: %v", err)
			return nil, err
		}
		s.logger.Tracef("init storage result: success")
	}

	s.logger.Tracef("init cluster info")
	{

		s.clusterConfig = sc.ClusterConfig
		if sc.ClusterConfig == nil {
			s.clients = make(map[string]point.Client, 0)
			s.logger.Warnf("cluster config is nil, no member found, skip cluster init")
		} else {
			s.clients = make(map[string]point.Client, len(sc.ClusterConfig.ClusterElement))
			for _, memberConfig := range sc.ClusterConfig.ClusterElement {
				if memberConfig.Id == s.ID {
					continue
					// skip self
				}
				if memberPoint, ok := point.GetPoint(memberConfig.Kind); ok {
					client, err := memberPoint.Client(memberConfig.Id, memberConfig.ClientConfig, s.logger)
					if err != nil {
						err = fmt.Errorf("err init cluster [%s]: %v", memberConfig.Id, err)
						s.logger.Errorf("%v", err)
						return nil, err
					}
					s.clients[memberConfig.Id] = client
				} else {
					err := fmt.Errorf("not found client [%s] kind: [%s]", memberConfig.Id, memberConfig.Kind)
					s.logger.Errorf("%v", err)
					return nil, err
				}
			}
		}
		s.logger.Tracef("init cluster info: success")
	}

	// Generate success, return instance and nil (of error).
	return s, nil
}

type LogConfig struct {
	TimestampFormat  string `json:"timestamp_format,omitempty" yaml:"timestamp_format" mapstructure:"timestamp_format"`
	DisableTimestamp bool   `json:"disable_timestamp,omitempty" yaml:"disable_timestamp" mapstructure:"disable_timestamp"`
	DisableColor     bool   `json:"disable_color,omitempty" yaml:"disable_color" mapstructure:"disable_color"`
	FullTimestamp    bool   `json:"full_timestamp,omitempty" yaml:"full_timestamp" mapstructure:"full_timestamp"`
	Level            string `json:"level,omitempty" yaml:"level" mapstructure:"level"`
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

// Server is the sraft core struct, if represent a client of a raft server or node. And the server is a status machine.
type Server struct {
	singleLocker *sync.Mutex

	// ==================================== Server base param

	// ID is the server identify, in the cluster, the ID is name of one server. The cluster can't have two server use
	// same ID value.
	ID       string
	logger   *logrus.Logger
	stopChan chan bool

	// ==================================== Behavior params
	serverImpl point.Server
	clients    map[string]point.Client

	dataStorage       storage.Storage
	logEntryPath      string
	dataPath          string
	clusterConfigPath string

	randomSeed         int64
	randomTimeout      int64
	randomTimeoutRange int64
	randomTimeoutValue int64
	heartbeatTicker    *time.Ticker

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

	clusterConfig *ClusterConfig
}

func (s *Server) Run() error {
	go s.LifeCycleRun()

	var err error
	go func() {
		err = s.serverImpl.Run()
	}()

	select {
	case <-s.stopChan:
		if err := s.serverImpl.Stop(); err != nil {
			return err
		}
	}
	return err
}

func (s *Server) Stop() {
	s.stopChan <- true
}

func (s *Server) LifeCycleRun() {
	for {
		select {
		case <-s.heartbeatTicker.C:
			// status convert
			switch s.status {
			case FollowerStatus:
				s.followerLifeCycleNext()
			case CandidateStatus:
				// keep
				s.candidateLifeCycleNext()
			case LeaderStatus:
				s.leaderLifeCycleNext()
			default:
				// if status err, reset to follower
				s.status = FollowerStatus
			}
		}
	}
}

func (s *Server) autoResult(body any) *point.ReceiveMessage {
	// not found id reject vote
	res, err := point.ReceiveMessageFromAny(body)
	if err != nil {
		s.logger.Errorf("Response RPC err: %v", err)
		return point.QuickErrorReceiveMessage(http.StatusInternalServerError, err)
	}
	return &res
}

func (s *Server) selfInfoDebugLog() {
	s.logger.Debugf("self info: term: [%d], vote for: [%s]", s.currentTerm, s.voteFor)
}

// ==================================== vote request rpc

type RequestVoteBody struct {
	RequestID            string
	RequestTerm          uint64
	RequestNextCommitted uint64
}
type RequestVoteResponse struct {
	Res bool
}

func (s *Server) RequestVote() (memberCount, voteCount, ignoreCount int) {
	memberCount = len(s.clients) + 1
	ignoreCount = 0 // todo: when cluster config update, some server should be ignore
	voteCount = 1

	s.voteFor = s.ID
	voteCountLock := sync.Mutex{}
	for id, clientImpl := range s.clients {
		idCopy := id
		clientImplCopy := clientImpl
		go func() {
			res, err := clientImplCopy.SendAny(RequestVotePath, RequestVoteBody{
				RequestID:            s.ID,
				RequestTerm:          s.currentTerm,
				RequestNextCommitted: s.nextCommittedAt,
			})
			if err != nil {
				s.logger.Errorf("request [%s] vote error: %v", idCopy, err)
				return
			}
			resObj := RequestVoteResponse{}
			if res.ToAny(&resObj) != nil {
				s.logger.Errorf("parse [%s] vote res error: %v", idCopy, err)
				return
			}
			if resObj.Res {
				voteCountLock.Lock()
				voteCount++
				defer voteCountLock.Unlock()
			}
		}()
	}
	timer := time.NewTimer(time.Duration(s.randomTimeoutValue) * time.Millisecond)
	select {
	case <-timer.C:
	}
	return
}

func (s *Server) RequestVoteRPC(path string, body point.SendMessage) *point.ReceiveMessage {
	s.singleLocker.Lock()
	defer s.singleLocker.Unlock()

	rvb := RequestVoteBody{}

	if err := body.ToAny(&rvb); err != nil {
		return point.QuickErrorReceiveMessage(http.StatusBadRequest, err)
	}

	s.logger.Infof("Receive vote request from [%s],"+
		" request server term: [%d], "+
		"request server last committed: [%d]",
		rvb.RequestID, rvb.RequestTerm, rvb.RequestNextCommitted)
	s.selfInfoDebugLog()
	// valid request
	if _, ok := s.clients[rvb.RequestID]; !ok || rvb.RequestTerm <= s.currentTerm || rvb.RequestNextCommitted < s.nextCommittedAt {
		return s.autoResult(RequestVoteResponse{
			Res: false,
		})
	}

	if s.status == CandidateStatus || s.status == LeaderStatus {
		s.status = FollowerStatus
	}
	s.currentTerm = rvb.RequestTerm
	s.voteFor = rvb.RequestID
	s.logger.Infof("Becaome follower, vote for: [%s], current term: [%d]", s.voteFor, s.currentTerm)
	s.heartbeatTicker.Reset(time.Duration(s.randomTimeoutValue) * time.Millisecond)
	return s.autoResult(RequestVoteResponse{
		Res: true,
	})
}

// ==================================== append entry rpc

type AppendEntryRPCBody struct {
	RequestID   string
	CurrentTerm uint64
}

type AppendEntryRPCResponse struct {
	Res bool
}

func (s *Server) SendAppendEntryRPC() {
	for id, clientImpl := range s.clients {
		idCopy := id
		clientImplCopy := clientImpl
		go func() {
			res, err := clientImplCopy.SendAny(AppendEntriesPath, AppendEntryRPCBody{
				RequestID:   s.ID,
				CurrentTerm: s.currentTerm,
			})
			if err != nil {
				s.logger.Errorf("request append entries to [%s] error: %v", idCopy, err)
				return
			}
			resObj := AppendEntryRPCResponse{}
			if res.ToAny(&resObj) != nil {
				s.logger.Errorf("parse [%s] append res error: %v", idCopy, err)
				return
			}
			if !resObj.Res {
				s.logger.Warnf("append to [%s] faild", idCopy)
			}
		}()
	}
}

func (s *Server) AppendEntryRPC(path string, body point.SendMessage) *point.ReceiveMessage {
	s.singleLocker.Lock()
	defer s.singleLocker.Unlock()

	requestBody := AppendEntryRPCBody{}
	if err := body.ToAny(&requestBody); err != nil {
		return point.QuickErrorReceiveMessage(http.StatusBadRequest, err)
	}

	s.logger.Infof("receive append entry request from : [%s], request server term: [%d]",
		requestBody.RequestID, requestBody.CurrentTerm)
	s.selfInfoDebugLog()
	if requestBody.RequestID != s.voteFor || requestBody.CurrentTerm < s.currentTerm {
		return s.autoResult(AppendEntryRPCResponse{
			Res: false,
		})
	}

	s.currentTerm = requestBody.CurrentTerm
	s.heartbeatTicker.Reset(time.Duration(s.randomTimeoutValue) * time.Millisecond)
	return s.autoResult(AppendEntryRPCResponse{
		Res: true,
	})
}

func (s *Server) followerLifeCycleNext() {
	s.singleLocker.Lock()
	defer s.singleLocker.Unlock()
	s.logger.Warnf("heartbeat time out, start a new leader election")
	s.status = CandidateStatus
	s.currentTerm++
	all, vote, ignoreCount := s.RequestVote()
	s.logger.Debugf("total vote member: [%d], ignore count: [%d], vote count: [%d]", all, ignoreCount, vote)
	if vote > (all-ignoreCount)/2 {
		s.logger.Info("Became leader")
		s.status = LeaderStatus
		s.SendAppendEntryRPC()
	}
}

func (s *Server) candidateLifeCycleNext() {
	s.singleLocker.Lock()
	defer s.singleLocker.Unlock()
	s.currentTerm++
	s.RandomTimeout()
	all, vote, ignoreCount := s.RequestVote()
	s.logger.Debugf("total vote member: [%d], ignore count: [%d], vote count: [%d]", all, ignoreCount, vote)
	if vote > (all-ignoreCount)/2 {
		s.logger.Info("Became leader")
		s.status = LeaderStatus
		s.SendAppendEntryRPC()
	}
}
func (s *Server) leaderLifeCycleNext() {
	s.singleLocker.Lock()
	defer s.singleLocker.Unlock()

	s.SendAppendEntryRPC()
}
