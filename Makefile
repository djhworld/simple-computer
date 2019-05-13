all:
	@@go build -o bin/simulator github.com/djhworld/simple-computer/cmd/simulator
	@@go build -o bin/assembler github.com/djhworld/simple-computer/cmd/assembler
	@@go build -o bin/generator github.com/djhworld/simple-computer/cmd/generator


test:
	go test ./...
