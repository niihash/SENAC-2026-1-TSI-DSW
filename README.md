# Lista de Tarefas

Aplicacao web de lista de tarefas com frontend em HTML, CSS e JavaScript, backend em Go e persistencia em MySQL.

## Estrutura do projeto

```text
backend/    API REST em Go
db/         Scripts SQL do banco de dados
frontend/   Interface web
```

## Requisitos

- Go 1.22 ou superior
- Docker ou Podman
- Git

## Configuracao do ambiente

O backend usa variaveis de ambiente para se conectar ao banco. O arquivo real de ambiente deve existir apenas na maquina local e nao deve ser commitado.

Crie o arquivo local a partir do exemplo versionado:

```powershell
copy backend\.env.example backend\.env
```

No Git Bash, use:

```bash
cp backend/.env.example backend/.env
```

Depois disso, confira se os valores do `backend/.env` estao de acordo com o seu ambiente local. O arquivo `backend/.env.example` serve apenas como modelo seguro para o projeto.

## Banco de dados

Suba o MySQL padronizado com Docker:

```powershell
docker compose up -d mysql
```

Ou, se estiver usando Podman:

```powershell
podman compose up -d mysql
```

O container cria o banco usado pela aplicacao e executa automaticamente o script `db/task.sql` na primeira inicializacao do volume.

Para recriar o banco do zero durante testes locais:

```powershell
docker compose down -v
docker compose up -d mysql
```

## Executar a aplicacao

Entre na pasta do backend:

```powershell
cd backend
```

Baixe ou atualize as dependencias:

```powershell
go mod tidy
```

Inicie a aplicacao:

```powershell
go run .
```

Acesse no navegador:

```text
http://localhost:8080
```

## Endpoints da API

```text
GET    /api/v1/tasks
POST   /api/v1/tasks
PUT    /api/v1/tasks/{id}
DELETE /api/v1/tasks/{id}
```

## Testes manuais sugeridos

Com a aplicacao rodando, valide pelo navegador se e possivel:

- Criar uma tarefa
- Listar as tarefas cadastradas
- Marcar uma tarefa como concluida
- Excluir uma tarefa

Tambem teste entradas invalidas, como campo vazio, texto muito longo e conteudo HTML, para confirmar as validacoes do backend e a exibicao segura no frontend.

