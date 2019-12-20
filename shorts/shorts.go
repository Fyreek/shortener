package shorts

import (
	"encoding/json"

	"github.com/fyreek/shortener/db"
	"github.com/fyreek/shortener/logging"
	"github.com/fyreek/shortener/security"
)

// Shorts is the internal data structure for the a shortening object
type Shorts struct {
	Short    string `json:"short" bson:"short"`
	ManageID string `json:"manageId" bson:"manageId"`
	URL      string `json:"url" bson:"url"`
	Visits   int    `bson:"visits"`
}

// Input defines the json input when a shortening request is received
type Input struct {
	URL      string `json:"url"`
	ManageID string `json:"manageId"`
}

// New creates a new shorts object from the provided url. length defines how long the short url is going to be
func New(url, manageID string, length int) *Shorts {
	if manageID == "" {
		manageID = security.GetUUIDString()
	}
	s := Shorts{
		Short:    security.GetRandomString(length),
		ManageID: manageID,
		URL:      url,
	}
	return &s
}

// Visit adds a new visit to the shorts object
func (s *Shorts) Visit(dBase db.Database) error {
	s.Visits = s.Visits + 1
	return s.Save(dBase)
}

// Save saves the current short to the database. It updates a existing one or creates a new one if it was not previously created
func (s *Shorts) Save(dBase db.Database) error {
	oldS := Shorts{}
	err := dBase.GetSingleEntry("shorts", "short", s.Short, &oldS)
	if err != nil {
		if err == db.ErrNoDocument {
			err = dBase.InsertSingleEntry("shorts", *s)
			if err != nil {
				logging.Log(logging.Failure, "Could not save document:", err)
				return err
			}
			return nil
		}
		return err
	}

	err = dBase.UpdateSingleEntry("shorts", "short", s.Short, *s)
	return err
}

// GetShort returns a single short by its short url
func GetShort(short string, dBase db.Database) (*Shorts, error) {
	s := Shorts{}
	err := dBase.GetSingleEntry("shorts", "short", short, &s)
	return &s, err
}

// GetShortsAll returns all shorts
func GetShortsAll(sort string, limit int, dBase db.Database) (*[]Shorts, error) {
	return getShorts(sort, "", "", limit, dBase)
}

// GetShortsForManageID returns all shorts beloging to a single manage id
func GetShortsForManageID(sort, manageID string, limit int, dBase db.Database) (*[]Shorts, error) {
	return getShorts(sort, "manageId", manageID, limit, dBase)
}

func getShorts(sort, column, value string, limit int, dBase db.Database) (*[]Shorts, error) {
	// Sort by most visited
	sortMap := make(map[string]interface{})
	if sort == "desc" {
		sortMap["visits"] = 1
	} else {
		sortMap["visits"] = -1
	}

	sSlice := make([]Shorts, 0)
	byteArraySlice, err := dBase.GetMultipleEntries("shorts", column, value, sortMap, limit)
	if err != nil {
		return &sSlice, err
	}
	for _, elem := range byteArraySlice {
		d := Shorts{}
		err := json.Unmarshal(elem, &d)
		if err != nil {
			return &sSlice, err
		}
		sSlice = append(sSlice, d)
	}
	return &sSlice, nil
}
