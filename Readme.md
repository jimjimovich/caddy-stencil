## Caddy Stencil

Stencil is a templating middleware for Caddy server. Stencil can process three types of input: raw HTML (or any other text-based format), raw HTML with JSON, YAML or TOML front matter, and valid JSON documents. Input can be from files or the output of another directive such as the [Proxy directive](https://caddyserver.com/docs/proxy).

Stencil processes input and runs it through pre-defined templates. Any JSON or front matter data is placed in the .Doc.data variable and the document body is placed in .Doc.body and are available to templates. The entire body of HTML input (without front matter) will is placed in the .Doc.body variable. The variable .Doc.title is either assigned by a root level JSON or front matter entry of "title" or automatically generated based on the file name.

### Syntax

```
stencil [basepath] {
	ext         extensions...
	template    [name] path
}
```

- **basepath** is the base path to match. Stencil will not activate if the request URL is not prefixed with this path. Default is site root.
- **extensions...** is a space-delimited list of file extensions to process with Stencil (defaults to .html, and .json).
- **template** defines a template with the given name to be at the given path. To specify the default template, omit name. Content can choose a template by using the name in its front matter or JSON.

### Processing HTML
Stencil can be used to inject raw HTML or text into templates. This may be useful for integrating legacy systems that don't have a JSON API.  The entire body of the document will be placed into the .Doc.body variable for use in your templates. 

**WARNING**: Injecting raw HTML into a template can be dangerous if the source of the HTML is from an untrusted source. Take precautions and make sure your input is trustworthy before injecting it into your template.  If you can't trust your input because you don't control it (for example, text input from a public API or website), be sure to use the [html, js, or urlquery functions](https://golang.org/pkg/text/template/#hdr-Functions) built into text/template to sanitize your input!


### Processing HTML with Front Matter
In addition to processing raw HTML (or text) as outlined above, Stencil will process documents with JSON, YAML or TOML front matter placed at the beginning of the document. The data in the front matter is placed in the .Doc.data variable to be used in your templates. The document body is placed in .Doc.body to be used in templates.

### Processing JSON Files and APIs
Stencil can be used to process valid JSON either from files or a live JSON API if used in conjunction with the [Proxy directive](https://caddyserver.com/docs/proxy). For Stencil to handle JSON files, the file name must contain the .json extension or, if using Proxy, must have either a .json extension or have a MIME type of "application/json".