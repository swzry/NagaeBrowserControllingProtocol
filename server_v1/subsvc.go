package server_v1

import (
	"context"
	"fmt"
	zgpf "git.swzry.com/zry/zry-go-program-framework/core"
)

var _ zgpf.ISubService = (*ZGPFSubSvc)(nil)

type ZGPFSubSvc struct {
	nbcpServer *NBCPServer
	rctx       context.Context
	cncl       context.CancelFunc
}

func NewZGPFSubSvc(nbcps *NBCPServer) *ZGPFSubSvc {
	s := &ZGPFSubSvc{
		nbcpServer: nbcps,
	}
	return s
}

func (s *ZGPFSubSvc) Prepare(ctx *zgpf.SubServiceContext) error {
	ctx.Info("preparing sub service...")
	pctx := ctx.GetParentContext()
	if pctx == nil {
		pctx = context.Background()
	}
	rctx, cncl := context.WithCancel(pctx)
	s.rctx = rctx
	s.cncl = cncl
	sessLogger := ctx.GetSubLog("session")
	s.nbcpServer.SetLogger(WrapModuleLoggerToServerLogger(sessLogger))
	return nil
}

func (s *ZGPFSubSvc) Run(ctx *zgpf.SubServiceContext) error {
	ctx.Info("starting sub service...")
	ctx.Debug("s.rctx:", s.rctx)
	return s.nbcpServer.Run(s.rctx)
}

func (s *ZGPFSubSvc) Stop(ctx *zgpf.SubServiceContext) {
	if s.cncl != nil {
		s.cncl()
	}
}

type ServerLoggerImpl struct {
	zgpf.IModuleLogger
}

func WrapModuleLoggerToServerLogger(ml zgpf.IModuleLogger) IServerLogger {
	return &ServerLoggerImpl{
		IModuleLogger: ml,
	}
}

func (i *ServerLoggerImpl) GetSessionLogger(sessionID int64) ISessionLogger {
	return i.IModuleLogger.GetSubLog(fmt.Sprintf("sid-%016x", sessionID))
}

func (i *ServerLoggerImpl) DiscardSessionLogger(sessionID int64) {
	i.IModuleLogger.CloseSubLog(fmt.Sprintf("sid-%016x", sessionID))
}
