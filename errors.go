package kevago

import "fmt"

var (
	ErrTryConnectFailed = fmt.Errorf("failed to create connections after %d tries", maxErrConnect)
)
