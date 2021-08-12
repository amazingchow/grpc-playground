// Package grpcconnpool provides a pool of grpc client conns
package grpcconnpool

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	ErrPoolClosed        = errors.New("grpc conn pool: client conn pool is closed")
	ErrPoolTimeout       = errors.New("grpc conn pool: timed out for client conn pool to response")
	ErrConnAlreadyClosed = errors.New("grpc conn pool: grpc client conn was already closed")
	ErrFullPool          = errors.New("grpc conn pool: put a grpc client conn into a full pool")
)

// GrpcConnFactory is a function type which cares about creating a new grpc client conn.
type GrpcConnFactory func() (*grpc.ClientConn, error)

// GrpcConnFactoryWithContext is a function type which cares about creating a new grpc client conn
// that accepts the context parameter which could be passed from Get or NewWithContext method.
type GrpcConnFactoryWithContext func(context.Context) (*grpc.ClientConn, error)

// Pool is the grpc client conn pool.
type Pool struct {
	mu sync.RWMutex

	factory GrpcConnFactoryWithContext

	conns       chan CliConn
	idleTimeout time.Duration
	maxLife     time.Duration
}

// CliConn is the wrapper for a grpc client conn.
type CliConn struct {
	pool *Pool

	*grpc.ClientConn
	lastUsed  time.Time
	initAt    time.Time
	unhealthy bool
}

// New creates a new grpc client conn pool with the given initial and maximum capacity,
// and the timeout for the idle conns.
func New(factory GrpcConnFactory, init, cap int, idleTimeout time.Duration,
	maxLife ...time.Duration) (*Pool, error) {
	return NewWithContext(context.Background(), func(ctx context.Context) (*grpc.ClientConn, error) { return factory() },
		init, cap, idleTimeout, maxLife...)
}

// NewWithContext creates a new grpc client conn pool with the given initial and maximum
// capacity, and the timeout for the idle conns.
// The context parameter would be passed to the factory method during initialization.
func NewWithContext(ctx context.Context, factory GrpcConnFactoryWithContext, init, cap int, idleTimeout time.Duration,
	maxLife ...time.Duration) (*Pool, error) {

	if init < 0 {
		init = 0
	}
	if cap <= 0 {
		cap = 1
	}
	if init > cap {
		init = cap
	}

	p := &Pool{
		factory:     factory,
		conns:       make(chan CliConn, cap),
		idleTimeout: idleTimeout,
	}
	if len(maxLife) > 0 {
		p.maxLife = maxLife[0]
	}

	for i := 0; i < init; i++ {
		conn, err := factory(ctx)
		if err != nil {
			return nil, err
		}

		p.conns <- CliConn{
			pool:       p,
			ClientConn: conn,
			lastUsed:   time.Now(),
			initAt:     time.Now(),
			unhealthy:  false,
		}
	}
	for i := 0; i < cap-init; i++ {
		p.conns <- CliConn{
			pool: p,
		}
	}

	return p, nil
}

// Close empties the grpc client conn pool calling Close on all its conns.
// The conn channel is then closed, and Get will not be allowed anymore.
func (p *Pool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		if conn.ClientConn == nil {
			continue
		}
		_ = conn.ClientConn.Close()
	}
}

// IsClosed returns true if the grpc client conn pool is closed.
func (p *Pool) IsClosed() bool {
	return p == nil || p.ClientConnChan() == nil
}

// Capacity returns the grpc client conn pool capacity.
func (p *Pool) Capacity() int {
	if p.IsClosed() {
		return 0
	}
	return cap(p.conns)
}

// Available returns the number of currently unused conns.
func (p *Pool) Available() int {
	if p.IsClosed() {
		return 0
	}
	return len(p.conns)
}

func (p *Pool) ClientConnChan() chan CliConn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.conns
}

// Get will return the next available conn.
// If capacity has not been reached, it will create a new one using the factory.
// Otherwise, it will wait till the next conn becomes available or a timeout happens.
// Note that a timeout of 0 is an indefinite wait.
func (p *Pool) Get(ctx context.Context) (*CliConn, error) {
	conn := p.ClientConnChan()
	if conn == nil {
		return nil, ErrPoolClosed
	}

	wrapper := CliConn{
		pool: p,
	}
	select {
	case wrapper = <-conn:
	case <-ctx.Done():
		return nil, ErrPoolTimeout
	}

	// If the conn has stayed idle too long, close the conn and create a new
	// one. It's safe to assume that there isn't any newer conn as the conn
	// we fetched is the first in the channel.
	idleTimeout := p.idleTimeout
	if wrapper.ClientConn != nil && idleTimeout > 0 &&
		wrapper.lastUsed.Add(idleTimeout).Before(time.Now()) {

		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	}

	var err error
	if wrapper.ClientConn == nil {
		wrapper.ClientConn, err = p.factory(ctx)
		if err != nil {
			// If there was an error, we want to put back a placeholder
			// conn in the channel.
			conn <- CliConn{
				pool: p,
			}
		}
		wrapper.lastUsed = time.Now()
		wrapper.initAt = time.Now()
	} else {
		wrapper.lastUsed = time.Now()
	}

	return &wrapper, err
}

// Unhealthy marks the grpc client conn as unhealthy, so that the conn
// gets reset when closed.
func (c *CliConn) Unhealthy() {
	c.unhealthy = true
}

// Close returns a conn to the pool. It is safe to call multiple time,
// but will return an error after first time.
func (c *CliConn) Close() error {
	if c == nil {
		return nil
	}
	if c.ClientConn == nil {
		return ErrConnAlreadyClosed
	}
	if c.pool.IsClosed() {
		return ErrPoolClosed
	}

	maxLife := c.pool.maxLife
	if maxLife > 0 && c.initAt.Add(maxLife).Before(time.Now()) {
		c.Unhealthy()
	}

	wrapper := CliConn{
		pool:       c.pool,
		ClientConn: c.ClientConn,
		lastUsed:   time.Now(),
		unhealthy:  false,
	}
	if c.unhealthy {
		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	} else {
		wrapper.initAt = c.initAt
	}
	select {
	case c.pool.conns <- wrapper:
	default:
		return ErrFullPool
	}

	c.ClientConn = nil
	return nil
}
