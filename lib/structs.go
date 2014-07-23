package lib

type Model struct{
    Id, Name string
}


func NewModel(id string, name string) *Model{
    return &Model{Id:id, Name:name}
}
