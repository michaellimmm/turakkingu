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
	ID         bson.ObjectID     `bson:"_id,omitempty"` // id
	TenantID   string            `bson:"tenant_id"`     // tenant id
	Url        string            `bson:"url"`           // original url
	Metadata   map[string]string `bson:"metadata"`      // metadata
	ShortID    string            `bson:"short_id"`      // short id
	BaseEntity `bson:",inline"`
}

func (f *Link) SetShortID() error {
	if f.ShortID != "" {
		return nil
	}

	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	f.ShortID = id
	return nil
}

func (f *Link) SetCreatedAt() {
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
}

func (f *Link) SetUpdatedAt() {
	f.UpdatedAt = time.Now().UTC()
}

func (f *Link) SoftDelete() {
	now := time.Now().UTC()
	f.DeletedAt = &now
	f.SetUpdatedAt()
}

func (f *Link) ConstructFixedUrl(baseUrl string) string {
	return fmt.Sprintf("%s/r/%s", baseUrl, f.ShortID)
}
