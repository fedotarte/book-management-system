package database

import (
	"book-management-system/config"
	"log"
	"os/exec"
)

// RunMigrations применяет миграции при старте сервера
func RunMigrations() {
	dsn := config.GetEnv("DATABASE_URL", "host=localhost user=postgres password=secret dbname=books port=5432 sslmode=disable")

	cmd := exec.Command("migrate", "-database", dsn, "-path", "migrations", "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Ошибка запуска миграций: %s\n%s", err, string(output))
	}

	log.Println("Миграции успешно применены!")
}
