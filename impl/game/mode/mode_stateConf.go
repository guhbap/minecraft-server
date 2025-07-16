package mode

import (
	"fmt"
	"sort"

	"github.com/golangmc/minecraft-server/apis/game"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/game/ents"
	"github.com/golangmc/minecraft-server/impl/game/registry/biome"
	"github.com/golangmc/minecraft-server/impl/prot/client"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

func HandleStateConfiguration(watcher util.Watcher, join chan base.PlayerAndConnection) {

	watcher.SubAs(func(packet *server.PacketICustomPayload, conn base.Connection) {
		fmt.Println("custom payload: ", packet.Channel, string(packet.Data))
		conn.SendPacket(&client.PacketOUpdateEnabledFeatures{})
		conn.SendPacket(&client.PacketOSelectKnownPacks{
			KnownPacks: []client.KnownPack{
				{
					Namespace: "minecraft:core",
					Id:        "1.21.4",
					Version:   "1.21.4",
				},
			},
		})
	})
	watcher.SubAs(func(packet *server.PacketISelectKnownPacks, conn base.Connection) {
		fmt.Println("select known packs: ", packet.KnownPacks)

		conn.SendPacket(&client.PacketORegistryData{
			Id: "minecraft:dimension_type",
			Entries: []client.RegistryEntry{
				{
					Id: "minecraft:overworld",
					Value: DimensionType{
						PiglinSafe:                  0,
						Natural:                     1,
						AmbientLight:                0.0,
						Infiniburn:                  "#minecraft:infiniburn_overworld",
						RespawnAnchorWorks:          0,
						HasSkylight:                 1,
						BedWorks:                    1,
						Effects:                     "minecraft:overworld",
						HasRaids:                    1,
						LogicalHeight:               384,
						CoordinateScale:             1.0,
						Ultrawarm:                   0,
						HasCeiling:                  0,
						MinY:                        -64,
						Height:                      384,
						MonsterSpawnLightLevel:      15,
						MonsterSpawnBlockLightLimit: 0,
					},
				},
			},
		})

		conn.SendPacket(&client.PacketORegistryData{
			Id: "minecraft:painting_variant",
			Entries: []client.RegistryEntry{
				{
					Id: "minecraft:backyard",
					Value: PaintingVariant{
						AssetId: "minecraft:backyard",
						Height:  2,
						Width:   2,
						Title:   `{"color": "gray", "translate": "painting.minecraft.skeleton.title"}`,
						Author:  `{"color": "gray", "translate": "painting.minecraft.skeleton.author"}`,
					},
				},
			},
		})

		conn.SendPacket(&client.PacketORegistryData{
			Id: "minecraft:wolf_variant",
			Entries: []client.RegistryEntry{
				{
					Id: "minecraft:wolf_ashen",
					Value: WolfVariant{
						WildTexture:  "minecraft:entity/wolf/wolf_ashen",
						TameTexture:  "minecraft:entity/wolf/wolf_ashen_tame",
						AngryTexture: "minecraft:entity/wolf/wolf_ashen_angry",
						Biomes:       "minecraft:snowy_taiga",
					},
				},
			},
		})

		conn.SendPacket(&client.PacketORegistryData{
			Id:      "minecraft:worldgen/biome",
			Entries: LoadBiomes(),
		})

		conn.SendPacket(&client.PacketORegistryData{
			Id: "minecraft:damage_type",
			Entries: CreateFakeDamageType(
				"in_fire",
				"campfire",
				"lightning_bolt",
				"on_fire",
				"lava",
				"hot_floor",
				"in_wall",
				"cramming",
				"drown",
				"starve",
				"cactus",
				"fall",
				"ender_pearl",
				"fly_into_wall",
				"out_of_world",
				"generic",
				"magic",
				"wither",
				"dragon_breath",
				"dry_out",
				"sweet_berry_bush",
				"freeze",
				"stalagmite",
				"falling_block",
				"falling_anvil",
				"falling_stalactite",
				"sting",
				"mob_attack",
				"mob_attack_no_aggro",
				"player_attack",
				"arrow",
				"trident",
				"mob_projectile",
				"spit",
				"wind_charge",
				"fireworks",
				"fireball",
				"unattributed_fireball",
				"wither_skull",
				"thrown",
				"indirect_magic",
				"thorns",
				"explosion",
				"player_explosion",
				"sonic_boom",
				"bad_respawn_point",
				"outside_border",
				"generic_kill",
				"mace_smash",
			),
		})

		conn.SendPacket(
			&client.PacketOUpdateTags{
				RawData: client.UpdateTagsDataRaw, // TODO: генерировать данные самостоятельно
			},
		)

		conn.SendPacket(&client.PacketOFinishConfiguration{})

	})

	watcher.SubAs(func(packet *server.PacketIClientInformation, conn base.Connection) {

		conn.SendPacket(&client.PacketOCustomPayload{
			Channel: "minecraft:brand",
			Data:    []byte("golangmc"),
		})

		// conn.SendPacket(&client.PacketOFinishConfiguration{})

		// conn.SetState(base.PLAY)
	})
	watcher.SubAs(func(packet *server.PacketIFinishConfiguration, conn base.Connection) {
		fmt.Println("finish configuration")
		conn.SetState(base.PLAY)
		join <- base.PlayerAndConnection{
			Player: ents.NewPlayer(&game.Profile{
				UUID: conn.Profile().UUID,
				Name: conn.Profile().Name,
			}, conn),
			Connection: conn,
		}
	})
}

func LoadBiomes() []client.RegistryEntry {

	entries := []client.RegistryEntry{}
	for _, bm := range biome.InvertedBiomes {
		// bm, ok := biome.InvertedBiomes["minecraft:"+id]
		// if !ok {
		// 	panic("biome not found: " + id)
		// }
		fmt.Println("send biome: ", bm.Name, bm)
		entries = append(entries, client.RegistryEntry{
			Id:    bm.Name,
			Value: bm,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Id < entries[j].Id
	})
	return entries
}

func CreateFakeDamageType(ids ...string) []client.RegistryEntry {
	dt := DamageType{
		MessageId:        "onFire",
		Scaling:          "always",
		Exhaustion:       0.1,
		DeathMessageType: "default",
	}

	entries := []client.RegistryEntry{}
	for _, id := range ids {
		entries = append(entries, client.RegistryEntry{
			Id:    "minecraft:" + id,
			Value: dt,
		})
	}
	return entries
}

func ptrInt32(i int32) *int32 {
	return &i
}

type PaintingVariant struct {
	AssetId string `nbt:"asset_id"`
	Height  int32  `nbt:"height"`
	Width   int32  `nbt:"width"`
	Title   string `nbt:"title"`
	Author  string `nbt:"author"`
}

type DimensionType struct {
	HasSkylight                 byte    `nbt:"has_skylight"`
	HasCeiling                  byte    `nbt:"has_ceiling"`
	Ultrawarm                   byte    `nbt:"ultrawarm"`
	Natural                     byte    `nbt:"natural"`
	CoordinateScale             float64 `nbt:"coordinate_scale"`
	BedWorks                    byte    `nbt:"bed_works"`
	RespawnAnchorWorks          byte    `nbt:"respawn_anchor_works"`
	MinY                        int32   `nbt:"min_y"`
	Height                      int32   `nbt:"height"`
	LogicalHeight               int32   `nbt:"logical_height"`
	Infiniburn                  string  `nbt:"infiniburn"`
	Effects                     string  `nbt:"effects"`
	AmbientLight                float32 `nbt:"ambient_light"`
	PiglinSafe                  byte    `nbt:"piglin_safe"`
	HasRaids                    byte    `nbt:"has_raids"`
	MonsterSpawnLightLevel      int32   `nbt:"monster_spawn_light_level"`
	MonsterSpawnBlockLightLimit int32   `nbt:"monster_spawn_block_light_limit"`
}

type WolfVariant struct {
	WildTexture  string `nbt:"wild_texture"`
	TameTexture  string `nbt:"tame_texture"`
	AngryTexture string `nbt:"angry_texture"`
	Biomes       string `nbt:"biomes"`
}

type DamageType struct {
	MessageId        string  `nbt:"message_id"`
	Scaling          string  `nbt:"scaling"`
	Exhaustion       float32 `nbt:"exhaustion"`
	DeathMessageType string  `nbt:"death_message_type"`
}
