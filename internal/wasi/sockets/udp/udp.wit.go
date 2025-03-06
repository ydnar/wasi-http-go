// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package udp represents the imported interface "wasi:sockets/udp@0.2.0".
package udp

import (
	"github.com/ydnar/wasi-http-go/internal/wasi/io/poll"
	"github.com/ydnar/wasi-http-go/internal/wasi/sockets/network"
	"go.bytecodealliance.org/cm"
)

// Pollable represents the imported type alias "wasi:sockets/udp@0.2.0#pollable".
//
// See [poll.Pollable] for more information.
type Pollable = poll.Pollable

// Network represents the imported type alias "wasi:sockets/udp@0.2.0#network".
//
// See [network.Network] for more information.
type Network = network.Network

// ErrorCode represents the type alias "wasi:sockets/udp@0.2.0#error-code".
//
// See [network.ErrorCode] for more information.
type ErrorCode = network.ErrorCode

// IPSocketAddress represents the type alias "wasi:sockets/udp@0.2.0#ip-socket-address".
//
// See [network.IPSocketAddress] for more information.
type IPSocketAddress = network.IPSocketAddress

// IPAddressFamily represents the type alias "wasi:sockets/udp@0.2.0#ip-address-family".
//
// See [network.IPAddressFamily] for more information.
type IPAddressFamily = network.IPAddressFamily

// IncomingDatagram represents the record "wasi:sockets/udp@0.2.0#incoming-datagram".
//
// A received datagram.
//
//	record incoming-datagram {
//		data: list<u8>,
//		remote-address: ip-socket-address,
//	}
type IncomingDatagram struct {
	_ cm.HostLayout
	// The payload.
	//
	// Theoretical max size: ~64 KiB. In practice, typically less than 1500 bytes.
	Data cm.List[uint8]

	// The source address.
	//
	// This field is guaranteed to match the remote address the stream was initialized
	// with, if any.
	//
	// Equivalent to the `src_addr` out parameter of `recvfrom`.
	RemoteAddress IPSocketAddress
}

// OutgoingDatagram represents the record "wasi:sockets/udp@0.2.0#outgoing-datagram".
//
// A datagram to be sent out.
//
//	record outgoing-datagram {
//		data: list<u8>,
//		remote-address: option<ip-socket-address>,
//	}
type OutgoingDatagram struct {
	_ cm.HostLayout
	// The payload.
	Data cm.List[uint8]

	// The destination address.
	//
	// The requirements on this field depend on how the stream was initialized:
	// - with a remote address: this field must be None or match the stream's remote address
	// exactly.
	// - without a remote address: this field is required.
	//
	// If this value is None, the send operation is equivalent to `send` in POSIX. Otherwise
	// it is equivalent to `sendto`.
	RemoteAddress cm.Option[IPSocketAddress]
}

// UDPSocket represents the imported resource "wasi:sockets/udp@0.2.0#udp-socket".
//
// A UDP socket handle.
//
//	resource udp-socket
type UDPSocket cm.Resource

// ResourceDrop represents the imported resource-drop for resource "udp-socket".
//
// Drops a resource handle.
//
//go:nosplit
func (self UDPSocket) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketResourceDrop((uint32)(self0))
	return
}

// AddressFamily represents the imported method "address-family".
//
// Whether this is a IPv4 or IPv6 socket.
//
// Equivalent to the SO_DOMAIN socket option.
//
//	address-family: func() -> ip-address-family
//
//go:nosplit
func (self UDPSocket) AddressFamily() (result IPAddressFamily) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_UDPSocketAddressFamily((uint32)(self0))
	result = (network.IPAddressFamily)((uint32)(result0))
	return
}

// FinishBind represents the imported method "finish-bind".
//
//	finish-bind: func() -> result<_, error-code>
//
//go:nosplit
func (self UDPSocket) FinishBind() (result cm.Result[ErrorCode, struct{}, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketFinishBind((uint32)(self0), &result)
	return
}

// LocalAddress represents the imported method "local-address".
//
// Get the current bound address.
//
// POSIX mentions:
// > If the socket has not been bound to a local name, the value
// > stored in the object pointed to by `address` is unspecified.
//
// WASI is stricter and requires `local-address` to return `invalid-state` when the
// socket hasn't been bound yet.
//
// # Typical errors
// - `invalid-state`: The socket is not bound to any local address.
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/getsockname.html>
// - <https://man7.org/linux/man-pages/man2/getsockname.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-getsockname>
// - <https://man.freebsd.org/cgi/man.cgi?getsockname>
//
//	local-address: func() -> result<ip-socket-address, error-code>
//
//go:nosplit
func (self UDPSocket) LocalAddress() (result cm.Result[IPSocketAddressShape, IPSocketAddress, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketLocalAddress((uint32)(self0), &result)
	return
}

// ReceiveBufferSize represents the imported method "receive-buffer-size".
//
// The kernel buffer space reserved for sends/receives on this socket.
//
// If the provided value is 0, an `invalid-argument` error is returned.
// Any other value will never cause an error, but it might be silently clamped and/or
// rounded.
// I.e. after setting a value, reading the same setting back may return a different
// value.
//
// Equivalent to the SO_RCVBUF and SO_SNDBUF socket options.
//
// # Typical errors
// - `invalid-argument`:     (set) The provided value was 0.
//
//	receive-buffer-size: func() -> result<u64, error-code>
//
//go:nosplit
func (self UDPSocket) ReceiveBufferSize() (result cm.Result[uint64, uint64, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketReceiveBufferSize((uint32)(self0), &result)
	return
}

// RemoteAddress represents the imported method "remote-address".
//
// Get the address the socket is currently streaming to.
//
// # Typical errors
// - `invalid-state`: The socket is not streaming to a specific remote address. (ENOTCONN)
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/getpeername.html>
// - <https://man7.org/linux/man-pages/man2/getpeername.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-getpeername>
// - <https://man.freebsd.org/cgi/man.cgi?query=getpeername&sektion=2&n=1>
//
//	remote-address: func() -> result<ip-socket-address, error-code>
//
//go:nosplit
func (self UDPSocket) RemoteAddress() (result cm.Result[IPSocketAddressShape, IPSocketAddress, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketRemoteAddress((uint32)(self0), &result)
	return
}

// SendBufferSize represents the imported method "send-buffer-size".
//
//	send-buffer-size: func() -> result<u64, error-code>
//
//go:nosplit
func (self UDPSocket) SendBufferSize() (result cm.Result[uint64, uint64, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketSendBufferSize((uint32)(self0), &result)
	return
}

// SetReceiveBufferSize represents the imported method "set-receive-buffer-size".
//
//	set-receive-buffer-size: func(value: u64) -> result<_, error-code>
//
//go:nosplit
func (self UDPSocket) SetReceiveBufferSize(value uint64) (result cm.Result[ErrorCode, struct{}, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	value0 := (uint64)(value)
	wasmimport_UDPSocketSetReceiveBufferSize((uint32)(self0), (uint64)(value0), &result)
	return
}

// SetSendBufferSize represents the imported method "set-send-buffer-size".
//
//	set-send-buffer-size: func(value: u64) -> result<_, error-code>
//
//go:nosplit
func (self UDPSocket) SetSendBufferSize(value uint64) (result cm.Result[ErrorCode, struct{}, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	value0 := (uint64)(value)
	wasmimport_UDPSocketSetSendBufferSize((uint32)(self0), (uint64)(value0), &result)
	return
}

// SetUnicastHopLimit represents the imported method "set-unicast-hop-limit".
//
//	set-unicast-hop-limit: func(value: u8) -> result<_, error-code>
//
//go:nosplit
func (self UDPSocket) SetUnicastHopLimit(value uint8) (result cm.Result[ErrorCode, struct{}, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	value0 := (uint32)(value)
	wasmimport_UDPSocketSetUnicastHopLimit((uint32)(self0), (uint32)(value0), &result)
	return
}

// StartBind represents the imported method "start-bind".
//
// Bind the socket to a specific network on the provided IP address and port.
//
// If the IP address is zero (`0.0.0.0` in IPv4, `::` in IPv6), it is left to the
// implementation to decide which
// network interface(s) to bind to.
// If the port is zero, the socket will be bound to a random free port.
//
// # Typical errors
// - `invalid-argument`:          The `local-address` has the wrong address family.
// (EAFNOSUPPORT, EFAULT on Windows)
// - `invalid-state`:             The socket is already bound. (EINVAL)
// - `address-in-use`:            No ephemeral ports available. (EADDRINUSE, ENOBUFS
// on Windows)
// - `address-in-use`:            Address is already in use. (EADDRINUSE)
// - `address-not-bindable`:      `local-address` is not an address that the `network`
// can bind to. (EADDRNOTAVAIL)
// - `not-in-progress`:           A `bind` operation is not in progress.
// - `would-block`:               Can't finish the operation, it is still in progress.
// (EWOULDBLOCK, EAGAIN)
//
// # Implementors note
// Unlike in POSIX, in WASI the bind operation is async. This enables
// interactive WASI hosts to inject permission prompts. Runtimes that
// don't want to make use of this ability can simply call the native
// `bind` as part of either `start-bind` or `finish-bind`.
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/bind.html>
// - <https://man7.org/linux/man-pages/man2/bind.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-bind>
// - <https://man.freebsd.org/cgi/man.cgi?query=bind&sektion=2&format=html>
//
//	start-bind: func(network: borrow<network>, local-address: ip-socket-address) ->
//	result<_, error-code>
//
//go:nosplit
func (self UDPSocket) StartBind(network_ Network, localAddress IPSocketAddress) (result cm.Result[ErrorCode, struct{}, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	network0 := cm.Reinterpret[uint32](network_)
	localAddress0, localAddress1, localAddress2, localAddress3, localAddress4, localAddress5, localAddress6, localAddress7, localAddress8, localAddress9, localAddress10, localAddress11 := lower_IPSocketAddress(localAddress)
	wasmimport_UDPSocketStartBind((uint32)(self0), (uint32)(network0), (uint32)(localAddress0), (uint32)(localAddress1), (uint32)(localAddress2), (uint32)(localAddress3), (uint32)(localAddress4), (uint32)(localAddress5), (uint32)(localAddress6), (uint32)(localAddress7), (uint32)(localAddress8), (uint32)(localAddress9), (uint32)(localAddress10), (uint32)(localAddress11), &result)
	return
}

// Stream represents the imported method "stream".
//
// Set up inbound & outbound communication channels, optionally to a specific peer.
//
// This function only changes the local socket configuration and does not generate
// any network traffic.
// On success, the `remote-address` of the socket is updated. The `local-address`
// may be updated as well,
// based on the best network path to `remote-address`.
//
// When a `remote-address` is provided, the returned streams are limited to communicating
// with that specific peer:
// - `send` can only be used to send to this destination.
// - `receive` will only return datagrams sent from the provided `remote-address`.
//
// This method may be called multiple times on the same socket to change its association,
// but
// only the most recently returned pair of streams will be operational. Implementations
// may trap if
// the streams returned by a previous invocation haven't been dropped yet before calling
// `stream` again.
//
// The POSIX equivalent in pseudo-code is:
//
//	if (was previously connected) {
//		connect(s, AF_UNSPEC)
//	}
//	if (remote_address is Some) {
//		connect(s, remote_address)
//	}
//
// Unlike in POSIX, the socket must already be explicitly bound.
//
// # Typical errors
// - `invalid-argument`:          The `remote-address` has the wrong address family.
// (EAFNOSUPPORT)
// - `invalid-argument`:          The IP address in `remote-address` is set to INADDR_ANY
// (`0.0.0.0` / `::`). (EDESTADDRREQ, EADDRNOTAVAIL)
// - `invalid-argument`:          The port in `remote-address` is set to 0. (EDESTADDRREQ,
// EADDRNOTAVAIL)
// - `invalid-state`:             The socket is not bound.
// - `address-in-use`:            Tried to perform an implicit bind, but there were
// no ephemeral ports available. (EADDRINUSE, EADDRNOTAVAIL on Linux, EAGAIN on BSD)
// - `remote-unreachable`:        The remote address is not reachable. (ECONNRESET,
// ENETRESET, EHOSTUNREACH, EHOSTDOWN, ENETUNREACH, ENETDOWN, ENONET)
// - `connection-refused`:        The connection was refused. (ECONNREFUSED)
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/connect.html>
// - <https://man7.org/linux/man-pages/man2/connect.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock2/nf-winsock2-connect>
// - <https://man.freebsd.org/cgi/man.cgi?connect>
//
//	%stream: func(remote-address: option<ip-socket-address>) -> result<tuple<incoming-datagram-stream,
//	outgoing-datagram-stream>, error-code>
//
//go:nosplit
func (self UDPSocket) Stream(remoteAddress cm.Option[IPSocketAddress]) (result cm.Result[TupleIncomingDatagramStreamOutgoingDatagramStreamShape, cm.Tuple[IncomingDatagramStream, OutgoingDatagramStream], ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	remoteAddress0, remoteAddress1, remoteAddress2, remoteAddress3, remoteAddress4, remoteAddress5, remoteAddress6, remoteAddress7, remoteAddress8, remoteAddress9, remoteAddress10, remoteAddress11, remoteAddress12 := lower_OptionIPSocketAddress(remoteAddress)
	wasmimport_UDPSocketStream((uint32)(self0), (uint32)(remoteAddress0), (uint32)(remoteAddress1), (uint32)(remoteAddress2), (uint32)(remoteAddress3), (uint32)(remoteAddress4), (uint32)(remoteAddress5), (uint32)(remoteAddress6), (uint32)(remoteAddress7), (uint32)(remoteAddress8), (uint32)(remoteAddress9), (uint32)(remoteAddress10), (uint32)(remoteAddress11), (uint32)(remoteAddress12), &result)
	return
}

// Subscribe represents the imported method "subscribe".
//
// Create a `pollable` which will resolve once the socket is ready for I/O.
//
// Note: this function is here for WASI Preview2 only.
// It's planned to be removed when `future` is natively supported in Preview3.
//
//	subscribe: func() -> pollable
//
//go:nosplit
func (self UDPSocket) Subscribe() (result Pollable) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_UDPSocketSubscribe((uint32)(self0))
	result = cm.Reinterpret[Pollable]((uint32)(result0))
	return
}

// UnicastHopLimit represents the imported method "unicast-hop-limit".
//
// Equivalent to the IP_TTL & IPV6_UNICAST_HOPS socket options.
//
// If the provided value is 0, an `invalid-argument` error is returned.
//
// # Typical errors
// - `invalid-argument`:     (set) The TTL value must be 1 or higher.
//
//	unicast-hop-limit: func() -> result<u8, error-code>
//
//go:nosplit
func (self UDPSocket) UnicastHopLimit() (result cm.Result[uint8, uint8, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_UDPSocketUnicastHopLimit((uint32)(self0), &result)
	return
}

// IncomingDatagramStream represents the imported resource "wasi:sockets/udp@0.2.0#incoming-datagram-stream".
//
//	resource incoming-datagram-stream
type IncomingDatagramStream cm.Resource

// ResourceDrop represents the imported resource-drop for resource "incoming-datagram-stream".
//
// Drops a resource handle.
//
//go:nosplit
func (self IncomingDatagramStream) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_IncomingDatagramStreamResourceDrop((uint32)(self0))
	return
}

// Receive represents the imported method "receive".
//
// Receive messages on the socket.
//
// This function attempts to receive up to `max-results` datagrams on the socket without
// blocking.
// The returned list may contain fewer elements than requested, but never more.
//
// This function returns successfully with an empty list when either:
// - `max-results` is 0, or:
// - `max-results` is greater than 0, but no results are immediately available.
// This function never returns `error(would-block)`.
//
// # Typical errors
// - `remote-unreachable`: The remote address is not reachable. (ECONNRESET, ENETRESET
// on Windows, EHOSTUNREACH, EHOSTDOWN, ENETUNREACH, ENETDOWN, ENONET)
// - `connection-refused`: The connection was refused. (ECONNREFUSED)
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/recvfrom.html>
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/recvmsg.html>
// - <https://man7.org/linux/man-pages/man2/recv.2.html>
// - <https://man7.org/linux/man-pages/man2/recvmmsg.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-recv>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-recvfrom>
// - <https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms741687(v=vs.85)>
// - <https://man.freebsd.org/cgi/man.cgi?query=recv&sektion=2>
//
//	receive: func(max-results: u64) -> result<list<incoming-datagram>, error-code>
//
//go:nosplit
func (self IncomingDatagramStream) Receive(maxResults uint64) (result cm.Result[cm.List[IncomingDatagram], cm.List[IncomingDatagram], ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	maxResults0 := (uint64)(maxResults)
	wasmimport_IncomingDatagramStreamReceive((uint32)(self0), (uint64)(maxResults0), &result)
	return
}

// Subscribe represents the imported method "subscribe".
//
// Create a `pollable` which will resolve once the stream is ready to receive again.
//
// Note: this function is here for WASI Preview2 only.
// It's planned to be removed when `future` is natively supported in Preview3.
//
//	subscribe: func() -> pollable
//
//go:nosplit
func (self IncomingDatagramStream) Subscribe() (result Pollable) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_IncomingDatagramStreamSubscribe((uint32)(self0))
	result = cm.Reinterpret[Pollable]((uint32)(result0))
	return
}

// OutgoingDatagramStream represents the imported resource "wasi:sockets/udp@0.2.0#outgoing-datagram-stream".
//
//	resource outgoing-datagram-stream
type OutgoingDatagramStream cm.Resource

// ResourceDrop represents the imported resource-drop for resource "outgoing-datagram-stream".
//
// Drops a resource handle.
//
//go:nosplit
func (self OutgoingDatagramStream) ResourceDrop() {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_OutgoingDatagramStreamResourceDrop((uint32)(self0))
	return
}

// CheckSend represents the imported method "check-send".
//
// Check readiness for sending. This function never blocks.
//
// Returns the number of datagrams permitted for the next call to `send`,
// or an error. Calling `send` with more datagrams than this function has
// permitted will trap.
//
// When this function returns ok(0), the `subscribe` pollable will
// become ready when this function will report at least ok(1), or an
// error.
//
// Never returns `would-block`.
//
//	check-send: func() -> result<u64, error-code>
//
//go:nosplit
func (self OutgoingDatagramStream) CheckSend() (result cm.Result[uint64, uint64, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	wasmimport_OutgoingDatagramStreamCheckSend((uint32)(self0), &result)
	return
}

// Send represents the imported method "send".
//
// Send messages on the socket.
//
// This function attempts to send all provided `datagrams` on the socket without blocking
// and
// returns how many messages were actually sent (or queued for sending). This function
// never
// returns `error(would-block)`. If none of the datagrams were able to be sent, `ok(0)`
// is returned.
//
// This function semantically behaves the same as iterating the `datagrams` list and
// sequentially
// sending each individual datagram until either the end of the list has been reached
// or the first error occurred.
// If at least one datagram has been sent successfully, this function never returns
// an error.
//
// If the input list is empty, the function returns `ok(0)`.
//
// Each call to `send` must be permitted by a preceding `check-send`. Implementations
// must trap if
// either `check-send` was not called or `datagrams` contains more items than `check-send`
// permitted.
//
// # Typical errors
// - `invalid-argument`:        The `remote-address` has the wrong address family.
// (EAFNOSUPPORT)
// - `invalid-argument`:        The IP address in `remote-address` is set to INADDR_ANY
// (`0.0.0.0` / `::`). (EDESTADDRREQ, EADDRNOTAVAIL)
// - `invalid-argument`:        The port in `remote-address` is set to 0. (EDESTADDRREQ,
// EADDRNOTAVAIL)
// - `invalid-argument`:        The socket is in "connected" mode and `remote-address`
// is `some` value that does not match the address passed to `stream`. (EISCONN)
// - `invalid-argument`:        The socket is not "connected" and no value for `remote-address`
// was provided. (EDESTADDRREQ)
// - `remote-unreachable`:      The remote address is not reachable. (ECONNRESET,
// ENETRESET on Windows, EHOSTUNREACH, EHOSTDOWN, ENETUNREACH, ENETDOWN, ENONET)
// - `connection-refused`:      The connection was refused. (ECONNREFUSED)
// - `datagram-too-large`:      The datagram is too large. (EMSGSIZE)
//
// # References
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/sendto.html>
// - <https://pubs.opengroup.org/onlinepubs/9699919799/functions/sendmsg.html>
// - <https://man7.org/linux/man-pages/man2/send.2.html>
// - <https://man7.org/linux/man-pages/man2/sendmmsg.2.html>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock2/nf-winsock2-send>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock2/nf-winsock2-sendto>
// - <https://learn.microsoft.com/en-us/windows/win32/api/winsock2/nf-winsock2-wsasendmsg>
// - <https://man.freebsd.org/cgi/man.cgi?query=send&sektion=2>
//
//	send: func(datagrams: list<outgoing-datagram>) -> result<u64, error-code>
//
//go:nosplit
func (self OutgoingDatagramStream) Send(datagrams cm.List[OutgoingDatagram]) (result cm.Result[uint64, uint64, ErrorCode]) {
	self0 := cm.Reinterpret[uint32](self)
	datagrams0, datagrams1 := cm.LowerList(datagrams)
	wasmimport_OutgoingDatagramStreamSend((uint32)(self0), (*OutgoingDatagram)(datagrams0), (uint32)(datagrams1), &result)
	return
}

// Subscribe represents the imported method "subscribe".
//
// Create a `pollable` which will resolve once the stream is ready to send again.
//
// Note: this function is here for WASI Preview2 only.
// It's planned to be removed when `future` is natively supported in Preview3.
//
//	subscribe: func() -> pollable
//
//go:nosplit
func (self OutgoingDatagramStream) Subscribe() (result Pollable) {
	self0 := cm.Reinterpret[uint32](self)
	result0 := wasmimport_OutgoingDatagramStreamSubscribe((uint32)(self0))
	result = cm.Reinterpret[Pollable]((uint32)(result0))
	return
}
