// Package ansible proxies ansible commands.
package ansible

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apex/apex/function"
	"github.com/apex/log"
)

// Dir in which ansible configs are stored
const Dir = "ansible"

// Proxy is a wrapper around ansible-playbook command.
type Proxy struct {
	Functions   []*function.Function
	Environment string
	Region      string
	Role        string
}

// Run ansible command in ansible directory.
func (p *Proxy) Run(args ...string) error {

	args = append(args, p.functionVars()...)
	args = append(args, p.ansibleArgs()...)

	log.WithFields(log.Fields{
		"args": args,
	}).Debug("ansible")

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("AWS_REGION=%s", p.Region))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(Dir)

	return cmd.Run()
}

// functionVars returns the function ARN's as ansibles -e arguments.
func (p *Proxy) functionVars() (args []string) {
	args = append(args, "-e")
	args = append(args, fmt.Sprintf("aws_region=%s", p.Region))

	args = append(args, "-e")
	args = append(args, fmt.Sprintf("apex_environment=%s", p.Environment))

	if p.Role != "" {
		args = append(args, "-e")
		args = append(args, fmt.Sprintf("apex_function_role=%s", p.Role))
	}

	for _, fn := range p.Functions {
		config, err := fn.GetConfig()
		if err != nil {
			log.Debugf("can't fetch function config: %s", err.Error())
			continue
		}

		args = append(args, "-e")
		args = append(args, fmt.Sprintf("apex_function_%s=%s", fn.Name, *config.Configuration.FunctionArn))
	}

	return args
}

// ansibleArgs returns the ansible args to run locally
func (p *Proxy) ansibleArgs() (args []string) {
	return append(args, []string{"-i", "localhost,", "-c", "local"}...)
}

// Output fetches output variable `name` from ansible.
func Output(environment, name string) (string, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ansible output %s", name))
	cmd.Dir = filepath.Join(Dir, environment)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}
