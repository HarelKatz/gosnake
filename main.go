package main

import (
	"flag"
	"fmt"
	output "gosnake/helpers"
	"os"
	"os/signal"
	"syscall"
)

// check if ran with apropriate permissions, changes UID and GID if so
func runAsRoot() error {

	if os.Geteuid() != 0 {
		return fmt.Errorf("cannot set UID/GID to root: not running as root")
	}

	// Setgid before Setuid (recommended)
	if err := syscall.Setuid(0); err != nil {
		fmt.Fprintf(output.Writer, "Setuid failed: %v\n", err)
		return err

	}
	if err := syscall.Setgid(0); err != nil {
		fmt.Fprintf(output.Writer, "Setgid failed: %v\n", err)
		return err
	}
	return nil
}

func daemonize() {
	// TODO
}

func main() {

	// Create a channel to receive OS signals
	recivedSignals := make(chan os.Signal, 1)

	// Catch all relevant signals
	signal.Notify(recivedSignals,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // kill, systemd stop
		syscall.SIGHUP,  // terminal hangup or reload
		syscall.SIGQUIT, // Ctrl+\
		syscall.SIGSEGV, // segfault (rarely caught in Go)
		syscall.SIGBUS,  // bus error
		syscall.SIGILL,  // illegal instruction
		syscall.SIGCHLD, // child process change
	)

	// a channel for the signals reciver gorutine to exit
	isGoroutineDone := make(chan bool, 1)

	// goroutine to handle signals
	go func() {
		sig := <-recivedSignals
		fmt.Fprintf(output.Writer, "Received signal: %s\n", sig)
		isGoroutineDone <- true
	}()

	// argv stuff
	toDaemonize := flag.Bool("daemonize", false, "If daemonized")
	logFile := flag.String("logFile", "stdout", "Log file location")

	flag.Parse()

	// Initialize the writer (stdout or file)
	err := output.Init(*logFile)
	if err != nil {
		panic(err)
	}

	// check if ran with apropriate permissions, changes UID and GID if so
	if err := runAsRoot(); err != nil {
		fmt.Fprintf(output.Writer, "%v\n", err)
		panic(err)
	}

	if *toDaemonize {
		fmt.Fprintf(output.Writer, "Daemonizing...\n")
		daemonize()

	} else {

		<-isGoroutineDone
	}

}
