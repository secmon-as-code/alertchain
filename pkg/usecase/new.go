package usecase

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

//go:embed templates/**
var embedTemplatesFS embed.FS

func NewPolicyDirectory(ctx context.Context, dir string) error {
	if err := copyEmbeddedFiles(ctx, embedTemplatesFS, "templates", dir); err != nil {
		return err
	}
	return nil
}

func copyEmbeddedFiles(ctx context.Context, efs embed.FS, srcDir, dstDir string) error {
	logger := ctxutil.Logger(ctx)
	return fs.WalkDir(efs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return goerr.Wrap(err, "failed to walk through embedded directory")
		}
		if d.IsDir() {
			return nil
		}

		srcPath := filepath.Clean(path)
		fd, err := efs.Open(srcPath)
		if err != nil {
			return goerr.Wrap(err, "failed to open embedded file")
		}
		defer utils.SafeClose(ctx, fd)

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return goerr.Wrap(err, "failed to get relative path")
		}
		dstPath := filepath.Join(dstDir, relPath)

		if err := os.MkdirAll(filepath.Dir(filepath.Clean(dstPath)), os.ModePerm); err != nil {
			return err
		}

		w, err := os.Create(dstPath)
		if err != nil {
			return goerr.Wrap(err, "failed to create file")
		}
		defer utils.SafeClose(ctx, w)

		if _, err := io.Copy(w, fd); err != nil {
			return goerr.Wrap(err, "failed to copy file")
		}
		logger.Info("Copy file", "path", dstPath)

		return nil
	})
}
