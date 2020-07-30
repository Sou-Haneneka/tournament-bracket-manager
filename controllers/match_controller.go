package controllers

import (
	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/bitspawngg/tournament-bracket-manager/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	Results []int    `json:"results"`
}

func (mc *MatchController) HandleGetMatchSchedule(c *gin.Context) {
	form := FormGetMatchSchedule{}
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
	} else { //Deal with results
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

		brackets, err := db.GetMatchesByTournament("4f3d9be9-226f-47f0-94f4-399c163fcd23")
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

}
