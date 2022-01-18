package system

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func BashCommand(cmd string) ([]byte, error) {
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	return result, err
}

func CommandOnceWithContext(ctx context.Context, exeName string, args ...string) ([]byte, error) {
	if strings.Contains(exeName, " ") {
		return nil, errors.New("exe name can't include space")
	}

	newArgs := make([]string, 0, 8)
	for _, v := range args {
		if strings.Contains(v, " ") {
			list := strings.Split(v, " ")
			newArgs = append(newArgs, list...)
		} else {
			newArgs = append(newArgs, v)
		}
	}

	cmd := exec.CommandContext(ctx, exeName, newArgs...)

	var buf bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return buf.Bytes(), fmt.Errorf("%s, detail info = %s", err, stderr.String())
	}

	if err := cmd.Wait(); err != nil {
		return buf.Bytes(), fmt.Errorf("err = %s, info = %s", err, stderr.String())
	}

	return buf.Bytes(), nil
}

func CommandAlwaysWithContext(ctx context.Context, outReader **bufio.Reader, exeName string, args ...string) error {
	if strings.Contains(exeName, " ") {
		return errors.New("exe name can't include space")
	}

	newArgs := make([]string, 0, 8)
	for _, v := range args {
		if strings.Contains(v, " ") {
			list := strings.Split(v, " ")
			newArgs = append(newArgs, list...)
		} else {
			newArgs = append(newArgs, v)
		}
	}

	cmd := exec.CommandContext(ctx, exeName, newArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	*outReader = bufio.NewReader(stdout)

	go cmd.Wait()

	return nil
}
