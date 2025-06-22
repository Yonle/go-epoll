## go-epoll [![Go Reference](https://pkg.go.dev/badge/github.com/Yonle/go-epoll.svg)](https://pkg.go.dev/github.com/Yonle/go-epoll)

an idiomatic go [`epoll(7)`](https://man7.org/linux/man-pages/man7/epoll.7.html) wrapper.

## example
```go
package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"
	"github.com/Yonle/go-epoll"
)

func main() {
	// Open hello.txt in non-blocking read-only mode
	fd, err := unix.Open("hello.txt", unix.O_RDONLY|unix.O_NONBLOCK, 0)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer unix.Close(fd)
	fmt.Printf("Opened hello.txt with fd %d\n", fd)

	// Create an epoll instance
	epi, err := epoll.NewInstance(0)
	if err != nil {
		log.Fatalf("failed to create epoll instance: %v", err)
	}
	defer unix.Close(epi.Fd)

	// Register the file for EPOLLIN events
	ev := epoll.MakeEvent(fd, unix.EPOLLIN)
	if err := epi.Add(fd, ev); err != nil {
		log.Fatalf("failed to add fd to epoll: %v", err)
	}

	// Wait for event
	events := make([]unix.EpollEvent, 1)
	fmt.Println("Waiting for events on hello.txt...")

	n, err := epi.Wait(events, -1)
	if err != nil {
		log.Fatalf("epoll_wait error: %v", err)
	}

	for i := 0; i < n; i++ {
		if events[i].Fd == int32(fd) && events[i].Events&unix.EPOLLIN != 0 {
			buf := make([]byte, 1024)
			n, err := unix.Read(fd, buf)
			if err != nil {
				log.Fatalf("read error: %v", err)
			}
			fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
		}
	}
})
```
