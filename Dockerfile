# Imagen base de Go
FROM golang:1.20.2 as builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del código fuente de tu aplicación al directorio de trabajo
COPY . .

# Compila tu aplicación Go
RUN go build -o main .
ENV APP_NAME=${APP_NAME}
ENV APP_PORT=${APP_PORT}
ENV APP_HOST=${APP_HOST}

RUN go build -o . main.go
# Expon el puerto 3001 en el contenedor
EXPOSE 3001

# Define el punto de entrada de la aplicación
CMD ["./main"]

# Crea un volumen para la carpeta "static"
VOLUME /app/static
