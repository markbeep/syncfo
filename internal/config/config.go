package config

import (
	"flag"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// -------------------------------
type Public struct {
	Database Database
	Logging  Logging
}

type Database struct {
	Path       string
	Migrations string
}

type Logging struct {
	Debug bool // optional
}

func (p Public) isComplete() bool {
	return p.Database.Path != "" && p.Database.Migrations != ""
}

func (p *Public) overwrite(o Public) {
	if o.Database.Path != "" {
		p.Database.Path = o.Database.Path
	}
	if o.Database.Migrations != "" {
		p.Database.Migrations = o.Database.Migrations
	}
	if o.Logging.Debug {
		p.Logging.Debug = o.Logging.Debug
	}
}

// -------------------------------
type Secret struct {
	Discord Discord
}

type Discord struct {
	Token     string
	TestGuild string `toml:"test_guild"` // optional
}

func (s Secret) isComplete() bool {
	return s.Discord.Token != ""
}

func (s *Secret) overwrite(o Secret) {
	if o.Discord.Token != "" {
		s.Discord.Token = o.Discord.Token
	}
	if o.Discord.TestGuild != "" {
		s.Discord.TestGuild = o.Discord.TestGuild
	}
}

var (
	publicPath = flagEnvString("public", "configs/public.toml", "Path to the public config file")
	secretPath = flagEnvString("secret", "configs/secret.toml", "Path to the secret config file")

	// public overwrites
	dbPath       = flagEnvString("database_path", "", "Overwrite path to the database")
	dbMigrations = flagEnvString("database_migrations", "", "Overwrite path to the database migrations")
	loggingDebug = flagEnvBool("logging_debug", false, "Overwrite logging debug flag")

	// secret overwrites
	discordToken     = flagEnvString("discord_token", "", "Overwrite Discord bot token")
	discordTestGuild = flagEnvString("discord_test_guild", "", "Overwrite Discord test guild ID")
)

func Parse() (Public, Secret, error) {
	flag.Parse()

	pub, err := loadPub()
	if err != nil {
		return Public{}, Secret{}, err
	}

	sec, err := loadSec()
	if err != nil {
		return Public{}, Secret{}, err
	}

	return pub, sec, nil
}

func loadPub() (Public, error) {
	overrides := Public{
		Database: Database{
			Path:       *dbPath,
			Migrations: *dbMigrations,
		},
		Logging: Logging{
			Debug: *loggingDebug,
		},
	}

	// allows all overrides to be set via command line flags, even if the config file is missing
	if overrides.isComplete() {
		return overrides, nil
	}

	dat, err := os.ReadFile(*publicPath)
	if err != nil {
		return Public{}, err
	}
	var pub Public
	if err := toml.Unmarshal(dat, &pub); err != nil {
		return Public{}, err
	}

	pub.overwrite(overrides)
	return pub, nil
}

func loadSec() (Secret, error) {
	overrides := Secret{
		Discord: Discord{
			Token:     *discordToken,
			TestGuild: *discordTestGuild,
		},
	}

	// allows all overrides to be set via command line flags, even if the config file is missing
	if overrides.isComplete() {
		return overrides, nil
	}

	dat, err := os.ReadFile(*secretPath)
	if err != nil {
		return Secret{}, err
	}
	var sec Secret
	if err := toml.Unmarshal(dat, &sec); err != nil {
		return Secret{}, err
	}

	sec.overwrite(overrides)
	return sec, nil
}

func flagEnvString(key, fallback, description string) *string {
	if val, ok := os.LookupEnv(key); ok {
		return flag.String(key, val, description)
	}
	return flag.String(key, fallback, description)
}

func flagEnvBool(key string, fallback bool, description string) *bool {
	if val, ok := os.LookupEnv(key); ok {
		if strings.ToLower(val) == "true" || val == "1" {
			return flag.Bool(key, true, description)
		}
		return flag.Bool(key, false, description)
	}
	return flag.Bool(key, fallback, description)
}
