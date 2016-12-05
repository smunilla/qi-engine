package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	config := Config{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- config dump:\n%s\n\n", string(d))

}

/*
Config Top-level config file object
*/
type Config struct {
	Deployment Deployment `yaml:"deployment"`
	Vars       map[string]string
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		Deployment Deployment `yaml:"deployment"`
	}

	if err := unmarshal(&params); err != nil {
		return err
	}
	var variables map[string]string
	if err := unmarshal(&variables); err != nil {
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
	}

	c.Deployment = params.Deployment
	c.Vars = variables

	return nil
}

/*
Deployment List of any Hosts and Roles
*/
type Deployment struct {
	Hosts []Host
	Roles Role
	Vars  map[string]string
}

func (d *Deployment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		Hosts []Host
		Roles Role
	}

	if err := unmarshal(&params); err != nil {
		return err
	}
	var variables map[string]string
	if err := unmarshal(&variables); err != nil {
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
	}

	d.Hosts = params.Hosts
	d.Roles = params.Roles
	d.Vars = variables

	return nil
}

/*
Host All information about a host
*/
type Host struct {
	ConnectTo      string   `yaml:"connect_to"`
	Hostname       string   `yaml:"hostname,omitempty"`
	PublicHostname string   `yaml:"public_hostname,omitempty"`
	IP             string   `yaml:"ip,omitempty"`
	PublicIP       string   `yaml:"public_ip,omitempty"`
	NodeLabels     string   `yaml:"node_labels,omitempty"`
	Roles          []string `yaml:"roles"`
	Vars           map[string]string
}

func (h *Host) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		ConnectTo      string   `yaml:"connect_to"`
		Hostname       string   `yaml:"hostname,omitempty"`
		PublicHostname string   `yaml:"public_hostname,omitempty"`
		IP             string   `yaml:"ip,omitempty"`
		PublicIP       string   `yaml:"public_ip,omitempty"`
		NodeLabels     string   `yaml:"node_labels,omitempty"`
		Roles          []string `yaml:"roles"`
	}

	if err := unmarshal(&params); err != nil {
		return err
	}
	var variables map[string]string
	if err := unmarshal(&variables); err != nil {
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
	}

	keysToDel := []string{"connect_to", "hostname", "public_hostname", "ip", "public_ip"}
	for _, v := range keysToDel {
		delete(variables, v)
	}

	h.ConnectTo = params.ConnectTo
	h.Hostname = params.Hostname
	h.PublicHostname = params.PublicHostname
	h.IP = params.IP
	h.PublicIP = params.PublicIP
	h.NodeLabels = params.NodeLabels
	h.Roles = params.Roles
	h.Vars = variables

	fmt.Printf("%s\n", variables)

	return nil
}

/*
Role A map of variables to be applied to a group of Hosts
*/
type Role map[string](map[string]string)

var data = `
ansible_callback_facts_yaml: /home/smunilla/.config/openshift/.ansible/callback_facts.yaml
ansible_config: /usr/share/atomic-openshift-utils/ansible.cfg
ansible_inventory_path: /home/smunilla/.config/openshift/hosts
ansible_log_path: /tmp/ansible.log
deployment:
  ansible_ssh_user: openshift
  hosts:
  - connect_to: 192.168.55.233
    hostname: armory-master-91c03.example.com
    ip: 192.168.55.233
    node_labels: '{''region'': ''infra''}'
    public_hostname: armory-master-91c03.example.com
    public_ip: 192.168.55.233
    foo: "bar"
    bar: false
    roles:
    - master
    - etcd
    - node
    - storage
  - connect_to: 192.168.55.8
    hostname: armory-node-compute-a4330.example.com
    ip: 192.168.55.8
    public_hostname: armory-node-compute-a4330.example.com
    public_ip: 192.168.55.8
    roles:
    - node
  master_routingconfig_subdomain: ''
  proxy_exclude_hosts: ''
  proxy_http: ''
  proxy_https: ''
  roles:
    etcd: {}
    master: {}
    node: {}
    storage: {}
variant: openshift-enterprise
variant_version: '3.3'
version: v2
`
