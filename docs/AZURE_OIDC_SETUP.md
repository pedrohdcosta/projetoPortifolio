# Configuração do Azure OIDC para GitHub Actions

Este documento descreve como resolver o erro "No subscriptions found" no GitHub Actions ao fazer deploy para o Azure usando OIDC (OpenID Connect).

## Erro Comum

```
Error: No subscriptions found for ***.
Error: Login failed with Error: The process '/usr/bin/az' failed with exit code 1.
```

## Causa Raiz

O Azure App Registration não está corretamente configurado com as Federated Credentials para o repositório GitHub.

## Solução Passo a Passo

### 1. Verificar o App Registration no Azure Portal

1. Acesse o [Azure Portal](https://portal.azure.com)
2. Navegue para **Microsoft Entra ID** → **App registrations**
3. Encontre o App Registration usado pelo workflow (verifique o Client ID nos secrets do GitHub)

### 2. Configurar Federated Credentials

1. No App Registration, vá para **Certificates & secrets**
2. Clique na aba **Federated credentials**
3. Clique em **Add credential**
4. Configure os seguintes valores:

   | Campo | Valor |
   |-------|-------|
   | Federated credential scenario | GitHub Actions deploying Azure resources |
   | Organization | `<seu-usuario-github>` (ex: `pedrohdcosta`) |
   | Repository | `<seu-repositorio>` (ex: `projetoPortifolio`) |
   | Entity type | `Branch` |
   | Branch name | `master` |
   | Name | Um nome descritivo (ex: `github-actions-deploy`) |

5. Clique em **Add**

### 3. Verificar a Role Assignment

O Service Principal precisa ter permissões na subscription do Azure:

1. Navegue para **Subscriptions** no Azure Portal
2. Selecione a subscription usada
3. Vá para **Access control (IAM)**
4. Clique em **Add** → **Add role assignment**
5. Selecione a role necessária (ex: `Contributor` para Web Apps)
6. Em **Members**, selecione **User, group, or service principal**
7. Procure pelo nome do seu App Registration
8. Clique em **Review + assign**

### 4. Verificar os Secrets no GitHub

No repositório GitHub, verifique se os secrets estão corretos:

- `AZUREAPPSERVICE_CLIENTID_...` → Client ID do App Registration
- `AZUREAPPSERVICE_TENANTID_...` → Tenant ID do Azure AD
- `AZUREAPPSERVICE_SUBSCRIPTIONID_...` → Subscription ID do Azure

Para encontrar esses valores:
1. **Client ID**: App Registration → Overview → Application (client) ID
2. **Tenant ID**: App Registration → Overview → Directory (tenant) ID
3. **Subscription ID**: Subscriptions → sua subscription → Overview → Subscription ID

### 5. Importante: Subject Claim

O subject claim no Azure deve corresponder exatamente ao formato esperado pelo GitHub:

- Para **branch**: `repo:<organization>/<repository>:ref:refs/heads/<branch>`
  - Exemplo: `repo:pedrohdcosta/projetoPortifolio:ref:refs/heads/master`
- Para **pull request**: `repo:<organization>/<repository>:pull_request`
- Para **environment**: `repo:<organization>/<repository>:environment:<environment-name>`

## Alternativa: Usar Service Principal com Secret

Se preferir não usar OIDC, você pode usar um Service Principal com secret:

1. Crie um secret no App Registration
2. Adicione as seguintes secrets no GitHub:
   - `AZURE_CLIENT_ID`
   - `AZURE_CLIENT_SECRET`
   - `AZURE_TENANT_ID`
   - `AZURE_SUBSCRIPTION_ID`

3. Modifique o workflow:

```yaml
- name: Login to Azure
  uses: azure/login@v2
  with:
    creds: '{"clientId":"${{ secrets.AZURE_CLIENT_ID }}","clientSecret":"${{ secrets.AZURE_CLIENT_SECRET }}","subscriptionId":"${{ secrets.AZURE_SUBSCRIPTION_ID }}","tenantId":"${{ secrets.AZURE_TENANT_ID }}"}'
```

## Referências

- [Azure Login GitHub Action](https://github.com/Azure/login)
- [Configure Azure credentials for GitHub Actions](https://learn.microsoft.com/en-us/azure/developer/github/connect-from-azure)
- [Use GitHub Actions to connect to Azure](https://learn.microsoft.com/en-us/azure/active-directory/workload-identities/workload-identity-federation-create-trust)
