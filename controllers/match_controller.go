package controllers

import (
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
)

type MatchController struct {
	log *logrus.Entry
	ms  *services.MatchService
}

func NewMatchController(log *logrus.Logger, ms *services.MatchService) *MatchController {
	return &MatchController{
		log: log.WithField("controller", "match"),
		ms:  ms,
	}
}

func (mc *MatchController) HandlePing(c *gin.Context) {
	mc.log.Info("handling ping")
	c.JSON(
		http.StatusOK,
		gin.H{
			"msg": "pong",
		},
	)
}

type FormGetMatchSchedule struct {
	Teams   []string `json:"teams"`
	Format  string   `json:"format"`
	Results []int    `json:"results"` //Add results
}

func (mc *MatchController) HandleGetMatchSchedule(c *gin.Context) {
	form := FormGetMatchSchedule{}

	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Info("Can not bind json")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	lenteams := len(form.Teams)
	if !(lenteams > 0 && lenteams&(lenteams-1) == 0) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "number of teams not a power of 2",
			},
		)
		return
	}
	db := mc.ms.GetDb()

	if form.Format != "" { //Generate initial table
		if form.Teams == nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"msg":   "failure",
					"error": "missing mandatory input parameter",
				},
			)
			return
		}

		brackets, err := services.GetMatchSchedule(form.Teams, form.Format)

		if err != nil {
			mc.log.Info("Get matches error")
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"msg":   "failure",
					"error": err.Error(),
				},
			)
			return
		}
		for _, match := range brackets {
			db.InsertMatch(match)
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"msg":  "success",
				"data": brackets,
			},
		)
		return

	}

	if form.Results == nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "missing mandatory input parameter",
			},
		)
		return
	} else { /*//Deal with results
		var FindReadymatches []models.Match
		var Readymatches []models.Match
		var Pendingmatches []models.Match
		var err error
		FindReadymatches, _ = db.GetMatchesByStatus("Ready")
		if len(FindReadymatches) == 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"msg":   "failure",
					"error": "no such match",
				},
			)
			return
		}
		Readymatches, err = db.GetMatchesByStatus("Ready")
		for i := 0; i < len(form.Results); i++ {

			Pendingmatches, err = db.GetMatchesByStatus("Pending")
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					gin.H{
						"msg":   "failure",
						"error": "get matches error",
					},
				)
				return
			}

			db.DB.Model(&Readymatches[i]).Update("Status", "Finished")
			db.DB.Model(&Readymatches[i]).Update("result", form.Results[i])

			if len(Pendingmatches) != 0 {
				if Pendingmatches[0].TeamOne == "Unknown" {
					if form.Results[i] == 1 {
						db.DB.Model(&Pendingmatches[0]).Update("TeamOne", Readymatches[i].TeamOne)
					} else {
						db.DB.Model(&Pendingmatches[0]).Update("TeamOne", Readymatches[i].TeamTwo)
					}
				} else {
					if form.Results[i] == 1 {
						db.DB.Model(&Pendingmatches[0]).Update("TeamTwo", Readymatches[i].TeamOne)
					} else {
						db.DB.Model(&Pendingmatches[0]).Update("TeamTwo", Readymatches[i].TeamTwo)
					}
				}
				for _, Pendingmatch := range Pendingmatches {
					if Pendingmatch.TeamOne != "Unknown" && Pendingmatch.TeamTwo != "Unknown" {
						db.DB.Model(&Pendingmatch).Update("Status", "Ready")
					}
				}
			}

		}

		brackets, err := db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23") //Get all matches
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"msg":   "failure",
					"error": "get matches error",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"msg":  "success",
				"data": brackets,
			},
		)*/
	}

}

type FormResults struct {
	TournamentID string `json:"tournamentId"`
	Table        int    `json:"table"`
	Round        int    `json:"round"`
	Result       int    `json:"result"` //Add results
}

func (mc *MatchController) HandleSingleResults(c *gin.Context) {
	form := FormResults{}
	var FindReadymatches []models.Match
	var Readymatches []models.Match
	var Pendingmatches []models.Match
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}

	db := mc.ms.GetDb()
	FindReadymatches, _ = db.GetMatchesByStatus("Ready")
	if len(FindReadymatches) == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "no such match",
			},
		)
		return
	}

	Readymatches, err := db.GetMatchesByStatus("Ready")

	Pendingmatches, err = db.GetMatchesByStatus("Pending")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}

	db.DB.Model(&Readymatches[0]).Update("Status", "Finished")
	db.DB.Model(&Readymatches[0]).Update("result", form.Result)

	if len(Pendingmatches) != 0 {
		if Pendingmatches[0].TeamOne == "Unknown" {
			if form.Result == 1 {
				db.DB.Model(&Pendingmatches[0]).Update("TeamOne", Readymatches[0].TeamOne)
			} else {
				db.DB.Model(&Pendingmatches[0]).Update("TeamOne", Readymatches[0].TeamTwo)
			}
		} else {
			if form.Result == 1 {
				db.DB.Model(&Pendingmatches[0]).Update("TeamTwo", Readymatches[0].TeamOne)
			} else {
				db.DB.Model(&Pendingmatches[0]).Update("TeamTwo", Readymatches[0].TeamTwo)
			}
		}
		for _, Pendingmatch := range Pendingmatches {
			if Pendingmatch.TeamOne != "Unknown" && Pendingmatch.TeamTwo != "Unknown" {
				db.DB.Model(&Pendingmatch).Update("Status", "Ready")
			}
		}
	}

	brackets, err := db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23") //Get all matches
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"msg":  "success",
			"data": brackets,
		},
	)
}
func ResultsToRank(results []int) int {
	var rank int
	for _, result := range results {
		rank *= 2
		rank += result
	}
	return rank + 1
}

func (mc *MatchController) HandleConsolationResults(c *gin.Context) {
	form := FormResults{}
	var ThisMatch *models.Match
	db := mc.ms.GetDb()
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Info("Can not bind json")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}
	brackets, err := db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23") //Get all matches
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}
	var maxtable int //get max table
	for _, match := range brackets {
		if match.Table > maxtable {
			maxtable = match.Table
		}

	}
	if form.Table <= 0 || form.Table > maxtable {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "wrong table",
			},
		)
		return
	}

	ThisMatch, err = db.GetMatch(form.TournamentID, form.Round, form.Table)
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}
	if ThisMatch.Result != 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "the result have already been updated",
			},
		)
		return
	}

	db.DB.Model(ThisMatch).Update("Result", form.Result)
	db.DB.Model(ThisMatch).Update("Status", "finished")
	GroupSize := maxtable / int(math.Pow(float64(2), float64(form.Round-1)))
	HalfSize := GroupSize / 2
	Group := (form.Table-1)/GroupSize + 1
	StartIndex := (Group-1)*GroupSize + 1
	BehindIndex := StartIndex + HalfSize
	PendingMatches, err := db.GetMatchesByStatus("Pending")
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}

	for _, match := range PendingMatches { //Winner
		if match.Table >= StartIndex && match.Table < BehindIndex {
			if match.TeamOne == "Unknown" {
				if form.Result == 1 {
					db.DB.Model(&match).Update("TeamOne", ThisMatch.TeamOne)
					break
				} else {
					db.DB.Model(&match).Update("TeamOne", ThisMatch.TeamTwo)
					break
				}
			} else {
				if form.Result == 1 {
					db.DB.Model(&match).Update("TeamTwo", ThisMatch.TeamOne)
					break
				} else {
					db.DB.Model(&match).Update("TeamTwo", ThisMatch.TeamTwo)
					break
				}
			}
			break
		}
	}

	PendingMatches, err = db.GetMatchesByStatus("Pending")
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}

	for _, Pendingmatch := range PendingMatches {
		if Pendingmatch.TeamOne != "Unknown" && Pendingmatch.TeamTwo != "Unknown" {
			db.DB.Model(&Pendingmatch).Update("Status", "Ready")
		}
	}

	for _, match := range PendingMatches { //Loser
		if match.Table >= BehindIndex && match.Table < BehindIndex+HalfSize {
			if match.TeamOne == "Unknown" {
				if form.Result == 1 {
					db.DB.Model(&match).Update("TeamOne", ThisMatch.TeamTwo)
					break
				} else {
					db.DB.Model(&match).Update("TeamOne", ThisMatch.TeamOne)
					break
				}
			} else {
				if form.Result == 1 {
					db.DB.Model(&match).Update("TeamTwo", ThisMatch.TeamTwo)
					break
				} else {
					db.DB.Model(&match).Update("TeamTwo", ThisMatch.TeamOne)
					break
				}
			}
			break
		}
	}

	PendingMatches, err = db.GetMatchesByStatus("Pending")
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}

	for _, Pendingmatch := range PendingMatches {
		if Pendingmatch.TeamOne != "Unknown" && Pendingmatch.TeamTwo != "Unknown" {
			db.DB.Model(&Pendingmatch).Update("Status", "Ready")
		}
	}

	brackets, err = db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23") //Get all matches
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"msg":  "success",
			"data": brackets,
		},
	)

}

type CommandForm struct {
	Cmd string `json:"command"` //Add results
}

func (mc *MatchController) HandleGetConsolationRank(c *gin.Context) {
	form := CommandForm{}
	db := mc.ms.GetDb()
	if err := c.ShouldBindJSON(&form); err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": err.Error(),
			},
		)
		return
	}
	if form.Cmd != "Get rank" {
		mc.log.Info("Unmounted command")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "invalid command",
			},
		)
		return
	}
	var teams []string
	matches, err := db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23")
	if err != nil {
		mc.log.Info("Get matches error")
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "failure",
				"error": "get matches error",
			},
		)
		return
	}
	for _, match := range matches {
		if match.Result == 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"msg":   "failure",
					"error": "the tournament is In progress",
				},
			)
			return
		}
	}
	var maxtable int //get max table
	for _, match := range matches {
		if match.Table > maxtable {
			maxtable = match.Table
		}
	}
	for i := 0; i < maxtable; i++ {
		teams = append(teams, matches[i].TeamOne)
		teams = append(teams, matches[i].TeamTwo)
	}
	var results []int
	var FoundMatches []models.Match
	var ranks []int

	for _, team := range teams {
		FoundMatches = []models.Match{}
		results = []int{}
		for _, match := range matches {
			if match.TeamOne == team || match.TeamTwo == team {
				FoundMatches = append(FoundMatches, match)
			}
		}
		for _, match := range FoundMatches {
			if match.TeamOne == team && match.Result == 1 || match.TeamTwo == team && match.Result == 2 {
				results = append(results, 0)
			} else {
				results = append(results, 1)
			}
		}
		ranks = append(ranks, ResultsToRank(results))
	}
	c.JSON(
		http.StatusBadRequest,
		gin.H{
			"msg":   "success",
			"teams": teams,
			"ranks": ranks,
		},
	)
	return

}
