package domain

import (
	"errors"
	"math"
	"slices"
	"strings"
	"time"
)

const (
	fuelTankCapacity                = 15
	initialFuelAmount               = fuelTankCapacity / 2
	travelFuelConsumption           = 1
	initialCoinsAmount              = 100
	refuelCoinCostPerUnit           = 10
	anchoringCoinCost               = 20
	migrationMinAcceptableKnowledge = 50
	migrationCoinCost               = 80
)

type Player struct {
	UserId             int32     `json:"-"`
	AtTerritory        string    `json:"atTerritory"`
	AtIsland           string    `json:"atIsland"`
	Anchored           bool      `json:"anchored"`
	Fuel               int32     `json:"fuel"`
	FuelCap            int32     `json:"fuelCap"`
	Coins              int32     `json:"coins"`
	BlueKeys           int32     `json:"blueKeys"`
	RedKeys            int32     `json:"redKeys"`
	GoldenKeys         int32     `json:"goldenKeys"`
	VisitedTerritories []string  `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

type FullPlayer struct {
	Player
	KnowledgeBars []KnowledgeBar `json:"knowledgeBars"`
}

const (
	PlayerUpdateEventInitial    = "initial"
	PlayerUpdateEventTravel     = "travel"
	PlayerUpdateEventRefuel     = "refuel"
	PlayerUpdateEventCorrection = "correction"
	PlayerUpdateEventAnchor     = "anchor"
	PlayerUpdateEventMigration  = "migration"
)

type PlayerUpdateEvent struct {
	Reason string
	Player *Player
}

type FullPlayerUpdateEvent struct {
	Reason string      `json:"reason"`
	Player *FullPlayer `json:"player"`
}

func NewPlayer(userId int32, startingTerritory *Territory) Player {
	return Player{
		UserId:             userId,
		AtTerritory:        startingTerritory.ID,
		AtIsland:           startingTerritory.StartIsland,
		Anchored:           true,
		Fuel:               initialFuelAmount,
		FuelCap:            fuelTankCapacity,
		Coins:              initialCoinsAmount,
		VisitedTerritories: []string{startingTerritory.ID},
	}
}

type Cost struct {
	Items []CostItem `json:"items"`
}

const (
	CostItemTypeFuel      = "fuel"
	CostItemTypeCoin      = "coin"
	CostItemTypeBlueKey   = "blueKey"
	CostItemTypeRedKey    = "redKey"
	CostItemTypeGoldenKey = "goldenKey"
)

type CostItem struct {
	Type   string `json:"type"`
	Amount int32  `json:"amount"`
}

func canAfford(player Player, cost Cost) bool {
	_, ok := deduceCost(player, cost)
	return ok
}

func deduceCost(player Player, cost Cost) (Player, bool) {
	for _, o := range cost.Items {
		switch o.Type {
		case CostItemTypeFuel:
			if player.Fuel < o.Amount {
				return player, false
			}
			player.Fuel -= o.Amount
		case CostItemTypeCoin:
			if player.Coins < o.Amount {
				return player, false
			}
			player.Coins -= o.Amount
		default:
			return player, false
		}
	}
	return player, true
}

type TravelCheckResult struct {
	Feasible bool `json:"feasible"`
	// Deprecated
	FuelCost   int32  `json:"fuelCost"`
	TravelCost Cost   `json:"travelCost"`
	Reason     string `json:"reason,omitempty"`
}

func TravelCheck(player Player, fromIsland, toIsland string, territory *Territory, isDestinationIslandUnlocked bool) (result TravelCheckResult) {
	result = TravelCheckResult{
		Feasible:   false,
		FuelCost:   travelFuelConsumption,
		TravelCost: Cost{Items: []CostItem{{Type: CostItemTypeFuel, Amount: travelFuelConsumption}}},
	}
	if !isDestinationIslandUnlocked {
		result.Reason = "باید قبل از سفر به این جزیره، سؤالات جزایر قبلی آن را پاسخ دهید."
		return
	}
	if player.AtIsland != fromIsland {
		result.Reason = "شما در جزیره اعلامی قرار ندارید."
		return
	}
	if !slices.ContainsFunc(territory.Edges, func(edge Edge) bool {
		return (edge.From == fromIsland && edge.To == toIsland) ||
			(edge.From == toIsland && edge.To == fromIsland)
	}) {
		result.Reason = "مسیر مستقیمی وجود ندارد."
		return
	}
	if !canAfford(player, result.TravelCost) {
		result.Reason = "سوخت کافی نیست."
		return
	}
	result.Feasible = true
	return
}

func Travel(player Player, fromIsland, toIsland string, territory *Territory, isDestinationIslandUnlocked bool) (*PlayerUpdateEvent, error) {
	check := TravelCheck(player, fromIsland, toIsland, territory, isDestinationIslandUnlocked)
	if !check.Feasible {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   check.Reason,
		}
	}

	player.Fuel -= travelFuelConsumption
	player.AtIsland = toIsland
	player.Anchored = false
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventTravel,
		Player: &player,
	}, nil
}

type RefuelCheckResult struct {
	MaxAvailableAmount int32  `json:"maxAvailableAmount"`
	CoinCostPerUnit    int32  `json:"coinCostPerUnit"`
	MaxReason          string `json:"maxReason"`
}

func RefuelCheck(player Player, territory *Territory) (result RefuelCheckResult) {
	idx := slices.IndexFunc(territory.RefuelIslands, func(ri RefuelIsland) bool {
		return ri.ID == player.AtIsland
	})
	if idx < 0 {
		result.MaxReason = "شما در حال حاضر در جزیره سوخت‌گیری قرار ندارید."
		return
	}
	result.CoinCostPerUnit = refuelCoinCostPerUnit
	fuelCapBound := player.FuelCap - player.Fuel
	coinBound := int32(math.MaxInt32)
	if result.CoinCostPerUnit > 0 {
		coinBound = player.Coins / result.CoinCostPerUnit
	}
	if coinBound < fuelCapBound {
		result.MaxAvailableAmount = coinBound
		result.MaxReason = "موجودی سکه شما تنها برای خرید این میزان سوخت کافی است."
	} else {
		result.MaxAvailableAmount = fuelCapBound
		result.MaxReason = "باک شما بیش از این مقدار گنجایش ندارد."
	}
	return
}

func Refuel(player Player, territory *Territory, amount int32) (*PlayerUpdateEvent, error) {
	check := RefuelCheck(player, territory)
	if amount <= 0 {
		return nil, Error{
			text:   "Invalid refuel amount",
			reason: ErrorReasonRuleViolation,
		}
	}
	if check.MaxAvailableAmount < amount {
		return nil, Error{
			text:   check.MaxReason,
			reason: ErrorReasonRuleViolation,
		}
	}
	player.Fuel += amount
	player.Coins -= amount * check.CoinCostPerUnit
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventRefuel,
		Player: &player,
	}, nil
}

type AnchorCheckResult struct {
	Feasible      bool   `json:"feasible"`
	AnchoringCost Cost   `json:"anchoringCost"`
	Reason        string `json:"reason,omitempty"`
}

func AnchorCheck(player Player, islandID string) (result AnchorCheckResult) {
	result.AnchoringCost = Cost{Items: []CostItem{{Type: CostItemTypeCoin, Amount: anchoringCoinCost}}}
	if player.AtIsland != islandID {
		result.Reason = "باید به جزیره سفر کنید تا بتوانید در آن لنگر بیندازید."
		return
	}
	if player.Anchored {
		result.Reason = "در حال حاضر در این جزیره لنگر انداخته اید."
		return
	}

	if !canAfford(player, result.AnchoringCost) {
		result.Reason = "دارایی شما برای لنگر انداختن کافی نیست."
		return
	}

	result.Feasible = true
	return
}

func Anchor(player Player, islandID string) (*PlayerUpdateEvent, error) {
	check := AnchorCheck(player, islandID)
	if !check.Feasible {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   check.Reason,
		}
	}
	var ok bool
	player, ok = deduceCost(player, check.AnchoringCost)
	if !ok {
		return nil, errors.New("logical error in anchor")
	}
	player.Anchored = true
	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventAnchor,
		Player: &player,
	}, nil
}

type MigrateCheckResult struct {
	KnowledgeCriteriaTerritory string                     `json:"knowledgeCriteriaTerritory"`
	KnowledgeValue             int32                      `json:"knowledgeValue"`
	TerritoryMigrationOptions  []TerritoryMigrationOption `json:"territoryMigrationOptions"`
}

const (
	TerritoryMigrationStatusResident  = "resident"
	TerritoryMigrationStatusUntouched = "untouched"
	TerritoryMigrationStatusVisited   = "visited"
)

type TerritoryMigrationOption struct {
	TerritoryID            string `json:"territoryId"`
	TerritoryName          string `json:"territoryName"`
	Status                 string `json:"status"`
	MinAcceptableKnowledge int32  `json:"minAcceptableKnowledge"`
	MigrationCost          Cost   `json:"migrationCost"`
	MustPayCost            bool   `json:"mustPayCost"`
	Feasible               bool   `json:"feasible"`
	Reason                 string `json:"reason,omitempty"`
}

func MigrateCheck(player Player, knowledgeBars []KnowledgeBar, currentTerritory Territory, territories []Territory) (result MigrateCheckResult) {
	atTerminalIsland := slices.ContainsFunc(currentTerritory.TerminalIslands, func(t TerminalIsland) bool {
		return t.ID == player.AtIsland
	})

	result.KnowledgeCriteriaTerritory = player.VisitedTerritories[len(player.VisitedTerritories)-1]
	for _, b := range knowledgeBars {
		if b.TerritoryID == result.KnowledgeCriteriaTerritory {
			result.KnowledgeValue = b.Value
			break
		}
	}

	for _, t := range territories {
		var knowledgeBar KnowledgeBar
		for _, b := range knowledgeBars {
			if t.ID == b.TerritoryID {
				knowledgeBar = b
				break
			}
		}
		result.TerritoryMigrationOptions = append(result.TerritoryMigrationOptions,
			getMigrationOption(player, result.KnowledgeValue, t, knowledgeBar, atTerminalIsland))
	}

	// order options based on state. break tie with territory id.
	order := []string{TerritoryMigrationStatusVisited, TerritoryMigrationStatusResident, TerritoryMigrationStatusUntouched}
	slices.SortFunc(result.TerritoryMigrationOptions, func(a, b TerritoryMigrationOption) int {
		diff := slices.Index(order, a.Status) - slices.Index(order, b.Status)
		if diff != 0 {
			return diff
		}
		return strings.Compare(a.TerritoryID, b.TerritoryID)
	})

	return
}

func getMigrationOption(player Player, knowledgeValue int32, territory Territory, knowledgeBar KnowledgeBar, atTerminalIsland bool) (option TerritoryMigrationOption) {
	option = TerritoryMigrationOption{
		TerritoryID:            territory.ID,
		TerritoryName:          territory.Name,
		Status:                 TerritoryMigrationStatusUntouched,
		MinAcceptableKnowledge: min(knowledgeBar.Total, migrationMinAcceptableKnowledge),
	}
	if player.AtTerritory == territory.ID {
		option.Status = TerritoryMigrationStatusResident
	} else if slices.Contains(player.VisitedTerritories, territory.ID) {
		option.Status = TerritoryMigrationStatusVisited
	}

	option.MigrationCost = Cost{Items: []CostItem{{Type: CostItemTypeCoin, Amount: migrationCoinCost}}}
	option.MustPayCost = option.Status == TerritoryMigrationStatusUntouched && knowledgeValue < option.MinAcceptableKnowledge

	if option.Status == TerritoryMigrationStatusResident {
		option.Reason = "شما در همین قلمرو قرار دارید."
		return
	}

	if !atTerminalIsland {
		option.Reason = "شما در جزیره شاهراه قرار ندارید."
		return
	}

	if option.MustPayCost {
		_, ok := deduceCost(player, option.MigrationCost)
		if !ok {
			option.Reason = "شما دانش یا سکه کافی برای مهاجرت ندارید."
			return
		}
	}

	option.Feasible = true
	return
}

func Migrate(player Player, knowledgeBars []KnowledgeBar, currentTerritory Territory, territories []Territory, toTerritory string) (*PlayerUpdateEvent, error) {
	check := MigrateCheck(player, knowledgeBars, currentTerritory, territories)
	idx := slices.IndexFunc(check.TerritoryMigrationOptions, func(state TerritoryMigrationOption) bool {
		return state.TerritoryID == toTerritory
	})
	if idx < 0 {
		return nil, Error{
			reason: ErrorReasonResourceNotFound,
			text:   "destination territory not found",
		}
	}
	chosenOption := check.TerritoryMigrationOptions[idx]
	if !chosenOption.Feasible {
		return nil, Error{
			reason: ErrorReasonRuleViolation,
			text:   chosenOption.Reason,
		}
	}

	if chosenOption.MustPayCost {
		var ok bool
		player, ok = deduceCost(player, chosenOption.MigrationCost)
		if !ok {
			return nil, Error{
				reason: ErrorReasonRuleViolation,
				text:   "cannot not afford cost",
			}
		}
	}

	idx = slices.IndexFunc(territories, func(t Territory) bool {
		return t.ID == toTerritory
	})
	if idx < 0 {
		return nil, errors.New("did not find toTerritory in territories")
	}
	destinationTerritory := territories[idx]

	player.AtTerritory = destinationTerritory.ID
	player.AtIsland = destinationTerritory.StartIsland
	if !slices.Contains(player.VisitedTerritories, toTerritory) {
		player.VisitedTerritories = append(player.VisitedTerritories, destinationTerritory.ID)
	}
	player.Anchored = false

	return &PlayerUpdateEvent{
		Reason: PlayerUpdateEventMigration,
		Player: &player,
	}, nil
}
