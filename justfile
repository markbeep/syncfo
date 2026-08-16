# https://just.systems

DB_PATH := `yq '.database.path' configs/public.toml`
MIGRATION_DIR := `yq '.database.migrations' configs/public.toml`

upgrade:
    mkdir -p .data
    GOOSE_MIGRATION_DIR={{ MIGRATION_DIR }} goose sqlite3 {{ DB_PATH }} up

downgrade:
    mkdir -p .data
    GOOSE_MIGRATION_DIR={{ MIGRATION_DIR }} goose sqlite3 {{ DB_PATH }} down
