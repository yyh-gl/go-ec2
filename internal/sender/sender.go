package sender

type (
	Sender interface {
		Send(materials Materials) error
	}

	Material struct {
		ID    string
		Name  string
		State string
	}

	Materials []Material
)
