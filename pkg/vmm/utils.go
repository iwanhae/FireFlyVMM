package vmm

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
)

func (v *VirtualMachineManager) prepareDirectories(ctx context.Context) error {
	if err := os.MkdirAll(v.SocketDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.TemplateDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.DataDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(v.LogDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (v *VirtualMachineManager) generateVMID() (string, error) {
	entries, err := os.ReadDir(v.DataDir)
	if err != nil {
		return "", err
	}
	max := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		v, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if max < v {
			max = v
		}
	}
	return fmt.Sprintf("%05d", max+1), nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
