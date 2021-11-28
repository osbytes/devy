package envars

import (
	"os"
)

var (
	DISCORD_BOT_TOKEN string
)

func LoadEnvs() {
	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
}
