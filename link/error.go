package link

type ErrorFoldedDirectory struct {
	Message   string
	Dot       string
	FoldedDir string
}

func (e *ErrorFoldedDirectory) Error() string {
	return e.Message
}
