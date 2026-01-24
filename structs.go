package main

type Material struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Cantegory   string `json:"cantegory"`
	Buy1        string `json:"buy1"`
	Buy2        string `json:"buy2"`
	Buy3        string `json:"buy3"`
	Buy4        string `json:"buy4"`
	Buy5        string `json:"buy5"`
	Buy6        string `json:"buy6"`
	Sell1       string `json:"sell1"`
	Sell2       string `json:"sell2"`
	Sell3       string `json:"sell3"`
	Sell4       string `json:"sell4"`
	Sell5       string `json:"sell5"`
	Sell6       string `json:"sell6"`
}

type Proceso struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	PesoInicial float64 `json:"peso_inicial"`
	PesoFinal   float64 `json:"peso_final"`
	Fecha       string  `json:"fecha"`
	FechaHora   string  `json:"fecha_hora"`
}

type MatProceso struct {
	ID      string  `json:"id"`
	UUID    string  `json:"uuid"`
	UuidMat string  `json:"uuid_mat"`
	Cant    float64 `json:"cant"`
}

type Infoproc struct {
	Totales    map[string]float64 `json:"totales"`
	Materiales []map[string]any   `json:"materiales"`
	Pesos      map[string]float64 `json:"pesos"`
}

type DatosProcesarCompra struct {
	Usd       float64 `json:"usd"`
	Bs        float64 `json:"bs"`
	PagoMovil float64 `json:"pagomovil"`
	UUID      string  `json:"uuid"`
}

type Report struct {
	Totales map[string]float64 `json:"totales"`
	Movs    []map[string]any   `json:"movs"`
}

type BuyProccess struct {
	Peso       float64              `json:"peso"`
	PesoFinal  float64              `json:"pesoFinal"`
	Nombre     string               `json:"nombre"`
	NewUUID    string               `json:"newUuid"`
	Materiales map[string]Materiall `json:"materiales"`
	Adelanto   float64              `json:"adelanto"`
	Montado    string               `json:"montado"`
	Pagos      Pago                 `json:"pago"`
	Fecha      string               `json:"fecha"`
	Hora       string               `json:"hora"`
	Cliente    Cliente              `json:"cliente"`
}

// Material representa cada entrada dentro de `materiales`.
type Materiall struct {
	Buy1        float64 `json:"buy1"`
	Buy2        float64 `json:"buy2"`
	Buy3        float64 `json:"buy3"`
	Buy4        float64 `json:"buy4"`
	Buy5        float64 `json:"buy5"`
	Buy6        float64 `json:"buy6"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Name        string  `json:"name"`
	Sell1       float64 `json:"sell1"`
	Sell2       float64 `json:"sell2"`
	Sell3       float64 `json:"sell3"`
	Sell4       float64 `json:"sell4"`
	Sell5       float64 `json:"sell5"`
	Sell6       float64 `json:"sell6"`
	UUID        string  `json:"uuid"`
	Peso        float64 `json:"peso"`
	Fecha       string  `json:"fecha"`
	Cant        float64 `json:"cant"`
}

type Pago struct {
	Usd        float64 `json:"usd"`
	Bs         float64 `json:"bs"`
	Pagomovil  float64 `json:"pagomovil"`
	UuidClient string  `json:"uuidClient"`
	Tasa       float64 `json:"tasa"`
	Tipo       string  `json:"tipo"`
	TotalBs    float64 `json:"totalbs"`
	TotalUsd   float64 `json:"totalusd"`
}

type User struct {
	Username  string  `json:"username"`
	Uuid      string  `json:"uuid"`
	Password  string  `json:"password"`
	Email     string  `json:"email"`
	Admin     int     `json:"admin"`
	Active    int     `json:"active"`
	Fecha     string  `json:"fecha"`
	Sesion    string  `json:"sesion"`
	USD       float64 `json:"usd"`
	BS        float64 `json:"bs"`
	PAGOMOVIL float64 `json:"pagomovil"`
}

type Cliente struct {
	Nombre    string            `json:"nombre"`
	Cedula    string            `json:"cedula"`
	Fecha     string            `json:"fecha"`
	Compras   map[string]string `json:"compras"`
	Ventas    map[string]string `json:"ventas"`
	PorPagar  float64           `json:"pagar"`
	PorCobrar float64           `json:"cobrar"`
	Telefono  string            `json:"telefono"`
	Direccion string            `json:"direccion"`
	Notas     string            `json:"notas"`
}

type MovCaj struct {
	USD       float64 `json:"usd"`
	BS        float64 `json:"bs"`
	PAGOMOVIL float64 `json:"pagomovil"`
	HoraFecha string  `json:"hora"`
	Concepto  string  `json:"concepto"`
}

type MovMat struct {
	UUID      string  `json:"uuid"`
	Username  string  `json:"username"`
	Fecha     string  `json:"fecha"`
	Cant      float64 `json:"cant"`
	OldCant   float64 `json:"oldcant"`
	NewCant   float64 `json:"newcant"`
	HoraFecha string  `json:"hora"`
	Concepto  string  `json:"concepto"`
}

type Dolar struct {
	Hora     string  `json:"hora"`
	Username string  `json:"username"`
	Precio   float64 `json:"precio"`
}

type MovCob struct {
	Username  string  `json:"username"`
	Fecha     string  `json:"fecha"`
	Monto     float64 `json:"monto"`
	OldMonto  float64 `json:"oldmonto"`
	NewMonto  float64 `json:"newmonto"`
	HoraFecha string  `json:"hora"`
	Concepto  string  `json:"concepto"`
	Cedula    string  `json:"cedula"`
}

type Index struct {
	Fecha     string `json:"fecha"`
	Prefix    string `json:"prefix"`
	HoraFecha string `json:"hora"`
	Info      string `json:"info"`
}
