package versions

import (
	"fmt"
	"runtime"
)

var (
	COMMIT    = "unknown"
	VERSION   = "unknown"
	BUILDDATE = "unknown"
)

func String() string {
	return fmt.Sprintf(`
-------------------------------------------------------------------------------
network-pinger: 
  Version:       %v
  Build:         %v
  Commit:        %v
  Go Version:    %v
  Arch:          %v
-------------------------------------------------------------------------------
`, VERSION, BUILDDATE, COMMIT, runtime.Version(), runtime.GOARCH)
}
