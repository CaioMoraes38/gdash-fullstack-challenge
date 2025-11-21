# 🌦️ GDASH Challenge 2025/02 — Weather Intelligence Platform

Uma plataforma full-stack completa para coleta, processamento, análise e visualização de dados climáticos em tempo real.
O projeto integra múltiplas linguagens (Python, Go, TypeScript) e serviços orquestrados via Docker, seguindo uma arquitetura moderna orientada a eventos e microsserviços.

## 🏗️  Arquitetura e Pipeline de Dados

O sistema foi projetado com foco em desacoplamento, resiliência e escalabilidade:

Coleta (Python)
Script agendado consome a API Open-Meteo, normaliza os dados da cidade alvo (ex: Birigui) e publica mensagens na fila.

Mensageria (RabbitMQ)
Garante desacoplamento total entre coleta e processamento usando AMQP.

Processamento (Go Worker)
Consome a fila com alta performance, valida os dados e publica na API.
Implementa Ack/Nack garantindo tolerância a falhas.

Núcleo da Plataforma (NestJS + MongoDB)
API RESTful que centraliza regras de negócio, persistência, autenticação e geração de insights.

Interface (React + shadcn/ui)
Dashboard moderno para visualização dos dados em tempo real.

🔁 Diagrama de Fluxo
[Python Collector] --(JSON)--> [RabbitMQ] --(AMQP)--> [Go Worker] --(HTTP)--> [NestJS API] --(Mongoose)--> [MongoDB]
                                                                                     ^
                                                                                     |
                                                                              [React Frontend]

## 🚀  Tecnologias Utilizadas
## Infraestrutura

Docker & Docker Compose

## Coleta de Dados

Python 3.10

Requests, Schedule

Message Broker

RabbitMQ (com painel administrativo)

## Worker

Go 1.23

AMQP 0.9.1

## Backend

NestJS

TypeScript

Mongoose

JWT

## Banco de Dados

MongoDB

Frontend

React + Vite

TailwindCSS

shadcn/ui

## ⚙️  Como Executar
✔️ Pré-requisitos

Docker e Docker Compose instalados

Git instalado

▶️ Execução Recomendada (Docker Compose)

Clone o repositório

git clone https://github.com/CaioMoraes38/gdash-challenge.git
cd gdash-challenge


Configure variáveis de ambiente

Renomeie .env.example → .env (se existir)

Confirme se o docker-compose.yml está correto

Suba todos os serviços

docker-compose up --build


Acesse as aplicações

Serviço	URL	Credenciais (padrão)
Frontend (Dashboard)	http://localhost:5173
	-
API (Backend)	http://localhost:3000
	-
RabbitMQ (Admin)	http://localhost:15672
	user / yourUser
Swagger Docs	http://localhost:3000/api
	-
🛠️ Execução Manual (Ambiente de Desenvolvimento)
1. Subir infraestrutura
docker-compose up -d mongodb rabbitmq

2. Backend (NestJS)
cd backend-api
npm install
npm run start:dev


Usuário padrão criado automaticamente:
admin@example.com / 123456

3. Worker (Go)
cd weather-worker
go mod tidy
go run main.go

4. Coletor (Python)
cd weather-collector
pip install -r requirements.txt
python main.py

## 🔌Endpoints Principais (API)
Método	Endpoint	Descrição
GET	/weather	Lista histórico climático
POST	/weather	Recebe dados do worker (uso interno)
GET	/weather/export/csv	Exporta dados em CSV
POST	/auth/login	Autenticação
GET	/users	CRUD de usuários (protegido)
🧠 Insights Automáticos de IA

O sistema gera insights com base nas leituras mais recentes, incluindo:

Tendência de alta/baixa de temperatura

Alertas de baixa umidade ou calor extremo

Cálculo de médias móveis para previsão de curto prazo

## 📝 Decisões de Projeto

## Monorepo
Simplifica execução e avaliação via Docker em um único comando.

## Worker em Go
Melhor desempenho para consumo de fila e validação em background.

## Backend em NestJS
Tipagem forte, testes fáceis e padrão corporativo escalável.

👤 Autor

Desenvolvido por Caio de Moraes Santos.


