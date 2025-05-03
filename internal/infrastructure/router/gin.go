package router

import (
	"glbackend/internal/adapters/api/action"
	"glbackend/internal/adapters/api/middleware"
	"net/http"

	"fmt"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/config"

	"glbackend/internal/repo"
	"glbackend/internal/usecase"
	"strconv"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type ginEngine struct {
	cfg         config.Config
	router      *gin.Engine
	log         logging.Logger
	transaction repo.GTransaction
	repo        repo.Repo
	ctxTimeout  time.Duration
	enforcer    *casbin.Enforcer
}

func newGinServer(
	cfg config.Config,
	log logging.Logger,
	db repo.GSQL,
	transtaction repo.GTransaction,
	ctxTimeout time.Duration,
) *ginEngine {
	enforcer := middleware.NewCasbin()
	return &ginEngine{
		cfg:         cfg,
		router:      gin.Default(),
		log:         log,
		transaction: transtaction,
		repo:        repo.NewRepo(db),
		ctxTimeout:  ctxTimeout,
		enforcer:    enforcer,
	}
}

func (g ginEngine) Listen() {
	gin.SetMode(gin.ReleaseMode)
	// gin.Recovery()
	g.router.Use(gin.Recovery())

	g.setupRoutes(g.router)

	g.router.Run(":" + strconv.Itoa(int(g.cfg.AppPort)))
}

func (g ginEngine) setupRoutes(router *gin.Engine) {
	router.Use(g.CORSMiddleware())
	root := router.Group("/api/")

	root.GET("/swag/swagger.json", func(c *gin.Context) {
		filePath := "./swag/swagger.json"
		c.File(filePath)
	})
	root.GET("/docs", g.buildAPIDocsAction(strconv.Itoa(int(g.cfg.AppPort))))

	products := root.Group("/products") //.Use(g.buildSessionMiddleware())
	{
		products.GET("", g.buildProductsFindAllAction)
		products.GET("/:id", g.buildProductsFindByIDAction)
		products.POST("", g.buildProductsCreateAction)
		products.PATCH("/:id", g.buildProductsUpdateAction)
		products.DELETE("/:id", g.buildProductsDeleteAction)
	}

	tags := root.Group("/tags") //.Use(g.buildSessionMiddleware())
	{
		tags.GET("", g.buildTagsFindAllAction)
		tags.POST("", g.buildTagsCreateAction)
		tags.PATCH("/:id", g.buildTagsUpdateAction)
		tags.DELETE("/:id", g.buildTagsDeleteAction)
	}

	categories := root.Group("/categories") //.Use(g.buildSessionMiddleware())
	{
		categories.GET("", g.buildCategoriesFindAllAction)
		categories.POST("", g.buildCategoriesCreateAction)
		categories.PATCH("/:id", g.buildCategoriesUpdateAction)
		categories.DELETE("/:id", g.buildCategoriesDeleteAction)
	}
}

// func paramsPageAndLimit(c *gin.Context) {
// 	q := c.Request.URL.Query()
// 	if _, exists := q["page"]; !exists {
// 		q.Set("page", "1")
// 	}
// 	if _, exists := q["limit"]; !exists {
// 		q.Set("limit", "10")
// 	}
// 	c.Request.URL.RawQuery = q.Encode()
// }

func paramsToQuery(c *gin.Context, params []string) {
	q := c.Request.URL.Query()
	for _, p := range params {
		q.Set(p, c.Param(p))
	}
	c.Request.URL.RawQuery = q.Encode()
}

/* ==================== ADMIN ==================== */

// @Summary	Получение списка продуктов
// @securityDefinitions.apikey	<BearerAuth>
// @Tags		products
// @Success	200		{array}	entities.FullProduct
// @Failure	400		{object}	response.Error
// @Router		/products [get]
func (g ginEngine) buildProductsFindAllAction(c *gin.Context) {
	var (
		uc = usecase.NewProductsFindAllUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewProductsFindAllAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Получение продукта по ID
// @securityDefinitions.apikey	<BearerAuth>
// @Tags		products
// @Param		id	path		string	true	"ID продукта"
// @Success	200	{object}	entities.FullProduct
// @Failure	400	{object}	response.Error
// @Router		/products/{id} [get]
func (g ginEngine) buildProductsFindByIDAction(c *gin.Context) {
	var (
		uc = usecase.NewProductsFindByIDUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewProductsFindByIDAction(uc, g.log)
	)

	paramsToQuery(c, []string{"id"})
	act.Execute(c.Writer, c.Request)
}

// @Summary	Создание продукта
// @Tags		products
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		input	body		usecase.ProductsCreateInput	true	"Данные продукта"
// @Success	200		{object}	entities.FullProduct
// @Failure	400		{object}	response.Error
// @Router		/products [post]
func (g ginEngine) buildProductsCreateAction(c *gin.Context) {
	var (
		uc = usecase.NewProductsCreateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewProductsCreateAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Обновление продукта
// @Tags		products
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id		path		string							true	"ID продукта"
// @Param		input	body		usecase.ProductsUpdateInput	true	"Данные продукта"
// @Success	200		{object}	entities.FullProduct
// @Failure	400		{object}	response.Error
// @Router		/products/{id} [patch]
func (g ginEngine) buildProductsUpdateAction(c *gin.Context) {
	var (
		uc = usecase.NewProductsUpdateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewProductsUpdateAction(uc, g.log)
	)

	paramsToQuery(c, []string{"id"})
	act.Execute(c.Writer, c.Request)
}

// @Summary	Удаление продукта
// @Tags		products
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"ID продукта"
// @Success	200	{object}	bool
// @Failure	400	{object}	response.Error
// @Router		/products/{id} [delete]
func (g ginEngine) buildProductsDeleteAction(c *gin.Context) {
	var (
		uc = usecase.NewProductsDeleteUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewProductsDeleteAction(uc, g.log)
	)

	paramsToQuery(c, []string{"id"})
	act.Execute(c.Writer, c.Request)
}

// @Summary	Получение списка тегов
// @securityDefinitions.apikey	<BearerAuth>
// @Tags		tags
// @Success	200		{array}	entities.Tag
// @Failure	400		{object}	response.Error
// @Router		/tags [get]
func (g ginEngine) buildTagsFindAllAction(c *gin.Context) {
	var (
		uc = usecase.NewTagsFindAllUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewTagsFindAllAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Создание тега
// @Tags		tags
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		input	body		usecase.TagsCreateInput	true	"Данные тега"
// @Success	200		{object}	entities.Tag
// @Failure	400		{object}	response.Error
// @Router		/tags [post]
func (g ginEngine) buildTagsCreateAction(c *gin.Context) {
	var (
		uc = usecase.NewTagsCreateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewTagsCreateAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Обновление тега
// @Tags		tags
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id	query		string	true	"ID тега"
// @Param		input	body		usecase.TagsUpdateInput	true	"Данные тега"
// @Success	200		{object}	entities.Tag
// @Failure	400		{object}	response.Error
// @Router		/tags [patch]
func (g ginEngine) buildTagsUpdateAction(c *gin.Context) {
	var (
		uc = usecase.NewTagsUpdateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewTagsUpdateAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Удаление тега
// @Tags		tags
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"ID тега"
// @Success	200	{object}	bool
// @Failure	400	{object}	response.Error
// @Router		/tags/{id} [delete]
func (g ginEngine) buildTagsDeleteAction(c *gin.Context) {
	var (
		uc = usecase.NewTagsDeleteUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewTagsDeleteAction(uc, g.log)
	)

	paramsToQuery(c, []string{"id"})
	act.Execute(c.Writer, c.Request)
}

// @Summary	Получение списка категорий
// @securityDefinitions.apikey	<BearerAuth>
// @Tags		categories
// @Success	200		{array}	entities.Category
// @Failure	400		{object}	response.Error
// @Router		/categories [get]
func (g ginEngine) buildCategoriesFindAllAction(c *gin.Context) {
	var (
		uc = usecase.NewCategoriesFindAllUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewCategoriesFindAllAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Создание категории
// @Tags		categories
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		input	body		usecase.CategoriesCreateInput	true	"Данные категории"
// @Success	200		{object}	entities.Category
// @Failure	400		{object}	response.Error
// @Router		/categories [post]
func (g ginEngine) buildCategoriesCreateAction(c *gin.Context) {
	var (
		uc = usecase.NewCategoriesCreateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewCategoriesCreateAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Обновление категории
// @Tags		categories
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id	query		string	true	"ID категории"
// @Param		input	body		usecase.CategoriesUpdateInput	true	"Данные категории"
// @Success	200		{object}	entities.Category
// @Failure	400		{object}	response.Error
// @Router		/categories [patch]
func (g ginEngine) buildCategoriesUpdateAction(c *gin.Context) {
	var (
		uc = usecase.NewCategoriesUpdateUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewCategoriesUpdateAction(uc, g.log)
	)

	act.Execute(c.Writer, c.Request)
}

// @Summary	Удаление категории
// @Tags		categories
// @securityDefinitions.apikey	<BearerAuth>
// @Accept		json
// @Produce	json
// @Param		id	path		string	true	"ID категории"
// @Success	200	{object}	bool
// @Failure	400	{object}	response.Error
// @Router		/categories/{id} [delete]
func (g ginEngine) buildCategoriesDeleteAction(c *gin.Context) {
	var (
		uc = usecase.NewCategoriesDeleteUsecase(
			g.repo,
			g.ctxTimeout,
		)
		act = action.NewCategoriesDeleteAction(uc, g.log)
	)

	paramsToQuery(c, []string{"id"})
	act.Execute(c.Writer, c.Request)
}

/* ==================== MIDDLEWARES ==================== */
// func (g *ginEngine) buildSessionMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		const logKey = "session_check"

// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			logging.NewError(
// 				g.log,
// 				errors.New("missing token"),
// 				logKey,
// 				http.StatusUnauthorized,
// 			).Log("missing token")

// 			response.NewError(c.Writer, errors.New("missing token"), http.StatusUnauthorized)
// 			c.Abort()
// 			return
// 		}
// 		// Ожидаемый формат: "Bearer <token>"
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			logging.NewError(
// 				g.log,
// 				errors.New("invalid authorization header format"),
// 				logKey,
// 				http.StatusUnauthorized,
// 			).Log("invalid authorization header format")

// 			response.NewError(c.Writer, errors.New("invalid authorization header format"), http.StatusUnauthorized)
// 			c.Abort()
// 			return
// 		}
// 		tokenString := parts[1]
// 		claims, err := utils.ValidateAccessToken(tokenString) // TODO: переделать это отдельным usecase или ValidateAccessToken в pkg вынести
// 		if err != nil {
// 			logging.NewError(
// 				g.log,
// 				errors.New("invalid token"),
// 				logKey,
// 				http.StatusUnauthorized,
// 			).Log("invalid token")

// 			response.NewError(c.Writer, errors.New("invalid token"), http.StatusUnauthorized)
// 			c.Abort()
// 			return
// 		}

// 		ctx := context.WithValue(c.Request.Context(), middleware.SessionContextKey, claims)
// 		c.Request = c.Request.WithContext(ctx)
// 	}
// }

// func (g *ginEngine) buildCasbinMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		const logKey = "casbin_middleware"
// 		session, ok := c.Request.Context().Value(middleware.SessionContextKey).(utils.JWTClaims)
// 		if !ok {
// 			logging.NewError(
// 				g.log,
// 				errors.New("missing session"),
// 				logKey,
// 				http.StatusForbidden,
// 			).Log("missing role in Authorization context")

// 			response.NewError(c.Writer, errors.New("missing role"), http.StatusUnauthorized)
// 			c.Abort()
// 			return
// 		}

// 		allowed, err := g.enforcer.Enforce(session.UserType, c.Request.URL.Path, c.Request.Method)
// 		fmt.Printf("session.TypeUser: %+v, allowed: %+v\n", session.UserType, allowed)
// 		if err != nil {
// 			logging.NewError(
// 				g.log,
// 				err,
// 				logKey,
// 				http.StatusBadRequest,
// 			).Log("missing role in Authorization header")

// 			response.NewError(c.Writer, err, http.StatusBadRequest)
// 			c.Abort()
// 			return
// 		}

// 		if !allowed {
// 			logging.NewError(
// 				g.log,
// 				errors.New("forbidden_access"),
// 				logKey,
// 				http.StatusForbidden,
// 			).Log("forbidden access for user role " + session.UserType)

// 			response.NewError(c.Writer, errors.New("forbidden_access"), http.StatusForbidden)
// 			c.Abort()
// 			return
// 		}
// 	}
// }

func (g ginEngine) CORSMiddleware() gin.HandlerFunc {
	// Список разрешенных доменов

	var allowedOrigins = []string{
		"http://localhost:83",
		"http://localhost:83/",
		"http://localhost:4200",
		"http://localhost:4200/",
		"https://glbackend.ru/",
		"https://glbackend.ru",
	}

	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		allowed := false

		// Проверяем, разрешён ли домен, указанный в Origin заголовке запроса
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		// if origin == "" {
		// 	referer := c.Request.Header.Get("Referer")
		// 	if referer != "" {
		// 		refererURL, err := url.Parse(referer)
		// 		if err == nil {
		// 			origin = refererURL.Scheme + "://" + refererURL.Host
		// 		}
		// 	}
		// }
		// fmt.Println("Origin/Referer:", origin) // лог для отладки

		// allowed = false

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed && origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

/* ==================== API REFERENCE ==================== */

/*
   ThemeDefault    ThemeId = "default"
   ThemeAlternate  ThemeId = "alternate"
   ThemeMoon       ThemeId = "moon"
   ThemePurple     ThemeId = "purple"
   ThemeSolarized  ThemeId = "solarized"
   ThemeBluePlanet ThemeId = "bluePlanet"
   ThemeDeepSpace  ThemeId = "deepSpace"
   ThemeSaturn     ThemeId = "saturn"
   ThemeKepler     ThemeId = "kepler"
   ThemeMars       ThemeId = "mars"
   ThemeNone       ThemeId = "none"
*/

func (g *ginEngine) buildAPIDocsAction(port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "http://localhost:" + port + "/api/swag/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "glbackend API (i love kitty)",
			},
			DarkMode: true,
			Theme:    "ThemeBluePlanet", // "bluePlanet",
		})

		if err != nil {
			fmt.Printf("%+v\n", err)
		}

		fmt.Fprintln(c.Writer, htmlContent)
	}
}
