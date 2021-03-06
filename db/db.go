package db

import (
	"fmt"
	"github.com/eputnam/health-check-server/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type DataStore struct {
	DB *gorm.DB
}

type SurveyDB struct {
	gorm.Model
	Team        string
	ResponseURL string
	Active      bool
	EndTime     int64
}

func (SurveyDB) TableName() string {
	return "surveys"
}

type ResponseDB struct {
	gorm.Model
	SurveyID   uint
	QuestionID uint
	Answer     int
	UserID     string
}

func (ResponseDB) TableName() string {
	return "responses"
}

type QuestionDB struct {
	gorm.Model
	Text     string
	SurveyID uint
}

func (QuestionDB) TableName() string {
	return "questions"
}

type TeamDB struct {
	gorm.Model
	Name string
}

func (teamdb TeamDB) toLogString() string {
	return fmt.Sprintf("Name: %s, ID=%d", teamdb.Name, teamdb.ID)
}

func (TeamDB) TableName() string {
	return "teams"
}

func (ds *DataStore) SaveTeam(team TeamDB) TeamDB {
	if result := ds.DB.Create(&team); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debugf("Successfully saved team %s", team.toLogString())
	return team
}

func (ds *DataStore) GetTeams() []TeamDB {
	var teams []TeamDB
	if result := ds.DB.Find(&teams, &TeamDB{}); nil != result.Error {
		panic(result.Error)
	}
	logrus.Debug("Successfully got list of teams")
	return teams
}

func NewStore(conf config.GlobalConfig) (*DataStore, error) {
	dbConf := conf.DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger(conf)})
	logrus.Infof("Connected to database at %s:%s", dbConf.Host, dbConf.Port)
	if nil != err {
		return nil, err
	}

	if err := db.AutoMigrate(&SurveyDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&ResponseDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&QuestionDB{}); nil != err {
		return nil, err
	}
	if err := db.AutoMigrate(&TeamDB{}); nil != err {
		return nil, err
	}

	return &DataStore{DB: db}, nil
}

func newLogger(conf config.GlobalConfig) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{LogLevel: conf.GetDbLogLevel(), Colorful: true})
}
