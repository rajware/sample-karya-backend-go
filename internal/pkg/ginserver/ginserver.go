package ginserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
)

// Server is an HTTP server which serves static files
// and a REST api to perform CRUD operations on Tasks.
// It depends on an implementation of tasks.TaskRepository
// to store and retrieve Task data.
//
// It handles the following routes:
//
// GET /tasks
// GET /tasks/:id
// POST /tasks
// PUT /tasks
// DELETE /tasks/:id
type Server struct {
	tasks  *tasks.Tasks
	router *gin.Engine
	port   int
}

type apiStatus struct {
	Data    interface{} `json:"data"`
	Error   int         `json:"error"`
	Message string      `json:"message"`
}

func succeed(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, apiStatus{
		Data:    data,
		Error:   0,
		Message: "success",
	})
}

func fail(c *gin.Context, statuscode int, err error) {
	c.AbortWithStatusJSON(
		statuscode,
		apiStatus{
			Data:    nil,
			Error:   statuscode,
			Message: err.Error(),
		},
	)
}

func (s *Server) getAllTasks(c *gin.Context) {
	alltasks, err := s.tasks.GetAll()
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	succeed(c, alltasks)
}

func (s *Server) getTaskByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	task, err := s.tasks.GetByID(uint(id))
	if errors.Is(err, tasks.ErrNotFound) {
		fail(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	succeed(c, task)
}

func (s *Server) addTask(c *gin.Context) {
	var newTask tasks.Task

	err := c.BindJSON(&newTask)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	addedTask, err := s.tasks.NewTask(newTask.Description, newTask.Deadline)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	succeed(c, addedTask)
}

func (s *Server) updateTask(c *gin.Context) {
	var theTask tasks.Task

	err := c.BindJSON(&theTask)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	updatedTask, err := s.tasks.Update(&theTask)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	succeed(c, updatedTask)
}

func (s *Server) deleteTaskByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	deletedTask, err := s.tasks.DeleteByID(uint(id))
	if errors.Is(err, tasks.ErrNotFound) {
		fail(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	succeed(c, deletedTask)
}

// Run starts listening and serving HTTP Requests. It also listens for
// SIGINT and SIGTERM, and stops listening if it does. Note: this method
// will block the calling goroutine indefinitely unless an error or an
// interrupt signal happens.
func (s *Server) Run() {
	if s.router == nil {
		panic("router not set up")
	}

	port := ":" + strconv.FormatInt(int64(s.port), 10)
	server := http.Server{Addr: port, Handler: s.router}

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

		slog.Info("Server shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			slog.Error("problem during HTTP server shutdown.", "Error", err)
			os.Exit(1)
		}

		close(serverClosed)
	}()

	// Start listening using the server
	slog.Info("Server starting.", "Port", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server failed.", "Error", err)
		os.Exit(1)
	}

	<-serverClosed

	slog.Info("Server shut down.")
}

// EnsureSubdirectory ensures the presence of a subdirectory
// under the current directory.
func EnsureSubdirectory(subdirname string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	datadir := filepath.Join(wd, subdirname)
	err = os.MkdirAll(datadir, 0755)
	if err != nil {
		return "", err
	}

	return datadir, nil
}

// New returns a Server instance. An implementation of tasks.TaskRepository
// is required, to ensure task data can be stored and retrieved.
func New(repo tasks.TaskRepository, port int) (*Server, error) {
	newServer := &Server{port: port}

	if repo == nil {
		return nil, errors.New("could not initialize data store")
	}
	newServer.tasks = tasks.New(repo)

	handler := gin.Default()

	// Set up static content serving
	staticDir, err := EnsureSubdirectory("wwwroot")
	if err != nil {
		return nil, errors.New("could not set up static serving:" + err.Error())
	}
	handler.Use(static.Serve("/", static.LocalFile(staticDir, true)))

	// Set up Tasks api
	handler.GET("/tasks", newServer.getAllTasks)
	handler.GET("/tasks/:id", newServer.getTaskByID)
	handler.POST("/tasks", newServer.addTask)
	handler.PUT("/tasks", newServer.updateTask)
	handler.DELETE("/tasks/:id", newServer.deleteTaskByID)

	newServer.router = handler
	return newServer, nil
}
