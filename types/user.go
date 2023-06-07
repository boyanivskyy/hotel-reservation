package types

// omitempty do not show empty id in JSON
type User struct {
	Id        string `bson:"_id" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}
