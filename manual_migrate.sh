export DB_NAME=$(gcloud secrets versions access latest --secret="db-name")
export DB_PASS=$(gcloud secrets versions access latest --secret="db-password")
export DB_USER=$(gcloud secrets versions access latest --secret="db-username")
export POSTGRESQL_URL="postgres://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path migrations up
