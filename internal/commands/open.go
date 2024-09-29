package commands

import (
	"log"
	"os"

	"github.com/GabrielDCelery/personal-notes-cli/internal/editor"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type OpenNotesConfig struct {
	Editor   string `mapstructure:"editor" validate:"required"`
	NotesDir string `mapstructure:"notesDir" validate:"required"`
}

func OpenNotes() {
	var config OpenNotesConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	if err := validator.New().Struct(config); err != nil {
		log.Fatalln(err)
	}
	if err := os.Chdir(config.NotesDir); err != nil {
		log.Fatalln(err)
	}
	if err := editor.OpenPathInEditor(config.NotesDir, config.Editor); err != nil {
		log.Fatalln(err)
	}

}
