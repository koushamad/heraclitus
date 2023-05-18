package application

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var Config = func() *viper.Viper {
	v := viper.New()

	v.SetConfigName(getConfigName())

	v.AddConfigPath("./config")

	v.AutomaticEnv()

	v.SetConfigType("yml")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	return v
}()

func getConfigName() string {
	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}

	execName := filepath.Base(exec)
	// Splitting the base name by "_" and taking the last word.
	words := strings.Split(execName, "_")
	lastWord := words[len(words)-1]

	return lastWord
}
