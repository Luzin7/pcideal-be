# PC Ideal - Backend

Este é o backend de uma plataforma voltada para entusiastas, consultores e compradores de hardware de computador, focada em facilitar a montagem de builds personalizadas e a busca por peças no mercado nacional.

---

## 💡 Visão Geral

O PC Ideal oferece uma API REST em Go, integrando-se a um scraper externo para manter sempre atualizadas as informações e preços de peças de hardware, tornando fácil pesquisar, comparar e planejar configurações de PCs.

O sistema foi pensado para ser modular, extensível e de fácil integração com frontends e serviços de terceiros.

---

## 🚀 Funcionalidades Principais

- **Consulta de Peças:** Busque peças individuais por ID ou modelo, ou faça um levantamento de todas disponíveis no banco.
- **Cadastro e Atualização Dinâmica:** Adicione novas peças ao sistema, com atualização automática de preços e especificações por meio de integração com scraper.
- **Modelagem Completa:** Cada peça é cadastrada com informações detalhadas (tipo, marca, modelo, specs completas, preço, loja de origem, etc).
- **Atualização Automática:** O sistema detecta quando os dados estão desatualizados e dispara atualizações em background.
- **Preparado para Builds:** Estrutura pronta para implementar montagem de builds personalizadas, gerenciamento de orçamento e recomendações.
- **Arquitetura Modular:** Separação clara entre controllers, services, repositórios, contratos e modelos.

---

## 🛠️ Estrutura do Projeto

```
cmd/
  main.go            # Inicialização da aplicação
internal/
  core/
    models/          # Modelos de dados (Part, Build, etc)
  contracts/         # Interfaces para scraper e repositórios
  http/
    controllers/     # Controllers das rotas REST
    services/        # Lógica de negócio
    routes/          # Definição de rotas
infra/
  database/          # Conexão e configuração do MongoDB
  repositories/      # Implementação dos repositórios (Mongo)
```

---

## 📦 Exemplos de Recursos

### Peça (`Part`)
```json
{
  "id": "abcdef123456",
  "type": "GPU",
  "brand": "NVIDIA",
  "model": "RTX 4060 Ti",
  "specs": {
    "memory_type": "GDDR6",
    "capacity": 8,
    "interface": "PCIe 4.0",
    "power_supply": 550
  },
  "price_cents": 320000,
  "url": "https://loja.com/produto/rtx4060ti",
  "store": "Loja do PC",
  "updated_at": "2025-05-23T00:00:00Z"
}
```

<!-- ### Build (`Build`)
```json
{
  "id": "build789",
  "user_id": "user123",
  "goal": "gaming",
  "budget": 5000,
  "parts": ["cpuId", "gpuId", "ramId", "ssdId"],
  "total_price": 4899.90,
  "created_at": "2025-05-23T00:00:00Z"
}
``` -->

---

## 🌐 Endpoints REST (Principais)

- `GET /api/parts/` – Lista todas as peças cadastradas.
- `GET /api/parts/:id` – Consulta uma peça por ID.
- (Próximos endpoints: busca por modelo, autenticação, builds, etc.)

---

## 🔌 Tecnologias e Integrações

- **Go:** Backend robusto, tipado e performático.
- **MongoDB:** Armazenamento NoSQL, flexível para diferentes tipos de peças.
- **Gin Gonic:** Framework web para rotas e middlewares.
- **Scraper HTTP:** Integração com serviço externo para atualizar specs e preços.

---

## 📚 Modelos de Dados

- **Part:** Representa uma peça de hardware, incluindo specs detalhadas.
- **Build:** Representa uma configuração de PC personalizada.
- **Specs:** Subdocumento com informações técnicas específicas para cada tipo de peça.