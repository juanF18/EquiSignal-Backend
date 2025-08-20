package db

import (
	"fmt"
	"log"

	"github.com/juanF18/EquiSignal-Backend/internal/config"
	"github.com/juanF18/EquiSignal-Backend/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect abre la conexión a CockroachDB con GORM
func ConnectCockroachDB(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=verify-full",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error conectando a CockroachDB: ", err)
	}

	err = db.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Fatal("❌ Error al migrar la base de datos: ", err)
	}

	DB = db
	log.Println("✅ Conectado a CockroachDB con éxito")
}
