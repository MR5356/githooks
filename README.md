# 编译
```shell
make build
```

# 发布
```shell
make all
```

# upx加壳压缩
```shell
make build
upx -9 -o githooks-upx githooks
```
