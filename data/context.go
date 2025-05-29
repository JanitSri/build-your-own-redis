package data

type RedisContext struct {
	RedisInfo *RedisInfo
	DataStore DataStore
}

func NewRedisContext(ri *RedisInfo, ds *RedisStore) *RedisContext {
	return &RedisContext{
		ri,
		ds,
	}
}

// https://redis.io/docs/latest/commands/info/
type RedisInfo struct {
	Server      *Server
	Clients     *Clients
	Memory      *Memory
	Persistence *Persistence
	Stats       *Stats
	Replication *Replication
	CPU         *CPU
	Modules     []*Module
	Keyspace    map[string]*KeySpace
	Cluster     *Cluster
	Latency     *Latency
}

type Server struct {
	RedisVersion      string `label:"redis_version"`
	RedisMode         string `label:"redis_mode"`
	OS                string `label:"os"`
	ArchBits          int    `label:"arch_bits"`
	ProcessID         int    `label:"process_id"`
	TCPPort           int    `label:"tcp_port"`
	UptimeInSeconds   int    `label:"uptime_in_seconds"`
	UptimeInDays      int    `label:"uptime_in_days"`
	RunID             string `label:"run_id"`
	Hz                int    `label:"hz"`
	ConfigEpoch       int    `label:"config_epoch"`
	MultiplexingAPI   string `label:"multiplexing_api"`
	GCCVersion        string `label:"gcc_version"`
	ProcessSupervised string `label:"process_supervised"`
}

type Clients struct {
	ConnectedClients         int `label:"connected_clients"`
	ClusterConnections       int `label:"cluster_connections"`
	MaxInputBuffer           int `label:"max_input_buffer"`
	MaxOutputBuffer          int `label:"max_output_buffer"`
	ClientRecentMaxInputBuf  int `label:"client_recent_max_input_buf"`
	ClientRecentMaxOutputBuf int `label:"client_recent_max_output_buf"`
	BlockedClients           int `label:"blocked_clients"`
}

type Memory struct {
	UsedMemory            int64   `label:"used_memory"`
	UsedMemoryHuman       string  `label:"used_memory_human"`
	UsedMemoryRss         int64   `label:"used_memory_rss"`
	UsedMemoryPeak        int64   `label:"used_memory_peak"`
	UsedMemoryPeakHuman   string  `label:"used_memory_peak_human"`
	MemFragmentationRatio float64 `label:"mem_fragmentation_ratio"`
	MaxMemory             int64   `label:"maxmemory"`
	MaxMemoryPolicy       string  `label:"maxmemory_policy"`
	Allocator             string  `label:"allocator"`
}

type Persistence struct {
	Loading                 int `label:"loading"`
	RDBChangesSinceLastSave int `label:"rdb_changes_since_last_save"`
	RDBBgsaveInProgress     int `label:"rdb_bgsave_in_progress"`
	AofEnabled              int `label:"aof_enabled"`
	AofRewriteInProgress    int `label:"aof_rewrite_in_progress"`
	AofRewriteScheduled     int `label:"aof_rewrite_scheduled"`
}

type Stats struct {
	TotalConnectionsReceived int64 `label:"total_connections_received"`
	TotalCommandsProcessed   int64 `label:"total_commands_processed"`
	InstantaneousOpsPerSec   int   `label:"instantaneous_ops_per_sec"`
	TotalNetInputBytes       int64 `label:"total_net_input_bytes"`
	TotalNetOutputBytes      int64 `label:"total_net_output_bytes"`
	RejectedConnections      int   `label:"rejected_connections"`
	ExpiredKeys              int64 `label:"expired_keys"`
	EvictedKeys              int64 `label:"evicted_keys"`
}

type Replication struct {
	Role                       string `label:"role"`
	ConnectedSlaves            int    `label:"connected_slaves"`
	MasterReplOffset           int64  `label:"master_repl_offset"`
	ReplBacklogSize            int64  `label:"repl_backlog_size"`
	ReplBacklogFirstByteOffset int64  `label:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         int64  `label:"repl_backlog_histlen"`
}

type CPU struct {
	UsedCPUSys          float64 `label:"used_cpu_sys"`
	UsedCPUUser         float64 `label:"used_cpu_user"`
	UsedCPUSysChildren  float64 `label:"used_cpu_sys_children"`
	UsedCPUUserChildren float64 `label:"used_cpu_user_children"`
}

type Module struct {
	Name    string            `label:"name"`
	Version int               `label:"version"`
	Args    map[string]string `label:"args"`
}

type KeySpace struct {
	DB      string `label:"db"`
	Keys    int    `label:"keys"`
	Expires int    `label:"expires"`
	AvgTTL  int64  `label:"avg_ttl"`
}

type Cluster struct {
	ClusterEnabled int `label:"cluster_enabled"`
}

type Latency struct {
	LatestLatencySpikeEvent string `label:"latest_latency_spike_event"`
	LatestLatencySpikeTime  int64  `label:"latest_latency_spike_time"`
}
