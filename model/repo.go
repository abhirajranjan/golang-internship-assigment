package model

type Repo interface {
	GetAllPlayerRankwise() ([]Player, error)
	GetPlayerByRank(int) (Player, error)
	GetRandomPlayer() (Player, error)
	DeletePlayer(int) error
	CreateNewPlayer(*Player) error
	UpdatePlayer(*Player) error
}
