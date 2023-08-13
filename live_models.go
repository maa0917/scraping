package main

import (
	"cloud.google.com/go/datastore"
	"time"
)

type Live struct {
	Band      string    `datastore:"Band"`
	Title     string    `datastore:"Title"`
	EventDate time.Time `datastore:"EventDate"`
	Venue     string    `datastore:"Venue"`
	URL       string    `datastore:"URL,noindex"`
}

type LatestLive struct {
	Live
	CreatedAt time.Time `datastore:"CreatedAt"`
}

type LiveMaster struct {
	Live
	ID        *datastore.Key `datastore:"__key__"`
	CreatedAt time.Time      `datastore:"CreatedAt"`
	UpdatedAt time.Time      `datastore:"UpdatedAt"`
}
