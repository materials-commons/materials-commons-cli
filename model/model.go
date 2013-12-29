package model

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"time"
)

type UserGroup struct {
	Id          string    `gorethink:"id,omitempty"`
	Owner       string    `gorethink:"owner"`
	Name        string    `gorethink:"name"`
	Description string    `gorethink:"description"`
	Birthtime   time.Time `gorethink:"birthtime"`
	MTime       time.Time `gorethink:"mtime"`
	Access      string    `gorethink:"access"`
	Users       []string  `gorethink:"users"`
}

func NewUserGroup(owner, name string) UserGroup {
	now := time.Now()
	return UserGroup{
		Owner:       owner,
		Name:        name,
		Description: name,
		Access:      "private",
		Birthtime:   now,
		MTime:       now,
	}
}

func MatchingUserGroups(query r.RqlTerm, session *r.Session) ([]UserGroup, error) {
	var results []UserGroup
	rows, err := query.Run(session)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var ug UserGroup
		rows.Scan(&ug)
		results = append(results, ug)
	}

	return results, nil
}

type DataFile struct {
	Id              string    `gorethink:"id,omitempty"`
	Name            string    `gorethink:"name"`
	Access          string    `gorethink:"access"`
	MarkedForReview bool      `gorethink:"marked_for_review"`
	Reviews         []string  `gorethink:"reviews"`
	Birthtime       time.Time `gorethink:"birthtime"`
	MTime           time.Time `gorethink:"mtime"`
	ATime           time.Time `gorethink:"atime"`
	Tags            []string  `gorethink:"tags"`
	MyTags          []string  `gorethink:"mytags"`
	Description     string    `gorethink:"description"`
	Notes           []string  `gorethink:"description"`
	Owner           string    `gorethink:"owner"`
	Process         string    `gorethink:"process"`
	Machine         string    `gorethink:"machine"`
	Checksum        string    `gorethink:"checksum"`
	Size            int64     `gorethink:"size"`
	Location        string    `gorethink:"location"`
	MediaType       string    `gorethink:"mediatype"`
	Conditions      []string  `gorethink:"conditions"`
	Text            string    `gorethink:"text"`
	MetaTags        []string  `gorethink:"metatags"`
	DataDirs        []string  `gorethink:"datadirs"`
	Parent          string    `gorethink:"parent"`
}

func NewDataFile(name, access, owner string) DataFile {
	now := time.Now()
	return DataFile{
		Name:        name,
		Access:      access,
		Owner:       owner,
		Description: name,
		Birthtime:   now,
		MTime:       now,
		ATime:       now,
	}
}

func GetDataFile(id string, session *r.Session) (DataFile, error) {
	var df DataFile
	result, err := r.Table("datafiles").Get(id).RunRow(session)
	switch {
	case err != nil:
		return df, err
	case result.IsNil():
		return df, fmt.Errorf("Unknown DataFile Id: %s", id)
	default:
		err := result.Scan(&df)
		return df, err
	}
}

type User struct {
	Id          string    `gorethink:"id,omitempty"`
	Name        string    `gorethink:"name"`
	Email       string    `gorethink:"email"`
	Fullname    string    `gorethink:"fullname"`
	Password    string    `gorethink:"password"`
	ApiKey      string    `gorethink:"apikey"`
	Birthtime   time.Time `gorethink:"birthtime"`
	MTime       time.Time `gorethink:"mtime"`
	Avatar      string    `gorethink:"avatar"`
	Description string    `gorethink:"description"`
	Affiliation string    `gorethink:"affiliation"`
	HomePage    string    `gorethink:"homepage"`
	Notes       []string  `gorethink:"notes"`
}

func NewUser(name, email, password, apikey string) User {
	now := time.Now()
	return User{
		Name:      name,
		Email:     email,
		Password:  password,
		ApiKey:    apikey,
		Birthtime: now,
		MTime:     now,
	}
}

func GetUser(id string, session *r.Session) (*User, error) {
	var u User
	if err := GetItem(id, "users", session, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetItem(id, table string, session *r.Session, obj interface{}) error {
	result, err := r.Table(table).Get(id).RunRow(session)
	switch {
	case err != nil:
		return err
	case result.IsNil():
		return fmt.Errorf("Unknown User Id: %s", id)
	default:
		err := result.Scan(obj)
		return err
	}
}
