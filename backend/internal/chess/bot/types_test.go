package bot_test

import (
	"testing"

	"github.com/alexbsec/AustralianChess/backend/internal/chess/bot"
	"github.com/stretchr/testify/require"
)

func TestDifficultyToBotId_Easy(t *testing.T) {
	id := bot.DifficultyToBotId(bot.EasyDifficulty)
	require.Equal(t, bot.EasyBotId, id)
}

func TestDifficultyToBotId_Medium(t *testing.T) {
	id := bot.DifficultyToBotId(bot.MediumDifficulty)
	require.Equal(t, bot.MediumBotId, id)
}

func TestDifficultyToBotId_Hard(t *testing.T) {
	id := bot.DifficultyToBotId(bot.HardDifficulty)
	require.Equal(t, bot.HardBotId, id)
}

func TestDifficultyToBotId_Unknown(t *testing.T) {
	id := bot.DifficultyToBotId(bot.Difficulty(99))
	require.Equal(t, "unknown", id)
}

func TestDifficultyConstants(t *testing.T) {
	require.Equal(t, bot.Difficulty(3), bot.EasyDifficulty)
	require.Equal(t, bot.Difficulty(4), bot.MediumDifficulty)
	require.Equal(t, bot.Difficulty(5), bot.HardDifficulty)
}
