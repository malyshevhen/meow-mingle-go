package main

import "github.com/go-sql-driver/mysql"

func main() {
	cfg := mysql.Config{
		User:                 Envs.DBUser,
		Passwd:               Envs.DBPasswd,
		Addr:                 Envs.DBAddress,
		DBName:               Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqlStorage := NewMySQLStorage(cfg)

	userService := NewUserService(sqlStorage.db)
	postService := NewPostService(sqlStorage.db)
	commentService := NewCommentService(sqlStorage.db)

	sCtx := InitSecurityContext(userService)

	server := NewApiServer(":8080", userService, postService, commentService, sCtx)
	server.Serve()
}
