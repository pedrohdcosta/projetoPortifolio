# Frontend Testing Guide

Este guia explica como executar e criar testes para o frontend do projeto Energy Controller.

## Configuração

O projeto utiliza **Vitest** como framework de testes, que é a solução recomendada para projetos Vite + Vue 3.

### Dependências Instaladas

- `vitest` - Framework de testes rápido e moderno
- `@vue/test-utils` - Utilitário oficial para testar componentes Vue
- `happy-dom` - Ambiente DOM leve para testes
- `@vitest/ui` - Interface visual para executar testes

## Scripts Disponíveis

```bash
# Executar testes em modo watch (re-executa quando arquivos mudam)
npm run test

# Executar testes uma vez (útil para CI/CD)
npm run test:run

# Executar testes com interface visual
npm run test:ui
```

## Estrutura de Testes

Os testes estão organizados em diretórios `__tests__` próximos aos arquivos que testam:

```
src/
├── api/
│   ├── __tests__/
│   │   └── devices.test.ts       # Testes das funções da API
│   └── devices.ts
├── pages/
│   ├── __tests__/
│   │   └── Devices.connect.test.ts  # Testes do componente Devices
│   └── Devices.vue
```

## Testes Implementados

### 1. Testes da API (`src/api/__tests__/devices.test.ts`)

Testa a função `testDeviceConnection()` que conecta ao dispositivo:

- ✅ Chama o endpoint correto (`/devices/:id/read`)
- ✅ Retorna o valor de potência da resposta
- ✅ Trata erros de rede
- ✅ Trata erro de dispositivo não configurado
- ✅ Trata erros de autenticação

**Exemplo de execução:**
```bash
npm run test:run -- src/api/__tests__/devices.test.ts
```

### 2. Testes do Componente (`src/pages/__tests__/Devices.connect.test.ts`)

Testa a funcionalidade do botão "Conectar":

#### Renderização baseada no modo:
- ✅ Mostra botão "Simular" quando modo simulação está ativo
- ✅ Mostra botão "Conectar" quando modo API está ativo
- ✅ Não mostra ambos os botões ao mesmo tempo

#### Funcionalidade de conexão:
- ✅ Chama `testDeviceConnection` ao clicar no botão
- ✅ Atualiza status do dispositivo para "online" após conexão bem-sucedida
- ✅ Mostra indicador de carregamento (⏳) enquanto conecta
- ✅ Exibe mensagem de erro quando conexão falha

**Exemplo de execução:**
```bash
npm run test:run -- src/pages/__tests__/Devices.connect.test.ts
```

## Como Escrever Novos Testes

### Teste de Função da API

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { suaFuncao } from '../seuArquivo'
import api from '../axios'

vi.mock('../axios', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

describe('Sua Funcionalidade', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('deve fazer algo específico', async () => {
    const mockResponse = { data: { resultado: 'esperado' } }
    vi.mocked(api.get).mockResolvedValue(mockResponse)

    const resultado = await suaFuncao(123)

    expect(api.get).toHaveBeenCalledWith('/endpoint/esperado')
    expect(resultado).toEqual({ resultado: 'esperado' })
  })
})
```

### Teste de Componente Vue

```typescript
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SeuComponente from '../SeuComponente.vue'

describe('SeuComponente', () => {
  it('deve renderizar corretamente', () => {
    const wrapper = mount(SeuComponente, {
      props: {
        propExemplo: 'valor'
      }
    })

    expect(wrapper.find('.classe-esperada').exists()).toBe(true)
  })

  it('deve reagir a clicks', async () => {
    const wrapper = mount(SeuComponente)
    
    await wrapper.find('button').trigger('click')
    
    expect(wrapper.emitted()).toHaveProperty('evento-esperado')
  })
})
```

## Cobertura de Código

Para gerar relatório de cobertura de código, você pode adicionar a flag `--coverage`:

```bash
npm run test:run -- --coverage
```

Isso gerará um relatório em HTML na pasta `coverage/` mostrando quais partes do código estão cobertas por testes.

## Boas Práticas

1. **Um teste por comportamento**: Cada teste deve verificar um comportamento específico
2. **Nomes descritivos**: Use nomes que expliquem o que está sendo testado
3. **Arrange-Act-Assert**: Organize seus testes em três partes:
   - Arrange: Configure o estado inicial
   - Act: Execute a ação que está testando
   - Assert: Verifique o resultado esperado
4. **Isolar dependências**: Use mocks para isolar o código sendo testado
5. **Testes rápidos**: Testes devem executar rapidamente para não atrasar o desenvolvimento

## Integração com CI/CD

Para executar testes em um pipeline de CI/CD, adicione:

```yaml
- name: Run Frontend Tests
  run: |
    cd Portifolio_front/energy-controller
    npm ci
    npm run test:run
```

## Recursos Adicionais

- [Documentação do Vitest](https://vitest.dev/)
- [Vue Test Utils](https://test-utils.vuejs.org/)
- [Guia de Testes Vue.js](https://vuejs.org/guide/scaling-up/testing.html)
