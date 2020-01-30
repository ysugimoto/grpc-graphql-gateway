package main

import (
	"context"
	"errors"
	"log"
	"net"

	"encoding/json"
	"io/ioutil"

	"google.golang.org/grpc"

	"github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/author"
	"github.com/ysugimoto/grpc-graphql-gateway/examples/basic/app/book"
)

// Thanks https://gist.github.com/nanotaboada/6396437
var localData = struct {
	Books   []*book.Book     `json:"books"`
	Authors []*author.Author `json:"authors"`
}{}

func init() {
	buf, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(buf, &localData); err != nil {
		log.Fatalln(err)
	}
}

type App struct{}

func (a *App) ListBooks(ctx context.Context, req *book.ListBooksRequest) (*book.ListBooksResponse, error) {
	return &book.ListBooksResponse{
		Books: localData.Books,
	}, nil
}

func (a *App) GetBook(ctx context.Context, req *book.GetBookRequest) (*book.Book, error) {
	id := req.GetId()
	for _, b := range localData.Books {
		if b.GetId() == id {
			return b, nil
		}
	}
	return nil, errors.New("book not found")
}

func (a *App) ListAuthors(ctx context.Context, req *author.ListAuthorsRequest) (*author.ListAuthorsResponse, error) {
	return &author.ListAuthorsResponse{
		Authors: localData.Authors,
	}, nil
}
func (a *App) GetAuthor(ctx context.Context, req *author.GetAuthorRequest) (*author.Author, error) {
	name := req.GetName()
	for _, a := range localData.Authors {
		if a.GetName() == name {
			return a, nil
		}
	}
	return nil, errors.New("author not found")
}

func main() {
	listener, err := net.Listen("tcp4", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	app := &App{}
	book.RegisterBookServiceServer(server, app)
	author.RegisterAuthorServiceServer(server, app)

	log.Println("Server starts on 0.0.0.0:8080...")
	server.Serve(listener)
}
