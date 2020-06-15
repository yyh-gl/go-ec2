package sender

type (
	Sender interface {
		Send(name string, materials Materials) error
	}

	Material struct {
		ID    string
		Name  string
		State string
	}

	Materials []Material
)
