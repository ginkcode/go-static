package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"
)

var sv *fiber.App

func main() {
	port := "3000"
	index, spa := false, false
	wg := &sync.WaitGroup{}
	app := &cli.App{
		Name:      "static-server",
		Usage:     "Start http server with static files!",
		ArgsUsage: "Path of static files (default: \".\")",
		Version:   "v1.0.0",
		Commands:  nil,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Value:       "3000",
				Destination: &port,
				Aliases:     []string{"p"},
				Usage:       "listening port of http server",
			},
			&cli.BoolFlag{
				Name:        "index",
				Destination: &index,
				Aliases:     []string{"i"},
				Usage:       "allow index all directory",
			},
			&cli.BoolFlag{
				Name:        "spa",
				Destination: &spa,
				Aliases:     []string{"s"},
				Usage:       "support SPA mode, this option will be ignored when allowing index",
			},
		},
		Action: func(c *cli.Context) error {
			wg.Add(1)
			colorReset := "\033[0m"
			colorRed := "\033[31m"
			colorBlue := "\033[36m"
			fmt.Println("Server is running at:", colorBlue+"http://localhost:"+port+colorReset)
			dir := c.Args().Get(0)
			if err := startServer(dir, ":"+port, index, spa); err != nil {
				fmt.Println(colorRed + err.Error() + colorReset)
				return err
			}
			return nil
		},
	}
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		if err := stopServer(ctx); err != nil {
			fmt.Println(err.Error())
		}
		cancel()
		wg.Done()
	}()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}

func startServer(dir, port string, index bool, spa bool) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if dir != "." && dir != "" {
		path = dir
	}
	sv = fiber.New(fiber.Config{
		Prefork:               false,
		DisableStartupMessage: true,
		ReduceMemoryUsage:     true,
	})
	sv.Static("/", path, fiber.Static{
		Browse: index,
	})
	if !index && spa {
		sv.Get("*", func(ctx *fiber.Ctx) error {
			if strings.Contains(string(ctx.Request().URI().Path()), ".") {
				return ctx.SendStatus(404)
			}
			f, e := os.Open(path + "/index.html")
			if e != nil {
				return ctx.SendStatus(404)
			}
			defer f.Close()
			ctx.Set("Content-Type", "text/html")
			return ctx.SendStream(f)
		})
	}
	return sv.Listen(port)
}

func stopServer(ctx context.Context) error {
	stop := make(chan struct{}, 1)
	go func(c chan struct{}) {
		if sv != nil {
			_ = sv.Shutdown()
		}
		c <- struct{}{}
	}(stop)
	fmt.Println("\nStopping http server...")
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout")
		case <-stop:
			return nil
		}
	}
}
