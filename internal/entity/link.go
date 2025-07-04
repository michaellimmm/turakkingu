package entity

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
- with omitempty:
  - If ID is zero value (primitive.NilObjectID), it's excluded from BSON
  - Mongodb will generate a new _id

- without omitempty:
  - Zero value gets serialized as ObjectId("000000000000000000000000")
  - This causes duplicate key errors on subsequent inserts

Best Practices
- Always use omitempty with _id field to allow MongoDB auto-generation
*/
type Link struct {
	ID         bson.ObjectID `bson:"_id,omitempty"` // id
	Name       string        `bson:"name"`          // name
	TenantID   string        `bson:"tenant_id"`     // tenant id
	Url        string        `bson:"url"`           // original url
	ShortID    string        `bson:"short_id"`      // short id
	BaseEntity `bson:",inline"`
}

func (l *Link) SetShortID() error {
	if l.ShortID != "" {
		return nil
	}

	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	l.ShortID = id
	return nil
}

func (l *Link) SetCreatedAt() {
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now().UTC()
	}
}

func (l *Link) SetUpdatedAt() {
	l.UpdatedAt = time.Now().UTC()
}

func (l *Link) SoftDelete() {
	now := time.Now().UTC()
	l.DeletedAt = &now
	l.SetUpdatedAt()
}

func (f *Link) ConstructFixedUrl(baseUrl string) string {
	return fmt.Sprintf("%s/r/%s", baseUrl, f.ShortID)
}
