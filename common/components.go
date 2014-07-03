package common

import (
	"fmt"
	"sort"
)

var Components = []Component{}

type Component struct {
	Command      string
	Callback     func(Command)
	Description  string
	AutoComplete string
}

func (c Component) String() string {
	return fmt.Sprintf("%s\t%s", c.Command, c.Description)
}

// make Components sortable - sometimes go is a little primative...
type ByCommand []Component

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command < a[j].Command }

// Register your components with klink
// Takes var args of common.Component
func Register(toReg ...Component) {
	for i := range toReg {
		Components = append(Components, toReg[i])
	}
	sort.Sort(ByCommand(Components))
}

// Returns a list of component names
func ComponentNames() []string {
	var names []string
	for _, component := range Components {
		names = append(names, component.Command)
	}
	return names
}
