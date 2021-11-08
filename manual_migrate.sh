# start cloud_sql_proxy in separate terminal
# cloud_sql_proxy -instances=tanzania-330508:europe-central2:webapp-db=tcp:0.0.0.0:5432
export DB_NAME=$(gcloud secrets versions access latest --secret="db-name")
export DB_PASS=$(gcloud secrets versions access latest --secret="db-password")
export DB_USER=$(gcloud secrets versions access latest --secret="db-username")
export POSTGRESQL_URL="postgres://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path migrations up