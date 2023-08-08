package memorydb

import (
	"math/rand"
	"score/model"
	"sync"
)

type playerwithrank struct {
	player model.Player
	rank   int // 1 based indexing
}

type memorydb struct {
	mu         *sync.RWMutex
	playerCard map[int]playerwithrank // map of player id to players with rank
	rank       []int                  // rank of player ids
	nextid     int                    // next id to give to new player
}

// validates that memorydb implements model.Repo
var _ model.Repo = (*memorydb)(nil)

func NewMemoryDB() *memorydb {
	return &memorydb{
		mu:         &sync.RWMutex{},
		playerCard: make(map[int]playerwithrank),
		rank:       make([]int, 0),
		nextid:     0,
	}
}

func (m *memorydb) GetAllPlayerRankwise() ([]model.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.playerCard) == 0 {
		return nil, model.ErrNoPlayer
	}

	players := make([]model.Player, len(m.playerCard))

	for pos, id := range m.rank {
		players[pos] = m.playerCard[id].player
	}
	return players, nil
}

func (m *memorydb) GetPlayerByRank(rank int) (model.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// if rank=1 then position in array will be 0 (0 based indexing)
	pos := rank - 1

	var player model.Player

	if len(m.rank) <= pos {
		return player, model.ErrRankDoensnotExist
	}

	return m.playerCard[m.rank[pos]].player, nil
}
func (m *memorydb) GetRandomPlayer() (model.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	n := rand.Intn(len(m.rank) - 1)
	// values in m.players cannot be altered hence safe to return them
	return m.playerCard[m.rank[n]].player, nil
}
func (m *memorydb) DeletePlayer(id int) error {
	m.mu.Lock()

	p, ok := m.playerCard[id]
	if !ok {
		return model.ErrInvalidPlayer
	}

	m.mu.Unlock()
	// player card stores rank with 1 based indexing
	m.deletePlayerFromRank(p.rank - 1)
	m.deletePlayerCardbyID(p.player.Id)
	return nil
}

func (m *memorydb) deletePlayerFromRank(rank int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// move every element one step forward
	copy(m.rank[rank:], m.rank[rank+1:])
	// update rank slice to exclude last element
	m.rank = m.rank[:len(m.rank)-1]
}

func (m *memorydb) deletePlayerCardbyID(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.playerCard, id)
}

func (m *memorydb) CreateNewPlayer(p *model.Player) error {
	m.mu.Lock()

	var player model.Player

	player.Country = p.Country
	player.Id = m.nextid
	m.nextid++
	player.Name = p.Name
	player.Score = p.Score

	m.mu.Unlock()

	m.createPlayer(player)
	return nil
}

func (m *memorydb) createPlayer(player model.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var rank int
	m.rank, rank = sortedinsert(m.rank, player.Id)

	m.playerCard[player.Id] = playerwithrank{
		player: player,
		rank:   rank + 1,
	}
}

func (m *memorydb) UpdatePlayer(id int, name string, score int, check map[string]interface{}) (model.Player, error) {
	var (
		player model.Player
	)

	m.mu.Lock()

	playerCard, ok := m.playerCard[id]
	if !ok {
		return player, model.ErrInvalidPlayer
	}

	m.mu.Unlock()

	// update the rank only if there is change in scores
	if playerCard.player.Score != score {
		m.deletePlayerFromRank(playerCard.rank)
	}

	// remove player card as they are immutable in maps (not referenced type)
	m.deletePlayerCardbyID(playerCard.player.Id)

	if _, ok := check["name"]; !ok {
		name = playerCard.player.Name
	}
	if _, ok := check["score"]; !ok {
		score = playerCard.player.Score
	}

	player = model.Player{
		Id:      id,
		Name:    name,
		Score:   score,
		Country: playerCard.player.Country,
	}

	m.createPlayer(player)
	return player, nil
}
