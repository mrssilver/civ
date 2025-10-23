

const std = @import("std");
const rand = std.rand;
const heap = std.heap;
const io = std.io;
const fmt = std.fmt;
const mem = std.mem;
const time = std.time;

// Game constants
const MAP_WIDTH = 20;
const MAP_HEIGHT = 15;
const MAX_PLAYERS = 8;
const MAX_CITIES = 50;
const MAX_UNITS = 100;
const START_YEAR = 4000; // 4000 BC
const END_YEAR = 2050;   // Game end year

// Terrain types
const TerrainType = enum {
    Ocean,
    Plains,
    Desert,
    Mountains,
    Forest,
    Hills,
    Tundra,
    Jungle,
};

// Unit types
const UnitType = enum {
    Settler,
    Warrior,
    Archer,
    Swordsman,
    Knight,
    Musketeer,
    Cannon,
    Tank,
};

// Building types
const BuildingType = enum {
    Monument,
    Granary,
    Library,
    Temple,
    Barracks,
    Walls,
    University,
    Factory,
};

// Technology types
const TechType = enum {
    Agriculture,
    Pottery,
    Writing,
    Mathematics,
    Construction,
    Philosophy,
    Engineering,
    Education,
    Gunpowder,
    Industrialization,
};

// Civilization types
const CivType = enum {
    Egypt,
    Greece,
    Rome,
    China,
    Persia,
    Inca,
    England,
    France,
};

// City structure
const City = struct {
    name: []const u8,
    population: u32,
    production: u32,
    food: u32,
    owner_id: u8,
    buildings: std.ArrayList(BuildingType),
    production_queue: std.ArrayList(union(enum) { Unit: UnitType, Building: BuildingType }),
};

// Unit structure
const Unit = struct {
    type: UnitType,
    health: u8,
    movement: u8,
    strength: u8,
    experience: u8,
    owner_id: u8,
    x: u8,
    y: u8,
};

// Player structure
const Player = struct {
    name: []const u8,
    civ_type: CivType,
    cities: std.ArrayList(*City),
    units: std.ArrayList(*Unit),
    techs: std.ArrayList(TechType),
    researching: ?TechType,
    gold: u32,
    happiness: u8,
    is_ai: bool,
    relations: [MAX_PLAYERS]i8, // Relations with other players
};

// Map tile structure
const Tile = struct {
    terrain: TerrainType,
    resource: ?[]const u8,
    improved: bool,
    city: ?*City,
    unit: ?*Unit,
    owner_id: ?u8,
};

// Game state structure
const Game = struct {
    year: u32,
    map: [MAP_HEIGHT][MAP_WIDTH]Tile,
    players: std.ArrayList(Player),
    current_player_index: u8,
    winner_id: ?u8,
    allocator: std.mem.Allocator,

    // Initialize game
    fn init(allocator: std.mem.Allocator, player_count: u8) !Game {
        var self = Game{
            .year = START_YEAR,
            .map = undefined,
            .players = std.ArrayList(Player).init(allocator),
            .current_player_index = 0,
            .winner_id = null,
            .allocator = allocator,
        };
        
        try self.generateMap();
        try self.createPlayers(player_count);
        
        return self;
    }

    // Generate game map
    fn generateMap(self: *Game) !void {
        const resources = [_][]const u8{ "Wheat", "Fish", "Gold", "Iron", "Horses" };
        var rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
        
        for (0..MAP_HEIGHT) |y| {
            for (0..MAP_WIDTH) |x| {
                const terrain = @enumFromInt(TerrainType, rng.random().intRangeAtMost(u8, 0, @typeInfo(TerrainType).Enum.fields.len - 1));
                const has_resource = rng.random().intRangeAtMost(u8, 0, 9) == 0;
                const resource = if (has_resource) resources[rng.random().intRangeAtMost(u8, 0, resources.len - 1)] else null;
                
                self.map[y][x] = Tile{
                    .terrain = terrain,
                    .resource = resource,
                    .improved = false,
                    .city = null,
                    .unit = null,
                    .owner_id = null,
                };
            }
        }
    }

    // Create players
    fn createPlayers(self: *Game, player_count: u8) !void {
        const civ_names = [_][]const u8{ "Egypt", "Greece", "Rome", "China", "Persia", "Inca", "England", "France" };
        var rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
        
        for (0..player_count) |i| {
            const is_ai = i > 0; // First player is human
            const name = civ_names[i];
            
            var player = Player{
                .name = name,
                .civ_type = @enumFromInt(CivType, i),
                .cities = std.ArrayList(*City).init(self.allocator),
                .units = std.ArrayList(*Unit).init(self.allocator),
                .techs = std.ArrayList(TechType).init(self.allocator),
                .researching = TechType.Agriculture,
                .gold = 100,
                .happiness = 100,
                .is_ai = is_ai,
                .relations = undefined,
            };
            
            // Initialize relations
            for (&player.relations) |*relation| {
                relation.* = 0;
            }
            
            // Add starting tech
            try player.techs.append(TechType.Agriculture);
            
            // Find valid starting position
            var start_x: u8 = undefined;
            var start_y: u8 = undefined;
            while (true) {
                start_x = rng.random().intRangeAtMost(u8, 0, MAP_WIDTH - 1);
                start_y = rng.random().intRangeAtMost(u8, 0, MAP_HEIGHT - 1);
                const tile = self.map[start_y][start_x];
                if (tile.terrain != .Ocean and tile.terrain != .Mountains) break;
            }
            
            // Create capital city
            const city_name = try fmt.allocPrint(self.allocator, "{s} Capital", .{name});
            var capital = try self.allocator.create(City);
            capital.* = .{
                .name = city_name,
                .population = 1,
                .production = 0,
                .food = 0,
                .owner_id = @intCast(u8, i),
                .buildings = std.ArrayList(BuildingType).init(self.allocator),
                .production_queue = std.ArrayList(union(enum) { Unit: UnitType, Building: BuildingType }).init(self.allocator),
            };
            try player.cities.append(capital);
            self.map[start_y][start_x].city = capital;
            self.map[start_y][start_x].owner_id = @intCast(u8, i);
            
            // Create starting units
            var settler = try self.allocator.create(Unit);
            settler.* = .{
                .type = .Settler,
                .health = 100,
                .movement = 2,
                .strength = 0,
                .experience = 0,
                .owner_id = @intCast(u8, i),
                .x = start_x,
                .y = start_y,
            };
            try player.units.append(settler);
            self.map[start_y][start_x].unit = settler;
            
            var warrior = try self.allocator.create(Unit);
            warrior.* = .{
                .type = .Warrior,
                .health = 100,
                .movement = 2,
                .strength = 10,
                .experience = 0,
                .owner_id = @intCast(u8, i),
                .x = start_x,
                .y = start_y,
            };
            try player.units.append(warrior);
            
            // Place warrior nearby
            const warrior_x = if (start_x < MAP_WIDTH - 1) start_x + 1 else start_x - 1;
            self.map[start_y][warrior_x].unit = warrior;
            warrior.x = warrior_x;
            warrior.y = start_y;
            
            try self.players.append(player);
        }
    }

    // Run game
    fn run(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        try stdout.print("🏛️ Welcome to Civilization!\n", .{});
        
        while (!self.checkGameOver()) {
            const player = &self.players.items[self.current_player_index];
            
            try stdout.print("\n======= {s}'s Turn ({d} BC) =======\n", .{ player.name, self.year });
            
            if (player.is_ai) {
                try self.aiTurn();
            } else {
                try self.playerTurn();
            }
            
            // Move to next player
            self.current_player_index = @intCast(u8, (self.current_player_index + 1) % self.players.items.len);
            
            // End of year
            if (self.current_player_index == 0) {
                try self.endYear();
            }
        }
        
        try self.displayWinner();
    }

    // AI player turn
    fn aiTurn(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        try stdout.print("🤖 {s} is taking its turn...\n", .{player.name});
        
        // Simple AI actions
        for (player.cities.items) |city| {
            // Randomly decide to produce a unit or building
            if (city.production_queue.items.len == 0) {
                const rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
                if (rng.random().boolean()) {
                    // Produce a unit
                    const unit_type = @enumFromInt(UnitType, rng.random().intRangeAtMost(u8, 0, @typeInfo(UnitType).Enum.fields.len - 1));
                    try city.production_queue.append(.{ .Unit = unit_type });
                    try stdout.print("🏭 Started producing {s}\n", .{@tagName(unit_type)});
                } else {
                    // Build a building
                    const building_type = @enumFromInt(BuildingType, rng.random().intRangeAtMost(u8, 0, @typeInfo(BuildingType).Enum.fields.len - 1));
                    try city.production_queue.append(.{ .Building = building_type });
                    try stdout.print("🏗️ Started building {s}\n", .{@tagName(building_type)});
                }
            }
        }
        
        // Move units randomly
        for (player.units.items) |unit| {
            const rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
            const dx = rng.random().intRangeAtMost(i8, -1, 1);
            const dy = rng.random().intRangeAtMost(i8, -1, 1);
            
            const new_x = @intCast(u8, (@as(i16, unit.x) + dx) % MAP_WIDTH);
            const new_y = @intCast(u8, (@as(i16, unit.y) + dy) % MAP_HEIGHT);
            
            // Check if move is valid
            const tile = self.map[new_y][new_x];
            if (tile.terrain != .Ocean and tile.terrain != .Mountains) {
                // Clear old position
                self.map[unit.y][unit.x].unit = null;
                
                // Set new position
                unit.x = new_x;
                unit.y = new_y;
                self.map[new_y][new_x].unit = unit;
                
                try stdout.print("🚶 Moved unit to ({d},{d})\n", .{ new_x, new_y });
            }
        }
    }

    // Human player turn
    fn playerTurn(self: *Game) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        while (true) {
            try stdout.print(
                \\\n🎮 Menu:
                \\1. View Map
                \\2. Manage Cities
                \\3. Move Units
                \\4. Found City
                \\5. Research Technology
                \\6. View Status
                \\7. End Turn
                \\Choose an option: 
            , .{});
            
            var input: [10]u8 = undefined;
            const len = try stdin.read(&input);
            const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
                try stdout.print("Invalid input\n", .{});
                continue;
            };
            
            switch (choice) {
                1 => try self.displayMap(),
                2 => try self.manageCities(),
                3 => try self.moveUnit(),
                4 => try self.foundCity(),
                5 => try self.researchTech(),
                6 => try self.displayPlayerStatus(),
                7 => {
                    try stdout.print("Ending turn...\n", .{});
                    return;
                },
                else => try stdout.print("Invalid choice\n", .{}),
            }
        }
    }

    // End of year processing
    fn endYear(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        
        self.year += 10;
        try stdout.print("\n📅 Year advanced to {d} BC\n", .{self.year});
        
        // Update all players
        for (self.players.items) |*player| {
            // Update cities
            for (player.cities.items) |city| {
                // City growth
                const rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
                city.population += rng.random().intRangeAtMost(u32, 0, 1);
                city.food += city.population * 2;
                
                // Process production
                if (city.production_queue.items.len > 0) {
                    const item = city.production_queue.items[0];
                    switch (item) {
                        .Unit => |unit_type| {
                            if (city.production >= 100) {
                                // Create unit
                                var unit = try self.allocator.create(Unit);
                                unit.* = .{
                                    .type = unit_type,
                                    .health = 100,
                                    .movement = 2,
                                    .strength = switch (unit_type) {
                                        .Warrior => 10,
                                        .Archer => 8,
                                        .Swordsman => 12,
                                        .Knight => 15,
                                        .Musketeer => 18,
                                        .Cannon => 20,
                                        .Tank => 30,
                                        else => 0,
                                    },
                                    .experience = 0,
                                    .owner_id = city.owner_id,
                                    .x = 0,
                                    .y = 0,
                                };
                                
                                // Find position near city
                                for (0..MAP_HEIGHT) |y| {
                                    for (0..MAP_WIDTH) |x| {
                                        if (self.map[y][x].owner_id == city.owner_id and self.map[y][x].unit == null) {
                                            unit.x = @intCast(u8, x);
                                            unit.y = @intCast(u8, y);
                                            self.map[y][x].unit = unit;
                                            try player.units.append(unit);
                                            try stdout.print("🏭 Produced {s} at ({d},{d})\n", .{ @tagName(unit_type), x, y });
                                            break;
                                        }
                                    }
                                }
                                
                                // Remove from queue
                                _ = city.production_queue.orderedRemove(0);
                                city.production = 0;
                            } else {
                                city.production += 10;
                            }
                        },
                        .Building => |building_type| {
                            if (city.production >= 50) {
                                try city.buildings.append(building_type);
                                try stdout.print("🏗️ Built {s}\n", .{@tagName(building_type)});
                                
                                // Remove from queue
                                _ = city.production_queue.orderedRemove(0);
                                city.production = 0;
                            } else {
                                city.production += 5;
                            }
                        },
                    }
                }
            }
            
            // Research technology
            const rng = rand.DefaultPrng.init(@bitCast(u64, time.nanoTimestamp()));
            if (rng.random().intRangeAtMost(u8, 0, 99) < 30) {
                if (player.researching) |tech| {
                    try player.techs.append(tech);
                    try stdout.print("🔬 Researched {s}\n", .{@tagName(tech)});
                    
                    // Move to next tech
                    const next_tech = @enumFromInt(TechType, @enumToInt(tech) + 1);
                    if (@enumToInt(next_tech) < @typeInfo(TechType).Enum.fields.len) {
                        player.researching = next_tech;
                    } else {
                        player.researching = null;
                    }
                }
            }
        }
    }

    // Check if game is over
    fn checkGameOver(self: *Game) bool {
        // Time victory
        if (self.year >= END_YEAR) {
            var highest_score: u32 = 0;
            for (self.players.items, 0..) |player, i| {
                const score = self.calculateScore(@intCast(u8, i));
                if (score > highest_score) {
                    highest_score = score;
                    self.winner_id = @intCast(u8, i);
                }
            }
            return true;
        }
        
        // Conquest victory
        for (self.players.items, 0..) |player, i| {
            if (player.cities.items.len == 0) continue;
            
            var all_conquered = true;
            for (self.players.items) |other| {
                if (other.cities.items.len > 0 and &player != &other) {
                    all_conquered = false;
                    break;
                }
            }
            
            if (all_conquered) {
                self.winner_id = @intCast(u8, i);
                return true;
            }
        }
        
        return false;
    }

    // Display winner
    fn displayWinner(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        
        if (self.winner_id) |winner_id| {
            const winner = &self.players.items[winner_id];
            try stdout.print(
                \\\n🏆🏆🏆 Game Over! 🏆🏆🏆
                \\🎉 Winner: {s}
                \\Year: {d} BC
                \\Cities: {d}
                \\Technologies: {d}
                \\
            , .{
                winner.name,
                self.year,
                winner.cities.items.len,
                winner.techs.items.len,
            });
        }
    }

    // Display map
    fn displayMap(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        try stdout.print("\n🗺️ World Map:\n", .{});
        
        for (self.map) |row| {
            for (row) |tile| {
                const symbol: u8 = switch (tile.terrain) {
                    .Ocean => '~',
                    .Plains => '.',
                    .Desert => 'd',
                    .Mountains => '^',
                    .Forest => '*',
                    .Hills => 'h',
                    .Tundra => 't',
                    .Jungle => 'j',
                };
                
                if (tile.city != null) {
                    try stdout.print("C", .{});
                } else if (tile.unit != null) {
                    try stdout.print("U", .{});
                } else {
                    try stdout.print("{c}", .{symbol});
                }
                try stdout.print(" ", .{});
            }
            try stdout.print("\n", .{});
        }
    }

    // Display player status
    fn displayPlayerStatus(self: *Game) !void {
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        try stdout.print(
            \\\n📊 Player Status:
            \\Name: {s}
            \\Civilization: {s}
            \\Gold: {d}
            \\Happiness: {d}
            \\Cities: {d}
            \\Units: {d}
            \\Technologies: {d}
            \\
        , .{
            player.name,
            @tagName(player.civ_type),
            player.gold,
            player.happiness,
            player.cities.items.len,
            player.units.items.len,
            player.techs.items.len,
        });
        
        if (player.researching) |tech| {
            try stdout.print("Researching: {s}\n", .{@tagName(tech)});
        }
    }

    // Manage cities
    fn manageCities(self: *Game) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        if (player.cities.items.len == 0) {
            try stdout.print("You have no cities!\n", .{});
            return;
        }
        
        try stdout.print("\n🏙️ Your Cities:\n", .{});
        for (player.cities.items, 0..) |city, i| {
            try stdout.print("{d}. {s} (Pop: {d})\n", .{ i + 1, city.name, city.population });
        }
        
        try stdout.print("Select a city: ", .{});
        var input: [10]u8 = undefined;
        const len = try stdin.read(&input);
        const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
            try stdout.print("Invalid input\n", .{});
            return;
        };
        
        if (choice < 1 or choice > player.cities.items.len) {
            try stdout.print("Invalid choice\n", .{});
            return;
        }
        
        const city = player.cities.items[choice - 1];
        
        while (true) {
            try stdout.print(
                \\\n🏙️ Managing {s}
                \\1. Produce Unit
                \\2. Build Building
                \\3. View Production Queue
                \\4. Back
                \\Choose an option: 
            , .{city.name});
            
            const len2 = try stdin.read(&input);
            const action = std.fmt.parseInt(u8, input[0..len2], 10) catch {
                try stdout.print("Invalid input\n", .{});
                continue;
            };
            
            switch (action) {
                1 => try self.produceUnit(city),
                2 => try self.buildBuilding(city),
                3 => try self.viewProductionQueue(city),
                4 => return,
                else => try stdout.print("Invalid choice\n", .{}),
            }
        }
    }

    // Produce unit
    fn produceUnit(self: *Game, city: *City) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        
        try stdout.print("\n⚔️ Available Units:\n", .{});
        const unit_types = std.enums.values(UnitType);
        for (unit_types, 0..) |unit_type, i| {
            try stdout.print("{d}. {s}\n", .{ i + 1, @tagName(unit_type) });
        }
        
        try stdout.print("Select a unit: ", .{});
        var input: [10]u8 = undefined;
        const len = try stdin.read(&input);
        const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
            try stdout.print("Invalid input\n", .{});
            return;
        };
        
        if (choice < 1 or choice > unit_types.len) {
            try stdout.print("Invalid choice\n", .{});
            return;
        }
        
        try city.production_queue.append(.{ .Unit = unit_types[choice - 1] });
        try stdout.print("🏭 Started producing {s}\n", .{@tagName(unit_types[choice - 1])});
    }

    // Build building
    fn buildBuilding(self: *Game, city: *City) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        
        try stdout.print("\n🏗️ Available Buildings:\n", .{});
        const building_types = std.enums.values(BuildingType);
        for (building_types, 0..) |building_type, i| {
            try stdout.print("{d}. {s}\n", .{ i + 1, @tagName(building_type) });
        }
        
        try stdout.print("Select a building: ", .{});
        var input: [10]u8 = undefined;
        const len = try stdin.read(&input);
        const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
            try stdout.print("Invalid input\n", .{});
            return;
        };
        
        if (choice < 1 or choice > building_types.len) {
            try stdout.print("Invalid choice\n", .{});
            return;
        }
        
        try city.production_queue.append(.{ .Building = building_types[choice - 1] });
        try stdout.print("🏗️ Started building {s}\n", .{@tagName(building_types[choice - 1])});
    }

    // View production queue
    fn viewProductionQueue(self: *Game, city: *City) !void {
        const stdout = io.getStdOut().writer();
        
        try stdout.print("\n🏭 Production Queue:\n", .{});
        if (city.production_queue.items.len == 0) {
            try stdout.print("Empty\n", .{});
            return;
        }
        
        for (city.production_queue.items, 0..) |item, i| {
            switch (item) {
                .Unit => |unit| try stdout.print("{d}. {s} (Unit)\n", .{ i + 1, @tagName(unit) }),
                .Building => |building| try stdout.print("{d}. {s} (Building)\n", .{ i + 1, @tagName(building) }),
            }
        }
    }

    // Move unit
    fn moveUnit(self: *Game) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        if (player.units.items.len == 0) {
            try stdout.print("You have no units!\n", .{});
            return;
        }
        
        try stdout.print("\n🚶 Your Units:\n", .{});
        for (player.units.items, 0..) |unit, i| {
            try stdout.print("{d}. {s} at ({d},{d})\n", .{ i + 1, @tagName(unit.type), unit.x, unit.y });
        }
        
        try stdout.print("Select a unit: ", .{});
        var input: [10]u8 = undefined;
        const len = try stdin.read(&input);
        const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
            try stdout.print("Invalid input\n", .{});
            return;
        };
        
        if (choice < 1 or choice > player.units.items.len) {
            try stdout.print("Invalid choice\n", .{});
            return;
        }
        
        const unit = player.units.items[choice - 1];
        
        try stdout.print("Enter direction (dx dy): ", .{});
        const len2 = try stdin.read(&input);
        var tokens = mem.tokenize(u8, input[0..len2], " ");
        const dx_str = tokens.next() orelse "";
        const dy_str = tokens.next() orelse "";
        
        const dx = std.fmt.parseInt(i8, dx_str, 10) catch {
            try stdout.print("Invalid dx\n", .{});
            return;
        };
        const dy = std.fmt.parseInt(i8, dy_str, 10) catch {
            try stdout.print("Invalid dy\n", .{});
            return;
        };
        
        const new_x = @intCast(u8, (@as(i16, unit.x) + dx) % MAP_WIDTH);
        const new_y = @intCast(u8, (@as(i16, unit.y) + dy) % MAP_HEIGHT);
        
        // Check if move is valid
        const tile = self.map[new_y][new_x];
        if (tile.terrain == .Ocean or tile.terrain == .Mountains) {
            try stdout.print("Cannot move to that terrain\n", .{});
            return;
        }
        
        // Clear old position
        self.map[unit.y][unit.x].unit = null;
        
        // Set new position
        unit.x = new_x;
        unit.y = new_y;
        self.map[new_y][new_x].unit = unit;
        
        try stdout.print("🚶 Moved unit to ({d},{d})\n", .{ new_x, new_y });
    }

    // Found a new city
    fn foundCity(self: *Game) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        // Find a settler
        var settler: ?*Unit = null;
        for (player.units.items) |unit| {
            if (unit.type == .Settler) {
                settler = unit;
                break;
            }
        }
        
        if (settler == null) {
            try stdout.print("You have no settler units!\n", .{});
            return;
        }
        
        const tile = self.map[settler.?.y][settler.?.x];
        if (tile.city != null) {
            try stdout.print("There is already a city here!\n", .{});
            return;
        }
        
        try stdout.print("Enter name for new city: ", .{});
        var input: [50]u8 = undefined;
        const len = try stdin.read(&input);
        const name = input[0..len];
        
        // Create new city
        var city = try self.allocator.create(City);
        city.* = .{
            .name = try self.allocator.dupe(u8, name),
            .population = 1,
            .production = 0,
            .food = 0,
            .owner_id = player.cities.items[0].owner_id,
            .buildings = std.ArrayList(BuildingType).init(self.allocator),
            .production_queue = std.ArrayList(union(enum) { Unit: UnitType, Building: BuildingType }).init(self.allocator),
        };
        try player.cities.append(city);
        self.map[settler.?.y][settler.?.x].city = city;
        self.map[settler.?.y][settler.?.x].owner_id = player.cities.items[0].owner_id;
        
        // Remove settler unit
        self.map[settler.?.y][settler.?.x].unit = null;
        for (player.units.items, 0..) |unit, i| {
            if (unit == settler.?) {
                _ = player.units.orderedRemove(i);
                break;
            }
        }
        
        try stdout.print("🏙️ Founded new city: {s}\n", .{name});
    }

    // Research technology
    fn researchTech(self: *Game) !void {
        const stdin = io.getStdIn().reader();
        const stdout = io.getStdOut().writer();
        const player = &self.players.items[self.current_player_index];
        
        try stdout.print("\n🔬 Available Technologies:\n", .{});
        const tech_types = std.enums.values(TechType);
        for (tech_types, 0..) |tech_type, i| {
            if (!mem.containsAtLeast(TechType, player.techs.items, 1, &.{tech_type})) {
                try stdout.print("{d}. {s}\n", .{ i + 1, @tagName(tech_type) });
            }
        }
        
        try stdout.print("Select a technology: ", .{});
        var input: [10]u8 = undefined;
        const len = try stdin.read(&input);
        const choice = std.fmt.parseInt(u8, input[0..len], 10) catch {
            try stdout.print("Invalid input\n", .{});
            return;
        };
        
        if (choice < 1 or choice > tech_types.len) {
            try stdout.print("Invalid choice\n", .{});
            return;
        }
        
        const tech = tech_types[choice - 1];
        player.researching = tech;
        try stdout.print("🔬 Started researching {s}\n", .{@tagName(tech)});
    }

    // Calculate player score
    fn calculateScore(self: *Game, player_id: u8) u32 {
        const player = &self.players.items[player_id];
        var score: u32 = 0;
        
        // City points
        score += @intCast(u32, player.cities.items.len) * 100;
        
        // Population points
        for (player.cities.items) |city| {
            score += city.population * 50;
        }
        
        // Tech points
        score += @intCast(u32, player.techs.items.len) * 50;
        
        // Territory points
        for (self.map) |row| {
            for (row) |tile| {
                if (tile.owner_id == player_id) {
                    score += 5;
                }
            }
        }
        
        return score;
    }
};

pub fn main() !void {
    var arena = heap.ArenaAllocator.init(heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    
    const stdout = io.getStdOut().writer();
    const stdin = io.getStdIn().reader();
    
    try stdout.print("Enter number of players (2-8): ", .{});
    var input: [10]u8 = undefined;
    const len = try stdin.read(&input);
    const player_count = std.fmt.parseInt(u8, input[0..len], 10) catch {
        try stdout.print("Invalid input. Using default 4 players.\n", .{});
        return;
    };
    
    if (player_count < 2 or player_count > 8) {
        try stdout.print("Invalid number of players. Using default 4 players.\n", .{});
        return;
    }
    
    var game = try Game.init(allocator, player_count);
    try game.run();
}
