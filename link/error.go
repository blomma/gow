package link

// ErrorFoldedDirectory is given when a folded directory is found
type ErrorFoldedDirectory struct {
	Message   string
	Dot       string
	FoldedDir string
}

func (e *ErrorFoldedDirectory) Error() string {
	return e.Message
}

// ErrorNotOwned is given when a symlink that is not owned is found
type ErrorNotOwned struct {
	Message string
}

func (e *ErrorNotOwned) Error() string {
	return e.Message
}
