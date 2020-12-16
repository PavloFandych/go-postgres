CREATE TABLE users (userid SERIAL PRIMARY KEY, name TEXT, age INT, location TEXT);
CREATE USER go_postgres_user WITH PASSWORD 'go_postgres_password';
GRANT ALL PRIVILEGES ON DATABASE "go_postgres" to go_postgres_user;
GRANT ALL PRIVILEGES ON TABLE "users" to go_postgres_user;
GRANT USAGE, SELECT ON SEQUENCE "users_userid_seq" to go_postgres_user;

curl http://localhost:8080/api/user
curl -d '{"name": "gopher","age":25,"location":"India"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/newuser