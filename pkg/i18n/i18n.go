package i18n

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	i18nManager struct {
		bundle          *i18n.Bundle
		localizer       *i18n.Localizer
		defaultMessages map[string]defaultMessage
	}

	messages       map[string]defaultMessage
	defaultMessage struct {
		other    string
		template []string
	}
)

const (
	PT_MESSAGES_PATH = "../resources/pt-messages.json"
	ES_MESSAGES_PATH = "../resources/es-messages.json"
	FAV_LANGUAGE_ENV = "TOMASTER_FAV_LANG"
)

var Lines i18nManager

func init() {
	Lines.bundle = setupBundle(language.English, PT_MESSAGES_PATH, ES_MESSAGES_PATH)

	// we might be filling it manually later
	Lines.defaultMessages = map[string]defaultMessage{
		"ID": {
			other: "abc",
		},
	}

	lang := getLang()
	buildLocalizer(Lines.bundle, lang)
}

func setupBundle(defaultLang language.Tag, JSONPaths ...string) *i18n.Bundle {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, path := range JSONPaths {
		bundle.LoadMessageFile(path)
	}

	return bundle
}

func buildLocalizer(bundle *i18n.Bundle, langs ...string) {
	Lines.localizer = i18n.NewLocalizer(bundle, langs...)
}

func (i i18nManager) getConfig(ID string, templateParameters ...string) *i18n.LocalizeConfig {
	var templateData map[string]interface{}
	for index, value := range i.defaultMessages[ID].template {
		templateData[value] = templateParameters[index]
	}

	return &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    ID,
			Other: i.defaultMessages[ID].other,
		},
		TemplateData: templateData,
	}
}

func (i i18nManager) Find(ID string, templateParameters ...string) string {
	config := i.getConfig(ID, templateParameters...)
	return Lines.localizer.MustLocalize(config) // panic if fails
}

func getLang() string {
	lang := os.Getenv(FAV_LANGUAGE_ENV)
	if lang != "" {
		return lang
	}

	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		lang = getLangOnWindows()
	case "darwin":
		lang = getLangOnMacos()
	case "linux":
		lang = getLangOnLinux()
	}

	if lang != "" {
		return lang
	}
	return language.English.String()
}

func getLangOnWindows() string {
	// Exec powershell Get-Culture on Windows.
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		langLocRaw := strings.TrimSpace(string(output))
		langLoc := strings.Split(langLocRaw, "-")
		lang := langLoc[0]
		return lang
	}

	return ""
}
func getLangOnMacos() string {
	// Exec shell Get-Culture on Macos
	cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
	output, err := cmd.Output()
	if err == nil {
		langLocRaw := strings.TrimSpace(string(output))
		langLoc := strings.Split(langLocRaw, "_")
		lang := langLoc[0]
		return lang
	}

	return ""
}
func getLangOnLinux() string {
	envlang, ok := os.LookupEnv("LANG")
	if ok {
		langLocRaw := strings.TrimSpace(envlang)
		langLocRaw = strings.Split(envlang, ".")[0]
		langLoc := strings.Split(langLocRaw, "_")
		lang := langLoc[0]
		return lang
	}

	return ""
}
