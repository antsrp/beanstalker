.PHONY: run
run: --setargs
	@go run cmd/main.go $(cmd)

.PHONY: build run-builded build-run
build-run: build run-builded
run-builded: --setargs
	@cd bin && ./beanstalker $(cmd)
build:
	@go build -o bin/beanstalker cmd/main.go

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: --setargs
--setargs:
ifneq ("$(host)","")
	$(eval cmd = -H $(host))	
endif
ifneq ("$(port)","")
	$(eval cmd += -P $(port))	
endif
