Listar Contêineres em execução no Docker
```
docker ps

```
Saída:
```
CONTAINER ID   IMAGE             COMMAND                  CREATED         STATUS         PORTS                    NAMES
8ff446c6af12   suicidiosos-app   "/wait-for-it.sh db:…"   2 minutes ago   Up 2 minutes   0.0.0.0:8080->8080/tcp   suicidiosos-app-1
3ad8a6d1271d   postgres:16.4     "docker-entrypoint.s…"   6 minutes ago   Up 6 minutes   0.0.0.0:5432->5432/tcp   suicidiosos-db-1
```



Abrir banco de dados no docker
````
docker exec -it suicidiosos-db-1 psql -U suicidiosos_ions -d suicidiosos_db
```

Para terminar alguma compisções e inicar novamente

````
docker-compose down
docker-compose up --build
```