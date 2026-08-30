#VARIABILI
BINARY_NAME=gialosai #eseguibile finale
SRC_DIR=./src
BIN_DIR=./bin

#TARGET

#default
all: compile

compile:
	@echo "Compilazione del progetto in corso..."
	@mkdir -p $(BIN_DIR)
	cd server && go build -o ../$(BIN_DIR)/$(BINARY_NAME) ./cmd/gialosai
	@echo "Compilazione completata! Eseguibile generato in: $(BIN_DIR)/$(BINARY_NAME)"
	
test:
	@echo "Esecuzione dei Test di Unità..."
	cd client && go test -v ./...
	cd server && go test -v ./...

clean:
	@echo "Pulizia in corso..."
	go clean
	rm -rf $(BIN_DIR)
	@echo "Pulizia completata!"

.PHONY: all compile test clean
