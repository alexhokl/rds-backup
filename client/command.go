package client

import (
	"fmt"
	"os/exec"

	"github.com/spf13/viper"
)

// Command indicates the requirements of executing a command
type Command interface {
	Execute(name string, args []string) (string, error)
}

// CommandLine executes commands
type CommandLine struct {
}

// Execute executes the specified commands
func (c *CommandLine) Execute(name string, args []string) (string, error) {
	if viper.GetBool("verbose") {
		fmt.Println("Command executed:", name, args)
	}
	byteOutput, err := exec.Command(name, args...).Output()
	return string(byteOutput), err

}
