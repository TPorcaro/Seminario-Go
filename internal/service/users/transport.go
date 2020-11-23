package users



import (
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"strconv"
)
//HTTPService ...
type HTTPService interface {
	Register (*gin.Engine)
}
type httpService struct {
	endpoints []*endpoint
}
type endpoint struct {
	method string
	path string
	function gin.HandlerFunc
}
//NewHTTPTransport ...
func NewHTTPTransport(s Service) HTTPService {
	endpoints:= makeEndpoints(s)
	return httpService{endpoints}
}
func makeEndpoints(s Service) []*endpoint{
	list := []*endpoint{}
	list = append(list, &endpoint{
		method : "GET",
		path : "/users",
		function: getAll(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/users/register",
		function: registerUser(s),
	})
	list = append(list, &endpoint{
		method : "GET",
		path : "/users/:id",
		function: getByID(s),
	})
	list = append(list, &endpoint{
		method : "DELETE",
		path : "/users/:id",
		function: deleteByID(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/users/login",
		function: login(s),
	})
	list = append(list, &endpoint{
		method : "POST",
		path : "/users/changepassword/:id",
		function: changePassword(s),
	})
	return list
}
func changePassword(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		var userData User
		body:= c.Request.Body
		x, _ := ioutil.ReadAll(body)
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12)
		json.Unmarshal([]byte(x), &userData)
		if userData.Password != ""{
			res, err := s.ChangePassword(ID, userData.Password)
			c.JSON(http.StatusOK, gin.H{
				"response" : res,
				"error" : err,
			})
		}else{
			c.JSON(http.StatusConflict, gin.H{
				"error" : "Porfavor escriba una contraseña",
			})		
		}
	}
}
func login(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		body := c.Request.Body
		var userData User
		x, _ := ioutil.ReadAll(body)
		json.Unmarshal([]byte(x), &userData)
		token, err := s.Login(userData.Email, userData.Password)
		c.JSON(http.StatusOK, gin.H{
			"token" : token,
			"error" : err,
		})
	}
}
func getByID(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12) 
		user, err := s.GetByID(ID)
		c.JSON(http.StatusOK, gin.H{
			"user": user,
			"response" : err,
		})
	}
}
func deleteByID(s Service) gin.HandlerFunc{
	return func (c*gin.Context){
		ID, _ := strconv.ParseInt(c.Param("id"), 6, 12) 
		res, err := s.DeleteByID(ID)
		c.JSON(http.StatusOK, gin.H{
			"user": res,
			"response": err,
		})
	}
}
func registerUser(s Service) gin.HandlerFunc{
	return func(c *gin.Context){
		body := c.Request.Body
		x, _ := ioutil.ReadAll(body)
		var userData User
		json.Unmarshal([]byte(x), &userData)
		user:= User{0,userData.Name,userData.Email,userData.Password}
		if user.Name != "" && user.Email != "" && user.Password != ""{
			c.JSON(http.StatusCreated, gin.H{
				"users": s.RegisterUser(user),
			})
		} else{
			c.JSON(http.StatusConflict, gin.H{
				"Error" : "Porfavor ingrese todos los campos",
			})
		}
	}
}
func getAll(s Service) gin.HandlerFunc{
	return func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"users" : s.GetAll(),
			})
	}
}
func (s httpService) Register( r *gin.Engine){
	for _, e:= range s.endpoints {
		r.Handle(e.method, e.path, e.function)
	}
}