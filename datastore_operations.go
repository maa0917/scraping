package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"time"
)

const (
	LatestLiveKind = "LatestLive"
	LiveMasterKind = "LiveMaster"
)

func saveLatestLives(ctx context.Context, client *datastore.Client, lives []Live) error {
	query := datastore.NewQuery(LatestLiveKind).KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("failed to get all keys: %w", err)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		return fmt.Errorf("failed to delete existing live: %w", err)
	}

	latestLives := make([]*LatestLive, len(lives))
	for i, live := range lives {
		latestLives[i] = &LatestLive{Live: live, CreatedAt: time.Now()}
	}

	keys = make([]*datastore.Key, len(latestLives))
	for i := range keys {
		keys[i] = datastore.IncompleteKey(LatestLiveKind, nil)
	}

	if _, err = client.PutMulti(ctx, keys, latestLives); err != nil {
		return fmt.Errorf("failed to bulk insert to datastore: %w", err)
	}

	return nil
}

func saveLiveMaster(ctx context.Context, client *datastore.Client, lives []Live) error {
	var liveMasters []LiveMaster
	keys, err := client.GetAll(ctx, datastore.NewQuery(LiveMasterKind), &liveMasters)
	if err != nil {
		return fmt.Errorf("get all keys error: %w", err)
	}

	masterMap := make(map[string]*datastore.Key, len(liveMasters))
	for i, master := range liveMasters {
		masterMap[master.Live.URL] = keys[i]
	}

	var entities []LiveMaster
	var putKeys []*datastore.Key

	for _, live := range lives {
		if key, exists := masterMap[live.URL]; exists {
			// 更新の場合
			putKeys = append(putKeys, key)
			entities = append(entities, LiveMaster{
				Live:      live,
				UpdatedAt: time.Now(),
			})
		} else {
			// 新規の場合
			putKeys = append(putKeys, datastore.IncompleteKey(LiveMasterKind, nil))
			entities = append(entities, LiveMaster{
				Live:      live,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	_, err = client.PutMulti(ctx, putKeys, entities)
	if err != nil {
		return fmt.Errorf("failed to bulk upsert to datastore: %w", err)
	}

	return nil
}
