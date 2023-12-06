import axios from "axios"

class NBCPWebComm {
  constructor(sessionInfo, windowName, webRemoteUrl) {
    this.sessionInfo = sessionInfo
    this.windowName = windowName
    this.webRemoteUrl = webRemoteUrl
    this.axios = axios
  }

  test() {
    console.log({"sessionInfo:": this.sessionInfo, "windowName:": this.windowName, "webRemoteUrl": this.webRemoteUrl})
  }

  getWindowName() {
    return this.windowName
  }

  endSession() {
    const requestData = {
      action_type: 'end_session', // or 'end_session'
      session_info: this.sessionInfo,   // replace with your actual session ID
    };
    this.axios.post(this.webRemoteUrl, requestData)
      .then(response => {
        if(response.data){
          if(response.data.ok){
            console.log("nbcp web remote end session ok.");
          }else{
            console.error('nbcp web remote end session error, with response data: ', response.data);
          }
        }else {
          console.error('nbcp web remote end session error, with no response data');
        }
        console.log(response.data);
      })
      .catch(error => {
        console.error('nbcp web remote end session error:', error);
      });
  }

  rpc(action, params) {
    const requestData = {
      action_type: 'rpc', // or 'end_session'
      session_info: this.sessionInfo,   // replace with your actual session ID

      rpc_payload: {
        rpc_action: action,  // replace with your actual RPC action name
        param: params,
      },
    };
    this.axios.post(this.webRemoteUrl, requestData)
      .then(response => {
        if(response.data){
          if(response.data.ok){
            console.log("nbcp web remote rpc ok.");
          }else{
            console.error('nbcp web remote rpc error, with response data: ', response.data);
          }
        }else {
          console.error('nbcp web remote rpc error, with no response data');
        }
        console.log(response.data);
      })
      .catch(error => {
        console.error('nbcp web remote rpc error:', error);
      });
  }

  getCurrentRootUrl() {
    var protocol = window.location.protocol;
    var host = window.location.host;
    return protocol + '//' + host + '/';
  }

}

function replaceWebSocketWithHttp(originalUrl) {
  const regex = /^(ws|wss):\/\//;
  return originalUrl.replace(regex, (match, protocol) => {
    if (protocol === 'ws') {
      return 'http://';
    } else if (protocol === 'wss') {
      return 'https://';
    }
    return match;
  });
}

function getCurrentRootUrl() {
  var protocol = window.location.protocol;
  var host = window.location.host;
  return protocol + '//' + host + '/';
}

export { NBCPWebComm, replaceWebSocketWithHttp, getCurrentRootUrl }

