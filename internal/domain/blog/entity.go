package blog

const ()

type Blog struct {
	ID          uint
	Sysname     string
	KeywordIDs  []uint
	TagIDs      []uint
	Name        string
	Description string
}

func (e *Blog) Validate() error {
	return nil
}
