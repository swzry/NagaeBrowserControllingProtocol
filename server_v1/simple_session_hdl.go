package server_v1

import (
	"fmt"
	"time"
)

var _ INBCPSessionHandler = (*SimpleSessionHandler)(nil)

type ISimpleSessionLogic interface {
	GetSessionKeepAliveTime() time.Duration
	OnSessionStartup(sid int64, slog ISessionLogger) error
	OnClientPostMsg(sid int64, slog ISessionLogger, topic string, exData map[string]interface{})
}

type SimpleSessionHandler struct {
	logic ISimpleSessionLogic
}

func (s *SimpleSessionHandler) GetProtocolVersion() (int, int) {
	return 1, 1
}

func (s *SimpleSessionHandler) GetSessionKeepAliveTime() time.Duration {
	return s.logic.GetSessionKeepAliveTime()
}

func (s *SimpleSessionHandler) HandleSession(ctx *NBCPSessionContext) error {
	osserrCh := make(chan error)
	go func() {
		osserr := s.logic.OnSessionStartup(ctx.GetSid(), ctx.GetLogger())
		osserrCh <- osserr
	}()
SessionLoop:
	for {
		select {
		case r, ok := <-ctx.GetDispatchRequestChannel():
			if !ok {
				return fmt.Errorf("failed read from dispatch request channel")
			}
			var rpcerr error
			if r.Parameters == nil {
				rpcerr = ctx.SendRPCRequest(r.Action, map[string]interface{}{})
			} else {
				rparam, ok2 := r.Parameters.(map[string]interface{})
				if !ok2 {
					ctx.GetLogger().Verbose("invalid rpc dispatch: parameter invalid")
					continue SessionLoop
				}
				rpcerr = ctx.SendRPCRequest(r.Action, rparam)
			}
			if rpcerr != nil {
				ctx.GetLogger().VerboseF("error in rpc dispatch: %v", rpcerr)
			}
		case r, ok := <-ctx.GetRPCReplyChannel():
			if !ok {
				return fmt.Errorf("failed read from rpc reply channel")
			}
			var drep *DispatchReplyMessage
			if r.Ok {
				drep = &DispatchReplyMessage{
					Error:     nil,
					ExtraData: r.ExtraData,
				}
			} else {
				drep = &DispatchReplyMessage{
					Error:     fmt.Errorf(r.ErrMsg),
					ExtraData: r.ExtraData,
				}
			}
			go func() { ctx.GetDispatchReplyChannel() <- drep }()
		case r, ok := <-ctx.GetClientPostChannel():
			if !ok {
				return fmt.Errorf("failed read from client post channel")
			}
			go s.logic.OnClientPostMsg(ctx.GetSid(), ctx.GetLogger(), r.Topic, r.ExData)
		case <-ctx.Done():
			return nil
		case osserr := <-osserrCh:
			if osserr != nil {
				ctx.Kill(osserr)
			}
		}
	}
	return nil
}

func NewSimpleSessionHandler(logic ISimpleSessionLogic) *SimpleSessionHandler {
	s := &SimpleSessionHandler{
		logic: logic,
	}
	return s
}
