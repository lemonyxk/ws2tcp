# ws2tcp
```shell
$ ws2tcp '127.0.0.1:12075 > ws://127.0.0.1:12076'
```

```shell
$ ws2tcp 'ws://127.0.0.1:12076 > 127.0.0.1:27017'
```

```shell
$ mongosh -host 127.0.0.1 -port 12075
Current Mongosh Log ID:	668e299016e623749b95477d
Connecting to:		mongodb://127.0.0.1:12075/?directConnection=true&appName=mongosh+2.2.10
Using MongoDB:		7.0.12
Using Mongosh:		2.2.10

For mongosh info see: https://docs.mongodb.com/mongodb-shell/

center [direct: secondary] test>
```
