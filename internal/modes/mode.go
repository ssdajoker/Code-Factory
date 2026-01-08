package modes

// Mode interface definition
type Mode interface {
	Name() string
	Description() string
}

// TODO: Implement INTAKE, REVIEW, CHANGE_ORDER, RESCUE modes
