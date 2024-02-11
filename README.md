一款检测 ThinkPHP 日志泄露的扫描工具，具体可查看：

http://furina.org.cn/2024/01/25/thinkphp-log-leakage

编译

```bash
go build thinkphp-log-leakage.go
```

扫描

```bash
./thinkphp-log-leakage -u http://exam.com
```

