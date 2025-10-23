


🏛️ Civilization Game in Zig: Features and Instructions

🎮 Game Features

1. Map System

  ◦ 20x15 grid world map

  ◦ 8 terrain types: Ocean, Plains, Desert, Mountains, Forest, Hills, Tundra, Jungle

  ◦ Resources: Wheat, Fish, Gold, Iron, Horses

2. Civilization Management

  ◦ City development: Population growth, production queues

  ◦ Unit movement: Settlers, Warriors, Archers, Knights, Tanks, etc.

  ◦ Technology research: Agriculture to Industrialization

3. Player Types

  ◦ Human player with interactive menu

  ◦ AI opponents with basic decision-making

  ◦ 8 civilizations: Egypt, Greece, Rome, China, Persia, Inca, England, France

4. Victory Conditions

  ◦ Time Victory: Highest score in 2050 AD

  ◦ Conquest Victory: Eliminate all opponents

5. Game Mechanics

  ◦ Turn-based system (10 years per turn)

  ◦ Production queues for units and buildings

  ◦ Technology research progression

  ◦ Territory control scoring

🕹️ How to Play

1. Compile and Run

zig build-exe civilization.zig
./civilization


2. Game Setup

  ◦ Enter number of players (2-8)

  ◦ First player is human, others are AI

3. Game Controls

  ◦ Main Menu:

    ▪       1. View Map

    ▪       2. Manage Cities

    ▪       3. Move Units

    ▪       4. Found City

    ▪       5. Research Technology

    ▪       6. View Status

    ▪       7. End Turn

4. City Management

  ◦ Produce units: Settlers, Warriors, Archers, etc.

  ◦ Construct buildings: Monuments, Granaries, Libraries, etc.

  ◦ View production queues

5. Unit Movement

  ◦ Select a unit and enter movement direction (dx dy)

  ◦ Units can't move to ocean or mountain tiles

6. Strategy Tips

  ◦ Expand early with Settlers

  ◦ Balance military and economic development

  ◦ Research technologies to unlock advanced units

  ◦ Control territory to increase your score

🧩 Game Structure

1. Data Structures

  ◦ Game: Overall game state

  ◦ Player: Civilization data

  ◦ City: City information

  ◦ Unit: Military units

  ◦ Tile: Map grid cells

2. Key Functions

  ◦ init(): Initialize game state

  ◦ generateMap(): Create random world map

  ◦ run(): Main game loop

  ◦ playerTurn(): Human player actions

  ◦ aiTurn(): Computer player logic

  ◦ endYear(): Yearly processing

  ◦ checkGameOver(): Victory conditions

3. Game Flow

  ◦ Players take turns managing their civilizations

  ◦ Each full round advances the game by 10 years

  ◦ Game ends when time runs out or one player conquers all others

⚙️ Zig Features Used

1. Memory Management

  ◦ Arena allocator for efficient memory management

  ◦ Explicit allocation and deallocation

2. Type Safety

  ◦ Strongly typed enums for game elements

  ◦ Compile-time checks for game logic

3. Error Handling

  ◦ Zig's error handling with try and catch

  ◦ Proper error propagation

4. Standard Library

  ◦ ArrayLists for dynamic collections

  ◦ Random number generation

  ◦ Input/output handling

5. Union Types

  ◦ Production queue with union of Unit and Building
🏛️ Civilization Game in C: Features and Instructions

🎮 Game Features

1. Map System

  ◦ 20x15 grid world map

  ◦ 8 terrain types: Ocean, Plains, Desert, Mountains, Forest, Hills, Tundra, Jungle

  ◦ Resources: Wheat, Fish, Gold, Iron, Horses

2. Civilization Management

  ◦ City development: Population growth, production queues

  ◦ Unit movement: Settlers, Warriors, Archers, Knights, Tanks, etc.

  ◦ Technology research: Agriculture to Industrialization

3. Player Types

  ◦ Human player with interactive menu

  ◦ AI opponents with basic decision-making

  ◦ 8 civilizations: Egypt, Greece, Rome, China, Persia, Inca, England, France

4. Victory Conditions

  ◦ Time Victory: Highest score in 2050 AD

  ◦ Conquest Victory: Eliminate all opponents

5. Game Mechanics

  ◦ Turn-based system (10 years per turn)

  ◦ Production queues for units and buildings

  ◦ Technology research progression

  ◦ Territory control scoring

🕹️ How to Play

1. Compile the Game

gcc civilization.c -o civilization


2. Run the Game

./civilization


3. Game Controls

  ◦ Main Menu:

    ▪       1. View Map

    ▪       2. Manage Cities

    ▪       3. Move Units

    ▪       4. Found City

    ▪       5. Research Technology

    ▪       6. View Status

    ▪       7. End Turn

4. City Management

  ◦ Produce units: Settlers, Warriors, Archers, etc.

  ◦ Construct buildings: Monuments, Granaries, Libraries, etc.

  ◦ View production queues

5. Unit Movement

  ◦ Select a unit and enter movement direction (x y)

  ◦ Units can't move to ocean or mountain tiles

6. Strategy Tips

  ◦ Expand early with Settlers

  ◦ Balance military and economic development

  ◦ Research technologies to unlock advanced units

  ◦ Control territory to increase your score

🧩 Game Structure

1. Data Structures

  ◦ Game: Overall game state

  ◦ Player: Civilization data

  ◦ City: City information

  ◦ Unit: Military units

  ◦ Tile: Map grid cells

2. Key Functions

  ◦ init_game(): Initialize game state

  ◦ generate_map(): Create random world map

  ◦ run_game(): Main game loop

  ◦ player_turn(): Human player actions

  ◦ ai_turn(): Computer player logic

  ◦ end_year(): Yearly processing

  ◦ check_game_over(): Victory conditions

3. Game Flow

  ◦ Players take turns managing their civilizations

  ◦ Each full round advances the game by 10 years

  ◦ Game ends when time runs out or one player conquers all others

🌟 Game Experience

1. Historical Simulation

  ◦ Start in 4000 BC

  ◦ Advance through technological eras

  ◦ Build your civilization from ancient times to modern era

2. Strategic Depth

  ◦ Balance expansion, military, and research

  ◦ Make decisions that affect your civilization's growth

  ◦ Compete against AI opponents with different strategies

3. Replayability

  ◦ Random map generation

  ◦ Multiple victory conditions

  ◦ Different civilizations to play

🚀 Getting Started

1. Compile the game:

gcc civilization.c -o civilization


2. Run the game:

./civilization


3. Game setup:

  ◦ Enter number of players (2-8)

  ◦ First player is human, others are AI

4. Gameplay:

  ◦ Build cities

  ◦ Train units

  ◦ Research technologies

  ◦ Expand your territory

  ◦ Achieve victory before 2050 AD

