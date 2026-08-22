package osloconquest

import "testing"

func TestAttackConquersTerritoryWithinDiceRangeWhenAttackerHasMoreUnits(t *testing.T) {
	room := playingRoomWithPlayers(MaxPlayers, "kolbotn_cp", t)
	attackerID := *room.ActivePlayerID
	defenderID := room.Players[1].ID
	attackingPlayer := playerByID(room.Players, attackerID)
	if attackingPlayer == nil {
		t.Fatal("attacking player not found")
	}
	fromID := MapNodeID("t1")
	toID := MapNodeID("t2")
	attackingPlayer.Position = &fromID
	attackingPlayer.ValidMoves = []MapNodeID{toID}
	attackingPlayer.MovesRemaining = 3
	room.Territories = map[MapNodeID]MapNode{
		fromID: {ID: fromID, Owner: &attackerID, Units: 3},
		toID:   {ID: toID, Owner: &defenderID, Units: 2},
	}

	room, err := Attack(room, attackerID, fromID, toID)
	if err != nil {
		t.Fatalf("Attack returned an error: %v", err)
	}

	attacker := room.Territories[fromID]
	if attacker.Units != 1 {
		t.Errorf("attacker units = %d, want 1", attacker.Units)
	}

	defender := room.Territories[toID]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want attacker ID")
	}
	if *defender.Owner != attackerID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, attackerID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}

	attackingPlayer = playerByID(room.Players, attackerID)
	if attackingPlayer == nil || attackingPlayer.Position == nil {
		t.Fatal("attacking player position = nil, want destination")
	}
	if *attackingPlayer.Position != toID {
		t.Errorf("attacking player position = %q, want %q", *attackingPlayer.Position, toID)
	}
	if attackingPlayer.MovesRemaining != 0 {
		t.Errorf("moves remaining = %d, want 0", attackingPlayer.MovesRemaining)
	}
	if len(attackingPlayer.ValidMoves) != 0 {
		t.Errorf("valid moves = %v, want []", attackingPlayer.ValidMoves)
	}
}

func TestAttackRejectsDestinationOutsideDiceRangeWithoutChangingState(t *testing.T) {
	room := playingRoomWithPlayers(MaxPlayers, "kolbotn_cp", t)
	attackerID := *room.ActivePlayerID
	defenderID := room.Players[1].ID
	attackingPlayer := playerByID(room.Players, attackerID)
	if attackingPlayer == nil {
		t.Fatal("attacking player not found")
	}
	fromID := MapNodeID("t1")
	attackingPlayer.Position = &fromID
	attackingPlayer.ValidMoves = []MapNodeID{"t2"}
	attackingPlayer.MovesRemaining = 3
	room.Territories = map[MapNodeID]MapNode{
		fromID: {ID: fromID, Owner: &attackerID, Units: 3},
		"t35": {ID: "t35", Owner: &defenderID, Units: 2},
	}

	updatedRoom, err := Attack(room, attackerID, fromID, "t35")

	if err != ErrInvalidAttack {
		t.Errorf("error = %v, want ErrInvalidAttack", err)
	}
	attacker := updatedRoom.Territories[fromID]
	if attacker.Units != 3 {
		t.Errorf("attacker units = %d, want 3", attacker.Units)
	}
	defender := updatedRoom.Territories["t35"]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want original defender ID")
	}
	if *defender.Owner != defenderID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, defenderID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}
}

func TestAttackRejectsWhenAttackerDoesNotHaveMoreUnitsWithoutChangingState(t *testing.T) {
	room := playingRoomWithPlayers(MaxPlayers, "kolbotn_cp", t)
	attackerID := *room.ActivePlayerID
	defenderID := room.Players[1].ID
	attackingPlayer := playerByID(room.Players, attackerID)
	if attackingPlayer == nil {
		t.Fatal("attacking player not found")
	}
	fromID := MapNodeID("t1")
	toID := MapNodeID("t2")
	attackingPlayer.Position = &fromID
	attackingPlayer.ValidMoves = []MapNodeID{toID}
	attackingPlayer.MovesRemaining = 3
	room.Territories = map[MapNodeID]MapNode{
		fromID: {ID: fromID, Owner: &attackerID, Units: 2},
		toID:   {ID: toID, Owner: &defenderID, Units: 2},
	}

	updatedRoom, err := Attack(room, attackerID, fromID, toID)

	if err != ErrAttackWouldLose {
		t.Errorf("error = %v, want ErrAttackWouldLose", err)
	}
	attacker := updatedRoom.Territories[fromID]
	if attacker.Units != 2 {
		t.Errorf("attacker units = %d, want 2", attacker.Units)
	}
	defender := updatedRoom.Territories[toID]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want original defender ID")
	}
	if *defender.Owner != defenderID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, defenderID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}
	attackingPlayer = playerByID(updatedRoom.Players, attackerID)
	if attackingPlayer == nil || attackingPlayer.Position == nil {
		t.Fatal("attacking player position = nil, want original position")
	}
	if *attackingPlayer.Position != fromID {
		t.Errorf("attacking player position = %q, want %q", *attackingPlayer.Position, fromID)
	}
	if attackingPlayer.MovesRemaining != 3 {
		t.Errorf("moves remaining = %d, want 3", attackingPlayer.MovesRemaining)
	}
	if len(attackingPlayer.ValidMoves) != 1 || attackingPlayer.ValidMoves[0] != toID {
		t.Errorf("valid moves = %v, want [%q]", attackingPlayer.ValidMoves, toID)
	}
}
