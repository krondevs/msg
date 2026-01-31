package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var nowords = map[string]struct{}{
	"la": {}, "le": {}, "el": {}, "los": {}, "las": {},
	"de": {}, "en": {}, "por": {}, "un": {}, "que": {},
	"a": {}, "ante": {}, "bajo": {}, "cabe": {}, "con": {},
	"contra": {}, "dentro": {}, "desde": {}, "entre": {}, "hacia": {},
	"hasta": {}, "para": {}, "segun": {}, "sobre": {}, "tras": {},
	"versus": {}, "al": {}, "del": {}, "lo": {}, "mi": {}, "mis": {},
	"tu": {}, "tus": {}, "su": {}, "sus": {}, "este": {}, "esta": {},
	"estos": {}, "estas": {}, "ese": {}, "esa": {}, "esos": {}, "esas": {},
	"nuestro": {}, "nuestra": {}, "vuestro": {}, "vuestra": {}, "todos": {},
	"todas": {}, "uno": {}, "una": {}, "unos": {}, "unas": {}, "como": {},
	"cuando": {}, "donde": {}, "porque": {}, "siempre": {}, "nunca": {},
	"aqui": {}, "alli": {}, "ya": {}, "todavia": {}, "solo": {}, "solamente": {},
	"casi": {}, "apenas": {}, "tambien": {}, "ademas": {}, "sin": {},
	"antes": {}, "despues": {}, "durante": {}, "encima": {}, "debajo": {},
	"traves": {}, "junto": {}, "cuanto": {}, "y": {}, "o": {},
	"e": {}, "i": {}, "u": {},
	"A": {}, "E": {}, "I": {}, "O": {}, "U": {},
	"b": {}, "c": {}, "d": {}, "f": {}, "g": {}, "h": {}, "j": {}, "k": {}, "l": {}, "m": {},
	"n": {}, "p": {}, "q": {}, "r": {}, "s": {}, "t": {}, "v": {}, "w": {}, "x": {}, "z": {},
	"B": {}, "C": {}, "D": {}, "F": {}, "G": {}, "H": {}, "J": {}, "K": {}, "L": {}, "M": {},
	"N": {}, "P": {}, "Q": {}, "R": {}, "S": {}, "T": {}, "V": {}, "W": {}, "X": {}, "Y": {}, "Z": {},
	"0": {}, "1": {}, "2": {}, "3": {}, "4": {},
	"5": {}, "6": {}, "7": {}, "8": {}, "9": {},
}

var wordRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func HashPassword(password string) (string, error) {
	// Genera el hash con un costo de 12 (recomendado para producción)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ValidatePassword valida si una contraseña coincide con su hash
func ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// FloatToString convierte un valor (float64, int, etc.) que esté dentro de la
// variable interface{} a su representación en cadena.
// Si el número es un float64 y no tiene parte decimal, devuelve solo la parte entera.
func FloatToString(v interface{}) string {
	switch val := v.(type) {

	case float64:
		// 1. si la parte fraccionaria es cero → entero
		if math.Mod(val, 1) == 0 {
			return fmt.Sprintf("%.0f", val)
		}
		// 2. número con decimales: usar la precisión mínima necesaria
		return strconv.FormatFloat(val, 'g', -1, 64)

	case float32:
		f := float64(val)
		if math.Mod(f, 1) == 0 {
			return fmt.Sprintf("%.0f", f)
		}
		return strconv.FormatFloat(f, 'g', -1, 32)

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(val)

	default:
		// Si no es un número, intentamos convertirlo con fmt.Sprintf
		return fmt.Sprintf("%v", val)
	}
}

func RateProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener la dirección IP del cliente
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = c.Request.RemoteAddr
		}
		now := UnixMillisecTime()

		dbMutex.Lock()
		next, exists := RATELIMIT[clientIP]
		if exists && next > now {
			dbMutex.Unlock()
			c.JSON(429, gin.H{
				"status":  "error",
				"message": "too many requests",
				"data":    "too many requests",
			})
			c.Abort()
			return
		}
		RATELIMIT[clientIP] = now + int64(RATELIMITMILISECS)
		dbMutex.Unlock()

		c.Set("ip", clientIP)
		c.Next()
	}
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"status": "error", "message": "token required", "data": "invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userUUID, err := ValidateJWT(tokenString, secret)
		if err != nil {
			fmt.Println(err)
			c.JSON(401, gin.H{"status": "error", "message": "token expired", "data": "invalid token"})
			c.Abort()
			return
		}
		/*now := UnixMillisecTime()

		dbMutex.Lock()
		next, exists := RATELIMIT[userUUID]
		if exists && next > now {
			dbMutex.Unlock()
			c.JSON(429, gin.H{
				"status":  "error",
				"message": "too many requests",
				"data":    "too many requests",
			})
			c.Abort()
			return
		}
		RATELIMIT[userUUID] = now + int64(RATELIMITMILISECS)
		dbMutex.Unlock()*/

		c.Set("uuid", userUUID)
		c.Next()
	}
}

func ShuffleStrings(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Tokenize(s string) []string {
	// 1. Reemplazar cualquier carácter no alfanumérico por un espacio.
	s = wordRe.ReplaceAllString(s, " ")

	// 2. Normalizar a minúsculas y eliminar espacios iniciales/finales.
	pp := strings.ToLower(strings.TrimSpace(s))

	// 3. Dividir en tokens (palabras).
	rawTokens := strings.Fields(pp)

	// 4. Filtrar stop‑words y construir el slice resultante.
	newTT := make([]string, 0, len(rawTokens))
	for _, tok := range rawTokens {
		if _, isStop := nowords[tok]; !isStop {
			newTT = append(newTT, tok)
		}
	}
	return newTT
}

func normalizeString(s string) string {
	var result strings.Builder

	for _, r := range s {
		switch r {
		case 'á', 'à', 'â', 'ä', 'ã', 'å', 'æ':
			result.WriteRune('a')
		case 'Á', 'À', 'Â', 'Ä', 'Ã', 'Å', 'Æ':
			result.WriteRune('A')
		case 'é', 'è', 'ê', 'ë', 'ẽ', 'ę', 'ė':
			result.WriteRune('e')
		case 'É', 'È', 'Ê', 'Ë', 'Ẽ', 'Ę', 'Ė':
			result.WriteRune('E')
		case 'í', 'ì', 'î', 'ï', 'ĩ', 'į', 'ı':
			result.WriteRune('i')
		case 'Í', 'Ì', 'Î', 'Ï', 'Ĩ', 'Į', 'I':
			result.WriteRune('I')
		case 'ó', 'ò', 'ô', 'ö', 'õ', 'ø', 'œ':
			result.WriteRune('o')
		case 'Ó', 'Ò', 'Ô', 'Ö', 'Õ', 'Ø', 'Œ':
			result.WriteRune('O')
		case 'ú', 'ù', 'û', 'ü', 'ũ', 'ų', 'ū':
			result.WriteRune('u')
		case 'Ú', 'Ù', 'Û', 'Ü', 'Ũ', 'Ų', 'Ū':
			result.WriteRune('U')
		case 'ç', 'ć', 'č':
			result.WriteRune('c')
		case 'Ç', 'Ć', 'Č':
			result.WriteRune('C')
		case 'ñ', 'ń', 'ņ', 'ň':
			result.WriteRune('n')
		case 'Ñ', 'Ń', 'Ņ', 'Ň':
			result.WriteRune('N')
		case 'ý', 'ỳ', 'ŷ', 'ÿ', 'ỹ', 'ȳ':
			result.WriteRune('y')
		case 'Ý', 'Ỳ', 'Ŷ', 'Ÿ', 'Ỹ', 'Ȳ':
			result.WriteRune('Y')
		case 'ß':
			result.WriteRune('s')
		case 'Đ', 'Ð', 'Ď':
			result.WriteRune('D')
		case 'ł', 'ľ', 'ĺ':
			result.WriteRune('l')
		case 'Ł', 'Ľ', 'Ĺ':
			result.WriteRune('L')
		case 'ř', 'ŕ', 'ŗ':
			result.WriteRune('r')
		case 'Ř', 'Ŕ', 'Ŗ':
			result.WriteRune('R')
		case 'š', 'ś', 'ş':
			result.WriteRune('s')
		case 'Š', 'Ś', 'Ş':
			result.WriteRune('S')
		case 'ž', 'ź', 'ż':
			result.WriteRune('z')
		case 'Ž', 'Ź', 'Ż':
			result.WriteRune('Z')
		case 'þ', 'ŧ':
			result.WriteRune('t')
		case 'Þ', 'Ŧ':
			result.WriteRune('T')
		case 'ħ':
			result.WriteRune('h')
		case 'Ħ':
			result.WriteRune('H')
		default:
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
				result.WriteRune(r)
			} else {
				// Reemplazar otros caracteres especiales por guión o espacio
				if r == '-' || r == '_' || r == '.' {
					result.WriteRune(r)
				} else {
					result.WriteRune(' ')
				}
			}
		}
	}
	normalized := strings.TrimSpace(result.String())
	for strings.Contains(normalized, "  ") {
		normalized = strings.ReplaceAll(normalized, "  ", " ")
	}

	return normalized
}

func ValidateJWT(tokenString, secret string) (string, error) {
	// Clave secreta para validar el token (debe ser la misma que se usó para generar)
	secretKey := []byte(secret)

	// Parsear y validar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Verificar que el algoritmo de firma sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("error al validar token: %w", err)
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return "", fmt.Errorf("token inválido")
	}

	// Extraer las claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("no se pudieron extraer las claims del token")
	}

	// Extraer el UUID del usuario
	userUUID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id no encontrado o inválido en el token")
	}

	return userUUID, nil
}

func GenerateJWT(userUUID, secret string, expirationHours int) (string, error) {
	// Clave secreta para firmar el token (en producción debe estar en variables de entorno)
	secretKey := []byte(secret)

	// Crear las claims del token
	claims := jwt.MapClaims{
		"user_id": userUUID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Crear el token con el algoritmo HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return tokenString, nil
}

func Encrypt(plaintext string, key string) (string, error) {
	keyBytes := []byte(key)
	if !isValidAESKeyLen(len(keyBytes)) {
		return "", fmt.Errorf("invalid key length: %d", len(keyBytes))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string, key string) (string, error) {
	keyBytes := []byte(key)
	if !isValidAESKeyLen(len(keyBytes)) {
		return "", fmt.Errorf("invalid key length: %d", len(keyBytes))
	}

	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, data := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func isValidAESKeyLen(n int) bool {
	return n == 16 || n == 24 || n == 32
}

func StringToHexa(s string) string {
	return hex.EncodeToString([]byte(s))
}

func SetupLogger() {
	// Crear o abrir el archivo de logs
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err)
	}

	// Configurar el logger para escribir en el archivo
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func errorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Procesar la solicitud
		c.Next()

		// Comprobar si hubo errores
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Println(err.Error())
			}
		}
	}
}

func GenerateRandomKey32() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(errorLogger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	return r
}

func RequestBearer(method, url, contentType string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

func HexaToString(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NumToInt64(i interface{}) int64 {
	switch v := i.(type) {
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case string:
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return num
		}
		return 0
	default:
		return 0
	}
}

func GetRandom(bajo, alto int) int {
	return rand.Intn(alto-bajo) + bajo // Generar el número aleatorio
}

func GetValues(m map[string]interface{}) []interface{} {
	values := reflect.ValueOf(m)
	vals := make([]interface{}, values.Len())

	for i, key := range values.MapKeys() {
		vals[i] = values.MapIndex(key).Interface()
	}

	return vals
}

func GetKeys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func StringToNum(valor string) (interface{}, error) {
	// Intentar convertir a entero
	if numeroInt, err := strconv.Atoi(valor); err == nil {
		return numeroInt, nil
	}

	// Intentar convertir a flotante
	if numeroFloat, err := strconv.ParseFloat(valor, 64); err == nil {
		return numeroFloat, nil
	}

	// Si ninguna conversión es exitosa, devolver un error
	return nil, fmt.Errorf("no se pudo convertir la cadena: %s", valor)
}

func NumToString(numero interface{}) string {
	switch v := numero.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 32)
	default:
		return "0"
	}
}

func Str(numero interface{}) string {
	return NumToString(numero)
}

func Int(number interface{}) int64 {
	return NumToInt64(number)
}

func Float(number interface{}) float64 {
	return NumToFloat64(number)
}

func Format(number any, prec int) string {
	nn := Float(number)
	kk := Int(nn)
	pp := Float(kk)
	gg := nn - pp
	if gg == 0 {
		return strconv.FormatInt(kk, 10)
	}
	return strconv.FormatFloat(nn, 'f', prec, 64)
}

func NumToFloat64(i interface{}) float64 {
	switch v := i.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
		return 0
	default:
		return 0
	}
}

func PrintRed(text ...string) {
	fondoRojo := color.New(color.FgWhite, color.BgRed)
	resultado := strings.Join(text, " ")
	fondoRojo.Println(resultado)
}

func PrintGreen(text ...string) {
	fondoVerde := color.New(color.FgBlack, color.BgGreen)
	resultado := strings.Join(text, " ")
	fondoVerde.Println(resultado)
}

func PrintYellow(text ...string) {
	fondoAmarillo := color.New(color.FgBlack, color.BgYellow)
	resultado := strings.Join(text, " ")
	fondoAmarillo.Println(resultado)
}

func PrintBlue(text ...string) {
	fondoAzul := color.New(color.FgWhite, color.BgBlue)
	resultado := strings.Join(text, " ")
	fondoAzul.Println(resultado)
}

func Input(prompt string) string {
	fmt.Print(prompt)                     // Imprimir el mensaje de aviso
	scanner := bufio.NewScanner(os.Stdin) // Crear un nuevo scanner
	scanner.Scan()                        // Leer la entrada del usuario
	return scanner.Text()                 // Devolver el texto ingresado
}

func Pause() {
	fmt.Print("Presiona Enter para salir...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Espera a que el usuario presione Enter
}

func CreateTorProxy(port string) (*http.Client, error) {
	proxyURL, err := url.Parse("socks5h://localhost:" + port)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 30 * time.Second,
	}
	return client, nil
}

func HexDecode(hexStr string) string {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return hexStr // Si hay error, devolver el string original
	}
	return string(decoded)
}

func Hoy() string {
	return time.Now().Format("2006-01-02")
}

func Manana() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05") // Formato de Go
}

func HourNow() string {
	return time.Now().Format("15:04:05") // Formato de Go
}

func UnixTime() int64 {
	return time.Now().Unix() // Devuelve el tiempo en segundos desde epoch
}

func UnixMillisecTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond) // Devuelve el tiempo en milisegundos desde epoch
}

func GenerarFechasConsecutivas(fechaInicio, fechaFin string) ([]string, error) {
	// Definir el formato de fecha
	formatoFecha := "2006-01-02"

	// Parsear las fechas de entrada
	inicio, err := time.Parse(formatoFecha, fechaInicio)
	if err != nil {
		return nil, fmt.Errorf("error al parsear fecha de inicio: %v", err)
	}

	fin, err := time.Parse(formatoFecha, fechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al parsear fecha de fin: %v", err)
	}

	// Verificar que la fecha de inicio sea anterior o igual a la fecha de fin
	if inicio.After(fin) {
		return nil, fmt.Errorf("la fecha de inicio debe ser anterior o igual a la fecha de fin")
	}

	// Slice para almacenar las fechas
	var fechas []string

	// Generar fechas consecutivas
	for fecha := inicio; !fecha.After(fin); fecha = fecha.AddDate(0, 0, 1) {
		fechas = append(fechas, fecha.Format(formatoFecha))
	}

	return fechas, nil
}

func UnixToDate(epoch int64) string {
	// Convertir el epoch a un objeto Time
	t := time.Unix(epoch, 0)
	// Formatear la fecha y hora
	return t.Format("2006-01-02")
}

func DiaHoy() string {
	fechaActual := time.Now()
	nombreDia := fechaActual.Weekday().String()

	diasSemana := map[string]string{
		"Monday":    "Lunes",
		"Tuesday":   "Martes",
		"Wednesday": "Miércoles",
		"Thursday":  "Jueves",
		"Friday":    "Viernes",
		"Saturday":  "Sábado",
		"Sunday":    "Domingo",
	}

	return diasSemana[nombreDia]
}

func GetHash(input string) string {
	// Crear un nuevo hash SHA-256
	hash := sha256.New()
	// Escribir el string en el hash
	hash.Write([]byte(input))
	// Obtener el resultado del hash
	return hex.EncodeToString(hash.Sum(nil))
}

type RequestLmdb struct {
	Query     string `json:"query"`
	KeyStore  string `json:"key"`
	Values    any    `json:"values"`
	Database  string `json:"db"`
	MasterKey string `json:"encrypt"`
}

type Resss struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Result     any    `json:"result"`
	StatusCode int
	ResultByte []byte
}

func Requests(method, url, contentType, bearerToken string, payload []byte, tor bool, port string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, 500, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	client := &http.Client{}
	if tor {
		client, err = CreateTorProxy(port)
		if err != nil {
			return nil, 500, err
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, 500, err
	}
	if resp.StatusCode == 500 {
		return nil, resp.StatusCode, fmt.Errorf("inertnal server error")
	}
	if resp.StatusCode == 400 {
		return nil, resp.StatusCode, fmt.Errorf("invalid request")
	}
	if resp.StatusCode == 401 {
		return nil, resp.StatusCode, fmt.Errorf("unautorized")
	}
	if resp.StatusCode == 429 {
		return nil, resp.StatusCode, fmt.Errorf("too many requests")
	}
	if resp.StatusCode == 404 {
		return nil, resp.StatusCode, fmt.Errorf("not found")
	}
	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf(resp.Status)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, 500, err
	}
	return body, 200, nil
}

func QueryBadger(queryType, key string, values any) (Resss, error) {
	if len(key) < 1 || key == "" {
		return Resss{}, fmt.Errorf("invalid key to make requests")
	}
	var datos = Resss{}
	envio := RequestLmdb{}
	envio.Database = SETTINGS["dbname"].(string)
	envio.KeyStore = key
	envio.Values = values
	envio.Query = queryType
	envio.MasterKey = ""
	jsonData, _ := json.Marshal(envio)
	resp, ss, err := Requests("POST", SETTINGS["dbhost"].(string), "application/json", SETTINGS["dbkey"].(string), jsonData, SETTINGS["tor"].(bool), SETTINGS["torport"].(string))
	if ss == 404 {
		datos.StatusCode = 404
		return datos, nil
	}
	if ss != 200 {
		err = json.Unmarshal(resp, &datos)
		if err != nil {
			return Resss{}, err
		}
		return Resss{}, fmt.Errorf(datos.Message)
	}
	if err != nil {
		return Resss{}, err
	}
	err = json.Unmarshal(resp, &datos)
	if err != nil {
		return Resss{}, err
	}
	datos.StatusCode = ss
	datos.ResultByte, err = json.Marshal(datos.Result)
	if err != nil {
		fmt.Println(err)
		return datos, err
	}
	return datos, nil
}
