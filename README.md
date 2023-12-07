# NagaeBrowserControllingProtocol

# 协议文档

见`protocol_docs_v1/nbcp-docs-v1.pdf`

# 测试

1. 编译`pkgtest/test-server-v1`
2. 编译前端`pkgtest/nbcp-test-frontend`
3. 创建目录`wwwroot`，将编译好的前端拷入
4. 创建配置文件`config.toml`

```toml
[server]
bind = ":9299"
wwwroot = "wwwroot"

[nbcp]
home = "http://localhost:9299"
```

5. 启动test-server-v1

```powershell
.\test-server-v1.exe
```

6. 下载编译NagaeSimpleWebBrowser

https://git.swzry.com/ProjectNagae/NagaeSimpleWebBrowser

github镜像：

https://github.com/swzry/NagaeBrowserControllingProtocol

7. 启动NagaeSimpleWebBrowser

```powershell
.\NagaeSimpleWebBrowser.exe with-nbcp --url ws://localhost:9299/nbcp
```