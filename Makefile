mockgen:
	mockgen -source=internal/server/interfaces.go -destination=./mocks/server_mock.go -package=mocks
	mockgen -source=internal/logging/logger.go -destination=./mocks/logging_mock.go -package=mocks
	mockgen -source=internal/global_errors/global_errors.go -destination=./mocks/global_errors_mock.go -package=mocks