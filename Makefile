
INSTALL_DIR:=~/go/bin/
BIN_NAME:=polycli
BUILD_DIR:=./out

GIT_SHA:=$(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG:=$(shell git describe --tags)
CUR_DATE:=$(shell date +%s)

# strip debug and supress warnings
LD_FLAGS=-s -w
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Version=$(GIT_TAG)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Commit=$(GIT_SHA)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.Date=$(CUR_DATE)\"
LD_FLAGS += -X \"github.com/maticnetwork/polygon-cli/version.BuiltBy=makefile\"
STATIC_LD_FLAGS=$(LD_FLAGS) -extldflags=-static

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build: $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BIN_NAME) main.go

cross: $(BUILD_DIR)
	env CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64	go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-arm64-$(BIN_NAME) main.go
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64				go build -ldflags "$(STATIC_LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/linux-amd64-$(BIN_NAME) main.go

        # # mac builds - this will be functional but will still have secp issues
	env GOOS=darwin GOARCH=arm64						go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-arm64-$(BIN_NAME) main.go
	env GOOS=darwin GOARCH=amd64						go build -ldflags "$(LD_FLAGS)" -tags netgo -o $(BUILD_DIR)/darwin-amd64-$(BIN_NAME) main.go

install: build
	$(RM) $(INSTALL_DIR)/$(BIN_NAME)
	cp $(BUILD_DIR)/$(BIN_NAME) $(INSTALL_DIR)

clean:
	$(RM) -r $(BUILD_DIR)

test:
	go test github.com/maticnetwork/polygon-cli/...
