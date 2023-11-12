package tag

const ()

type Tag struct {
	ID      uint
	Sysname string
	Value   string
}

func (e *Tag) Validate() error {
	return nil
}
