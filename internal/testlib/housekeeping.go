package testlib

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/config"
)

// CleanTestDirectory will remove all directories created when a server is
// started, such as storage and cache directory. Used to clean up directories
// where tests are run.
func CleanTestDirectory(r *app.Registry, cfg *config.Config) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cleanStorages(cwd, cfg.Private.Storage)
	cleanCacheStorage(cwd, cfg.Private.Cache, r)
	cleanLogs(cwd, cfg.Private.Log, r)
	cleanTestingDirectory(cwd)
}

func cleanStorages(cwd string, cfg config.StorageConfig) {
	for _, disk := range cfg.Disks {
		if strings.ToLower(disk.Driver) != "local" {
			continue
		}

		var p string
		if len(disk.Dir) > 0 {
			p = path.Join(cwd, disk.Dir)
		} else {
			p = path.Join(cwd, "./storage")
		}

		err := os.RemoveAll(p)
		if err != nil {
			panic(err)
		}
	}

	publicPath := path.Join(cwd, "./storage/public")
	err := os.RemoveAll(publicPath)
	if err != nil {
		panic(err)
	}

	removeIfEmpty(path.Join(cwd, "./storage"))
}

func cleanCacheStorage(cwd string, cfg config.CacheConfig, r *app.Registry) {
	if strings.ToLower(cfg.Engine) == "badger" && !cfg.Badger.InMemory {
		if err := r.Cache.Close(); err != nil {
			panic(err)
		}
		if err := os.RemoveAll(path.Join(cwd, cfg.Badger.Path)); err != nil {
			panic(err)
		}
	}
	removeIfEmpty(path.Join(cwd, "./cache"))
}

func cleanTestingDirectory(cwd string) {
	err := os.RemoveAll(path.Join(cwd, "./testing"))
	if err != nil {
		panic(err)
	}
}

func cleanLogs(cwd string, cfg config.LogConfig, r *app.Registry) {
	for _, loggers := range r.Log.Writers {
		for _, logger := range loggers {
			logger.Close()
		}
	}
	for _, log := range cfg.Writers {
		if strings.TrimSpace(log.LogRotatingFileWriterConfig.Filepath) != "" {
			if err := os.RemoveAll(path.Join(cwd, log.LogRotatingFileWriterConfig.Filepath)); err != nil {
				panic(err)
			}
		}
		if strings.TrimSpace(log.LogSingleFileWriterConfig.Filepath) != "" {
			if err := os.RemoveAll(path.Join(cwd, log.LogSingleFileWriterConfig.Filepath)); err != nil {
				panic(err)
			}
		}
	}
	removeIfEmpty(path.Join(cwd, "./logs"))
}

func removeIfEmpty(dir string) {
	e, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		panic(err)
	}
	if len(e) == 0 {
		err := os.Remove(dir)
		if err != nil {
			panic(err)
		}
	}
}
