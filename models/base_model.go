package models

type BaseModel struct {
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at"`
}
