package entityMetadata

import (
	"github.com/golangmc/minecraft-server/apis/buff"
)

type EntityField struct {
	Index byte
	Type  byte
	Value any
}

func GetBaseFields() baseFields {
	return baseFields{
		EntityMetaBitMask:   EntityField{0, 0, 0},
		AirTicks:            EntityField{1, 1, 300},
		CustomName:          EntityField{2, 6, nil},
		IsCustomNameVisible: EntityField{3, 8, false},
		IsSilent:            EntityField{4, 8, false},
		HasNoGravity:        EntityField{5, 8, false},
		Pose:                EntityField{6, 21, 0},
		TicksFrozen:         EntityField{7, 1, 0},
	}
}

type baseFields struct {
	EntityMetaBitMask   EntityField
	AirTicks            EntityField
	CustomName          EntityField
	IsCustomNameVisible EntityField
	IsSilent            EntityField
	HasNoGravity        EntityField
	Pose                EntityField
	TicksFrozen         EntityField
}

func GetLivingEntityFields() livingEntityFields {
	return livingEntityFields{
		HandStates:                EntityField{8, 0, 0},
		Health:                    EntityField{9, 3, 1.0},
		PotionEffectColor:         EntityField{10, 18, 0},
		IsPotionEffectAmbient:     EntityField{11, 8, false},
		NumberOfArrowsInBody:      EntityField{12, 1, 0},
		NumberOfBeeStingersInBody: EntityField{13, 1, 0},
		LocationOfBed:             EntityField{14, 11, nil},
	}
}

type livingEntityFields struct {
	HandStates                EntityField
	Health                    EntityField
	PotionEffectColor         EntityField
	IsPotionEffectAmbient     EntityField
	NumberOfArrowsInBody      EntityField
	NumberOfBeeStingersInBody EntityField
	LocationOfBed             EntityField
}

func GetPlayerFields() playerFields {
	return playerFields{
		AdditionalHearts:                 EntityField{15, 3, 0},
		Score:                            EntityField{16, 1, 0},
		TheDisplayedSkinPartsBitMaskThat: EntityField{17, 0, 0},
		MainHand:                         EntityField{18, 0, 1},
		LeftShoulderEntityData:           EntityField{19, 16, nil},
		RightShoulderEntityData:          EntityField{20, 16, nil},
	}
}

type playerFields struct {
	AdditionalHearts                 EntityField
	Score                            EntityField
	TheDisplayedSkinPartsBitMaskThat EntityField
	MainHand                         EntityField
	LeftShoulderEntityData           EntityField
	RightShoulderEntityData          EntityField
}

// 0 	Byte 	Byte
// 1 	VarInt 	VarInt
// 2 	VarLong 	VarLong
// 3 	Float 	Float
// 4 	String 	String (32767)
// 5 	Text Component 	Text Component
// 6 	Optional Text Component 	(Boolean, Optional Text Component) 	Text Component is present if the Boolean is set to true.
// 7 	Slot 	Slot
// 8 	Boolean 	Boolean
// 9 	Rotations 	(Float, Float, Float) 	rotation on x, rotation on y, rotation on z
// 10 	Position 	Position
// 11 	Optional Position 	(Boolean, Optional Position) 	Position is present if the Boolean is set to true.
// 12 	Direction 	VarInt Enum 	Down = 0, Up = 1, North = 2, South = 3, West = 4, East = 5
// 13 	Optional Living Entity Reference 	(Boolean, Optional UUID) 	UUID is present if the Boolean is set to true.
// 14 	Block State 	VarInt 	An ID in the block state registry.
// 15 	Optional Block State 	VarInt 	0 for absent (air is unrepresentable); otherwise, an ID in the block state registry.
// 16 	NBT 	NBT
// 17 	Particle 	(VarInt, Varies) 	particle type (an ID in the minecraft:particle_type registry), particle data (See Particles.)
// 18 	Particles 	(VarInt, Array of (VarInt, Varies)) 	length-prefixed list of particle defintions (as above).
// 19 	Villager Data 	(VarInt, VarInt, VarInt) 	villager type, villager profession, level (See below.)
// 20 	Optional VarInt 	VarInt 	0 for absent; 1 + actual value otherwise. Used for entity IDs.
// 21 	Pose 	VarInt Enum 	STANDING = 0, FALL_FLYING = 1, SLEEPING = 2, SWIMMING = 3, SPIN_ATTACK = 4, SNEAKING = 5, LONG_JUMPING = 6, DYING = 7, CROAKING = 8, USING_TONGUE = 9, SITTING = 10, ROARING = 11, SNIFFING = 12, EMERGING = 13, DIGGING = 14, (1.21.3: SLIDING = 15, SHOOTING = 16, INHALING = 17)
// 22 	Cat Variant 	VarInt 	An ID in the minecraft:cat_variant registry.
// 23 	Cow Variant 	VarInt 	An ID in the minecraft:cow_variant registry.
// 24 	Wolf Variant 	VarInt 	An ID in the minecraft:wolf_variant registry.
// 25 	Wolf Sound Variant 	VarInt 	An ID in the minecraft:wolf_sound_variant registry.
// 26 	Frog Variant 	VarInt 	An ID in the minecraft:frog_variant registry.
// 27 	Pig Variant 	VarInt 	An ID in the minecraft:pig_variant registry.
// 28 	Chicken Variant 	VarInt 	An ID in the minecraft:chicken_variant registry.
// 29 	Optional Global Position 	(Boolean, Optional Identifier, Optional Position) 	dimension identifier, position; only if the Boolean is set to true.
// 30 	Painting Variant 	ID or Painting Variant 	An ID in the minecraft:painting_variant registry, or an inline definition.
// 31 	Sniffer State 	VarInt Enum 	IDLING = 0, FEELING_HAPPY = 1, SCENTING = 2, SNIFFING = 3, SEARCHING = 4, DIGGING = 5, RISING = 6
// 32 	Armadillo State 	VarInt Enum 	IDLE = 0, ROLLING = 1, SCARED = 2, UNROLLING = 3
// 33 	Vector3 	(Float, Float, Float) 	x, y, z
// 34 	Quaternion 	(Float, Float, Float, Float) 	x, y, z, w

func (f *EntityField) Push(buf buff.Buffer) {
	buf.PushByt(f.Index)
	buf.PushByt(f.Type)

	switch f.Type {
	case 0:
		buf.PushByt(f.Value.(byte))
	case 1:
		buf.PushVrI(f.Value.(int32))
	case 2:
		buf.PushVrL(f.Value.(int64))
	case 3:
		buf.PushF32(f.Value.(float32))
	case 4:
		buf.PushF64(f.Value.(float64))
	case 5:
		panic("TextComponent not implemented")
	case 6:
		if f.Value != nil {
			panic("Opt TextComponent not implemented")
		} else {
			buf.PushBit(false)
		}
	case 7:
		panic("Slot not implemented")
	case 8:
		buf.PushBit(f.Value.(bool))
	case 9:
		panic("Rotations not implemented")
	case 10:
		panic("Position not implemented")
	case 11:
		panic("Opt Position not implemented")
	case 12:
		panic("Direction not implemented")
	case 13:
		panic("Opt Living Entity Reference not implemented")
	case 14:
		panic("Block State not implemented")
	case 15:
		panic("Opt Block State not implemented")
	case 16:
		panic("NBT not implemented")
	case 17:
		panic("Particle not implemented")
	case 18:
		panic("Particles not implemented")
	case 19:
		panic("Villager Data not implemented")
	case 20:
		panic("Opt VarInt not implemented")
	case 21:
		buf.PushVrI(f.Value.(int32))
	case 22:
		panic("Cat Variant not implemented")
	case 23:
		panic("Cow Variant not implemented")
	case 24:
		panic("Wolf Variant not implemented")
	case 25:
		panic("Wolf Sound Variant not implemented")
	case 26:
		panic("Frog Variant not implemented")
	case 27:
		panic("Pig Variant not implemented")
	case 28:
		panic("Chicken Variant not implemented")
	case 29:
		panic("Opt Global Position not implemented")
	case 30:
		panic("Painting Variant not implemented")
	case 31:
		panic("Sniffer State not implemented")
	case 32:
		panic("Armadillo State not implemented")
	case 33:
		panic("Vector3 not implemented")
	case 34:
		panic("Quaternion not implemented")
	}
}
