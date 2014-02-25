package common

var Components = []Component {}

type Component struct {
    Command string
    Description string
    Callback func(Command)
}

// Register your component with klink
func Register(component Component) {
    Components = append(Components, component)
}
