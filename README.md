# 🛍️ E-commerce Backend com Go e Fiber

Este repositório contém o código-fonte de um **backend para e-commerce**, desenvolvido com o objetivo de **aprender e praticar a linguagem Go** utilizando o framework **Fiber**. O projeto inclui funcionalidades essenciais para a operação de um e-commerce, seguindo boas práticas de organização de código, princípios de design (como SOLID), e arquitetura limpa.

---

## 🚀 Funcionalidades Implementadas

- **Autenticação e Autorização**
    - Registro e login de usuários.
    - Tokens JWT para autenticação segura.
- **Gestão de Produtos**
    - Cadastro, edição, listagem e exclusão de produtos.
- **Pedidos**
    - Criação de pedidos a partir do carrinho.
    - Histórico de pedidos por usuário.
- **Categorias**
    - Criação de categorias a partir de pedidos.
    - Histórico de criação de categorias por usuário.
- **Usuários**
    - Criação de usuários.
    - Usuarios podem ou nao ser vendedores também.
- **Administração**
    - CRUD de produtos e categorias.
    - Gerenciamento de usuários.

---

## 🛠️ Tecnologias Utilizadas

- **Linguagem**: Go
- **Framework Web**: Fiber
- **Banco de Dados**: PostgreSQL
- **ORM**: GORM
- **Autenticação**: JWT
- **Gerenciamento de Configuração**: Viper
- **Validação de Dados**: Validator
- **Logs**: Logrus
- **Testes**: Go Testing (com exemplos simples para aprendizado)

---

## 📂 Estrutura do Projeto

```plaintext
.
├── config/             # Configuração inicial do aplicativo
├── infra/              # Arquivos de configuração
├── internal/           # Lógica do domínio
│   ├── handlers/       # Handlers para rotas
│   ├── services/       # Lógica de negócios
│   ├── repositories/   # Acesso ao banco de dados
│   ├── dto/            # Arquivos de formatação de dados
│   ├── helper/         # Funções auxiliares
│   ├── api/            # configuração de rotas
│   ├── domain/         # Objetos 
├── pkg/                # Pacotes reutilizáveis
│   ├── notification/   # Lógica de disparo de notificação     
└── main.go             # Ponto de entrada da aplicação
```

---
## 🌱 Objetivo do Projeto
#### Este projeto foi desenvolvido para:

1. **Aprender Go**: Explorar a sintaxe e os conceitos fundamentais da linguagem.
2. **Praticar o uso de Fiber**: Entender como criar aplicações web rápidas e eficientes.
3. **Implementar Boas Práticas**: Aplicar princípios de design como SOLID e organização modular.
4. **Simular um cenário real**: Trabalhar em um backend com funcionalidades comuns no mundo real.

---

## 🧰 Como Executar o Projeto
### Pré-requisitos
- Go (versão 1.19 ou superior)
- PostgreSQL
- Go Fiber
- Ferramenta para gerenciamento de variáveis de ambiente (ex.: dotenv)

---
## 📝 Próximos Passos

- Adicionar testes unitários e de integração.
- Implementar upload de imagens para produtos.
- Adicionar um sistema de pagamento simulado.
- Melhorar a documentação da API usando Swagger.

---

## 🤝 Contribuições
Se você tiver sugestões ou melhorias, sinta-se à vontade para abrir uma issue ou enviar um pull request! 😊
