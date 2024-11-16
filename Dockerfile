# Usar imagem oficial do Golang
FROM golang:1.23.2

# Definir diretório de trabalho
WORKDIR /muralexistencial

# Copiar arquivos de dependências
COPY go.mod ./
COPY go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código-fonte
COPY . .

# Compilar aplicativo Go com otimizações para produção
RUN CGO_ENABLED=0 GOOS=linux go build -o app_muralexistencial -ldflags="-s -w"

# Expor porta 8080
EXPOSE 8080

# Comando para rodar o executável
CMD ["./app_muralexistencial"]

# Baixar o script wait-for-it
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh

# Tornar o script executável
RUN chmod +x /wait-for-it.sh
