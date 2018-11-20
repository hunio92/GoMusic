package Database

import (
	"context"
	"fmt"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ConfigFireStore struct {
	ProjectID string
}

func newSessionStore(config ConfigFireStore) (*firestore.Client, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %v", err)
	}
	return client, nil
}

func (fs *Repo) AddSession(ctx context.Context, sessionID, email string) error {
	_, _, err := fs.Session.Collection("sessions").Add(ctx, map[string]interface{}{
		"sessionID": sessionID,
		"email": email,
	})
	if err != nil {
		return fmt.Errorf("failed to add field: %v", err)
	}

	return nil
}

func (fs *Repo) GetSessionEmail(ctx context.Context, sessionID string) (string, error) {
	it := fs.Session.Collection("sessions").
		Where("sessionID", "==", sessionID).
		Documents(ctx)
	defer it.Stop()

	doc, err := it.Next();
	if err == iterator.Done {
		return "", nil
	}
	type SessionData struct {
		SessionID string
		Email string
	}
	var s SessionData
	doc.DataTo(&s)
	fmt.Println("data: ", s)
	return s.Email, nil
}

func (fs *Repo) IsValidSession(ctx context.Context, sessionID string) bool {
	it := fs.Session.Collection("sessions").
		Where("sessionID", "==", sessionID).
		Documents(ctx)
	defer it.Stop()

	if _, err := it.Next(); err == iterator.Done {
		return false
	}
	return true
}