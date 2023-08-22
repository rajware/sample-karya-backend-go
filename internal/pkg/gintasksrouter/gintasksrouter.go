package gintasksrouter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
)

type Server struct {
	tasks  *tasks.Tasks
	router *gin.Engine
	Port   int
}

type ApiStatus struct {
	Data    interface{} `json:"data"`
	Error   int         `json:"error"`
	Message string      `json:"message"`
}

func success(data interface{}) ApiStatus {
	return ApiStatus{
		Data:    data,
		Error:   0,
		Message: "success",
	}
}

func (s *Server) getAllTasks(c *gin.Context) {
	alltasks, err := s.tasks.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ApiStatus{
			Data:    nil,
			Error:   http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, success(alltasks))
}

func (s *Server) getTaskById(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiStatus{
			Data:    nil,
			Error:   http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := s.tasks.GetById(uint(id))
	if errors.Is(err, tasks.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, ApiStatus{
			Data:    nil,
			Error:   http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ApiStatus{
			Data:    nil,
			Error:   http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, success(task))
}

func (s *Server) addTask(c *gin.Context) {
	var newTask tasks.Task

	err := c.BindJSON(&newTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiStatus{
			Data:    nil,
			Error:   http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	addedTask, err := s.tasks.NewTask(newTask.Description, newTask.Deadline)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ApiStatus{
			Data:    nil,
			Error:   http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, success(addedTask))
}

func (s *Server) updateTask(c *gin.Context) {
	var theTask tasks.Task

	err := c.BindJSON(&theTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiStatus{
			Data:    nil,
			Error:   http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	updatedTask, err := s.tasks.Update(&theTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ApiStatus{
			Data:    nil,
			Error:   http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, success(updatedTask))
}

func (s *Server) deleteTaskById(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ApiStatus{
			Data:    nil,
			Error:   http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	deletedTask, err := s.tasks.DeleteById(uint(id))
	if errors.Is(err, tasks.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, ApiStatus{
			Data:    nil,
			Error:   http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ApiStatus{
			Data:    nil,
			Error:   http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, success(deletedTask))
}

func (r *Server) Run() {
	if r.router == nil {
		panic("router not set up")
	}

	port := ":" + strconv.FormatInt(int64(r.Port), 10)

	server := http.Server{Addr: port, Handler: r.router}

	// Use a channel to signal server closure
	serverClosed := make(chan struct{})

	// Handle OS signals for graceful shutdown
	go func() {
		signalReceived := make(chan os.Signal, 1)

		// Handle SIGINT
		signal.Notify(signalReceived, os.Interrupt)
		// Handle SIGTERM
		signal.Notify(signalReceived, syscall.SIGTERM)

		// Wait for signal
		<-signalReceived

		log.Println("Server shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Fatalf("Error during HTTP server shutdown: %v.", err)
		}

		close(serverClosed)
	}()

	// Start listening using the server
	log.Printf("Server starting on port %v...\n", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("The server failed with the following error: %v.\n", err)
	}

	<-serverClosed

	log.Println("Server shut down.")
}

var defaultServer = &Server{Port: 8080}

func New() *Server {

	repo := gormrepo.New()
	if repo == nil {
		panic("could not initialize data store")
	}
	defaultServer.tasks = tasks.New(repo)
	defaultServer.router = gin.Default()
	defaultServer.router.GET("/tasks", defaultServer.getAllTasks)
	defaultServer.router.GET("/tasks/:id", defaultServer.getTaskById)
	defaultServer.router.POST("/tasks", defaultServer.addTask)
	defaultServer.router.PUT("/tasks", defaultServer.updateTask)
	defaultServer.router.DELETE("/tasks/:id", defaultServer.deleteTaskById)

	return defaultServer
}
