package data

type MinecraftVersion int

const (
	MC1_12_2 MinecraftVersion = iota
	MC1_13_2
	MC1_14_4
	MC1_15_2
	MC1_21_4
)

var CurrentProtocol = MC1_21_4

var protocolVersion = map[MinecraftVersion]int{
	MC1_12_2: 340,
	MC1_13_2: 404,
	MC1_14_4: 498,
	MC1_15_2: 578,
	MC1_21_4: 769,
}

func (m MinecraftVersion) Protocol() int {
	return protocolVersion[m]
}

func (m MinecraftVersion) String() string {
	switch m {
	case MC1_12_2:
		return "1.12.2"
	case MC1_13_2:
		return "1.13.2"
	case MC1_14_4:
		return "1.14.4"
	case MC1_15_2:
		return "1.15.2"
	case MC1_21_4:
		return "1.21.4"
	default:
		return "Unknown"
	}
}
