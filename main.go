package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var SETTINGS = map[string]any{}
var DB = ""

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
	DB = SETTINGS["database"].(string)
	go executeBadger()
	fmt.Println("iniciando...")
	time.Sleep(3 * time.Second)
	r := GinRouter()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/despacho", func(c *gin.Context) {
		c.HTML(200, "despacho.html", nil)
	})
	r.GET("/entrega", func(c *gin.Context) {
		c.HTML(200, "entrega.html", nil)
	})
	r.GET("/dash", func(ctx *gin.Context) {
		hash := ctx.Param("hash")
		dat, _ := AssocSecure("SELECT * FROM users WHERE sesion = ?", hash)
		if len(dat) < 1 {
			ctx.Redirect(http.StatusFound, "/")
			return
		}
		jwt, err := GenerateJWT(dat[0]["uuid"].(string), SETTINGS["jwt"].(string), 10)
		if err != nil {
			ctx.Redirect(http.StatusFound, "/")
			return
		}
		ctx.HTML(200, "index.html", gin.H{
			"Hash":      hash,
			"Username":  dat[0]["username"].(string),
			"Timestamp": time.Now().Unix(),
			"Jwt":       jwt,
		})
	})

	r.POST("/login", login)
	r.POST("/materiales", materiales)
	r.POST("/createProc", createProc)
	r.POST("/insertProc", insertProc)
	r.POST("/updateProc", updateProc)
	r.POST("/updateMatProc", updateMatProc)
	r.POST("/listProcs", listProcs)
	r.POST("/getAccount", getAccount)
	r.POST("/finalize", finalize)

	api := r.Group("/api", AuthMiddleware(SETTINGS["jwt"].(string)))
	{
		api.POST("/menu", menu)
		api.POST("/createMaterial", createMaterial)
		api.POST("/getAccount", getAccount)
		api.POST("/procesarCompra", procesarCompra)
		api.POST("/insertCaja", insertCaja)
		api.POST("/getCaja", getCaja)
		api.POST("/updateMaterial", updateMaterial)
		api.POST("/qttBuyReport", qttBuyReport)
		api.POST("/pagar", pagar)
		api.POST("/finalizeP", finalizeP)
		api.POST("/cantidades", cantidades)
		api.POST("/despachar", despachar)
		api.POST("/consulta", consulta)
		api.POST("/usuarios", usuarios)
		api.POST("/updateProduct", updateProduct)
		api.POST("/registerUser", registerUser)
		api.POST("/updateUser", updateUser)
		api.POST("/pagarCedula", pagarCedula)
		api.POST("/cobrar", cobrar)
		api.POST("/newCobrar", newCobrar)
		api.POST("/archivar", archivar)
		api.POST("/showPagos", showPagos)
		api.POST("/showMovsMat", showMovsMat)
	}

	//BuildDatabase(SETTINGS["database"].(string), "sql.sql")
	insertAdmin()
	fmt.Println("server ir running...", SETTINGS["port"])
	r.Run("0.0.0.0:" + SETTINGS["port"].(string))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func executeBadger() {
	if runtime.GOOS == "windows" {
		if !fileExists("lmdb_server.exe") {
			return
		}
		runnnnn("./lmdb_server.exe")
	} else {
		if !fileExists("lmdb_server") {
			return
		}
		runnnnn("./lmdb_server")
	}
}

func showMovsMat(ctx *gin.Context) {
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("LIKE", "MOV_"+datos["cedula"], SETFILE, 100000)
	var movs map[string]MovMat
	json.Unmarshal(dat.ResultByte, &movs)
	ctx.JSON(200, gin.H{"status": "error", "message": "ok", "data": movs})
}

func showPagos(ctx *gin.Context) {
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("LIKE", "COBRO_"+datos["cedula"], SETFILE, 100000)
	var movs map[string]MovCob
	json.Unmarshal(dat.ResultByte, &movs)
	ctx.JSON(200, gin.H{"status": "error", "message": "ok", "data": movs})
}

func archivar(ctx *gin.Context) {
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("LIKE", datos["prefix"], SETFILE, 1000000)
	var results map[string]any
	json.Unmarshal(dat.ResultByte, &results)
	for k, v := range results {
		QueryBadger("INSERT", "AR_"+k, SETFILE, v)
		QueryBadger("DELETE", k, SETFILE, "")
		CreateIndex("ARCHIVADO", "AR_"+k, "ARCHIVADO DE "+k+" A AR_"+k)
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "datos archivados correctamente", "data": ""})
}

func newCobrar(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var cliente Cliente
	err := ctx.ShouldBindJSON(&cliente)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	cuentasCobrar := map[string]Cliente{}
	dat, _ := QueryBadger("SELECT", "COBRAR", SETFILE, "")
	json.Unmarshal(dat.ResultByte, &cuentasCobrar)
	cli := "CLIENT_" + cliente.Cedula
	dat, _ = QueryBadger("SELECT", cli, SETFILE, "")
	mov := Prefix("COBRO", cliente.Cedula)
	if dat.StatusCode == 404 {
		cliente.Fecha = Hoy()
		cliente.Compras = map[string]string{}
		cliente.Ventas = map[string]string{}
		cuentasCobrar[cli] = cliente
		QueryBadger("INSERT", cli, SETFILE, cliente)
		QueryBadger("UPDATE", "COBRAR", SETFILE, cuentasCobrar)
		mm := MovCob{
			Username:  uuiduser.(string),
			Cedula:    cliente.Cedula,
			Monto:     cliente.PorPagar,
			OldMonto:  0,
			NewMonto:  cliente.PorPagar,
			Fecha:     Hoy(),
			HoraFecha: DateTime(),
			Concepto:  "REGISTRO DE MONTO",
		}
		QueryBadger("INSERT", mov, SETFILE, mm)
		CreateIndex("I_COBRO", mov, "")
		ctx.JSON(200, gin.H{"status": "success", "message": "Cuenta creada con exito", "data": ""})
		return
	}

	monto := cliente.PorPagar
	telefono := cliente.Telefono
	direccion := cliente.Direccion
	notas := cliente.Notas

	json.Unmarshal(dat.ResultByte, &cliente)
	mm := MovCob{
		Username:  uuiduser.(string),
		Cedula:    cliente.Cedula,
		Monto:     monto,
		OldMonto:  cliente.PorPagar,
		NewMonto:  cliente.PorPagar + monto,
		Fecha:     Hoy(),
		HoraFecha: DateTime(),
		Concepto:  "REGISTRO DE MONTO",
	}
	cliente.PorPagar += monto
	cliente.Telefono = telefono
	cliente.Direccion = direccion
	cliente.Notas = notas
	cuentasCobrar[cli] = cliente
	if cliente.PorPagar <= 0 {
		cliente.PorPagar = 0
		delete(cuentasCobrar, cli)
	}
	QueryBadger("INSERT", mov, SETFILE, mm)
	CreateIndex("I_COBRO", mov, "")
	QueryBadger("UPDATE", cli, SETFILE, cliente)
	QueryBadger("UPDATE", "COBRAR", SETFILE, cuentasCobrar)
	ctx.JSON(200, gin.H{"status": "success", "message": "Cuenta creada con exito", "data": ""})
}

func Prefix(prefix, key string) string {
	return prefix + "_" + key + "_" + Hoy() + "_" + Str(UnixMillisecTime())
}

func cobrar(ctx *gin.Context) {
	cuentasCobrar := map[string]Cliente{}
	dat, _ := QueryBadger("SELECT", "COBRAR", SETFILE, "")
	if dat.StatusCode == 404 {
		QueryBadger("INSERT", "COBRAR", SETFILE, cuentasCobrar)
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": cuentasCobrar})
		return
	}
	json.Unmarshal(dat.ResultByte, &cuentasCobrar)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": cuentasCobrar})
}

func pagarCedula(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var ced map[string]string
	ctx.ShouldBindJSON(&ced)
	var clientes map[string]Cliente
	fmt.Println(ced["cedula"])
	dat, _ := QueryBadger("LIKE", "CLIENT_"+ced["cedula"], SETFILE, 10000)
	json.Unmarshal(dat.ResultByte, &clientes)
	procs := map[string]BuyProccess{}
	for _, v := range clientes {
		for i := range v.Compras {
			dat, _ := QueryBadger("SELECT", i, SETFILE, "")
			var proc BuyProccess
			json.Unmarshal(dat.ResultByte, &proc)
			proc.Cliente = v
			procs[i] = proc
		}
	}
	tasa, _ := QueryBadger("SELECT", "DOLAR", SETFILE, "")
	var dolar Dolar
	json.Unmarshal(tasa.ResultByte, &dolar)
	var user User
	dat, _ = QueryBadger("SELECT", uuiduser.(string), SETFILE, "")
	json.Unmarshal(dat.ResultByte, &user)
	datos := map[string]float64{
		"usd":       user.USD,
		"bs":        user.BS,
		"pagomovil": user.PAGOMOVIL,
		"tasa":      dolar.Precio,
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": procs, "caja": datos})
}

func updateUser(ctx *gin.Context) {
	var user User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	pass, _ := HashPassword(user.Password)
	dat, _ := QueryBadger("SELECT", user.Uuid, SETFILE, "")
	json.Unmarshal(dat.ResultByte, &user)
	user.Password = pass
	QueryBadger("UPDATE", user.Username, SETFILE, user)
	QueryBadger("UPDATE", user.Uuid, SETFILE, user)
	ctx.JSON(200, gin.H{"status": "success", "message": "usuario actualizado correctamente", "data": ""})
}

func registerUser(ctx *gin.Context) {
	var user User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	uuid := Str(UnixMillisecTime())
	uuid = "USER_" + uuid
	pass, _ := HashPassword(user.Password)
	user.Fecha = Hoy()
	user.Active = 1
	user.Admin = 0
	user.Password = pass
	dat, _ := QueryBadger("SELECT", user.Username, SETFILE, "")
	if dat.StatusCode == 200 {
		ctx.JSON(401, gin.H{"status": "error", "message": "nombre de usuario ya en uso", "data": ""})
		return
	}
	QueryBadger("INSERT", user.Username, SETFILE, user)
	QueryBadger("INSERT", uuid, SETFILE, user)
	ctx.JSON(200, gin.H{"status": "success", "message": "usuario registrado correctamente", "data": ""})
}

func updateProduct(ctx *gin.Context) {
	uuidUser, _ := ctx.Get("uuid")
	var material Materiall
	ctx.ShouldBindJSON(&material)
	dat, _ := QueryBadger("SELECT", material.UUID, SETFILE, "")
	var mm Materiall
	json.Unmarshal(dat.ResultByte, &mm)
	mm.Buy1 = material.Buy1
	mm.Buy2 = material.Buy2
	mm.Buy3 = material.Buy3
	mm.Buy4 = material.Buy4
	mm.Buy5 = material.Buy5
	mm.Buy6 = material.Buy6
	mm.Sell1 = material.Sell1
	mm.Sell2 = material.Sell2
	mm.Sell3 = material.Sell3
	mm.Sell4 = material.Sell4
	mm.Sell5 = material.Sell5
	mm.Sell6 = material.Sell6
	QueryBadger("UPDATE", mm.UUID, SETFILE, mm)
	mov := MovMat{
		UUID:      mm.UUID,
		Username:  uuidUser.(string),
		Fecha:     Hoy(),
		Cant:      0,
		OldCant:   0,
		NewCant:   0,
		HoraFecha: HourNow(),
		Concepto:  "MODIFICACION DE PRECIOS",
	}
	QueryBadger("INSERT", "MOV_"+mm.UUID+"_"+Str(UnixMillisecTime()), SETFILE, mov)
	//Execute(DB, "UPDATE materiales SET name = ?, buy1 = ?, buy2 = ?, buy3 = ?, sell1 = ?, sell2 = ?, sell3 = ? WHERE uuid = ?", material.Name, material.Buy1, material.Buy2, material.Buy3, material.Sell1, material.Sell2, material.Sell3, material.UUID)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func usuarios(ctx *gin.Context) {
	dat, _ := QueryBadger("LIKE", "USER", SETFILE, 1000) //AssocSecure(DB, "SELECT * FROM users WHERE id > ?", 1)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": dat.Result})
}

func consulta(ctx *gin.Context) {
	var datos map[string]string
	ctx.ShouldBindJSON(&datos)
	fmt.Println(datos)
	fechas, _ := GenerarFechasConsecutivas(datos["fecha1"], datos["fecha2"])
	procesosCompras := map[string]BuyProccess{}
	procesosVentas := map[string]BuyProccess{}
	for _, i := range fechas {
		dat, _ := QueryBadger("SELECT", "I_PROC_"+i, SETFILE, "")
		if dat.StatusCode == 404 {
			continue
		}
		if dat.StatusCode != 200 {
			continue
		}
		index := map[string]Index{}
		json.Unmarshal(dat.ResultByte, &index)
		for k := range index {
			dat, _ = QueryBadger("SELECT", k, SETFILE, "")
			if dat.StatusCode == 404 {
				continue
			}
			if dat.StatusCode != 200 {
				continue
			}
			proceso := BuyProccess{}
			json.Unmarshal(dat.ResultByte, &proceso)
			procesosCompras[k] = proceso
		}
	}
	//compras, _ := AssocSecure(DB, "SELECT * FROM compras WHERE fecha BETWEEN ? AND ?", datos["fecha1"], datos["fecha2"])
	//ventas, _ := AssocSecure(DB, "SELECT * FROM ventas WHERE fecha BETWEEN ? AND ?", datos["fecha1"], datos["fecha2"])
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "", "compras": procesosCompras, "ventas": procesosVentas})
}

func despachar(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos map[string]any
	ctx.ShouldBindJSON(&datos)
	dat, _ := QueryBadger("SELECT", datos["uuid"].(string), SETFILE, "")
	var mm Materiall
	json.Unmarshal(dat.ResultByte, &mm)
	cant := Float(datos["cantidad"])
	mov := MovMat{
		UUID:      mm.UUID,
		Username:  uuiduser.(string),
		Fecha:     Hoy(),
		Cant:      cant,
		OldCant:   mm.Cant,
		NewCant:   mm.Cant - cant,
		HoraFecha: HourNow(),
		Concepto:  "VENTA DE MATERIAL",
	}
	mm.Cant -= cant
	QueryBadger("UPDATE", mm.UUID, SETFILE, mm)
	mv := Prefix("MOV", mm.UUID)
	QueryBadger("INSERT", mv, SETFILE, mov)
	CreateIndex("I_MOV", mv, "")
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func cantidades(ctx *gin.Context) {
	dat, _ := QueryBadger("LIKE", "MAT", SETFILE, 1000) //AssocSecure(DB, "SELECT * FROM materiales WHERE name != ?", "")
	//fmt.Println(dat)

	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": dat.Result})
}

func finalizeP(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos Pago
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	fecha := Hoy()
	hora_fecha := DateTime()
	tas := datos.Tasa
	claveProceso := Prefix("COMPRA", datos.UuidClient) //"COMPRA_" + fecha + "_" + datos.UuidClient
	claveCaja := Prefix("CAJ", uuiduser.(string))      //"CAJ_" + fecha + "_" + uuiduser.(string) + Str(UnixMillisecTime())
	dat, err := QueryBadger("SELECT", datos.UuidClient, SETFILE, "")
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	if dat.StatusCode == 404 {
		ctx.JSON(404, gin.H{"status": "error", "message": "proceso n o encontrado", "data": ""})
		return
	}
	var ff BuyProccess
	json.Unmarshal(dat.ResultByte, &ff)
	dat, _ = QueryBadger("SELECT", uuiduser.(string), SETFILE, "")
	var user User
	json.Unmarshal(dat.ResultByte, &user)
	if datos.Usd > user.USD {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient usd balance", "data": ""})
		return
	}
	if datos.Bs > user.BS {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient bs balance", "data": ""})
		return
	}
	if datos.Pagomovil > user.PAGOMOVIL {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient pagomovil balance", "data": ""})
		return
	}
	tttt := (((datos.Pagomovil + datos.Bs) / tas) + datos.Usd)
	if tttt <= 0 {
		ctx.JSON(400, gin.H{"status": "error", "message": "imposible operacion", "data": ""})
		return
	}
	ff.Pagos = datos
	ff.Hora = hora_fecha
	ff.Fecha = fecha
	user.BS -= datos.Bs
	user.USD -= datos.Usd
	user.PAGOMOVIL -= datos.Pagomovil
	mov := MovCaj{
		USD:       -datos.Usd,
		BS:        -datos.Bs,
		PAGOMOVIL: -datos.Pagomovil,
		HoraFecha: hora_fecha,
		Concepto:  "COMPRA DE MATERIAL",
	}
	claveMovCliente := Prefix("MOV", ff.Cliente.Cedula) //"MOV_" + ff.Cliente.Cedula + "_" + fecha + "_" + Str(UnixMillisecTime())

	QueryBadger("UPDATE", user.Uuid, SETFILE, user)
	QueryBadger("UPDATE", user.Username, SETFILE, user)

	QueryBadger("INSERT", claveCaja, SETFILE, mov)
	CreateIndex("I_CAJ", claveCaja, "")
	if datos.Tipo == "ADELANTO" {
		ff.Adelanto = (((datos.Pagomovil + datos.Bs) / tas) + datos.Usd)
		ff.Pagos.TotalBs = datos.TotalBs - (datos.Pagomovil + datos.Bs)
		ff.Pagos.TotalUsd = datos.TotalUsd - datos.Usd
		QueryBadger("UPDATE", datos.UuidClient, SETFILE, ff)
		ff.Materiales = map[string]Materiall{}
		QueryBadger("INSERT", claveMovCliente, SETFILE, ff)
		CreateIndex("I_CLIENT", claveMovCliente, "")
		// rutina de si el cliente debe plata
		cuentasCobrar := map[string]Cliente{}
		dat, _ := QueryBadger("SELECT", "COBRAR", SETFILE, "")
		json.Unmarshal(dat.ResultByte, &cuentasCobrar)
		cli := "CLIENT_" + ff.Cliente.Cedula
		_, ex := cuentasCobrar[cli]
		if ex {
			dat, _ = QueryBadger("SELECT", cli, SETFILE, "")
			var cliente Cliente
			json.Unmarshal(dat.ResultByte, &cliente)
			monto := tttt
			if cliente.PorPagar < tttt {
				monto = cliente.PorPagar
			}
			mm := MovCob{
				Username:  uuiduser.(string),
				Cedula:    cliente.Cedula,
				Monto:     monto,
				OldMonto:  cliente.PorPagar,
				NewMonto:  cliente.PorPagar - monto,
				Fecha:     Hoy(),
				HoraFecha: DateTime(),
				Concepto:  "REGISTRO DE MONTO",
			}
			cliente.PorPagar -= monto
			cuentasCobrar[cli] = cliente
			if cliente.PorPagar <= 0 {
				cliente.PorPagar = 0
				delete(cuentasCobrar, cli)
			}
			move := Prefix("COBRO", cliente.Cedula)
			QueryBadger("INSERT", move, SETFILE, mm)
			CreateIndex("I_COBRO", move, "")

			QueryBadger("UPDATE", cli, SETFILE, cliente)
			QueryBadger("UPDATE", "COBRAR", SETFILE, cuentasCobrar)
		}
		//
		ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
		return
	}
	QueryBadger("INSERT", claveProceso, SETFILE, ff)
	CreateIndex("I_PROC", claveProceso, "")
	QueryBadger("INSERT", claveMovCliente, SETFILE, ff)
	CreateIndex("I_CLIENT", claveMovCliente, "")

	QueryBadger("DELETE", datos.UuidClient, SETFILE, "")
	for k, v := range ff.Materiales {
		mat, _ := QueryBadger("SELECT", k, SETFILE, "")
		var mate Materiall
		json.Unmarshal(mat.ResultByte, &mate)

		keyMovmat := Prefix("MOV", k)
		mMat := MovMat{
			Fecha:    fecha,
			UUID:     k,
			Cant:     v.Peso,
			OldCant:  mate.Cant,
			NewCant:  mate.Cant + v.Peso,
			Concepto: "COMPRA DE MATERIAL",
			Username: uuiduser.(string),
		}
		mate.Cant += v.Peso
		QueryBadger("UPDATE", k, SETFILE, mate)
		QueryBadger("INSERT", keyMovmat, SETFILE, mMat)
		CreateIndex("I_MAT", keyMovmat, "")
	}
	cli := "CLIENT_" + ff.Cliente.Cedula
	dat, _ = QueryBadger("SELECT", cli, SETFILE, "")
	var cliente Cliente
	json.Unmarshal(dat.ResultByte, &cliente)
	delete(cliente.Compras, datos.UuidClient)
	QueryBadger("UPDATE", cli, SETFILE, cliente)
	// rutina de si el cliente debe plata
	/*cuentasCobrar := map[string]Cliente{}
	dat, _ = QueryBadger("SELECT", "COBRAR", SETFILE, "")
	json.Unmarshal(dat.ResultByte, &cuentasCobrar)
	_, ex := cuentasCobrar[cli]
	if ex {
		dat, _ = QueryBadger("SELECT", cli, SETFILE, "")
		var cliente Cliente
		json.Unmarshal(dat.ResultByte, &cliente)
		monto := tttt
		if cliente.PorPagar < tttt {
			monto = cliente.PorPagar
		}
		mm := MovCob{
			Username:  uuiduser.(string),
			Cedula:    cliente.Cedula,
			Monto:     monto,
			OldMonto:  cliente.PorPagar,
			NewMonto:  cliente.PorPagar - monto,
			Fecha:     Hoy(),
			HoraFecha: DateTime(),
			Concepto:  "REGISTRO DE MONTO",
		}
		cliente.PorPagar -= monto
		cuentasCobrar[cli] = cliente
		if cliente.PorPagar <= 0 {
			cliente.PorPagar = 0
			delete(cuentasCobrar, cli)
		}
		move := Prefix("COBRO", cliente.Cedula)
		QueryBadger("INSERT", move, SETFILE, mm)
		QueryBadger("UPDATE", cli, SETFILE, cliente)
		QueryBadger("UPDATE", "COBRAR", SETFILE, cuentasCobrar)
	}
	*/
	//
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func CreateIndex(prefix, key, info string) {
	fecha := Hoy()
	index := Index{
		Fecha:     fecha,
		HoraFecha: HourNow(),
		Info:      info,
		Prefix:    prefix,
	}
	mapaIndex := map[string]Index{}
	dat, _ := QueryBadger("SELECT", prefix+"_"+fecha, SETFILE, "")
	if dat.StatusCode == 404 {
		mapaIndex[key] = index
		QueryBadger("INSERT", prefix+"_"+fecha, SETFILE, mapaIndex)
		return
	}
	json.Unmarshal(dat.ResultByte, &mapaIndex)
	mapaIndex[key] = index
	QueryBadger("UPDATE", prefix+"_"+fecha, SETFILE, mapaIndex)
}

func pagar(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	dat, err := QueryBadger("LIKE", "PROC", SETFILE, 10000)
	if err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": "error", "data": ""})
		return
	}
	//procs := dat.Result
	//fmt.Println(procs)
	procs := map[string]BuyProccess{}
	json.Unmarshal(dat.ResultByte, &procs)
	for k, v := range procs {
		cli := "CLIENT_" + v.Cliente.Cedula
		dat, _ = QueryBadger("SELECT", cli, SETFILE, "")
		var cliente Cliente
		json.Unmarshal(dat.ResultByte, &cliente)
		v.Cliente = cliente
		procs[k] = v
	}
	tasa, _ := QueryBadger("SELECT", "DOLAR", SETFILE, "")
	var dolar Dolar
	json.Unmarshal(tasa.ResultByte, &dolar)
	var user User
	dat, _ = QueryBadger("SELECT", uuiduser.(string), SETFILE, "")
	json.Unmarshal(dat.ResultByte, &user)
	datos := map[string]float64{
		"usd":       user.USD,
		"bs":        user.BS,
		"pagomovil": user.PAGOMOVIL,
		"tasa":      dolar.Precio,
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": procs, "caja": datos})
}

func finalize(ctx *gin.Context) {
	var pp BuyProccess
	err := ctx.ShouldBindJSON(&pp)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "danger", "message": err.Error(), "data": ""})
		return
	}
	cli := "CLIENT_" + pp.Cliente.Cedula
	dat, _ := QueryBadger("SELECT", cli, SETFILE, "")
	pp.NewUUID = "PROC_" + pp.NewUUID
	if dat.StatusCode == 404 {
		pp.Cliente.Compras = map[string]string{pp.NewUUID: pp.NewUUID}
		pp.Cliente.Fecha = Hoy()
		pp.Cliente.Ventas = map[string]string{}
		QueryBadger("INSERT", cli, SETFILE, pp.Cliente)
	} else {
		var cliente Cliente
		json.Unmarshal(dat.ResultByte, &cliente)
		cliente.Compras[pp.NewUUID] = pp.NewUUID
		QueryBadger("UPDATE", cli, SETFILE, cliente)
	}
	QueryBadger("INSERT", pp.NewUUID, SETFILE, pp)
	ctx.JSON(200, gin.H{"status": "success", "message": "exitoso", "data": ""})
}

func listProcs(ctx *gin.Context) {
	dats, _ := AssocSecure(DB, "SELECT * FROM proceso WHERE id > ?", 0)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": dats})
}

func qttBuyReport(ctx *gin.Context) {
	var datos map[string]string
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	sumarize, err := AssocSecure(DB, "SELECT SUM(cant) AS kilos, SUM(usd) AS usd, SUM(bs) AS bs, SUM(pagomovil) AS pagomovil, SUM(total_pagado) AS total_pagado FROM compras WHERE fecha BETWEEN ? AND ?;", datos["fecha1"], datos["fecha2"])
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	kilos := float64(0)
	usd := float64(0)
	bs := float64(0)
	pagomovil := float64(0)
	totalPagado := float64(0)
	if len(sumarize) > 0 {
		kilos = Float(sumarize[0]["kilos"])
		usd = Float(sumarize[0]["usd"])
		bs = Float(sumarize[0]["bs"])
		pagomovil = Float(sumarize[0]["pagomovil"])
		totalPagado = Float(sumarize[0]["total_pagado"])
	}
	dat, err := AssocSecure(DB, "SELECT * FROM compras WHERE fecha BETWEEN ? AND ?", datos["fecha1"], datos["fecha2"])
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	data := Report{
		Totales: map[string]float64{
			"usd":          usd,
			"bs":           bs,
			"pagomovil":    pagomovil,
			"kilos":        kilos,
			"total_pagado": totalPagado,
		},
		Movs: dat,
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": data})
}

func getCaja(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	caja, _ := QueryBadger("SELECT", uuiduser.(string), SETFILE, "") //AssocSecure(DB, "SELECT SUM(usd), SUM(bs), SUM(pagomovil) FROM caja_chica WHERE username = ?", uuiduser)
	user := User{}
	json.Unmarshal(caja.ResultByte, &user)
	datos := map[string]float64{
		"usd":       user.USD,
		"bs":        user.BS,
		"pagomovil": user.PAGOMOVIL,
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": datos})

}

func updateMatProc(ctx *gin.Context) {
	var datos MatProceso
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	dat, err := AssocSecure(DB, "SELECT * FROM materiales WHERE uuid = ?", datos.UuidMat)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(dat) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "material not found", "data": ""})
		return
	}
	det, err := AssocSecure(DB, "SELECT * FROM proceso WHERE uuid = ?", datos.UUID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(det) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "proccess not found", "data": ""})
		return
	}
	_, err = Execute(DB, "UPDATE mat_proceso SET cant = ?, total1 = ?, total2 = ?, total3 = ? WHERE id = ?", datos.Cant, datos.Cant*Float(dat[0]["buy1"]), datos.Cant*Float(dat[0]["buy2"]), datos.Cant*Float(dat[0]["buy3"]), datos.ID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "success"})
}

func updateMaterial(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos Materiall
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	fecha := Hoy()
	sql := `
    UPDATE materiales 
    SET  name        = ?,
         description = ?,
         cantegory   = ?,
         buy1        = ?,
         buy2        = ?,
         buy3        = ?,
         sell1       = ?,
         sell2       = ?,
         sell3       = ?
    WHERE uuid = ?`

	_, err = Execute(sql,
		datos.Name,
		datos.Description,

		datos.Buy1,
		datos.Buy2,
		datos.Buy3,
		datos.Sell1,
		datos.Sell2,
		datos.Sell3,
		datos.UUID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	_, err = Execute(DB, "INSERT INTO cant_materiales (uuid, tipo, cant, fecha, username) VALUES (?, ?, ?, ?, ?)", datos.UUID, "ACTUALIZADO", 0.0, fecha, uuiduser)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": datos.UUID})
}

func updateProc(ctx *gin.Context) {
	var datos Proceso
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	_, err = Execute(DB, "UPDATE proceso SET peso_final = ? WHERE uuid = ?", datos.PesoFinal, datos.UUID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func insertCaja(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos DatosProcesarCompra
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	dat, err := QueryBadger("SELECT", uuiduser.(string), SETFILE, "")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	var user User
	err = json.Unmarshal(dat.ResultByte, &user)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	user.BS += datos.Bs
	user.USD += datos.Usd
	user.PAGOMOVIL += datos.PagoMovil
	//Execute(DB, "INSERT INTO caja_chica (usd, bs, pagomovil, fecha, fecha_hora, username) VALUES (?,?,?,?,?,?)", datos.Usd, datos.Bs, datos.PagoMovil, Hoy(), DateTime(), uuiduser)
	QueryBadger("UPDATE", user.Uuid, SETFILE, user)
	QueryBadger("UPDATE", user.Username, SETFILE, user)
	index := Prefix("CAJ", user.Uuid) //"CAJ_" + Hoy() + "_" + user.Uuid + "_" + Str(UnixMillisecTime())
	mov := MovCaj{}
	mov.BS = datos.Bs
	mov.USD = datos.Usd
	mov.PAGOMOVIL = datos.PagoMovil
	mov.HoraFecha = DateTime()
	QueryBadger("INSERT", index, SETFILE, mov)
	CreateIndex("I_CAJ", index, "")
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

func procesarCompra(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos DatosProcesarCompra
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	det, err := AssocSecure(DB, "SELECT * FROM proceso WHERE uuid = ?", datos.UUID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(det) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "proccess not found", "data": ""})
		return
	}
	materials, _ := AssocSecure(DB, "SELECT * FROM mat_proceso WHERE uuid = ?", datos.UUID)
	if len(materials) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "any materials to proccess", "data": ""})
		return
	}
	caja, _ := AssocSecure(DB, "SELECT SUM(usd), SUM(bs), SUM(pagomovil) FROM caja_chica WHERE username = ?", uuiduser)
	usd := float64(0)
	bs := float64(0)
	pagomovil := float64(0)
	if len(caja) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "no money to pay", "data": ""})
		return
	}
	usd, bs, pagomovil = Float(caja[0]["SUM(usd)"]), Float(caja[0]["SUM(bs)"]), Float(caja[0]["SUM(pagomovil)"])
	if datos.Usd > usd {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient usd balance", "data": ""})
		return
	}
	if datos.Bs > bs {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient bs balance", "data": ""})
		return
	}
	if datos.PagoMovil > pagomovil {
		ctx.JSON(401, gin.H{"status": "error", "message": "insuficient pagomovil balance", "data": ""})
		return
	}
	fecha := Hoy()
	hora_fecha := DateTime()
	tasa, _ := AssocSecure(DB, "SELECT dolar FROM tasa WHERE id = ?", 1)
	tas := float64(1)
	if len(tasa) > 0 {
		tas = Float(tasa[0]["dolar"])
	}
	for cont, i := range materials {
		if cont == 0 {
			Execute(DB, "INSERT INTO compras (uuid, uuid_mat, nombre, cant, usd, bs, pagomovil, total_pagado, fecha, fecha_hora, username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", datos.UUID, i["uuid_mat"], i["nombre"], i["cant"], datos.Usd, datos.Bs, datos.PagoMovil, (((datos.PagoMovil + datos.Bs) / tas) + datos.Usd), fecha, hora_fecha, uuiduser)

			Execute(DB, "INSERT INTO caja_chica (usd, bs, pagomovil, fecha, fecha_hora, username) VALUES (?,?,?,?,?,?)", -datos.Usd, -datos.Bs, -datos.PagoMovil, fecha, hora_fecha, uuiduser)
		} else {
			Execute(DB, "INSERT INTO compras (uuid, uuid_mat, nombre, cant, usd, bs, pagomovil, total_pagado, fecha, fecha_hora, username) VALUES (?,?,?,?,?,?,?,?,?,?,?)", datos.UUID, i["uuid_mat"], i["nombre"], i["cant"], 0, 0, 0, 0, fecha, hora_fecha, uuiduser)
		}
		Execute(DB, "INSERT INTO cant_materiales (uuid, tipo, cant, fecha, username) VALUES (?, ?, ?, ?, ?)", i["uuid_mat"], "COMPRA", i["cant"], fecha, uuiduser)

	}
	Execute(DB, "DELETE FROM mat_proceso WHERE uuid = ?", datos.UUID)
	Execute(DB, "DELETE FROM proceso WHERE uuid = ?", datos.UUID)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": ""})
}

// que salga aparte el simbolo de la divisa en la caja y en la etiqueta
func getAccount(ctx *gin.Context) {
	var datos map[string]string
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	det, err := AssocSecure(DB, "SELECT * FROM proceso WHERE uuid = ?", datos["uuid"])
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(det) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "proccess not found", "data": ""})
		return
	}
	suma, _ := AssocSecure(DB, "SELECT SUM(total1), SUM(total2), SUM(total3), SUM(cant) FROM mat_proceso WHERE uuid = ?", datos["uuid"])
	total1 := float64(0)
	total2 := float64(0)
	total3 := float64(0)
	kilos := float64(0)
	peso_inicial := Float(det[0]["peso_inicial"])
	peso_final := Float(det[0]["peso_final"])
	materials, _ := AssocSecure(DB, "SELECT * FROM mat_proceso WHERE uuid = ?", datos["uuid"])
	if len(suma) > 0 {
		total1 = Float(suma[0]["SUM(total1)"])
		total2 = Float(suma[0]["SUM(total2)"])
		total3 = Float(suma[0]["SUM(total3)"])
		kilos = Float(suma[0]["SUM(cant)"])
	}
	datum := Infoproc{
		Totales: map[string]float64{
			"total1": total1,
			"total2": total2,
			"total3": total3,
		},
		Materiales: materials,
		Pesos: map[string]float64{
			"peso_inicial": peso_inicial,
			"peso_final":   peso_final,
			"kilos":        kilos,
		},
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": datum})
}

func insertProc(ctx *gin.Context) {
	var datos MatProceso
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	dat, err := AssocSecure(DB, "SELECT * FROM materiales WHERE uuid = ?", datos.UuidMat)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(dat) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "material not found", "data": ""})
		return
	}
	det, err := AssocSecure(DB, "SELECT * FROM proceso WHERE uuid = ?", datos.UUID)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if len(det) < 1 {
		ctx.JSON(404, gin.H{"status": "error", "message": "proccess not found", "data": ""})
		return
	}
	_, err = Execute(DB, "INSERT INTO mat_proceso (uuid, uuid_mat, nombre, cant, total1, total2, total3) VALUES (?, ?, ?, ?, ?,?,?)", datos.UUID, datos.UuidMat, dat[0]["name"], datos.Cant, datos.Cant*Float(dat[0]["buy1"]), datos.Cant*Float(dat[0]["buy2"]), datos.Cant*Float(dat[0]["buy3"]))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": "success"})
}

func createProc(ctx *gin.Context) {
	var datos Proceso
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	uuid := uuid.New().String()
	fecha := Hoy()
	_, err = Execute(DB, "INSERT INTO proceso (uuid, name, peso_inicial, peso_final, fecha, fecha_hora) VALUES (?, ?, ?, ?, ?,?)", uuid, datos.Name, datos.PesoInicial, datos.PesoFinal, fecha, DateTime())
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": uuid})
}

func createMaterial(ctx *gin.Context) {
	uuiduser, _ := ctx.Get("uuid")
	var datos Materiall
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	uuid := Str(UnixMillisecTime())
	uuid = "MAT_" + uuid
	fecha := Hoy()
	datos.Fecha = fecha
	datos.UUID = uuid
	mov := Str(UnixMillisecTime() + 100)
	mov = "MOV_" + uuid + "_" + fecha + "_" + mov
	mMat := MovMat{
		Fecha:    fecha,
		UUID:     uuid,
		Cant:     0,
		OldCant:  0,
		NewCant:  0,
		Concepto: "REGISTRO DE PRODUCTO",
		Username: uuiduser.(string),
	}
	QueryBadger("INSERT", uuid, SETFILE, datos)
	QueryBadger("INSERT", mov, SETFILE, mMat)
	CreateIndex("I_MAT", mov, "")
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": uuid})
}

func materiales(ctx *gin.Context) {
	dat, _ := QueryBadger("LIKE", "MAT", SETFILE, 1000) //AssocSecure(DB, "SELECT * FROM materiales WHERE id > ? AND name != ?", 0, "")
	/*mats := make(map[string]any)
	for _, i := range dat {
		mats[i["uuid"].(string)] = i
	}*/
	//mats := []Materiall{}
	//fmt.Println(dat)
	//json.Unmarshal(dat.ResultByte, &mats)
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": dat.Result})
}

func menu(ctx *gin.Context) {
	//uuid, _ := ctx.Get("uuid")
	text := `
<nav class="navbar navbar-expand-lg bg-dark navbar-dark sticky-top">
  <div class="container-fluid">
    <a class="navbar-brand" href="#">SISTEMA</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse"
            data-bs-target="#navMenu" aria-controls="navMenu"
            aria-expanded="false" aria-label="Alternar navegación">
      <span class="navbar-toggler-icon"></span>
    </button>

    <div class="collapse navbar-collapse" id="navMenu">
      <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
        <li class="nav-item"><a class="nav-link active" href="javascript:void(0)" onclick="pagar()">Pagar</a></li>
		<li class="nav-item"><a class="nav-link active" href="javascript:void(0)" onclick="cobrar()">Cobrar</a></li>
		<li class="nav-item"><a class="nav-link active" href="javascript:void(0)" onclick="cajachica()">Caja</a></li>
		<li class="nav-item"><a class="nav-link active" href="javascript:void(0)" onclick="usuarios()">Usuarios</a></li>
        <li class="nav-item"><a class="nav-link" href="javascript:void(0)" onclick="entrega()">Entrega</a></li>
        <li class="nav-item"><a class="nav-link" href="javascript:void(0)" onclick="reporte()">Reportes</a></li>
        <li class="nav-item"><a class="nav-link" href="javascript:void(0)" onclick="registerProd()">Registrar Productos</a></li>
      </ul>
    </div>
  </div>
</nav>
	`
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": text})
}

func login(ctx *gin.Context) {
	var datos map[string]any
	err := ctx.ShouldBindJSON(&datos)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	if datos["username"] == nil {
		fmt.Println("username is nill")
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	dats, err := QueryBadger("SELECT", datos["username"].(string), SETFILE, "")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	//dats, err := AssocSecure(SETTINGS["database"].(string), "SELECT * FROM users WHERE username = ?", datos["username"])
	if dats.StatusCode == 404 {
		ctx.JSON(404, gin.H{"status": "error", "message": "username not found", "data": ""})
		return
	}
	var user User
	err = json.Unmarshal(dats.ResultByte, &user)
	if !ValidatePassword(datos["password"].(string), user.Password) {
		ctx.JSON(401, gin.H{"status": "error", "message": "unauthorized", "data": ""})
		return
	}
	jwt, err := GenerateJWT(user.Uuid, SETTINGS["jwt"].(string), 9)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"status": "error", "message": err, "data": ""})
		return
	}
	ctx.JSON(200, gin.H{"status": "success", "message": "ok", "data": jwt})
}

func insertAdmin() {
	dat, err := QueryBadger("SELECT", "admin", SETFILE, "") //AssocSecure(SETTINGS["database"].(string), "SELECT * FROM users WHERE username = ?", "admin@app.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	if dat.StatusCode == 404 {
		uuid := Str(UnixMillisecTime())
		uuid = "USER_" + uuid
		pass, _ := HashPassword("123456")
		hhs := Str(UnixMillisecTime())
		//Execute(SETTINGS["database"].(string), "INSERT INTO users (username, uuid, password, admin, active, fecha) VALUES (?,?,?,?,?,?)", "admin@app.com", uuid, pass, 1, 1, Hoy())
		user := User{
			Username: "admin",
			Password: pass,
			Uuid:     uuid,
			Fecha:    Hoy(),
			Active:   1,
			Admin:    1,
		}
		_, err := QueryBadger("INSERT", "admin", SETFILE, user)
		if err != nil {
			fmt.Println(err)
			return
		}
		QueryBadger("INSERT", uuid, SETFILE, user)
		QueryBadger("INSERT", "DOLAR", SETFILE, Dolar{Hora: DateTime(), Precio: 1})
		QueryBadger("INSERT", "MOV_DOLAR_"+hhs, SETFILE, Dolar{Hora: DateTime(), Precio: 1, Username: uuid})
		//QueryBadger("INSERT", "CAJ_"+uuid, SETFILE, []any{})
	}
}
