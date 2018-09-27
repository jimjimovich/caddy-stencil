// Copyright 2018 Jim Mendenhall
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stencil

import (
	"bytes"
	"mime"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("stencil", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

// setup configures a new Stencil middleware instance.
func Setup(c *caddy.Controller) error {
	stconfigs, err := stencilParse(c)
	if err != nil {
		return err
	}

	cfg := httpserver.GetConfig(c)

	// Add json mime type in case it is not available on the system
	mime.AddExtensionType(".json", "application/json")

	st := Stencil{
		Root:    cfg.Root,
		FileSys: http.Dir(cfg.Root),
		Configs: stconfigs,
		BufPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		st.Next = next
		return st
	})

	return nil
}

func stencilParse(c *caddy.Controller) ([]*Config, error) {
	var stconfigs []*Config

	for c.Next() {
		st := &Config{
			Extensions:    make(map[string]struct{}),
			Template:      GetDefaultTemplate(),
			TemplateFiles: make(map[string]*CachedFileInfo),
		}

		// Get the path scope
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			st.PathScope = "/"
		case 1:
			st.PathScope = args[0]
		default:
			return stconfigs, c.ArgErr()
		}

		// Load any other configuration parameters
		for c.NextBlock() {
			if err := loadParams(c, st); err != nil {
				return stconfigs, err
			}
		}

		// If no extensions were specified, assume some defaults
		if len(st.Extensions) == 0 {
			st.Extensions[".html"] = struct{}{}
			st.Extensions[".json"] = struct{}{}
		}

		stconfigs = append(stconfigs, st)
	}

	return stconfigs, nil
}

func loadParams(c *caddy.Controller, stc *Config) error {
	cfg := httpserver.GetConfig(c)

	switch c.Val() {
	case "ext":
		for _, ext := range c.RemainingArgs() {
			stc.Extensions[ext] = struct{}{}
		}
		return nil
	case "template":
		tArgs := c.RemainingArgs()
		switch len(tArgs) {
		default:
			return c.ArgErr()
		case 1:
			fpath := filepath.ToSlash(filepath.Clean(cfg.Root + string(filepath.Separator) + tArgs[0]))

			if err := SetTemplate(stc.Template, "", fpath); err != nil {
				return c.Errf("default template parse error: %v", err)
			}

			stc.TemplateFiles[""] = &CachedFileInfo{
				Path: fpath,
			}
			return nil
		case 2:
			fpath := filepath.ToSlash(filepath.Clean(cfg.Root + string(filepath.Separator) + tArgs[1]))

			if err := SetTemplate(stc.Template, tArgs[0], fpath); err != nil {
				return c.Errf("template parse error: %v", err)
			}

			stc.TemplateFiles[tArgs[0]] = &CachedFileInfo{
				Path: fpath,
			}
			return nil
		}
	default:
		return c.Err("Expected valid stencil configuration")
	}
}
