using System.Text.Json.Nodes;
using Websocket.Client;
using System.Reflection;
using System.Net.WebSockets;
using System.Text.Json;

[assembly: AssemblyVersion("1.0.*")]

namespace nbcpClientLibv1
{
    public class NBCPClient
    {
        private Uri serverUri;
        private int majorPV, minorPV;
        private bool isConnected;
        private NBCPLogicHandler logicHandler;
        private Thread workThread;
        private bool isWorking;
        private WebsocketClient? wsClient;
        private ManualResetEvent mreWorkThread = new ManualResetEvent(false);
        private string sessionInfo = "";
        private string sessionId = "";
        private string serverId = "";
        private bool waitForSIMReply = false;
        private bool isWaitForNewSessionReply = false;

        public delegate void ErrorHandlerDelegate(string errLocation, Exception exception);
        public delegate void InfoLogHandlerDelegate(string location, string msg);
        public delegate void ClientExitDelegate();
        public event ErrorHandlerDelegate? ErrorHandler;
        public event InfoLogHandlerDelegate? InfoLogHandler;
        public event ClientExitDelegate? OnClientExit;

        public NBCPClient(string server, int majorProtocolVer, int minorProtocolVer, NBCPLogicHandler handler) {
            this.serverUri = new Uri(server);
            this.majorPV = majorProtocolVer;
            this.minorPV = minorProtocolVer;
            this.logicHandler = handler;
            this.workThread = new Thread(this.workThreadProc);
        }

        private void workThreadProc() {
            InfoLogHandler?.Invoke("work-thread", "work thread starting.");
            while (isWorking)
            {
                sessionProc();
                if (!isWorking)
                {
                    break;
                }
                bool willReAttach = this.logicHandler.WillReAttach();
                if (!willReAttach)
                {
                    this.sessionInfo = "";
                    this.sessionId = "";
                    this.serverId = "";
                    break;
                }
                else
                {
                    InfoLogHandler?.Invoke("work-thread", "doing reattaching...");
                }
            }
            InfoLogHandler?.Invoke("work-thread", "work thread end.");
            OnClientExit?.Invoke();
        }

        private void sessionProc()
        {
            wsClient = new WebsocketClient(this.serverUri);
            wsClient.IsReconnectionEnabled = false;
            wsClient.MessageReceived.Subscribe(msg =>
            {
                wsMsgProcess(msg);
            });
            wsClient.DisconnectionHappened.Subscribe(msg =>
            {
                InfoLogHandler?.Invoke("ws-client", "connection end.");
                mreWorkThread.Set();
                wsClient = null;
            });
            wsClient.ReconnectionHappened.Subscribe(msg => {
                InfoLogHandler?.Invoke("ws-client", "connected.");
                waitForSIMReply = true;
                if (sessionInfo == "")
                {
                    sendNewSessionRequest();
                }
                else
                {
                    sendAttachSessionRequest();
                }
            });
            mreWorkThread.Reset();
            try
            {
                InfoLogHandler?.Invoke("ws-client", string.Format("connecting to '{0}'...", this.serverUri.ToString()));
                wsClient.Start();
            }catch (Exception ex)
            {
                this.ErrorHandler?.Invoke("ws_connect_error", ex);
                return;
            }
            mreWorkThread.WaitOne();
        }

        private void sendNewSessionRequest()
        {
            InfoLogHandler?.Invoke("session-mgr", "apply for new session...");
            SessionInMessage sim = new SessionInMessage();
            sim.Action = "new";
            sim.ProtocolMajorVersion = this.majorPV;
            sim.ProtocolMinorVersion = this.minorPV;
            sim.SIDInfo = "";
            this.isWaitForNewSessionReply = true;
            this.sendSIM(sim);
        }

        private void sendAttachSessionRequest()
        {
            InfoLogHandler?.Invoke("session-mgr", string.Format("attach for session '{0}'...", this.sessionInfo));
            SessionInMessage sim = new SessionInMessage();
            sim.Action = "attach";
            sim.ProtocolMajorVersion = this.majorPV;
            sim.ProtocolMinorVersion = this.minorPV;
            sim.SIDInfo = this.sessionInfo;
            this.isWaitForNewSessionReply = false;
            this.sendSIM(sim);
        }

        private void sendSIM(SessionInMessage sim) {
            string jsonStr;
            try {
                jsonStr = System.Text.Json.JsonSerializer.Serialize(sim);
                InfoLogHandler?.Invoke("session-in-msg/serialize", string.Format("result: {0}.", jsonStr));
            }catch (Exception ex){
                Fail("json-serde/sim", ex);
                return;
            }
            if(wsClient != null)
            {
                wsClient.Send(jsonStr);
            }
        }

        private void wsMsgProcess(ResponseMessage msg) {
            string repType = analyzeReplyType(msg.Text);
            InfoLogHandler?.Invoke("debug/ws-msg-in/pre-proc", string.Format("type='{0}', raw='{1}'", repType, msg.Text));
            switch (repType)
            {
                case "": return;
                case "panic": processPanicMessage(msg.Text); return;
                case "reply": break;
                case "rpc": break;
                case "end_session": this.Stop(); return; 
                default: Fail("proc-msg/type-check", new UnknownReplyMessageTypeException(repType)); return;
            }
            if(waitForSIMReply)
            {
                InfoLogHandler?.Invoke("debug/ws-msg-in/sim", msg.Text);
                if (this.isWaitForNewSessionReply)
                {
                    processNewSessionReply(msg.Text);
                }
                else
                {
                    processAttachSessionReply(msg.Text);
                }
            }
            else
            {
                if(repType == "rpc")
                {
                    InfoLogHandler?.Invoke("debug/ws-msg-in/rpc", msg.Text);
                    ProcessRPC(msg.Text);
                }
            }
        }

        private void processPanicMessage(string payload) {
            try
            {
                var res = System.Text.Json.JsonSerializer.Deserialize<SessionPanicMessage>(payload);
                if(res == null)
                {
                    Fail("proc-msg/panic-proc", new InvalidReplyMessageException(payload));
                    return;
                }
                Fail("proc-msg/panic-proc", res.ToException());
                return;
            }catch(Exception ex)
            {
                Fail("proc-msg/panic-proc", ex);
                return;
            }
        }

        private string analyzeReplyType(string payload) {
            try
            {
                var res = System.Text.Json.JsonSerializer.Deserialize<ReplyTypeCheck>(payload);
                if(res == null)
                {
                    Fail("reply-type-check", new InvalidReplyMessageException(payload));
                    return "";
                }
                return res.ActionType;
            }catch(Exception ex)
            {
                Fail("reply-type-check", ex);
                return "";
            }
        }

        private void processNewSessionReply(string payload)
        {
            try
            {
                var res = System.Text.Json.JsonSerializer.Deserialize<NewSessionResultMessage>(payload);
                if(res == null)
                {
                    Fail("proc-sim-reply/new", new InvalidReplyMessageException(payload));
                    return;
                }
                if(res.ReplyFor != "new_session")
                {
                    Fail("proc-sim-reply/new", new InvalidReplyForTypeException(res.ReplyFor, "new_session"));
                    return;
                }
                if(res.Status != "ok")
                {
                    Fail("proc-sim-reply/new", new ServerReplyStatusNotOkException(res.Status));
                    return;
                }
                this.sessionInfo = res.SessionInfo;
                this.sessionId = res.SessionId;
                this.serverId = res.ServerId;
                this.waitForSIMReply = false;
                InfoLogHandler?.Invoke("proc-sim-reply/new", string.Format("new session 0x{0:X16}@{1} (Raw='{2}').", this.sessionId, this.serverId, this.sessionInfo));
            }catch(Exception ex){
                Fail("proc-sim-reply/new", ex);
                return;
            }
        }

        private void processAttachSessionReply(string payload)
        {
            try
            {
                var res = System.Text.Json.JsonSerializer.Deserialize<AttachSessionResultMessage>(payload);
                if(res == null)
                {
                    Fail("proc-sim-reply/attach", new InvalidReplyMessageException(payload));
                    return;
                }
                if(res.ReplyFor != "attach_session")
                {
                    Fail("proc-sim-reply/attach", new InvalidReplyForTypeException(res.ReplyFor, "new_session"));
                    return;
                }
                if(res.Status != "ok")
                {
                    Fail("proc-sim-reply/attach", new ServerReplyStatusNotOkException(res.Status));
                    return;
                }
                this.waitForSIMReply = false;
                InfoLogHandler?.Invoke("proc-sim-reply/attach", string.Format("attached session 0x{0:X16}@{1}.", this.sessionId, this.serverId));
            }catch(Exception ex){
                Fail("proc-sim-reply/attach", ex);
                return;
            }
        }

        private void Fail(string loc, WebSocketCloseStatus code, Exception ex)
        {
            ErrorHandler?.Invoke(loc, ex);
            if(wsClient != null)
            {
                mreWorkThread.Set();
                wsClient.Stop(code, string.Format("ws closed by failure (location: {0}): {1}", loc, ex.Message));
            }
        }

        private static object ConvertJsonValue(JsonElement jsonElement)
        {
            switch (jsonElement.ValueKind)
            {
                case JsonValueKind.Undefined:
                case JsonValueKind.Null:
                    return null;

                case JsonValueKind.Number:
                    return jsonElement.GetDouble();

                case JsonValueKind.String:
                    string? str = jsonElement.GetString();
                    if (str == null)
                    {
                        return "";
                    }
                    else
                    {
                        return str;
                    }

                case JsonValueKind.True:
                    return true;

                case JsonValueKind.False:
                    return false;

                case JsonValueKind.Object:
                    // Recursively convert nested objects
                    return jsonElement.EnumerateObject()
                        .ToDictionary(property => property.Name, property => ConvertJsonValue(property.Value));

                case JsonValueKind.Array:
                    // Recursively convert array elements
                    return jsonElement.EnumerateArray().Select(ConvertJsonValue).ToList();

                default:
                    throw new NotSupportedException($"Unsupported JsonValueKind: {jsonElement.ValueKind}");
            }
        }

        private void ProcessRPC(string payload)
        {
            try
            {
                var res = System.Text.Json.JsonSerializer.Deserialize<RPCRequestMessage>(payload);
                if(res == null)
                {
                    Fail("proc-rpc", new InvalidRPCRequestMessageException(payload));
                    return;
                }
                Dictionary<string, object> paramFin = new Dictionary<string, object>();
                if (res.Param != null)
                {
                    foreach(KeyValuePair<string, object> kvp in res.Param)
                    {
                        if(kvp.Value is JsonElement)
                        {
                            paramFin[kvp.Key] = ConvertJsonValue((JsonElement)kvp.Value);
                        }
                        else
                        {
                            paramFin[kvp.Key] = kvp.Value;
                        }
                    }
                }
                var rpcRes = this.logicHandler.RPCHandler(res.RPCAction, paramFin);
                if(rpcRes != null)
                {
                    InfoLogHandler?.Invoke("process-rpc/result", string.Format("isSuccess: {0}, ErrMsg: {1}.", rpcRes.isSuccess, rpcRes.ErrMsg));
                }
                else
                {
                    InfoLogHandler?.Invoke("process-rpc/result", "null result");
                }
                RPCReplyMessage rpcReplyMsg;
                if(rpcRes != null)
                {
                    rpcReplyMsg = rpcRes.ToRPCReplyMessage();
                }
                else
                {
                    rpcReplyMsg = new RPCReplyMessage();
                    rpcReplyMsg.payload.Ok = true;
                }
                InfoLogHandler?.Invoke("process-rpc/rpc-reply", string.Format("msg_type: {0}, rpc_reply.ok: {1}, rpc_reply.err_msg: {2}.", rpcReplyMsg.MsgType, rpcReplyMsg.payload.Ok, rpcReplyMsg.payload.ErrMsg));
                SendRPCReply(rpcReplyMsg);
            }catch(Exception ex){
                Fail("proc-rpc", ex);
                return;
            }
        }

        private void SendRPCReply(RPCReplyMessage rpm) {
            string jsonStr;
            try {
                jsonStr = System.Text.Json.JsonSerializer.Serialize(rpm);
                InfoLogHandler?.Invoke("rpc-reply/serialize", string.Format("result: {0}.", jsonStr));
            }catch (Exception ex){
                Fail("json-serde/rpc-reply", ex);
                return;
            }
            if(wsClient != null)
            {
                wsClient.Send(jsonStr);
            }
        }

        public void SendClientPostMsg(ClientPostMessage cpm) {
            string jsonStr = "";
            try {
                jsonStr = System.Text.Json.JsonSerializer.Serialize(cpm);
                InfoLogHandler?.Invoke("client-post/serialize", string.Format("result: {0}.", jsonStr));
            }catch (Exception ex){
                Fail("json-serde/client-post-msg", ex);
                return;
            }
            if(wsClient != null)
            {
                wsClient.Send(jsonStr);
            }
        }

        private void Fail(string loc, Exception ex)
        {
            ErrorHandler?.Invoke(loc, ex);
            if(wsClient != null)
            {
                mreWorkThread.Set();
                wsClient.Stop(WebSocketCloseStatus.InvalidPayloadData, string.Format("ws closed by failure (location: {0}): {1}", loc, ex.Message));
            }
        }

        public string GetSessionId()
        {
            return this.sessionId;
        }

        public string GetServerId()
        {
            return this.serverId;
        }

        public string GetSessionInfo()
        {
            return this.sessionInfo;
        }

        public void Stop()
        {
            isWorking = false;
            if(wsClient != null)
            {
                mreWorkThread.Set();
                wsClient = null;
            }
        }

        public void Start()
        {
            isWorking =true;
            this.workThread.Start();
        }
    }
}