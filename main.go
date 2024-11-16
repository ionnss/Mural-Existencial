// main.go
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Estrutura para representar um post
type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// Estrutura do servidor que contém o banco de dados e os templates
type Server struct {
	db   *sql.DB
	tmpl *template.Template
}

func main() {
	// Configura as variáveis de ambiente para conectar ao banco de dados no Docker
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Cria a string de conexão com o banco de dados
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	// Cria a tabela de posts se não existir
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Carrega e configura os templates HTML
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Inicializa a estrutura do servidor com o banco de dados e os templates
	server := &Server{
		db:   db,
		tmpl: tmpl,
	}

	// Define as rotas e as funções que irão lidar com elas
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/sobre", server.handleSobre)
	http.HandleFunc("/post", server.handleCreatePost)
	http.HandleFunc("/posts", server.handleLoadPosts)

	// Servir arquivos estáticos (como CSS e JavaScript) da pasta "static"
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Inicia o servidor HTTP na porta 8080
	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Função para carregar a página inicial
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "index.html", nil)
}

// Função para carregar a página "Sobre"
func (s *Server) handleSobre(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "sobre.html", nil)
}

// Função para criar um novo post
func (s *Server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método da requisição é POST
	if r.Method != http.MethodPost {
		log.Println("Erro: Método não permitido")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtém os valores do título e do conteúdo do post
	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" || content == "" {
		log.Println("Erro: Campo vazio")
		http.Error(w, "Conteúdo não pode estar vazio", http.StatusBadRequest)
		return
	}

	// Inserir post no banco de dados
	_, err := s.db.Exec(
		"INSERT INTO posts (title, content) VALUES ($1, $2)",
		title, content,
	)
	if err != nil {
		log.Printf("Erro ao inserir no banco de dados: %v", err)
		http.Error(w, "Erro ao criar post", http.StatusInternalServerError)
		return
	}

	log.Println("Post criado com sucesso")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Função para carregar os posts e exibi-los na página
func (s *Server) handleLoadPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`
        SELECT id, title, content, created_at 
        FROM posts 
        ORDER BY created_at DESC 
        LIMIT 50
    `)
	if err != nil {
		http.Error(w, "Erro ao carregar posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Cria uma lista para armazenar os posts
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			log.Printf("Erro ao scanear post: %v", err)
			continue // Ignora o erro e continua com o próximo post
		}
		posts = append(posts, post)
	}

	log.Printf("Total de posts carregados: %d", len(posts))

	// Renderiza o template `posts.html` passando a lista de posts
	if err := s.tmpl.ExecuteTemplate(w, "posts.html", posts); err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
		return
	}
}
