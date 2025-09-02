# EquiSignal-Backend 📈

**EquiSignal-Backend** es un sistema de recomendaciones de acciones desarrollado en Go, que utiliza Clean Architecture para proporcionar un backend robusto y escalable. El sistema analiza datos de acciones de fuentes externas y genera recomendaciones inteligentes basadas en ratings, precios objetivo y comportamientos del mercado.

## 🏗️ Arquitectura del Proyecto

El proyecto sigue los principios de **Clean Architecture**, separando claramente las responsabilidades en diferentes capas:

```
EquiSignal-Backend/
├── cmd/                          # Punto de entrada de la aplicación
│   └── app/
│       └── main.go              # Configuración principal y arranque del servidor
├── internal/                     # Código interno de la aplicación
│   ├── algorithms/              # Algoritmos de negocio
│   │   └── stock/
│   │       └── recommender.go   # Sistema de recomendaciones de acciones
│   ├── application/             # Capa de aplicación (casos de uso)
│   │   └── stock_service.go     # Servicios de lógica de negocio
│   ├── config/                  # Configuración de la aplicación
│   │   └── config.go           # Manejo de variables de entorno
│   ├── domain/                  # Entidades de dominio
│   │   └── models/
│   │       └── stock.go        # Modelo de datos para acciones
│   ├── infrastructure/          # Capa de infraestructura
│   │   ├── db/
│   │   │   └── cockroachdb.go  # Conexión a CockroachDB
│   │   └── repository/
│   │       └── stock_repository.go # Acceso a datos
│   └── interface/              # Capa de interfaz (adapters)
│       ├── dto/
│       │   └── stock_dto.go    # Data Transfer Objects
│       ├── external/
│       │   └── external_api.go # Integración con APIs externas
│       ├── handlers/
│       │   └── stock_handler.go # Controladores HTTP
│       └── http/
│           ├── external_routes.go
│           ├── routes.go       # Configuración de rutas principales
│           └── stock_routes.go # Rutas específicas para acciones
├── go.mod                      # Dependencias del módulo Go
├── go.sum                      # Checksums de dependencias
└── README.md                   # Documentación del proyecto
```

## 🚀 Características Principales

- **Sistema de Recomendaciones Inteligente**: Algoritmo que evalúa acciones basado en ratings, precios objetivo y temporalidad
- **API RESTful**: Interfaz HTTP robusta construida con Gin Framework
- **Base de Datos Distribuida**: Integración con CockroachDB para alta disponibilidad
- **Integración con APIs Externas**: Consumo de datos de mercado en tiempo real
- **CORS Configurado**: Soporte para aplicaciones frontend
- **Clean Architecture**: Código mantenible y testeable

## 🛠️ Stack Tecnológico

- **Lenguaje**: Go 1.24.6
- **Framework Web**: Gin Gonic
- **Base de Datos**: CockroachDB (PostgreSQL compatible)
- **ORM**: GORM
- **Configuración**: Godotenv para variables de entorno
- **CORS**: Gin-contrib/cors
- **UUIDs**: Google UUID

## 📋 Prerrequisitos

- Go 1.24.6 o superior
- CockroachDB o PostgreSQL
- Variables de entorno configuradas

## 🔧 Instalación y Configuración

1. **Clonar el repositorio**:

   ```bash
   git clone https://github.com/juanF18/EquiSignal-Backend.git
   cd EquiSignal-Backend
   ```

2. **Instalar dependencias**:

   ```bash
   go mod download
   ```

3. **Configurar variables de entorno**:
   Crear un archivo `.env` en la raíz del proyecto:

   ```env
   # Database
   DATABASE_URL=postgresql://username:password@host:port/database

   # API Configuration
   FRONTEND_URL=http://localhost:3000
   PORT=8080

   # External APIs
   EXTERNAL_API_KEY=your_api_key_here
   ```

4. **Ejecutar la aplicación**:
   ```bash
   go run cmd/app/main.go
   ```

## 🏛️ Descripción de Capas

### 1. **Capa de Dominio** (`domain/`)

- **Responsabilidad**: Define las entidades de negocio y reglas fundamentales
- **Componentes**:
  - `models/stock.go`: Entidad Stock con todos sus atributos (ticker, company, brokerage, action, ratings, etc.)

### 2. **Capa de Aplicación** (`application/`)

- **Responsabilidad**: Contiene la lógica de negocio y casos de uso
- **Componentes**:
  - `stock_service.go`: Servicios que orquestan las operaciones con acciones

### 3. **Capa de Infraestructura** (`infrastructure/`)

- **Responsabilidad**: Implementa detalles técnicos (base de datos, APIs externas)
- **Componentes**:
  - `db/cockroachdb.go`: Configuración y conexión a la base de datos
  - `repository/stock_repository.go`: Implementación de acceso a datos

### 4. **Capa de Interfaz** (`interface/`)

- **Responsabilidad**: Maneja la comunicación externa (HTTP, APIs)
- **Componentes**:
  - `handlers/stock_handler.go`: Controladores HTTP
  - `external/external_api.go`: Integración con APIs externas
  - `dto/stock_dto.go`: Objetos de transferencia de datos
  - `http/routes.go`: Configuración de rutas

### 5. **Algoritmos** (`algorithms/`)

- **Responsabilidad**: Contiene la lógica de recomendaciones
- **Componentes**:
  - `stock/recommender.go`: Sistema de scoring para recomendaciones de acciones

## 📊 Sistema de Recomendaciones

El algoritmo de recomendaciones evalúa las acciones basándose en:

- **Ratings**: Strong Buy (+4), Buy (+3), Hold (+0)
- **Precios Objetivo**: Diferencia entre precio actual y objetivo
- **Temporalidad**: Prioriza recomendaciones más recientes
- **Brokerage**: Considera la fuente de la recomendación

### Ejemplo de Scoring:

```go
type StockRecommendation struct {
    Ticker     string    // Símbolo de la acción
    Company    string    // Nombre de la empresa
    Score      int       // Puntaje calculado
    Reason     string    // Justificación del puntaje
    Rating     string    // Rating asignado
    TargetFrom string    // Precio objetivo inicial
    TargetTo   string    // Precio objetivo final
    Time       time.Time // Timestamp de la recomendación
}
```

## 🔌 Endpoints API

### Salud del Sistema

- `GET /health` - Verificar el estado del servidor

### Acciones (Stock Routes)

- Endpoints definidos en `internal/interface/http/stock_routes.go`

## 🗄️ Modelo de Datos

### Entidad Stock

```go
type Stock struct {
    ID         uuid.UUID // Identificador único
    Ticker     string    // Símbolo de la acción (ej: AAPL)
    Company    string    // Nombre de la empresa
    Brokerage  string    // Casa de corretaje
    Action     string    // Acción recomendada
    RatingFrom string    // Rating inicial
    RatingTo   string    // Rating actualizado
    TargetFrom string    // Precio objetivo inicial
    TargetTo   string    // Precio objetivo actualizado
    Time       time.Time // Timestamp de la recomendación
    CreatedAt  time.Time // Fecha de creación
    UpdatedAt  time.Time // Fecha de actualización
}
```

## 🚀 Desarrollo

### Estructura de Branches

- `main`: Rama principal de producción
- `development`: Rama de desarrollo activa

### Comandos Útiles

```bash
# Ejecutar tests
go test ./...

# Construir la aplicación
go build -o bin/equisignal cmd/app/main.go

# Ejecutar con hot reload (requiere air)
air

# Formatear código
go fmt ./...
```

## 🤝 Contribuciones

1. Fork el repositorio
2. Crear una rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit los cambios (`git commit -am 'Add nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear un Pull Request

## 👥 Equipo

- **Desarrollador Principal**: [juanF18](https://github.com/juanF18)

## 🔗 Enlaces Relacionados

- [CockroachDB Documentation](https://www.cockroachlabs.com/docs/)
- [Gin Framework](https://gin-gonic.com/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**EquiSignal-Backend** - Sistema inteligente de recomendaciones de acciones 📈
