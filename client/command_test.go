package client_test

import (
	"fmt"
	"strings"
	"testing"
)

type test struct {
	commands []string
}

// CommandLine executes commands
type MockCommandLine struct {
	Commands []string
}

// Execute executes the specified commands
func (c *MockCommandLine) Execute(name string, args []string) (string, error) {
	builder := strings.Builder{}
	builder.WriteString(name)
	for _, a := range args {
		builder.WriteString(fmt.Sprintf(" %s", a))
	}
	fmt.Println(builder.String())
	c.Commands = append(c.Commands, builder.String())
	return "", nil
}

func testCommands(t *testing.T, testIndex int, actualCommands []string, expectedCommands []string) {
	if len(actualCommands) != len(expectedCommands) {
		t.Errorf("test %d: number of commands not matching (expect %d, actual %d)", testIndex, len(expectedCommands), len(actualCommands))
	}
	for j, c := range actualCommands {
		if c != expectedCommands[j] {
			t.Errorf("test %d: command %d: does not match (expect %s, actual %s)", testIndex, j, expectedCommands[j], c)
		}
	}
}
