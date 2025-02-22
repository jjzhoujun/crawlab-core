package process

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/sys_exec"
	"github.com/crawlab-team/go-trace"
	"math/rand"
	"os/exec"
	"time"
)

const (
	SignalCreate = iota
	SignalStart
	SignalStopped
	SignalError
	SignalExited
	SignalReachedMaxErrors
)

type Daemon struct {
	// settings
	maxErrors   int
	exitTimeout time.Duration

	// internals
	errors   int
	errMsg   string
	newCmdFn func() *exec.Cmd
	cmd      *exec.Cmd
	stopped  bool
	ch       chan int
}

func (d *Daemon) GetMaxErrors() (maxErrors int) {
	return d.maxErrors
}

func (d *Daemon) SetMaxErrors(maxErrors int) {
	d.maxErrors = maxErrors
}

func (d *Daemon) GetExitTimeout() (timeout time.Duration) {
	return d.exitTimeout
}

func (d *Daemon) SetExitTimeout(timeout time.Duration) {
	d.exitTimeout = timeout
}

func (d *Daemon) GetCmd() (cmd *exec.Cmd) {
	return d.cmd
}

func (d *Daemon) GetCh() (ch chan int) {
	return d.ch
}

func (d *Daemon) Start() (err error) {
	for {
		// command
		d.cmd = d.newCmdFn()
		d.ch <- SignalCreate

		// attempt to run
		_ = d.cmd.Start()
		d.ch <- SignalStart

		if err := d.cmd.Wait(); err != nil {
			// stopped
			d.ch <- SignalStopped
			if d.stopped {
				log.Infof("daemon stopped")
				return nil
			}

			// error
			d.ch <- SignalError
			d.errMsg = err.Error()
			trace.PrintError(err)
		}

		// exited
		d.ch <- SignalExited
		d.errMsg = errors.ErrorProcessDaemonProcessExited.Error()

		// increment errors
		d.errors++

		// validate if error count exceeds max errors
		if d.errors >= d.maxErrors {
			log.Infof("reached max errors: %d", d.maxErrors)
			d.ch <- SignalReachedMaxErrors
			return errors.ErrorProcessReachedMaxErrors
		}

		// re-attempt
		waitSec := rand.Intn(5)
		log.Infof("re-attempt to start process in %d seconds...", waitSec)
		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}

func (d *Daemon) Stop() {
	d.stopped = true
	opts := &sys_exec.KillProcessOptions{
		Timeout: d.exitTimeout,
		Force:   false,
	}
	_ = sys_exec.KillProcess(d.cmd, opts)
}

func NewProcessDaemon(newCmdFn func() *exec.Cmd, opts ...DaemonOption) (d interfaces.ProcessDaemon) {
	// daemon
	d = &Daemon{
		maxErrors:   5,
		exitTimeout: 15 * time.Second,
		errors:      0,
		errMsg:      "",
		newCmdFn:    newCmdFn,
		stopped:     false,
		ch:          make(chan int),
	}

	// apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}
