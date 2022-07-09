
INSTALL_DIR:=~/go/bin/
BIN_NAME:=polycli

build:
	go build -o $(BIN_NAME) main.go

install: build
	cp $(BIN_NAME) $(INSTALL_DIR)


