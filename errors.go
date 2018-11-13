package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// All errors are defined as constants
const (
	cAlreadyRunning  = `Web server already running`
	cNoConfiguration = `Web server configuration is missing or nil`
)

// Constants are specified in objects as a fixed address all the time the application is running
// Error with error can be compared by a content, by a address, etc
var (
	errSingleton       = &Error{}
	errAlreadyRunning  = err(cAlreadyRunning)
	errNoConfiguration = err(cNoConfiguration)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors All errors of a known state that can return functions of the package
func Errors() *Error { return errSingleton }

// ERRORS:

// ErrAlreadyRunning Error: Web server already running
func ErrAlreadyRunning() error { return &errAlreadyRunning }

// ErrNoConfiguration Error: Web server configuration is missing or nil
func ErrNoConfiguration() error { return &errNoConfiguration }
