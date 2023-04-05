package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antsrp/beanstalker/config"
	"github.com/antsrp/beanstalker/usecases/events"
	"github.com/antsrp/beanstalker/usecases/queue"

	"github.com/TwiN/go-color"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
)

func main() {
	svc := micro.NewService(
		micro.Name("beanstalk wrapper"),
		micro.Version("latest"),
	)

	app := svc.Options().Cmd.App()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "host",
			Aliases:  []string{"H"},
			Value:    "localhost",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "port",
			Aliases:  []string{"p", "P"},
			Value:    11300,
			Required: false,
		},
	}
	app.Action = func(*cli.Context) error {

		cfg, err := config.Load()

		if err != nil {
			return err
		}

		producer, err := queue.NewProducer(cfg.Host, cfg.Port, cfg.GetSettings())
		if err != nil {
			return err
		}

		listener := events.NewListener(producer, cfg)

		r := bufio.NewReader(os.Stdin)

		for {
			fmt.Println(cfg.Current() + "Command:")
			s, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return err
				}
			}
			res, err := listener.Parse(s)
			if err != nil {
				fmt.Println(color.Ize(color.Red, err.Error()))
			} else {
				if res != nil {
					fmt.Println(color.Ize(color.Cyan, *res))
				} else {
					fmt.Println(color.Ize(color.Cyan, "OK"))
				}
			}
		}
		producer.Close()
		return nil
	}

	if err := app.Run(os.Args[1:]); err != nil {
		log.Println("error while running: ", err.Error())
	}
}
