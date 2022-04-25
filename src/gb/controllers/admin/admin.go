package admin

//TODO: Implement organization for http methods
//		So that when publishing events its easier to know which to use

import (
	"errors"
	"strconv"
	"strings"

	"insanitygaming.net/bans/src/gb"
	admcontrol "insanitygaming.net/bans/src/gb/controllers/groups/admin"
	servercontrol "insanitygaming.net/bans/src/gb/controllers/groups/server"
	webcontrol "insanitygaming.net/bans/src/gb/controllers/groups/web"
	"insanitygaming.net/bans/src/gb/models/admin"
	adm "insanitygaming.net/bans/src/gb/models/groups/admin"
	"insanitygaming.net/bans/src/gb/models/groups/server"
	"insanitygaming.net/bans/src/gb/models/groups/web"
)

func Find(app *gb.GB, id uint) (*admin.Admin, error) {
	var admin admin.Admin
	row, err := app.Database().QueryRow(app.Context(), "SELECT admin_id, name, email, created_at, adm_groups, web_groups, svr_groups, immunity FROM gb_admin WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("Admin not found")
	}
	var admgroups, webgroups, svrgroups string

	row.Scan(&admin.Id, &admin.Username, &admin.Email, &admin.CreatedAt, &admgroups, &webgroups, &svrgroups, &admin.Immunity)

	admin.AdminGroup = parseAdminGropsFromList(app, parseIntsFromString(admgroups))
	admin.WebGroups = parseWebGroupsFromList(app, parseIntsFromString(webgroups))
	admin.ServerGroups = parseServerGroupsFromList(app, parseIntsFromString(svrgroups))

	return &admin, nil
}

func FindByName(app *gb.GB, username string) (*admin.Admin, error) {
	var admin admin.Admin
	row, err := app.Database().QueryRow(app.Context(), "SELECT admin_id, name, email, auths, created_at, adm_groups, web_groups, svr_groups, immunity FROM gb_admin WHERE username = ?", username)
	if err != nil {
		return nil, errors.New("Admin not found")
	}
	var admgroups, webgroups, svrgroups string

	row.Scan(&admin.Id, &admin.Username, &admin.Email, &admin.CreatedAt, &admgroups, &webgroups, &svrgroups, &admin.Immunity)

	admin.AdminGroup = parseAdminGropsFromList(app, parseIntsFromString(admgroups))
	admin.WebGroups = parseWebGroupsFromList(app, parseIntsFromString(webgroups))
	admin.ServerGroups = parseServerGroupsFromList(app, parseIntsFromString(svrgroups))
	return &admin, nil
}

func FindByServerId(app *gb.GB, id uint) ([]*admin.Admin, error) {
	var admins []*admin.Admin
	rows, err := app.Database().Query(app.Context(), "SELECT admin_id, name, email, auths, created_at, adm_groups, web_groups, svr_groups, immunity FROM gb_admin WHERE FIND_IN_SET(?,servers)", id)
	if err != nil {
		return nil, errors.New("Admin not found")
	}
	for rows.Next() {
		var admin admin.Admin

		var admgroups, webgroups, svrgroups string

		rows.Scan(&admin.Id, &admin.Username, &admin.Email, &admin.CreatedAt, &admgroups, &webgroups, &svrgroups, &admin.Immunity)

		admin.AdminGroup = parseAdminGropsFromList(app, parseIntsFromString(admgroups))
		admin.WebGroups = parseWebGroupsFromList(app, parseIntsFromString(webgroups))
		admin.ServerGroups = parseServerGroupsFromList(app, parseIntsFromString(svrgroups))
		admins = append(admins, &admin)
	}
	return admins, nil
}

func FindByServerGroup(app *gb.GB, group uint) ([]*admin.Admin, error) {
	var admins []*admin.Admin
	rows, err := app.Database().Query(app.Context(), "SELECT admin_id, name, email, auths, created_at, adm_groups, web_groups, svr_groups, immunity FROM gb_admin WHERE FIND_IN_SET(?, server_groups) = ?", group)
	if err != nil {
		return nil, errors.New("Admin not found")
	}
	for rows.Next() {
		var admin admin.Admin
		var admgroups, webgroups, svrgroups string

		rows.Scan(&admin.Id, &admin.Username, &admin.Email, &admin.CreatedAt, &admgroups, &webgroups, &svrgroups, &admin.Immunity)

		admin.AdminGroup = parseAdminGropsFromList(app, parseIntsFromString(admgroups))
		admin.WebGroups = parseWebGroupsFromList(app, parseIntsFromString(webgroups))
		admin.ServerGroups = parseServerGroupsFromList(app, parseIntsFromString(svrgroups))
		admins = append(admins, &admin)
	}
	return admins, nil
}

func FindByApp(app *gb.GB, service string, id string) (*admin.Admin, error) {
	var admin *admin.Admin
	if !app.EventBus().HasCallback("get:admin:app:" + service) {
		return nil, errors.New("app not found")
	}
	app.EventBus().Publish("get:admin:app:"+service, id, &admin)
	app.EventBus().WaitAsync()
	return admin, nil
}

func parseIntsFromString(groups string) []uint {
	var ids []uint
	for _, group := range strings.Split(groups, ",") {
		id, err := strconv.ParseUint(group, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	return ids
}

func parseAdminGropsFromList(app *gb.GB, groups []uint) []adm.Group {
	var adminGroups []adm.Group
	for _, group := range groups {
		adminGroup, err := admcontrol.Find(app, uint(group))
		if err != nil {
			continue
		}
		adminGroups = append(adminGroups, *adminGroup)
	}
	return adminGroups
}

func parseWebGroupsFromList(app *gb.GB, groups []uint) []web.Group {
	var webGroups []web.Group
	for _, group := range groups {
		webGroup, err := webcontrol.Find(app, uint(group))
		if err != nil {
			continue
		}
		webGroups = append(webGroups, *webGroup)
	}
	return webGroups
}

func parseServerGroupsFromList(app *gb.GB, groups []uint) []server.Group {
	var serverGroups []server.Group
	for _, group := range groups {
		serverGroup, err := servercontrol.Find(app, uint(group))
		if err != nil {
			continue
		}
		serverGroups = append(serverGroups, *serverGroup)
	}
	return serverGroups
}
