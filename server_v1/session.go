package server_v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type INBCPSessionHandler interface {
	GetSessionKeepAliveTime() time.Duration
	GetProtocolVersion() (int, int)
	HandleSession(ctx *NBCPSessionContext) error
}

type NBCPSessionEntry struct {
	routineContext     context.Context
	routineCancel      context.CancelFunc
	sessionID          int64
	sessionContext     *NBCPSessionContext
	attached           bool
	sessionTimeoutTime time.Time
	dead               bool
	lock               sync.RWMutex
	lifetime           time.Duration
	logger             ISessionLogger
	wsconn             *websocket.Conn
	repCh              chan *RPCReplyMessage
	cmCh               chan *ClientPostMessage
	dispatchReqCh      chan *DispatchRequestMessage
	dispatchRepCh      chan *DispatchReplyMessage
	waitReattachLock   sync.Mutex
}

func NewSessionEntry(
	lifetime time.Duration,
	sessionID int64,
	parentCtx context.Context,
	logger ISessionLogger,
) *NBCPSessionEntry {
	rctx, rcncl := context.WithCancel(parentCtx)
	sctx := &NBCPSessionContext{}
	sent := &NBCPSessionEntry{
		lifetime:           lifetime,
		routineContext:     rctx,
		routineCancel:      rcncl,
		sessionID:          sessionID,
		sessionContext:     sctx,
		dead:               false,
		attached:           true,
		sessionTimeoutTime: time.Now().Add(lifetime),
		logger:             logger,
		repCh:              make(chan *RPCReplyMessage),
		cmCh:               make(chan *ClientPostMessage),
		dispatchReqCh:      make(chan *DispatchRequestMessage),
		dispatchRepCh:      make(chan *DispatchReplyMessage),
	}
	sctx.sessionEntry = sent
	return sent
}

func (e *NBCPSessionEntry) GetSessionID() int64 {
	// No need for lock, because it will not change after create
	return e.sessionID
}

func (e *NBCPSessionEntry) GetLogger() ISessionLogger {
	// No need for lock, because it will not change after create
	return e.logger
}

func (e *NBCPSessionEntry) ReplyIn(m *RPCReplyMessage) {
	defer e.lock.RUnlock()
	e.lock.RLock()
	if e.dead {
		return
	}
	e.repCh <- m
}

func (e *NBCPSessionEntry) ClientPostMsgIn(m *ClientPostMessage) {
	defer e.lock.RUnlock()
	e.lock.RLock()
	if e.dead {
		return
	}
	e.cmCh <- m
}

func (e *NBCPSessionEntry) CheckAlive() bool {
	e.lock.RLock()
	if e.dead {
		e.lock.RUnlock()
		return false
	}
	if e.attached {
		e.lock.RUnlock()
		return true
	}
	if time.Now().Before(e.sessionTimeoutTime) {
		e.lock.RUnlock()
		return true
	}
	e.lock.RUnlock()
	e.KillSession()
	e.lock.Lock()
	e.dead = true
	e.lock.Unlock()
	return false
}

func (e *NBCPSessionEntry) KillSession() {
	defer e.lock.Unlock()
	e.lock.Lock()
	if e.dead {
		return
	}
	if e.routineCancel != nil {
		e.routineCancel()
	}
	e.waitReattachLock.TryLock()
	e.waitReattachLock.Unlock()
}

func (e *NBCPSessionEntry) IsDead() bool {
	defer e.lock.RUnlock()
	e.lock.RLock()
	return e.dead
}

func (e *NBCPSessionEntry) ReAttached(wsconn *websocket.Conn) {
	defer e.lock.Unlock()
	e.lock.Lock()
	e.attached = true
	e.sessionTimeoutTime = time.Now().Add(e.lifetime)
	e.wsconn = wsconn
	e.waitReattachLock.TryLock()
	e.waitReattachLock.Unlock()
}

func (e *NBCPSessionEntry) Detached() {
	defer e.lock.Unlock()
	e.lock.Lock()
	e.attached = false
	e.sessionTimeoutTime = time.Now().Add(e.lifetime)
	e.wsconn = nil
	e.waitReattachLock.Lock()
}

func (e *NBCPSessionEntry) RunSession(sessHdl INBCPSessionHandler, wsconn *websocket.Conn) {
	e.lock.RLock()
	rctx := e.routineContext
	sctx := e.sessionContext
	e.lock.RUnlock()
	e.lock.Lock()
	e.wsconn = wsconn
	e.lock.Unlock()
	if rctx == nil {
		e.logger.VerboseF("can not start session: routine context is nil")
		return
	}
	rsch := make(chan error)
	go func() {
		rserr := sessHdl.HandleSession(sctx)
		rsch <- rserr
	}()
	err := <-rsch
	_ = wsconn.WriteJSON(&gin.H{
		"action_type": "end_session",
	})
	if err != nil {
		e.logger.Verbose("session handler end with error: ", err)
	} else {
		e.logger.Verbose("session handler end normally")
	}
	close(e.repCh)
	close(e.dispatchReqCh)
	close(e.dispatchRepCh)
	defer e.lock.Unlock()
	e.lock.Lock()
	if e.wsconn != nil {
		_ = e.wsconn.Close()
	}
	e.attached = false
	e.dead = true
}

func (e *NBCPSessionEntry) Dispatch(m *DispatchRequestMessage) (error, *DispatchReplyMessage) {
	defer e.lock.RUnlock()
	e.lock.RLock()
	if e.dead {
		return fmt.Errorf("dispatch to a dead session"), nil
	}
	e.dispatchReqCh <- m
	reply := <-e.dispatchRepCh
	return nil, reply
}

type NBCPSessionContext struct {
	sessionEntry *NBCPSessionEntry
}

func (c *NBCPSessionContext) GetDispatchRequestChannel() chan *DispatchRequestMessage {
	defer c.sessionEntry.lock.RUnlock()
	c.sessionEntry.lock.RLock()
	if c.sessionEntry.dead {
		return nil
	}
	return c.sessionEntry.dispatchReqCh
}

func (c *NBCPSessionContext) GetDispatchReplyChannel() chan *DispatchReplyMessage {
	defer c.sessionEntry.lock.RUnlock()
	c.sessionEntry.lock.RLock()
	if c.sessionEntry.dead {
		return nil
	}
	return c.sessionEntry.dispatchRepCh
}

func (c *NBCPSessionContext) GetRPCReplyChannel() chan *RPCReplyMessage {
	defer c.sessionEntry.lock.RUnlock()
	c.sessionEntry.lock.RLock()
	if c.sessionEntry.dead {
		return nil
	}
	return c.sessionEntry.repCh
}

func (c *NBCPSessionContext) GetClientPostChannel() chan *ClientPostMessage {
	defer c.sessionEntry.lock.RUnlock()
	c.sessionEntry.lock.RLock()
	if c.sessionEntry.dead {
		return nil
	}
	return c.sessionEntry.cmCh
}

func (c *NBCPSessionContext) SendRPCRequest(action string, param map[string]interface{}) error {
	defer c.sessionEntry.lock.RUnlock()
	c.sessionEntry.lock.RLock()
	if c.sessionEntry.dead {
		return fmt.Errorf("dead session")
	}
	if !c.sessionEntry.attached {
		// block here to wait attach or dead
		c.sessionEntry.waitReattachLock.Lock()
		c.sessionEntry.waitReattachLock.Unlock()
		if c.sessionEntry.dead {
			return fmt.Errorf("dead session")
		}
	}
	if c.sessionEntry.wsconn == nil {
		return fmt.Errorf("failed write to ws: wsconn is nil")
	}
	jd := &RPCRequestMessage{
		ActionType: "rpc",
		RPCAction:  action,
		Parameters: param,
	}
	err := c.sessionEntry.wsconn.WriteJSON(jd)
	if err != nil {
		return fmt.Errorf("failed write to ws: %v", err)
	}
	return nil
}

func (c *NBCPSessionContext) GetLogger() ISessionLogger {
	return c.sessionEntry.logger
}

func (c *NBCPSessionContext) GetSid() int64 {
	return c.sessionEntry.sessionID
}

func (c *NBCPSessionContext) Done() <-chan struct{} {
	return c.sessionEntry.routineContext.Done()
}

func (c *NBCPSessionContext) Kill(err error) {
	if err != nil {
		c.sessionEntry.logger.VerboseF("session killed by error: %v", err)
	} else {
		c.sessionEntry.logger.Verbose("session killed by normal sequence")
	}
	c.sessionEntry.KillSession()
}
