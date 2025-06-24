package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type Tracking struct {
	ID                bson.ObjectID     `bson:"_id,omitempty"` // id
	TrackingSettingID bson.ObjectID     `bson:"tracking_setting_id" json:"tracking_setting_id"`
	Url               string            `bson:"url" json:"url"` // LP page url to know lp
	EndUserID         string            `bson:"end_user_id" json:"end_user_id"`
	GeneratedFrom     string            `bson:"generated_from" json:"generated_from"`
	Metadata          map[string]string `bson:"metadata" json:"metadata"`
	BaseEntity        `bson:",inline"`
}
