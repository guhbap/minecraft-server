package subtypes

type entityType struct {
	Index  int32
	Width  float64
	Height float64
}

var EntityTypesRegistry = map[string]entityType{
	"minecraft:acacia_boat": {0, 1.375, 0.5625},
	// Acacia Boat 	1.375 	0.5625
	"minecraft:acacia_chest_boat": {1, 1.375, 0.5625},
	// Acacia Boat with Chest 	1.375 	0.5625
	"minecraft:allay": {2, 0.35, 0.6},
	// Allay 	0.35 	0.6
	"minecraft:area_effect_cloud": {3, 2.0, 0.5},
	// Area Effect Cloud 	2.0 * Radius 	0.5
	"minecraft:armadillo": {4, 0.7, 0.65},
	// Armadillo 	0.7 	0.65
	"minecraft:armor_stand": {5, 0.5, 1.975},
	// Armor Stand 	normal: 0.5 marker: 0.0 small: 0.25 	normal: 1.975 marker: 0.0 small: 0.9875
	"minecraft:arrow": {6, 0.5, 0.5},
	// Arrow 	0.5 	0.5
	"minecraft:axolotl": {7, 0.75, 0.42},
	// Axolotl 	0.75 	0.42
	"minecraft:bamboo_chest_raft": {8, 1.375, 0.5625},
	// Bamboo Raft with Chest 	1.375 	0.5625
	"minecraft:bamboo_raft": {9, 1.375, 0.5625},
	// Bamboo Raft 	1.375 	0.5625
	"minecraft:bat": {10, 0.5, 0.9},
	// Bat 	0.5 	0.9
	"minecraft:bee": {11, 0.7, 0.6},
	// Bee 	0.7 	0.6
	"minecraft:birch_boat": {12, 1.375, 0.5625},
	// Birch Boat 	1.375 	0.5625
	"minecraft:birch_chest_boat": {13, 1.375, 0.5625},
	// Birch Boat with Chest 	1.375 	0.5625
	"minecraft:blaze": {14, 0.6, 1.8},
	// Blaze 	0.6 	1.8
	"minecraft:block_display": {15, 0.0, 0.0},
	// Block Display 	0.0 	0.0
	"minecraft:bogged": {16, 0.6, 1.99},
	// Bogged 	0.6 	1.99
	"minecraft:breeze": {17, 0.6, 1.77},
	// Breeze 	0.6 	1.77
	"minecraft:breeze_wind_charge": {18, 0.3125, 0.3125},
	// Wind Charge 	0.3125 	0.3125
	"minecraft:camel": {19, 1.7, 2.375},
	// Camel 	1.7 	2.375
	"minecraft:cat": {20, 0.6, 0.7},
	// Cat 	0.6 	0.7
	"minecraft:cave_spider": {21, 0.7, 0.5},
	// Cave Spider 	0.7 	0.5
	"minecraft:cherry_boat": {22, 1.375, 0.5625},
	// Cherry Boat 	1.375 	0.5625
	"minecraft:cherry_chest_boat": {23, 1.375, 0.5625},
	// Cherry Boat with Chest 	1.375 	0.5625
	"minecraft:chest_minecart": {24, 0.98, 0.7},
	// Minecart with Chest 	0.98 	0.7
	"minecraft:chicken": {25, 0.4, 0.7},
	// Chicken 	0.4 	0.7
	"minecraft:cod": {26, 0.5, 0.3},
	// Cod 	0.5 	0.3
	"minecraft:command_block_minecart": {27, 0.98, 0.7},
	// Minecart with Command Block 	0.98 	0.7
	"minecraft:cow": {28, 0.9, 1.4},
	// Cow 	0.9 	1.4
	"minecraft:creaking": {29, 0.9, 2.7},
	// Creaking 	0.9 	2.7
	"minecraft:creeper": {30, 0.6, 1.7},
	// Creeper 	0.6 	1.7
	"minecraft:dark_oak_boat": {31, 1.375, 0.5625},
	// Dark Oak Boat 	1.375 	0.5625
	"minecraft:dark_oak_chest_boat": {32, 1.375, 0.5625},
	// Dark Oak Boat with Chest 	1.375 	0.5625
	"minecraft:dolphin": {33, 0.9, 0.6},
	// Dolphin 	0.9 	0.6
	"minecraft:donkey": {34, 1.3964844, 1.5},
	// Donkey 	1.3964844 	1.5
	"minecraft:dragon_fireball": {35, 1.0, 1.0},
	// Dragon Fireball 	1.0 	1.0
	"minecraft:drowned": {36, 0.6, 1.95},
	// Drowned 	0.6 	1.95
	"minecraft:egg": {37, 0.25, 0.25},
	// Thrown Egg 	0.25 	0.25
	"minecraft:elder_guardian": {38, 1.9975, 1.9975},
	// Elder Guardian 	1.9975 (2.35 * guardian) 	1.9975 (2.35 * guardian)
	"minecraft:enderman": {39, 0.6, 2.9},
	// Enderman 	0.6 	2.9
	"minecraft:endermite": {40, 0.4, 0.3},
	// Endermite 	0.4 	0.3
	"minecraft:ender_dragon": {41, 16.0, 8.0},
	// Ender Dragon 	16.0 	8.0
	"minecraft:ender_pearl": {42, 0.25, 0.25},
	// Thrown Ender Pearl 	0.25 	0.25
	"minecraft:end_crystal": {43, 2.0, 2.0},
	// End Crystal 	2.0 	2.0
	"minecraft:evoker": {44, 0.6, 1.95},
	// Evoker 	0.6 	1.95
	"minecraft:evoker_fangs": {45, 0.5, 0.8},
	// Evoker Fangs 	0.5 	0.8
	"minecraft:experience_bottle": {46, 0.25, 0.25},
	// Thrown Bottle o' Enchanting 	0.25 	0.25
	"minecraft:experience_orb": {47, 0.5, 0.5},
	// Experience Orb 	0.5 	0.5
	"minecraft:eye_of_ender": {48, 0.25, 0.25},
	// Eye of Ender 	0.25 	0.25
	"minecraft:falling_block": {49, 0.98, 0.98},
	// Falling Block 	0.98 	0.98
	"minecraft:fireball": {50, 1.0, 1.0},
	// Fireball 	1.0 	1.0
	"minecraft:firework_rocket": {51, 0.25, 0.25},
	// Firework Rocket 	0.25 	0.25
	"minecraft:fox": {52, 0.6, 0.7},
	// Fox 	0.6 	0.7
	"minecraft:frog": {53, 0.5, 0.5},
	// Frog 	0.5 	0.5
	"minecraft:furnace_minecart": {54, 0.98, 0.7},
	// Minecart with Furnace 	0.98 	0.7
	"minecraft:ghast": {55, 4.0, 4.0},
	// Ghast 	4.0 	4.0
	"minecraft:happy_ghast": {56, 4.0, 4.0},
	// Happy Ghast 	4.0 	4.0
	"minecraft:giant": {57, 3.6, 12.0},
	// Giant 	3.6 	12.0
	"minecraft:glow_item_frame": {58, 0.75, 0.75},
	// Glow Item Frame 	0.75 or 0.0625 (depth) 	0.75
	"minecraft:glow_squid": {59, 0.8, 0.8},
	// Glow Squid 	0.8 	0.8
	"minecraft:goat": {60, 1.3, 0.9},
	// Goat 	1.3 	0.9
	"minecraft:guardian": {61, 0.85, 0.85},
	// Guardian 	0.85 	0.85
	"minecraft:hoglin": {62, 1.3964844, 1.4},
	// Hoglin 	1.3964844 	1.4
	"minecraft:hopper_minecart": {63, 0.98, 0.7},
	// Minecart with Hopper 	0.98 	0.7
	"minecraft:horse": {64, 1.3964844, 1.6},
	// Horse 	1.3964844 	1.6
	"minecraft:husk": {65, 0.6, 1.95},
	// Husk 	0.6 	1.95
	"minecraft:illusioner": {66, 0.6, 1.95},
	// Illusioner 	0.6 	1.95
	"minecraft:interaction": {67, 0.0, 0.0},
	// Interaction 	0.0 	0.0
	"minecraft:iron_golem": {68, 1.4, 2.7},
	// Iron Golem 	1.4 	2.7
	"minecraft:item": {69, 0.25, 0.25},
	// Item 	0.25 	0.25
	"minecraft:item_display": {70, 0.0, 0.0},
	// Item Display 	0.0 	0.0
	"minecraft:item_frame": {71, 0.75, 0.75},
	// Item Frame 	0.75 or 0.0625 (depth) 	0.75
	"minecraft:jungle_boat": {72, 1.375, 0.5625},
	// Jungle Boat 	1.375 	0.5625
	"minecraft:jungle_chest_boat": {73, 1.375, 0.5625},
	// Jungle Boat with Chest 	1.375 	0.5625
	"minecraft:leash_knot": {74, 0.375, 0.5},
	// Leash Knot 	0.375 	0.5
	"minecraft:lightning_bolt": {75, 0.0, 0.0},
	// Lightning Bolt 	0.0 	0.0
	"minecraft:llama": {76, 0.9, 1.87},
	// Llama 	0.9 	1.87
	"minecraft:llama_spit": {77, 0.25, 0.25},
	// Llama Spit 	0.25 	0.25
	"minecraft:magma_cube": {78, 0.5202, 0.5202},
	// Magma Cube 	0.5202 * size 	0.5202 * size
	"minecraft:mangrove_boat": {79, 1.375, 0.5625},
	// Mangrove Boat 	1.375 	0.5625
	"minecraft:mangrove_chest_boat": {80, 1.375, 0.5625},
	// Mangrove Boat with Chest 	1.375 	0.5625
	"minecraft:marker": {81, 0.0, 0.0},
	// Marker 	0.0 	0.0
	"minecraft:minecart": {82, 0.98, 0.7},
	// Minecart 	0.98 	0.7
	"minecraft:mooshroom": {83, 0.9, 1.4},
	// Mooshroom 	0.9 	1.4
	"minecraft:mule": {84, 1.3964844, 1.6},
	// Mule 	1.3964844 	1.6
	"minecraft:oak_boat": {85, 1.375, 0.5625},
	// Oak Boat 	1.375 	0.5625
	"minecraft:oak_chest_boat": {86, 1.375, 0.5625},
	// Oak Boat with Chest 	1.375 	0.5625
	"minecraft:ocelot": {87, 0.6, 0.7},
	// Ocelot 	0.6 	0.7
	"minecraft:ominous_item_spawner": {88, 0.25, 0.25},
	// Ominous Item Spawner 	0.25 	0.25
	"minecraft:painting": {89, 0.0, 0.0},
	// Painting 	type width or 0.0625 (depth) 	type height
	"minecraft:pale_oak_boat": {90, 1.375, 0.5625},
	// Pale Oak Boat 	1.375 	0.5625
	"minecraft:pale_oak_chest_boat": {91, 1.375, 0.5625},
	// Pale Oak Boat with Chest 	1.375 	0.5625
	"minecraft:panda": {92, 1.3, 1.25},
	// Panda 	1.3 	1.25
	"minecraft:parrot": {93, 0.5, 0.9},
	// Parrot 	0.5 	0.9
	"minecraft:phantom": {94, 0.9, 0.5},
	// Phantom 	0.9 	0.5
	"minecraft:pig": {95, 0.9, 0.9},
	// Pig 	0.9 	0.9
	"minecraft:piglin": {96, 0.6, 1.95},
	// Piglin 	0.6 	1.95
	"minecraft:piglin_brute": {97, 0.6, 1.95},
	// Piglin Brute 	0.6 	1.95
	"minecraft:pillager": {98, 0.6, 1.95},
	// Pillager 	0.6 	1.95
	"minecraft:polar_bear": {99, 1.4, 1.4},
	// Polar Bear 	1.4 	1.4
	"minecraft:splash_potion": {100, 0.25, 0.25},
	// Splash Potion 	0.25 	0.25
	"minecraft:lingering_potion": {101, 0.25, 0.25},
	// Lingering Potion 	0.25 	0.25
	"minecraft:pufferfish": {102, 0.7, 0.7},
	// Pufferfish 	0.7 	0.7
	"minecraft:rabbit": {103, 0.4, 0.5},
	// Rabbit 	0.4 	0.5
	"minecraft:ravager": {104, 1.95, 2.2},
	// Ravager 	1.95 	2.2
	"minecraft:salmon": {105, 0.7, 0.4},
	// Salmon 	0.7 	0.4
	"minecraft:sheep": {106, 0.9, 1.3},
	// Sheep 	0.9 	1.3
	"minecraft:shulker": {107, 1.0, 1.0},
	// Shulker 	1.0 	1.0-2.0 (depending on peek)
	"minecraft:shulker_bullet": {108, 0.3125, 0.3125},
	// Shulker Bullet 	0.3125 	0.3125
	"minecraft:silverfish": {109, 0.4, 0.3},
	// Silverfish 	0.4 	0.3
	"minecraft:skeleton": {110, 0.6, 1.99},
	// Skeleton 	0.6 	1.99
	"minecraft:skeleton_horse": {111, 1.3964844, 1.6},
	// Skeleton Horse 	1.3964844 	1.6
	"minecraft:slime": {112, 0.5202, 0.5202},
	// Slime 	0.5202 * size 	0.5202 * size
	"minecraft:small_fireball": {113, 0.3125, 0.3125},
	// Small Fireball 	0.3125 	0.3125
	"minecraft:sniffer": {114, 1.9, 1.75},
	// Sniffer 	1.9 	1.75
	"minecraft:snowball": {115, 0.25, 0.25},
	// Snowball 	0.25 	0.25
	"minecraft:snow_golem": {116, 0.7, 1.9},
	// Snow Golem 	0.7 	1.9
	"minecraft:spawner_minecart": {117, 0.98, 0.7},
	// Minecart with Monster Spawner 	0.98 	0.7
	"minecraft:spectral_arrow": {118, 0.5, 0.5},
	// Spectral Arrow 	0.5 	0.5
	"minecraft:spider": {119, 1.4, 0.9},
	// Spider 	1.4 	0.9
	"minecraft:spruce_boat": {120, 1.375, 0.5625},
	// Spruce Boat 	1.375 	0.5625
	"minecraft:spruce_chest_boat": {121, 1.375, 0.5625},
	// Spruce Boat with Chest 	1.375 	0.5625
	"minecraft:squid": {122, 0.8, 0.8},
	// Squid 	0.8 	0.8
	"minecraft:stray": {123, 0.6, 1.99},
	// Stray 	0.6 	1.99
	"minecraft:strider": {124, 0.9, 1.7},
	// Strider 	0.9 	1.7
	"minecraft:tadpole": {125, 0.4, 0.3},
	// Tadpole 	0.4 	0.3
	"minecraft:text_display": {126, 0.0, 0.0},
	// Text Display 	0.0 	0.0
	"minecraft:tnt": {127, 0.98, 0.98},
	// Primed TNT 	0.98 	0.98
	"minecraft:tnt_minecart": {128, 0.98, 0.7},
	// Minecart with TNT 	0.98 	0.7
	"minecraft:trader_llama": {129, 0.9, 1.87},
	// Trader Llama 	0.9 	1.87
	"minecraft:trident": {130, 0.5, 0.5},
	// Trident 	0.5 	0.5
	"minecraft:tropical_fish": {131, 0.5, 0.4},
	// Tropical Fish 	0.5 	0.4
	"minecraft:turtle": {132, 1.2, 0.4},
	// Turtle 	1.2 	0.4
	"minecraft:vex": {133, 0.4, 0.8},
	// Vex 	0.4 	0.8
	"minecraft:villager": {134, 0.6, 1.95},
	// Villager 	0.6 	1.95
	"minecraft:vindicator": {135, 0.6, 1.95},
	// Vindicator 	0.6 	1.95
	"minecraft:wandering_trader": {136, 0.6, 1.95},
	// Wandering Trader 	0.6 	1.95
	"minecraft:warden": {137, 0.9, 2.9},
	// Warden 	0.9 	2.9
	"minecraft:wind_charge": {138, 0.3125, 0.3125},
	// Wind Charge 	0.3125 	0.3125
	"minecraft:witch": {139, 0.6, 1.95},
	// Witch 	0.6 	1.95
	"minecraft:wither": {140, 0.9, 3.5},
	// Wither 	0.9 	3.5
	"minecraft:wither_skeleton": {141, 0.7, 2.4},
	// Wither Skeleton 	0.7 	2.4
	"minecraft:wither_skull": {142, 0.3125, 0.3125},
	// Wither Skull 	0.3125 	0.3125
	"minecraft:wolf":             {143, 0.6, 0.85},
	"minecraft:zoglin":           {144, 1.3964844, 1.4},
	"minecraft:zombie":           {145, 0.6, 1.95},
	"minecraft:zombie_horse":     {146, 1.3964844, 1.6},
	"minecraft:zombie_villager":  {145, 0.6, 1.95},
	"minecraft:zombified_piglin": {146, 0.6, 1.95},
	"minecraft:player":           {147, 0.6, 1.8},
	"minecraft:fishing_bobber":   {150, 0.25, 0.25},
}
