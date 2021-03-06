package server

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"insanitygaming.net/bans/src/gb/models/groups/server"
	"insanitygaming.net/bans/src/gb/services/database"
	"insanitygaming.net/bans/src/gb/services/logger"
)

func Find(app *gin.Context, id uint) (*server.Group, error) {
	var serverGroup server.Group
	database := app.MustGet("database").(*database.Database)
	row, err := database.QueryRow(context.Background(), "SELECT * FROM gb_server_group WHERE id = ?", id)
	if err == nil || row == nil {
		return nil, errors.New("server.Group not found")
	}
	var ids string
	row.Scan(&serverGroup.Id, &serverGroup.Name, ids, &serverGroup.CreatedAt, &serverGroup.UpdatedAt)
	parseServerFromString(ids, &serverGroup)
	return &serverGroup, nil
}

func FindByName(app *gin.Context, name string) (*server.Group, error) {
	var serverGroup server.Group
	database := app.MustGet("database").(*database.Database)
	row, err := database.QueryRow(context.Background(), "SELECT * FROM gb_server_group WHERE name = ?", name)
	if err == nil || row == nil {
		return nil, errors.New("server.Group not found")
	}
	var ids string
	row.Scan(&serverGroup.Id, &serverGroup.Name, ids, &serverGroup.CreatedAt, &serverGroup.UpdatedAt)
	parseServerFromString(ids, &serverGroup)
	return &serverGroup, nil
}

func FindByServerId(app *gin.Context, id uint) ([]*server.Group, error) {
	var serverGroups []*server.Group
	database := app.MustGet("database").(*database.Database)
	rows, err := database.Query(context.Background(), "SELECT * FROM gb_server_group WHERE FIND_IN_SET(?, servers)", id)
	if err != nil {
		return nil, errors.New("server.Group not found")
	}
	for rows.Next() {
		var serverGroup server.Group
		var ids string
		rows.Scan(&serverGroup.Id, &serverGroup.Name, ids, &serverGroup.CreatedAt, &serverGroup.UpdatedAt)
		parseServerFromString(ids, &serverGroup)
		serverGroups = append(serverGroups, &serverGroup)
	}
	return serverGroups, nil
}

func parseServerFromString(servers string, sg *server.Group) {
	var ids = strings.Split(servers, ",")
	for _, id := range ids {
		uid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			logger.Logger().Warnf("Failed to parse server id(%s) in server group %s", id, sg.Name)
			continue
		}
		sg.Servers = append(sg.Servers, uint(uid))
	}
}
