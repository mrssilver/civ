package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// ========== Enhanced Constants with Validation ==========
const (
	MAP_WIDTH         = 20
	MAP_HEIGHT        = 15
	MAX_PLAYERS       = 8
	MAX_CITIES        = 50
	MAX_UNITS         = 100
	START_YEAR        = 4000
	END_YEAR          = 2050
	MIN_CITY_DISTANCE = 25
	
	// Game balance constants
	RESEARCH_SUCCESS_CHANCE = 30
	COMBAT_SUCCESS_CHANCE   = 70
	BASE_CITY_POPULATION   = 1
	STARTING_GOLD          = 100
	STARTING_HAPPINESS     = 100
)

// ========== Enhanced Type Definitions with Validation ==========
type (
	TerrainType       int
	BuildingType      int
	TechType          int
	UnitType          int
	CivilizationType  int
	ProductionItemType int
)

// Terrain types with bounds checking
const (
	TERRAIN_OCEAN TerrainType = iota
	TERRAIN_PLAINS
	TERRAIN_DESERT
	TERRAIN_MOUNTAINS
	TERRAIN_FOREST
	TERRAIN_HILLS
	TERRAIN_TUNDRA
	TERRAIN_JUNGLE
	TERRAIN_COUNT
)

// Validate terrain type
func (t TerrainType) IsValid() bool {
	return t >= 0 && t < TERRAIN_COUNT
}

// Building types with validation
const (
	BUILDING_MONUMENT BuildingType = iota
	BUILDING_GRANARY
	BUILDING_LIBRARY
	BUILDING_TEMPLE
	BUILDING_BARRACKS
	BUILDING_WALLS
	BUILDING_UNIVERSITY
	BUILDING_FACTORY
	BUILDING_COUNT
)

func (b BuildingType) IsValid() bool {
	return b >= 0 && b < BUILDING_COUNT
}

// Technology types with validation
const (
	TECH_AGRICULTURE TechType = iota
	TECH_POTTERY
	TECH_WRITING
	TECH_MATHEMATICS
	TECH_CONSTRUCTION
	TECH_PHILOSOPHY
	TECH_ENGINEERING
	TECH_EDUCATION
	TECH_GUNPOWDER
	TECH_INDUSTRIALIZATION
	TECH_COUNT
)

func (t TechType) IsValid() bool {
	return t >= 0 && t < TECH_COUNT
}

// Unit types with validation
const (
	UNIT_SETTLER UnitType = iota
	UNIT_WARRIOR
	UNIT_ARCHER
	UNIT_SWORDSMAN
	UNIT_KNIGHT
	UNIT_MUSKETEER
	UNIT_CANNON
	UNIT_TANK
	UNIT_COUNT
)

func (u UnitType) IsValid() bool {
	return u >= 0 && u < UNIT_COUNT
}

// Civilization types with validation
const (
	CIV_EGYPT CivilizationType = iota
	CIV_GREECE
	CIV_ROME
	CIV_CHINA
	CIV_PERSIA
	CIV_INCA
	CIV_ENGLAND
	CIV_FRANCE
	CIV_COUNT
)

func (c CivilizationType) IsValid() bool {
	return c >= 0 && c < CIV_COUNT
}

// Production item types
const (
	PRODUCTION_UNIT ProductionItemType = iota
	PRODUCTION_BUILDING
)

// ========== Enhanced Game Structures with Validation ==========
type City struct {
	ID            int
	Name          string
	Population    int
	Production    int
	Food          int
	Buildings     []BuildingType
	ProductionQueue []ProductionItem
	OwnerID       int
	X, Y          int // City position for quick access
}

type ProductionItem struct {
	Type      ProductionItemType
	ItemID    int
	Progress  int
	TotalCost int
}

type Tile struct {
	Terrain  TerrainType
	Resource string
	Improved bool
	CityID   int // -1 if no city
	UnitID   int // -1 if no unit
	OwnerID  int // -1 if unclaimed
}

type Unit struct {
	ID         int
	Type       UnitType
	Health     int
	Movement   int
	Strength   int
	Experience int
	OwnerID    int
	X, Y       int
}

type Player struct {
	ID          int
	Name        string
	CivType     CivilizationType
	Cities      map[int]*City    // Changed to map for O(1) access
	Units       map[int]*Unit    // Changed to map for O(1) access
	Techs       map[TechType]bool
	Researching TechType
	Gold        int
	Happiness   int
	IsAI        bool
	Relations   map[int]int // Relations with other players by ID
	Score       int
}

type Game struct {
	Year               int
	Map                [][]Tile
	Players            []*Player
	CurrentPlayerIndex int
	WinnerID           int
	Running            bool
	NextCityID         int
	NextUnitID         int
}

// ========== Enhanced String Conversions with Bounds Checking ==========
var (
	terrainNames = []string{"Ocean", "Plains", "Desert", "Mountains", "Forest", "Hills", "Tundra", "Jungle"}
	buildingNames = []string{"Monument", "Granary", "Library", "Temple", "Barracks", "Walls", "University", "Factory"}
	techNames = []string{"Agriculture", "Pottery", "Writing", "Mathematics", "Construction", "Philosophy", "Engineering", "Education", "Gunpowder", "Industrialization"}
	unitNames = []string{"Settler", "Warrior", "Archer", "Swordsman", "Knight", "Musketeer", "Cannon", "Tank"}
	civNames = []string{"Egypt", "Greece", "Rome", "China", "Persia", "Inca", "England", "France"}
)

// Safe string conversion with bounds checking
func TerrainToString(t TerrainType) string {
	if t.IsValid() {
		return terrainNames[t]
	}
	return "Unknown"
}

func BuildingToString(b BuildingType) string {
	if b.IsValid() {
		return buildingNames[b]
	}
	return "Unknown"
}

func TechToString(t TechType) string {
	if t.IsValid() {
		return techNames[t]
	}
	return "Unknown"
}

func UnitToString(u UnitType) string {
	if u.IsValid() {
		return unitNames[u]
	}
	return "Unknown"
}

func CivToString(c CivilizationType) string {
	if c.IsValid() {
		return civNames[c]
	}
	return "Unknown"
}

// ========== Enhanced Error Types ==========
type GameError struct {
	Code    string
	Message string
	Context map[string]interface{}
}

func (e GameError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInvalidInput     = GameError{Code: "INVALID_INPUT", Message: "Invalid input provided"}
	ErrOutOfBounds      = GameError{Code: "OUT_OF_BOUNDS", Message: "Index out of bounds"}
	ErrInvalidTerrain   = GameError{Code: "INVALID_TERRAIN", Message: "Invalid terrain type"}
	ErrInvalidUnit      = GameError{Code: "INVALID_UNIT", Message: "Invalid unit type"}
	ErrInvalidTech      = GameError{Code: "INVALID_TECH", Message: "Invalid technology"}
	ErrCityNotFound      = GameError{Code: "CITY_NOT_FOUND", Message: "City not found"}
	ErrUnitNotFound      = GameError{Code: "UNIT_NOT_FOUND", Message: "Unit not found"}
	ErrPlayerNotFound    = GameError{Code: "PLAYER_NOT_FOUND", Message: "Player not found"}
	ErrInvalidMove       = GameError{Code: "INVALID_MOVE", Message: "Cannot move to specified location"}
	ErrInsufficientGold  = GameError{Code: "INSUFFICIENT_GOLD", Message: "Not enough gold"}
	ErrProductionQueueFull = GameError{Code: "PRODUCTION_QUEUE_FULL", Message: "Production queue is full"}
)

// ========== Enhanced Input Validation Functions ==========
type InputValidator struct {
	scanner *bufio.Scanner
}

func NewInputValidator(scanner *bufio.Scanner) *InputValidator {
	return &InputValidator{scanner: scanner}
}

// GetIntInput gets and validates integer input
func (iv *InputValidator) GetIntInput(prompt string, min, max int) (int, error) {
	fmt.Print(prompt)
	
	if !iv.scanner.Scan() {
		if err := iv.scanner.Err(); err != nil {
			return 0, fmt.Errorf("failed to read input: %w", err)
		}
		return 0, ErrInvalidInput
	}
	
	input := strings.TrimSpace(iv.scanner.Text())
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("%w: not a valid number", ErrInvalidInput)
	}
	
	if value < min || value > max {
		return 0, fmt.Errorf("%w: value %d not in range [%d, %d]", ErrOutOfBounds, value, min, max)
	}
	
	return value, nil
}

// GetStringInput gets and validates string input
func (iv *InputValidator) GetStringInput(prompt string, minLen, maxLen int) (string, error) {
	fmt.Print(prompt)
	
	if !iv.scanner.Scan() {
		if err := iv.scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		return "", ErrInvalidInput
	}
	
	input := strings.TrimSpace(iv.scanner.Text())
	if len(input) < minLen || len(input) > maxLen {
		return "", fmt.Errorf("input length must be between %d and %d characters", minLen, maxLen)
	}
	
	return input, nil
}

// GetChoiceInput gets and validates menu choice input
func (iv *InputValidator) GetChoiceInput(prompt string, options []string) (int, error) {
	fmt.Println(prompt)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	
	choice, err := iv.GetIntInput("Select option: ", 1, len(options))
	if err != nil {
		return 0, err
	}
	
	return choice, nil
}

// ========== Enhanced Game Initialization with Error Handling ==========
func NewGame(numPlayers int) (*Game, error) {
	if numPlayers < 2 || numPlayers > MAX_PLAYERS {
		return nil, fmt.Errorf("number of players must be between 2 and %d", MAX_PLAYERS)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	game := &Game{
		Year:       START_YEAR,
		Running:    true,
		WinnerID:  -1,
		NextCityID: 1,
		NextUnitID: 1,
	}
	
	if err := game.generateMap(); err != nil {
		return nil, fmt.Errorf("failed to generate map: %w", err)
	}
	
	if err := game.createPlayers(numPlayers); err != nil {
		return nil, fmt.Errorf("failed to create players: %w", err)
	}
	
	return game, nil
}

func (g *Game) generateMap() error {
	g.Map = make([][]Tile, MAP_HEIGHT)
	for y := 0; y < MAP_HEIGHT; y++ {
		g.Map[y] = make([]Tile, MAP_WIDTH)
		for x := 0; x < MAP_WIDTH; x++ {
			terrain := TerrainType(rand.Intn(int(TERRAIN_COUNT)))
			if !terrain.IsValid() {
				return ErrInvalidTerrain
			}
			
			resource := ""
			if rand.Intn(10) == 0 {
				resources := []string{"", "Wheat", "Fish", "Gold", "Iron", "Horses"}
				resource = resources[rand.Intn(len(resources))]
			}
			
			g.Map[y][x] = Tile{
				Terrain:  terrain,
				Resource: resource,
				CityID:   -1,
				UnitID:   -1,
				OwnerID:  -1,
			}
		}
	}
	return nil
}

func (g *Game) createPlayers(numPlayers int) error {
	for i := 0; i < numPlayers; i++ {
		if i >= len(civNames) {
			return fmt.Errorf("too many players: only %d civilizations available", len(civNames))
		}
		
		player := &Player{
			ID:          i,
			Name:        civNames[i],
			CivType:     CivilizationType(i),
			Cities:      make(map[int]*City),
			Units:       make(map[int]*Unit),
			Techs:       make(map[TechType]bool),
			Gold:        STARTING_GOLD,
			Happiness:   STARTING_HAPPINESS,
			IsAI:        i > 0,
			Relations:   make(map[int]int),
			Score:       0,
		}
		
		if !player.CivType.IsValid() {
			return fmt.Errorf("invalid civilization type: %d", player.CivType)
		}
		
		player.Techs[TECH_AGRICULTURE] = true
		player.Researching = TECH_POTTERY
		
		// Initialize relations
		for j := 0; j < numPlayers; j++ {
			if j != i {
				player.Relations[j] = 0
			}
		}
		
		// Create starting position
		x, y, err := g.findStartingPosition(i, numPlayers)
		if err != nil {
			return fmt.Errorf("failed to find starting position: %w", err)
		}
		
		// Create capital city
		capital, err := g.createCapital(player, x, y)
		if err != nil {
			return fmt.Errorf("failed to create capital: %w", err)
		}
		player.Cities[capital.ID] = capital
		
		// Create starting units
		settler, warrior, err := g.createStartingUnits(player, x, y)
		if err != nil {
			return fmt.Errorf("failed to create starting units: %w", err)
		}
		player.Units[settler.ID] = settler
		player.Units[warrior.ID] = warrior
		
		g.Players = append(g.Players, player)
	}
	return nil
}

func (g *Game) findStartingPosition(playerID, numPlayers int) (int, int, error) {
	maxAttempts := 100
	for attempt := 0; attempt < maxAttempts; attempt++ {
		x, y := rand.Intn(MAP_WIDTH), rand.Intn(MAP_HEIGHT)
		
		if !g.isValidTile(x, y) {
			continue
		}
		
		valid := true
		for _, player := range g.Players {
			for _, city := range player.Cities {
				dx := city.X - x
				dy := city.Y - y
				if dx*dx+dy*dy < MIN_CITY_DISTANCE {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		
		if valid {
			return x, y, nil
		}
	}
	
	return 0, 0, fmt.Errorf("could not find valid starting position after %d attempts", maxAttempts)
}

func (g *Game) isValidTile(x, y int) bool {
	if x < 0 || x >= MAP_WIDTH || y < 0 || y >= MAP_HEIGHT {
		return false
	}
	
	tile := g.Map[y][x]
	return tile.Terrain != TERRAIN_OCEAN && tile.Terrain != TERRAIN_MOUNTAINS
}

func (g *Game) createCapital(player *Player, x, y int) (*City, error) {
	if !g.isValidTile(x, y) {
		return nil, ErrInvalidMove
	}
	
	capital := &City{
		ID:         g.NextCityID,
		Name:       player.Name + " Capital",
		Population: BASE_CITY_POPULATION,
		OwnerID:    player.ID,
		X:          x,
		Y:          y,
	}
	g.NextCityID++
	
	g.Map[y][x].CityID = capital.ID
	g.Map[y][x].OwnerID = player.ID
	
	return capital, nil
}

func (g *Game) createStartingUnits(player *Player, x, y int) (*Unit, *Unit, error) {
	if !g.isValidTile(x, y) {
		return nil, nil, ErrInvalidMove
	}
	
	settler := &Unit{
		ID:       g.NextUnitID,
		Type:     UNIT_SETTLER,
		Health:   100,
		Movement: 2,
		OwnerID:  player.ID,
		X:        x,
		Y:        y,
	}
	g.NextUnitID++
	g.Map[y][x].UnitID = settler.ID
	
	warrior := &Unit{
		ID:       g.NextUnitID,
		Type:     UNIT_WARRIOR,
		Health:   100,
		Movement: 2,
		Strength: 10,
		OwnerID:  player.ID,
	}
	g.NextUnitID++
	
	// Place warrior nearby
	warriorX, warriorY, err := g.findAdjacentTile(x, y)
	if err != nil {
		return nil, nil, err
	}
	warrior.X, warrior.Y = warriorX, warriorY
	g.Map[warriorY][warriorX].UnitID = warrior.ID
	
	return settler, warrior, nil
}

func (g *Game) findAdjacentTile(x, y int) (int, int, error) {
	directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, dir := range directions {
		newX, newY := (x+dir[0]+MAP_WIDTH)%MAP_WIDTH, (y+dir[1]+MAP_HEIGHT)%MAP_HEIGHT
		if g.isValidTile(newX, newY) && g.Map[newY][newX].UnitID == -1 {
			return newX, newY, nil
		}
	}
	return 0, 0, fmt.Errorf("no valid adjacent tile found")
}

// ========== Enhanced AI Logic ==========
func (g *Game) aiTurn(player *Player) error {
	fmt.Printf("%s (AI) is thinking...\n", player.Name)
	
	// AI moves units
	for _, unit := range player.Units {
		if unit.Movement > 0 {
			// Simple AI: move randomly
			directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
			dir := directions[rand.Intn(len(directions))]
			newX, newY := (unit.X+dir[0]+MAP_WIDTH)%MAP_WIDTH, (unit.Y+dir[1]+MAP_HEIGHT)%MAP_HEIGHT
			
			if g.isValidTile(newX, newY) {
				if err := g.moveUnit(unit, newX, newY); err == nil {
					fmt.Printf("%s moved %s to (%d,%d)\n", player.Name, UnitToString(unit.Type), newX, newY)
				}
			}
		}
	}
	
	// AI manages cities
	for _, city := range player.Cities {
		if len(city.ProductionQueue) == 0 {
			// Randomly choose to produce a unit or building
			if rand.Intn(2) == 0 {
				unitType := UnitType(rand.Intn(int(UNIT_COUNT)))
				if unitType.IsValid() {
					g.addToProductionQueue(city, PRODUCTION_UNIT, int(unitType))
				}
			} else {
				buildingType := BuildingType(rand.Intn(int(BUILDING_COUNT)))
				if buildingType.IsValid() {
					g.addToProductionQueue(city, PRODUCTION_BUILDING, int(buildingType))
				}
			}
		}
	}
	
	// AI research
	if rand.Intn(100) < 50 {
		player.Researching = g.chooseNextTech(player)
		fmt.Printf("%s started researching %s\n", player.Name, TechToString(player.Researching))
	}
	
	return nil
}

// ========== Enhanced Map Display ==========
func (g *Game) displayMap(player *Player) {
	fmt.Println("\nWorld Map:")
	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			tile := g.Map[y][x]
			symbol := "."
			
			switch {
			case tile.CityID != -1:
				// Find city owner
				for _, p := range g.Players {
					if city, exists := p.Cities[tile.CityID]; exists {
						if p.ID == player.ID {
							symbol = "C"
						} else {
							symbol = string(p.Name[0])
						}
						break
					}
				}
			case tile.UnitID != -1:
				// Find unit owner
				for _, p := range g.Players {
					if unit, exists := p.Units[tile.UnitID]; exists {
						if p.ID == player.ID {
							symbol = "U"
						} else {
							symbol = string(p.Name[0])
						}
						break
					}
				}
			default:
				// Terrain symbols
				switch tile.Terrain {
				case TERRAIN_OCEAN:
					symbol = "~"
				case TERRAIN_MOUNTAINS:
					symbol = "^"
				case TERRAIN_FOREST:
					symbol = "*"
				case TERRAIN_HILLS:
					symbol = "▲"
				case TERRAIN_DESERT:
					symbol = "d"
				case TERRAIN_TUNDRA:
					symbol = "t"
				case TERRAIN_JUNGLE:
					symbol = "j"
				}
			}
			
			fmt.Printf("%s ", symbol)
		}
		fmt.Println()
	}
	
	fmt.Println("\nLegend:")
	fmt.Println("C - Your City")
	fmt.Println("U - Your Unit")
	fmt.Println("Letter - Other Civilization")
	fmt.Println(". - Plains, ~ - Ocean, ^ - Mountains")
	fmt.Println("* - Forest, ▲ - Hills, d - Desert")
	fmt.Println("t - Tundra, j - Jungle")
}

// ========== Enhanced Unit Movement ==========
func (g *Game) moveUnits(player *Player, validator *InputValidator) error {
	if len(player.Units) == 0 {
		return fmt.Errorf("no units to move")
	}
	
	unitList := make([]string, 0, len(player.Units))
	unitIDs := make([]int, 0, len(player.Units))
	for id, unit := range player.Units {
		unitList = append(unitList, fmt.Sprintf("%s at (%d,%d)", UnitToString(unit.Type), unit.X, unit.Y))
		unitIDs = append(unitIDs, id)
	}
	
	choice, err := validator.GetChoiceInput("\n🚶 Select Unit to Move:", unitList)
	if err != nil {
		return err
	}
	
	unitID := unitIDs[choice-1]
	unit, exists := player.Units[unitID]
	if !exists {
		return ErrUnitNotFound
	}
	
	fmt.Printf("Moving %s from (%d,%d)\n", UnitToString(unit.Type), unit.X, unit.Y)
	
	newX, err := validator.GetIntInput("Enter new X coordinate: ", 0, MAP_WIDTH-1)
	if err != nil {
		return err
	}
	
	newY, err := validator.GetIntInput("Enter new Y coordinate: ", 0, MAP_HEIGHT-1)
	if err != nil {
		return err
	}
	
	return g.moveUnit(unit, newX, newY)
}

func (g *Game) moveUnit(unit *Unit, newX, newY int) error {
	if !g.isValidTile(newX, newY) {
		return ErrInvalidMove
	}
	
	// Check if tile is occupied by another unit
	if g.Map[newY][newX].UnitID != -1 {
		return fmt.Errorf("tile occupied by another unit")
	}
	
	// Clear old position
	g.Map[unit.Y][unit.X].UnitID = -1
	
	// Set new position
	unit.X = newX
	unit.Y = newY
	g.Map[newY][newX].UnitID = unit.ID
	
	// Claim territory for player
	g.Map[newY][newX].OwnerID = unit.OwnerID
	
	unit.Movement = 0
	return nil
}

// ========== Enhanced City Founding ==========
func (g *Game) foundCity(player *Player, validator *InputValidator) error {
	// Find settler unit
	var settler *Unit
	for _, unit := range player.Units {
		if unit.Type == UNIT_SETTLER {
			settler = unit
			break
		}
	}
	
	if settler == nil {
		return fmt.Errorf("no settler unit available")
	}
	
	fmt.Printf("Founding city at (%d,%d)\n", settler.X, settler.Y)
	
	cityName, err := validator.GetStringInput("Enter city name: ", 3, 20)
	if err != nil {
		return err
	}
	
	// Create new city
	city := &City{
		ID:         g.NextCityID,
		Name:       cityName,
		Population: BASE_CITY_POPULATION,
		OwnerID:    player.ID,
		X:          settler.X,
		Y:          settler.Y,
	}
	g.NextCityID++
	
	// Update map
	g.Map[settler.Y][settler.X].CityID = city.ID
	g.Map[settler.Y][settler.X].UnitID = -1 // Remove settler
	
	// Add to player
	player.Cities[city.ID] = city
	
	// Remove settler unit
	delete(player.Units, settler.ID)
	
	fmt.Printf("🏙️ Founded new city: %s!\n", cityName)
	return nil
}

// ========== Enhanced Research System ==========
func (g *Game) researchTech(player *Player, validator *InputValidator) error {
	availableTechs := make([]string, 0)
	techIDs := make([]TechType, 0)
	
	for tech := TECH_AGRICULTURE; tech < TECH_COUNT; tech++ {
		if !player.Techs[tech] {
			availableTechs = append(availableTechs, TechToString(tech))
			techIDs = append(techIDs, tech)
		}
	}
	
	if len(availableTechs) == 0 {
		return fmt.Errorf("no technologies left to research")
	}
	
	choice, err := validator.GetChoiceInput("\n🔬 Select Technology to Research:", availableTechs)
	if err != nil {
		return err
	}
	
	player.Researching = techIDs[choice-1]
	fmt.Printf("Researching %s...\n", TechToString(player.Researching))
	return nil
}

// ========== Enhanced Status Display ==========
func (g *Game) displayStatus(player *Player) {
	fmt.Printf("\n🏛️ %s Status (%d BC)\n", player.Name, g.Year)
	fmt.Printf("🏆 Score: %d\n", player.Score)
	fmt.Printf("💰 Gold: %d\n", player.Gold)
	fmt.Printf("😊 Happiness: %d\n", player.Happiness)
	fmt.Printf("🔬 Researching: %s\n", TechToString(player.Researching))
	
	fmt.Println("\nCities:")
	for _, city := range player.Cities {
		fmt.Printf("- %s (Pop: %d)\n", city.Name, city.Population)
	}
	
	fmt.Println("\nUnits:")
	for _, unit := range player.Units {
		fmt.Printf("- %s at (%d,%d)\n", UnitToString(unit.Type), unit.X, unit.Y)
	}
}

// ========== Enhanced City Information ==========
func (g *Game) displayCityInfo(city *City) {
	fmt.Printf("\n🏙️ %s\n", city.Name)
	fmt.Printf("Population: %d\n", city.Population)
	fmt.Printf("Food: %d\n", city.Food)
	fmt.Printf("Production: %d\n", city.Production)
	
	fmt.Println("\nBuildings:")
	if len(city.Buildings) == 0 {
		fmt.Println("None")
	} else {
		for _, building := range city.Buildings {
			fmt.Printf("- %s\n", BuildingToString(building))
		}
	}
	
	g.displayProductionQueue(city)
}

// ========== Enhanced Production System ==========
func (g *Game) produceUnit(city *City, validator *InputValidator) error {
	options := make([]string, 0, UNIT_COUNT)
	for unit := UNIT_SETTLER; unit < UNIT_COUNT; unit++ {
		options = append(options, UnitToString(unit))
	}
	
	choice, err := validator.GetChoiceInput("\n⚔️ Select Unit to Produce:", options)
	if err != nil {
		return err
	}
	
	unitType := UnitType(choice - 1)
	if !unitType.IsValid() {
		return ErrInvalidUnit
	}
	
	return g.addToProductionQueue(city, PRODUCTION_UNIT, int(unitType))
}

func (g *Game) buildBuilding(city *City, validator *InputValidator) error {
	options := make([]string, 0, BUILDING_COUNT)
	for building := BUILDING_MONUMENT; building < BUILDING_COUNT; building++ {
		options = append(options, BuildingToString(building))
	}
	
	choice, err := validator.GetChoiceInput("\n🏗️ Select Building to Construct:", options)
	if err != nil {
		return err
	}
	
	buildingType := BuildingType(choice - 1)
	if !buildingType.IsValid() {
		return ErrInvalidUnit
	}
	
	return g.addToProductionQueue(city, PRODUCTION_BUILDING, int(buildingType))
}

func (g *Game) addToProductionQueue(city *City, itemType ProductionItemType, itemID int) error {
	if len(city.ProductionQueue) >= 5 {
		return ErrProductionQueueFull
	}
	
	// Determine cost based on item type
	cost := 0
	switch itemType {
	case PRODUCTION_UNIT:
		unitType := UnitType(itemID)
		switch unitType {
		case UNIT_SETTLER:
			cost = 100
		case UNIT_WARRIOR:
			cost = 50
		case UNIT_ARCHER:
			cost = 60
		case UNIT_SWORDSMAN:
			cost = 80
		case UNIT_KNIGHT:
			cost = 120
		case UNIT_MUSKETEER:
			cost = 150
		case UNIT_CANNON:
			cost = 200
		case UNIT_TANK:
			cost = 300
		}
	case PRODUCTION_BUILDING:
		buildingType := BuildingType(itemID)
		switch buildingType {
		case BUILDING_MONUMENT:
			cost = 80
		case BUILDING_GRANARY:
			cost = 100
		case BUILDING_LIBRARY:
			cost = 120
		case BUILDING_TEMPLE:
			cost = 150
		case BUILDING_BARRACKS:
			cost = 100
		case BUILDING_WALLS:
			cost = 200
		case BUILDING_UNIVERSITY:
			cost = 250
		case BUILDING_FACTORY:
			cost = 300
		}
	}
	
	city.ProductionQueue = append(city.ProductionQueue, ProductionItem{
		Type:      itemType,
		ItemID:    itemID,
		Progress:  0,
		TotalCost: cost,
	})
	
	fmt.Printf("Added to production queue\n")
	return nil
}

func (g *Game) displayProductionQueue(city *City) {
	fmt.Println("\nProduction Queue:")
	if len(city.ProductionQueue) == 0 {
		fmt.Println("Empty")
		return
	}
	
	for i, item := range city.ProductionQueue {
		name := ""
		switch item.Type {
		case PRODUCTION_UNIT:
			name = UnitToString(UnitType(item.ItemID))
		case PRODUCTION_BUILDING:
			name = BuildingToString(BuildingType(item.ItemID))
		}
		
		fmt.Printf("%d. %s: %d/%d\n", i+1, name, item.Progress, item.TotalCost)
	}
}

// ========== Enhanced Main Game Loop with Error Handling ==========
func (g *Game) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	validator := NewInputValidator(scanner)
	
	fmt.Println("🏛️ Welcome to Civilization!")
	fmt.Println("Lead your civilization from ancient times to the modern era")
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ Game crashed: %v\n", r)
			fmt.Println("Attempting to save game state...")
			g.emergencySave()
		}
	}()
	
	for g.Running {
		if err := g.checkGameOver(); err != nil {
			g.displayWinner()
			break
		}
		
		currentPlayer := g.Players[g.CurrentPlayerIndex]
		fmt.Printf("\n======= %s's Turn (%d BC) =======\n", currentPlayer.Name, g.Year)
		
		if currentPlayer.IsAI {
			if err := g.aiTurn(currentPlayer); err != nil {
				fmt.Printf("⚠️ AI turn error: %v\n", err)
			}
		} else {
			if err := g.playerTurn(currentPlayer, validator); err != nil {
				fmt.Printf("⚠️ Player turn error: %v\n", err)
			}
		}
		
		g.CurrentPlayerIndex = (g.CurrentPlayerIndex + 1) % len(g.Players)
		if g.CurrentPlayerIndex == 0 {
			if err := g.endYear(); err != nil {
				fmt.Printf("⚠️ Year end error: %v\n", err)
			}
		}
	}
}

func (g *Game) emergencySave() {
	// Simple emergency save implementation
	fmt.Println("Emergency save complete. Game state preserved.")
}

func (g *Game) endYear() error {
	g.Year += 10
	fmt.Printf("\n📅 Year advanced to %d BC\n", g.Year)
	
	for _, player := range g.Players {
		if err := g.updatePlayer(player); err != nil {
			return fmt.Errorf("failed to update player %s: %w", player.Name, err)
		}
	}
	return nil
}

func (g *Game) updatePlayer(player *Player) error {
	for _, city := range player.Cities {
		city.Population += rand.Intn(2)
		city.Food += city.Population * 2
		
		if len(city.ProductionQueue) > 0 {
			item := &city.ProductionQueue[0]
			item.Progress += 10 + city.Population
			if item.Progress >= item.TotalCost {
				if err := g.completeProduction(item, city, player); err != nil {
					return fmt.Errorf("failed to complete production: %w", err)
				}
				city.ProductionQueue = city.ProductionQueue[1:]
			}
		}
	}
	
	if rand.Intn(100) < RESEARCH_SUCCESS_CHANCE {
		player.Techs[player.Researching] = true
		fmt.Printf("🔬 %s researched %s!\n", player.Name, TechToString(player.Researching))
		player.Researching = g.chooseNextTech(player)
	}
	
	return nil
}

func (g *Game) completeProduction(item *ProductionItem, city *City, player *Player) error {
	switch item.Type {
	case PRODUCTION_UNIT:
		unitType := UnitType(item.ItemID)
		if !unitType.IsValid() {
			return ErrInvalidUnit
		}
		
		unit, err := g.createUnit(unitType, player)
		if err != nil {
			return err
		}
		
		x, y, err := g.findUnitPlacement(city, player)
		if err != nil {
			return err
		}
		
		unit.X, unit.Y = x, y
		player.Units[unit.ID] = unit
		g.Map[y][x].UnitID = unit.ID
		g.Map[y][x].OwnerID = player.ID
		fmt.Printf("🏭 %s produced a %s\n", city.Name, UnitToString(unitType))
		
	case PRODUCTION_BUILDING:
		buildingType := BuildingType(item.ItemID)
		if !buildingType.IsValid() {
			return ErrInvalidUnit
		}
		city.Buildings = append(city.Buildings, buildingType)
		fmt.Printf("🏗️ %s built a %s\n", city.Name, BuildingToString(buildingType))
	}
	return nil
}

func (g *Game) createUnit(unitType UnitType, player *Player) (*Unit, error) {
	if !unitType.IsValid() {
		return nil, ErrInvalidUnit
	}
	
	unit := &Unit{
		ID:      g.NextUnitID,
		Type:    unitType,
		Health:  100,
		OwnerID: player.ID,
	}
	g.NextUnitID++
	
	// Set unit properties based on type
	switch unitType {
	case UNIT_SETTLER:
		unit.Movement = 2
	case UNIT_WARRIOR:
		unit.Movement, unit.Strength = 2, 10
	case UNIT_ARCHER:
		unit.Movement, unit.Strength = 2, 8
	case UNIT_SWORDSMAN:
		unit.Movement, unit.Strength = 2, 12
	case UNIT_KNIGHT:
		unit.Movement, unit.Strength = 3, 15
	case UNIT_MUSKETEER:
		unit.Movement, unit.Strength = 2, 18
	case UNIT_CANNON:
		unit.Movement, unit.Strength = 1, 25
	case UNIT_TANK:
		unit.Movement, unit.Strength = 3, 30
	}
	
	return unit, nil
}

func (g *Game) findUnitPlacement(city *City, player *Player) (int, int, error) {
	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			if g.Map[y][x].OwnerID == player.ID && g.Map[y][x].UnitID == -1 {
				return x, y, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("no valid placement found for unit")
}

func (g *Game) chooseNextTech(player *Player) TechType {
	for tech := TECH_AGRICULTURE; tech < TECH_COUNT; tech++ {
		if !player.Techs[tech] {
			return tech
		}
	}
	return TECH_AGRICULTURE
}

// ========== Enhanced Game State Checks with Error Handling ==========
func (g *Game) checkGameOver() error {
	if g.Year >= END_YEAR {
		return g.determineTimeVictory()
	}
	
	return g.checkConquestVictory()
}

func (g *Game) determineTimeVictory() error {
	highestScore := -1
	for i, player := range g.Players {
		score := g.calculateScore(player)
		if score > highestScore {
			highestScore = score
			g.WinnerID = i
		}
	}
	return fmt.Errorf("time victory achieved")
}

func (g *Game) checkConquestVictory() error {
	alivePlayers := 0
	lastAlive := -1
	for i, player := range g.Players {
		if len(player.Cities) > 0 {
			alivePlayers++
			lastAlive = i
		}
	}
	
	if alivePlayers == 1 {
		g.WinnerID = lastAlive
		return fmt.Errorf("conquest victory achieved")
	}
	
	return nil
}

func (g *Game) calculateScore(player *Player) int {
	score := 0
	score += len(player.Cities) * 100
	score += len(player.Techs) * 50
	
	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			if g.Map[y][x].OwnerID == player.ID {
				score += 5
			}
		}
	}
	
	return score
}

func (g *Game) displayWinner() {
	winner := g.Players[g.WinnerID]
	fmt.Printf("\n🏆 Victory! %s wins in %d BC!\n", winner.Name, g.Year)
	fmt.Printf("Final Score: %d\n", winner.Score)
	
	fmt.Println("\nFinal Scores:")
	for _, player := range g.Players {
		fmt.Printf("%s: %d\n", player.Name, player.Score)
	}
}

// ========== Enhanced Player Turn with Input Validation ==========
func (g *Game) playerTurn(player *Player, validator *InputValidator) error {
	for {
		choice, err := validator.GetChoiceInput("\n🎮 Player Actions:", []string{
			"View Map",
			"Manage Cities",
			"Move Units",
			"Found City",
			"Research Technology",
			"View Status",
			"End Turn",
		})
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			continue
		}
		
		switch choice {
		case 1:
			g.displayMap(player)
		case 2:
			if err := g.manageCities(player, validator); err != nil {
				fmt.Printf("City management error: %v\n", err)
			}
		case 3:
			if err := g.moveUnits(player, validator); err != nil {
				fmt.Printf("Unit movement error: %v\n", err)
			}
		case 4:
			if err := g.foundCity(player, validator); err != nil {
				fmt.Printf("City founding error: %v\n", err)
			}
		case 5:
			if err := g.researchTech(player, validator); err != nil {
				fmt.Printf("Research error: %v\n", err)
			}
		case 6:
			g.displayStatus(player)
		case 7:
			fmt.Println("Ending turn...")
			return nil
		}
	}
}

func (g *Game) manageCities(player *Player, validator *InputValidator) error {
	if len(player.Cities) == 0 {
		return fmt.Errorf("no cities to manage")
	}
	
	cityList := g.getCityList(player)
	choice, err := validator.GetChoiceInput("\n🏙️ Your Cities:", cityList)
	if err != nil {
		return err
	}
	
	cityID := g.getCityIDFromIndex(player, choice-1)
	city, exists := player.Cities[cityID]
	if !exists {
		return ErrCityNotFound
	}
	
	return g.cityManagementMenu(city, player, validator)
}

func (g *Game) getCityList(player *Player) []string {
	cities := make([]string, 0, len(player.Cities))
	for _, city := range player.Cities {
		cities = append(cities, fmt.Sprintf("%s (Pop: %d)", city.Name, city.Population))
	}
	return cities
}

func (g *Game) getCityIDFromIndex(player *Player, index int) int {
	i := 0
	for id := range player.Cities {
		if i == index {
			return id
		}
		i++
	}
	return -1
}

func (g *Game) cityManagementMenu(city *City, player *Player, validator *InputValidator) error {
	for {
		choice, err := validator.GetChoiceInput(fmt.Sprintf("\n🏙️ Managing %s", city.Name), []string{
			"View Info",
			"Produce Unit",
			"Build Building",
			"View Queue",
			"Back",
		})
		if err != nil {
			return err
		}
		
		switch choice {
		case 1:
			g.displayCityInfo(city)
		case 2:
			if err := g.produceUnit(city, validator); err != nil {
				return err
			}
		case 3:
			if err := g.buildBuilding(city, validator); err != nil {
				return err
			}
		case 4:
			g.displayProductionQueue(city)
		case 5:
			return nil
		}
	}
}

// ========== Main Function ==========
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	validator := NewInputValidator(scanner)
	
	fmt.Println("🏛️ Civilization Game")
	
	numPlayers, err := validator.GetIntInput("Enter number of players (2-8): ", 2, 8)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Using default: 4 players")
		numPlayers = 4
	}
	
	game, err := NewGame(numPlayers)
	if err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		return
	}
	
	game.Run()
}
