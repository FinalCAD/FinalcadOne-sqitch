package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/FinalCAD/FinalcadOne-sqitch/internal/configsqitch"
	"github.com/FinalCAD/FinalcadOne-sqitch/internal/utils"
)

func exitHandler() {
	var err error
	if e := recover(); e != nil {
		switch x := e.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			err = errors.New("unknown panic error")
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func initLog() {
	debug, _ := utils.GetenvBool("DEBUG", false)
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	if debug {
		lvl.Set(slog.LevelDebug)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})).With(slog.Bool("debug", debug))
	slog.SetDefault(logger)
}

func main() {
	defer exitHandler()
	initLog()
	slog.Debug("Starting config sqitch")
	configSqitch, err := configsqitch.GetConfig()
	if err != nil {
		panic(err.Error())
	}
	err = configsqitch.WriteConfig(configSqitch)
	if err != nil {
		panic(err.Error())
	}
	slog.Info("Configuration successfuly generated")
	slog.Debug("End config sqitch")
}
