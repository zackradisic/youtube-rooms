package ws

import "fmt"

// ActionInvoker manages and invokes actions
type ActionInvoker struct {
	actions map[string]action
}

type action interface {
	Execute(data interface{}) error
}

// InvokeAction invokes an action specified by name
func (a *ActionInvoker) InvokeAction(name string, data interface{}) error {
	if action, ok := a.actions[name]; ok {
		if err := action.Execute(data); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("could not find action: (%s)", name)
}

func (a *ActionInvoker) registerAction(name string, action action) {
	a.actions[name] = action
}

// NewActionInvoker returns a new ActionInvoker
func NewActionInvoker() *ActionInvoker {
	a := &ActionInvoker{
		actions: make(map[string]action),
	}

	return a
}

type PlayAction struct {
}
