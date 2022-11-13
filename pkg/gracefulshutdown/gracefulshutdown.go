// package gracefulshutdown contains graceful shutdown functions.
package gracefulshutdown

import (
	"os"
	"os/signal"
)

// InterceptOSSignals usually intercepts a set [os.Interrupt, syscall.SIGTERM, syscall.SIGKILL].
func InterceptOSSignals(interceptFn func(), sig ...os.Signal) {
	go func() {
		c := make(chan os.Signal, len(sig))
		signal.Notify(c, sig...)
		<-c
		interceptFn()
	}()
}
