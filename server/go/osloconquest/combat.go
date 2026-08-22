package osloconquest

import (
	"errors"
	"slices"
)

var (
	ErrInvalidAttack           = errors.New("invalid attack: destination is not an enemy territory")
	ErrInvalidAttackNotOwned   = errors.New("invalid attack: attacking territory not owned by attacker")
	ErrInvalidAttackNotPresent = errors.New("invalid attack: attacker not present in attacking territory")
	ErrAttackWouldLose         = errors.New("attack would result in loss: attacker has fewer units than defender")
)

func Attack(room Room, actorID PlayerID, fromID, toID MapNodeID) (Room, error) {
	if room.ActivePlayerID == nil || *room.ActivePlayerID != actorID {
		return room, ErrNotActivePlayer
	}

	if room.Phase != PhasePlaying {
		return room, ErrNotPlaying
	}
	room = room.Clone()

	attackingTerritory := room.Territories[fromID]
	if attackingTerritory.Owner == nil || *attackingTerritory.Owner != actorID {
		return room, ErrInvalidAttackNotOwned
	}
	attackingPlayer := playerByID(room.Players, actorID)
	if attackingPlayer == nil {
		return room, ErrPlayerNotFound
	}
	if attackingPlayer.Position == nil ||
		*attackingPlayer.Position != fromID {
		return room, ErrInvalidAttackNotPresent
	}

	defendingTerritory := room.Territories[toID]
	if defendingTerritory.Owner != nil && *defendingTerritory.Owner == actorID {
		return room, ErrInvalidAttack
	}

	if !slices.Contains(attackingPlayer.ValidMoves, toID) {
		return room, ErrInvalidAttack
	}

	if attackingTerritory.Units <= defendingTerritory.Units {
		return room, ErrAttackWouldLose
	}

	attackingPlayer.Position = &toID
	attackingPlayer.MovesRemaining = 0
	attackingPlayer.ValidMoves = nil

	defendingTerritory.Owner = &actorID
	defendingTerritory.Units = attackingTerritory.Units - 1
	attackingTerritory.Units = 1

	room.Territories[fromID] = attackingTerritory
	room.Territories[toID] = defendingTerritory

	return room, nil
}
