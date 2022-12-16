package sraft

type Status string

const (
	FollowerStatus  Status = "follower"
	CandidateStatus Status = "candidate"
	LeaderStatus    Status = "leader"
)
