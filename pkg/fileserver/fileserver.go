package fileserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
)

// FileServer handles serving dynamically generated images via HTTP
type FileServer struct {
	server    *http.Server
	basePath  string
	port      int
	running   bool
	urlPrefix string
	mutex     sync.Mutex
	prefs     *preferences.Preferences
}

// NewFileServer creates a new file server instance
func NewFileServer(prefs *preferences.Preferences) *FileServer {
	return &FileServer{
		basePath:  "", // Will be set in Start() method
		port:      0,  // Will be determined dynamically
		running:   false,
		urlPrefix: "/images/",
		prefs:     prefs,
	}
}

// Start initializes and starts the file server
func (fs *FileServer) Start() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.running {
		return nil // Already running
	}

	// Determine storage location
	basePath, err := fs.getStorageLocation()
	if err != nil {
		return fmt.Errorf("failed to determine storage location: %w", err)
	}
	fs.basePath = basePath

	// Ensure directory exists
	if err := os.MkdirAll(fs.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}

	// Create sample files in the base path
	if err := CreateSampleFiles(fs.basePath); err != nil {
		log.Printf("Warning: failed to create sample files: %v", err)
		// Continue even if sample creation fails - this is non-critical
	} else {
		log.Printf("Sample files created successfully in %s", filepath.Join(fs.basePath, "samples"))
	}

	// Find available port
	port, err := findAvailablePort(8090)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}
	fs.port = port

	// Create server
	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir(fs.basePath))

	// Apply middleware for security and logging
	handler := LoggingMiddleware(SecurityMiddleware(http.StripPrefix(fs.urlPrefix, fileHandler)))
	mux.Handle(fs.urlPrefix, handler)

	// Configure server
	fs.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", fs.port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		fs.running = true
		log.Printf("File server started at http://127.0.0.1:%d serving files from %s", fs.port, fs.basePath)
		if err := fs.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("File server error: %v", err)
		}
		fs.running = false
	}()

	return nil
}

// Stop gracefully shuts down the file server
func (fs *FileServer) Stop() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if !fs.running || fs.server == nil {
		return nil // Not running, nothing to do
	}

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Stopping file server on port %d", fs.port)

	// Shutdown the server
	err := fs.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down file server: %w", err)
	}

	fs.running = false
	return nil
}

// UpdateBasePath changes the base directory and restarts the server
func (fs *FileServer) UpdateBasePath(newPath string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Verify that the new path exists
	if _, err := os.Stat(newPath); err != nil {
		if os.IsNotExist(err) {
			// Try to create the directory
			if err := os.MkdirAll(newPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory at new path: %w", err)
			}
		} else {
			return fmt.Errorf("error checking new path: %w", err)
		}
	}

	// If server is running, we need to restart it with new path
	if fs.running {
		// Shutdown the current server
		log.Printf("Restarting file server with new base path: %s", newPath)

		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server
		if err := fs.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("error shutting down file server: %w", err)
		}

		fs.running = false
	}

	// Update the base path
	fs.basePath = newPath

	// If server was running, restart it
	if fs.server != nil {
		return fs.Start()
	}

	return nil
}

// GetBasePath returns the current base path of the file server
func (fs *FileServer) GetBasePath() string {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return fs.basePath
}

// Helper function to find an available port
func findAvailablePort(basePort int) (int, error) {
	for port := basePort; port < basePort+100; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", basePort, basePort+100)
}

// getStorageLocation determines where images should be stored
func (fs *FileServer) getStorageLocation() (string, error) {
	// Check if there are recent projects
	if len(fs.prefs.App.RecentProjects) > 0 {
		// Take the directory of the most recent project
		recentPath := fs.prefs.App.RecentProjects[0]
		// Extract directory from file path
		dir := filepath.Dir(recentPath)
		if _, err := os.Stat(dir); err == nil {
			return filepath.Join(dir, "images"), nil
		}
	}

	// Fall back to default location
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	appId := preferences.GetAppId()
	return filepath.Join(userDir, "Documents", appId.AppName, "images"), nil
}

// GetURL returns the URL for accessing a specific image
func (fs *FileServer) GetURL(relativePath string) string {
	// No need for locking here - we're just reading values
	if !fs.running || fs.port == 0 {
		return "" // Server not running, can't generate URL
	}

	// Clean the path to prevent directory traversal
	relativePath = path.Clean(relativePath)

	// Ensure the path doesn't start with a slash
	relativePath = strings.TrimPrefix(relativePath, "/")

	// Construct the URL
	return fmt.Sprintf("http://127.0.0.1:%d%s%s",
		fs.port, fs.urlPrefix, relativePath)
}
