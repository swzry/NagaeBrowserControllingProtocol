using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace nbcpClientLibv1
{
    public class RPCResult
    {
        public bool isSuccess;
        public string? ErrMsg;
        public Dictionary<string, object>? ExtraInfo;

        public RPCResult(bool success, string? ErrMsg, Dictionary<string, object>? extraData) {
            this.isSuccess = success;
            this.ErrMsg = ErrMsg;
            this.ExtraInfo = extraData;
        }

        public static RPCResult Success() {
            return new RPCResult(true, "", null);
        }

        public static RPCResult Success(Dictionary<string, object> extraData) {
            return new RPCResult(true, "", extraData);
        }

        public static RPCResult Fail(string errMsg)
        {
            return new RPCResult(false, errMsg, null);
        }

        public static RPCResult Fail(string errMsg, Dictionary<string, object> extraData)
        {
            return new RPCResult(false, errMsg, extraData);
        }

        public RPCReplyMessage ToRPCReplyMessage()
        {
            if (this.isSuccess)
            {
                if(this.ExtraInfo != null)
                {
                    return new RPCReplyMessage(this.ExtraInfo);
                }
                else
                {
                    var m = new RPCReplyMessage();
                    m.payload.Ok = true;
                    return m;
                }
            }
            else
            {
                string emsg = "";
                if(this.ErrMsg != null)
                {
                    emsg = this.ErrMsg;
                }
                if(this.ExtraInfo == null)
                {
                    return new RPCReplyMessage(emsg, new Dictionary<string, object>());
                }
                else
                {
                    return new RPCReplyMessage(emsg, this.ExtraInfo);
                }
            }
        }
    }

    public class ClientPostMessagePayload
    {
        [JsonPropertyName("topic")]
        public string Topic { get; set; }

        [JsonPropertyName("extra_data")]
        public Dictionary<string, object> ExtraData { get; set; }

        public ClientPostMessagePayload(string topic)
        {
            Topic = topic;
            ExtraData = new Dictionary<string, object>();
        }

        public ClientPostMessagePayload(string topic, Dictionary<string, object> extraData)
        {
            Topic = topic;
            ExtraData = extraData;
        }
    }

    public class ClientPostMessage
    {
        [JsonPropertyName("msg_type")]
        public string MsgType { get; set; }

        [JsonPropertyName("client_post")]
        public ClientPostMessagePayload payload { get; set; }

        public ClientPostMessage(string topic) {
            this.MsgType = "client_post";
            this.payload = new ClientPostMessagePayload(topic);
        }

        public ClientPostMessage(string topic, Dictionary<string, object> extraData) {
            this.MsgType = "client_post";
            this.payload = new ClientPostMessagePayload(topic, extraData);
        }
    }

    public class RPCReplyMsgPayload
    {
        [JsonPropertyName("ok")]
        public bool Ok { get; set; }
        [JsonPropertyName("err_msg")]
        public string ErrMsg { get; set; }
        [JsonPropertyName("extra_data")]
        public Dictionary<string, object> ExtraData { get; set; }

        public RPCReplyMsgPayload() { 
            ErrMsg = "";
            ExtraData = new Dictionary<string, object>();
        }

        public RPCReplyMsgPayload(Dictionary<string, object> extraData) {
            this.Ok = true; 
            this.ExtraData = extraData;
            this.ErrMsg = "";
        }

        public RPCReplyMsgPayload(string ErrMsg)
        {
            this.Ok = false;
            this.ErrMsg = ErrMsg;
            this.ExtraData = new Dictionary<string, object>();
        }

        public RPCReplyMsgPayload(string ErrMsg, Dictionary<string, object> extraData)
        {
            this.Ok = false;
            this.ErrMsg = ErrMsg;
            this.ExtraData = extraData;
        }
    }

    public class RPCReplyMessage
    {
        [JsonPropertyName("msg_type")]
        public string MsgType { get; set; }

        [JsonPropertyName("rpc_reply")]
        public RPCReplyMsgPayload payload { get; set; }
        public RPCReplyMessage() {
            this.MsgType = "rpc_reply";
            this.payload = new RPCReplyMsgPayload();
        }

        public RPCReplyMessage(Dictionary<string, object> extraData) {
            this.MsgType = "rpc_reply";
            this.payload = new RPCReplyMsgPayload(extraData);
        }

        public RPCReplyMessage(string ErrMsg)
        {
            this.MsgType = "rpc_reply";
            this.payload = new RPCReplyMsgPayload(ErrMsg);
        }

        public RPCReplyMessage(string ErrMsg, Dictionary<string, object> extraData)
        {
            this.MsgType = "rpc_reply";
            this.payload = new RPCReplyMsgPayload(ErrMsg, extraData);
        }
    }

    public class RPCRequestMessage
    {
        [JsonPropertyName("rpc_action")]
        public string RPCAction { get; set; }
        [JsonPropertyName("param")]
        public Dictionary<string, object> Param { get; set; }

        public RPCRequestMessage() {
            RPCAction = "";
            Param = new Dictionary<string, object>();
        }
    }

    public class SessionInMessage
    {
        [JsonPropertyName("msg_type")]
        public string MsgType { get; set; }

        [JsonPropertyName("action")]
        public string Action { get; set; }

        [JsonPropertyName("sid_info")]
        public string SIDInfo { get; set; }

        [JsonPropertyName("protocol_major_ver")]
        public int ProtocolMajorVersion { get; set; }

        [JsonPropertyName("protocol_minor_ver")]
        public int ProtocolMinorVersion { get; set; }

        public SessionInMessage()
        {
            this.MsgType = "session_in";
            this.Action = "";
            this.SIDInfo = "";
        }
    }

    public class ReplyTypeCheck
    {
        [JsonPropertyName("action_type")]
        public string ActionType { get; set; }

        public ReplyTypeCheck() {
            ActionType = "";
        }
    }

    public class NewSessionResultMessage
    {
        [JsonPropertyName("action_type")]
        public string ActionType { get; set; }

        [JsonPropertyName("reply_for")]
        public string ReplyFor { get; set; }

        [JsonPropertyName("session_info")]
        public string SessionInfo { get; set; }

        [JsonPropertyName("server_id")]
        public string ServerId { get; set; }

        [JsonPropertyName("session_id")]
        public string SessionId { get; set; }


        [JsonPropertyName("status")]
        public string Status { get; set; }

        public NewSessionResultMessage() {
            ActionType = "";
            ReplyFor = "";
            SessionId = "";
            SessionInfo = "";
            ServerId = "";
            Status = "";
        }

    }

    public class AttachSessionResultMessage
    {
        [JsonPropertyName("action_type")]
        public string ActionType { get; set; }

        [JsonPropertyName("reply_for")]
        public string ReplyFor { get; set; }

        [JsonPropertyName("status")]
        public string Status { get; set; }

        public AttachSessionResultMessage()
        {
            ActionType = "";
            ReplyFor = "";
            Status = "";
        }
    }

    public class SessionPanicMessage
    {
        [JsonPropertyName("action_type")]
        public string ActionType { get; set; }


        [JsonPropertyName("panic_type")]
        public string PanicType { get; set; }

        [JsonPropertyName("panic_msg")]
        public string PanicMsg { get; set; }

        public SessionPanicMessage()
        {
            ActionType = "";
            PanicType = "";
            PanicMsg = "";
        }

        public Exception ToException()
        {
            switch(this.PanicType)
            {
                case "internal_server_error": return new ServerInternalErrorException(this);
                case "invalid_incoming_message": return new InvalidIncomeMsgException(this);
                case "invalid_sid_info": return new InvalidSidInfoException(this);
                case "attach_previous_server_session": return new AttachPreviousServerSessionException(this);
                case "attach_with_non_int64_sid": return new AttachNonInt64SessionIDException(this);
                case "attach_session_not_found": return new AttachSessionIDNotFoundException(this);
                case "invalid_msg_type": return new InvalidMsgTypeException(this);
                case "invalid_action": return new InvalidActionException(this);
                case "protocol_major_version_not_matched": return new ProtocolMajorVersionNotMatchedException(this);
                case "protocol_minor_version_not_supported": return new ProtocolMinorVersionNotSupportedException(this);
                default: return new UnknownSessionPanicException(this);
            }
        }
    }

    public class UnknownSessionPanicException : Exception {
        public UnknownSessionPanicException(SessionPanicMessage spm) :
            base(string.Format("server reply an error (type='{0}'): {1}", spm.PanicType, spm.PanicMsg))
        { }
    }

    public class ServerInternalErrorException : Exception {
        public ServerInternalErrorException(SessionPanicMessage spm) :
            base(string.Format("server internal error: {0}", spm.PanicMsg))
        { }
    }

    public class  InvalidIncomeMsgException: Exception {
        public InvalidIncomeMsgException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  InvalidSidInfoException: Exception {
        public InvalidSidInfoException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  AttachPreviousServerSessionException: Exception {
        public AttachPreviousServerSessionException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  AttachNonInt64SessionIDException : Exception {
        public AttachNonInt64SessionIDException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  AttachSessionIDNotFoundException: Exception {
        public AttachSessionIDNotFoundException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  InvalidMsgTypeException: Exception {
        public InvalidMsgTypeException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  InvalidActionException: Exception {
        public InvalidActionException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  ProtocolMajorVersionNotMatchedException: Exception {
        public ProtocolMajorVersionNotMatchedException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class  ProtocolMinorVersionNotSupportedException: Exception {
        public ProtocolMinorVersionNotSupportedException(SessionPanicMessage spm) :
            base(string.Format("server reply an error: {0}", spm.PanicMsg))
        { }
    }

    public class InvalidReplyMessageException: Exception {
        public InvalidReplyMessageException(string msg) :
            base(string.Format("invalid reply msg: {0}", msg))
        { }
    }

    public class InvalidRPCRequestMessageException: Exception {
        public InvalidRPCRequestMessageException(string msg) :
            base(string.Format("invalid RPC msg: {0}", msg))
        { }
    }

    public class UnknownReplyMessageTypeException: Exception {
        public UnknownReplyMessageTypeException(string msgType) :
            base(string.Format("invalid reply msg type: '{0}'", msgType))
        { }
    }

    public class InvalidReplyForTypeException: Exception {
        public InvalidReplyForTypeException(string actual, string expected) :
            base(string.Format("invalid 'reply_for' field: expected '{0}', got '{1}'", expected, actual))
        { }
    }

    public class ServerReplyStatusNotOkException: Exception {
        public ServerReplyStatusNotOkException(string actual) :
            base(string.Format("field 'status' value is not 'ok': got '{0}'", actual))
        { }
    }
}
