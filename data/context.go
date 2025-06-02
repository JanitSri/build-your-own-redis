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
	Replication *Replication
	CPU         *CPU
	Sentinel    *Sentinel
	Cluster     *Cluster
	Debug       *Debug
}

type Server struct {
	RedisVersion           string `label:"redis_version"`
	RedisGitSha1           string `label:"redis_git_sha1"`
	RedisGitDirty          string `label:"redis_git_dirty"`
	RedisBuildId           string `label:"redis_build_id"`
	RedisMode              string `label:"redis_mode"`
	OS                     string `label:"os"`
	ArchBits               int    `label:"arch_bits"`
	MultiplexingApi        string `label:"multiplexing_api"`
	AtomicvarApi           string `label:"atomicvar_api"`
	GccVersion             string `label:"gcc_version"`
	ProcessId              int    `label:"process_id"`
	ProcessSupervised      string `label:"process_supervised"`
	RunId                  string `label:"run_id"`
	TcpPort                int    `label:"tcp_port"`
	ServerTimeUsec         int64  `label:"server_time_usec"`
	UptimeInSeconds        int64  `label:"uptime_in_seconds"`
	UptimeInDays           int64  `label:"uptime_in_days"`
	Hz                     int    `label:"hz"`
	ConfiguredHz           int    `label:"configured_hz"`
	LruClock               int64  `label:"lru_clock"`
	Executable             string `label:"executable"`
	ConfigFile             string `label:"config_file"`
	IoThreadsActive        int    `label:"io_threads_active"`
	ShutdownInMilliseconds *int64 `label:"shutdown_in_milliseconds"`
}

type Clients struct {
	ConnectedClients            int `label:"connected_clients"`
	ClusterConnections          int `label:"cluster_connections"`
	MaxClients                  int `label:"maxclients"`
	ClientRecentMaxInputBuffer  int `label:"client_recent_max_input_buffer"`
	ClientRecentMaxOutputBuffer int `label:"client_recent_max_output_buffer"`
	BlockedClients              int `label:"blocked_clients"`
	TrackingClients             int `label:"tracking_clients"`
	PubsubClients               int `label:"pubsub_clients"`
	WatchingClients             int `label:"watching_clients"`
	ClientsInTimeoutTable       int `label:"clients_in_timeout_table"`
	TotalWatchedKeys            int `label:"total_watched_keys"`
	TotalBlockingKeys           int `label:"total_blocking_keys"`
	TotalBlockingKeysOnNoKey    int `label:"total_blocking_keys_on_nokey"`
}

type Memory struct {
	UsedMemory                      int64   `label:"used_memory"`
	UsedMemoryHuman                 string  `label:"used_memory_human"`
	UsedMemoryRss                   int64   `label:"used_memory_rss"`
	UsedMemoryRssHuman              string  `label:"used_memory_rss_human"`
	UsedMemoryPeak                  int64   `label:"used_memory_peak"`
	UsedMemoryPeakHuman             string  `label:"used_memory_peak_human"`
	UsedMemoryPeakPerc              string  `label:"used_memory_peak_perc"`
	UsedMemoryOverhead              int64   `label:"used_memory_overhead"`
	UsedMemoryStartup               int64   `label:"used_memory_startup"`
	UsedMemoryDataset               int64   `label:"used_memory_dataset"`
	UsedMemoryDatasetPerc           string  `label:"used_memory_dataset_perc"`
	TotalSystemMemory               int64   `label:"total_system_memory"`
	TotalSystemMemoryHuman          string  `label:"total_system_memory_human"`
	UsedMemoryLua                   int64   `label:"used_memory_lua"`
	UsedMemoryVmEval                int64   `label:"used_memory_vm_eval"`
	UsedMemoryLuaHuman              string  `label:"used_memory_lua_human"`
	UsedMemoryScriptsEval           int64   `label:"used_memory_scripts_eval"`
	NumberOfCachedScripts           int     `label:"number_of_cached_scripts"`
	NumberOfFunctions               int     `label:"number_of_functions"`
	NumberOfLibraries               int     `label:"number_of_libraries"`
	UsedMemoryVmFunctions           int64   `label:"used_memory_vm_functions"`
	UsedMemoryVmTotal               int64   `label:"used_memory_vm_total"`
	UsedMemoryVmTotalHuman          string  `label:"used_memory_vm_total_human"`
	UsedMemoryFunctions             int64   `label:"used_memory_functions"`
	UsedMemoryScripts               int64   `label:"used_memory_scripts"`
	UsedMemoryScriptsHuman          string  `label:"used_memory_scripts_human"`
	Maxmemory                       int64   `label:"maxmemory"`
	MaxmemoryHuman                  string  `label:"maxmemory_human"`
	MaxmemoryPolicy                 string  `label:"maxmemory_policy"`
	MemFragmentationRatio           float64 `label:"mem_fragmentation_ratio"`
	MemFragmentationBytes           int64   `label:"mem_fragmentation_bytes"`
	AllocatorFragRatio              float64 `label:"allocator_frag_ratio"`
	AllocatorFragBytes              int64   `label:"allocator_frag_bytes"`
	AllocatorRssRatio               float64 `label:"allocator_rss_ratio"`
	AllocatorRssBytes               int64   `label:"allocator_rss_bytes"`
	RssOverheadRatio                float64 `label:"rss_overhead_ratio"`
	RssOverheadBytes                int64   `label:"rss_overhead_bytes"`
	AllocatorAllocated              int64   `label:"allocator_allocated"`
	AllocatorActive                 int64   `label:"allocator_active"`
	AllocatorResident               int64   `label:"allocator_resident"`
	AllocatorMuzzy                  int64   `label:"allocator_muzzy"`
	MemNotCountedForEvict           int64   `label:"mem_not_counted_for_evict"`
	MemClientsSlaves                int64   `label:"mem_clients_slaves"`
	MemClientsNormal                int64   `label:"mem_clients_normal"`
	MemClusterLinks                 int64   `label:"mem_cluster_links"`
	MemAofBuffer                    int64   `label:"mem_aof_buffer"`
	MemReplicationBacklog           int64   `label:"mem_replication_backlog"`
	MemTotalReplicationBuffers      int64   `label:"mem_total_replication_buffers"`
	MemAllocator                    string  `label:"mem_allocator"`
	MemOverheadDbHashtableRehashing int64   `label:"mem_overhead_db_hashtable_rehashing"`
	ActiveDefragRunning             string  `label:"active_defrag_running"`
	LazyfreePendingObjects          int64   `label:"lazyfree_pending_objects"`
	LazyfreedObjects                int64   `label:"lazyfreed_objects"`
}

type Persistence struct {
	Loading                  int    `label:"loading"`
	AsyncLoading             int    `label:"async_loading"`
	CurrentCowPeak           int64  `label:"current_cow_peak"`
	CurrentCowSize           int64  `label:"current_cow_size"`
	CurrentCowSizeAge        int64  `label:"current_cow_size_age"`
	CurrentForkPerc          string `label:"current_fork_perc"`
	CurrentSaveKeysProcessed int64  `label:"current_save_keys_processed"`
	CurrentSaveKeysTotal     int64  `label:"current_save_keys_total"`
	RdbChangesSinceLastSave  int64  `label:"rdb_changes_since_last_save"`
	RdbBgsaveInProgress      int    `label:"rdb_bgsave_in_progress"`
	RdbLastSaveTime          int64  `label:"rdb_last_save_time"`
	RdbLastBgsaveStatus      string `label:"rdb_last_bgsave_status"`
	RdbLastBgsaveTimeSec     int64  `label:"rdb_last_bgsave_time_sec"`
	RdbCurrentBgsaveTimeSec  int64  `label:"rdb_current_bgsave_time_sec"`
	RdbLastCowSize           int64  `label:"rdb_last_cow_size"`
	AofEnabled               int    `label:"aof_enabled"`
	AofRewriteInProgress     int    `label:"aof_rewrite_in_progress"`
	AofRewriteScheduled      int    `label:"aof_rewrite_scheduled"`
	AofLastRewriteTimeSec    int64  `label:"aof_last_rewrite_time_sec"`
	AofCurrentRewriteTimeSec int64  `label:"aof_current_rewrite_time_sec"`
	AofLastBgrewriteStatus   string `label:"aof_last_bgrewrite_status"`
	AofLastWriteStatus       string `label:"aof_last_write_status"`
	AofLastCowSize           int64  `label:"aof_last_cow_size"`
	ModuleForkInProgress     int    `label:"module_fork_in_progress"`
	ModuleForkLastCowSize    int64  `label:"module_fork_last_cow_size"`
	AofRewrites              int64  `label:"aof_rewrites"`
	RdbSaves                 int64  `label:"rdb_saves"`
}

type Replication struct {
	Role                       string `label:"role"`
	ConnectedSlaves            int    `label:"connected_slaves"`
	MasterReplid               string `label:"master_replid"`
	MasterReplid2              string `label:"master_replid2"`
	MasterReplOffset           int64  `label:"master_repl_offset"`
	SecondReplOffset           int64  `label:"second_repl_offset"`
	ReplBacklogActive          int    `label:"repl_backlog_active"`
	ReplBacklogSize            int64  `label:"repl_backlog_size"`
	ReplBacklogFirstByteOffset int64  `label:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         int64  `label:"repl_backlog_histlen"`
}

type CPU struct {
	UsedCpuSys            float64 `label:"used_cpu_sys"`
	UsedCpuUser           float64 `label:"used_cpu_user"`
	UsedCpuSysChildren    float64 `label:"used_cpu_sys_children"`
	UsedCpuUserChildren   float64 `label:"used_cpu_user_children"`
	UsedCpuSysMainThread  float64 `label:"used_cpu_sys_main_thread"`
	UsedCpuUserMainThread float64 `label:"used_cpu_user_main_thread"`
}

type Sentinel struct {
	SentinelMasters              int `label:"sentinel_masters"`
	SentinelTilt                 int `label:"sentinel_tilt"`
	SentinelTiltSinceSeconds     int `label:"sentinel_tilt_since_seconds"`
	SentinelRunningScripts       int `label:"sentinel_running_scripts"`
	SentinelScriptsQueueLength   int `label:"sentinel_scripts_queue_length"`
	SentinelSimulateFailureFlags int `label:"sentinel_simulate_failure_flags"`
}

type Cluster struct {
	ClusterEnabled int `label:"cluster_enabled"`
}

type Debug struct {
	EventloopDurationAofSum  int64 `label:"eventloop_duration_aof_sum"`
	EventloopDurationCronSum int64 `label:"eventloop_duration_cron_sum"`
	EventloopDurationMax     int64 `label:"eventloop_duration_max"`
	EventloopCmdPerCycleMax  int64 `label:"eventloop_cmd_per_cycle_max"`
	AllocatorAllocatedLua    int64 `label:"allocator_allocated_lua"`
	AllocatorActiveLua       int64 `label:"allocator_active_lua"`
	AllocatorResidentLua     int64 `label:"allocator_resident_lua"`
	AllocatorFragBytesLua    int64 `label:"allocator_frag_bytes_lua"`
}
