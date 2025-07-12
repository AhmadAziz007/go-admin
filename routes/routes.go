package routes

import (
	"go-admin/controller"
	"go-admin/middlewares"
	"go-admin/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, minioService *service.MinioService) {
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	roleService := service.NewRoleService(db)
	roleController := controller.NewRoleController(roleService)

	productService := service.NewProductService(db, minioService)
	productController := controller.NewProductController(productService)

	transactionService := service.NewTransactionService(db)
	transactionController := controller.NewTransactionController(transactionService)

	customerService := service.NewCustomerService(db)
	customerController := controller.NewCustomerController(customerService)

	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middlewares.IsAuthenticated)

	app.Put("/api/users/info", controller.UpdateInfo)
	app.Put("/api/users/password", controller.UpdatePassword)

	app.Get("/api/user", controller.User)
	app.Post("/api/logout", controller.Logout)

	//app.Get("/api/users", userController.AllUsers)
	//app.Post("/api/users", userController.CreateUser)
	//app.Get("/api/users/:id", userController.GetUser)
	//app.Put("/api/users/:id", userController.UpdateUser)
	//app.Delete("/api/users/:id", userController.DeleteUser)

	//users
	app.Get("/api/users", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "users"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return userController.AllUsers(c)
	})

	app.Post("/api/users", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "users"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return userController.CreateUser(c)
	})

	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "users"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return userController.GetUser(c)
	})

	app.Put("/api/users/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "users"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return userController.UpdateUser(c)
	})

	app.Delete("/api/users/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "users"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return userController.DeleteUser(c)
	})

	//app.Get("/api/roles", roleController.AllRoles)
	//app.Post("/api/roles", roleController.CreateRole)
	//app.Get("/api/roles/:id", roleController.GetRole)
	//app.Put("/api/roles/:id", roleController.UpdateRole)
	//app.Delete("/api/roles/:id", roleController.DeleteRole)

	//roles
	app.Get("/api/roles", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "roles"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return roleController.AllRoles(c)
	})

	app.Post("/api/roles", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "roles"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return roleController.CreateRole(c)
	})

	app.Get("/api/roles/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "roles"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return roleController.GetRole(c)
	})

	app.Put("/api/roles/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "roles"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return roleController.UpdateRole(c)
	})

	app.Delete("/api/roles/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "roles"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return roleController.DeleteRole(c)
	})

	app.Get("/api/dropdown/customers", customerController.DropdownCustomers)

	app.Get("/api/customers", customerController.AllCustomers)
	app.Post("/api/customers", customerController.CreateCustomer)
	app.Get("/api/customers/:id", customerController.GetCustomer)
	app.Put("/api/customers/:id", customerController.UpdateCustomer)
	app.Delete("/api/customers/:id", customerController.DeleteCustomer)

	//products
	app.Post("/api/products", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "products"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return productController.Create(c)
	})

	app.Put("/api/products/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "products"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return productController.Update(c)
	})

	app.Delete("/api/products/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "products"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return productController.Delete(c)
	})

	app.Get("/api/products/:id", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "products"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return productController.GetByID(c)
	})

	app.Get("/api/products", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "products"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return productController.GetAll(c)
	})

	app.Get("/api/permissions", controller.AllPermissions)

	//app.Get("/api/transactions/searchProduct", transactionController.SearchProduct)
	//app.Post("/api/transactions/addToCart", transactionController.AddToCart)
	//app.Delete("/api/transactions/destroyCart", transactionController.DestroyCart)
	//app.Get("/api/transactions/getCart", transactionController.GetCart)
	//app.Post("/api/transactions/payOrder", transactionController.PayOrder)

	//transaction
	app.Get("/api/transactions/searchProduct", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "transactions"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return transactionController.SearchProduct(c)
	})

	app.Post("/api/transactions/addToCart", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "transactions"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return transactionController.AddToCart(c)
	})

	app.Delete("/api/transactions/destroyCart", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "transactions"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return transactionController.DestroyCart(c)
	})

	app.Get("/api/transactions/getCart", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "transactions"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return transactionController.GetCart(c)
	})

	app.Post("/api/transactions/payOrder", func(c *fiber.Ctx) error {
		if err := middlewares.IsAuthorized(c, "transactions"); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return transactionController.PayOrder(c)
	})

}
