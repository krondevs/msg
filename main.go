package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
var RATELIMIT = make(map[string]int64)
var RATELIMITMILISECS = 10

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
		c.HTML(200, "index.html", map[string]string{"Group": ""})
	})

	//r.POST("/login")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/app", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", map[string]string{"Group": ""})
	})

	r.GET("/app/:token", func(ctx *gin.Context) {
		token := ctx.Param("token")
		fmt.Println(token, "token")
		if token == "" {
			ctx.HTML(200, "index.html", map[string]string{"Group": ""})
			return
		}
		dat, _ := QueryBadger("SELECT", token, "")
		if dat.StatusCode == 404 {
			fmt.Println("not found")
			ctx.HTML(200, "index.html", map[string]string{"Group": ""})
			return
		}
		var dato map[string]string
		json.Unmarshal(dat.ResultByte, &dato)
		dat, _ = QueryBadger("SELECT", dato["group"], "")
		if dat.StatusCode == 404 {
			fmt.Println("not found")
			ctx.HTML(200, "index.html", map[string]string{"Group": ""})
			return
		}
		var group Chat
		json.Unmarshal(dat.ResultByte, &group)
		QueryBadger("DELETE", token, "")
		link := uuid.New().String()
		group.Link = link
		QueryBadger("INSERT", link, map[string]string{"group": group.ID})
		QueryBadger("UPDATE", group.ID, group)
		ctx.HTML(200, "index.html", map[string]string{"Group": group.ID})
	})

	r.POST("/signin", login)

	r.POST("/register", register)

	r.GET("/join/:token", join)

	r.GET("/download/:filename", downloadFile)

	api := r.Group("/api", AuthMiddleware(SETTINGS["jwt"].(string)))
	{
		api.POST("/sendTyping", sendTyping)
		api.POST("/verifyTyping", verifyTyping)
		api.POST("/getChats", getChats)
		api.POST("/validateSession", validateSession)
		api.POST("/getGroupsList", getGroupsList)
		api.POST("/addNewGroup", addNewGroup)
		api.POST("/loadGroupChat", loadGroupChat)
		api.POST("/sendMessage", sendMessage)
		api.POST("/join", join)
		api.POST("/cerrarGrupo", cerrarGrupo)
		api.POST("/uploadFile", uploadFile)
		api.POST("/suspender", suspender)
		api.POST("/expulsar", expulsar)
		api.POST("/eliminarMsg", eliminarMsg)
		api.POST("/salirGrupo", salirGrupo)
	}

	//BuildDatabase(SETTINGS["database"].(string), "sql.sql")
	insertAdmin()
	fmt.Println("server ir running...", SETTINGS["port"])
	r.Run("0.0.0.0:" + SETTINGS["port"].(string))
}

func salirGrupo(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	det, _ := QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(det.ResultByte, &user)
	dat, _ := QueryBadger("SELECT", datos["group"], "")
	var grp Chat
	json.Unmarshal(dat.ResultByte, &grp)
	delete(grp.Members, uuiduser.(string))
	QueryBadger("UPDATE", grp.ID, grp)
	delete(user.Chats, grp.ID)
	UpdateUser(uuiduser.(string), user)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func eliminarMsg(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("SELECT", datos["group"], "")
	datosGrupo := Chat{}
	if dat.StatusCode == 404 {
		fmt.Println("no hay registros")
		ctx.JSON(404, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	err := json.Unmarshal(dat.ResultByte, &datosGrupo)
	if err != nil {
		fmt.Println("no hay registros")
		ctx.JSON(404, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	_, ex := datosGrupo.Members[uuiduser.(string)]
	if !ex {
		ctx.JSON(401, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	dat, _ = QueryBadger("SELECT", datos["msg"], "")
	if dat.StatusCode == 404 {
		fmt.Println("no hay registros")
		ctx.JSON(404, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	chat := Message{}
	json.Unmarshal(dat.ResultByte, &chat)
	chat.DeletedAt = DateTime()
	QueryBadger("UPDATE", datos["msg"], chat)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func expulsar(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	det, _ := QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(det.ResultByte, &user)

	dat, _ := QueryBadger("SELECT", datos["group"], "")
	var grp Chat
	json.Unmarshal(dat.ResultByte, &grp)
	_, ex := grp.Owners[uuiduser.(string)]
	if !ex {
		ctx.JSON(401, gin.H{"status": "error", "message": "unauthorized", "data": ""})
		return
	}
	delete(grp.Members, datos["user"])
	QueryBadger("UPDATE", grp.ID, grp)
	det, _ = QueryBadger("SELECT", uuiduser.(string), "")
	json.Unmarshal(det.ResultByte, &user)
	delete(user.Chats, grp.ID)
	UpdateUser(datos["user"], user)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func suspender(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	det, _ := QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(det.ResultByte, &user)

	dat, _ := QueryBadger("SELECT", datos["group"], "")
	var grp Chat
	json.Unmarshal(dat.ResultByte, &grp)
	_, ex := grp.Owners[uuiduser.(string)]
	if !ex {
		ctx.JSON(401, gin.H{"status": "error", "message": "unauthorized", "data": ""})
		return
	}
	if grp.Members[datos["user"]] == "SUSPENDED" {
		grp.Members[datos["user"]] = ""
		QueryBadger("UPDATE", grp.ID, grp)
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	grp.Members[datos["user"]] = "SUSPENDED"
	QueryBadger("UPDATE", grp.ID, grp)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	c.FileAttachment(filepath.Join("static", "uploads", filename), filename)
}

func uploadFile(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	gg := ctx.PostForm("current")
	fmt.Println(gg)
	file, err := ctx.FormFile("fileInput")
	fmt.Println(file.Filename, err)
	if err != nil {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	if file.Filename == "" {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	gg = strings.TrimSpace(gg)
	if gg == "" {
		fmt.Println(gg)
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	dat, _ := QueryBadger("SELECT", gg, "")
	var grp Chat
	json.Unmarshal(dat.ResultByte, &grp)
	_, ex := grp.Members[uuiduser.(string)]
	if !ex {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	if grp.Members[uuiduser.(string)] == "SUSPENDED" || grp.IsBlocked == true || grp.LockedByAdmin == true {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	nombreArchivo := uuid.New().String() + "_" + file.Filename
	nombreArchivo = strings.ToLower(nombreArchivo)
	descargable := nombreArchivo
	nn := strings.ToLower(file.Filename)
	dst := filepath.Join("static", "uploads", nombreArchivo)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error al guardar el archivo: ", "data": ""})
		return
	}
	det, _ := QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(det.ResultByte, &user)
	dat, _ = QueryBadger("SELECT", "MSGS_"+gg, "")
	hh := []string{}
	idmsg := "MSG_" + uuid.New().String()
	if filepath.Ext(nn) == ".jpg" || filepath.Ext(nn) == ".png" || filepath.Ext(nn) == ".jpeg" {
		nombreArchivo = `<img src="/` + dst + `" style="max-width: 80%;"><br><br>` + file.Filename
	} else if filepath.Ext(nn) == ".mp4" {
		nombreArchivo = `<video width="80%" height="auto" controls><source src="/` + dst + `" type="video/mp4">Tu navegador no soporta el elemento de video.</video><br><br>` + file.Filename

	} else if filepath.Ext(nn) == ".pdf" {
		nombreArchivo = `<a href="/` + dst + `" target="_blank">` + file.Filename + `</a><br><br>` + file.Filename
	} else {
		nombreArchivo = `<a href="/download/` + descargable + `">` + file.Filename + `</a><br><br>` + file.Filename
	}
	msg := Message{
		ID:        idmsg,
		GroupID:   gg,
		FromUser:  uuiduser.(string),
		CreatedAt: DateTime(),
		Text:      nombreArchivo,
		EditedAt:  "",
		DeletedAt: "",
		Status:    "OK",
		Filename:  "",
		MediaType: "",
		Apodo:     user.Apodo,
	}
	if dat.StatusCode == 404 {
		hh = append(hh, idmsg)
		QueryBadger("INSERT", "MSGS_"+gg, hh)
		QueryBadger("INSERT", idmsg, msg)
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	json.Unmarshal(dat.ResultByte, &hh)
	hh = append(hh, idmsg)
	QueryBadger("UPDATE", "MSGS_"+gg, hh)
	QueryBadger("INSERT", idmsg, msg)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func cerrarGrupo(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("SELECT", datos["group"], "")
	if dat.StatusCode == 404 {
		ctx.JSON(404, gin.H{"status": "error", "message": "grupo no existe", "data": ""})
		return
	}
	var group Chat
	json.Unmarshal(dat.ResultByte, &group)
	_, ex := group.Owners[uuiduser.(string)]
	if !ex {
		ctx.JSON(401, gin.H{"status": "error", "message": "no estas autorizado", "data": ""})
		return
	}
	if !group.IsBlocked {
		group.IsBlocked = true
	} else {
		group.IsBlocked = false
	}
	QueryBadger("UPDATE", group.ID, group)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func join(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	if datos["group"] == "" {
		ctx.JSON(400, gin.H{"status": "error", "message": "", "data": ""})
		return
	}
	dat, _ := QueryBadger("SELECT", datos["group"], "")
	if dat.StatusCode == 404 {
		ctx.JSON(404, gin.H{"status": "error", "message": "", "data": ""})
		return
	}
	var group Chat
	json.Unmarshal(dat.ResultByte, &group)
	_, ex := group.Members[uuiduser.(string)]
	if ex {
		ctx.JSON(401, gin.H{"status": "error", "message": "", "data": ""})
		return
	}
	group.Members[uuiduser.(string)] = ""
	dat, _ = QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(dat.ResultByte, &user)
	user.Chats[datos["group"]] = ""
	QueryBadger("UPDATE", group.ID, group)
	UpdateUser(uuiduser.(string), user)
	ctx.JSON(200, gin.H{"status": "success", "message": "Tes has unido al grupo de forma exitosa", "data": group.ID})
}

func register(c *gin.Context) {
	var form RegisterForm

	// Parse form data (campos normales)
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	// Manejo de archivos
	files := []struct {
		fieldName string
		jsonField *string
	}{
		{"foto_perfil", &form.FotoPerfil},
		{"selfie_doc", &form.SelfieDoc},
		{"foto_documento", &form.FotoDocumento},
		{"recibo_servicio", &form.ReciboServicio},
	}

	for _, fileConfig := range files {
		file, err := c.FormFile(fileConfig.fieldName)
		if err != nil {
			// Puedes manejar archivos opcionales o requeridos según tu lógica
			continue
		}

		// Guardar archivo
		nombreArchivo := uuid.New().String() + "_" + file.Filename
		dst := filepath.Join("static", "uploads", nombreArchivo)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error al guardar el archivo: " + fileConfig.fieldName, "data": ""})
			return
		}

		// Asignar nombre de archivo a la struct (para JSON)
		*fileConfig.jsonField = nombreArchivo
	}

	// Aquí puedes guardar form en base de datos o hacer lo que necesites
	//PrintData(form)
	form.Chats = map[string]string{}
	form.KycStatus = "PENDING"
	form.UserType = "USER"
	form.CreatedAt = DateTime()
	form.Password, _ = HashPassword(form.Password)
	newUser := "USER_" + uuid.New().String()
	form.ID = newUser
	_, err := InsertUser(form, newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Usuario registrado exitosamente",
		"data":    "",
	})
}

func UpdateUser(uuid string, form User) error {
	_, err := QueryBadger("UPDATE", form.Apodo, form)
	if err != nil {
		return err
	}
	_, err = QueryBadger("UPDATE", uuid, form)
	if err != nil {
		QueryBadger("DELETE", form.Apodo, "")
		return err
	}
	return err
}

func InsertUser(form RegisterForm, uuid string) (string, error) {
	form.Chats = map[string]string{}
	_, err := QueryBadger("INSERT", form.Apodo, form)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	_, err = QueryBadger("INSERT", uuid, form)

	return uuid, nil
}

func DetectarEnlacesCorreosTelefonos(texto string) bool {
	// Expresión regular para detectar enlaces (URLs)
	regexEnlaces := regexp.MustCompile(`https?://[^\s]+`)

	// Expresión regular para detectar correos electrónicos
	regexCorreos := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Expresión regular para detectar números de teléfono (varias formas comunes)
	regexTelefonos := regexp.MustCompile(`(?:\+?(\d{1,3}))?[-.\s]?$?(\d{1,4})$?[-.\s]?(\d{1,4})[-.\s]?(\d{1,9})`)

	// Verificar si hay coincidencias en el texto
	return regexEnlaces.MatchString(texto) ||
		regexCorreos.MatchString(texto) ||
		regexTelefonos.MatchString(texto)
}

func sendMessage(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	datos["msg"] = strings.TrimSpace(datos["msg"])
	if datos["msg"] == "" {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	enlace := DetectarEnlacesCorreosTelefonos(datos["msg"])
	if enlace {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	det, _ := QueryBadger("SELECT", uuiduser.(string), "")
	var user User
	json.Unmarshal(det.ResultByte, &user)

	dat, _ := QueryBadger("SELECT", datos["group"], "")
	var grp Chat
	json.Unmarshal(dat.ResultByte, &grp)
	_, ex := grp.Members[uuiduser.(string)]
	if !ex {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	if grp.Members[uuiduser.(string)] == "SUSPENDED" {
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	dat, _ = QueryBadger("SELECT", "MSGS_"+datos["group"], "")
	hh := []string{}
	idmsg := "MSG_" + uuid.New().String()
	msg := Message{
		ID:        idmsg,
		GroupID:   datos["group"],
		FromUser:  uuiduser.(string),
		CreatedAt: DateTime(),
		Text:      datos["msg"],
		EditedAt:  "",
		DeletedAt: "",
		Status:    "OK",
		Filename:  "",
		MediaType: "",
		Apodo:     user.Apodo,
	}
	if dat.StatusCode == 404 {
		hh = append(hh, idmsg)
		QueryBadger("INSERT", "MSGS_"+datos["group"], hh)
		QueryBadger("INSERT", idmsg, msg)
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	json.Unmarshal(dat.ResultByte, &hh)
	hh = append(hh, idmsg)
	QueryBadger("UPDATE", "MSGS_"+datos["group"], hh)
	QueryBadger("INSERT", idmsg, msg)

	grp.Cant += 1
	for k := range grp.Members {
		if k != uuiduser.(string) {
			if grp.Members[k] == "" {
				grp.Members[k] = "UNREAD"
			}
		}
	}
	QueryBadger("UPDATE", datos["group"], grp)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func loadGroupChat(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos MsgPagination
	ctx.ShouldBindJSON(&datos)
	//PrintData(datos)
	chats := []string{}
	msgs := []Message{}
	dat, _ := QueryBadger("SELECT", datos.UuidMsg, "")
	datosGrupo := Chat{}
	if dat.StatusCode == 404 {
		fmt.Println("no hay registros")
		ctx.JSON(200, gin.H{"status": "success", "message": "", "data": msgs, "me": uuiduser, "datosGrupo": Chat{}})
		return
	}
	err := json.Unmarshal(dat.ResultByte, &datosGrupo)
	if err != nil {
		fmt.Println("no hay registros")
		ctx.JSON(200, gin.H{"status": "success", "message": "", "data": msgs, "me": uuiduser, "datosGrupo": Chat{}})
		return
	}
	_, ex := datosGrupo.Members[uuiduser.(string)]
	if !ex {
		ctx.JSON(200, gin.H{"status": "success", "message": "", "data": msgs, "me": uuiduser, "datosGrupo": Chat{}})
		return
	}
	if datosGrupo.Members[uuiduser.(string)] == "UNREAD" {
		datosGrupo.Members[uuiduser.(string)] = ""
		QueryBadger("UPDATE", datos.UuidMsg, datosGrupo)
	}
	name := datosGrupo.Name
	members := datosGrupo.Members
	lok := datosGrupo.IsBlocked
	lockedByAdmin := datosGrupo.LockedByAdmin
	_, ex = datosGrupo.Owners[uuiduser.(string)]
	if !ex {
		datosGrupo = Chat{}
		datosGrupo.Owners = map[string]string{}
		datosGrupo.DenuncedBy = map[string]string{}
		datosGrupo.Name = name
		datosGrupo.Members = members
		datosGrupo.LockedByAdmin = lockedByAdmin
		datosGrupo.IsBlocked = lok

	}
	dat, _ = QueryBadger("SELECT", "MSGS_"+datos.UuidMsg, "")
	if dat.StatusCode == 404 {
		fmt.Println("no hay registros")
		ctx.JSON(200, gin.H{"status": "success", "message": "", "data": msgs, "me": uuiduser, "datosGrupo": datosGrupo})
		return
	}
	json.Unmarshal(dat.ResultByte, &chats)
	//PrintData(chats)
	chats = getLast100(chats, datos.Max)
	//PrintData(chats)

	chat := Message{}
	for _, i := range chats {
		dat, _ := QueryBadger("SELECT", i, "")
		json.Unmarshal(dat.ResultByte, &chat)
		if chat.DeletedAt == "" {
			msgs = append(msgs, chat)
		}
		/*if chat.FromUser != uuiduser.(string) && chat.ReSendts == 0 {
			chat.ReSendts = 1
			QueryBadger("UPDATE", i, chat)
			fmt.Println("read")
		}*/
	}

	ctx.JSON(200, gin.H{"status": "success", "message": "", "data": msgs, "me": uuiduser, "datosGrupo": datosGrupo})
}

func getLast100(slice []string, cont int) []string {
	if len(slice) <= cont {
		return slice
	}
	return slice[len(slice)-100:]
}

func addNewGroup(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var chat Chat
	err := ctx.ShouldBindJSON(&chat)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	chat.Members = map[string]string{uuiduser.(string): ""}
	chat.DenuncedBy = map[string]string{}
	chat.Owners = map[string]string{uuiduser.(string): ""}
	PrintData(chat)
	chat.City = strings.ToUpper(strings.TrimSpace(chat.City))
	chat.Country = strings.ToUpper(strings.TrimSpace(chat.Country))
	chat.Specialty = strings.ToUpper(strings.TrimSpace(chat.Specialty))
	dat, _ := QueryBadger("SELECT", "PROPS", "")
	dd := CategoriesEtc{}
	if dat.StatusCode == 404 {
		dd = CategoriesEtc{
			City:      map[string]string{chat.City: ""},
			Country:   map[string]string{chat.Country: ""},
			Specialty: map[string]string{chat.Specialty: ""},
		}
		QueryBadger("INSERT", "PROPS", dd)
	} else {
		json.Unmarshal(dat.ResultByte, &dd)
		_, ex := dd.City[chat.City]
		if !ex {
			dd.City[chat.City] = ""
		}
		_, ex = dd.Specialty[chat.Specialty]
		if !ex {
			dd.Specialty[chat.Specialty] = ""
		}
		_, ex = dd.Country[chat.Country]
		if !ex {
			dd.Country[chat.Country] = ""
		}
		QueryBadger("UPDATE", "PROPS", dd)
	}
	chttr := Str(UnixMillisecTime())
	chat.ID = "GROUP_" + chttr
	link := uuid.New().String()
	chat.Link = link
	QueryBadger("INSERT", "GROUP_"+chttr, chat)
	QueryBadger("INSERT", link, map[string]string{"group": "GROUP_" + chttr})
	var user User
	dat, _ = QueryBadger("SELECT", uuiduser.(string), "")
	json.Unmarshal(dat.ResultByte, &user)
	user.Chats["GROUP_"+chttr] = ""
	UpdateUser(uuiduser.(string), user)
	ctx.JSON(200, gin.H{"status": "error", "message": "ok", "data": ""})
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
	chats := map[string]Chat{}
	var user User
	dat, _ := QueryBadger("SELECT", uuiduser.(string), "")
	json.Unmarshal(dat.ResultByte, &user)
	for k := range user.Chats {
		dat, _ = QueryBadger("SELECT", k, "")
		var group Chat
		if dat.StatusCode == 200 {
			json.Unmarshal(dat.ResultByte, &group)
			chats[k] = group
		}
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": chats, "me": uuiduser})
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
		pass, _ := HashPassword("123456")
		uuid := "USER_" + uuid.New().String()
		user := RegisterForm{
			Apodo:     "admin",
			Password:  pass,
			CreatedAt: DateTime(),
			UserType:  "SUPERUSER",
			Chats:     map[string]string{},
			ID:        uuid,
		}
		InsertUser(user, uuid)
	}
}
