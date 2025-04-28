BIN=server
SRC=main.go

build:
	go build -o dist/$(BIN) $(SRC)

run: build
	./dist/$(BIN)

clean:
	rm -f $(BIN)
