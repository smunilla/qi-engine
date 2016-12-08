package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("please provide a file!")
	}

	configFilePath := flag.Args()[0]
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(content)

	config := readConfig(data)
	//dumpConfig(config)

	inventoryTemplate.Execute(os.Stdout, config)

}

func readConfig(data string) Config {
	config := Config{}

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}

func dumpConfig(config Config) {
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- config dump:\n%s\n\n", string(d))
}

func isInList(roles []string, role string) bool {
	for _, v := range roles {
		if v == role {
			return true
		}
	}
	return false
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
	Roles Role `yaml:"roles"`
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

	keysToDel := []string{"connect_to", "hostname", "public_hostname", "ip", "public_ip", "node_labels"}
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

	return nil
}

/*
Role A map of variables to be applied to a group of Hosts
*/
type Role map[string](map[string]string)

var (
	templateGroup = template.New("")
	inventoryText = `
[OSEv3:children]
{{range $k, $v := .Deployment.Roles}}{{$k}}
{{end}}

[OSEv3:vars]
{{range $k, $v := .Deployment.Vars}}{{$k}}={{$v}}
{{end}}

{{$hosts := .Deployment.Hosts -}}
{{$roles := .Deployment.Roles -}}
{{range $role_name, $role_vars := $roles}}[{{$role_name}}]
  {{- range $hosts -}}
    {{- if IsInList .Roles $role_name -}} 
      {{- template "host" . -}} {{range $rk, $rv := $role_vars}}{{$rk}}={{$rv}}{{end -}}
    {{- end -}}
{{end}}

{{end -}}
`
	hostText = `
{{.ConnectTo}} {{if .IP}}openshift_ip={{.IP}}{{end -}} 
               {{- if .PublicIP}} openshift_public_ip={{.PublicIP}} {{end -}}
			   {{- if .Hostname}}openshift_hostname={{.Hostname}} {{end -}}
			   {{- if .PublicHostname}}openshift_public_hostname={{.PublicHostname}} {{end -}}
			   {{- range $k, $v := .Vars}}{{$k}}={{$v}} {{end -}}
`
	inventoryTemplate = template.Must(templateGroup.New("config").Funcs(template.FuncMap{
		"IsInList": isInList}).Parse(inventoryText))
	hostTemplate = template.Must(templateGroup.New("host").Parse(hostText))
)
