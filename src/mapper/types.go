package mapper

type Command string

const (
	Add            Command = "add"
	Update         Command = "update"
	Delete         Command = "delete"
	MarkInProgress Command = "mark-in-progress"
	MarkDone       Command = "mark-done"
	List           Command = "list"
)
