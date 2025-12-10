# Tabela Comparativa: RFC vs Implementação Atual

**Data:** 03/12/2025  
**RFC Versão:** 1.4  
**Branch:** copilot/analyze-repository-branch

---

## Legenda de Status

| Ícone | Status | Descrição |
|-------|--------|-----------|
| ✅ | Completo | Implementado conforme RFC |
| ⚠️ | Parcial | Implementado, mas incompleto ou com desvios |
| ❌ | Ausente | Não implementado |
| 🔴 | Crítico | Bloqueador para MVP |
| 🟡 | Alto | Importante para MVP |
| 🟢 | Médio | Desejável |
| 🔵 | Baixo | Pós-MVP |

---

## 1. Requisitos Funcionais

| ID | Requisito RFC | Prioridade RFC | Status | Prioridade Impl | Observações |
|----|---------------|----------------|--------|-----------------|-------------|
| RF01 | Cadastro/Login (Admin, Usuário) | Alta | ⚠️ | 🟢 | ✅ Auth JWT implementado<br>❌ Roles não diferenciados |
| RF02 | CRUD Dispositivos IoT | Alta | ✅ | - | Completo com validação de propriedade |
| RF03 | Receber telemetria MQTT | Alta | ❌ | 🔴 | **BLOQUEADOR:** Nenhum código MQTT<br>Broker não configurado |
| RF04 | Consumir API REST dispositivo | Média | ⚠️ | 🟡 | ✅ Tapo implementado<br>❌ Shelly, Tuya, Tasmota ausentes<br>❌ Polling automático ausente |
| RF05 | Exibir consumo real-time (<60s) | Alta | ⚠️ | 🟡 | ✅ Dashboard + gráficos<br>❌ WebSocket ausente (apenas polling manual) |
| RF06 | Visualizar histórico período | Alta | ✅ | - | Completo com paginação |
| RF07 | Definir limites e alertas | Média | ⚠️ | 🟡 | ✅ UI frontend<br>❌ Backend não implementado<br>❌ E-mail não implementado |
| RF08 | Projeções heurísticas | Média | ❌ | 🟢 | Nenhum código de analytics |
| RF09 | Relatórios PDF mensais | Baixa | ❌ | 🔵 | Nenhuma biblioteca PDF |
| RF10 | Controle remoto tomadas | Média | ⚠️ | 🟡 | ✅ REST (Tapo) implementado<br>❌ MQTT não implementado |

### Métricas RF
- **Completos:** 2/10 (20%)
- **Parciais:** 4/10 (40%)
- **Ausentes:** 4/10 (40%)
- **Bloqueadores:** 1 (RF03)

---

## 2. Requisitos Não Funcionais

| ID | Requisito RFC | Meta RFC | Status | Medido | Gap |
|----|---------------|----------|--------|--------|-----|
| RNF01 | Latência API p95 | < 2000ms | ⚠️ | Não medido | Sem APM |
| RNF02 | Dispositivos simultâneos | ≥ 50 | ⚠️ | Não testado | Sem testes de carga |
| RNF03 | Disponibilidade mensal | ≥ 99.5% | ⚠️ | Não medido | Sem monitoramento |
| RNF04 | Criptografia TLS | 100% | ⚠️ | Parcial | ✅ HTTPS<br>❌ MQTT TLS |
| RNF05 | Cobertura testes | ≥ 80% | ❌ | 0% | **BLOQUEADOR:** Nenhum teste |
| RNF06 | Acessibilidade WCAG | 2.1 AA | ⚠️ | Não auditado | Sem testes |

### Métricas RNF
- **Atendidos:** 0/6 (0%)
- **Parciais:** 5/6 (83%)
- **Não atendidos:** 1/6 (17%)
- **Bloqueadores:** 1 (RNF05)

---

## 3. Stack Tecnológica

| Camada | RFC Especifica | Implementado | Status | Notas |
|--------|----------------|--------------|--------|-------|
| **Dispositivo** | ESP32 + SCT-013 | Tapo Smart Plug | ⚠️ | Diferente, mas válido |
| **Mensageria** | EMQX 5 (MQTT v3.1.1) | ❌ Ausente | ❌ | **CRÍTICO** |
| **API Device** | REST + OAuth2 | Tapo REST | ⚠️ | OAuth2 não implementado |
| **Backend** | Go 1.22 (Gin) | Go 1.23.2 (Gin) | ✅ | Versão mais recente |
| **Base dados** | PostgreSQL 15 / TimescaleDB 2.15 | PostgreSQL (sem TimescaleDB) | ⚠️ | Sem otimizações para séries temporais |
| **Frontend** | Vue 3 + Vite (PWA) | Vue 3.5 + Vite | ✅ | PWA não configurado |
| **Cache** | (não especificado) | ❌ Sem Redis | ⚠️ | Mencionado em C4, não implementado |
| **Infra Dev** | Docker Compose | Docker Compose (parcial) | ⚠️ | Apenas backend, sem serviços |
| **Infra Prod** | Kubernetes + Helm | Azure Static Web Apps | ⚠️ | Diferente, mas funcional |

### Dependências Go (`go.mod`)

| Biblioteca | RFC Necessário | Implementado | Versão |
|------------|----------------|--------------|--------|
| `gin-gonic/gin` | ✅ | ✅ | v1.10.0 |
| `pgx/v5` (PostgreSQL) | ✅ | ✅ | v5.6.0 |
| `golang-jwt/jwt` | ✅ | ✅ | v5.2.1 |
| `paho.mqtt.golang` | ✅ | ❌ | - |
| `tess1o/tapo-go` | - | ✅ | v0.1.1 |
| `go-redis/redis` | - | ❌ | - |
| PDF library | - | ❌ | - |

### Dependências Frontend (`package.json`)

| Biblioteca | RFC Necessário | Implementado | Versão |
|------------|----------------|--------------|--------|
| `vue` | ✅ | ✅ | 3.5.18 |
| `vue-router` | ✅ | ✅ | 4.5.1 |
| `pinia` | ✅ | ✅ | 3.0.3 |
| `axios` | ✅ | ✅ | 1.12.2 |
| `chart.js` | ✅ | ✅ | 4.5.1 |
| PWA plugin | ✅ | ❌ | - |
| WebSocket client | - | ❌ | - |

---

## 4. Banco de Dados

### Tabelas Implementadas vs RFC

| Entidade RFC | Tabela Implementada | Status | Campos Faltantes |
|--------------|---------------------|--------|------------------|
| `Usuario` | `app_user` | ✅ | ❌ `role` (Admin/User) |
| `Dispositivo` | `device` | ✅ | ⚠️ `tipoConexao` (inferido via metadata) |
| `LeituraEnergia` | `telemetry` | ✅ | - |
| `AlertaConsumo` | ❌ **Ausente** | ❌ | **Tabela completa ausente** |

### Schema Detalhado

#### ✅ `app_user` (5/5 campos esperados)
```sql
id BIGSERIAL                 ✅
name TEXT                    ✅
email TEXT UNIQUE            ✅
password_hash TEXT           ✅
created_at TIMESTAMPTZ       ✅ (bonus)
-- FALTA: role TEXT (admin/user)
```

#### ✅ `device` (7/7 campos + extras)
```sql
id BIGSERIAL                 ✅
user_id BIGINT               ✅
name TEXT                    ✅
room TEXT                    ✅ (bonus)
type TEXT                    ✅
status TEXT                  ✅
metadata TEXT                ✅ (para tipoConexao + config)
created_at TIMESTAMPTZ       ✅ (bonus)
last_seen TIMESTAMPTZ        ✅ (bonus)
power_state BOOLEAN          ✅ (bonus, não em main.go mas em repo)
```

#### ✅ `telemetry` (5/5 campos esperados)
```sql
id BIGSERIAL                 ✅
device_id BIGINT             ✅
power DOUBLE PRECISION       ✅
voltage DOUBLE PRECISION     ✅
current DOUBLE PRECISION     ✅
timestamp TIMESTAMPTZ        ✅
-- Índices:
-- idx_telemetry_device_id  ✅
-- idx_telemetry_timestamp  ✅
```

#### ❌ `alert_consumo` (0/4 campos - tabela ausente)
```sql
-- ESPERADO:
id BIGSERIAL
device_id BIGINT
limite DOUBLE PRECISION
ativo BOOLEAN

-- RECOMENDADO (expandido):
id, device_id, user_id
threshold_type (warning/danger)
threshold_value
notification_email
active
created_at
```

### Otimizações TimescaleDB

| Recurso | RFC | Implementado | Impacto |
|---------|-----|--------------|---------|
| Extensão TimescaleDB | ✅ | ❌ | **Alto:** Performance com séries temporais |
| Hypertable | ✅ | ❌ | **Alto:** Particionamento automático |
| Retention policy | - | ❌ | **Médio:** Gestão de espaço |
| Compression | - | ❌ | **Médio:** Armazenamento otimizado |
| Continuous aggregates | - | ❌ | **Alto:** Dashboards rápidos |

---

## 5. Endpoints Backend

### Autenticação

| Endpoint | Método | RFC | Implementado | Autenticação |
|----------|--------|-----|--------------|--------------|
| `/api/auth/signup` | POST | ✅ | ✅ | Pública |
| `/api/auth/login` | POST | ✅ | ✅ | Pública |
| `/api/auth/me` | GET | - | ✅ | JWT |

### Dispositivos

| Endpoint | Método | RFC | Implementado | Autenticação |
|----------|--------|-----|--------------|--------------|
| `/api/devices` | GET | ✅ | ✅ | JWT |
| `/api/devices` | POST | ✅ | ✅ | JWT |
| `/api/devices/:id` | GET | ✅ | ✅ | JWT |
| `/api/devices/:id` | PUT | ✅ | ✅ | JWT |
| `/api/devices/:id` | DELETE | ✅ | ✅ | JWT |
| `/api/devices/:id/toggle` | POST | ✅ | ✅ | JWT |
| `/api/devices/:id/read` | GET | ✅ | ✅ | JWT |

### Telemetria

| Endpoint | Método | RFC | Implementado | Autenticação |
|----------|--------|-----|--------------|--------------|
| `/api/telemetry` | GET | ✅ | ✅ | JWT |
| `/api/telemetry` | POST | ✅ | ✅ | JWT |
| `/api/devices/:id/telemetry` | GET | ✅ | ✅ | JWT |
| `/api/telemetry/simulate` | POST | - | ✅ | JWT |

### Alertas (Faltantes)

| Endpoint | Método | RFC | Implementado | Prioridade |
|----------|--------|-----|--------------|------------|
| `/api/alerts` | GET | ✅ | ❌ | 🟡 Alta |
| `/api/alerts` | POST | ✅ | ❌ | 🟡 Alta |
| `/api/alerts/:id` | PUT | ✅ | ❌ | 🟡 Alta |
| `/api/alerts/:id` | DELETE | ✅ | ❌ | 🟡 Alta |
| `/api/alerts/history` | GET | - | ❌ | 🟢 Média |

### Analytics (Faltantes)

| Endpoint | Método | RFC | Implementado | Prioridade |
|----------|--------|-----|--------------|------------|
| `/api/devices/:id/analytics` | GET | ✅ | ❌ | 🟢 Média |
| `/api/devices/:id/projections` | GET | ✅ | ❌ | 🟢 Média |

### Relatórios (Faltantes)

| Endpoint | Método | RFC | Implementado | Prioridade |
|----------|--------|-----|--------------|------------|
| `/api/reports/:device_id/pdf` | GET | ✅ | ❌ | 🔵 Baixa |

### Real-time (Faltantes)

| Endpoint | Protocolo | RFC | Implementado | Prioridade |
|----------|-----------|-----|--------------|------------|
| `/ws/telemetry` | WebSocket | ✅ | ❌ | 🟡 Alta |

**Total Endpoints:**
- **Implementados:** 14/25 (56%)
- **Faltantes:** 11/25 (44%)

---

## 6. Frontend - Páginas e Componentes

### Páginas

| Página | RFC | Implementado | Funcionalidade | Status |
|--------|-----|--------------|----------------|--------|
| Login | ✅ | ✅ | Auth | Completo |
| Register | ✅ | ✅ | Auth | Completo |
| Dashboard | ✅ | ✅ | Consumo real-time | ⚠️ Sem WebSocket |
| Devices | ✅ | ✅ | CRUD dispositivos | Completo |
| Thresholds | ✅ | ✅ | Config limites | ⚠️ Apenas frontend |
| Profile | - | ✅ | Dados usuário | Bonus |

### Componentes

| Componente | RFC | Implementado | Uso |
|------------|-----|--------------|-----|
| `ConsumptionChart.vue` | ✅ | ✅ | Gráfico Chart.js |
| `TelemetryTable.vue` | ✅ | ✅ | Tabela histórico |

### Stores (Pinia)

| Store | Implementado | Uso |
|-------|--------------|-----|
| `auth.ts` | ✅ | JWT + user state |

### Services/API

| Service | Implementado | Endpoints |
|---------|--------------|-----------|
| `axios.ts` | ✅ | HTTP client com interceptors |
| `devices.ts` | ✅ | API devices + telemetry |

### Utilitários

| Util | Implementado | Uso |
|------|--------------|-----|
| `thresholds.ts` | ✅ | localStorage para limites (⚠️ deveria ser backend) |
| `simulationMode.ts` | ✅ | Mock data para testes |

---

## 7. Integração com Dispositivos

### Fabricantes Suportados

| Fabricante | RFC Menciona | Implementado | Arquivo | Funcionalidades |
|------------|--------------|--------------|---------|-----------------|
| **TP-Link Tapo** | ✅ | ✅ | `integrations/tapo/tapo.go` | ✅ SetPower<br>✅ ReadPower |
| **Shelly** | ✅ | ❌ | - | - |
| **Tuya** | ✅ | ❌ | - | - |
| **Tasmota** | ✅ | ❌ | - | - |
| **Genérico MQTT** | ✅ | ❌ | - | - |

### Funcionalidades Tapo

```go
// internal/integrations/tapo/tapo.go
type Connection struct {
    IP       string
    Username string
    Password string
}

func SetPower(ctx context.Context, conn Connection, on bool) error    ✅
func ReadPower(ctx context.Context, conn Connection) (float64, error)  ✅
```

**Status:** 1/5 fabricantes (20%)

---

## 8. CI/CD e Deploy

### GitHub Workflows

| Workflow | Arquivo | Propósito | Status |
|----------|---------|-----------|--------|
| Azure Deploy | `deploy.yml` | Deploy backend Azure | ✅ |
| Azure Static Web Apps | `azure-static-web-apps-*.yml` | Deploy frontend | ✅ |
| Tests | - | Rodar testes + coverage | ❌ |
| Build | - | Build + lint | ❌ |

### Configuração Missing

- ❌ Workflow de testes automáticos
- ❌ Coverage check (≥ 80%)
- ❌ Linting automático
- ❌ Security scan (Snyk, Dependabot)

---

## 9. Diagramas

### Conformidade com RFC

| Diagrama | Arquivo PlantUML | PNG | RFC Seção | Status |
|----------|------------------|-----|-----------|--------|
| Casos de Uso | `useCaseDiagram.puml` | ✅ | 3.1.3 | ✅ Atualizado |
| Classes | `classDiagram.puml` | ✅ | 6.1 | ⚠️ Falta `AlertaConsumo` |
| Sequência | `sequenceDiagram.puml` | ✅ | 6.1 | ⚠️ MQTT não implementado |
| C4 Contexto | `C4diagrams.puml` | ✅ | 3.2.1 | ⚠️ MQTT não implementado |
| C4 Contêiner | `C4diagrams.puml` | ✅ | 3.2.2 | ⚠️ Redis não implementado |

---

## 10. Cronograma vs Realidade

| Sprint | RFC Planejamento | Status Atual | Gap |
|--------|------------------|--------------|-----|
| **1** | Repo, CI/CD, Vue, MQTT POC, REST POC | ⚠️ 60% | ❌ MQTT POC ausente |
| **2** | CRUD users/devices, telemetry mock | ✅ 100% | - |
| **3** | TimescaleDB, dashboard realtime | ⚠️ 50% | ❌ TimescaleDB<br>❌ Realtime |
| **4** | Alertas, controle remoto | ⚠️ 40% | ❌ Alertas backend<br>❌ MQTT control |
| **5** | Projeções, PDF | ❌ 0% | **Não iniciado** |
| **6** | Hardening, testes, docs | ❌ 0% | **Não iniciado** |

**Status Geral do Cronograma:** ~40% completo

---

## 11. Segurança e Compliance

| Aspecto | RFC | Implementado | Status |
|---------|-----|--------------|--------|
| **Autenticação** | JWT + OAuth2 | JWT | ⚠️ OAuth2 ausente |
| **Senhas** | bcrypt | bcrypt | ✅ |
| **HTTPS** | TLS 1.3 | Azure HTTPS | ✅ |
| **MQTT TLS** | mTLS | N/A | ❌ MQTT não implementado |
| **RBAC** | Roles Admin/User | Não | ❌ |
| **Data encryption at rest** | AES-256 | Depende Azure | ⚠️ |
| **LGPD Compliance** | Sim | Parcial | ⚠️ Sem política de retenção |

---

## 12. Resumo Executivo

### Pontos Fortes ✅
1. **Arquitetura bem estruturada** - Separação clara de responsabilidades
2. **Auth completo** - JWT funcionando
3. **CRUD devices** - Implementação sólida
4. **Frontend moderno** - Vue 3 + Chart.js
5. **Deploy automatizado** - Azure CI/CD configurado
6. **Integração básica IoT** - Tapo funcionando

### Bloqueadores Críticos 🔴
1. **MQTT não implementado** - Funcionalidade core ausente
2. **Testes 0%** - RNF05 não atendido (meta: 80%)
3. **Alertas incompletos** - Apenas frontend

### Gaps Significativos 🟡
4. **Real-time ausente** - Sem WebSocket
5. **TimescaleDB não usado** - Apenas PostgreSQL
6. **Alertas sem backend** - Não persiste, não envia e-mail
7. **Multi-device** - Apenas Tapo suportado

### Melhorias Desejáveis 🟢
8. **Analytics ausente** - Sem projeções/médias
9. **Redis não configurado** - Cache ausente
10. **PDF não implementado** - Relatórios ausentes

---

## 13. Score Final

### Por Categoria

| Categoria | Peso | Score | Ponderado |
|-----------|------|-------|-----------|
| Requisitos Funcionais | 40% | 50% | 20% |
| Requisitos Não Funcionais | 30% | 20% | 6% |
| Stack Tecnológica | 15% | 60% | 9% |
| Documentação | 10% | 70% | 7% |
| Segurança | 5% | 50% | 2.5% |

**Score Total:** **44.5% / 100%**

### Distribuição de Esforço Restante

```
█████████░░░░░░░░░░░ 45% - Implementado
███░░░░░░░░░░░░░░░░░ 15% - MQTT + Testes (Sprint 1)
██░░░░░░░░░░░░░░░░░░ 10% - Alertas + Real-time (Sprint 2)
██░░░░░░░░░░░░░░░░░░ 10% - Analytics + TimescaleDB (Sprint 3)
███░░░░░░░░░░░░░░░░░ 15% - Multi-device + PDF (Sprint 4)
█░░░░░░░░░░░░░░░░░░░ 5% - Hardening (Sprint 5)
```

---

## 14. Recomendação Final

### Caminho Crítico para MVP

**Prioridade Absoluta (2 semanas):**
1. Implementar MQTT completo (RF03)
2. Criar suite de testes ≥80% (RNF05)

**Alta Prioridade (4 semanas):**
3. Sistema de alertas backend + e-mail (RF07)
4. WebSocket para real-time (RF05)
5. Migrar para TimescaleDB

**Pós-MVP (6 semanas):**
6. Analytics e projeções (RF08)
7. Multi-device support (RF04)
8. PDF reports (RF09)

**Estimativa Total:** 12 semanas (~3 meses) para MVP completo conforme RFC

---

**Documento gerado:** 03/12/2025 (DD/MM/YYYY - 3 de Dezembro de 2025)  
**Próxima revisão:** Após Sprint 1 (implementação MQTT + testes)
