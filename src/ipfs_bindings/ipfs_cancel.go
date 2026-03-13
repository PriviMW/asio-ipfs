package main

import (
    "context"
    "log"
)

//#include <stdint.h>
import "C"

func withCancel(n *Node, cancel_signal C.uint64_t) (context.Context) {
	ctx, cancel := context.WithCancel(n.ctx)
	n.cancel_signals[cancel_signal] = cancel
	return ctx
}

//export go_asio_ipfs_cancellation_allocate
func go_asio_ipfs_cancellation_allocate() C.uint64_t {
	log.Println("go_asio_ipfs_cancellation_allocate: entry")
	n := g_node
	if n == nil {
		log.Println("go_asio_ipfs_cancellation_allocate: g_node is nil!")
		return C.uint64_t(1<<64 - 1)
	}

	ret := n.next_cancel_signal_id
	n.next_cancel_signal_id += 1
	log.Printf("go_asio_ipfs_cancellation_allocate: returning signal %d", int(ret))
	return ret
}

//export go_asio_ipfs_cancellation_free
func go_asio_ipfs_cancellation_free(cancel_signal C.uint64_t) {
	n := g_node
	if n == nil { return }
	delete(n.cancel_signals, cancel_signal)
}

//export go_asio_ipfs_cancel
func go_asio_ipfs_cancel(cancel_signal C.uint64_t) {
	n := g_node
	if n == nil { return }

	cancel, ok := n.cancel_signals[cancel_signal]
	if !ok { return }

	cancel()
}
