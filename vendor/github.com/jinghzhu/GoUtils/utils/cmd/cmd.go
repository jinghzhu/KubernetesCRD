package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jinghzhu/GoUtils/array"
)

// ValidMounts return true if all of the mountpoins are legal
func ValidMounts(ms []Mount) (bool, error) {
	errMsg := "Mount data is not correct. "
	if len(ms) == 0 {
		return true, nil
	}
	mountData := make(map[string][]string)
	for _, m := range ms {
		if m.Server == "" || !strings.Contains(m.Share, "/") {
			fmt.Println(errMsg)
			return false, errors.New(errMsg)
		}
		if data, ok := mountData[m.Server]; ok {
			if !array.Include(data, m.Share) {
				return false, nil
			}
		} else {
			mountPoints, err := GetMountPoints(m.Server)
			if err != nil {
				errMsg += err.Error()
				fmt.Println(errMsg)
				return false, err
			}
			mountData[m.Server] = mountPoints
			if !array.Include(mountPoints, m.Share) {
				return false, nil
			}
		}
	}
	return true, nil
}

// GetMountPoints returns all moutpoins in a string array
func GetMountPoints(server string) ([]string, error) {
	cmd := exec.Command(CmdShowmount, CmdShowmountOptE, server)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Start()
	done := make(chan error)
	timeout := make(chan bool, 1)

	go func() {
		done <- cmd.Wait()
		close(done)
	}()
	go func() {
		time.Sleep(ShowmountTimeout)
		timeout <- true
		close(timeout)
	}()
	select {
	case <-timeout:
		errMsg := "Command showmout timeout for server " + server
		fmt.Println(errMsg)
		err := cmd.Process.Kill()
		if err != nil {
			errMsg1 := "Error to kill showmount process. " + err.Error()
			fmt.Println(errMsg1)
			return nil, errors.New(errMsg + " " + errMsg1)
		}
		return nil, errors.New(errMsg)
	case err := <-done:
		if err != nil {
			errMsg := fmt.Sprintf("Error in getting mount info of server(%s): %s %s\n", server, err.Error(), buf.String())
			fmt.Println(errMsg)
			return nil, errors.New(errMsg)
		}
	}

	s := strings.TrimSpace(buf.String())
	// The fist line of showmount -e <server> is Exports list on <server>
	firstLine := strings.Index(s, "\n")
	sArr := strings.Split(s[firstLine+1:], "\n")
	for i := 0; i < len(sArr); i++ {
		index := strings.Index(sArr[i], " ")
		temp := sArr[i]
		sArr[i] = temp[:index]
	}

	fmt.Println("Mount points of server " + server)
	fmt.Println(sArr)

	return sArr, nil
}
