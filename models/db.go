package models

//<https://api.github.com/repositories/15773229/stargazers?clientid=3d032602cc3318f720bf&clientsecret=8e7f4a421ce4183bcc4ffb16991942c18cd0c9b6&page=2>; rel="next", <https://api.github.com/repositories/15773229/stargazers?clientid=3d032602cc3318f720bf&clientsecret=8e7f4a421ce4183bcc4ffb16991942c18cd0c9b6&page=63>; rel="last"

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tosone/logging"
)

var engine *gorm.DB

// Connect connect to the database
func Connect() (err error) {
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
	var gormLogger Logger
	engine.SetLogger(gormLogger)
	return
}

// Logger default logger
type Logger struct{}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

var sqlRegexp = regexp.MustCompile(`\?`)
var numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)

func (logger Logger) Print(values ...interface{}) {
	var formattedValues []string
	var sql string
	if len(values) > 1 {
		if values[0] == "sql" {
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			var affected string
			if len(values) > 5 {
				affected = strconv.FormatInt(values[5].(int64), 10)
			}
			logging.WithFields(logging.Fields{
				"cost":     fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0),
				"trace":    values[1],
				"affected": affected,
			}).Debug(sql)
		}
	}
}
