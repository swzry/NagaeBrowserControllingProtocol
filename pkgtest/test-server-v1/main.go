package main

import (
	nbcp "git.swzry.com/ProjectNagae/NagaeBrowserControllingProtocol/server_v1"
	"git.swzry.com/zry/zry-go-program-framework/easy_toml_config"
	"git.swzry.com/zry/zry-go-program-framework/svcfw"
	"git.swzry.com/zry/zry-go-program-framework/websubsvc"
	"net/http"
	"os"
)

var NBCPServer *nbcp.NBCPServer
var SessionHdl *nbcp.SimpleSessionHandler
var WebLogic *WebLogicClass
var SessionLogic *SessionLogicClass
var WWWRootFS websubsvc.IFilesystem

var App *svcfw.AppFramework
var SubSvcNBCP *nbcp.ZGPFSubSvc
var SubSvcWeb *websubsvc.WebSubService

func main() {
	App = svcfw.NewAppFramework(true, "main")
	App.InitConsoleLogBackend(os.Stdout, "")

	cwd, err := os.Getwd()
	if err == nil {
		App.Info("Work Dir:", cwd)
	} else {
		App.Warn("failed get work dir: ", err)
	}

	App.Info("preparing phase1...")
	App.MustPrepare("config", func() error {
		return easy_toml_config.LoadConfigFromFile("config.toml", &Config)
	})

	SessionLogic = NewSessionLogic()
	SessionHdl = nbcp.NewSimpleSessionHandler(SessionLogic)
	NBCPServer = nbcp.NewNBCPServer(SessionHdl)
	WebLogic = NewWebLogic(Config.Server.BindAddr)

	/*
			App.MustPrepare("wwwroot", func() error {
			wrabs, xerr := filepath.Abs(Config.Server.WWWRoot)
			if xerr != nil {
				return fmt.Errorf("failed get abs path of wwwroot: %v", xerr)
			}
			App.Info("wwwroot: ", wrabs)
			wrfs, ok := os.DirFS(wrabs).(websubsvc.IFilesystem)
			if !ok {
				return fmt.Errorf("type assertion failed for os.DirFS to websubsvc.IFilesystem")
			}
			WWWRootFS = wrfs
			return nil
		})
	*/
	WWWRootFS = websubsvc.NewOsFsWithPrefix(Config.Server.WWWRoot)

	App.Info("init sub services...")

	SubSvcNBCP = nbcp.NewZGPFSubSvc(NBCPServer)
	SubSvcWeb = websubsvc.NewWebSubService(WebLogic)
	App.AddSubSvc("nbcp", SubSvcNBCP)
	App.AddSubSvc("web", SubSvcWeb)
	App.AddSubSvc("sig-quit", svcfw.NewWatchSignalExitSubServiceWithDefault())

	App.Info("preparing phase2...")
	err = App.Prepare()
	if err != nil {
		App.Panic("error in preparing sub services: ", err)
	}

	App.Info("start running sub services...")
	err = App.Run()

	App.Info("all sub services end.")
	if err != nil {
		App.Panic("error in running sub services: ", err)
	}
}

type WebLogicClass struct {
	bindAddr string
}

func NewWebLogic(bindAddr string) *WebLogicClass {
	l := &WebLogicClass{
		bindAddr: bindAddr,
	}
	return l
}

func (w WebLogicClass) GetHttpServer(ctx *websubsvc.WebSubServiceContext) *http.Server {
	return ctx.MakeHttpServer(w.bindAddr, nil)
}

func (w WebLogicClass) Prepare(ctx *websubsvc.WebSubServiceContext) error {
	ctx.DefaultMiddleware("access", "err500")
	ctx.Info("prepare URL router...")
	ctx.GetRootRouter().GET("/nbcp", NBCPServer.HandleNBCP)
	ctx.GetRootRouter().POST("/nbcp_web_remote", NBCPServer.HandleNBCPWebCommRemote)
	ctx.EnableFrontend(WWWRootFS, true)
	return nil
}
