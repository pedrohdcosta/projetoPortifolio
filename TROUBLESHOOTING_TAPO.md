# Guia de Solução de Problemas - Integração Tapo

## Erro 403: "handshake 2 failed with status code: 403"

Este é o erro mais comum ao conectar dispositivos Tapo e geralmente está relacionado à **autenticação**.

---

## 🔍 Causas Comuns

### 1. **Formato Incorreto das Credenciais** ⚠️ MAIS COMUM
O campo `username` deve ser o **EMAIL da sua conta Tapo Cloud**, não um nome de usuário.

**❌ Errado:**
```json
{
  "tapo": {
    "ip": "192.168.1.100",
    "username": "meu_usuario",
    "password": "minhasenha"
  }
}
```

**✅ Correto:**
```json
{
  "tapo": {
    "ip": "192.168.1.100",
    "username": "seu-email@exemplo.com",
    "password": "senha_da_conta_tapo"
  }
}
```

### 2. **Credenciais Incorretas**
- Use a **mesma senha** que você usa no aplicativo Tapo (iOS/Android)
- Não é a senha do Wi-Fi do dispositivo
- Não é uma senha específica do dispositivo

### 3. **Versão do Firmware**
- Firmware mais recente (2023+) pode ter mudanças no protocolo de autenticação
- Verifique no app Tapo se há atualizações disponíveis

### 4. **Problemas de Rede**
- Dispositivo deve estar acessível pela rede do backend
- Backend e dispositivo devem estar na mesma rede (ou rotas configuradas)
- Firewall pode estar bloqueando portas 80/443

---

## 🛠️ Soluções

### Solução 1: Verificar Credenciais (PRIMEIRO PASSO)

1. **Abra o aplicativo Tapo no celular**
2. **Confirme que consegue controlar o dispositivo pelo app**
3. **Use o EMAIL da conta** (não username)
4. **No frontend, ao adicionar o dispositivo, configure:**

```json
{
  "name": "Tomada Sala",
  "room": "Sala",
  "type": "smart_plug",
  "metadata": {
    "tapo": {
      "ip": "192.168.1.100",
      "username": "seu-email@gmail.com",    ← EMAIL aqui!
      "password": "sua-senha-tapo"
    }
  }
}
```

### Solução 2: Verificar Conectividade de Rede

**No servidor (onde o backend roda):**

```bash
# Testar se o IP do dispositivo é alcançável
ping 192.168.1.100

# Testar porta HTTP
curl http://192.168.1.100

# Verificar se há resposta (mesmo que seja erro HTTP, confirma conectividade)
```

### Solução 3: Atualizar Firmware do Dispositivo

1. Abra o app Tapo
2. Vá em **Configurações** do dispositivo
3. Procure por **Atualização de Firmware**
4. Instale se houver atualizações

### Solução 4: Atualizar Biblioteca (Para Desenvolvedores)

Se nenhuma solução acima funcionar, pode ser incompatibilidade de versão:

```bash
cd Portifolio_back
go get -u github.com/tess1o/tapo-go@latest
go mod tidy
```

### Solução 5: Biblioteca Alternativa

Se o problema persistir, considere usar biblioteca alternativa:

```bash
# Trocar biblioteca
go get github.com/insomniacslk/tapo
```

---

## 📋 Checklist de Diagnóstico

Use este checklist para diagnosticar o problema:

- [ ] **Credenciais:** Estou usando EMAIL (não username) no campo "username"?
- [ ] **Senha:** É a mesma senha que uso no app Tapo?
- [ ] **App Funciona:** Consigo controlar o dispositivo pelo aplicativo móvel?
- [ ] **Rede:** Backend e dispositivo estão na mesma rede?
- [ ] **IP Correto:** O IP do dispositivo está correto? (Verificar no app ou router)
- [ ] **Ping:** Consigo fazer `ping` no IP do dispositivo do servidor?
- [ ] **Firmware:** Dispositivo tem firmware atualizado?
- [ ] **Firewall:** Há algum firewall bloqueando?

---

## 🔬 Testando a Conexão

### Teste Manual via API

Use este comando para testar a conexão manualmente:

```bash
# Substitua os valores
curl -X POST http://localhost:8080/api/devices/:id/toggle \
  -H "Authorization: Bearer SEU_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### Verificar Logs do Backend

```bash
# Se rodando com Docker
docker-compose logs -f backend

# Se rodando localmente
# Os logs devem mostrar o erro detalhado com sugestões
```

---

## 💡 Mensagem de Erro Melhorada

A partir da versão atual, o backend retorna erro mais descritivo:

```json
{
  "error": "failed to control physical device",
  "detail": "authentication failed (403): verify credentials and try these fixes:
    1. Use your Tapo Cloud EMAIL (not username) in the 'username' field
    2. Use your Tapo Cloud PASSWORD (the one you use in the mobile app)
    3. Ensure device firmware is up to date
    4. Check if device IP '192.168.1.100' is reachable from the backend
    Original error: ..."
}
```

---

## 📞 Outros Erros Comuns

### Erro: "timeout" ou "connection refused"
**Causa:** Dispositivo não alcançável pela rede  
**Solução:** Verificar IP, conectividade de rede, firewall

### Erro: "device not found"
**Causa:** IP incorreto ou dispositivo offline  
**Solução:** Confirmar IP no app Tapo, verificar se dispositivo está ligado

### Erro: "invalid credentials"
**Causa:** Email ou senha incorretos  
**Solução:** Confirmar credenciais no app Tapo

---

## 🎯 Exemplo Completo de Configuração

### 1. Descobrir IP do Dispositivo

**Opção A: Pelo App Tapo**
1. Abra app Tapo
2. Toque no dispositivo
3. Vá em **Configurações** (ícone de engrenagem)
4. Procure **Informações do Dispositivo**
5. Anote o endereço IP

**Opção B: Pelo Roteador**
1. Acesse interface web do roteador (geralmente 192.168.1.1)
2. Procure **Dispositivos Conectados** ou **DHCP Clients**
3. Encontre dispositivo com nome "Tapo_P110_XXXX"

### 2. Adicionar Dispositivo no Frontend

```json
POST /api/devices
Authorization: Bearer <seu_token_jwt>
Content-Type: application/json

{
  "name": "Tomada da TV",
  "room": "Sala de Estar",
  "type": "smart_plug",
  "metadata": {
    "tapo": {
      "ip": "192.168.1.105",
      "username": "pedro@email.com",
      "password": "minhaSenhaTapo123"
    }
  }
}
```

### 3. Testar Controle

```bash
# Ligar dispositivo
POST /api/devices/1/toggle

# Ler consumo
GET /api/devices/1/read
```

---

## 📚 Referências

- [Tapo-Go Library (GitHub)](https://github.com/tess1o/tapo-go)
- [TP-Link Tapo Support](https://www.tp-link.com/support/)
- [RFC - Seção 3.3: Stack Tecnológica](../README.MD#33-stack-tecnológica)

---

## ❓ Ainda com Problemas?

Se após seguir todas as soluções o problema persistir:

1. **Confirme versão do firmware** no app Tapo
2. **Tente com outro dispositivo Tapo** (se disponível)
3. **Verifique modelo do dispositivo** - alguns modelos podem ter restrições
4. **Considere criar uma issue** no repositório com:
   - Modelo do dispositivo (ex: Tapo P110, P100)
   - Versão do firmware
   - Logs completos do erro
   - Output de `ping` e `curl` para o IP

---

**Última atualização:** 05/12/2025  
**Versão do Backend:** Go 1.23.2 + tapo-go v0.1.1
