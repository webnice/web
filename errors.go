package web // import "gopkg.in/webnice/web.v1"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "fmt"

var (
	_ErrAlreadyRunning  = fmt.Errorf("Web server already running")
	_ErrNoConfiguration = fmt.Errorf("Web server configuration is missing or nil")
)

// ErrAlreadyRunning Error: Web server already running
func ErrAlreadyRunning() error { return _ErrAlreadyRunning }

// ErrNoConfiguration Error: Web server configuration is missing or nil
func ErrNoConfiguration() error { return _ErrNoConfiguration }
