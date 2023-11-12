package keyword

const ()

type Keyword struct {
	ID      uint
	Sysname string
	Value   string
}

func (e *Keyword) Validate() error {
	return nil
}
