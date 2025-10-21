package user

import (
	//"database/sql"
	pb "robust-backend/microservices/auth/gen/user"
	"sync"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	mu    sync.Mutex
	users map[string]*pb.User
	//db *sql.DB
}
