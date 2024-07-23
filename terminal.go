package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
)

func test() error {
	fmt.Println("Begin starting bash")

	// Create arbitrary command.
	c := exec.Command("ls")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	fmt.Println("started bash")

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			// w, h, err := pty.Getsize(os.Stdin)
			// log.Printf("size is: %d , %d", w, h)
			// if err != nil {
			// 	log.Printf("error getting size: %s", err)
			// }
			// if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
			// 	log.Printf("error resizing pty: %s", err)
			// }
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// term.MakeRaw()
	// // Set stdin in raw mode.
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	// go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func runTerm() {
	if err := test(); err != nil {
		log.Fatal(err)
	}
}
