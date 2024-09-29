package commands

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GabrielDCelery/personal-notes-cli/internal/editor"
	"github.com/go-playground/validator/v10"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

type CreteNoteConfig struct {
	Editor       string `mapstructure:"editor" validate:"required"`
	Author       string `mapstructure:"author" validate:"required"`
	InboxDir     string `mapstructure:"inboxDir" validate:"required"`
	TemplatePath string `mapstructure:"templatePath" validate:"required"`
}

func CreateNote(title string) {
	var createNoteConfig CreteNoteConfig
	if err := viper.Unmarshal(&createNoteConfig); err != nil {
		log.Fatalln(err)
	}
	if err := validator.New().Struct(createNoteConfig); err != nil {
		log.Fatalln(err)
	}
	now := time.Now().UTC()
	fileName := createFileNameFromTitle(title, now)
	fileDate := now.Format("2006-01-02T15-04-05Z")
	notePath := createNoteConfig.InboxDir + "/" + fileName
	fmt.Println("Will create a note with the following settings:")
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Date: %s\n", fileDate)
	fmt.Printf("Path: %s\n", notePath)
	prompt := promptui.Prompt{
		Label: "Do you want to continue? [yN]",
	}
	promptAnswer, err := prompt.Run()
	if err != nil {
		log.Fatalln(err)
	}
	if promptAnswer != "y" {
		fmt.Printf("Exitting...\n")
		return
	}
	templateAsBytes, err := os.ReadFile(createNoteConfig.TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}
	template := string(templateAsBytes)
	template = strings.ReplaceAll(template, "{{ title }}", title)
	template = strings.ReplaceAll(template, "{{ author }}", createNoteConfig.Author)
	template = strings.ReplaceAll(template, "{{ date }}", now.Format("2006-01-02T15-04-05Z"))
	if err := os.WriteFile(notePath, []byte(template), 0644); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Finished creating note: %s\n", notePath)
	if err := editor.OpenPathInEditor(notePath, createNoteConfig.Editor); err != nil {
		log.Fatalln(err)
	}
}

func createFileNameFromTitle(title string, createdAt time.Time) string {
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "-")
	title = strings.ToLower(title)
	title = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(title, "")
	formattedDate := createdAt.Format("20060102150405")
	return title + "-" + formattedDate + ".md"
}
