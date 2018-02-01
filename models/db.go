package models

//<https://api.github.com/repositories/15773229/stargazers?clientid=3d032602cc3318f720bf&clientsecret=8e7f4a421ce4183bcc4ffb16991942c18cd0c9b6&page=2>; rel="next", <https://api.github.com/repositories/15773229/stargazers?clientid=3d032602cc3318f720bf&clientsecret=8e7f4a421ce4183bcc4ffb16991942c18cd0c9b6&page=63>; rel="last"

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tosone/GithubTraveler/models/logger"
	"github.com/tosone/logging"
)

var engine *gorm.DB

// Connect connect to the database
func Connect() {
	var err error
	var dialString string
	if viper.GetString("Database.Engine") == "sqlite3" {
		dialString = viper.GetString("Database.Path")
	} else {
		dialString = fmt.Sprintf("%s://%s@%s:%s:%s",
			viper.GetString("Database.Engine"),
			viper.GetString("Database.Username"),
			viper.GetString("Database.Password"),
			viper.GetString("Database.Host"),
			viper.GetString("Database.Port"),
		)
	}

	if engine, err = gorm.Open(viper.GetString("Database.Engine"), dialString); err != nil {
		logging.WithFields(logging.Fields{"engine": viper.GetString("Database.Engine"), "dialString": dialString}).Panic(err.Error())
	}
	if err = engine.AutoMigrate(new(Log), new(User), new(Repo)).Error; err != nil {
		logging.Panic(err.Error())
	}

	engine.LogMode(true)
	gormLogger := logger.Logger{}
	engine.SetLogger(gormLogger)
}
