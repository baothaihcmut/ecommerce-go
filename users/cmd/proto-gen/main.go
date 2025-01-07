package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	GO_PACKAGE = "github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	goOut      = "./internal/adapter/grpc/proto/"
	protoFiles = []string{"../libs/proto/shared.proto", "../libs/proto/user.proto"}
)

func main() {
	//create proto dir temp
	protoDir := "./proto-temp/"
	err := os.MkdirAll(protoDir, 0755)
	if err != nil {
		fmt.Println("Error creating proto temp dir")
		panic(err)
	}
	for _, file := range protoFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file", file)
			panic(err)
		}
		newContent := strings.ReplaceAll(string(content), "{{GO_PACKAGE}}", GO_PACKAGE)
		tempFile := protoDir + filepath.Base(file)
		err = os.WriteFile(tempFile, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("Error writing temp file")
			panic(err)
		}
		cmd := exec.Command(
			"protoc",
			fmt.Sprintf("--proto_path=%s", protoDir),
			tempFile,
			fmt.Sprintf("--go_out=%s", goOut),
			fmt.Sprintf("--go-grpc_out=%s", goOut),
			"--go_opt=paths=source_relative",
			"--go-grpc_opt=paths=source_relative",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error running protoc")
			fmt.Println(file)
			panic(err)
		}
	}
	//remove proto dir temp
	err = os.RemoveAll(protoDir)
	if err != nil {
		fmt.Println("Error removing proto temp dir")
		panic(err)
	}
}
