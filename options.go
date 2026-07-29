package loggertron

import "io"

/* Option defines a functional option for the logger
 */
type Option func(*Logger)

/* WithOutPut return a confirmation function that sets the output of logs.
 */
func WithOutput(output io.Writer) Option {
	return func(lgr *Logger) {
		lgr.output = output
	}
}

func WithFormatter(f Formatter) Option {
	return func(lgr *Logger) {
		lgr.formatter = f
	}
}
