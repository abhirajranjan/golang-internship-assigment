package memorydb

import (
	"math/rand"
	"score/model"
	"sync"

	"github.com/pkg/errors"
)

type memorydb struct {
	mu         *sync.RWMutex
	playerCard map[int]model.Player // map of player id to players with rank
	rank       []int                // rank of player ids
	nextid     int                  // next id to give to new player
}

// validates that memorydb implements model.Repo
var _ model.Repo = (*memorydb)(nil)

func NewMemoryDB() *memorydb {
	return &memorydb{
		mu:         &sync.RWMutex{},
		playerCard: make(map[int]model.Player),
		rank:       make([]int, 0),
		nextid:     1,
	}
}

func (m *memorydb) GetAllPlayerRankwise() ([]model.Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	players := make([]model.Player, len(m.playerCard))

	if len(m.playerCard) == 0 {
		return players, model.ErrNoPlayer
	}

	for pos, id := range m.rank {
		p, ok := m.playerCard[id]
		if !ok {
			return nil, errors.Errorf("ranging error %d\n\n%#v\n\n%#v", id, m.playerCard, m.rank)
		}
		players[pos] = p
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

	return m.playerCard[m.rank[pos]], nil
}
func (m *memorydb) GetRandomPlayer() (model.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.rank) == 0 {
		return model.Player{}, model.ErrNoPlayer
	}

	n := rand.Intn(len(m.rank))
	// values in m.players cannot be altered hence safe to return them
	return m.playerCard[m.rank[n]], nil
}
func (m *memorydb) DeletePlayer(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.playerCard[id]
	if !ok {
		return model.ErrInvalidPlayer
	}

	// player card stores rank with 1 based indexing
	m.deletePlayerFromRank(p.Id)
	m.deletePlayerCardbyID(id)
	return nil
}

func (m *memorydb) deletePlayerFromRank(id int) {
	var rank int = len(m.rank)

	for idx, i := range m.rank {
		if i == id {
			rank = idx
			break
		}
	}

	if rank == len(m.rank) {
		return
	}

	// move every element one step forward
	copy(m.rank[rank:], m.rank[rank+1:])
	// update rank slice to exclude last element
	m.rank = m.rank[:len(m.rank)-1]
}

func (m *memorydb) deletePlayerCardbyID(id int) {
	delete(m.playerCard, id)
}

func (m *memorydb) CreateNewPlayer(p *model.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var player model.Player

	player.Country = p.Country
	player.Id = m.nextid
	p.Id = player.Id
	m.nextid++
	player.Name = p.Name
	player.Score = p.Score

	m.createPlayer(player)
	return nil
}

func (m *memorydb) createPlayer(player model.Player) {
	m.rank, _ = sortedinsert(m.rank, player.Score, player.Id, func(rank int) int {
		card := m.playerCard[rank]
		return card.Score
	})

	m.playerCard[player.Id] = player
}

func (m *memorydb) UpdatePlayer(id int, name string, score int, check map[string]interface{}) (model.Player, error) {
	var (
		player model.Player
	)

	m.mu.Lock()
	defer m.mu.Unlock()

	playerCard, ok := m.playerCard[id]
	if !ok {
		return player, model.ErrInvalidPlayer
	}

	// update the rank only if there is change in scores
	if playerCard.Score != score {
		m.deletePlayerFromRank(playerCard.Id)
	}

	// remove player card as they are immutable in maps (not referenced type)
	m.deletePlayerCardbyID(id)

	if _, ok := check["name"]; !ok {
		name = playerCard.Name
	}
	if _, ok := check["score"]; !ok {
		score = playerCard.Score
	}

	player = model.Player{
		Id:      id,
		Name:    name,
		Score:   score,
		Country: playerCard.Country,
	}

	m.createPlayer(player)
	return player, nil
}
