package midi

import "fmt"

// ErrPortClosed should be returned from a driver when trying to write to a closed port.
var ErrPortClosed = fmt.Errorf("ERROR: port is closed")
