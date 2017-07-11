package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	config := api.DefaultConfig()
	config.Address = "localhost:8500"
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	ag := client.Agent()
	if err = RegisterService(ag, "sendgrid", "api.sendgrid.com", "sendgrid", "https", 443); err != nil {
		fmt.Printf("Encountered error registering a service with consul -> %s\n", err)
	}

	if err = RegisterService(ag, "twilio", "api.twilio.com", "twilio", "https", 443); err != nil {
		fmt.Printf("Encountered error registering a service with consul -> %s\n", err)
	}

	select {}
}

func RegisterService(ag *api.Agent, service, hostname, node, protocol string, port int) error {
	rand.Seed(time.Now().UnixNano())

	sid := rand.Intn(65534)
	serviceID := service + "-" + strconv.Itoa(sid)

	consulService := api.AgentServiceRegistration{
		ID:   serviceID,
		Name: service,
		Tags: []string{time.Now().Format("Jan 02 15:04:05.000 MST")},
		Port: port,
		Check: &api.AgentServiceCheck{
			Script:   "curl --connect-timeout=5 " + protocol + "://" + hostname + ":" + strconv.Itoa(port),
			Interval: "10s",
			Timeout:  "8s",
			TTL:      "",
			HTTP:     protocol + "://" + hostname + ":" + strconv.Itoa(port),
			Status:   "passing",
		},
		Checks: api.AgentServiceChecks{},
	}
	err := ag.ServiceRegister(&consulService)
	if err != nil {
		return err
	}
	return nil
}
