package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
)

func newDatastoreClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore client: %w", err)
	}
	return client, nil
}
