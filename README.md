# kube-exec

## 批量在 pod 里执行命令

- 每个 deployment 只取一个正常运行状态的 pod 执行
- 默认执行容器为`app`

## 构建并执行命令

```
go build && ./kube-exec
```
