# Lista de Tarefas

Aplicacao web de lista de tarefas com frontend em HTML/CSS/JavaScript, backend em Go e persistencia em MySQL.

## Estrutura

```text
frontend/   Interface web
backend/    API RESTful em Go
db/         Scripts de banco de dados
```

## Requisitos locais

- Go 1.22 ou superior
- Docker ou Podman para subir o MySQL padronizado

## Banco de dados com container

Suba o MySQL padronizado:

```powershell
docker compose up -d mysql
```

Ou, usando Podman:

```powershell
podman compose up -d mysql
```

O container usa MySQL 8.4, cria o banco `todo_list` e executa o script `db/task.sql` automaticamente.

## Variaveis de ambiente do backend

Crie o arquivo local de ambiente a partir do exemplo:

```powershell
copy backend\.env.example backend\.env
```

O backend carrega automaticamente o arquivo `backend/.env`. O exemplo ja esta alinhado com o container MySQL:

```env
DB_USER=root
DB_PASSWORD=todo_root
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=todo_list
CORS_ALLOWED_ORIGIN=http://localhost:8080
```

Se alguma variavel ja estiver definida no sistema, ela tem prioridade sobre o `.env`.

## Executar

Entre na pasta do backend, baixe as dependencias e suba a aplicacao:

```powershell
cd backend
go mod tidy
go run .
```

Abra no navegador:

```text
http://localhost:8080
```

## Endpoints

```text
GET    /api/v1/tasks
POST   /api/v1/tasks
PUT    /api/v1/tasks/{id}
DELETE /api/v1/tasks/{id}
```

## Observacao

O `compose.yaml` padroniza apenas o MySQL. A aplicacao Go continua rodando localmente durante o desenvolvimento.
