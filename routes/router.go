package routes

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type RouterServiceInterface interface {
	RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListController)
	RegisterHandlerToGroup(group *gin.RouterGroup, path string, method string, handler gin.HandlerFunc)
}

type RouterService struct {
	app *gin.Engine
}

func NewRouterService(app *gin.Engine) (svc *RouterService) {
	return &RouterService{
		app: app,
	}
}

func (svc *RouterService) RegisterControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.BasicController) {
	group.GET(basePath, ctr.Get)
	group.PUT(basePath, ctr.Put)
	group.POST(basePath, ctr.Post)
	group.DELETE(basePath, ctr.Delete)
}

func (svc *RouterService) RegisterListControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListController) {
	group.GET(basePath+"/:id", ctr.Get)
	group.GET(basePath, ctr.GetList)
	group.PUT(basePath, ctr.Put)
	group.PUT(basePath+"/batch", ctr.PutList)
	group.POST(basePath+"/:id", ctr.Post)
	group.POST(basePath, ctr.PostList)
	group.DELETE(basePath+"/:id", ctr.Delete)
	group.DELETE(basePath, ctr.DeleteList)
}

func (svc *RouterService) RegisterActionControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ActionController) {
	for _, action := range ctr.Actions() {
		routerPath := path.Join(basePath, action.Path)
		switch action.Method {
		case http.MethodGet:
			group.GET(routerPath, action.HandlerFunc)
		case http.MethodPut:
			group.PUT(routerPath, action.HandlerFunc)
		case http.MethodPost:
			group.POST(routerPath, action.HandlerFunc)
		case http.MethodDelete:
			group.DELETE(routerPath, action.HandlerFunc)
		}
	}
}

func (svc *RouterService) RegisterListActionControllerToGroup(group *gin.RouterGroup, basePath string, ctr controllers.ListActionController) {
	svc.RegisterListControllerToGroup(group, basePath, ctr)
	svc.RegisterActionControllerToGroup(group, basePath, ctr)
}

func (svc *RouterService) RegisterHandlerToGroup(group *gin.RouterGroup, path string, method string, handler gin.HandlerFunc) {
	switch method {
	case http.MethodGet:
		group.GET(path, handler)
	case http.MethodPut:
		group.PUT(path, handler)
	case http.MethodPost:
		group.POST(path, handler)
	case http.MethodDelete:
		group.DELETE(path, handler)
	default:
		log.Warn(fmt.Sprintf("%s is not a valid http method", method))
	}
}

func InitRoutes(app *gin.Engine) (err error) {
	// routes groups
	groups := NewRouterGroups(app)

	// router service
	svc := NewRouterService(app)

	// node
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/nodes", controllers.NodeController)

	// project
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/projects", controllers.ProjectController)

	// user
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/users", controllers.UserController)

	// spider
	svc.RegisterListActionControllerToGroup(groups.AuthGroup, "/spiders", controllers.SpiderController)

	// task
	svc.RegisterListActionControllerToGroup(groups.AuthGroup, "/tasks", &controllers.TaskController)

	// tag
	svc.RegisterListControllerToGroup(groups.AuthGroup, "/tags", controllers.TagController)

	// login
	svc.RegisterActionControllerToGroup(groups.AnonymousGroup, "/", controllers.LoginController)

	// color
	svc.RegisterActionControllerToGroup(groups.AuthGroup, "/colors", controllers.ColorController)

	return nil
}
