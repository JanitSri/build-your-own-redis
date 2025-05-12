package parser

import "errors"

var InvalidRespDataTypeError = errors.New("invalid redis resp type")
var InvalidRedisCommandError = errors.New("invalid redis command")
var InvalidSetCommandFlag = func(flag string) error {
	return errors.New("invalid flag for set command: " + flag)
}

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
	PING   = "PING"
	ECHO   = "ECHO"
	GET    = "GET"
	SET    = "SET"
	CONFIG = "CONFIG"

	// SET COMMAND FLAGS
	PX = "PX"

	REDIS_TERMINATOR = "\r\n"
	PONG             = SIMPLE_STRING + "PONG" + REDIS_TERMINATOR
	OK               = SIMPLE_STRING + "OK" + REDIS_TERMINATOR
	NULL_BULK_STRING = BULK_STRING + "-1" + REDIS_TERMINATOR
)
