package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	injectVariableNameFlag = flag.String("injectvar", "", `The inject variable name, will be used as -ldflags "-X INJECT_VAR_NAME=" (could for example be main.GitHash if you have a GitHash variable in your main package)`)
	outFlag                = flag.String("out", "", "This flag is directly passed to `go build -o=`")
)

func runCommandGetCombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed to run, error: %s. Output: %s", err.Error(), string(out))
	}
	return out, nil
}

func main() {
	flag.Parse()
	if len(*outFlag) == 0 ||
		len(*injectVariableNameFlag) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	tmpGitHashBytes, err := runCommandGetCombinedOutput(exec.Command("git", "rev-parse", "HEAD"))
	if err != nil {
		log.Fatal(err)
	}
	gitHash := strings.TrimSpace(string(tmpGitHashBytes))

	envExpandedOutFlag := os.ExpandEnv(*outFlag)
	ldFlags := fmt.Sprintf("-X %s=%s", *injectVariableNameFlag, gitHash)

	fmt.Println(fmt.Sprintf("Got git hash: %s", gitHash))
	fmt.Println(fmt.Sprintf("Using go build out flag: %s", envExpandedOutFlag))
	fmt.Println(fmt.Sprintf("Using ldflags: %s", ldFlags))

	goBuildCmd := exec.Command("go", "build", "-ldflags", ldFlags, "-o", envExpandedOutFlag)

	// Proxy the Stdio
	goBuildCmd.Stdin = os.Stdin
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr

	err = goBuildCmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = goBuildCmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
