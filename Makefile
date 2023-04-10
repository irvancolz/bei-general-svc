run_migrate:
	go run ./db/migrate/migrate.go

run_seed:
	psql -U postgres postgres < db/seed.sql