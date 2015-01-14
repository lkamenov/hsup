package main

import (
	"fmt"
	"os"
)

type DynoDriver interface {
	Build(*Release) error
	Start(*Executor) error
	Stop(*Executor) error
	Wait(*Executor) *ExitStatus
}

type ExitStatus struct {
	code int
	err  error
}

type Release struct {
	appName string
	config  map[string]string
	slugURL string
	version int

	// docker dyno driver properties
	imageName string
}

func (r *Release) Name() string {
	return fmt.Sprintf("%v-%v", r.appName, r.version)
}

func (r *Release) ConfigSlice() []string {
	var c []string
	for k, v := range r.config {
		c = append(c, k+"="+v)
	}
	return c
}

func FindDynoDriver(name string) (DynoDriver, error) {
	switch name {
	case "simple":
		return &SimpleDynoDriver{}, nil
	case "docker":
		return &DockerDynoDriver{}, nil
	case "abspath":
		return &AbsPathDynoDriver{}, nil
	case "libcontainer":
		newRoot := os.Getenv("HSUP_NEWROOT")
		if newRoot == "" {
			return nil, fmt.Errorf("HSUP_NEWROOT empty")
		}

		hostname := os.Getenv("HSUP_HOSTNAME")
		if hostname == "" {
			return nil, fmt.Errorf("HSUP_HOSTNAME empty")
		}

		user := os.Getenv("HSUP_USER")
		if user == "" {
			return nil, fmt.Errorf("HSUP_USER empty")
		}

		return &LibContainerDynoDriver{
			NewRoot:  newRoot,
			User:     user,
			Hostname: hostname,
		}, nil
	default:
		return nil, fmt.Errorf("could not locate driver. "+
			"specified by the user: %v", name)
	}
}
