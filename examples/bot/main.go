package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/auroradevllc/valourgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	token := os.Getenv("AUTH_TOKEN")

	if token == "" {
		log.Fatal("No token provided")
	}

	c, err := valourgo.NewClient(token)

	if err != nil {
		log.WithError(err).Fatal("Unable to create client")
	}

	// Close will close all open connections to nodes
	defer c.Close()

	// Add event handlers
	c.AddHandler(messageCreateHandler)

	// Open main node connection, optional but recommended
	// This node will contain most planets.
	// If you don't do this, JoinAllChannels will do it for you
	if err := c.Open(context.Background()); err != nil {
		log.WithError(err).Fatal("Unable to connect to main node for realtime updates")
	}

	err = c.JoinAllChannels(context.Background())

	if err != nil {
		log.WithError(err).Fatal("Unable to join channels")
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	<-ch
}

func messageCreateHandler(e *valourgo.MessageCreateEvent) {
	log.Info("Message created in planet " + e.PlanetID.String() + " by " + e.AuthorID.String() + ": " + e.Content)
}
