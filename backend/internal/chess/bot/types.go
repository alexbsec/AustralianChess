package bot

type Difficulty int

const (
	EasyDifficulty   Difficulty = 3
	MediumDifficulty Difficulty = 4
	HardDifficulty   Difficulty = 5
)



const (
	EasyBotId string = "EasyBot"
	MediumBotId string = "MediumBot"
	HardBotId string = "HardBot"
)

func DifficultyToBotId(diff Difficulty) string {
	switch diff {
	case EasyDifficulty:
		return EasyBotId
	case MediumDifficulty:
		return MediumBotId
	case HardDifficulty:
		return HardBotId
	}

	return "unknown"
}
