package system

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// docker output format see https://github.com/docker/cli/blob/v20.10.0-beta1/cli/command/formatter/container.go

// dockerfile中ADD、COPY等命令中的src是相对与此执行程序所在目录的相对路径
func BuildImageFromDockerfile(ctx context.Context, dockerfile string, name string, version string) error {
	_, err := CommandOnceWithContext(ctx, "docker", "build", "--file", dockerfile, "-t", fmt.Sprintf("%s:%s", name, version), ".")
	if err != nil {
		return err
	}

	return nil
}

// docker version --format '{{.Server.Version}}'
func GetDockerVersion(ctx context.Context) (string, error) {
	data, err := CommandOnceWithContext(ctx, "docker", "version", "--format={{.Server.Version}}")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type ContainerInfo struct {
	ID         string `json:"ID"`
	Image      string `json:"Image"`
	Command    string `json:"Command"`
	CreatedAt  string `json:"CreatedAt"`
	RunningFor string `json:"RunningFor"`
	Ports      string `json:"Ports"`
	State      string `json:"State"`
	Status     string `json:"Status"`
	Size       string `json:"Size"`
	Names      string `json:"Names"`
	Labels     string `json:"Labels"`
	Mounts     string `json:"Mounts"`
	Networks   string `json:"Networks"`
}

func ContainerAllList(ctx context.Context) ([]ContainerInfo, error) {
	data, err := CommandOnceWithContext(ctx, "docker", "ps", "-a", "--no-trunc", `--format='{{json .}}'`)
	if err != nil {
		return nil, err
	}

	data = bytes.Trim(data, "\n\r ")
	list := bytes.Split(data, []byte("\n"))

	out := make([]ContainerInfo, 0, 32)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, v := range list {
		v = bytes.Trim(v, "' \n\r")

		var temp ContainerInfo
		err = json.Unmarshal(v, &temp)
		if err != nil {
			fmt.Printf("err = %s\n", err)
			continue
		}

		out = append(out, temp)
	}

	return out, nil
}

type ContainerState struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	CpuPerc      float64 `json:"cpuPerc"`
	MemPerc      float64 `json:"memperc"`
	MemUsage     float64 `json:"memUsage"`
	MemUsageUnit string  `json:"memUsageUnit"`
	MemLimit     float64 `json:"memLimit"`
	MemLimitUnit string  `json:"memLimitUnit"`
}

type containerState struct {
	ID        string `json:"ID"`
	Container string `json:"Container"`
	Name      string `json:"Name"`
	CPUPerc   string `json:"CPUPerc"`
	MemPerc   string `json:"MemPerc"`
	MemUsage  string `json:"MemUsage"`
	NetIO     string `json:"NetIO"`
	BlockIO   string `json:"BlockIO"`
	PIDs      string `json:"PIDs"`
}

func parseContainerState(state containerState) (out ContainerState, err error) {
	out.Id = state.ID
	out.Name = state.Name

	out.CpuPerc, err = strconv.ParseFloat(strings.Trim(state.CPUPerc, `%`), 64)
	if err != nil {
		fmt.Printf("CpuPerc = %s ParseFloat err = %s\n", strings.Trim(state.CPUPerc, `%`), err)
		return out, err
	}

	out.MemPerc, err = strconv.ParseFloat(strings.Trim(state.MemPerc, `%`), 64)
	if err != nil {
		fmt.Printf("MemPerc = %s ParseFloat err = %s\n", strings.Trim(state.MemPerc, `%`), err)
		return out, err
	}

	list := strings.Split(state.MemUsage, "/")
	if len(list) != 2 {
		fmt.Printf("memory usage invalid\n")
		return out, errors.New("memory usage invalid")
	}

	memUseStr := strings.Trim(list[0], " ")
	cleanMemUseStr := strings.TrimSuffix(memUseStr, "KiB")
	cleanMemUseStr = strings.TrimSuffix(cleanMemUseStr, "MiB")
	cleanMemUseStr = strings.TrimSuffix(cleanMemUseStr, "GiB")
	cleanMemUseStr = strings.TrimSuffix(cleanMemUseStr, "B")

	memLimitStr := strings.Trim(list[1], " ")
	cleanMemLimitStr := strings.TrimSuffix(memLimitStr, "KiB")
	cleanMemLimitStr = strings.TrimSuffix(cleanMemLimitStr, "MiB")
	cleanMemLimitStr = strings.TrimSuffix(cleanMemLimitStr, "GiB")
	cleanMemLimitStr = strings.TrimSuffix(cleanMemLimitStr, "B")

	out.MemUsage, err = strconv.ParseFloat(cleanMemUseStr, 64)
	if err != nil {
		fmt.Printf("MemUsage = %s ParseFloat err = %s\n", cleanMemUseStr, err)
		return out, err
	}

	if strings.Contains(memUseStr, "GiB") {
		out.MemUsageUnit = "GiB"
	} else if strings.Contains(memUseStr, "MiB") {
		out.MemUsageUnit = "MiB"
	} else if strings.Contains(memUseStr, "KiB") {
		out.MemUsageUnit = "KiB"
	} else if strings.Contains(memUseStr, "B") {
		out.MemUsageUnit = "B"
	} else {
		fmt.Printf("mem usage unit err = %s\n", memUseStr)
		return out, err
	}

	out.MemLimit, err = strconv.ParseFloat(cleanMemLimitStr, 64)
	if err != nil {
		fmt.Printf("memLimit = %s ParseFloat err = %s\n", cleanMemLimitStr, err)
		return out, err
	}

	if strings.Contains(memLimitStr, "GiB") {
		out.MemLimitUnit = "GiB"
	} else if strings.Contains(memLimitStr, "MiB") {
		out.MemLimitUnit = "MiB"
	} else if strings.Contains(memLimitStr, "KiB") {
		out.MemLimitUnit = "KiB"
	} else if strings.Contains(memLimitStr, "B") {
		out.MemLimitUnit = "B"
	} else {
		fmt.Printf("mem limit unit err = %s\n", memUseStr)
		return out, err
	}

	return out, nil
}

func ContainerUsage(ctx context.Context) (<-chan ContainerState, error) {
	var outReader *bufio.Reader

	err := CommandAlwaysWithContext(ctx, &outReader, "docker", "stats", "--all", "--no-trunc", `--format='{{json .}}'`)
	if err != nil {
		return nil, err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	outCh := make(chan ContainerState, 1)

	go func() {
		for {
			line, err := outReader.ReadBytes('\n')
			if err != nil {
				if io.EOF != err {
					fmt.Printf("readBytes err = %s\n", err)
				}

				return
			}

			begin := strings.IndexAny(string(line), "{")
			if begin >= 0 {
				line = line[begin:]
			}

			line = bytes.Trim(line, "' \n\r")

			var usage containerState
			err = json.Unmarshal(line, &usage)
			if err != nil {
				fmt.Printf("json Unmarshal err = %s\n", err)
				continue
			}

			state, err := parseContainerState(usage)
			if err != nil {
				continue
			}

			outCh <- state
		}
	}()

	return outCh, nil
}

func StopContainer(ctx context.Context, id string) error {
	data, err := CommandOnceWithContext(ctx, "docker", "stop", id)
	if err != nil {
		return err
	}

	if !strings.Contains(id, string(bytes.Trim(data, "\n\r "))) {
		fmt.Printf("data = %s\n", data)
		return errors.New("stop failed")
	}

	return nil
}

func StartContainer(ctx context.Context, id string) error {
	data, err := CommandOnceWithContext(ctx, "docker", "start", id)
	if err != nil {
		return err
	}

	if !strings.Contains(id, string(bytes.Trim(data, "\n\r "))) {
		fmt.Printf("data = %s\n", data)
		return errors.New("stop failed")
	}

	return nil
}

func RestartContainer(ctx context.Context, id string) error {
	data, err := CommandOnceWithContext(ctx, "docker", "restart", id)
	if err != nil {
		return err
	}

	if !strings.Contains(id, string(bytes.Trim(data, "\n\r "))) {
		fmt.Printf("data = %s\n", data)
		return errors.New("stop failed")
	}

	return nil
}

func RemoveContainer(ctx context.Context, id string) error {
	data, err := CommandOnceWithContext(ctx, "docker", "rm", id)
	if err != nil {
		return err
	}

	if !strings.Contains(id, string(bytes.Trim(data, "\n\r "))) {
		fmt.Printf("data = %s\n", data)
		return errors.New("stop failed")
	}

	return nil
}

func CreateContainer(ctx context.Context, baseImage string, name string, arg ...string) (string, error) {
	argList := make([]string, 0, 32)
	argList = append(argList, "run")
	argList = append(argList, "--net=host")
	argList = append(argList, "--restart")
	argList = append(argList, "always")
	argList = append(argList, "-d")

	for _, v := range arg {
		str := strings.Trim(v, " ")
		if strings.Contains(str, " ") {
			return "", errors.New("arg can't include empty space")
		}

		argList = append(argList, str)
	}

	argList = append(argList, "--name")
	argList = append(argList, strings.Trim(name, " "))
	argList = append(argList, strings.Trim(baseImage, " "))

	data, err := CommandOnceWithContext(ctx, "docker", argList...)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
