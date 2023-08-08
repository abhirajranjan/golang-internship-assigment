package model

type Repo interface {
	GetAllPlayerRankwise() ([]Player, error)
	GetPlayerByRank(int) (Player, error)
	GetRandomPlayer() (Player, error)
	DeletePlayer(int) error
	CreateNewPlayer(*Player) error
	UpdatePlayer(id int, name string, score int, check map[string]interface{}) (Player, error)
}
