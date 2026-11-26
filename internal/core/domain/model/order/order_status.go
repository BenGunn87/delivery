package order

type Status string

const (
	StatusCreated   Status = "created"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
)

func (s Status) Equal(other Status) bool {
	return s == other
}

func (s Status) String() string {
	return string(s)
}
