package main

type Message struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	FromUser  string `json:"from_user"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
	EditedAt  string `json:"edited_at"`
	DeletedAt string `json:"deleted_at"`
	ReplyToID string `json:"reply_to_id"`
	ResendAt  string `json:"resend_at"`
	Status    string `json:"status"`
	ReSendts  int64  `json:"resents"`
	Filename  string `json:"filename"`
	MediaType string `json:"mediatype"`
	Apodo     string `json:"apodo"`
}

type RegisterForm struct {
	Nombre         string            `form:"nombre" binding:"required"`
	Apellido       string            `form:"apellido" binding:"required"`
	Cedula         string            `form:"cedula" binding:"required"`
	Apodo          string            `form:"apodo" binding:"required"`
	Correo         string            `form:"correo" binding:"required,email"`
	Password       string            `form:"password" binding:"required"`
	Telefono       string            `form:"telefono"`
	Pais           string            `form:"pais"`
	Estado         string            `form:"estado"`
	Ciudad         string            `form:"ciudad"`
	Direccion      string            `form:"direccion"`
	FotoPerfil     string            `json:"foto_perfil"`
	SelfieDoc      string            `json:"selfie_doc"`
	FotoDocumento  string            `json:"foto_documento"`
	ReciboServicio string            `json:"recibo_servicio"`
	KycStatus      string            `json:"kycstatus"`
	Chats          map[string]string `json:"chats"`
	CreatedAt      string            `json:"created_at"`
	UserType       string            `json:"user_type"`
	ID             string            `json:"id"`
}

type RegisterFormClient struct {
	Nombre         string            `form:"nombre" binding:"required"`
	Apellido       string            `form:"apellido"`
	Cedula         string            `form:"cedula"`
	Apodo          string            `form:"apodo"`
	Correo         string            `form:"correo" binding:"required, email"`
	Password       string            `form:"password" binding:"required"`
	Telefono       string            `form:"telefono"`
	Pais           string            `form:"pais"`
	Estado         string            `form:"estado"`
	Ciudad         string            `form:"ciudad"`
	Direccion      string            `form:"direccion"`
	FotoPerfil     string            `json:"foto_perfil"`
	SelfieDoc      string            `json:"selfie_doc"`
	FotoDocumento  string            `json:"foto_documento"`
	ReciboServicio string            `json:"recibo_servicio"`
	KycStatus      string            `json:"kycstatus"`
	Chats          map[string]string `json:"chats"`
	CreatedAt      string            `json:"created_at"`
	UserType       string            `json:"user_type"`
	ID             string            `json:"id"`
}

type Multimedia struct {
	Current string `form:"current" binding:"required"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UserType  string `json:"user_type"`
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Cedula    string `json:"cedula"`
	Apodo     string `json:"apodo"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`

	Pais      string `json:"pais"`
	Estado    string `json:"estado"`
	Ciudad    string `json:"ciudad"`
	Direccion string `json:"direccion"`

	FotoPerfil     string            `json:"foto_perfil"`
	SelfieDoc      string            `json:"selfie_doc"`
	FotoDocumento  string            `json:"foto_documento"`
	ReciboServicio string            `json:"recibo_servicio"`
	KycStatus      string            `json:"kycstatus"`
	Chats          map[string]string `json:"chats"`
}

type Session struct {
	ID        string `json:"id"`
	Exp       int64  `json:"exp"`
	CreatedAt string `json:"created_at"`
}

type UserInfo struct {
	UserType   string `json:"user_type"`
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Cedula     string `json:"cedula"`
	Apodo      string `json:"apodo"`
	Correo     string `json:"correo"`
	Contrasena string `json:"password"`
	Telefono   string `json:"telefono"`

	Pais      string `json:"pais"`
	Estado    string `json:"estado"`
	Ciudad    string `json:"ciudad"`
	Direccion string `json:"direccion"`

	FotoPerfil     string `json:"foto_perfil"`
	SelfieDoc      string `json:"selfie_doc"`
	FotoDocumento  string `json:"foto_documento"`
	ReciboServicio string `json:"recibo_servicio"`
	KycStatus      string `json:"kycstatus"`
}

// / estos son los metadatos de los grupos, por ejemplo el usuario tiene una clave que almacena un mapa de estos, es decir a todos los chats que pertenece
type MetaInfoMsg struct {
	CreatedBy      string            `json:"created_by"`
	CreatedAt      string            `json:"created_at"`
	ID             string            `json:"id"`
	Denunces       int64             `json:"denunces"`
	DenuncedBy     map[string]string `json:"denunced_by"`
	Administrators map[string]string `json:"administrators"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Specialty      string            `json:"specialty"`
	Country        string            `json:"country"`
	City           string            `json:"city"`
}

type MetaInfoChats struct {
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Specialty   string `json:"specialty"`
	Country     string `json:"country"`
	City        string `json:"city"`
}

type UserActivity struct {
	Chats map[string]MetaInfoMsg `json:"chats"`
}

type Chat struct {
	ID            string            `json:"id"`
	Members       map[string]string `json:"members"`
	Avatar        string            `json:"avatar"`
	Cant          int64             `json:"cant"`
	CreatedAt     string            `json:"created_at"`
	CreatedBy     string            `json:"created_by"`
	Denunces      int64             `json:"denunces"`
	DenuncedBy    map[string]string `json:"denunced_by"`
	IsBlocked     bool              `json:"is_blocked"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Specialty     string            `json:"specialty"`
	Country       string            `json:"country"`
	City          string            `json:"city"`
	Intermediary  string            `json:"intermediary"`
	Contact       string            `json:"contact"`
	Owners        map[string]string `json:"owner"`
	Link          string            `json:"link"`
	LockedByAdmin bool              `json:"locked"`
	Tags          map[string]string `json:"tags"`
}

type CategoriesEtc struct {
	City      map[string]string `json:"cities"`
	Specialty map[string]string `json:"specialty"`
	Country   map[string]string `json:"country"`
}

type MsgPagination struct {
	Max     int    `json:"max"`
	UuidMsg string `json:"msgs"`
}

type PassChange struct {
	ConfirmPassword string `json:"confirmPassword"`
	NewPassword     string `json:"newPassword"`
	CurrentPassword string `json:"currentPassword"`
}
