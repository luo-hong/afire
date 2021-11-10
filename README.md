# afire

##项目编译

1. 生成二进制文件
2. bin/afire_man 就是可运行的bin文件

##启动服务

```bash
./bin/afire_man -c ./configs/manager.toml
```

命令
```bash
supervisorctl start manager #启动
supervisorctl restart manager #重启
supervisorctl stop manager #停止
```

##项目位置 阿里云服务器

1. 项目代码位置 `~/afire`
2. 服务日志和配置文件位置 `~/manager`

## version.go的处理

`version.go` 用于自动生成编译信息