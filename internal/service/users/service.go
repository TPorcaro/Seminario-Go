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
	RegisterUser(User) string
	GetAll() []*User
	GetByID(int64) (*User, string)
	DeleteByID(int64) (string, error)
	Login(string, string) (uuid.UUID, error)
	ChangePassword(int64, string) (string, error)
}
type service struct {
	db   *sqlx.DB
	conf *config.Config
}

//New ...
func New(db *sqlx.DB, c *config.Config) (Service, error) {
	return service{db, c}, nil
}
func (s service) RegisterUser(u User) string {
	insertUser := `INSERT INTO users (name, email, password) VALUES(?,?,?)`
	data := s.db.MustExec(insertUser, u.Name, u.Email, u.Password)
	if data != nil {
		return "Se registro correctamente a " + u.Email
	}
	return "Error al registrarse"
}
func (s service) GetAll() []*User {
	var list []*User
	if err := s.db.Select(&list, "SELECT * FROM users"); err != nil {
		panic(err)
	}
	return list
}
func (s service) GetByID(ID int64) (*User, string) {
	var user User
	err := s.db.QueryRowx("SELECT * FROM users WHERE id = ?", ID).StructScan(&user)
	if err != nil {
		return nil, "Ese id no corresponde"
	}
	return &user, "Se trajo correctamente"
}
func (s service) DeleteByID(ID int64) (string,error) {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", ID)
	if err != nil {
		return "Error al borrar", err
	}
	return "Se borro correctamente", err
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
func (s service) ChangePassword(id int64, newPass string)(string, error){
	var user User
	err := s.db.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&user)
	if user.Name==""{
		return "No existe un usuario con ese mail", err
	}
	_, err = s.db.Exec("UPDATE users SET password = $1", newPass)
	if err != nil{
		return "Error al cambiar la contraseña", err
	}
	return "Su contraseña ha sido cambiada satisfactoriamente", nil
}
