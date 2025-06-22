// Package epoll provides an idiomatic Go wrapper around Linux's epoll(7) API.
//
// This module simplifies usage of epoll by wrapping unix.Epoll* functions
// into a small, clean struct interface. It offers basic operations for creating
// an epoll instance and managing file descriptors with add, modify, and delete.
//
// Unlike high-level event loops, this package exposes epoll behavior directly,
// offering full control to the user. You manage lifecycle. You own the file descriptors.
// No opinionated abstractions. Just raw, minimal tooling.
//
// Designed for those who want to write fast, event-driven servers or socket logic
// using Go's syscall layer with explicit control over epoll.
//
// Note: It is up to the caller to manage the lifetime and closure of the epoll file descriptor.
package epoll

import (
	"golang.org/x/sys/unix"
)

// The epoll instance
type Instance struct {
	Fd int
}

// Make a new epoll instance
func NewInstance(flags int) (i *Instance, err error) {
	i = &Instance{}
	i.Fd, err = unix.EpollCreate1(flags)
	return
}

// Add en entry to the interest list of the epoll file descriptor
func (i *Instance) Add(fd int, ev *unix.EpollEvent) (err error) {
	err = i.Ctl(unix.EPOLL_CTL_ADD, fd, ev)
	return
}

// Change the settings associated with [fd] in the interest list to the new settings specified in [unix.EpollEvent]
func (i *Instance) Mod(fd int, ev *unix.EpollEvent) (err error) {
	err = i.Ctl(unix.EPOLL_CTL_MOD, fd, ev)
	return
}

// Remove (deregister the target file descriptor [fd] from the interest list.
func (i *Instance) Del(fd int, ev *unix.EpollEvent) (err error) {
	err = i.Ctl(unix.EPOLL_CTL_DEL, fd, ev)
	return
}

// Control interface for an epoll file descriptor from the interest list. The [ev] event argument is ignored and can be [nil].
func (i *Instance) Ctl(op, fd int, ev *unix.EpollEvent) (err error) {
	err = unix.EpollCtl(i.Fd, op, fd, ev)
	return
}

// Wait for an I/O event on an epoll file descriptor
func (i *Instance) Wait(events []unix.EpollEvent, timeout int) (n int, err error) {
	n, err = unix.EpollWait(i.Fd, events, timeout)
	return
}

// Make an object of [EpollEvent]
func MakeEvent(fd int, events uint32) (ev *unix.EpollEvent) {
	ev = &unix.EpollEvent{
		Events: events,
		Fd:     int32(fd),
	}
	return
}
