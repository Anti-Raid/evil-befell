package mewld

import "time"

// The final store of the ClusterMap list as well as a instance store
type InstanceList struct {
	LastClusterStartedAt time.Time    `json:"LastClusterStartedAt"`
	Map                  []ClusterMap `json:"Map"`            // The list of clusters (ClusterMap) which defines how mewld will start clusters
	Instances            []*Instance  `json:"Instances"`      // The list of instances (Instance) which are running
	ShardCount           uint64       `json:"ShardCount"`     // The number of shards in ``mewld``
	GatewayBot           *GatewayBot  `json:"GetGatewayBot"`  // The response from Get Gateway Bot
	Dir                  string       `json:"Dir"`            // The base directory instances will use when loading clusters
	RollRestarting       bool         `json:"RollRestarting"` // whether or not we are roll restarting (rolling restart)
	FullyUp              bool         `json:"FullyUp"`        // whether or not we are fully up
}

// Represents a "cluster" of instances.
type ClusterMap struct {
	ID     int      // The clusters ID
	Name   string   // The friendly name of the cluster
	Shards []uint64 // The shard numbers/IDs of the cluster
}

type SessionStartLimit struct {
	Total          uint64 `json:"total"`           // Total number of session starts the current user is allowed
	Remaining      uint64 `json:"remaining"`       // Remaining number of session starts the current user is allowed
	ResetAfter     uint64 `json:"reset_after"`     // Number of milliseconds after which the limit resets
	MaxConcurrency uint64 `json:"max_concurrency"` // Number of identify requests allowed per 5 seconds
}

// Represents a response from the 'Get Gateway Bot' API
type GatewayBot struct {
	Url               string            `json:"url"`
	Shards            uint64            `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// Represents a instance of a cluster
type Instance struct {
	StartedAt        time.Time     `json:"StartedAt"`        // The time the instance was last started
	SessionID        string        `json:"SessionID"`        // Internally used to identify the instance
	ClusterID        int           `json:"ClusterID"`        // ClusterID from clustermap
	Shards           []uint64      `json:"Shards"`           // Shards that this instance is responsible for currently, should be equal to clustermap
	Active           bool          `json:"Active"`           // Whether or not this instance is active
	ClusterHealth    []ShardHealth `json:"ClusterHealth"`    // Cache of shard health from a ping
	CurrentlyKilling bool          `json:"CurrentlyKilling"` // Whether or not we are currently killing this instance
	LockClusterTime  *time.Time    `json:"LockClusterTime"`  // Time at which we last locked the cluster
	LaunchedFully    bool          `json:"LaunchedFully"`    // Whether or not we have launched the instance fully (till launch_next)
	LastChecked      time.Time     `json:"LastChecked"`      // The last time the shard was checked for health.
}

type ShardHealth struct {
	ShardID uint64  `json:"shard_id"` // The shard ID
	Up      bool    `json:"up"`       // Whether or not the shard is up
	Latency float64 `json:"latency"`  // Latency of the shard (optional, send if possible)
	Guilds  uint64  `json:"guilds"`   // The number of guilds in the shard. Is optional
	Users   uint64  `json:"users"`    // The number of users in the shard. Is optional
}
