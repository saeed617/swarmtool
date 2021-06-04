package swarmtool

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"log"
)

type (
	// Connection is a connection to systemd's dbus endpoint.
	Connection interface {
		StartUnit(ctx context.Context, name string) error
		StopUnit(ctx context.Context, name string) error
		Status(ctx context.Context, name string) (string, error)
		Close()
	}
	// DbusConnection implements Connection.
	DbusConnection struct {
		Conn *dbus.Conn
	}
	// Dockerd manages docker daemon.
	Dockerd struct {
		DbusConn Connection
	}
)

const dockerServiceName = "docker.service"

// StartUnit starts specified unit.
func (d *DbusConnection) StartUnit(ctx context.Context, name string) error {
	ch := make(chan string)
	jobId, err := d.Conn.StartUnitContext(ctx, name, "replace", ch)
	if err != nil {
		return err
	}
	log.Printf("starting %s with job id %d ...", name, jobId)
	jobResult := <-ch
	if jobResult != "done" {
		return errors.New(fmt.Sprintf("failed to start %s with status %s", name, jobResult))
	}
	log.Printf("%s started", name)
	return nil
}

// StopUnit stops specified unit.
func (d *DbusConnection) StopUnit(ctx context.Context, name string) error {
	ch := make(chan string)
	jobId, err := d.Conn.StopUnitContext(ctx, name, "replace", ch)
	if err != nil {
		return err
	}
	log.Printf("stoping %s with job id %d ...", name, jobId)
	jobResult := <-ch
	if jobResult != "done" {
		return errors.New(fmt.Sprintf("failed to stop %s with status %s", name, jobResult))
	}
	log.Printf("%s stopped.", name)
	return nil
}

// Status returns status of specified unit.
func (d *DbusConnection) Status(ctx context.Context, name string) (string, error) {
	services, err := d.Conn.ListUnitsByNamesContext(ctx, []string{name})
	if err != nil {
		return "", err
	}
	if len(services) < 1 {
		return "", errors.New(fmt.Sprintf("service %s not found", name))
	}
	service := services[0]
	return service.ActiveState, nil
}

func (d *DbusConnection) Close() {
	d.Conn.Close()
}

// Start starts docker service.
func (d *Dockerd) Start() error {
	ctx := context.Background()
	return d.DbusConn.StartUnit(ctx, dockerServiceName)
}

// Stop stops docker service.
func (d *Dockerd) Stop() error {
	ctx := context.Background()
	return d.DbusConn.StopUnit(ctx, dockerServiceName)
}

// Status returns status of docker service.
func (d *Dockerd) Status() (string, error) {
	ctx := context.Background()
	return d.DbusConn.Status(ctx, dockerServiceName)
}

// IsActive checks if docker service is running.
func (d *Dockerd) IsActive() (bool, error) {
	status, err := d.Status()
	return status == "active", err
}
