package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Rumm1/eduhub-backend/internal/platform/password"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	email := getEnv("SUPER_ADMIN_EMAIL", "superadmin@eduhub.kz")
	plainPassword := getEnv("SUPER_ADMIN_PASSWORD", "SuperAdmin123!")
	fullName := getEnv("SUPER_ADMIN_FULL_NAME", "EduHub Super Admin")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	roleID, err := getOrCreateSuperAdminRole(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}

	userID, err := getOrCreateSuperAdminUser(ctx, tx, email, plainPassword, fullName)
	if err != nil {
		log.Fatal(err)
	}

	if err := assignRoleToUser(ctx, tx, userID, roleID); err != nil {
		log.Fatal(err)
	}

	if err := assignAllPermissionsToRole(ctx, tx, roleID); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seed completed successfully")
	fmt.Println("SUPER_ADMIN email:", email)
	fmt.Println("SUPER_ADMIN password:", plainPassword)
}

func getOrCreateSuperAdminRole(ctx context.Context, tx pgx.Tx) (uuid.UUID, error) {
	var roleID uuid.UUID

	err := tx.QueryRow(ctx, `
SELECT id
FROM roles
WHERE organization_id IS NULL AND code = 'SUPER_ADMIN'
LIMIT 1
`).Scan(&roleID)

	if err == nil {
		return roleID, nil
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	roleID = uuid.New()

	_, err = tx.Exec(ctx, `
INSERT INTO roles (
id,
organization_id,
name,
code,
description,
is_system
)
VALUES ($1, NULL, $2, $3, $4, true)
`, roleID, "Super Admin", "SUPER_ADMIN", "Platform super administrator")

	if err != nil {
		return uuid.Nil, err
	}

	return roleID, nil
}

func getOrCreateSuperAdminUser(ctx context.Context, tx pgx.Tx, email string, plainPassword string, fullName string) (uuid.UUID, error) {
	var userID uuid.UUID

	err := tx.QueryRow(ctx, `
SELECT id
FROM users
WHERE organization_id IS NULL AND email = $1
LIMIT 1
`, email).Scan(&userID)

	if err == nil {
		return userID, nil
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	hashedPassword, err := password.Hash(plainPassword)
	if err != nil {
		return uuid.Nil, err
	}

	userID = uuid.New()

	_, err = tx.Exec(ctx, `
INSERT INTO users (
id,
organization_id,
email,
password_hash,
full_name,
status
)
VALUES ($1, NULL, $2, $3, $4, 'active')
`, userID, email, hashedPassword, fullName)

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func assignRoleToUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, roleID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, userID, roleID)

	return err
}

func assignAllPermissionsToRole(ctx context.Context, tx pgx.Tx, roleID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
INSERT INTO role_permissions (role_id, permission_id)
SELECT $1, id
FROM permissions
ON CONFLICT DO NOTHING
`, roleID)

	return err
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
