package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var SETTINGS = map[string]any{}
var DB = ""
var dbMutex = sync.RWMutex{}
var TYPYING = make(map[string]int64)

const SETFILE = "./settings.json"

func main() {
	set, err := os.ReadFile(SETFILE)
	if err != nil {
		PrintRed(err.Error())
		return
	}
	err = json.Unmarshal(set, &SETTINGS)
	if err != nil {
		PrintRed(err.Error())
		return
	}
	//go executeBadger()
	fmt.Println("iniciando...")
	time.Sleep(3 * time.Second)
	r := GinRouter()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "search.html", nil)
	})

	//r.POST("/login")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/app", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.POST("/signin", login)

	api := r.Group("/api", AuthMiddleware(SETTINGS["jwt"].(string)))
	{
		api.POST("/sendTyping", sendTyping)
		api.POST("/verifyTyping", verifyTyping)
		api.POST("/getChats", getChats)
		api.POST("/validateSession", validateSession)
		api.POST("/getGroupsList", getGroupsList)
		api.POST("/addNewGroup", addNewGroup)
	}

	//BuildDatabase(SETTINGS["database"].(string), "sql.sql")
	insertAdmin()
	fmt.Println("server ir running...", SETTINGS["port"])
	r.Run("0.0.0.0:" + SETTINGS["port"].(string))
}

func addNewGroup(ctx *gin.Context) {
	var chat Chat
	err := ctx.ShouldBindJSON(&chat)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	chat.Members = map[string]string{}
	chat.DenuncedBy = map[string]string{}
	PrintData(chat)
	chat.City = strings.ToUpper(strings.TrimSpace(chat.City))
	chat.Country = strings.ToUpper(strings.TrimSpace(chat.Country))
	chat.Specialty = strings.ToUpper(strings.TrimSpace(chat.Specialty))
	dat, _ := QueryBadger("SELECT", "PROPS", "")
	dd := CategoriesEtc{}
	if dat.StatusCode == 404 {
		dd = CategoriesEtc{
			City:      map[string]string{chat.City: chat.City},
			Country:   map[string]string{chat.Country: chat.Country},
			Specialty: map[string]string{chat.Specialty: chat.Specialty},
		}
		QueryBadger("INSERT", "PROPS", dd)
	} else {
		json.Unmarshal(dat.ResultByte, &dd)
		_, ex := dd.City[chat.City]
		if !ex {
			dd.City[chat.City] = chat.City
		}
		_, ex = dd.Specialty[chat.Specialty]
		if !ex {
			dd.Specialty[chat.Specialty] = chat.Specialty
		}
		_, ex = dd.Country[chat.Country]
		if !ex {
			dd.Country[chat.Country] = chat.Country
		}
		QueryBadger("UPDATE", "PROPS", dd)
	}
	ctx.JSON(500, gin.H{"status": "error", "message": "ok", "data": ""})
}

func PrintData(data any) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(data)
		return
	}
	fmt.Println(string(jsonData))
}

func getGroupsList(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var chats map[string]MetaInfoChats
	dat, err := QueryBadger("SELECT", "CHATS_"+uuiduser.(string), "")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	if dat.StatusCode == 404 {
		chats = map[string]MetaInfoChats{}
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": chats})
		return
	}
	err = json.Unmarshal(dat.ResultByte, &chats)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": chats})
}

func validateSession(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func getChats(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	dat, err := QueryBadger("SELECT", "UserActivity_"+uuiduser.(string), "")
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	userActivity := UserActivity{}
	json.Unmarshal(dat.ResultByte, &userActivity)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": userActivity})
}

func sendTyping(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	dbMutex.Lock()
	TYPYING[uuiduser.(string)] = UnixTime() + 10
	dbMutex.Unlock()
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func verifyTyping(ctx *gin.Context) {
	//uuiduser, _ := ctx.Get("uuid")
	var username map[string]string
	err := ctx.ShouldBindJSON(&username)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": "NO"})
		return
	}
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	hora, ex := TYPYING[username["user"]]
	if !ex {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "NO"})
		return
	}
	if hora < UnixTime() {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "NO"})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "SI"})
}

func Prefix(prefix, key string) string {
	return prefix + "_" + key + "_" + Str(UnixMillisecTime())
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func login(ctx *gin.Context) {
	var datos map[string]string
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if datos["apodo"] == "" {
		fmt.Println("username is nill")
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	dats, err := QueryBadger("SELECT", datos["apodo"], "")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if dats.StatusCode == 404 {
		ctx.JSON(404, gin.H{"status": "error", "message": "username not found", "data": ""})
		return
	}
	var user User
	err = json.Unmarshal(dats.ResultByte, &user)
	if !ValidatePassword(datos["password"], user.Password) {
		ctx.JSON(401, gin.H{"status": "error", "message": "unauthorized", "data": ""})
		return
	}
	jwt, err := GenerateJWT(user.ID, SETTINGS["jwt"].(string), int(Int(SETTINGS["duration"])))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": jwt})
}

func insertAdmin() {
	dat, err := QueryBadger("CREATE", "sss", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(1 * time.Second)
	dat, err = QueryBadger("SELECT", "admin", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	if dat.StatusCode == 404 {
		uuid := uuid.New().String()
		uuid = "USER_" + uuid
		pass, _ := HashPassword("123456")
		user := User{
			Username:  "admin",
			Password:  pass,
			ID:        uuid,
			CreatedAt: DateTime(),
		}
		_, err := QueryBadger("INSERT", "admin", user)
		if err != nil {
			fmt.Println(err)
			return
		}
		QueryBadger("INSERT", uuid, user)
		info := UserInfo{
			UserType: "SUPERUSER",
			Apodo:    "webmaster",
		}
		QueryBadger("INSERT", "INFOUSER_"+uuid, info)
	}
}
