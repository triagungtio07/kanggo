mock:
	go generate ./pkg/usecase/user
	go generate ./pkg/storage/user
	go generate ./pkg/usecase/product
	go generate ./pkg/storage/product
	go generate ./pkg/usecase/order
	go generate ./pkg/storage/order

test:
	go test ./pkg/usecase/product -v -cover -covermode=atomic
	go test ./pkg/usecase/user -v -cover -covermode=atomic
	go test ./pkg/usecase/order -v -cover -covermode=atomic
	go test ./pkg/handler/order -v -cover -covermode=atomic
	go test ./pkg/handler/user -v -cover -covermode=atomic
	go test ./pkg/handler/product -v -cover -covermode=atomic
