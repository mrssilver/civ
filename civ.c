#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <stdbool.h>

// Game constants
#define MAP_WIDTH 20
#define MAP_HEIGHT 15
#define MAX_PLAYERS 8
#define MAX_CITIES 50
#define MAX_UNITS 100
#define START_YEAR 4000 // 4000 BC
#define END_YEAR 2050  // Game end year

// Terrain types
typedef enum {
    TERRAIN_OCEAN,
    TERRAIN_PLAINS,
    TERRAIN_DESERT,
    TERRAIN_MOUNTAINS,
    TERRAIN_FOREST,
    TERRAIN_HILLS,
    TERRAIN_TUNDRA,
    TERRAIN_JUNGLE,
    TERRAIN_COUNT
} TerrainType;

// Unit types
typedef enum {
    UNIT_SETTLER,
    UNIT_WARRIOR,
    UNIT_ARCHER,
    UNIT_SWORDSMAN,
    UNIT_KNIGHT,
    UNIT_MUSKETEER,
    UNIT_CANNON,
    UNIT_TANK,
    UNIT_COUNT
} UnitType;

// Building types
typedef enum {
    BUILDING_MONUMENT,
    BUILDING_GRANARY,
    BUILDING_LIBRARY,
    BUILDING_TEMPLE,
    BUILDING_BARRACKS,
    BUILDING_WALLS,
    BUILDING_UNIVERSITY,
    BUILDING_FACTORY,
    BUILDING_COUNT
} BuildingType;

// Technology types
typedef enum {
    TECH_AGRICULTURE,
    TECH_POTTERY,
    TECH_WRITING,
    TECH_MATHEMATICS,
    TECH_CONSTRUCTION,
    TECH_PHILOSOPHY,
    TECH_ENGINEERING,
    TECH_EDUCATION,
    TECH_GUNPOWDER,
    TECH_INDUSTRIALIZATION,
    TECH_COUNT
} TechType;

// Civilization types
typedef enum {
    CIV_EGYPT,
    CIV_GREECE,
    CIV_ROME,
    CIV_CHINA,
    CIV_PERSIA,
    CIV_INCA,
    CIV_ENGLAND,
    CIV_FRANCE,
    CIV_COUNT
} CivType;

// City structure
typedef struct {
    char name[50];
    int population;
    int production;
    int food;
    int owner_id;
    BuildingType buildings[20];
    int building_count;
    UnitType production_queue[10];
    int queue_count;
    int production_progress;
} City;

// Map tile structure
typedef struct {
    TerrainType terrain;
    char resource[20];
    bool improved;
    int city_id; // -1 if no city
    int unit_id; // -1 if no unit
    int owner_id; // -1 if unclaimed
} Tile;

// Unit structure
typedef struct {
    UnitType type;
    int health;
    int movement;
    int strength;
    int experience;
    int owner_id;
    int x, y; // Position on map
} Unit;

// Player structure
typedef struct {
    char name[50];
    CivType civ_type;
    City cities[MAX_CITIES];
    int city_count;
    Unit units[MAX_UNITS];
    int unit_count;
    bool techs[TECH_COUNT];
    TechType researching;
    int gold;
    int happiness;
    bool is_ai;
    int relations[MAX_PLAYERS]; // Relations with other players
} Player;

// Game state structure
typedef struct {
    int year;
    Tile map[MAP_HEIGHT][MAP_WIDTH];
    Player players[MAX_PLAYERS];
    int player_count;
    int current_player;
    int winner_id;
} Game;

// Terrain names
const char* terrain_names[] = {
    "Ocean", "Plains", "Desert", "Mountains", 
    "Forest", "Hills", "Tundra", "Jungle"
};

// Unit names
const char* unit_names[] = {
    "Settler", "Warrior", "Archer", "Swordsman",
    "Knight", "Musketeer", "Cannon", "Tank"
};

// Building names
const char* building_names[] = {
    "Monument", "Granary", "Library", "Temple",
    "Barracks", "Walls", "University", "Factory"
};

// Tech names
const char* tech_names[] = {
    "Agriculture", "Pottery", "Writing", "Mathematics",
    "Construction", "Philosophy", "Engineering", "Education",
    "Gunpowder", "Industrialization"
};

// Civ names
const char* civ_names[] = {
    "Egypt", "Greece", "Rome", "China",
    "Persia", "Inca", "England", "France"
};

// Function prototypes
void init_game(Game* game, int player_count);
void generate_map(Game* game);
void create_players(Game* game, int player_count);
void run_game(Game* game);
void player_turn(Game* game);
void ai_turn(Game* game);
void end_year(Game* game);
bool check_game_over(Game* game);
void display_winner(Game* game);
void display_map(Game* game);
void display_player_status(Game* game, int player_id);
void manage_cities(Game* game);
void move_unit(Game* game);
void found_city(Game* game);
void research_tech(Game* game);
void produce_unit(City* city, UnitType type);
void build_building(City* city, BuildingType type);
int calculate_score(Game* game, int player_id);

int main() {
    srand(time(NULL));
    
    int player_count;
    printf("Enter number of players (2-8): ");
    scanf("%d", &player_count);
    
    if (player_count < 2 || player_count > 8) {
        printf("Invalid number of players. Using default 4 players.\n");
        player_count = 4;
    }
    
    Game game;
    init_game(&game, player_count);
    run_game(&game);
    
    return 0;
}

// Initialize game state
void init_game(Game* game, int player_count) {
    game->year = START_YEAR;
    game->player_count = player_count;
    game->current_player = 0;
    game->winner_id = -1;
    
    generate_map(game);
    create_players(game, player_count);
}

// Generate game map
void generate_map(Game* game) {
    for (int y = 0; y < MAP_HEIGHT; y++) {
        for (int x = 0; x < MAP_WIDTH; x++) {
            // Random terrain
            TerrainType terrain = rand() % TERRAIN_COUNT;
            
            // Add resources (10% chance)
            const char* resources[] = {"", "Wheat", "Fish", "Gold", "Iron", "Horses"};
            char resource[20] = "";
            if (rand() % 10 == 0) {
                strcpy(resource, resources[rand() % (sizeof(resources)/sizeof(resources[0]))]);
            }
            
            game->map[y][x] = (Tile){
                .terrain = terrain,
                .resource = "",
                .improved = false,
                .city_id = -1,
                .unit_id = -1,
                .owner_id = -1
            };
            
            if (strlen(resource) > 0) {
                strcpy(game->map[y][x].resource, resource);
            }
        }
    }
}

// Create players
void create_players(Game* game, int player_count) {
    for (int i = 0; i < player_count; i++) {
        Player player;
        strcpy(player.name, civ_names[i]);
        player.civ_type = i;
        player.city_count = 0;
        player.unit_count = 0;
        player.gold = 100;
        player.happiness = 100;
        player.is_ai = i > 0; // First player is human
        
        // Initialize techs
        for (int j = 0; j < TECH_COUNT; j++) {
            player.techs[j] = false;
        }
        player.techs[TECH_AGRICULTURE] = true; // Starting tech
        player.researching = TECH_POTTERY;
        
        // Initialize relations
        for (int j = 0; j < MAX_PLAYERS; j++) {
            player.relations[j] = 0;
        }
        
        // Find a valid starting position
        int start_x, start_y;
        do {
            start_x = rand() % MAP_WIDTH;
            start_y = rand() % MAP_HEIGHT;
        } while (game->map[start_y][start_x].terrain == TERRAIN_OCEAN || 
                 game->map[start_y][start_x].terrain == TERRAIN_MOUNTAINS);
        
        // Create capital city
        City capital;
        snprintf(capital.name, sizeof(capital.name), "%s Capital", player.name);
        capital.population = 1;
        capital.production = 0;
        capital.food = 0;
        capital.owner_id = i;
        capital.building_count = 0;
        capital.queue_count = 0;
        capital.production_progress = 0;
        
        player.cities[player.city_count++] = capital;
        game->map[start_y][start_x].city_id = player.city_count - 1;
        game->map[start_y][start_x].owner_id = i;
        
        // Create starting units
        Unit settler = {UNIT_SETTLER, 100, 2, 0, 0, i, start_x, start_y};
        Unit warrior = {UNIT_WARRIOR, 100, 2, 10, 0, i, start_x, start_y};
        
        player.units[player.unit_count++] = settler;
        player.units[player.unit_count++] = warrior;
        game->map[start_y][start_x].unit_id = player.unit_count - 1;
        
        // Place warrior nearby
        int warrior_x = (start_x + 1) % MAP_WIDTH;
        int warrior_y = start_y;
        game->map[warrior_y][warrior_x].unit_id = player.unit_count - 1;
        warrior.x = warrior_x;
        warrior.y = warrior_y;
        player.units[player.unit_count - 1] = warrior;
        
        game->players[i] = player;
    }
}

// Main game loop
void run_game(Game* game) {
    while (!check_game_over(game)) {
        Player* current = &game->players[game->current_player];
        
        printf("\n======= %s's Turn (%d BC) =======\n", current->name, game->year);
        
        if (current->is_ai) {
            ai_turn(game);
        } else {
            player_turn(game);
        }
        
        // Move to next player
        game->current_player = (game->current_player + 1) % game->player_count;
        
        // End of year processing
        if (game->current_player == 0) {
            end_year(game);
        }
    }
    
    display_winner(game);
}

// AI player turn
void ai_turn(Game* game) {
    Player* player = &game->players[game->current_player];
    printf("\n🤖 %s's turn (AI)\n", player->name);
    
    // Simple AI behavior
    for (int i = 0; i < player->city_count; i++) {
        City* city = &player->cities[i];
        
        // Randomly decide to produce a unit or building
        if (city->queue_count == 0) {
            if (rand() % 2 == 0) {
                // Produce a unit
                UnitType unit = rand() % UNIT_COUNT;
                city->production_queue[city->queue_count++] = unit;
                city->production_progress = 100;
                printf("🏭 Started producing %s\n", unit_names[unit]);
            } else {
                // Build a building
                BuildingType building = rand() % BUILDING_COUNT;
                city->production_queue[city->queue_count++] = building;
                city->production_progress = 50;
                printf("🏗️ Started building %s\n", building_names[building]);
            }
        }
    }
    
    // Move units randomly
    for (int i = 0; i < player->unit_count; i++) {
        Unit* unit = &player->units[i];
        int dx = (rand() % 3) - 1; // -1, 0, or 1
        int dy = (rand() % 3) - 1;
        
        int new_x = (unit->x + dx + MAP_WIDTH) % MAP_WIDTH;
        int new_y = (unit->y + dy + MAP_HEIGHT) % MAP_HEIGHT;
        
        // Check if move is valid
        if (game->map[new_y][new_x].terrain != TERRAIN_OCEAN && 
            game->map[new_y][new_x].terrain != TERRAIN_MOUNTAINS) {
            
            // Clear old position
            game->map[unit->y][unit->x].unit_id = -1;
            
            // Set new position
            unit->x = new_x;
            unit->y = new_y;
            game->map[new_y][new_x].unit_id = i;
            
            printf("🚶 Moved unit to (%d, %d)\n", new_x, new_y);
        }
    }
    
    printf("🤖 End of turn\n");
}

// Player turn
void player_turn(Game* game) {
    int choice;
    
    while (1) {
        printf("\n🎮 Player Menu:\n");
        printf("1. View Map\n");
        printf("2. Manage Cities\n");
        printf("3. Move Units\n");
        printf("4. Found City\n");
        printf("5. Research Technology\n");
        printf("6. View Status\n");
        printf("7. End Turn\n");
        printf("Choose an action: ");
        scanf("%d", &choice);
        
        switch (choice) {
            case 1:
                display_map(game);
                break;
            case 2:
                manage_cities(game);
                break;
            case 3:
                move_unit(game);
                break;
            case 4:
                found_city(game);
                break;
            case 5:
                research_tech(game);
                break;
            case 6:
                display_player_status(game, game->current_player);
                break;
            case 7:
                printf("Ending turn...\n");
                return;
            default:
                printf("Invalid choice\n");
        }
    }
}

// End of year processing
void end_year(Game* game) {
    game->year += 10;
    
    printf("\n📅 Year %d BC\n", game->year);
    
    // Update all players
    for (int p = 0; p < game->player_count; p++) {
        Player* player = &game->players[p];
        
        // Update cities
        for (int c = 0; c < player->city_count; c++) {
            City* city = &player->cities[c];
            
            // City growth
            city->population += rand() % 2; // 0 or 1
            city->food += city->population * 2;
            
            // Process production queue
            if (city->queue_count > 0) {
                if (city->production_progress <= 0) {
                    // Production complete
                    UnitType unit = city->production_queue[0];
                    
                    // Create new unit
                    Unit new_unit;
                    new_unit.type = unit;
                    new_unit.health = 100;
                    new_unit.owner_id = p;
                    new_unit.x = -1;
                    new_unit.y = -1;
                    
                    // Find position near city
                    for (int y = city->population; y >= 0; y--) {
                        for (int x = city->population; x >= 0; x--) {
                            int check_x = (x + city->population) % MAP_WIDTH;
                            int check_y = (y + city->population) % MAP_HEIGHT;
                            
                            if (game->map[check_y][check_x].unit_id == -1) {
                                new_unit.x = check_x;
                                new_unit.y = check_y;
                                game->map[check_y][check_x].unit_id = player->unit_count;
                                break;
                            }
                        }
                        if (new_unit.x != -1) break;
                    }
                    
                    if (new_unit.x != -1) {
                        player->units[player->unit_count++] = new_unit;
                        printf("🏭 %s produced a %s\n", city->name, unit_names[unit]);
                    }
                    
                    // Remove from queue
                    for (int i = 0; i < city->queue_count - 1; i++) {
                        city->production_queue[i] = city->production_queue[i + 1];
                    }
                    city->queue_count--;
                } else {
                    city->production_progress -= 10; // Work on production
                }
            }
        }
        
        // Research technology
        if (rand() % 100 < 30) { // 30% chance to complete research
            player->techs[player->researching] = true;
            printf("🔬 %s researched %s\n", player->name, tech_names[player->researching]);
            player->researching = (player->researching + 1) % TECH_COUNT;
        }
    }
}

// Check if game is over
bool check_game_over(Game* game) {
    // Time victory
    if (game->year >= END_YEAR) {
        int highest_score = 0;
        for (int i = 0; i < game->player_count; i++) {
            int score = calculate_score(game, i);
            if (score > highest_score) {
                highest_score = score;
                game->winner_id = i;
            }
        }
        return true;
    }
    
    // Conquest victory
    for (int i = 0; i < game->player_count; i++) {
        if (game->players[i].city_count == 0) continue; // Player eliminated
        
        bool all_conquered = true;
        for (int j = 0; j < game->player_count; j++) {
            if (i != j && game->players[j].city_count > 0) {
                all_conquered = false;
                break;
            }
        }
        
        if (all_conquered) {
            game->winner_id = i;
            return true;
        }
    }
    
    return false;
}

// Display winner
void display_winner(Game* game) {
    Player* winner = &game->players[game->winner_id];
    printf("\n🏆🏆🏆 Game Over! 🏆🏆🏆\n");
    printf("🎉 Winner: %s\n", winner->name);
    printf("Age: %d BC | Score: %d\n", game->year, calculate_score(game, game->winner_id));
    
    // Display scores
    printf("\nFinal Scores:\n");
    for (int i = 0; i < game->player_count; i++) {
        printf("%s: %d\n", game->players[i].name, calculate_score(game, i));
    }
}

// Display map
void display_map(Game* game) {
    printf("\n🗺️ World Map:\n");
    for (int y = 0; y < MAP_HEIGHT; y++) {
        for (int x = 0; x < MAP_WIDTH; x++) {
            Tile tile = game->map[y][x];
            char symbol;
            
            switch (tile.terrain) {
                case TERRAIN_OCEAN: symbol = '~'; break;
                case TERRAIN_PLAINS: symbol = '.'; break;
                case TERRAIN_DESERT: symbol = 'd'; break;
                case TERRAIN_MOUNTAINS: symbol = '^'; break;
                case TERRAIN_FOREST: symbol = '*'; break;
                case TERRAIN_HILLS: symbol = 'h'; break;
                case TERRAIN_TUNDRA: symbol = 't'; break;
                case TERRAIN_JUNGLE: symbol = 'j'; break;
                default: symbol = '?';
            }
            
            if (tile.unit_id != -1) {
                Unit unit = game->players[game->map[y][x].owner_id].units[tile.unit_id];
                switch (unit.type) {
                    case UNIT_SETTLER: symbol = 'S'; break;
                    case UNIT_WARRIOR: symbol = 'W'; break;
                    case UNIT_ARCHER: symbol = 'A'; break;
                    case UNIT_SWORDSMAN: symbol = 's'; break;
                    case UNIT_KNIGHT: symbol = 'K'; break;
                    case UNIT_MUSKETEER: symbol = 'M'; break;
                    case UNIT_CANNON: symbol = 'C'; break;
                    case UNIT_TANK: symbol = 'T'; break;
                }
            }
            
            if (tile.city_id != -1) {
                symbol = 'C';
            }
            
            printf("%c ", symbol);
        }
        printf("\n");
    }
}

// Display player status
void display_player_status(Game* game, int player_id) {
    Player* player = &game->players[player_id];
    printf("\n📊 %s's Status\n", player->name);
    printf("Gold: %d\n", player->gold);
    printf("Happiness: %d\n", player->happiness);
    printf("Researching: %s\n", tech_names[player->researching]);
    
    printf("\n🏙️ Cities (%d):\n", player->city_count);
    for (int i = 0; i < player->city_count; i++) {
        printf("- %s (Pop: %d)\n", player->cities[i].name, player->cities[i].population);
    }
    
    printf("\n⚔️ Units (%d):\n", player->unit_count);
    for (int i = 0; i < player->unit_count; i++) {
        printf("- %s at (%d, %d)\n", 
               unit_names[player->units[i].type], 
               player->units[i].x, player->units[i].y);
    }
    
    printf("\n🔬 Technologies:\n");
    for (int i = 0; i < TECH_COUNT; i++) {
        if (player->techs[i]) {
            printf("- %s\n", tech_names[i]);
        }
    }
}

// Manage cities
void manage_cities(Game* game) {
    Player* player = &game->players[game->current_player];
    
    if (player->city_count == 0) {
        printf("You have no cities!\n");
        return;
    }
    
    printf("\n🏙️ Your Cities:\n");
    for (int i = 0; i < player->city_count; i++) {
        printf("%d. %s (Pop: %d)\n", i+1, player->cities[i].name, player->cities[i].population);
    }
    
    int city_choice;
    printf("Select a city: ");
    scanf("%d", &city_choice);
    city_choice--; // Convert to 0-based index
    
    if (city_choice < 0 || city_choice >= player->city_count) {
        printf("Invalid city selection.\n");
        return;
    }
    
    City* city = &player->cities[city_choice];
    
    int action;
    do {
        printf("\nManaging %s\n", city->name);
        printf("1. Produce Unit\n");
        printf("2. Build Building\n");
        printf("3. View Production Queue\n");
        printf("4. Back\n");
        printf("Choose an action: ");
        scanf("%d", &action);
        
        switch (action) {
            case 1: {
                printf("\nAvailable Units:\n");
                for (int i = 0; i < UNIT_COUNT; i++) {
                    printf("%d. %s\n", i+1, unit_names[i]);
                }
                
                int unit_choice;
                printf("Select a unit to produce: ");
                scanf("%d", &unit_choice);
                unit_choice--; // Convert to 0-based index
                
                if (unit_choice >= 0 && unit_choice < UNIT_COUNT) {
                    city->production_queue[city->queue_count++] = unit_choice;
                    city->production_progress = 100;
                    printf("Started producing %s\n", unit_names[unit_choice]);
                } else {
                    printf("Invalid unit selection.\n");
                }
                break;
            }
            case 2: {
                printf("\nAvailable Buildings:\n");
                for (int i = 0; i < BUILDING_COUNT; i++) {
                    printf("%d. %s\n", i+1, building_names[i]);
                }
                
                int building_choice;
                printf("Select a building to construct: ");
                scanf("%d", &building_choice);
                building_choice--; // Convert to 0-based index
                
                if (building_choice >= 0 && building_choice < BUILDING_COUNT) {
                    city->production_queue[city->queue_count++] = building_choice;
                    city->production_progress = 50;
                    printf("Started building %s\n", building_names[building_choice]);
                } else {
                    printf("Invalid building selection.\n");
                }
                break;
            }
            case 3: {
                printf("\nProduction Queue:\n");
                for (int i = 0; i < city->queue_count; i++) {
                    printf("%d. %s (%d%% complete)\n", i+1, 
                           unit_names[city->production_queue[i]], 
                           (100 - city->production_progress));
                }
                break;
            }
            case 4:
                return;
            default:
                printf("Invalid option\n");
        }
    } while (action != 4);
}

// Move unit
void move_unit(Game* game) {
    Player* player = &game->players[game->current_player];
    
    if (player->unit_count == 0) {
        printf("You have no units!\n");
        return;
    }
    
    printf("\n⚔️ Your Units:\n");
    for (int i = 0; i < player->unit_count; i++) {
        printf("%d. %s at (%d, %d)\n", 
               i+1, unit_names[player->units[i].type],
               player->units[i].x, player->units[i].y);
    }
    
    int unit_choice;
    printf("Select a unit to move: ");
    scanf("%d", &unit_choice);
    unit_choice--; // Convert to 0-based index
    
    if (unit_choice < 0 || unit_choice >= player->unit_count) {
        printf("Invalid unit selection.\n");
        return;
    }
    
    Unit* unit = &player->units[unit_choice];
    
    int dx, dy;
    printf("Enter movement direction (x y): ");
    scanf("%d %d", &dx, &dy);
    
    int new_x = (unit->x + dx + MAP_WIDTH) % MAP_WIDTH;
    int new_y = (unit->y + dy + MAP_HEIGHT) % MAP_HEIGHT;
    
    // Check if move is valid
    if (game->map[new_y][new_x].terrain == TERRAIN_OCEAN || 
        game->map[new_y][new_x].terrain == TERRAIN_MOUNTAINS) {
        printf("Cannot move to that terrain.\n");
        return;
    }
    
    // Clear old position
    game->map[unit->y][unit->x].unit_id = -1;
    
    // Set new position
    unit->x = new_x;
    unit->y = new_y;
    game->map[new_y][new_x].unit_id = unit_choice;
    
    printf("Moved unit to (%d, %d)\n", new_x, new_y);
}

// Found a new city
void found_city(Game* game) {
    Player* player = &game->players[game->current_player];
    
    // Find a settler
    int settler_id = -1;
    for (int i = 0; i < player->unit_count; i++) {
        if (player->units[i].type == UNIT_SETTLER) {
            settler_id = i;
            break;
        }
    }
    
    if (settler_id == -1) {
        printf("You have no settler units!\n");
        return;
    }
    
    Unit* settler = &player->units[settler_id];
    
    // Check if position is valid
    if (game->map[settler->y][settler->x].city_id != -1) {
        printf("There is already a city here!\n");
        return;
    }
    
    char city_name[50];
    printf("Enter name for new city: ");
    scanf("%s", city_name);
    
    // Create new city
    City new_city;
    snprintf(new_city.name, sizeof(new_city.name), "%s", city_name);
    new_city.population = 1;
    new_city.production = 0;
    new_city.food = 0;
    new_city.owner_id = game->current_player;
    new_city.building_count = 0;
    new_city.queue_count = 0;
    new_city.production_progress = 0;
    
    player->cities[player->city_count] = new_city;
    game->map[settler->y][settler->x].city_id = player->city_count;
    game->map[settler->y][settler->x].owner_id = game->current_player;
    player->city_count++;
    
    // Remove settler unit
    game->map[settler->y][settler->x].unit_id = -1;
    for (int i = settler_id; i < player->unit_count - 1; i++) {
        player->units[i] = player->units[i + 1];
    }
    player->unit_count--;
    
    printf("🏙️ Founded new city: %s\n", city_name);
}

// Research technology
void research_tech(Game* game) {
    Player* player = &game->players[game->current_player];
    
    printf("\n🔬 Available Technologies:\n");
    for (int i = 0; i < TECH_COUNT; i++) {
        if (!player->techs[i]) {
            printf("%d. %s\n", i+1, tech_names[i]);
        }
    }
    
    int tech_choice;
    printf("Select a technology to research: ");
    scanf("%d", &tech_choice);
    tech_choice--; // Convert to 0-based index
    
    if (tech_choice >= 0 && tech_choice < TECH_COUNT && !player->techs[tech_choice]) {
        player->researching = tech_choice;
        printf("Started researching %s\n", tech_names[tech_choice]);
    } else {
        printf("Invalid technology selection.\n");
    }
}

// Calculate player score
int calculate_score(Game* game, int player_id) {
    Player* player = &game->players[player_id];
    int score = 0;
    
    // City points
    score += player->city_count * 100;
    
    // Population points
    for (int i = 0; i < player->city_count; i++) {
        score += player->cities[i].population * 50;
    }
    
    // Tech points
    for (int i = 0; i < TECH_COUNT; i++) {
        if (player->techs[i]) {
            score += 50;
        }
    }
    
    // Territory points
    for (int y = 0; y < MAP_HEIGHT; y++) {
        for (int x = 0; x < MAP_WIDTH; x++) {
            if (game->map[y][x].owner_id == player_id) {
                score += 5;
            }
        }
    }
    
    return score;
}

