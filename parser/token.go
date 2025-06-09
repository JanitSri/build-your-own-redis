package parser

const (
	// RESP Data Types
	SIMPLE_STRING   = "+"
	SIMPLE_ERROR    = "-"
	INTEGER         = ":"
	BULK_STRING     = "$"
	ARRAY           = "*"
	BOOLEAN         = "#"
	DOUBLE          = ","
	BIG_NUMBER      = "("
	BULK_ERROR      = "!"
	VERBATIM_STRING = "="
	MAPS            = "%"
	ATTRIBUTES      = "|"
	SETS            = "~"
	PUSH            = ">"

	// Redis Commands
	PING        = "PING"
	ECHO        = "ECHO"
	GET         = "GET"
	SET         = "SET"
	CONFIG      = "CONFIG"
	KEYS        = "KEYS"
	INFO        = "INFO"
	REPLICATION = "REPLICATION"

	// SET COMMAND FLAGS
	PX = "PX"

	REDIS_TERMINATOR = "\r\n"
	PONG             = SIMPLE_STRING + "PONG" + REDIS_TERMINATOR
	OK               = SIMPLE_STRING + "OK" + REDIS_TERMINATOR
	NULL_BULK_STRING = BULK_STRING + "-1" + REDIS_TERMINATOR
)
