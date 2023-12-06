package server_v1

type DispatchRequestMessage struct {
	Action     string
	Parameters interface{}
}

type DispatchReplyMessage struct {
	Error     error
	ExtraData interface{}
}
