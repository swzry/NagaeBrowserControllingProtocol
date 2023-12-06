package main

import (
	"fmt"
	nbcp "git.swzry.com/ProjectNagae/NagaeBrowserControllingProtocol/server_v1"
	"time"
)

var _ nbcp.ISimpleSessionLogic = (*SessionLogicClass)(nil)

var AllWinClosed = fmt.Errorf("all browser windows closed")

type SessionLogicClass struct {
}

func (l *SessionLogicClass) OnClientPostMsg(sid int64, slog nbcp.ISessionLogger, topic string, exData map[string]interface{}) {
	switch topic {
	case "webview2_core_loaded":
		slog.Info("webview2 core loaded.")
		//_ = l.DoRPCUtil(sid, slog, "open_dev_tool",
		//	map[string]interface{}{
		//		"name": "test1",
		//	})
	case "all_window_closed":
		slog.Info("all browser windows closed. session will end.")
		l.EndSession(sid, slog)
	}
}

func NewSessionLogic() *SessionLogicClass {
	l := &SessionLogicClass{}
	return l
}

func (l *SessionLogicClass) GetSessionKeepAliveTime() time.Duration {
	return time.Minute
}

func (l *SessionLogicClass) DoRPCUtil(sid int64, slog nbcp.ISessionLogger, fn string, param map[string]interface{}) error {
	err, drep := NBCPServer.DispatchToSession(sid, &nbcp.DispatchRequestMessage{
		Action:     fn,
		Parameters: param,
	})
	if err != nil {
		slog.WarnF("failed dispatch '%s': %v", fn, err)
		return fmt.Errorf("failed dispatch '%s': %v", fn, err)
	}
	if drep != nil && drep.Error != nil {
		slog.WarnF("error in '%s': %v", fn, drep.Error)
		return fmt.Errorf("error in '%s': %v", fn, drep.Error)
	}
	return nil
}

func (l *SessionLogicClass) EndSession(sid int64, slog nbcp.ISessionLogger) {
	err := NBCPServer.KillSession(sid)
	if err != nil {
		slog.Warn("failed end session.")
	}
}

func (l *SessionLogicClass) OnSessionStartup(sid int64, slog nbcp.ISessionLogger) error {
	err := l.DoRPCUtil(sid, slog, "new_window",
		map[string]interface{}{
			"name":  "test1",
			"title": "TEST HOME",
			"url":   Config.NBCP.HomeURL,
		})
	if err != nil {
		return err
	}
	return nil
}
