package user

import (
	"db"

	"models/comment"
	"models/picture"

	"github.com/jinzhu/gorm"
)

// User account
type User struct {
	gorm.Model

	Nickname     string
	Email        string
	Password     string
	Pictures     []picture.Picture `gorm:"ForeignKey:UserID"`
	CommentsLeft []comment.Comment `gorm:"ForeignKey:AuthorID"`
}

func init() {
	db.Get().AutoMigrate(&User{})
}

// AddUser adds new user to the database
func (usr *User) AddUser() error {
	return db.Get().Create(usr).Error
}

// Код ниже уязвим к SQL injection и так-то не использует ORM

// func addUser(u User) {
// 	exec(fmt.Sprintf(`INSERT INTO users (
// 		phone,
// 		email,
// 		name,
// 		surname,
// 		gender,
// 		birth_date,
// 		registration_date,
// 		created_at,
// 		updated_at
// 	) values (
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s',
// 		'%s'
// 	)`,
// 		u.Phone,
// 		u.Email,
// 		u.Name,
// 		u.Surname,
// 		u.Gender,
// 		formatTime(u.BirthDate),
// 		formatTime(u.RegistrationDate),
// 		formatTime(u.CreatedAt),
// 		formatTime(u.UpdatedAt),
// 	))
// }
