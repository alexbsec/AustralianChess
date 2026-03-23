package config

import (
	"os"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestSetEnvVars_SetsMissingValues(t *testing.T) {
	os.Clearenv()

	input := map[string]any{
		"db_host": "localhost",
		"db_port": 5432,
	}

	setEnvVars(input)

	require.Equal(t, "localhost", os.Getenv("DB_HOST"))
	require.Equal(t, "5432", os.Getenv("DB_PORT"))
}


func TestSetEnvVars_DoesNotOverrideExisting(t *testing.T) {
	os.Setenv("DB_HOST", "existing")

	input := map[string]any{
		"db_host": "new-value",
	}

	setEnvVars(input)

	require.Equal(t, "existing", os.Getenv("DB_HOST"))
}


func TestLoadMappedEnvVars_BindsEnv(t *testing.T) {
	v := &Environments{
		DBHost: "host",
		DBPort: 5432,
	}

	err := loadMappedEnvVars(v)

	require.NoError(t, err)

	require.NoError(t, viper.BindEnv("DB_HOST"))
}


func TestLoadMappedEnvVars_DecodeError(t *testing.T) {
	err := mapstructure.Decode(nil, nil)
	require.Error(t, err)
}


func TestLoadEnvs_Success(t *testing.T) {
	// backup + cleanup
	original := os.Getenv("DB_HOST")
	defer os.Setenv("DB_HOST", original)

	os.Unsetenv("DB_HOST")

	tmpFile := ".env"

	content := `DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=user
DB_PASSWORD=pass
DB_NAME=testdb`

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	defer os.Remove(tmpFile)

	env := LoadEnvs()

	require.Equal(t, "localhost", env.DBHost)
	require.Equal(t, 5432, env.DBPort)
	require.Equal(t, "user", env.DBUsername)
	require.Equal(t, "pass", env.DBPassword)
	require.Equal(t, "testdb", env.DBName)
}

func TestLoadEnvs_PanicsWhenNoConfig(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
	}()

	// force missing file
	viper.SetConfigFile("nonexistent.env")

	LoadEnvs()
}
