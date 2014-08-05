package lib

type Model struct {
	Id   string
	Name map[string]bool
}

type Product struct {
	Id          string `json:"id"`
	Category_id string `json:"category_id"`
	Model_id    string `json:"model_id"`
	Name        string `json:"name"`
}

type MatchData struct {
	Callback_url            string
	Callback_model_id_param string
	Products                []Product
}

func NewModel(id string, name map[string]bool) *Model {
	return &Model{Id: id, Name: name}
}
