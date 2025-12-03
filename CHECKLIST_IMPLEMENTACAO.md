# Checklist de Implementação - RFC Plataforma Energia IoT

**Data:** 03/12/2025  
**Status Atual:** 50% Completo  
**Meta:** 100% MVP Conforme RFC v1.4

---

## 🎯 Resumo de Status

- ✅ **Completo:** 6/10 Requisitos Funcionais
- ⚠️ **Parcial:** 3/10 Requisitos Funcionais  
- ❌ **Pendente:** 1/10 Requisitos Funcionais
- 🔴 **Crítico:** 3 bloqueadores identificados

---

## 📋 Requisitos Funcionais (RF)

### ✅ RF01: Cadastro/Login de Usuários
- [x] Endpoint de registro
- [x] Endpoint de login
- [x] JWT implementado
- [x] Middleware de autenticação
- [x] Hash de senhas (bcrypt)
- [ ] ⚠️ Sistema de roles (Admin/User) - **ADICIONAR**

---

### ✅ RF02: CRUD de Dispositivos IoT
- [x] Criar dispositivo
- [x] Listar dispositivos
- [x] Obter dispositivo
- [x] Atualizar dispositivo
- [x] Deletar dispositivo
- [x] Validação de propriedade

---

### ❌ RF03: Receber Telemetria MQTT - **BLOQUEADOR CRÍTICO**
**Prioridade:** 🔴 CRÍTICA

- [ ] Adicionar EMQX ao docker-compose.yml
- [ ] Criar `internal/mqtt/client.go`
- [ ] Implementar subscriber para tópicos
- [ ] Processar mensagens MQTT → telemetry table
- [ ] Configurar TLS mútuo
- [ ] Adicionar retry logic
- [ ] Testes de integração MQTT

**Referências:**
- RFC Seção 3.3: EMQX 5 (MQTT v3.1.1)
- Diagrama de Sequência: Fluxo MQTT

**Dependências:**
```yaml
# docker-compose.yml
mqtt:
  image: emqx/emqx:5.0
  ports:
    - "1883:1883"
    - "8883:8883"  # TLS
    - "18083:18083"  # Dashboard
```

```bash
# Go dependencies
go get github.com/eclipse/paho.mqtt.golang
```

---

### ⚠️ RF04: Consumir Dados via API REST do Dispositivo
**Prioridade:** 🟡 ALTA

- [x] Integração Tapo (básica)
- [x] Endpoint `/devices/:id/read`
- [ ] Worker background para polling automático
- [ ] Integração Shelly API
- [ ] Integração Tuya API
- [ ] Integração Tasmota API
- [ ] Abstração genérica `DeviceAdapter` interface

**Estrutura Sugerida:**
```go
// internal/integrations/adapter.go
type DeviceAdapter interface {
    ReadTelemetry(ctx context.Context) (*Telemetry, error)
    SetPower(ctx context.Context, on bool) error
}

// Implementações:
// - internal/integrations/tapo/adapter.go
// - internal/integrations/shelly/adapter.go
// - internal/integrations/tuya/adapter.go
```

---

### ⚠️ RF05: Exibir Consumo em Tempo Real (<60s)
**Prioridade:** 🟡 ALTA

- [x] Endpoint GET `/api/telemetry`
- [x] Gráfico no dashboard
- [ ] **WebSocket para real-time push**
- [ ] Broadcast de novas leituras
- [ ] Reconexão automática no frontend
- [ ] Rate limiting

**Implementação Sugerida:**
```go
// cmd/api/main.go
r.GET("/ws/telemetry", telemetryHandler.WebSocket)

// Frontend
const ws = new WebSocket('ws://localhost:8080/ws/telemetry')
ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  updateChart(data)
}
```

---

### ✅ RF06: Visualizar Histórico por Período
- [x] Query com filtros de data
- [x] Paginação
- [x] Ordenação
- [x] Visualização gráfica

---

### ⚠️ RF07: Definir Limites e Enviar Alertas
**Prioridade:** 🟡 ALTA

- [x] UI de configuração (frontend)
- [x] localStorage thresholds
- [ ] **Tabela `alert_consumo` no banco**
- [ ] **CRUD de alertas (backend)**
- [ ] **Worker para processar alertas**
- [ ] **Integração SMTP (e-mail)**
- [ ] **Push notifications (opcional)**
- [ ] Histórico de alertas disparados

**Schema SQL Proposto:**
```sql
CREATE TABLE alert_config (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES device(id),
    user_id BIGINT NOT NULL REFERENCES app_user(id),
    threshold_type TEXT NOT NULL CHECK(threshold_type IN ('warning', 'danger')),
    threshold_value DOUBLE PRECISION NOT NULL,
    notification_email BOOLEAN DEFAULT true,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE alert_history (
    id BIGSERIAL PRIMARY KEY,
    alert_config_id BIGINT NOT NULL REFERENCES alert_config(id),
    telemetry_id BIGINT NOT NULL REFERENCES telemetry(id),
    triggered_at TIMESTAMPTZ DEFAULT NOW(),
    acknowledged BOOLEAN DEFAULT false
);
```

**Backend:**
```go
// internal/alerts/worker.go
func (w *Worker) ProcessAlerts(ctx context.Context) {
    // Query recentes telemetry readings
    // Check against alert_config
    // Send email if threshold exceeded
    // Insert into alert_history
}
```

---

### ❌ RF08: Projeções Heurísticas (Médias Móveis)
**Prioridade:** 🟢 MÉDIA

- [ ] Criar package `internal/analytics/`
- [ ] Calcular média móvel 7 dias
- [ ] Calcular média móvel 30 dias
- [ ] Projeção de consumo mensal
- [ ] Estimativa de custo (R$/kWh)
- [ ] Endpoint `GET /api/devices/:id/analytics`
- [ ] Dashboard com projeções

**Exemplo de Cálculo:**
```go
// internal/analytics/projections.go
func (s *Service) CalculateMovingAverage(deviceID int64, days int) (float64, error) {
    start := time.Now().AddDate(0, 0, -days)
    readings, err := s.repo.GetTelemetryRange(deviceID, start, time.Now())
    // Calcular média
}

func (s *Service) ProjectMonthlyConsumption(deviceID int64) (*Projection, error) {
    // Usar média dos últimos 7 dias * 30
    // Multiplicar por tarifa
}
```

---

### ❌ RF09: Relatórios Mensais em PDF
**Prioridade:** 🔵 BAIXA (mas no MVP)

- [ ] Integrar biblioteca `go-pdf/fpdf` ou `jung-kurt/gofpdf`
- [ ] Template de relatório mensal
- [ ] Gráficos como imagens (Chart.js → canvas → png)
- [ ] Endpoint `GET /api/reports/:device_id/pdf?month=2024-11`
- [ ] Job agendado mensal (cron)
- [ ] Envio por e-mail (opcional)

**Estrutura:**
```go
// internal/reports/pdf_generator.go
func (g *Generator) GenerateMonthlyReport(deviceID int64, month time.Month) ([]byte, error)

// internal/reports/scheduler.go
func (s *Scheduler) ScheduleMonthlyReports() // Cron job
```

---

### ⚠️ RF10: Controle Remoto de Tomadas
**Prioridade:** 🟡 ALTA

- [x] Toggle via REST (Tapo)
- [x] Endpoint `POST /devices/:id/toggle`
- [ ] **Toggle via MQTT**
- [ ] Abstração para múltiplos fabricantes
- [ ] Logs de comandos enviados
- [ ] Confirmação de execução

**MQTT Command Flow:**
```
Backend PUBLISH → energy/{device_id}/command/power {"on": true}
Device SUBSCRIBE → executa
Device PUBLISH → energy/{device_id}/status {"power": true}
Backend SUBSCRIBE → atualiza DB
```

---

## 🔒 Requisitos Não Funcionais (RNF)

### ❌ RNF01: Latência API p95 < 2000ms
- [ ] Implementar testes de carga (k6, locust)
- [ ] APM (Prometheus + Grafana)
- [ ] Otimizar queries lentas
- [ ] Cache Redis para dados frequentes

---

### ⚠️ RNF02: Suporte a 50 Dispositivos Simultâneos
- [ ] Teste de carga (50 devices, 1 leitura/segundo)
- [ ] Validar índices do banco
- [ ] Connection pooling adequado

---

### ⚠️ RNF03: Disponibilidade ≥ 99.5%
- [ ] Health checks (`/health`, `/readiness`)
- [ ] Monitoramento uptime (UptimeRobot, Pingdom)
- [ ] Alertas de downtime

---

### ⚠️ RNF04: Criptografia TLS
- [x] HTTPS em produção (Azure)
- [x] JWT
- [ ] **MQTT TLS mútuo**

---

### ❌ RNF05: Cobertura de Testes ≥ 80% - **BLOQUEADOR CRÍTICO**
**Prioridade:** 🔴 CRÍTICA

**Status Atual:** 0% (nenhum teste implementado)

#### Backend (Go)
- [ ] `internal/auth/service_test.go`
- [ ] `internal/devices/repo_test.go`
- [ ] `internal/devices/handler_test.go`
- [ ] `internal/telemetry/repo_test.go`
- [ ] `internal/mqtt/client_test.go`
- [ ] Mocks para DB (`testify/mock`)
- [ ] CI: falhar se cobertura < 80%

```bash
# Executar testes
go test -v -cover ./...

# Gerar relatório HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

#### Frontend (Vue)
- [ ] Configurar Vitest ou Jest
- [ ] Testes unitários para stores (Pinia)
- [ ] Testes de componentes (Testing Library)
- [ ] E2E testes (Playwright, Cypress)

---

### ⚠️ RNF06: Acessibilidade WCAG 2.1 AA
- [ ] Auditoria com axe-core
- [ ] Corrigir issues encontrados
- [ ] Adicionar labels semânticos
- [ ] Keyboard navigation
- [ ] Contraste de cores adequado

---

## 🏗️ Arquitetura e Infraestrutura

### Docker Compose Completo
- [x] Backend service
- [ ] **PostgreSQL/TimescaleDB**
- [ ] **EMQX (MQTT broker)**
- [ ] **Redis**
- [ ] Frontend (opcional, para dev local)

```yaml
services:
  postgres:
    image: timescale/timescaledb:latest-pg15
    environment:
      POSTGRES_DB: energydb
      POSTGRES_USER: energyuser
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  mqtt:
    image: emqx/emqx:5.0
    ports:
      - "1883:1883"
      - "8883:8883"
      - "18083:18083"
    volumes:
      - ./emqx/data:/opt/emqx/data
      - ./emqx/log:/opt/emqx/log

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build: ./Portifolio_back
    depends_on:
      - postgres
      - mqtt
      - redis
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - MQTT_BROKER_URL=tcp://mqtt:1883
      - REDIS_URL=redis://redis:6379/0
    ports:
      - "8080:8080"

volumes:
  pgdata:
  redis_data:
```

---

### Migração para TimescaleDB
- [ ] Atualizar docker-compose (imagem timescaledb)
- [ ] Script SQL de migração:

```sql
-- Migration: Enable TimescaleDB
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Convert telemetry to hypertable
SELECT create_hypertable('telemetry', 'timestamp', if_not_exists => TRUE);

-- Retenção automática (6 meses)
SELECT add_retention_policy('telemetry', INTERVAL '6 months');

-- Compressão automática (dados > 7 dias)
ALTER TABLE telemetry SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'device_id'
);

SELECT add_compression_policy('telemetry', INTERVAL '7 days');

-- Continuous aggregates para dashboards
CREATE MATERIALIZED VIEW telemetry_hourly
WITH (timescaledb.continuous) AS
SELECT
  device_id,
  time_bucket('1 hour', timestamp) AS hour,
  AVG(power) AS avg_power,
  MAX(power) AS max_power,
  MIN(power) AS min_power
FROM telemetry
GROUP BY device_id, hour;

SELECT add_continuous_aggregate_policy('telemetry_hourly',
  start_offset => INTERVAL '1 month',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '1 hour');
```

- [ ] Atualizar queries no backend para usar views

---

### Variáveis de Ambiente Completas
Adicionar ao `.env.example`:

```env
# Database
DATABASE_URL=postgresql://energyuser:energypass@localhost:5432/energydb?sslmode=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# MQTT
MQTT_BROKER_URL=tcp://localhost:1883
MQTT_USERNAME=admin
MQTT_PASSWORD=public
MQTT_CLIENT_ID=energy-backend
MQTT_TLS_ENABLED=false

# Redis
REDIS_URL=redis://localhost:6379/0
REDIS_PASSWORD=

# Email Alerts (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=alerts@energy.com
SMTP_PASSWORD=
SMTP_FROM=noreply@energy.com

# Features
ENABLE_TIMESCALEDB=true
ENABLE_MQTT=true
ENABLE_EMAIL_ALERTS=true

# App
APP_PORT=8080
APP_ENV=development
```

---

## 📊 Diagramas - Atualizações Necessárias

### Diagrama de Classes
- [ ] Adicionar classe `AlertConfig`
- [ ] Adicionar classe `AlertHistory`
- [ ] Expandir `ApiDeviceAdapter` para interface genérica

---

### Diagrama de Sequência
- [ ] Adicionar fluxo completo MQTT (pub/sub)
- [ ] Adicionar fluxo de alertas
- [ ] Adicionar fluxo WebSocket

---

### Diagrama C4
- [ ] Adicionar Redis ao contêiner
- [ ] Adicionar EMQX ao contêiner
- [ ] Atualizar para TimescaleDB

---

## 🧪 Testes e Qualidade

### Testes Backend
```bash
# Suite completa
- [ ] internal/auth/service_test.go (signup, login, JWT)
- [ ] internal/auth/middleware_test.go
- [ ] internal/devices/repo_test.go (CRUD)
- [ ] internal/devices/handler_test.go (HTTP endpoints)
- [ ] internal/telemetry/repo_test.go
- [ ] internal/mqtt/client_test.go (mock broker)
- [ ] internal/alerts/worker_test.go
- [ ] internal/integrations/tapo/tapo_test.go

# Testes de integração
- [ ] tests/integration/auth_flow_test.go
- [ ] tests/integration/device_crud_test.go
- [ ] tests/integration/mqtt_ingestion_test.go
```

### Testes Frontend
```bash
# Vitest config
- [ ] vitest.config.ts
- [ ] src/stores/auth.test.ts
- [ ] src/api/axios.test.ts
- [ ] src/components/ConsumptionChart.test.ts

# E2E
- [ ] e2e/login.spec.ts
- [ ] e2e/dashboard.spec.ts
- [ ] e2e/devices.spec.ts
```

### CI/CD
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: 1.23
      - run: go test -v -coverprofile=coverage.out ./...
      - run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage $COVERAGE% is below 80%"
            exit 1
          fi

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 20
      - run: npm ci
        working-directory: Portifolio_front/energy-controller
      - run: npm run test -- --coverage
        working-directory: Portifolio_front/energy-controller
```

---

## 🚀 Plano de Implementação

### Sprint 1: MQTT + Testes (2 semanas)
**Objetivo:** Eliminar bloqueadores críticos

#### Semana 1: MQTT
- [ ] Dia 1-2: Setup EMQX no docker-compose
- [ ] Dia 3-4: Implementar `internal/mqtt/client.go`
- [ ] Dia 5: Integração com telemetry service
- [ ] Dia 6-7: Testes de integração MQTT

#### Semana 2: Testes
- [ ] Dia 1-3: Testes auth + devices (meta 80%)
- [ ] Dia 4-5: Testes telemetry + mqtt
- [ ] Dia 6: CI/CD com coverage check
- [ ] Dia 7: Buffer / fixes

---

### Sprint 2: Alertas + Real-time (2 semanas)
**Objetivo:** Completar RF07 e RF05

#### Semana 1: Sistema de Alertas
- [ ] Dia 1: Schema `alert_config` + `alert_history`
- [ ] Dia 2-3: CRUD alertas (backend)
- [ ] Dia 4-5: Worker de processamento
- [ ] Dia 6: Integração SMTP
- [ ] Dia 7: UI backend alertas

#### Semana 2: Real-time
- [ ] Dia 1-2: WebSocket backend
- [ ] Dia 3-4: WebSocket frontend
- [ ] Dia 5: Testes + otimizações
- [ ] Dia 6-7: Integração frontend ↔ backend

---

### Sprint 3: Analytics + TimescaleDB (2 semanas)
**Objetivo:** RF08 e otimizações

#### Semana 1: TimescaleDB
- [ ] Dia 1-2: Migração PostgreSQL → TimescaleDB
- [ ] Dia 3-4: Hypertables + continuous aggregates
- [ ] Dia 5: Otimizar queries
- [ ] Dia 6-7: Testes de performance

#### Semana 2: Analytics
- [ ] Dia 1-3: Implementar cálculos (médias, projeções)
- [ ] Dia 4-5: Endpoint `/analytics`
- [ ] Dia 6: Dashboard com projeções
- [ ] Dia 7: Testes

---

### Sprint 4: Multi-Device + PDF (2 semanas)
**Objetivo:** Abstrações e relatórios

#### Semana 1: Abstração Devices
- [ ] Interface `DeviceAdapter`
- [ ] Implementações Shelly, Tuya, Tasmota
- [ ] Worker de polling REST

#### Semana 2: PDF Reports
- [ ] Biblioteca PDF
- [ ] Template relatório
- [ ] Endpoint + scheduler

---

### Sprint 5: Hardening (2 semanas)
**Objetivo:** Produção-ready

- [ ] Testes de carga (50 devices)
- [ ] APM (Prometheus)
- [ ] Health checks
- [ ] Acessibilidade WCAG
- [ ] Documentação completa
- [ ] Security audit

---

## 📝 Documentação Necessária

- [ ] **API.md** - Documentação OpenAPI/Swagger completa
- [ ] **DEPLOYMENT.md** - Guia de deploy (Docker, Azure, K8s)
- [ ] **ARCHITECTURE.md** - Decisões arquiteturais
- [ ] **CONTRIBUTING.md** - Guia de contribuição
- [ ] **TESTING.md** - Guia de testes
- [ ] Atualizar README.md com setup completo

---

## ✅ Critérios de Aceitação do MVP

### Funcional
- [x] Usuário pode se cadastrar e fazer login
- [x] Usuário pode adicionar/editar/remover dispositivos
- [ ] **Sistema recebe telemetria via MQTT automaticamente**
- [x] Sistema recebe telemetria via REST (Tapo)
- [ ] **Dashboard atualiza em tempo real (WebSocket)**
- [ ] **Sistema envia alertas por e-mail quando limites são excedidos**
- [x] Usuário pode controlar dispositivos remotamente (REST)
- [ ] **Usuário pode controlar dispositivos via MQTT**

### Não Funcional
- [ ] **Cobertura de testes ≥ 80%**
- [ ] API responde em < 2s (p95)
- [ ] Sistema suporta 50 dispositivos simultâneos
- [ ] Dados persistem em TimescaleDB
- [ ] Deploy automatizado via CI/CD

### Documentação
- [x] RFC v1.4 (completo)
- [x] Diagramas atualizados
- [ ] API documentada (Swagger)
- [ ] README com setup
- [ ] Guia de deploy

---

## 🎯 Meta Final

**MVP Completo:** 10/10 Requisitos Funcionais + 6/6 RNF atendidos  
**Prazo Estimado:** 5 sprints (10 semanas)  
**Prioridade Imediata:** MQTT + Testes (Sprint 1)

---

**Última atualização:** 03/12/2025  
**Responsável:** Análise Automatizada
