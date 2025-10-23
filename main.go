//remake
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

// ========== Constants ==========
const (
	mapWidth         = 20
	mapHeight        = 15
	maxPlayers       = 8
	maxCities        = 50
	maxUnits         = 100
	startYear        = 4000
	endYear          = 2050
	minCityDistance  = 25
	maxProductionQueue = 5
	
	// Game balance constants
	researchSuccessChance = 30
	combatSuccessChance   = 70
	baseCityPopulation   = 1
	startingGold          = 100
	startingHappiness     = 100
)

// ========== Type Definitions ==========
type (
	terrainType       int
	buildingType      int
	techType          int
	unitType          int
	civilizationType  int
	productionItemType int
)

// Terrain types
const (
	terrainOcean terrainType = iota
	terrainPlains
	terrainDesert
	terrainMountains
	terrainForest
	terrainHills
	terrainTundra
	terrainJungle
	terrainCount
)

func (t terrainType) isValid() bool {
	return t >= 0 && t < terrainCount
}

// Building types
const (
	buildingMonument buildingType = iota
	buildingGranary
	buildingLibrary
	buildingTemple
	buildingBarracks
	buildingWalls
	buildingUniversity
	buildingFactory
	buildingCount
)

func (b buildingType) isValid() bool {
	return b >= 0 && b < buildingCount
}

// Technology types
const (
	techAgriculture techType = iota
	techPottery
	techWriting
	techMathematics
	techConstruction
	techPhilosophy
	techEngineering
	techEducation
	techGunpowder
	techIndustrialization
	techCount
)

func (t techType) isValid() bool {
	return t >= 0 && t < techCount
}

// Unit types
const (
	unitSettler unitType = iota
	unitWarrior
	unitArcher
	unitSwordsman
	unitKnight
	unitMusketeer
	unitCannon
	unitTank
	unitCount
)

func (u unitType) isValid() bool {
	return u >= 0 && u < unitCount
}

// Civilization types
const (
	civEgypt civilizationType = iota
	civGreece
	civRome
	civChina
	civPersia
	civInca
	civEngland
	civFrance
	civCount
)

func (c civilizationType) isValid() bool {
	return c >= 0 && c < civCount
}

// Production item types
const (
	productionUnit productionItemType = iota
	productionBuilding
)

// ========== Game Structures ==========
type city struct {
	ID            int
	Name          string
	Population    int
	Production    int
	Food          int
	Buildings     []buildingType
	ProductionQueue []productionItem
	OwnerID       int
	X, Y          int
}

type productionItem struct {
	Type      productionItemType
	ItemID    int
	Progress  int
	TotalCost int
	Name      string
}

type tile struct {
	Terrain  terrainType
	Resource string
	Improved bool
	CityID   int
	UnitID   int
	OwnerID  int
}

type unit struct {
	ID         int
	Type       unitType
	Health     int
	Movement   int
	Strength   int
	Experience int
	OwnerID    int
	X, Y       int
}

type player struct {
	ID          int
	Name        string
	CivType     civilizationType
	Cities      map[int]*city
	Units       map[int]*unit
	Techs       map[techType]bool
	Researching techType
	Gold        int
	Happiness   int
	IsAI        bool
	Relations   map[int]int
	Score       int
	CityCount   int
	UnitCount   int
}

type game struct {
	Year               int
	Map                [][]tile
	Players            []*player
	CurrentPlayerIndex int
	WinnerID           int
	Running            bool
	NextCityID         int
	NextUnitID         int
	TurnCount          int
}

// ========== String Conversions ==========
var (
	terrainNames = [terrainCount]string{"Ocean", "Plains", "Desert", "Mountains", "Forest", "Hills", "Tundra", "Jungle"}
	buildingNames = [buildingCount]string{"Monument", "Granary", "Library", "Temple", "Barracks", "Walls", "University", "Factory"}
	techNames = [techCount]string{"Agriculture", "Pottery", "Writing", "Mathematics", "Construction", "Philosophy", "Engineering", "Education", "Gunpowder", "Industrialization"}
	unitNames = [unitCount]string{"Settler", "Warrior", "Archer", "Swordsman", "Knight", "Musketeer", "Cannon", "Tank"}
	civNames = [civCount]string{"Egypt", "Greece", "Rome", "China", "Persia", "Inca", "England", "France"}
)

func terrainToString(t terrainType) string {
	if t.isValid() {
		return terrainNames[t]
	}
	return "Unknown"
}

func buildingToString(b buildingType) string {
	if b.isValid() {
		return buildingNames[b]
	}
	return "Unknown"
}

func techToString(t techType) string {
	if t.isValid() {
		return techNames[t]
	}
	return "Unknown"
}

func unitToString(u unitType) string {
	if u.isValid() {
		return unitNames[u]
	}
	return "Unknown"
}

func civToString(c civilizationType) string {
	if c.isValid() {
		return civNames[c]
	}
	return "Unknown"
}

// ========== Error Handling ==========
type gameError struct {
	Code    string
	Message string
}

func (e gameError) Error() string {
	return e.Code + ": " + e.Message
}

var (
	errInvalidInput     = gameError{Code: "INVALID_INPUT", Message: "invalid input provided"}
	errOutOfBounds      = gameError{Code: "OUT_OF_BOUNDS", Message: "index out of bounds"}
	errInvalidTerrain   = gameError{Code: "INVALID_TERRAIN", Message: "invalid terrain type"}
	errInvalidUnit      = gameError{Code: "INVALID_UNIT", Message: "invalid unit type"}
	errInvalidTech      = gameError{Code: "INVALID_TECH", Message: "invalid technology"}
	errCityNotFound     = gameError{Code: "CITY_NOT_FOUND", Message: "city not found"}
	errUnitNotFound     = gameError{Code: "UNIT_NOT_FOUND", Message: "unit not found"}
	errInvalidMove      = gameError{Code: "INVALID_MOVE", Message: "cannot move to specified location"}
	errProductionQueueFull = gameError{Code: "PRODUCTION_QUEUE_FULL", Message: "production queue is full"}
)

// ========== Input Validation ==========
type inputValidator struct {
	scanner *bufio.Scanner
}

func newInputValidator(scanner *bufio.Scanner) *inputValidator {
	return &inputValidator{
		scanner: scanner,
	}
}

func (iv *inputValidator) getIntInput(prompt string, min, max int) (int, error) {
	fmt.Print(prompt)
	
	if !iv.scanner.Scan() {
		if err := iv.scanner.Err(); err != nil {
			return 0, fmt.Errorf("failed to read input: %w", err)
		}
		return 0, errInvalidInput
	}
	
	input := strings.TrimSpace(iv.scanner.Text())
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("%w: not a valid number", errInvalidInput)
	}
	
	if value < min || value > max {
		return 0, fmt.Errorf("%w: value %d not in range [%d, %d]", errOutOfBounds, value, min, max)
	}
	
	return value, nil
}

func (iv *inputValidator) getStringInput(prompt string, minLen, maxLen int) (string, error) {
	fmt.Print(prompt)
	
	if !iv.scanner.Scan() {
		if err := iv.scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		return "", errInvalidInput
	}
	
	input := strings.TrimSpace(iv.scanner.Text())
	if len(input) < minLen || len(input) > maxLen {
		return "", fmt.Errorf("input length must be between %d and %d characters", minLen, maxLen)
	}
	
	return input, nil
}

func (iv *inputValidator) getChoiceInput(prompt string, options []string) (int, error) {
	fmt.Println(prompt)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	
	choice, err := iv.getIntInput("Select option: ", 1, len(options))
	if err != nil {
		return 0, err
	}
	
	return choice, nil
}

// ========== Game Initialization ==========
func newGame(numPlayers int) (*game, error) {
	if numPlayers < 2 || numPlayers > maxPlayers {
		return nil, fmt.Errorf("number of players must be between 2 and %d", maxPlayers)
	}
	
	rand.Seed(time.Now().UnixNano())
	
	game := &game{
		Year:       startYear,
		Running:    true,
		WinnerID:  -1,
		NextCityID: 1,
		NextUnitID: 1,
		TurnCount:  0,
	}
	
	if err := game.generateMap(); err != nil {
		return nil, fmt.Errorf("failed to generate map: %w", err)
	}
	
	if err := game.createPlayers(numPlayers); err != nil {
		return nil, fmt.Errorf("failed to create players: %w", err)
	}
	
	return game, nil
}

func (g *game) generateMap() error {
	g.Map = make([][]tile, mapHeight)
	for y := 0; y < mapHeight; y++ {
		g.Map[y] = make([]tile, mapWidth)
		for x := 0; x < mapWidth; x++ {
			terrain := terrainType(rand.Intn(int(terrainCount)))
			if !terrain.isValid() {
				return errInvalidTerrain
			}
			
			resource := ""
			if rand.Intn(10) == 0 {
				resources := []string{"Wheat", "Fish", "Gold", "Iron", "Horses"}
				resource = resources[rand.Intn(len(resources))]
			}
			
			g.Map[y][x] = tile{
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

func (g *game) createPlayers(numPlayers int) error {
	for i := 0; i < numPlayers; i++ {
		if i >= len(civNames) {
			return fmt.Errorf("too many players: only %d civilizations available", len(civNames))
		}
		
		player := &player{
			ID:          i,
			Name:        civNames[i],
			CivType:     civilizationType(i),
			Cities:      make(map[int]*city, maxCities),
			Units:       make(map[int]*unit, maxUnits),
			Techs:       make(map[techType]bool),
			Gold:        startingGold,
			Happiness:   startingHappiness,
			IsAI:        i > 0,
			Relations:   make(map[int]int),
			Score:       0,
			CityCount:   0,
			UnitCount:   0,
		}
		
		if !player.CivType.isValid() {
			return fmt.Errorf("invalid civilization type: %d", player.CivType)
		}
		
		player.Techs[techAgriculture] = true
		player.Researching = techPottery
		
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
		player.CityCount++
		
		// Create starting units
		settler, warrior, err := g.createStartingUnits(player, x, y)
		if err != nil {
			return fmt.Errorf("failed to create starting units: %w", err)
		}
		player.Units[settler.ID] = settler
		player.Units[warrior.ID] = warrior
		player.UnitCount += 2
		
		g.Players = append(g.Players, player)
	}
	return nil
}

func (g *game) findStartingPosition(playerID, numPlayers int) (int, int, error) {
	maxAttempts := 100
	directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	
	for attempt := 0; attempt < maxAttempts; attempt++ {
		x, y := rand.Intn(mapWidth), rand.Intn(mapHeight)
		
		if !g.isValidTile(x, y) {
			continue
		}
		
		valid := true
		// Check distance from existing cities
		for _, player := range g.Players {
			for _, city := range player.Cities {
				dx := city.X - x
				dy := city.Y - y
				if dx*dx+dy*dy < minCityDistance {
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
	
	// Fallback: find any valid position
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.isValidTile(x, y) {
				return x, y, nil
			}
		}
	}
	
	return 0, 0, fmt.Errorf("could not find valid starting position")
}

func (g *game) isValidTile(x, y int) bool {
	return x >= 0 && x < mapWidth && y >= 0 && y < mapHeight &&
		g.Map[y][x].Terrain != terrainOcean && 
		g.Map[y][x].Terrain != terrainMountains
}

func (g *game) createCapital(player *player, x, y int) (*city, error) {
	if !g.isValidTile(x, y) {
		return nil, errInvalidMove
	}
	
	capital := &city{
		ID:         g.NextCityID,
		Name:       player.Name + " Capital",
		Population: baseCityPopulation,
		OwnerID:    player.ID,
		X:          x,
		Y:          y,
	}
	g.NextCityID++
	
	g.Map[y][x].CityID = capital.ID
	g.Map[y][x].OwnerID = player.ID
	
	return capital, nil
}

func (g *game) createStartingUnits(player *player, x, y int) (*unit, *unit, error) {
	if !g.isValidTile(x, y) {
		return nil, nil, errInvalidMove
	}
	
	settler := &unit{
		ID:       g.NextUnitID,
		Type:     unitSettler,
		Health:   100,
		Movement: 2,
		Strength: 5,
		OwnerID:  player.ID,
		X:        x,
		Y:        y,
	}
	g.NextUnitID++
	g.Map[y][x].UnitID = settler.ID
	
	warrior := &unit{
		ID:       g.NextUnitID,
		Type:     unitWarrior,
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

func (g *game) findAdjacentTile(x, y int) (int, int, error) {
	directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	// Shuffle directions for better distribution
	for i := range directions {
		j := rand.Intn(i + 1)
		directions[i], directions[j] = directions[j], directions[i]
	}
	
	for _, dir := range directions {
		newX, newY := (x+dir[0]+mapWidth)%mapWidth, (y+dir[1]+mapHeight)%mapHeight
		if g.isValidTile(newX, newY) && g.Map[newY][newX].UnitID == -1 {
			return newX, newY, nil
		}
	}
	return 0, 0, fmt.Errorf("no valid adjacent tile found")
}

// ========== AI Logic ==========
func (g *game) aiTurn(player *player) error {
	fmt.Printf("%s (AI) is thinking...\n", player.Name)
	
	// AI moves units
	for _, unit := range player.Units {
		if unit.Movement > 0 {
			directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
			dir := directions[rand.Intn(len(directions))]
			newX, newY := (unit.X+dir[0]+mapWidth)%mapWidth, (unit.Y+dir[1]+mapHeight)%mapHeight
			
			if g.isValidTile(newX, newY) && g.Map[newY][newX].UnitID == -1 {
				if err := g.moveUnit(unit, newX, newY); err == nil {
					fmt.Printf("%s moved %s to (%d,%d)\n", player.Name, unitToString(unit.Type), newX, newY)
				}
			}
		}
	}
	
	// AI manages cities
	for _, city := range player.Cities {
		if len(city.ProductionQueue) == 0 {
			if rand.Intn(2) == 0 {
				unitType := unitType(rand.Intn(int(unitCount)))
				if unitType.isValid() {
					g.addToProductionQueue(city, productionUnit, int(unitType))
				}
			} else {
				buildingType := buildingType(rand.Intn(int(buildingCount)))
				if buildingType.isValid() {
					g.addToProductionQueue(city, productionBuilding, int(buildingType))
				}
			}
		}
	}
	
	// AI research
	if rand.Intn(100) < 50 {
		player.Researching = g.chooseNextTech(player)
		fmt.Printf("%s started researching %s\n", player.Name, techToString(player.Researching))
	}
	
	return nil
}

// ========== Map Display ==========
func (g *game) displayMap(player *player) {
	fmt.Println("\nWorld Map:")
	terrainSymbols := [terrainCount]string{"~", ".", "d", "^", "*", "▲", "t", "j"}
	
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			tile := g.Map[y][x]
			symbol := terrainSymbols[tile.Terrain]
			
			if tile.CityID != -1 {
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
			} else if tile.UnitID != -1 {
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
			}
			
			fmt.Printf("%s ", symbol)
		}
		fmt.Println()
	}
	
	fmt.Println("\nLegend:")
	fmt.Println("C - Your City, U - Your Unit")
	fmt.Println("Letter - Other Civilization")
	fmt.Println(". - Plains, ~ - Ocean, ^ - Mountains")
	fmt.Println("* - Forest, ▲ - Hills, d - Desert")
	fmt.Println("t - Tundra, j - Jungle")
}

// ========== Unit Movement ==========
func (g *game) moveUnits(player *player, validator *inputValidator) error {
	if player.UnitCount == 0 {
		return fmt.Errorf("no units to move")
	}
	
	unitList := make([]string, 0, player.UnitCount)
	unitIDs := make([]int, 0, player.UnitCount)
	for id, unit := range player.Units {
		unitList = append(unitList, fmt.Sprintf("%s at (%d,%d)", unitToString(unit.Type), unit.X, unit.Y))
		unitIDs = append(unitIDs, id)
	}
	
	choice, err := validator.getChoiceInput("\n🚶 Select Unit to Move:", unitList)
	if err != nil {
		return err
	}
	
	unitID := unitIDs[choice-1]
	unit, exists := player.Units[unitID]
	if !exists {
		return errUnitNotFound
	}
	
	fmt.Printf("Moving %s from (%d,%d)\n", unitToString(unit.Type), unit.X, unit.Y)
	
	newX, err := validator.getIntInput("Enter new X coordinate: ", 0, mapWidth-1)
	if err != nil {
		return err
	}
	
	newY, err := validator.getIntInput("Enter new Y coordinate: ", 0, mapHeight-1)
	if err != nil {
		return err
	}
	
	return g.moveUnit(unit, newX, newY)
}

func (g *game) moveUnit(unit *unit, newX, newY int) error {
	if !g.isValidTile(newX, newY) {
		return errInvalidMove
	}
	
	if g.Map[newY][newX].UnitID != -1 {
		return fmt.Errorf("tile occupied by another unit")
	}
	
	// Clear old position
	g.Map[unit.Y][unit.X].UnitID = -1
	
	// Set new position
	unit.X, unit.Y = newX, newY
	g.Map[newY][newX].UnitID = unit.ID
	g.Map[newY][newX].OwnerID = unit.OwnerID
	
	unit.Movement = 0
	return nil
}

// ========== City Founding ==========
func (g *game) foundCity(player *player, validator *inputValidator) error {
	var settler *unit
	for _, unit := range player.Units {
		if unit.Type == unitSettler {
			settler = unit
			break
		}
	}
	
	if settler == nil {
		return fmt.Errorf("no settler unit available")
	}
	
	fmt.Printf("Founding city at (%d,%d)\n", settler.X, settler.Y)
	
	cityName, err := validator.getStringInput("Enter city name: ", 3, 20)
	if err != nil {
		return err
	}
	
	city := &city{
		ID:         g.NextCityID,
		Name:       cityName,
		Population: baseCityPopulation,
		OwnerID:    player.ID,
		X:          settler.X,
		Y:          settler.Y,
	}
	g.NextCityID++
	
	g.Map[settler.Y][settler.X].CityID = city.ID
	g.Map[settler.Y][settler.X].UnitID = -1
	
	player.Cities[city.ID] = city
	player.CityCount++
	delete(player.Units, settler.ID)
	player.UnitCount--
	
	fmt.Printf("🏙️ Founded new city: %s!\n", cityName)
	return nil
}

// ========== Research System ==========
func (g *game) researchTech(player *player, validator *inputValidator) error {
	availableTechs := make([]string, 0, techCount)
	techIDs := make([]techType, 0, techCount)
	
	for tech := techAgriculture; tech < techCount; tech++ {
		if !player.Techs[tech] {
			availableTechs = append(availableTechs, techToString(tech))
			techIDs = append(techIDs, tech)
		}
	}
	
	if len(availableTechs) == 0 {
		return fmt.Errorf("no technologies left to research")
	}
	
	choice, err := validator.getChoiceInput("\n🔬 Select Technology to Research:", availableTechs)
	if err != nil {
		return err
	}
	
	player.Researching = techIDs[choice-1]
	fmt.Printf("Researching %s...\n", techToString(player.Researching))
	return nil
}

// ========== Status Display ==========
func (g *game) displayStatus(player *player) {
	fmt.Printf("\n🏛️ %s Status (%d BC)\n", player.Name, g.Year)
	fmt.Printf("🏆 Score: %d\n", player.Score)
	fmt.Printf("💰 Gold: %d\n", player.Gold)
	fmt.Printf("😊 Happiness: %d\n", player.Happiness)
	fmt.Printf("🔬 Researching: %s\n", techToString(player.Researching))
	
	fmt.Printf("\nCities (%d):\n", player.CityCount)
	for _, city := range player.Cities {
		fmt.Printf("- %s (Pop: %d)\n", city.Name, city.Population)
	}
	
	fmt.Printf("\nUnits (%d):\n", player.UnitCount)
	for _, unit := range player.Units {
		fmt.Printf("- %s at (%d,%d)\n", unitToString(unit.Type), unit.X, unit.Y)
	}
}

// ========== Production System ==========
func (g *game) produceUnit(city *city, validator *inputValidator) error {
	options := make([]string, unitCount)
	for i := 0; i < int(unitCount); i++ {
		options[i] = unitToString(unitType(i))
	}
	
	choice, err := validator.getChoiceInput("\n⚔️ Select Unit to Produce:", options)
	if err != nil {
		return err
	}
	
	unitType := unitType(choice - 1)
	if !unitType.isValid() {
		return errInvalidUnit
	}
	
	return g.addToProductionQueue(city, productionUnit, int(unitType))
}

func (g *game) buildBuilding(city *city, validator *inputValidator) error {
	options := make([]string, buildingCount)
	for i := 0; i < int(buildingCount); i++ {
		options[i] = buildingToString(buildingType(i))
	}
	
	choice, err := validator.getChoiceInput("\n🏗️ Select Building to Construct:", options)
	if err != nil {
		return err
	}
	
	buildingType := buildingType(choice - 1)
	if !buildingType.isValid() {
		return errInvalidUnit
	}
	
	return g.addToProductionQueue(city, productionBuilding, int(buildingType))
}

func (g *game) addToProductionQueue(city *city, itemType productionItemType, itemID int) error {
	if len(city.ProductionQueue) >= maxProductionQueue {
		return errProductionQueueFull
	}
	
	var cost int
	var name string
	
	switch itemType {
	case productionUnit:
		unitType := unitType(itemID)
		cost = g.getUnitCost(unitType)
		name = unitToString(unitType)
	case productionBuilding:
		buildingType := buildingType(itemID)
		cost = g.getBuildingCost(buildingType)
		name = buildingToString(buildingType)
	}
	
	city.ProductionQueue = append(city.ProductionQueue, productionItem{
		Type:      itemType,
		ItemID:    itemID,
		Progress:  0,
		TotalCost: cost,
		Name:      name,
	})
	
	fmt.Printf("Added %s to production queue (Cost: %d)\n", name, cost)
	return nil
}

func (g *game) getUnitCost(unitType unitType) int {
	costs := map[unitType]int{
		unitSettler:  100,
		unitWarrior:   50,
		unitArcher:    60,
		unitSwordsman: 80,
		unitKnight:   120,
		unitMusketeer:150,
		unitCannon:   200,
		unitTank:     300,
	}
	return costs[unitType]
}

func (g *game) getBuildingCost(buildingType buildingType) int {
	costs := map[buildingType]int{
		buildingMonument:  80,
		buildingGranary:  100,
		buildingLibrary:  120,
		buildingTemple:   150,
		buildingBarracks:100,
		buildingWalls:   200,
		buildingUniversity:250,
		buildingFactory: 300,
	}
	return costs[buildingType]
}

// ========== Main Game Loop ==========
func (g *game) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	validator := newInputValidator(scanner)
	
	fmt.Println("🏛️ Welcome to Civilization!")
	fmt.Println("Lead your civilization from ancient times to the modern era")
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ Game crashed: %v\n", r)
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

func (g *game) emergencySave() {
	fmt.Println("Emergency save complete. Game state preserved.")
}

func (g *game) endYear() error {
	g.Year += 10
	g.TurnCount++
	fmt.Printf("\n📅 Year advanced to %d BC\n", g.Year)
	
	for _, player := range g.Players {
		if err := g.updatePlayer(player); err != nil {
			return fmt.Errorf("failed to update player %s: %w", player.Name, err)
		}
	}
	return nil
}

func (g *game) updatePlayer(player *player) error {
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
	
	if rand.Intn(100) < researchSuccessChance {
		player.Techs[player.Researching] = true
		fmt.Printf("🔬 %s researched %s!\n", player.Name, techToString(player.Researching))
		player.Researching = g.chooseNextTech(player)
	}
	
	return nil
}

func (g *game) completeProduction(item *productionItem, city *city, player *player) error {
	switch item.Type {
	case productionUnit:
		unitType := unitType(item.ItemID)
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
		player.UnitCount++
		g.Map[y][x].UnitID = unit.ID
		g.Map[y][x].OwnerID = player.ID
		fmt.Printf("🏭 %s produced a %s\n", city.Name, item.Name)
		
	case productionBuilding:
		buildingType := buildingType(item.ItemID)
		city.Buildings = append(city.Buildings, buildingType)
		fmt.Printf("🏗️ %s built a %s\n", city.Name, item.Name)
	}
	return nil
}

func (g *game) createUnit(unitType unitType, player *player) (*unit, error) {
	if !unitType.isValid() {
		return nil, errInvalidUnit
	}
	
	unit := &unit{
		ID:      g.NextUnitID,
		Type:    unitType,
		Health:  100,
		OwnerID: player.ID,
	}
	g.NextUnitID++
	
	// Set unit properties
	switch unitType {
	case unitSettler:
		unit.Movement, unit.Strength = 2, 5
	case unitWarrior:
		unit.Movement, unit.Strength = 2, 10
	case unitArcher:
		unit.Movement, unit.Strength = 2, 8
	case unitSwordsman:
		unit.Movement, unit.Strength = 2, 12
	case unitKnight:
		unit.Movement, unit.Strength = 3, 15
	case unitMusketeer:
		unit.Movement, unit.Strength = 2, 18
	case unitCannon:
		unit.Movement, unit.Strength = 1, 25
	case unitTank:
		unit.Movement, unit.Strength = 3, 30
	}
	
	return unit, nil
}

func (g *game) findUnitPlacement(city *city, player *player) (int, int, error) {
	// Check adjacent tiles first
	directions := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, dir := range directions {
		x, y := (city.X+dir[0]+mapWidth)%mapWidth, (city.Y+dir[1]+mapHeight)%mapHeight
		if g.Map[y][x].OwnerID == player.ID && g.Map[y][x].UnitID == -1 {
			return x, y, nil
		}
	}
	
	// Fallback: any player-owned tile
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.Map[y][x].OwnerID == player.ID && g.Map[y][x].UnitID == -1 {
				return x, y, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("no valid placement found for unit")
}

func (g *game) chooseNextTech(player *player) techType {
	for tech := techAgriculture; tech < techCount; tech++ {
		if !player.Techs[tech] {
			return tech
		}
	}
	return techAgriculture
}

// ========== Game State Checks ==========
func (g *game) checkGameOver() error {
	if g.Year >= endYear {
		return g.determineTimeVictory()
	}
	return g.checkConquestVictory()
}

func (g *game) determineTimeVictory() error {
	highestScore := -1
	for i, player := range g.Players {
		score := g.calculateScore(player)
		player.Score = score
		if score > highestScore {
			highestScore = score
			g.WinnerID = i
		}
	}
	return fmt.Errorf("time victory achieved")
}

func (g *game) checkConquestVictory() error {
	alivePlayers := 0
	lastAlive := -1
	for i, player := range g.Players {
		if player.CityCount > 0 {
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

func (g *game) calculateScore(player *player) int {
	score := player.CityCount * 100
	score += len(player.Techs) * 50
	score += player.UnitCount * 10
	
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			if g.Map[y][x].OwnerID == player.ID {
				score += 5
			}
		}
	}
	
	return score
}

func (g *game) displayWinner() {
	winner := g.Players[g.WinnerID]
	fmt.Printf("\n🏆 Victory! %s wins in %d BC!\n", winner.Name, g.Year)
	fmt.Printf("Final Score: %d\n", winner.Score)
	
	fmt.Println("\nFinal Scores:")
	for _, player := range g.Players {
		fmt.Printf("%s: %d\n", player.Name, player.Score)
	}
}

// ========== Player Turn ==========
func (g *game) playerTurn(player *player, validator *inputValidator) error {
	for {
		choice, err := validator.getChoiceInput("\n🎮 Player Actions:", []string{
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

func (g *game) manageCities(player *player, validator *inputValidator) error {
	if player.CityCount == 0 {
		return fmt.Errorf("no cities to manage")
	}
	
	cityList := make([]string, 0, player.CityCount)
	for _, city := range player.Cities {
		cityList = append(cityList, fmt.Sprintf("%s (Pop: %d)", city.Name, city.Population))
	}
	
	choice, err := validator.getChoiceInput("\n🏙️ Your Cities:", cityList)
	if err != nil {
		return err
	}
	
	// Get city by index
	var selectedCity *city
	i := 0
	for _, c := range player.Cities {
		if i == choice-1 {
			selectedCity = c
			break
		}
		i++
	}
	
	if selectedCity == nil {
		return errCityNotFound
	}
	
	return g.cityManagementMenu(selectedCity, player, validator)
}

func (g *game) cityManagementMenu(city *city, player *player, validator *inputValidator) error {
	for {
		choice, err := validator.getChoiceInput(fmt.Sprintf("\n🏙️ Managing %s", city.Name), []string{
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

func (g *game) displayCityInfo(city *city) {
	fmt.Printf("\n🏙️ %s\n", city.Name)
	fmt.Printf("Population: %d\n", city.Population)
	fmt.Printf("Food: %d\n", city.Food)
	fmt.Printf("Production: %d\n", city.Production)
	
	fmt.Println("\nBuildings:")
	if len(city.Buildings) == 0 {
		fmt.Println("None")
	} else {
		for _, building := range city.Buildings {
			fmt.Printf("- %s\n", buildingToString(building))
		}
	}
	
	g.displayProductionQueue(city)
}

func (g *game) displayProductionQueue(city *city) {
	fmt.Println("\nProduction Queue:")
	if len(city.ProductionQueue) == 0 {
		fmt.Println("Empty")
		return
	}
	
	for i, item := range city.ProductionQueue {
		fmt.Printf("%d. %s: %d/%d\n", i+1, item.Name, item.Progress, item.TotalCost)
	}
}

// ========== Main Function ==========
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	validator := newInputValidator(scanner)
	
	fmt.Println("🏛️ Civilization Game")
	
	numPlayers, err := validator.getIntInput("Enter number of players (2-8): ", 2, 8)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Using default: 4 players")
		numPlayers = 4
	}
	
	game, err := newGame(numPlayers)
	if err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		return
	}
	
	game.Run()
}
