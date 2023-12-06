package server_v1

import "fmt"

type SessionInMessage struct {
	MsgType              string `json:"msg_type"`
	Action               string `json:"action"`
	SIDInfo              string `json:"sid_info"`
	ProtocolMajorVersion int    `json:"protocol_major_ver"`
	ProtocolMinorVersion int    `json:"protocol_minor_ver"`
}

func (m *SessionInMessage) Verify(pvMajor, pvMinor int) (error, string) {
	if m.MsgType != "session_in" {
		return fmt.Errorf("failed verify SessionInMessage: invalid message sequence: msg_type 'session_in' expected, found '%s'", m.MsgType), "invalid_msg_type"
	}
	if m.Action != "new" && m.Action != "attach" {
		return fmt.Errorf("failed verify SessionInMessage: invalid message content: unknown action '%s'", m.Action), "invalid_action"
	}
	if m.ProtocolMajorVersion != pvMajor {
		return fmt.Errorf("failed verify SessionInMessage: protocol not supported: protocol major version %d(got) != %d(expected)", m.ProtocolMajorVersion, pvMajor), "protocol_major_version_not_matched"
	}
	if m.ProtocolMinorVersion < pvMinor {
		return fmt.Errorf("failed verify SessionInMessage: protocol not supported: protocol major version %d(got) < %d(minimal)", m.ProtocolMinorVersion, pvMinor), "protocol_minor_version_not_supported"
	}
	return nil, ""
}

type ClientWsMessage struct {
	// MsgType can be 'rpc_reply' or 'client_post'
	MsgType string `json:"msg_type"`
	// RPCReplyMessage if MsgType is 'rpc_reply', this will be present.
	RPCReplyMessage *RPCReplyMessage `json:"rpc_reply"`
	// ClientMessage if MsgType is 'client_msg', this will be present.
	ClientPostMsg *ClientPostMessage `json:"client_post"`
}

type ClientPostMessage struct {
	Topic  string                 `json:"topic"`
	ExData map[string]interface{} `json:"extra_data"`
}

type RPCReplyMessage struct {
	Ok        bool                   `json:"ok"`
	ErrMsg    string                 `json:"err_msg"`
	ExtraData map[string]interface{} `json:"extra_data"`
}

type RPCRequestMessage struct {
	ActionType string                 `json:"action_type"`
	RPCAction  string                 `json:"rpc_action"`
	Parameters map[string]interface{} `json:"param"`
}

type NBCPWebRemoteRequest struct {
	// ActionType Can be 'rpc' or 'end_session'
	ActionType  string `json:"action_type"`
	SessionInfo string `json:"session_info"`
	// RPCPayload when ActionType is 'rpc', this will be use.
	RPCPayload *struct {
		RPCActionName string                 `json:"rpc_action"`
		Parameters    map[string]interface{} `json:"param"`
	} `json:"rpc_payload"`
}
