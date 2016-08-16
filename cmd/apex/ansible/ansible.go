package cf

import (
	"strings"

	"github.com/apex/apex/ansible"
	"github.com/apex/apex/cmd/apex/root"
	"github.com/apex/apex/stats"
	"github.com/tj/cobra"
)

// example output.
const example = `
    Preview using dry run
    $ apex ansible apex-playbook.yml --check

    Apply changes
    $ apex ansible apex-playbook.yml`

// Command config.
var Command = &cobra.Command{
	Use:     "ansible",
	Short:   "ansible management",
	Example: example,
	RunE:    run,
}

// Initialize.
func init() {
	root.Register(Command)
}

// Run command.
func run(c *cobra.Command, args []string) error {
	stats.Track("Ansible", nil)

	err := root.Project.LoadFunctions()

	// Hack to prevent initial `apex ansible` from failing,
	// as we load functions to expose their ARNs.
	if err != nil {
		if !strings.Contains(err.Error(), "Role: zero value") {
			return err
		}
	}

	p := &ansible.Proxy{
		Functions:   root.Project.Functions,
		Region:      *root.Session.Config.Region,
		Environment: root.Project.InfraEnvironment,
		Role:        root.Project.Role,
	}

	return p.Run(args...)
}
