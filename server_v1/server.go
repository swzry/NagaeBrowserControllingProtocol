package server_v1

import (
	"context"
	"fmt"
	"github.com/GUAIK-ORG/go-snowflake/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/puzpuzpuz/xsync"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NBCPServer struct {
	sessHdl    INBCPSessionHandler
	snflk      *snowflake.Snowflake
	wsUpgrader *websocket.Upgrader
	logger     IServerLogger
	sessions   *xsync.MapOf[int64, *NBCPSessionEntry]
	baseCtx    context.Context
	serverID   string
}

func NewNBCPServer(sessionHandler INBCPSessionHandler) *NBCPServer {
	wsup := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	s := &NBCPServer{
		sessHdl:    sessionHandler,
		snflk:      nil,
		wsUpgrader: wsup,
		logger:     &EmptyLogger{},
		sessions:   xsync.NewIntegerMapOf[int64, *NBCPSessionEntry](),
	}
	return s
}

func (s *NBCPServer) Run(ctx context.Context) error {
	sn, err := snowflake.NewSnowflake(0, 0)
	if err != nil {
		return fmt.Errorf("failed init snowflake generator: %v", err)
	}
	svrid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed generating server ID: %v", err)
	}
	s.serverID = svrid
	s.snflk = sn
	s.baseCtx = ctx
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.sessHdl.GetSessionKeepAliveTime() / 2):
			s.clearTimeoutSessions()
		}
	}
	return nil
}

func (s *NBCPServer) clearTimeoutSessions() {
	cleanList := make([]int64, 0)
	s.sessions.Range(func(key int64, value *NBCPSessionEntry) bool {
		if !value.CheckAlive() {
			cleanList = append(cleanList, key)
		}
		return true
	})
	for _, k := range cleanList {
		v, ok := s.sessions.Load(k)
		if ok {
			if v.IsDead() {
				s.sessions.Delete(k)
			}
		}
	}
}

func (s *NBCPServer) SetLogger(logger IServerLogger) {
	s.logger = logger
}

func (s *NBCPServer) HandleNBCPWebCommRemote(ctx *gin.Context) {
	var jdata NBCPWebRemoteRequest
	err := ctx.BindJSON(&jdata)
	if err != nil {
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: invalid json",
		})
		s.logger.Verbose("bad request: failed bind json: ", err)
		return
	}
	switch jdata.ActionType {
	case "rpc":
		s.doWebRemoteRPC(&jdata, ctx)
		return
	case "end_session":
		s.doWebRemoteEndSession(&jdata, ctx)
		return
	default:
		ctx.JSON(400, gin.H{
			"ok":          false,
			"err_msg":     "bad request: invalid action",
			"action_name": jdata.ActionType,
		})
		s.logger.Verbose("bad request: invalid action: ", jdata.ActionType)
		return
	}
}

func (s *NBCPServer) doWebRemoteRPC(jdata *NBCPWebRemoteRequest, ctx *gin.Context) {
	if jdata.RPCPayload == nil {
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: no rpc payload",
		})
		s.logger.VerboseF("bad request: no rpc payload (session_info=%s)", jdata.SessionInfo)
		return
	}
	req := &DispatchRequestMessage{
		Action:     jdata.RPCPayload.RPCActionName,
		Parameters: jdata.RPCPayload.Parameters,
	}
	sidInfoSplit := strings.Split(jdata.SessionInfo, "@")
	if len(sidInfoSplit) != 2 {
		s.logger.VerboseF("bad request: invalid session_info format '%s': can not split well", jdata.SessionInfo)
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: invalid session_info format",
		})
		return
	}
	if sidInfoSplit[1] != s.serverID {
		s.logger.VerboseF("bad request: invalid session_info '%s': server id not match", jdata.SessionInfo)
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: invalid session_info: server id not match",
		})
		return
	}
	sessid, err := strconv.ParseInt(sidInfoSplit[0], 16, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"ok":             false,
			"err_msg":        "bad request: invalid session id format: parse int64 failed",
			"session_id_raw": sidInfoSplit[0],
		})
		s.logger.VerboseF("bad request: invalid session id format: parse int64 failed: '%s'", sidInfoSplit[0])
		return
	}
	xerr, resp := s.DispatchToSession(sessid, req)
	if xerr != nil {
		ctx.JSON(500, &gin.H{
			"ok":      false,
			"err_msg": "rpc_failed",
		})
		s.logger.VerboseF("rpc failed at session_id=0x%016X, rpc_action='%s': %v", sessid, jdata.RPCPayload.RPCActionName, xerr)
		return
	}
	if resp != nil && resp.Error != nil {
		ctx.JSON(500, &gin.H{
			"ok":          false,
			"err_msg":     "rpc_error",
			"raw_rpc_err": resp.Error,
			"ext_info":    resp.ExtraData,
		})
		s.logger.VerboseF("rpc error at session_id=0x%016X, rpc_action='%s': %v", sessid, jdata.RPCPayload.RPCActionName, resp.Error)
		return
	} else {
		ctx.JSON(200, &gin.H{
			"ok": true,
		})
		return
	}
}

func (s *NBCPServer) doWebRemoteEndSession(jdata *NBCPWebRemoteRequest, ctx *gin.Context) {
	sidInfoSplit := strings.Split(jdata.SessionInfo, "@")
	if len(sidInfoSplit) != 2 {
		s.logger.VerboseF("bad request: invalid session_info format '%s': can not split well", jdata.SessionInfo)
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: invalid session_info format",
		})
		return
	}
	if sidInfoSplit[1] != s.serverID {
		s.logger.VerboseF("bad request: invalid session_info '%s': server id not match", jdata.SessionInfo)
		ctx.JSON(400, gin.H{
			"ok":      false,
			"err_msg": "bad request: invalid session_info: server id not match",
		})
		return
	}
	sessid, err := strconv.ParseInt(sidInfoSplit[0], 16, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"ok":             false,
			"err_msg":        "bad request: invalid session id format: parse int64 failed",
			"session_id_raw": sidInfoSplit[0],
		})
		s.logger.VerboseF("bad request: invalid session id format: parse int64 failed: '%s'", sidInfoSplit[0])
		return
	}
	err = s.KillSession(sessid)
	if err != nil {
		ctx.JSON(500, &gin.H{
			"ok":      false,
			"err_msg": "end session failed",
		})
		s.logger.VerboseF("end session 0x%016X failed: %v", sessid, err)
		return
	}
	ctx.JSON(200, &gin.H{
		"ok": true,
	})
	return
}

func (s *NBCPServer) HandleNBCP(ctx *gin.Context) {
	wsconn, err := s.wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		s.logger.VerboseF("failed upgrading to websocket (client=%s): %v", ctx.Request.RemoteAddr, err)
		ctx.JSON(500, gin.H{
			"err_msg": "failed upgrading to websocket",
			"status":  "ws-upgrade-error",
		})
		return
	}
	defer wsconn.Close()
	mch := make(chan *SessionInMessage)
	go func() {
		var sim SessionInMessage
		wcerr := wsconn.ReadJSON(&sim)
		if wcerr != nil {
			s.logger.VerboseF("failed handle SessionInMessage (client=%s): %v", ctx.Request.RemoteAddr, wcerr)
			_ = wsconn.WriteJSON(&gin.H{
				"action_type": "panic",
				"panic_type":  "invalid_incoming_message",
				"panic_msg":   "invalid incoming message is rejected by server",
			})
			_ = wsconn.Close()
			return
		}
		mch <- &sim
	}()
	select {
	case <-ctx.Done():
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "end_session",
		})
		return
	case simmsg := <-mch:
		pvMajor, pvMinor := s.sessHdl.GetProtocolVersion()
		xerr, ptype := simmsg.Verify(pvMajor, pvMinor)
		if xerr != nil {
			s.logger.VerboseF("failed verify SessionInMessage (client=%s): %v", ctx.Request.RemoteAddr, xerr)
			_ = wsconn.WriteJSON(&gin.H{
				"action_type": "panic",
				"panic_type":  ptype,
				"panic_msg":   xerr.Error(),
			})
			return
		}
		var sent *NBCPSessionEntry
		if simmsg.Action == "new" {
			sent = s.newSession(ctx, wsconn)
		} else {
			sent = s.attachSession(ctx, wsconn, simmsg.SIDInfo)
		}
		if sent == nil {
			return
		}
		defer sent.Detached()
		s.processIncomingMsg(sent, wsconn)
		return
	}
}

func (s *NBCPServer) newSession(ctx *gin.Context, wsconn *websocket.Conn) *NBCPSessionEntry {
	if s.baseCtx == nil {
		s.logger.ErrorF("failed handle SessionInMessage (client=%s): NBCP Server main routine not running", ctx.Request.RemoteAddr)
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "panic",
			"panic_type":  "internal_server_error",
			"panic_msg":   "NBCP server internal error",
		})
		return nil
	}
	sid := s.snflk.NextVal()
	sessLog := s.logger.GetSessionLogger(sid)
	sent := NewSessionEntry(s.sessHdl.GetSessionKeepAliveTime(), sid, s.baseCtx, sessLog)
	s.sessions.Store(sid, sent)
	sessLog.VerboseF("session created at '%s' (client=%s)", time.Now().Format(time.RFC3339Nano), ctx.Request.RemoteAddr)
	err := wsconn.WriteJSON(&gin.H{
		"action_type":  "reply",
		"reply_for":    "new_session",
		"session_id":   fmt.Sprintf("%016X", sid),
		"server_id":    s.serverID,
		"session_info": fmt.Sprintf("%016X@%s", sid, s.serverID),
		"status":       "ok",
	})
	if err != nil {
		sessLog.Verbose("new session created but failed reply to client: %v", err)
		return nil
	}
	go sent.RunSession(s.sessHdl, wsconn)
	return sent
}

func (s *NBCPServer) attachSession(ctx *gin.Context, wsconn *websocket.Conn, sidInfo string) *NBCPSessionEntry {
	sidInfoSplit := strings.Split(sidInfo, "@")
	if len(sidInfoSplit) != 2 {
		s.logger.ErrorF("failed attached to session '%s' (client=%s): invalid sid info: split failed", sidInfo, ctx.Request.RemoteAddr)
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "panic",
			"panic_type":  "invalid_sid_info",
			"panic_msg":   "invalid sid_info for attaching",
		})
		return nil
	}
	if sidInfoSplit[1] != s.serverID {
		s.logger.ErrorF("failed attached to session '%s' (client=%s): server id mot match", sidInfo, ctx.Request.RemoteAddr)
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "panic",
			"panic_type":  "attach_previous_server_session",
			"panic_msg":   "can not attach this session: server id not matched. maybe this is a session from previous server.",
		})
		return nil
	}
	sid, err := strconv.ParseInt(sidInfoSplit[0], 16, 64)
	if err != nil {
		s.logger.ErrorF("failed attached to session '%s' (client=%s): failed parse sid as int64", sidInfo, ctx.Request.RemoteAddr)
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "panic",
			"panic_type":  "attach_with_non_int64_sid",
			"panic_msg":   "can not attach this session: can not parse sid part as int64",
		})
		return nil
	}
	sent, ok := s.sessions.Load(sid)
	if !ok {
		s.logger.ErrorF("failed attached to session '%s' (client=%s): session not found", sidInfo, ctx.Request.RemoteAddr)
		_ = wsconn.WriteJSON(&gin.H{
			"action_type": "panic",
			"panic_type":  "attach_session_not_found",
			"panic_msg":   "can not attach this session: not found. maybe it is expired.",
		})
		return nil
	}
	err = wsconn.WriteJSON(&gin.H{
		"action_type": "reply",
		"reply_for":   "attach_session",
		"status":      "ok",
	})
	if err != nil {
		sent.GetLogger().VerboseF("session can attach but failed reply to client (client=%s): %v", err)
		return nil
	}
	sent.ReAttached(wsconn)
	return sent
}

func (s *NBCPServer) processIncomingMsg(sent *NBCPSessionEntry, wsconn *websocket.Conn) {
	for {
		var cm ClientWsMessage
		err := wsconn.ReadJSON(&cm)
		if err != nil {
			sent.GetLogger().Verbose("failed recv from ws: ", err)
			return
		}
		switch cm.MsgType {
		case "rpc_reply":
			if cm.RPCReplyMessage != nil {
				sent.ReplyIn(cm.RPCReplyMessage)
			} else {
				sent.GetLogger().Verbose("null field 'rpc_reply' for msg type 'rpc_reply'.")
			}
		case "client_post":
			if cm.ClientPostMsg == nil {
				sent.GetLogger().Verbose("null field 'client_msg' for msg type 'client_msg'.")
			} else {
				sent.ClientPostMsgIn(cm.ClientPostMsg)
			}
		default:
			sent.GetLogger().VerboseF("invalid client ws msg type '%s'.", cm.MsgType)
		}
	}
}

func (s *NBCPServer) DispatchToSession(sessionID int64, dm *DispatchRequestMessage) (error, *DispatchReplyMessage) {
	sent, ok := s.sessions.Load(sessionID)
	if !ok {
		return fmt.Errorf("session '0x%016X' not found", sessionID), nil
	}
	return sent.Dispatch(dm)
}

func (s *NBCPServer) KillSession(sessionID int64) error {
	sent, ok := s.sessions.Load(sessionID)
	if !ok {
		return fmt.Errorf("session '0x%016X' not found", sessionID)
	}
	sent.KillSession()
	return nil
}
