package platform

type Platform struct {
	Name        string   `bson:"name,omitempty"`
	Platform    string   `bson:"platform,omitempty"`
	Directories []string `bson:"platform,omitempty"`
}
