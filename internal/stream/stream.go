package stream

// Opcode is used to specify allowed server operations.
type Opcode uint8

// Result represents the result of the attempted operation.
type result uint8

const (
	// Undefined represents an undefined operation.
	Undefined Opcode = iota
	// Find represents the directive to find snippets.
	Find
	// List represents the action of listing out snippets.
	List
	// Insert represents the operaion of inserting snippets.
	Insert
	// Delete represents the directive to delete snippets.
	Delete
	// Reload represents the directive to reload the snippet container.
	Reload
)

const (
	// Success means that the action passed without any issues.
	Success result = iota
	// Failure means that the action failed to be executed.
	Failure
)

// Request defines the data format for the server request.
type Request struct {
	Operation Opcode `json:"operation"`
	Body      []byte `json:"body"`
}

// Reply defines the data format for ther server reply.
type Reply struct {
	Result result `json:"result"`
	Body   []byte `json:"body"`
}
