/*
The server status define. In the raft, there are three status of server: leader, follower and
candidate. But when the server member change, it cloud import 4th status: non-vote follower. It's a
specify status of follower, and same as 5th non-vote leader.

But the non-vote info should be store by leader. So for the cluster, the status follower is same as
follower, leader too.

The status can trans between different status. We can see non-vote follower as common follower, and
non-vote leader as common leader.
 */

const SERVER_STATUS_FOLLOWER:  u8 = 0; // Follower
const SERVER_STATUS_CANDIDATE: u8 = 1; // Candidate
const SERVER_STATUS_LEADER:    u8 = 2; // Leader

/*
The server events contain the server lifecycle: start, ready stopping and etc. Every status change
should with a event, like: shutdown-cmd. In there, the event is the server event, not the cluster
event. About the cluster event, should define in the log-entry like configuration of cluster.
 */
const SERVER_EVENT_TYPE_ELECTION_TIME_OUT: u8 = 0;
const SERVER_EVENT_TYPE_HEARTBEAT_TIMEOUT: u8 = 1;
const SERVER_EVENT_TYPE_MEMBER_CHANGE:     u8 = 2;


