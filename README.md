# financial-cli

Este projeto é uma **releitura do [Financial-Manager](https://github.com/Ch3did/Financial-Manager)**, originalmente feito em **Python**, agora **reescrito em Go**, com arquitetura baseada em DDD


## Funcionalidades

- Importa transações a partir de arquivos OFX
- Categorização interativa das transações
- Registro dos imports para controle de qual arquivo importado veio cada uma das trasações
- Resumo mensal de gastos com barras de progresso
- Criação e listagem de categorias personalizadas

## Requisitos

- Go 1.23 ou superior instalado

## Instalação

Use o script de instalação para configurar a aplicação no seu sistema Linux:

```bash
sudo ./install.sh
```    

Isso vai:

- Copiar a aplicação para `/opt/financial-cli`
- Criar o comando `cx` no `/usr/local/bin` para executar a CLI
- Executar automaticamente o comando `cx init` para criar o banco de dados

## Categorias pré-definidas

  Se quiser adicionar categorias padrão automaticamente, utilize o script add_categories.sh disponível na pasta documents:

  ```bash
   ./documents/add_categories.sh
   ```

## Uso

Após a instalação, você pode executar os seguintes comandos no terminal:

- `cx import -p caminho/do/arquivo.ofx`  
  Importa transações a partir do arquivo OFX especificado.

- `cx add-category`  
  Adiciona uma nova categoria de gastos manualmente, com entrada interativa.

- `cx home`  
  Exibe um resumo dos gastos do mês atual, com barras de progresso.

- `cx categories`  
  Lista todas as categorias registradas no sistema.

- `cx init`  
  Inicializa ou reseta o banco de dados da aplicação.

> **Notas:**  
> O banco de dados será criado automaticamente na pasta do usuário, dentro do diretório do projeto (`/home/usuario/fcli.db`), garantindo que seus dados fiquem localmente armazenados.


