package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
)

// EducationalPage contains the fields used for an educational (awareness) page.
// Unlike a landing Page, it never captures input: it is the content shown to a
// recipient after they interact with the landing page, so a campaign can
// redirect users to a chosen awareness message.
type EducationalPage struct {
	Id           int64     `json:"id" gorm:"column:id; primaryKey"`
	UserId       int64     `json:"-" gorm:"column:user_id"`
	Name         string    `json:"name"`
	HTML         string    `json:"html" gorm:"column:html"`
	ModifiedDate time.Time `json:"modified_date"`
}

// ErrEducationalPageNameNotSpecified is thrown if the name of the educational
// page is blank.
var ErrEducationalPageNameNotSpecified = errors.New("Educational page name not specified")

// TableName overrides the table name used by GORM.
func (p EducationalPage) TableName() string {
	return "educational_pages"
}

// Validate ensures that an educational page contains the appropriate details.
func (p *EducationalPage) Validate() error {
	if p.Name == "" {
		return ErrEducationalPageNameNotSpecified
	}
	return ValidateTemplate(p.HTML)
}

// GetEducationalPages returns the educational pages owned by the given user.
func GetEducationalPages(uid int64) ([]EducationalPage, error) {
	ps := []EducationalPage{}
	err := db.Where("user_id=?", uid).Find(&ps).Error
	if err != nil {
		log.Error(err)
		return ps, err
	}
	return ps, err
}

// GetEducationalPage returns the educational page, if it exists, specified by
// the given id and user_id.
func GetEducationalPage(id int64, uid int64) (EducationalPage, error) {
	p := EducationalPage{}
	err := db.Where("user_id=? and id=?", uid, id).First(&p).Error
	if err != nil {
		log.Error(err)
	}
	return p, err
}

// GetEducationalPageByName returns the educational page, if it exists,
// specified by the given name and user_id.
func GetEducationalPageByName(n string, uid int64) (EducationalPage, error) {
	p := EducationalPage{}
	err := db.Where("user_id=? and name=?", uid, n).First(&p).Error
	if err != nil {
		log.Error(err)
	}
	return p, err
}

// PostEducationalPage creates a new educational page in the database.
func PostEducationalPage(p *EducationalPage) error {
	if err := p.Validate(); err != nil {
		log.Error(err)
		return err
	}
	err := db.Save(p).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// PutEducationalPage edits an existing educational page in the database.
// Per the PUT Method RFC, it presumes all data for a page is provided.
func PutEducationalPage(p *EducationalPage) error {
	if err := p.Validate(); err != nil {
		return err
	}
	err := db.Where("id=?", p.Id).Save(p).Error
	if err != nil {
		log.Error(err)
	}
	return err
}

// DeleteEducationalPage deletes an existing educational page in the database.
// An error is returned if a page with the given user id and page id is not
// found.
func DeleteEducationalPage(id int64, uid int64) error {
	err := db.Where("user_id=?", uid).Delete(EducationalPage{Id: id}).Error
	if err != nil {
		log.Error(err)
	}
	return err
}
