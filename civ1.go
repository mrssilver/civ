//3阶导数的英文是 third derivative。

//在更正式或书面的情况下，你可能会看到：

//• Third-order derivative

//•  f'''(x)  （对函数  f(x)  的三阶导数）

//•  \frac{d^3y}{dx^3}  （莱布尼茨记法）

//更高阶的导数（n阶）通常称为 nth derivative 或 nth-order derivative。
//compl
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

// 🏛️ 游戏常量
const (
	MAP_WIDTH    = 20
	MAP_HEIGHT   = 15
	MAX_PLAYERS  = 8
	START_YEAR   = 4000 // 公元前4000年
	END_YEAR     = 2050 // 游戏结束年份
)

// 🌍 地形类型
type TerrainType int

const (
	TERRAIN_OCEAN TerrainType = iota
	TERRAIN_PLAINS
	TERRAIN_DESERT
	TERRAIN_MOUNTAINS
	TERRAIN_FOREST
	TERRAIN_HILLS
	TERRAIN_TUNDRA
	TERRAIN_JUNGLE
)

// 🏛️ 建筑类型
type BuildingType int

const (
	BUILDING_MONUMENT BuildingType = iota
	BUILDING_GRANARY
	BUILDING_LIBRARY
	BUILDING_TEMPLE
	BUILDING_BARRACKS
	BUILDING_WALLS
	BUILDING_UNIVERSITY
	BUILDING_FACTORY
)

// 🔬 科技类型
type TechType int

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
)

// ⚔️ 单位类型
type UnitType int

const (
	UNIT_SETTLER UnitType = iota
	UNIT_WARRIOR
	UNIT_ARCHER
	UNIT_SWORDSMAN
	UNIT_KNIGHT
	UNIT_MUSKETEER
	UNIT_CANNON
	UNIT_TANK
)

// 🏛️ 文明类型
type CivilizationType int

const (
	CIV_EGYPT CivilizationType = iota
	CIV_GREECE
	CIV_ROME
	CIV_CHINA
	CIV_PERSIA
	CIV_INCA
	CIV_ENGLAND
	CIV_FRANCE
)

// 🏛️ 城市结构
type City struct {
	Name          string
	Population    int
	Production    int
	Food          int
	Buildings     []BuildingType
	ProductionQueue []interface{} // 可以是UnitType或BuildingType
}

// 🌍 地图单元格
type Tile struct {
	Terrain  TerrainType
	Resource string
	Improved bool // 是否被改良
	City     *City
	Unit     *Unit
	Owner    *Player
}

// ⚔️ 单位结构
type Unit struct {
	Type       UnitType
	Health     int
	Movement   int
	Strength   int
	Experience int
}

// 🏛️ 玩家结构
type Player struct {
	Name         string
	CivType      CivilizationType
	Cities       []*City
	Units        []*Unit
	Techs        map[TechType]bool
	Researching  TechType
	Gold         int
	Happiness    int
	IsAI         bool
	Relations    map[*Player]int // 与其他玩家的关系 (-100到100)
}

// 🎮 游戏状态
type Game struct {
	Year               int
	Map                [][]Tile
	Players            []*Player
	CurrentPlayerIndex int
	Winner             *Player
}

// 🌍 生成地图
func (g *Game) generateMap() {
	g.Map = make([][]Tile, MAP_HEIGHT)
	for y := 0; y < MAP_HEIGHT; y++ {
		g.Map[y] = make([]Tile, MAP_WIDTH)
		for x := 0; x < MAP_WIDTH; x++ {
			// 随机地形
			terrain := TerrainType(rand.Intn(8))

			// 添加资源
			resources := []string{"", "Wheat", "Fish", "Gold", "Iron", "Horses"}
			resource := ""
			if rand.Intn(10) == 0 { // 10%几率有资源
				resource = resources[rand.Intn(len(resources))]
			}

			g.Map[y][x] = Tile{
				Terrain:  terrain,
				Resource: resource,
			}
		}
	}
}

// 🏛️ 创建玩家
func (g *Game) createPlayers(numPlayers int) {
	civNames := []string{
		"埃及", "希腊", "罗马", "中国", "波斯", "印加", "英格兰", "法国",
	}

	for i := 0; i < numPlayers; i++ {
		player := &Player{
			Name:        civNames[i],
			CivType:     CivilizationType(i),
			Techs:       make(map[TechType]bool),
			Gold:        100,
			Happiness:   100,
			IsAI:        i > 0, // 第一个玩家是人类
			Relations:   make(map[*Player]int),
		}

		// 初始科技
		player.Techs[TECH_AGRICULTURE] = true
		player.Researching = TECH_POTTERY

		// 随机位置建立首都
		x, y := rand.Intn(MAP_WIDTH), rand.Intn(MAP_HEIGHT)
		for g.Map[y][x].Terrain == TERRAIN_OCEAN || g.Map[y][x].Terrain == TERRAIN_MOUNTAINS {
			x, y = rand.Intn(MAP_WIDTH), rand.Intn(MAP_HEIGHT)
		}

		capital := &City{
			Name:       player.Name + "首都",
			Population: 1,
		}
		player.Cities = append(player.Cities, capital)
		g.Map[y][x].City = capital
		g.Map[y][x].Owner = player

		// 初始单位
		settler := &Unit{Type: UNIT_SETTLER, Health: 100, Movement: 2}
		warrior := &Unit{Type: UNIT_WARRIOR, Health: 100, Movement: 2, Strength: 10}
		player.Units = append(player.Units, settler, warrior)
		g.Map[y][x].Unit = settler

		// 放置战士在附近
		nearX, nearY := (x+1)%MAP_WIDTH, y
		g.Map[nearY][nearX].Unit = warrior

		g.Players = append(g.Players, player)
	}
}

// 🏛️ 初始化游戏
func NewGame(numPlayers int) *Game {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		Year: START_YEAR,
	}

	game.generateMap()
	game.createPlayers(numPlayers)

	return game
}

// 🎮 主游戏循环
func (g *Game) Run() {
//	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🏛️ 欢迎来到文明游戏!")
	fmt.Println("你将带领一个文明从古代走向现代")

	for {
		// 检查游戏是否结束
		if g.checkGameOver() {
			g.displayWinner()
			return
		}

		// 获取当前玩家
		currentPlayer := g.Players[g.CurrentPlayerIndex]

		fmt.Printf("\n======= %s 的回合 (年份: %d BC) =======\n", currentPlayer.Name, g.Year)

		if currentPlayer.IsAI {
			g.aiTurn(currentPlayer)
		} else {
			g.playerTurn(currentPlayer)
		}

		// 移动到下一个玩家
		g.CurrentPlayerIndex = (g.CurrentPlayerIndex + 1) % len(g.Players)

		// 所有玩家完成回合后推进时间
		if g.CurrentPlayerIndex == 0 {
			g.endTurn()
		}
	}
}

// 📅 结束回合
func (g *Game) endTurn() {
	g.Year += 10 // 每回合推进10年

	// 更新所有玩家状态
	for _, player := range g.Players {
		// 城市增长
		for _, city := range player.Cities {
			city.Population += rand.Intn(2) // 随机增长
			city.Food += city.Population * 2 // 生产食物

			// 处理生产队列
			if len(city.ProductionQueue) > 0 {
				item := city.ProductionQueue[0]
				switch v := item.(type) {
				case UnitType:
					city.Production -= 10
					if city.Production <= 0 {
						// 生产完成
						unit := &Unit{Type: v}
						switch v {
						case UNIT_SETTLER:
							unit.Movement = 2
						case UNIT_WARRIOR:
							unit.Movement = 2
							unit.Strength = 10
						}
						player.Units = append(player.Units, unit)
						city.ProductionQueue = city.ProductionQueue[1:]
						fmt.Printf("🏭 %s 生产了 %s\n", city.Name, unitTypeToString(v))
					}
				case BuildingType:
					city.Production -= 5
					if city.Production <= 0 {
						city.Buildings = append(city.Buildings, v)
						city.ProductionQueue = city.ProductionQueue[1:]
						fmt.Printf("🏗️ %s 建造了 %s\n", city.Name, buildingTypeToString(v))
					}
				}
			}
		}

		// 科技研究
		if rand.Intn(100) < 30 { // 30%几率完成研究
			player.Techs[player.Researching] = true
			fmt.Printf("🔬 %s 研究完成: %s\n", player.Name, techTypeToString(player.Researching))
			player.Researching = TechType((int(player.Researching) + 1) % 10)
		}
	}

	fmt.Printf("\n📅 年份推进至 %d BC\n", g.Year)
}

// 🏆 检查游戏结束
func (g *Game) checkGameOver() bool {
	// 时间结束
	if g.Year >= END_YEAR {
		// 按分数决定胜利者
		highestScore := 0
		var winner *Player
		for _, player := range g.Players {
			score := g.calculateScore(player)
			if score > highestScore {
				highestScore = score
				winner = player
			}
		}
		g.Winner = winner
		return true
	}

	// 征服胜利
	for _, player := range g.Players {
		if len(player.Cities) == 0 {
			continue // 已被消灭
		}

		allConquered := true
		for _, other := range g.Players {
			if other != player && len(other.Cities) > 0 {
				allConquered = false
				break
			}
		}

		if allConquered {
			g.Winner = player
			return true
		}
	}

	return false
}

// 🏆 显示胜利者
func (g *Game) displayWinner() {
	fmt.Println("\n🏆🏆🏆 游戏结束! 🏆🏆🏆")
	fmt.Printf("🎉 胜利者: %s\n", g.Winner.Name)
	fmt.Printf("年份: %d BC | 城市: %d | 科技: %d\n",
		g.Year, len(g.Winner.Cities), len(g.Winner.Techs))

	// 显示分数
	fmt.Println("\n📊 最终分数:")
	for _, player := range g.Players {
		score := g.calculateScore(player)
		fmt.Printf("- %s: %d分\n", player.Name, score)
	}
}

// 🧮 计算分数
func (g *Game) calculateScore(player *Player) int {
	score := 0

	// 城市分数
	score += len(player.Cities) * 100

	// 科技分数
	score += len(player.Techs) * 50

	// 领土分数
	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			if g.Map[y][x].Owner == player {
				score += 5
			}
		}
	}

	return score
}

// 🤖 AI玩家回合
func (g *Game) aiTurn(player *Player) {
	fmt.Printf("\n🤖 %s 的回合\n", player.Name)

	// 简单AI策略
	for _, city := range player.Cities {
		// 随机生产单位或建筑
		if len(city.ProductionQueue) == 0 {
			if rand.Intn(2) == 0 {
				// 生产单位
				unitTypes := []UnitType{UNIT_WARRIOR, UNIT_ARCHER, UNIT_SETTLER}
				city.ProductionQueue = append(city.ProductionQueue, unitTypes[rand.Intn(len(unitTypes))])
				city.Production = 100
			} else {
				// 建造建筑
				buildingTypes := []BuildingType{BUILDING_GRANARY, BUILDING_MONUMENT, BUILDING_BARRACKS}
				city.ProductionQueue = append(city.ProductionQueue, buildingTypes[rand.Intn(len(buildingTypes))])
				city.Production = 50
			}
		}
	}

	// 移动单位
	for _, unit := range player.Units {
		// 寻找单位位置
		var x, y int
		for y = 0; y < MAP_HEIGHT; y++ {
			for x = 0; x < MAP_WIDTH; x++ {
				if g.Map[y][x].Unit == unit {
					goto found
				}
			}
		}
	found:

		// 随机移动
		dx, dy := rand.Intn(3)-1, rand.Intn(3)-1
		newX, newY := (x+dx+MAP_WIDTH)%MAP_WIDTH, (y+dy+MAP_HEIGHT)%MAP_HEIGHT

		// 检查是否可以移动
		if g.Map[newY][newX].Terrain != TERRAIN_OCEAN && g.Map[newY][newX].Terrain != TERRAIN_MOUNTAINS {
			if g.Map[newY][newX].Unit == nil {
				// 移动单位
				g.Map[y][x].Unit = nil
				g.Map[newY][newX].Unit = unit
				fmt.Printf("🚶 %s 单位移动到 (%d,%d)\n", player.Name, newX, newY)
			} else if g.Map[newY][newX].Owner != player {
				// 攻击
				fmt.Printf("⚔️ %s 单位攻击 (%d,%d)\n", player.Name, newX, newY)
				if rand.Intn(100) < 70 { // 70%胜率
					// 获胜
					g.Map[newY][newX].Unit = nil
					g.Map[newY][newX].Owner = player
					fmt.Printf("✅ 攻击成功! 占领 (%d,%d)\n", newX, newY)
				} else {
					// 失败
					g.Map[y][x].Unit = nil
					fmt.Printf("❌ 攻击失败! 单位被消灭\n")
				}
			}
		}
	}

	fmt.Println("🤖 结束回合")
}

// 🎮 玩家回合
func (g *Game) playerTurn(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n🎮 请选择行动:")
		fmt.Println("1. 管理城市")
		fmt.Println("2. 移动单位")
		fmt.Println("3. 建立城市")
		fmt.Println("4. 研究科技")
		fmt.Println("5. 外交关系")
		fmt.Println("6. 查看地图")
		fmt.Println("7. 查看状态")
		fmt.Println("8. 结束回合")

		fmt.Print("请输入选项: ")
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("无效输入")
			continue
		}

		switch choice {
		case 1:
			g.manageCities(player)
		case 2:
			g.moveUnits(player)
		case 3:
			g.foundCity(player)
		case 4:
			g.researchTech(player)
		case 5:
			g.diplomacy(player)
		case 6:
			g.displayMap()
		case 7:
			g.displayStatus(player)
		case 8:
			fmt.Println("结束回合")
			return
		default:
			fmt.Println("无效选项")
		}
	}
}

// 🏙️ 管理城市
func (g *Game) manageCities(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	if len(player.Cities) == 0 {
		fmt.Println("你没有城市")
		return
	}

	// 选择城市
	fmt.Println("\n🏙️ 你的城市:")
	for i, city := range player.Cities {
		fmt.Printf("%d. %s (人口: %d)\n", i+1, city.Name, city.Population)
	}

	fmt.Print("选择城市编号: ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(player.Cities) {
		fmt.Println("无效选择")
		return
	}

	city := player.Cities[index-1]

	// 城市管理菜单
	for {
		fmt.Printf("\n🏙️ 管理城市: %s\n", city.Name)
		fmt.Println("1. 查看信息")
		fmt.Println("2. 生产单位")
		fmt.Println("3. 建造建筑")
		fmt.Println("4. 返回")

		fmt.Print("请输入选项: ")
		input, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("无效输入")
			continue
		}

		switch choice {
		case 1:
			g.displayCityInfo(city)
		case 2:
			g.produceUnit(city)
		case 3:
			g.buildBuilding(city)
		case 4:
			return
		default:
			fmt.Println("无效选项")
		}
	}
}

// 🚶 移动单位
func (g *Game) moveUnits(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	if len(player.Units) == 0 {
		fmt.Println("你没有单位")
		return
	}

	// 选择单位
	fmt.Println("\n🚶 你的单位:")
	for i, unit := range player.Units {
		fmt.Printf("%d. %s\n", i+1, unitTypeToString(unit.Type))
	}

	fmt.Print("选择单位编号: ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(player.Units) {
		fmt.Println("无效选择")
		return
	}

	unit := player.Units[index-1]

	// 寻找单位位置
	var x, y int
	for y = 0; y < MAP_HEIGHT; y++ {
		for x = 0; x < MAP_WIDTH; x++ {
			if g.Map[y][x].Unit == unit {
				goto found
			}
		}
	}
	fmt.Println("单位未在地图上")
	return

found:
	fmt.Printf("单位位置: (%d,%d)\n", x, y)
	fmt.Print("输入移动方向 (w上, s下, a左, d右): ")
	input, _ = reader.ReadString('\n')
	direction := strings.TrimSpace(input)

	dx, dy := 0, 0
	switch direction {
	case "w":
		dy = -1
	case "s":
		dy = 1
	case "a":
		dx = -1
	case "d":
		dx = 1
	default:
		fmt.Println("无效方向")
		return
	}

	newX, newY := (x+dx+MAP_WIDTH)%MAP_WIDTH, (y+dy+MAP_HEIGHT)%MAP_HEIGHT

	// 检查是否可以移动
	if g.Map[newY][newX].Terrain == TERRAIN_OCEAN || g.Map[newY][newX].Terrain == TERRAIN_MOUNTAINS {
		fmt.Println("无法移动到该地形")
		return
	}

	if g.Map[newY][newX].Unit != nil {
		if g.Map[newY][newX].Owner == player {
			fmt.Println("该位置已有友方单位")
			return
		} else {
			// 攻击
			fmt.Printf("⚔️ 攻击 (%d,%d) 的单位\n", newX, newY)
			if rand.Intn(100) < 70 { // 70%胜率
				// 获胜
				g.Map[newY][newX].Unit = nil
				g.Map[newY][newX].Owner = player
				fmt.Printf("✅ 攻击成功! 占领 (%d,%d)\n", newX, newY)
			} else {
				// 失败
				g.Map[y][x].Unit = nil
				fmt.Printf("❌ 攻击失败! 单位被消灭\n")
			}
			return
		}
	}

	// 移动单位
	g.Map[y][x].Unit = nil
	g.Map[newY][newX].Unit = unit
	fmt.Printf("✅ 单位移动到 (%d,%d)\n", newX, newY)
}

// 🏙️ 建立城市
func (g *Game) foundCity(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	// 寻找定居者
	var settler *Unit
	for _, unit := range player.Units {
		if unit.Type == UNIT_SETTLER {
			settler = unit
			break
		}
	}

	if settler == nil {
		fmt.Println("没有可用的定居者")
		return
	}

	// 寻找定居者位置
	var x, y int
	for y = 0; y < MAP_HEIGHT; y++ {
		for x = 0; x < MAP_WIDTH; x++ {
			if g.Map[y][x].Unit == settler {
				goto found
			}
		}
	}
	fmt.Println("定居者未在地图上")
	return

found:
	// 检查是否可以建立城市
	if g.Map[y][x].City != nil {
		fmt.Println("该位置已有城市")
		return
	}

	fmt.Print("输入新城市名称: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// 建立城市
	city := &City{
		Name:       name,
		Population: 1,
	}
	player.Cities = append(player.Cities, city)
	g.Map[y][x].City = city
	g.Map[y][x].Owner = player
	g.Map[y][x].Unit = nil // 定居者消失

	// 移除定居者
	for i, unit := range player.Units {
		if unit == settler {
			player.Units = append(player.Units[:i], player.Units[i+1:]...)
			break
		}
	}

	fmt.Printf("🏙️ 建立了新城市: %s\n", name)
}

// 🔬 研究科技
func (g *Game) researchTech(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🔬 可研究的科技:")
	techs := []TechType{
		TECH_POTTERY,
		TECH_WRITING,
		TECH_MATHEMATICS,
		TECH_CONSTRUCTION,
		TECH_PHILOSOPHY,
		TECH_ENGINEERING,
		TECH_EDUCATION,
		TECH_GUNPOWDER,
		TECH_INDUSTRIALIZATION,
	}

	for i, tech := range techs {
		if !player.Techs[tech] {
			fmt.Printf("%d. %s\n", i+1, techTypeToString(tech))
		}
	}

	fmt.Print("选择要研究的科技编号 (0取消): ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(techs) {
		fmt.Println("取消研究")
		return
	}

	tech := techs[index-1]
	player.Researching = tech
	fmt.Printf("🔬 开始研究: %s\n", techTypeToString(tech))
}

// 🤝 外交关系
func (g *Game) diplomacy(player *Player) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🤝 外交关系:")
	for i, other := range g.Players {
		if other != player {
			relation := player.Relations[other]
			status := "中立"
			if relation > 50 {
				status = "友好"
			} else if relation < -50 {
				status = "敌对"
			}
			fmt.Printf("%d. %s: %s (%d)\n", i+1, other.Name, status, relation)
		}
	}

	fmt.Print("选择外交对象编号 (0取消): ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(g.Players)-1 {
		fmt.Println("取消外交")
		return
	}

	// 获取目标玩家
	targetIndex := index
	if targetIndex >= playerIndex(g, player) {
		targetIndex++ // 跳过自己
	}
	target := g.Players[targetIndex-1]

	// 外交行动
	fmt.Println("\n外交行动:")
	fmt.Println("1. 宣战")
	fmt.Println("2. 和平协议")
	fmt.Println("3. 贸易协定")
	fmt.Println("4. 返回")

	fmt.Print("请选择行动: ")
	input, _ = reader.ReadString('\n')
	action, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return
	}

	switch action {
	case 1:
		player.Relations[target] = -100
		fmt.Printf("⚔️ 向 %s 宣战!\n", target.Name)
	case 2:
		player.Relations[target] = 50
		fmt.Printf("🕊️ 与 %s 签订和平协议\n", target.Name)
	case 3:
		player.Relations[target] += 20
		fmt.Printf("🤝 与 %s 签订贸易协定\n", target.Name)
	}
}

// 🗺️ 显示地图
func (g *Game) displayMap() {
	fmt.Println("\n🗺️ 世界地图:")

	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			tile := g.Map[y][x]

			// 显示符号
			symbol := ""
			switch tile.Terrain {
			case TERRAIN_OCEAN:
				symbol = "🌊"
			case TERRAIN_PLAINS:
				symbol = "🌾"
			case TERRAIN_DESERT:
				symbol = "🏜️"
			case TERRAIN_MOUNTAINS:
				symbol = "⛰️"
			case TERRAIN_FOREST:
				symbol = "🌲"
			case TERRAIN_HILLS:
				symbol = "🏞️"
			case TERRAIN_TUNDRA:
				symbol = "❄️"
			case TERRAIN_JUNGLE:
				symbol = "🌴"
			}

			// 显示城市
			if tile.City != nil {
				symbol = "🏙️"
			}

			// 显示单位
			if tile.Unit != nil {
				switch tile.Unit.Type {
				case UNIT_SETTLER:
					symbol = "👨‍🌾"
				case UNIT_WARRIOR:
					symbol = "⚔️"
				case UNIT_ARCHER:
					symbol = "🏹"
				case UNIT_SWORDSMAN:
					symbol = "🗡️"
				case UNIT_KNIGHT:
					symbol = "🐎"
				}
			}

			fmt.Print(symbol)
		}
		fmt.Println()
	}
}

// 📊 显示玩家状态
func (g *Game) displayStatus(player *Player) {
	fmt.Printf("\n📊 %s 的状态\n", player.Name)
	fmt.Printf("年份: %d BC\n", g.Year)
	fmt.Printf("黄金: %d\n", player.Gold)
	fmt.Printf("快乐度: %d\n", player.Happiness)
	fmt.Printf("研究中的科技: %s\n", techTypeToString(player.Researching))

	fmt.Println("\n🏙️ 城市:")
	for _, city := range player.Cities {
		fmt.Printf("- %s (人口: %d)\n", city.Name, city.Population)
	}

	fmt.Println("\n🔬 已掌握的科技:")
	for tech := range player.Techs {
		fmt.Printf("- %s\n", techTypeToString(tech))
	}
}

// 🏙️ 显示城市信息
func (g *Game) displayCityInfo(city *City) {
	fmt.Printf("\n🏙️ 城市: %s\n", city.Name)
	fmt.Printf("人口: %d\n", city.Population)
	fmt.Printf("食物: %d\n", city.Food)
	fmt.Printf("生产力: %d\n", city.Production)

	fmt.Println("\n🏗️ 建筑:")
	for _, building := range city.Buildings {
		fmt.Printf("- %s\n", buildingTypeToString(building))
	}

	fmt.Println("\n🏭 生产队列:")
	for _, item := range city.ProductionQueue {
		switch v := item.(type) {
		case UnitType:
			fmt.Printf("- %s\n", unitTypeToString(v))
		case BuildingType:
			fmt.Printf("- %s\n", buildingTypeToString(v))
		}
	}
}

// ⚔️ 生产单位
func (g *Game) produceUnit(city *City) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n⚔️ 可生产的单位:")
	units := []UnitType{
		UNIT_SETTLER,
		UNIT_WARRIOR,
		UNIT_ARCHER,
		UNIT_SWORDSMAN,
		UNIT_KNIGHT,
	}

	for i, unit := range units {
		fmt.Printf("%d. %s\n", i+1, unitTypeToString(unit))
	}

	fmt.Print("选择单位编号 (0取消): ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(units) {
		fmt.Println("取消生产")
		return
	}

	city.ProductionQueue = append(city.ProductionQueue, units[index-1])
	city.Production = 100
	fmt.Printf("🏭 开始生产: %s\n", unitTypeToString(units[index-1]))
}

// 🏗️ 建造建筑
func (g *Game) buildBuilding(city *City) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🏗️ 可建造的建筑:")
	buildings := []BuildingType{
		BUILDING_MONUMENT,
		BUILDING_GRANARY,
		BUILDING_LIBRARY,
		BUILDING_TEMPLE,
		BUILDING_BARRACKS,
	}

	for i, building := range buildings {
		fmt.Printf("%d. %s\n", i+1, buildingTypeToString(building))
	}

	fmt.Print("选择建筑编号 (0取消): ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index < 1 || index > len(buildings) {
		fmt.Println("取消建造")
		return
	}

	city.ProductionQueue = append(city.ProductionQueue, buildings[index-1])
	city.Production = 50
	fmt.Printf("🏗️ 开始建造: %s\n", buildingTypeToString(buildings[index-1]))
}

// 🧪 地形类型转字符串
func terrainTypeToString(terrain TerrainType) string {
	names := []string{
		"海洋", "平原", "沙漠", "山脉", "森林", "丘陵", "冻土", "丛林",
	}
	return names[terrain]
}

// 🏗️ 建筑类型转字符串
func buildingTypeToString(building BuildingType) string {
	names := []string{
		"纪念碑", "粮仓", "图书馆", "神庙", "兵营", "城墙", "大学", "工厂",
	}
	return names[building]
}

// 🔬 科技类型转字符串
func techTypeToString(tech TechType) string {
	names := []string{
		"农业", "制陶术", "书写", "数学", "建筑学", "哲学", "工程学", "教育", "火药", "工业化",
	}
	return names[tech]
}

// ⚔️ 单位类型转字符串
func unitTypeToString(unit UnitType) string {
	names := []string{
		"定居者", "战士", "弓箭手", "剑士", "骑士", "火枪手", "加农炮", "坦克",
	}
	return names[unit]
}

// 🏛️ 文明类型转字符串
func civTypeToString(civ CivilizationType) string {
	names := []string{
		"埃及", "希腊", "罗马", "中国", "波斯", "印加", "英格兰", "法国",
	}
	return names[civ]
}

// 🏛️ 获取玩家索引
func playerIndex(g *Game, player *Player) int {
	for i, p := range g.Players {
		if p == player {
			return i
		}
	}
	return -1
}

// 🎮 主函数
func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🏛️ 文明游戏")
	fmt.Print("请输入玩家数量 (2-8): ")
	input, _ := reader.ReadString('\n')
	numPlayers, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || numPlayers < 2 || numPlayers > 8 {
		numPlayers = 4
		fmt.Println("使用默认玩家数量: 4")
	}
	game := NewGame(numPlayers)
	game.Run()
}
