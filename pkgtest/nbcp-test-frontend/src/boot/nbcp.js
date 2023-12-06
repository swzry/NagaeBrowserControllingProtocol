import { boot } from 'quasar/wrappers'
import { NBCPWebComm, replaceWebSocketWithHttp } from "nbcp-web-comm"

export default boot(async ({ app }) => {
  const sessInfo = await chrome.webview.hostObjects.nbcp.getSessionInfo();
  const winName = await chrome.webview.hostObjects.nbcp.getWindowName();
  const nbcpBase = await chrome.webview.hostObjects.nbcp.getNBCPServerUrl();
  const nbcpWebRemote = replaceWebSocketWithHttp(nbcpBase) + "_web_remote"

  const nbcpWebComm = new NBCPWebComm(sessInfo, winName, nbcpWebRemote);
  app.config.globalProperties.$nbcp = nbcpWebComm;
})
