## Caddy Stencil

**Stencil is currently NOT stable (or even working!!!) and not recommended for use in production!**

Stencil is a templating middleware for Caddy server. It can be used with static files or combined or other middlewares (for example, proxy) to insert HTML, JSON, or HTML with JSON front matter into a pre-defined template.

Stencil is based on the Markdown plugin for Caddy but plans to focus on the integration of backend services and APIs and does not process Markdown.

### Syntax

```
stencil [basepath] {
	ext         extensions...
	template    [name] path
}
```

- **basepath** is the base path to match. Markdown will not activate if the request URL is not prefixed with this path. Default is site root.
- **extensions...** is a space-delimited list of file extensions to process with Stencil (defaults to .html, and .json).
- **template** defines a template with the given name to be at the given path. To specify the default template, omit name. Content can choose a template by using the name in its front matter or JSON.
- **templatedir** sets the default path with the given defaultpath to be used when listing templates.