package users

import (
	"entrega/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

//User ...
type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

//Service ...
type Service interface {
	RegisterUser(User) (int64,error)
	GetAll() []*User
	GetByID(int64) (*User, error)
	DeleteByID(int64) error
	Login(string, string) (uuid.UUID, error)
	ChangePassword(int64, string) (error)
}
type service struct {
	db   *sqlx.DB
	conf *config.Config
}

//New ...
func New(db *sqlx.DB, c *config.Config) (Service, error) {
	return service{db, c}, nil
}
// Return un int mayor a 0 si lo creo, si no un -1
func (s service) RegisterUser(u User) (int64,error) {
	insertUser := `INSERT INTO users (name, email, password) VALUES(?,?,?)`
	data, err := s.db.MustExec(insertUser, u.Name, u.Email, u.Password).LastInsertId()
	if data > 0 {
		return data, nil
	}
	return -1, err
}
func (s service) GetAll() []*User {
	var list []*User
	if err := s.db.Select(&list, "SELECT * FROM users"); err != nil {
		panic(err)
	}
	return list
}
func (s service) GetByID(ID int64) (*User, error) {
	var user User
	err := s.db.QueryRowx("SELECT * FROM users WHERE id = ?", ID).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s service) DeleteByID(ID int64) (error) {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}
func (s service) Login(email string, pass string) (uuid.UUID,error){
	var user User
	 err := s.db.QueryRowx("SELECT * FROM users WHERE email = $1 AND password = $2",email, pass).StructScan(&user)
	if user.Name == "" {
		return uuid.UUID{}, err
	}
	token := uuid.NewV4()
	return token, nil
}
func (s service) ChangePassword(id int64, newPass string)(error){
	var user User
	err := s.db.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&user)
	if user.Name==""{
		return  err
	}
	_, err = s.db.Exec("UPDATE users SET password = $1", newPass)
	if err != nil{
		return  err
	}
	return  nil
}
