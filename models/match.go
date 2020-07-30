/*

 */

package models

import (
	_ "github.com/satori/go.uuid"
)

type Match struct {
	TournamentID string `json:"tournamentId" gorm:"primary_key"`
	Round        int    `json:"round" gorm:"primarykey"`
	Table        int    `json:"table" gorm:"primarykey"`
	TeamOne      string `json:"teamOne"`
	TeamTwo      string `json:"teamTwo"`
	Status       string `json:"status"`
	Result       int    `json:"result"` // 1 if Team One wins, 2 if Team Two wins, -1 if no winner
}

func (db DB) CreateMatches(matches []Match) error {
	return db.DB.Create(matches).Error
}

func (db DB) GetMatch(tournamentId string, round, table int) (*Match, error) {
	match := Match{}
	err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (db DB) GetMatchesByTournament(tournamentId string) ([]Match, error) {
	matches := make([]Match, 0)
	err := db.DB.Order("round").Where("tournament_id = ?", tournamentId).Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (db DB) GetMatchesByStatus(status string) ([]Match, error) {
	matches := make([]Match, 0)
	err := db.DB.Where("status = ?", status).Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (db DB) DeleteMatch(tournamentId string, round, table int) error {
	match := Match{}
	err := db.DB.Where(`"tournament_id" = ? AND "round" = ? AND "table" = ?`, tournamentId, round, table).First(&match).Error
	if err != nil {
		return err
	}
	err = db.DB.Delete(&match).Error
	if err != nil {
		return err
	}
	return nil
}
func (db DB) InsertMatch(match Match) error {
	err := db.DB.Create(&match).Error
	if err != nil {
		return err
	}
	return nil
}

func (db DB) GetAllMatches() ([]Match, error) {
	matches := make([]Match, 0)
	err := db.DB.Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}
