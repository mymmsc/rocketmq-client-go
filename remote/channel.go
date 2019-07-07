package remote

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/mymmsc/go-rocketmq-client/v1/executor"

	"github.com/mymmsc/go-rocketmq-client/v1/log"
)

const defaultBufferSize = 16 * 1024
const defaultReqRespBufferSize = 100

// ChannelState channel state type
type ChannelState int

const (
	// StateDisconnected disconnected status
	StateDisconnected int32 = iota
	// StateConnected connected status
	StateConnected
	// StateClosing closing status
	StateClosing
)

var errBadState = errors.New("bad state")

// ErrDisconnected disconnected error
var ErrDisconnected = errors.New("disconnected")

// Handler event handler
type Handler interface {
	OnActive(ctx *ChannelContext)
	OnDeactive(ctx *ChannelContext)
	OnClose(ctx *ChannelContext)
	OnError(ctx *ChannelContext, err error)
	OnMessage(ctx *ChannelContext, cmd *Command)
}

// ChannelConfig contains channel's configuration
type ChannelConfig struct {
	ClientConfig

	Encoder
	PacketReader
	Decoder
	Handler

	logger log.Logger
}

// ChannelContext channel's context
type ChannelContext struct {
	Address string
	Conn    net.Conn
}

func (ctx *ChannelContext) String() string {
	return "[" + ctx.Conn.LocalAddr().String() + "->" + ctx.Conn.RemoteAddr().String() + "]"
}

type runnable struct {
	c *channel
	d []byte
}

func (r runnable) Run() {
	o, err := r.c.Decode(r.d)
	if err == nil {
		r.c.OnMessage(&r.c.ctx, o)
	} else {
		r.c.OnError(&r.c.ctx, err)
	}
}

type channel struct {
	ChannelConfig

	executor *executor.GoroutinePoolExecutor
	state    int32
	ctx      ChannelContext

	exitChan chan struct{}
}

func newChannel(addr string, conf ChannelConfig) (*channel, error) {
	if conf.Decoder == nil {
		return nil, errors.New("new channel error:empty decoder")
	}
	if conf.Encoder == nil {
		return nil, errors.New("new channel error:empty encoder")
	}
	if conf.Handler == nil {
		return nil, errors.New("new channel error:empty handler")
	}

	if conf.logger == nil {
		return nil, errors.New("new channel error:empty logger")
	}

	if conf.DialTimeout <= 0 {
		conf.DialTimeout = time.Second * 3
	}

	conn, err := net.DialTimeout("tcp4", addr, conf.DialTimeout)
	if err != nil {
		return nil, fmt.Errorf("new channel dail %s, error:%s", addr, err)
	}

	executor, err := executor.NewPoolExecutor(
		"channel-executor:"+addr, 20, 20, time.Hour, executor.NewLinkedBlockingQueue(),
	)
	if err != nil {
		return nil, err
	}

	ch := &channel{
		ChannelConfig: conf,

		exitChan: make(chan struct{}),

		ctx:   ChannelContext{Address: addr, Conn: conn},
		state: StateConnected,

		executor: executor,
	}
	ch.OnActive(&ch.ctx)
	go ch.ioloop()
	return ch, nil
}

func (c *channel) ioloop() {
	var (
		zeroTime time.Time
		err      error
		d        []byte
	)
	buf := bufio.NewReaderSize(c.ctx.Conn, 1<<23)
	for {
		if atomic.LoadInt32(&c.state) != StateConnected {
			break
		}

		if c.ReadTimeout > 0 {
			c.ctx.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
		} else {
			c.ctx.Conn.SetReadDeadline(zeroTime)
		}

		d, err = c.Read(buf)
		if err == ErrNeedContent {
			c.logger.Warnf("decode not enough content:%s", err.Error())
			continue
		}

		if err != nil {
			c.logger.Errorf("decode error:%s exit loop", err)
			break
		}

		if d != nil {
			c.executor.Execute(runnable{c: c, d: d})
		}
	}
	c.close()

	if err != nil {
		if err == io.EOF {
			c.OnDeactive(&c.ctx)
		} else {
			c.OnError(&c.ctx, err)
		}
	}
}

// SendSync send data sync
func (c *channel) SendSync(cmd *Command) error {
	if c.WriteTimeout > 0 {
		c.ctx.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	} else {
		c.ctx.Conn.SetWriteDeadline(time.Time{})
	}
	bs, err := c.Encode(cmd)
	if err != nil {
		return err
	}

	if _, err := c.ctx.Conn.Write(bs); err != nil {
		c.logger.Errorf("SendSync write error:%v", err)
		c.close()
		c.OnError(&c.ctx, err)
		return err
	}
	return nil
}

func (c *channel) getState() int32 {
	return atomic.LoadInt32(&c.state)
}

func (c *channel) close() {
	if !atomic.CompareAndSwapInt32(&c.state, StateConnected, StateClosing) {
		return
	}

	c.ctx.Conn.Close()
	close(c.exitChan)
}
