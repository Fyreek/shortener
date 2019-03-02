package shorts

import "github.com/fyreek/shortener/random"

// Shorts is the internal data structure for the a shortening object
type Shorts struct {
	Short  string `json:"short" bson:"short"`
	URL    string `json:"url" bson:"url"`
	Visits int    `bson:"visits"`
}

// Input defines the json input when a shortening request is received
type Input struct {
	URL string `json:"url"`
}

// New creates a new shorts object from the provided url. length defines how long the short url is going to be
func New(url string, length int) *Shorts {
	s := Shorts{
		Short: random.String(length),
		URL:   url,
	}
	return &s
}

// Visit adds a new visit to the shorts object
func (s *Shorts) Visit() {
	s.Visits = s.Visits + 1
}
