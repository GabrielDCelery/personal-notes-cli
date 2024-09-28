package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

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
	err := viper.Unmarshal(&createNoteConfig)
	if err != nil {
		log.Fatalln(err)
	}
	err = validator.New().Struct(createNoteConfig)
	if err != nil {
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
	template, err := readFileAsString(createNoteConfig.TemplatePath)
	if err != nil {
		log.Fatalln(err)
	}
	template = strings.ReplaceAll(template, "{{ title }}", title)
	template = strings.ReplaceAll(template, "{{ author }}", createNoteConfig.Author)
	template = strings.ReplaceAll(template, "{{ date }}", now.Format("2006-01-02T15-04-05Z"))
	err = writeStringToFile(notePath, template)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Finished creating note: %s\n", notePath)
	openNoteInEditorCommand := exec.Command(createNoteConfig.Editor, notePath)
	openNoteInEditorCommand.Stdin = os.Stdin
	openNoteInEditorCommand.Stdout = os.Stdout
	openNoteInEditorCommand.Stderr = os.Stderr
	err = openNoteInEditorCommand.Run()
	if err != nil {
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

func isValidDirectory(path string) bool {
	dirInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return dirInfo.IsDir()
}

func isValidFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFileAsString(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
