package lib

type Model struct{
    id, name string
}

func NewModel(id string, name string) *Model{
    return &Model{id:id, name:name}
}
