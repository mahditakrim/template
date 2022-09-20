package entity

type Book struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Writer  string `json:"writer"`
	PageNum uint   `json:"page_num"`
}
