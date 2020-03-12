package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/ysugimoto/grpc-graphql-gateway/example/starwars/spec/starwars"
	"google.golang.org/grpc"
)

var characters map[int64]*starwars.Character

// nolint: funlen
func init() {
	characters = map[int64]*starwars.Character{
		1000: {
			Id:         1000,
			Name:       "Luke Skywalker",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		1001: {
			Id:         1001,
			Name:       "Darth Vader",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		1002: {
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
		1003: {
			Id:         1003,
			Name:       "Leia Organa",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Alderaa",
			Type:       starwars.Type_HUMAN,
		},
		1004: {
			Id:        1004,
			Name:      "Wilhuff Tarkin",
			AppearsIn: []starwars.Episode{1},
			Type:      starwars.Type_HUMAN,
		},
		2000: {
			Id:              2000,
			Name:            "C-3PO",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Protocol",
			Type:            starwars.Type_DROID,
		},
		2001: {
			Id:              2001,
			Name:            "R2-D2",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Astromech",
			Type:            starwars.Type_DROID,
		},
	}
	characters[1000].Friends = append(characters[1000].Friends,
		&starwars.Character{
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:         1003,
			Name:       "Leia Organa",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Alderaa",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:              2000,
			Name:            "C-3PO",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Protocol",
			Type:            starwars.Type_DROID,
		},
		&starwars.Character{
			Id:              2001,
			Name:            "R2-D2",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Astromech",
			Type:            starwars.Type_DROID,
		},
	)
	characters[1001].Friends = append(characters[1001].Friends,
		&starwars.Character{
			Id:        1004,
			Name:      "Wilhuff Tarkin",
			AppearsIn: []starwars.Episode{1},
			Type:      starwars.Type_HUMAN,
		},
	)
	characters[1002].Friends = append(characters[1002].Friends,
		&starwars.Character{
			Id:         1000,
			Name:       "Luke Skywalker",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:         1003,
			Name:       "Leia Organa",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Alderaa",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:              2001,
			Name:            "R2-D2",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Astromech",
			Type:            starwars.Type_DROID,
		},
	)
	characters[1003].Friends = append(characters[1003].Friends,
		&starwars.Character{
			Id:         1000,
			Name:       "Luke Skywalker",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:              2000,
			Name:            "C-3PO",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Protocol",
			Type:            starwars.Type_DROID,
		},
		&starwars.Character{
			Id:              2001,
			Name:            "R2-D2",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Astromech",
			Type:            starwars.Type_DROID,
		},
	)
	characters[1004].Friends = append(characters[1004].Friends,
		&starwars.Character{
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
	)
	characters[2000].Friends = append(characters[2000].Friends,
		&starwars.Character{
			Id:         1000,
			Name:       "Luke Skywalker",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:         1003,
			Name:       "Leia Organa",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Alderaa",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:              2001,
			Name:            "R2-D2",
			AppearsIn:       []starwars.Episode{1, 2, 3},
			PrimaryFunction: "Astromech",
			Type:            starwars.Type_DROID,
		},
	)
	characters[2001].Friends = append(characters[2001].Friends,
		&starwars.Character{
			Id:         1000,
			Name:       "Luke Skywalker",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Tatooine",
			Type:       starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:        1002,
			Name:      "Han Solo",
			AppearsIn: []starwars.Episode{1, 2, 3},
			Type:      starwars.Type_HUMAN,
		},
		&starwars.Character{
			Id:         1003,
			Name:       "Leia Organa",
			AppearsIn:  []starwars.Episode{1, 2, 3},
			HomePlanet: "Alderaa",
			Type:       starwars.Type_HUMAN,
		},
	)
}

type Server struct{}

func (s *Server) GetHero(
	ctx context.Context, req *starwars.GetHeroRequest) (*starwars.Character, error) {

	if req.GetEpisode() == starwars.Episode_EMPIRE {
		return characters[1000], nil
	}
	return characters[2001], nil
}

func (s *Server) GetHuman(
	ctx context.Context, req *starwars.GetHumanRequest) (*starwars.Character, error) {

	h, ok := characters[req.GetId()]
	if !ok {
		return nil, errors.New("character not found")
	}
	if h.GetType() != starwars.Type_HUMAN {
		return nil, errors.New("character is not a human")
	}
	return h, nil
}

func (s *Server) GetDroid(
	ctx context.Context, req *starwars.GetDroidRequest) (*starwars.Character, error) {

	d, ok := characters[req.GetId()]
	if !ok {
		return nil, errors.New("character not found")
	}
	if d.GetType() != starwars.Type_DROID {
		return nil, errors.New("character is not a droid")
	}
	return d, nil
}

func (s *Server) ListHumans(
	ctx context.Context, req *starwars.ListEmptyRequest) (*starwars.ListHumansResponse, error) {

	cs := make([]*starwars.Character, 0)
	for _, c := range characters {
		if c.GetType() != starwars.Type_HUMAN {
			continue
		}
		cs = append(cs, c)
	}
	log.Println(cs)
	return &starwars.ListHumansResponse{
		Humans: cs,
	}, nil
}

func (s *Server) ListDroids(
	ctx context.Context, req *starwars.ListEmptyRequest) (*starwars.ListDroidsResponse, error) {

	var cs []*starwars.Character
	for _, c := range characters {
		if c.GetType() != starwars.Type_DROID {
			continue
		}
		cs = append(cs, c)
	}
	return &starwars.ListDroidsResponse{
		Droids: cs,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	srv := &Server{}
	g := grpc.NewServer()
	starwars.RegisterStartwarsServiceServer(g, srv)
	if err := g.Serve(listener); err != nil {
		log.Println(err)
	}
}
