package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "fmt"

var (
	errAlreadyRunning  = fmt.Errorf("Web server already running")
	errNoConfiguration = fmt.Errorf("Web server configuration is missing or nil")
)

// ErrAlreadyRunning Error: Web server already running
func ErrAlreadyRunning() error { return errAlreadyRunning }

// ErrNoConfiguration Error: Web server configuration is missing or nil
func ErrNoConfiguration() error { return errNoConfiguration }
