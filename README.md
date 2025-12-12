# PC Ideal - Backend

> API inteligente para recomendação de builds de PC personalizadas

API REST em Go que utiliza IA generativa e estratégias de orçamento dinâmicas para recomendar configurações de PC otimizadas com base no perfil de uso e orçamento do usuário.

---

## Funcionalidades

- **Consulta de peças** com atualização automática de preços via scraper externo
- **Recomendações de builds** com 3 perfis (Econômica, Balanceada, Performance)
- **Estratégias de orçamento inteligentes** adaptadas ao tipo de uso (Gaming, Work, Office)
- **Análise por IA** usando Google Gemini para cada build gerada
- **Validação de compatibilidade** entre componentes

---

## Arquitetura

Projeto estruturado seguindo **Clean Architecture**:

- **Domain**: Entidades e regras de negócio (Part, Build, BudgetStrategy)
- **Use Cases**: Lógica de negócio e algoritmos de seleção de componentes
- **Infrastructure**: HTTP controllers, repositórios MongoDB, integrações externas
- **Dependency Injection**: Configurada no `main.go`

---

## Como Funciona

O sistema gera 3 builds (Econômica, Balanceada, Performance) seguindo estas etapas:

1. **Seleção de estratégia de orçamento** baseada no tipo de uso e valor
2. **Busca de componentes compatíveis** seguindo ordem de dependências (CPU → Motherboard → GPU → PSU → RAM → SSD)
3. **Validação de compatibilidade** (socket, tipo de memória, potência)
4. **Análise por IA** gerando recomendações personalizadas

---

## Endpoints

```http
GET /api/parts/              # Lista todas as peças
GET /api/parts/:id           # Detalhes de uma peça
POST /api/builds/recommendations  # Gera recomendações de builds
```

**Exemplo de requisição:**
```json
{
  "budget": 5000.00,
  "usage_type": "GAMING",
  "cpu_preference": "AMD",
  "gpu_preference": "NVIDIA"
}
```

---

## Stack

- **Go 1.24** com Gin Framework
- **MongoDB** para persistência
- **Google Gemini AI** para análises
- **Scraper API** para atualização de preços

---

## Executando o Projeto

```bash
# Clone e configure
git clone https://github.com/Luzin7/pcideal-be.git
cd pcideal-be

# Configure .env com as credenciais necessárias
PORT=8080
DATABASE_URL=mongodb://localhost:27017
PCIDEAL_DB_NAME=pcideal
SCRAPER_API_URL=...
SCRAPER_API_KEY=...
GOOGLE_AI_API_KEY=...

# Inicie MongoDB
docker-compose up -d

# Execute
go run cmd/main.go
```

## Testes

```bash
go test ./...                    # Todos os testes
go test -cover ./...             # Com cobertura
```

---

**PC Ideal** - Builds inteligentes para todos os orçamentos