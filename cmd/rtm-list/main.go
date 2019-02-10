package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/directionless/go-rememberthemilk/rememberthemilk"
	"github.com/go-kit/kit/log/level"
	"github.com/kolide/kit/logutil"
	"github.com/peterbourgon/ff"
)

func main() {

	fs := flag.NewFlagSet("my-program", flag.ExitOnError)
	var (
		apiKey    = fs.String("api-key", "", "API Key")
		apiSecret = fs.String("api-secret", "", "API Shared Secret")

		debug = fs.Bool("debug", false, "log debug information")
		_     = fs.String("config", "", "config file (optional)")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("RTM"),
	)

	logger := logutil.NewCLILogger(*debug)

	rtm, err := rememberthemilk.New()
	if err != nil {
		level.Error(logger).Log(
			"msg", "Failed to create RTM client",
			"err", err,
		)
		os.Exit(1)
	}

	// Need to load before SetAuth, so we don't overwrite the existing creds
	if err := rtm.LoadAuth(); err != nil {
		level.Error(logger).Log(
			"msg", "Failed to create load RTM creds",
			"err", err,
		)
		os.Exit(1)
	}

	if *apiKey != "" && *apiSecret != "" {
		if err := rtm.SetAuth(*apiKey, *apiSecret); err != nil {
			level.Error(logger).Log(
				"msg", "Failed to create set auth",
				"err", err,
			)
			os.Exit(1)
		}
	}

	if err := rtm.EnsureAuth(); err != nil {
		level.Error(logger).Log(
			"msg", "Unable to login. No saved creds, and/or invalid creds supplied",
			"err", err,
		)
		os.Exit(1)
	}

	lists, err := rtm.GetList()
	if err != nil {
		level.Error(logger).Log(
			"msg", "Unable to get lists",
			"err", err,
		)
		os.Exit(1)
	}

	testTasks := &rememberthemilk.TasklistResponse{}
	for _, list := range lists {
		if list.Name == "test list" {
			listId := fmt.Sprintf("%d", list.ID)
			if err := rtm.Req("rtm.tasks.getList", testTasks, rememberthemilk.Param("list_id", listId)); err != nil {
				level.Error(logger).Log(
					"msg", "Unable to get tasks",
					"err", err,
				)
				os.Exit(1)

			}
		}

	}

	spew.Dump(testTasks)

}
