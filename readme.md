# Chat Application Backend

## Descripción

Este es un backend completo para una aplicación de chat en tiempo real con funcionalidades avanzadas como gestión de grupos, mensajería, autenticación y control de acceso. Utiliza Go con el framework Gin para crear una API RESTful robusta.

## Características Principales

- **Autenticación JWT**: Sistema seguro de inicio de sesión y gestión de sesiones
- **Gestión de Grupos**: Creación, edición, cierre y eliminación de grupos de chat
- **Mensajería en Tiempo Real**: Envío y recepción de mensajes con verificación de lectura
- **Sistema de Archivos**: Subida y descarga de archivos multimedia (imágenes, videos, documentos)
- **Control de Acceso**: Suspensión, expulsión y gestión de usuarios dentro de grupos
- **Moderación**: Denuncias y bloqueo de grupos por administradores
- **Tipado en Tiempo Real**: Indicadores de escritura en tiempo real

## Estructura del Proyecto

```
.
├── main.go                 # Punto de entrada principal
├── settings.json          # Configuración de la aplicación
├── templates/             # Plantillas HTML
│   ├── search.html
│   ├── login.html
│   └── index.html
├── static/                # Archivos estáticos
│   └── uploads/           # Directorio para archivos subidos
└── README.md              # Este archivo
```

## Configuración

### Archivo settings.json

```json
{
  "port": "8080",
  "jwt": "tu_secreto_jwt_aqui",
  "duration": 3600,
  "database": "badger"
}
```

## Endpoints API

### Autenticación
- `POST /signin` - Iniciar sesión
- `POST /register` - Registrar nuevo usuario

### Mensajería
- `POST /api/sendTyping` - Enviar indicador de escritura
- `POST /api/verifyTyping` - Verificar si alguien está escribiendo
- `POST /api/getChats` - Obtener lista de chats del usuario
- `POST /api/validateSession` - Validar sesión activa
- `POST /api/getGroupsList` - Obtener lista de grupos del usuario
- `POST /api/addNewGroup` - Crear nuevo grupo
- `POST /api/loadGroupChat` - Cargar chat de grupo
- `POST /api/sendMessage` - Enviar mensaje
- `POST /api/join` - Unirse a un grupo
- `POST /api/cerrarGrupo` - Cerrar/abrir grupo
- `POST /api/uploadFile` - Subir archivo
- `POST /api/suspender` - Suspender usuario en grupo
- `POST /api/expulsar` - Expulsar usuario de grupo
- `POST /api/eliminarMsg` - Eliminar mensaje
- `POST /api/salirGrupo` - Salir del grupo

## Funcionalidades Especiales

### Gestión de Usuarios
- Registro con verificación KYC (Know Your Customer)
- Sistema de roles (Usuario, Administrador, Superusuario)
- Control de acceso por estado de cuenta

### Seguridad
- Hashing de contraseñas con bcrypt
- Validación de enlaces, correos y teléfonos en mensajes
- Protección contra spam automático
- Rate limiting para prevenir abusos

### Interfaz de Usuario
- Plantillas HTML responsivas
- Página de búsqueda y login
- Interfaz de chat completa con historial
- Gestión de grupos y usuarios

## Requisitos del Sistema

- Go 1.19+
- Gin Framework
- BadgerDB (para almacenamiento)
- JWT para autenticación
- bcrypt para hash de contraseñas

## Instalación

1. Clonar el repositorio:
```bash
git clone <url-del-repositorio>
cd chat-app-backend
```

2. Instalar dependencias:
```bash
go mod tidy
```

3. Configurar archivo `settings.json` con los parámetros adecuados

4. Ejecutar la aplicación:
```bash
go run main.go
```

## Uso

1. Acceder a la aplicación en `http://localhost:8080`
2. Registrar un nuevo usuario o iniciar sesión con las credenciales del administrador (admin/123456)
3. Crear grupos de chat y comenzar a chatear
4. Utilizar el sistema de archivos para compartir multimedia

## Notas Importantes

- El sistema incluye un usuario administrador preconfigurado
- Todos los mensajes se almacenan con timestamps y estado de lectura
- Se implementa protección contra enlaces/spam automáticos
- El sistema maneja múltiples usuarios por grupo con diferentes estados (activo, suspendido)
- Se pueden crear grupos por ciudad, país y especialidad

## Licencia

Este proyecto está bajo la licencia MIT. 
