## Queue Worker

### Packages Used

- github.com/hibiken/asynq
- github.com/redis/go-redis/v9

### How to run
Run redis server first

### Using hibiken/asynq
```
go test ./pkg/xasync -run TestQueue_Hibiken -v
go test ./pkg/xasync -run TestWorker_Hibiken -v
```

### Using go-redis/v9
```
go test ./pkg/xasync -run TestQueue_Async -v
go test ./pkg/xasync -run TestWorker_Async -v
```
