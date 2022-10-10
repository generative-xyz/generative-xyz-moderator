package protoc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const API_OUTPUT = "api/"

var (
	includes = []string{
		"-Iproto",
		"-Ithird_party/googleapis",
		"-Ithird_party/envoyproxy",
	}
	outs = []string{
		"--go_out=%s",
		"--go-grpc_out=%s",
		`--validate_out="lang=go:%s"`,
		`--openapiv2_out=%s`,
		`--grpc-gateway_out=%s`,
	}
	opts = []string{
		`--openapiv2_opt`,
		`logtostderr=true`,
		`--grpc-gateway_opt`,
		`logtostderr=true`,
	}
)

func buildArgs(tempDir string, input string, includes []string) []string {
	args := []string{"protoc"}
	args = append(args, includes...)
	for i := 0; i < len(outs); i++ {
		args = append(args, fmt.Sprintf(outs[i], tempDir))
	}
	args = append(args, opts...)
	args = append(args, input)
	return args
}

func Action(c *cli.Context) error {
	rootPath, _ := os.Getwd()
	output := "tmp"
	tempPath := path.Join(rootPath, output)
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err := os.Mkdir(tempPath, os.FileMode(0777))
		if err != nil {
			fmt.Println(err)
		}
	}

	input := "proto/*.proto"

	args := buildArgs(output, input, includes)
	command := strings.Join(args, " ")
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	err = filepath.Walk(output,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			command := fmt.Sprintf(`cp %s %s`, path, API_OUTPUT)
			cmd := exec.Command("/bin/sh", "-c", command)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return errors.New(string(out))
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	cmd = exec.Command("/bin/sh", "-c",
		// TODO shuold dynamic here
		fmt.Sprintf("cp %s%s swaggerUI", API_OUTPUT, "api.swagger.json"),
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	cmd = exec.Command("/bin/sh", "-c", "rm -rf tmp/")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
