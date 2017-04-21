package reaper

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/Sirupsen/logrus"
)

func ReapZombieChildren() {
	// Only reap processes if we are taking over init's duty,
	// aka we are running as pid 1 inside a docker container.
	if os.Getpid() != 1 {
		return
	}

	go func() {
		var notifications = make(chan os.Signal, 1)

		go func() {
			var sigs = make(chan os.Signal, 3)
			signal.Notify(sigs, syscall.SIGCHLD)

			for {
				var sig = <-sigs
				select {
				case notifications <- sig: /*  published it.  */
				default:
					/* Notifications channel full - drop it to the
					 * floor. This ensures we don't fill up the SIGCHLD
					 * queue. The reaper just waits for any child
					 * process (pid=-1), so we aren't loosing it!
					 */
				}
			}
		}()

		for {
			var sig = <-notifications
			logrus.Infof(" - Received signal: %v\n", sig)

			for {
				// Reap zombie
				var wstatus syscall.WaitStatus
				pid, err := syscall.Wait4(-1, &wstatus, 0, nil)
				for syscall.EINTR == err {
					pid, err = syscall.Wait4(-1, &wstatus, 0, nil)
				}

				if syscall.ECHILD == err {
					break
				}

				logrus.Infof(" - Grim reap: pid=%d, wstatus=%+v\n", pid, wstatus)
			}
		}
	}()
}
