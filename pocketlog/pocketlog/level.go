package pocketlog

/* Level represents an available logging level
*/
type Level byte
/* byte (an alias for uint8) for the level type to be memory-efficient, as we only need to store small numbers
*/

const (
/* LevelDebug represents the lowest level of log, mostly used for debugging purposes
*/
	LevelDebug Level = iota
/* LevelInfo represents the logging level for valuable insights.
*/
	LevelInfo
/* LevelError represents the highest logging level for tracing errors.
*/
	LevelError
)

/* String() method returns the string representation of the logging level.
*/
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "[DEBUG]"
	case LevelInfo:
		return "[INFO]"
	case LevelError:
		return "[ERROR]"
	default:
		return "[UNKNOWN]"
	}
}