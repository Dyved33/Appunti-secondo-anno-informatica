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
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "Compilazione completata! Eseguibile generato in: $(BIN_DIR)/$(BINARY_NAME)"

run: compile
	@echo "Avvio dell'applicazione..."
	$(BIN_DIR)/$(BINARY_NAME)

test:
	@echo "Esecuzione dei Test di Unità..."
	go test -v ./...

clean:
	@echo "Pulizia in corso..."
	go clean
	rm -rf $(BIN_DIR)
	@echo "Pulizia completata!"

.PHONY: all compile run test clean
