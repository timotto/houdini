package houdini

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"code.cloudfoundry.org/garden"
)

func (c *container) setup() error {
	for _, bm := range c.spec.BindMounts {
		if bm.Mode == garden.BindMountModeRO {
			return errors.New("read-only bind mounts are unsupported")
		}

		dest := filepath.Join(c.workDir, bm.DstPath)
		_, err := os.Stat(dest)
		if err == nil {
			err = os.Remove(dest)
			if err != nil {
				return fmt.Errorf("failed to remove destination for bind mount: %s", err)
			}
		}

		err = os.MkdirAll(filepath.Dir(dest), 0755)
		if err != nil {
			return fmt.Errorf("failed to create parent dir for bind mount: %s", err)
		}

		// darwin hard-links support directories
		err = syscall.Link(bm.SrcPath, dest)
		if err != nil {
			return fmt.Errorf("failed to create hardlink for bind mount: %s", err)
		}
	}

	return nil
}