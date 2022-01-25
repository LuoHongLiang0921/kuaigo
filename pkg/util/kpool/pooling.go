// @Description 连接池

package kpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var nowFunc = time.Now // for testing

// ErrPoolExhausted is returned from a pool connection method (Do, Send,
// Receive, Flush, Err) when the maximum number of database connections in the
// pool has been reached.
var ErrPoolExhausted = errors.New("connection pool exhausted")

var (
	errPoolClosed = errors.New("connection pool closed")
	errConnClosed = errors.New("connection closed")
)

// Conn represents a connection to a server.
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value when the connection is not usable.
	// Err() error
}

// Pool maintains a pool of connections. The application calls the Get method
// to get a connection from the pool and the connection's Close method to
// return the connection's resources to the pool.
//
// Demo in example_test.go
//
type Pool struct {
	// Name of pool's unique ID
	Name string
	// Dial is an application supplied function for creating and configuring a
	// connection.
	//
	// The connection returned from Dial must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	Dial func() (Conn, error)

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c Conn, t time.Time) error

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration

	chInitialized uint32 // set to 1 when field ch is initialized

	mu     sync.Mutex    // mu protects the following fields
	closed bool          // set to true when the pool is closed.
	active int           // the number of open connections in the pool
	ch     chan struct{} // limits open connections when p.Wait is true
	idle   idleList      // idle connections
}

// ActiveCount returns the number of connections in the pool. The count
// includes idle connections and connections in use.
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

// IdleCount returns the number of idle connections in the pool.
func (p *Pool) IdleCount() int {
	p.mu.Lock()
	idle := p.idle.count
	p.mu.Unlock()
	return idle
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.active -= p.idle.count
	pc := p.idle.front
	p.idle.count = 0
	p.idle.front, p.idle.back = nil, nil
	if p.ch != nil {
		close(p.ch)
	}
	p.mu.Unlock()
	for ; pc != nil; pc = pc.next {
		pc.C.Close()
	}
	return nil
}

// PoolCtx Pool Get Context
type PoolCtx interface {
	Done() <-chan struct{}
	Err() error
}

// Get prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) Get(ctx PoolCtx) (*PoolConn, error) {

	// Handle limit for p.Wait == true.
	if p.Wait && p.MaxActive > 0 {
		p.lazyInit()
		if ctx == nil {
			<-p.ch
		} else {
			select {
			case <-p.ch:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	p.mu.Lock()

	// Prune stale connections at the back of the idle list.
	if p.IdleTimeout > 0 {
		n := p.idle.count
		for i := 0; i < n && p.idle.back != nil && p.idle.back.t.Add(p.IdleTimeout).Before(nowFunc()); i++ {
			pc := p.idle.back
			p.idle.popBack()
			p.mu.Unlock()
			pc.C.Close()
			p.mu.Lock()
			p.active--
		}
	}

	// Get idle connection from the front of idle list.
	for p.idle.front != nil {
		pc := p.idle.front
		p.idle.popFront()
		p.mu.Unlock()
		if (p.TestOnBorrow == nil || p.TestOnBorrow(pc.C, pc.t) == nil) &&
			(p.MaxConnLifetime == 0 || nowFunc().Sub(pc.created) < p.MaxConnLifetime) {
			return pc, nil
		}
		pc.C.Close()
		p.mu.Lock()
		p.active--
	}

	// Check for pool closed before dialing a new connection.
	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("get on closed pool")
	}

	// Handle limit for p.Wait == false.
	if !p.Wait && p.MaxActive > 0 && p.active >= p.MaxActive {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}

	p.active++
	p.mu.Unlock()
	c, err := p.Dial()
	if err != nil {
		c = nil
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		p.mu.Unlock()
	}
	return &PoolConn{C: c, created: nowFunc()}, err
}

func (p *Pool) lazyInit() {
	// Fast path.
	if atomic.LoadUint32(&p.chInitialized) == 1 {
		return
	}
	// Slow path.
	p.mu.Lock()
	if p.chInitialized == 0 {
		p.ch = make(chan struct{}, p.MaxActive)
		if p.closed {
			close(p.ch)
		} else {
			for i := 0; i < p.MaxActive; i++ {
				p.ch <- struct{}{}
			}
		}
		atomic.StoreUint32(&p.chInitialized, 1)
	}
	p.mu.Unlock()
}

// Put PoolConn to pool
func (p *Pool) Put(pc *PoolConn, forceClose bool) error {
	p.mu.Lock()
	if !p.closed && !forceClose {
		pc.t = nowFunc()
		p.idle.pushFront(pc)
		if p.idle.count > p.MaxIdle {
			pc = p.idle.back
			p.idle.popBack()
		} else {
			pc = nil
		}
	}

	if pc != nil {
		if pc.C != nil {
			p.mu.Unlock()
			pc.C.Close()
			p.mu.Lock()
		}
		p.active--
	}
	if p.ch != nil && !p.closed {
		p.ch <- struct{}{}
	}
	p.mu.Unlock()
	return nil
}

type idleList struct {
	count       int
	front, back *PoolConn
}

// PoolConn return PoolConn
type PoolConn struct {
	C          Conn
	t          time.Time
	created    time.Time
	next, prev *PoolConn
}

func (l *idleList) pushFront(pc *PoolConn) {
	pc.next = l.front
	pc.prev = nil
	if l.count == 0 {
		l.back = pc
	} else {
		l.front.prev = pc
	}
	l.front = pc
	l.count++
	return
}

func (l *idleList) popFront() {
	pc := l.front
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.next.prev = nil
		l.front = pc.next
	}
	pc.next, pc.prev = nil, nil
}

func (l *idleList) popBack() {
	pc := l.back
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.prev.next = nil
		l.back = pc.prev
	}
	pc.next, pc.prev = nil, nil
}
