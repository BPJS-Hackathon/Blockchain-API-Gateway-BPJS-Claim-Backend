# Blockchain-API-Gateway-BPJS-Claim-Backend
# Blockchain-API-Gateway-BPJS-Claim-Backend

# BPJS Blockchain Gateway (Developer 3)

ENV (example):
DATABASE_URL=postgres://user:pass@localhost:5432/bpjs?sslmode=disable
JWT_SECRET=donutsharethissecret
PORT=8080

go run ./app

Endpoints:
POST /api/login { username, password }
POST /api/claims/submit
POST /api/claims/approve (Authorization: Bearer <token>)
GET  /api/claims/:id
GET  /api/claims?status=submitted
