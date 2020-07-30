package services

import (
	"os"
	"testing"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	"github.com/sirupsen/logrus"
)

func setupDB(t *testing.T) *models.DB {
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		t.Fatal("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
	if !exists {
		t.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path)
	if err := db.Connect(); err != nil {
		t.Fatal("db connection failed")
	}
	// db.DB.Migrator().DropTable(&models.Match{})
	db.DB.AutoMigrate(&models.Match{})
	return db
}

func setupLogger(t *testing.T) *logrus.Logger {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal("failed to open file for log")
	}

	logger := logrus.New()
	logger.Out = file
	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

/*func TestGetMatchSchedule(t *testing.T) {
	matchOne := models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 1, TeamOne: "A", TeamTwo: "B", Status: "Ready", Result: 0}
	matchTwo := models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 2, TeamOne: "C", TeamTwo: "D", Status: "Ready", Result: 0}
	matchThree := models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 2, Table: 1, TeamOne: "A", TeamTwo: "C", Status: "Pending", Result: 0}
	type args struct {
		teams  []string
		format string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Match
		wantErr bool
	}{
		{"1.1-happy path", args{[]string{"teamA", "teamB"}, "SINGLE"}, []models.Match{matchOne, matchTwo, matchThree}, false},
		{"1.2-wrong number of teams", args{[]string{"teamA", "teamB", "teamC"}, "SINGLE"}, nil, true},
		{"1.3-unsupported format", args{[]string{"teamA", "teamB", "teamC", "teamD"}, "CONSOLATION"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := services.GetMatchSchedule(tt.args.teams, tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchSchedule() got = %v, want %v", got, tt.want)
			}
		})
	}
}*/
