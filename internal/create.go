package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

func Create(title string) {
	now := time.Now().UTC()
	fileName := createFileNameFromTitle(title, now)
	fileDate := now.Format("2006-01-02T15-04-05Z")
	err, envVariables := getOsEnvVariablesForNoteCreation()
	if err != nil {
		log.Fatalln(err)
	}
	notePath := envVariables.PERSONAL_NOTES_INBOX_DIR + "/" + fileName
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
	template, err := readFileAsString(envVariables.PERSONAL_NOTES_TEMPLATE)
	if err != nil {
		log.Fatalln(err)
	}
	template = strings.ReplaceAll(template, "{{ title }}", title)
	template = strings.ReplaceAll(template, "{{ author }}", envVariables.PERSONAL_NOTES_DEFAULT_AUTHOR)
	template = strings.ReplaceAll(template, "{{ date }}", now.Format("2006-01-02T15-04-05Z"))
	err = writeStringToFile(notePath, template)
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

type OsEnvVariablesForNoteCreation struct {
	PERSONAL_NOTES_DEFAULT_AUTHOR string
	PERSONAL_NOTES_INBOX_DIR      string
	PERSONAL_NOTES_TEMPLATE       string
}

func getOsEnvVariablesForNoteCreation() (error, *OsEnvVariablesForNoteCreation) {
	PERSONAL_NOTES_DEFAULT_AUTHOR := os.Getenv("PERSONAL_NOTES_DEFAULT_AUTHOR")
	if len(PERSONAL_NOTES_DEFAULT_AUTHOR) == 0 {
		return fmt.Errorf("PERSONAL_NOTES_DEFAULT_AUTHOR has not been specified"), &OsEnvVariablesForNoteCreation{}
	}
	PERSONAL_NOTES_INBOX_DIR := os.Getenv("PERSONAL_NOTES_INBOX_DIR")
	if !isValidDirectory(PERSONAL_NOTES_INBOX_DIR) {
		return fmt.Errorf("PERSONAL_NOTES_INBOX_DIR is not a valid directory"), &OsEnvVariablesForNoteCreation{}
	}
	PERSONAL_NOTES_TEMPLATE := os.Getenv("PERSONAL_NOTES_TEMPLATE")
	if !isValidFile(PERSONAL_NOTES_TEMPLATE) {
		return fmt.Errorf("PERSONAL_NOTES_TEMPLATE is not a valid file"), &OsEnvVariablesForNoteCreation{}
	}
	return nil, &OsEnvVariablesForNoteCreation{
		PERSONAL_NOTES_DEFAULT_AUTHOR,
		PERSONAL_NOTES_INBOX_DIR,
		PERSONAL_NOTES_TEMPLATE,
	}
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
