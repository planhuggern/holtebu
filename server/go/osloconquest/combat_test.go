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

func TestAttackConquersNeutralTerritory(t *testing.T) {
	scenario := newAttackScenario(t, 3, 0)
	neutralTerritory := scenario.room.Territories[scenario.toID]
	neutralTerritory.Owner = nil
	scenario.room.Territories[scenario.toID] = neutralTerritory

	updatedRoom, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, scenario.toID)
	if err != nil {
		t.Fatalf("Attack returned an error: %v", err)
	}

	conqueredTerritory := updatedRoom.Territories[scenario.toID]
	if conqueredTerritory.Owner == nil {
		t.Fatal("conquered territory owner = nil, want attacker ID")
	}
	if *conqueredTerritory.Owner != scenario.attackerID {
		t.Errorf("conquered territory owner = %q, want %q", *conqueredTerritory.Owner, scenario.attackerID)
	}
	if conqueredTerritory.Units != 2 {
		t.Errorf("conquered territory units = %d, want 2", conqueredTerritory.Units)
	}
}

func TestAttackRejectsSourceThatIsNotAttackingPlayerPosition(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)
	remoteTerritoryID := MapNodeID("t3")
	scenario.room.Territories[remoteTerritoryID] = MapNode{
		ID:    remoteTerritoryID,
		Owner: &scenario.attackerID,
		Units: 3,
	}

	updatedRoom, err := Attack(
		scenario.room,
		scenario.attackerID,
		remoteTerritoryID,
		scenario.toID,
	)

	if err != ErrInvalidAttackNotPresent {
		t.Errorf("error = %v, want ErrInvalidAttackNotPresent", err)
	}
	remoteTerritory := updatedRoom.Territories[remoteTerritoryID]
	if remoteTerritory.Units != 3 {
		t.Errorf("remote territory units = %d, want 3", remoteTerritory.Units)
	}
	defendingTerritory := updatedRoom.Territories[scenario.toID]
	if defendingTerritory.Owner == nil {
		t.Fatal("defending territory owner = nil, want original defender ID")
	}
	if *defendingTerritory.Owner != scenario.defenderID {
		t.Errorf("defending territory owner = %q, want %q", *defendingTerritory.Owner, scenario.defenderID)
	}
}

func TestAttackRejectsPlayerWithoutPosition(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)
	attackingPlayer := playerByID(scenario.room.Players, scenario.attackerID)
	if attackingPlayer == nil {
		t.Fatal("attacking player not found")
	}
	attackingPlayer.Position = nil

	updatedRoom, err := Attack(
		scenario.room,
		scenario.attackerID,
		scenario.fromID,
		scenario.toID,
	)

	if err != ErrInvalidAttackNotPresent {
		t.Errorf("error = %v, want ErrInvalidAttackNotPresent", err)
	}
	if updatedRoom.Territories[scenario.fromID].Units != 3 {
		t.Errorf("attacker units = %d, want 3", updatedRoom.Territories[scenario.fromID].Units)
	}
	defendingTerritory := updatedRoom.Territories[scenario.toID]
	if defendingTerritory.Owner == nil {
		t.Fatal("defending territory owner = nil, want original defender ID")
	}
	if *defendingTerritory.Owner != scenario.defenderID {
		t.Errorf("defending territory owner = %q, want %q", *defendingTerritory.Owner, scenario.defenderID)
	}
}

func TestAttackRejectsSetupPhase(t *testing.T) {
	room := fillRoomWithPlayers(MaxPlayers)
	actorID := *room.ActivePlayerID

	updatedRoom, err := Attack(room, actorID, "t1", "t2")

	if err != ErrNotPlaying {
		t.Errorf("error = %v, want ErrNotPlaying", err)
	}
	if updatedRoom.Phase != PhaseSetup {
		t.Errorf("phase = %q, want %q", updatedRoom.Phase, PhaseSetup)
	}
}

func TestAttackSetsWinnerWhenConquestReachesSixtyPercentThreshold(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)
	scenario.room.Territories["t3"] = MapNode{
		ID:    "t3",
		Owner: &scenario.attackerID,
		Units: 1,
	}
	for _, territoryID := range []MapNodeID{"t4", "t5"} {
		scenario.room.Territories[territoryID] = MapNode{
			ID:    territoryID,
			Owner: &scenario.defenderID,
			Units: 1,
		}
	}

	updatedRoom, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, scenario.toID)
	if err != nil {
		t.Fatalf("Attack returned an error: %v", err)
	}

	if updatedRoom.WinnerID == nil {
		t.Fatal("winner ID = nil, want attacker ID")
	}
	if *updatedRoom.WinnerID != scenario.attackerID {
		t.Errorf("winner ID = %q, want %q", *updatedRoom.WinnerID, scenario.attackerID)
	}
}

func TestAttackDoesNotSetWinnerBelowControlThreshold(t *testing.T) {
	scenario := newAttackScenario(t, 3, 2)
	for _, territoryID := range []MapNodeID{"t3", "t4", "t5"} {
		scenario.room.Territories[territoryID] = MapNode{
			ID:    territoryID,
			Owner: &scenario.defenderID,
			Units: 1,
		}
	}

	updatedRoom, err := Attack(scenario.room, scenario.attackerID, scenario.fromID, scenario.toID)
	if err != nil {
		t.Fatalf("Attack returned an error: %v", err)
	}

	if updatedRoom.WinnerID != nil {
		t.Errorf("winner ID = %q, want nil", *updatedRoom.WinnerID)
	}
	if updatedRoom.Phase != PhasePlaying {
		t.Errorf("phase = %q, want %q", updatedRoom.Phase, PhasePlaying)
	}
}
