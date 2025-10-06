package init

import (
	"GoFlix/common/model/database"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		"root", "root", "127.0.0.1", "4000", "goflix",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	err = db.AutoMigrate(
		&database.User{},
		&database.Comment{},
		&database.CommentCount{},
		&database.Follower{},
		&database.Following{},
		&database.FollowerNums{},
		&database.FollowingNums{},
		&database.Like{},
		&database.LikeCount{},
	)
	if err != nil {
		panic(err.Error())
	}
}
