package main

import (
	desc "chat_server/pkg/chat_v1"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	grpcPort          = 50052
	chatTable         = "chat"
	chatsToUsersTable = "chats_to_users"
	messageTable      = "message"
	dbDsn             = "host=localhost port=54322 dbname=chat-service user=dev_course password=1801 sslmode=disable"
)

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	if len(req.Usernames) == 0 {
		return nil, errors.New("Usernames shouldn't be null")
	}

	pool, err := pgx.Connect(ctx, dbDsn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to database %s", err)
	}

	insertBuilder := squirrel.Insert(chatTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns("name").
		Values("chat name").
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build query %s", err)
	}

	var chatId int64
	err = pool.QueryRow(ctx, query, args...).Scan(&chatId)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute %s", err)
	}

	for _, v := range req.Usernames {
		insertBuider := squirrel.Insert(chatsToUsersTable).
			PlaceholderFormat(squirrel.Dollar).
			Columns("chatId", "username").
			Values(chatId, v)

		query, args, err := insertBuider.ToSql()
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to build query %s", err)
		}

		_, err = pool.Exec(ctx, query, args...)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to execute %s", err)
		}
	}

	return &desc.CreateChatResponse{Id: chatId}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	pool, err := pgx.Connect(ctx, dbDsn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to database %s", err)
	}

	deleteBuilder := squirrel.Delete(chatsToUsersTable).
		Where(squirrel.Eq{"chatId": req.Id})

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build query %s", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute %s", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	pool, err := pgx.Connect(ctx, dbDsn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to database %s", err)
	}

	insertBuilder := squirrel.Insert(messageTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns("\"from\"", "text").
		Values(req.From, req.Text)

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build query %s", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute %s", err)
	}
	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	fmt.Println("Server has been started")
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
