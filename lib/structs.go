package lib

type Model struct {
	Id   string
	Name map[string]bool
}

func NewModel(id string, name map[string]bool) *Model {
	return &Model{Id: id, Name: name}
}
