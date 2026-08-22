package osloconquest

import "testing"

type attackScenario struct {
	room       Room
	attackerID PlayerID
	defenderID PlayerID
	fromID     MapNodeID
	toID       MapNodeID
}

func newAttackScenario(t *testing.T, attackerUnits, defenderUnits int) attackScenario {
	t.Helper()

	room := playingRoomWithPlayers(MaxPlayers, "kolbotn_cp", t)
	attackerID := *room.ActivePlayerID
	defenderID := room.Players[1].ID
	fromID := MapNodeID("t1")
	toID := MapNodeID("t2")
	attackingPlayer := playerByID(room.Players, attackerID)
	if attackingPlayer == nil {
		t.Fatal("attacking player not found")
	}
	attackingPlayer.Position = &fromID
	attackingPlayer.ValidMoves = []MapNodeID{toID}
	attackingPlayer.MovesRemaining = 3
	room.Territories = map[MapNodeID]MapNode{
		fromID: {ID: fromID, Owner: &attackerID, Units: attackerUnits},
		toID:   {ID: toID, Owner: &defenderID, Units: defenderUnits},
	}

	return attackScenario{
		room:       room,
		attackerID: attackerID,
		defenderID: defenderID,
		fromID:     fromID,
		toID:       toID,
	}
}

func TestAttackConquersTerritoryWithinDiceRangeWhenAttackerHasMoreUnits(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)

	room, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, scenario.toID)
	if err != nil {
		t.Fatalf("Attack returned an error: %v", err)
	}

	attacker := room.Territories[scenario.fromID]
	if attacker.Units != 1 {
		t.Errorf("attacker units = %d, want 1", attacker.Units)
	}

	defender := room.Territories[scenario.toID]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want attacker ID")
	}
	if *defender.Owner != scenario.attackerID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, scenario.attackerID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}

	attackingPlayer := playerByID(room.Players, scenario.attackerID)
	if attackingPlayer == nil || attackingPlayer.Position == nil {
		t.Fatal("attacking player position = nil, want destination")
	}
	if *attackingPlayer.Position != scenario.toID {
		t.Errorf("attacking player position = %q, want %q", *attackingPlayer.Position, scenario.toID)
	}
	if attackingPlayer.MovesRemaining != 0 {
		t.Errorf("moves remaining = %d, want 0", attackingPlayer.MovesRemaining)
	}
	if len(attackingPlayer.ValidMoves) != 0 {
		t.Errorf("valid moves = %v, want []", attackingPlayer.ValidMoves)
	}
}

func TestAttackRejectsDestinationOutsideDiceRangeWithoutChangingState(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)
	outOfRangeID := MapNodeID("t35")
	scenario.room.Territories[outOfRangeID] = MapNode{
		ID:    outOfRangeID,
		Owner: &scenario.defenderID,
		Units: 2,
	}

	updatedRoom, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, outOfRangeID)

	if err != ErrInvalidAttack {
		t.Errorf("error = %v, want ErrInvalidAttack", err)
	}
	attacker := updatedRoom.Territories[scenario.fromID]
	if attacker.Units != 3 {
		t.Errorf("attacker units = %d, want 3", attacker.Units)
	}
	defender := updatedRoom.Territories[outOfRangeID]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want original defender ID")
	}
	if *defender.Owner != scenario.defenderID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, scenario.defenderID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}
}

func TestAttackRejectsWhenAttackerDoesNotHaveMoreUnitsWithoutChangingState(t *testing.T) {
	scenario := newAttackScenario(t, 2, 2)

	updatedRoom, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, scenario.toID)

	if err != ErrAttackWouldLose {
		t.Errorf("error = %v, want ErrAttackWouldLose", err)
	}
	attacker := updatedRoom.Territories[scenario.fromID]
	if attacker.Units != 2 {
		t.Errorf("attacker units = %d, want 2", attacker.Units)
	}
	defender := updatedRoom.Territories[scenario.toID]
	if defender.Owner == nil {
		t.Fatal("defender owner = nil, want original defender ID")
	}
	if *defender.Owner != scenario.defenderID {
		t.Errorf("defender owner = %q, want %q", *defender.Owner, scenario.defenderID)
	}
	if defender.Units != 2 {
		t.Errorf("defender units = %d, want 2", defender.Units)
	}
	attackingPlayer := playerByID(updatedRoom.Players, scenario.attackerID)
	if attackingPlayer == nil || attackingPlayer.Position == nil {
		t.Fatal("attacking player position = nil, want original position")
	}
	if *attackingPlayer.Position != scenario.fromID {
		t.Errorf("attacking player position = %q, want %q", *attackingPlayer.Position, scenario.fromID)
	}
	if attackingPlayer.MovesRemaining != 3 {
		t.Errorf("moves remaining = %d, want 3", attackingPlayer.MovesRemaining)
	}
	if len(attackingPlayer.ValidMoves) != 1 || attackingPlayer.ValidMoves[0] != scenario.toID {
		t.Errorf("valid moves = %v, want [%q]", attackingPlayer.ValidMoves, scenario.toID)
	}
}
