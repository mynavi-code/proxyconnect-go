# proxyconnect-go

Tool for tunneling SSH through HTTP proxies by Go

# Usage

```
$ http_proxy=http://username:password@proxyhost:proxyport ssh -oProxyCommand='proxyconnect.py %h %p' targethost
```

# Additional Infomation

This is successor to proxyconnect

https://github.com/hirosuzuki/proxyconnect
