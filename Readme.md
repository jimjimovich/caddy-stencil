## Caddy Stencil

**Stencil is currently NOT stable (or even working!!!) and not recommended for use in production!**

Stencil is a templating middleware for Caddy server. Stencil can process three types of input: raw HTML (or any other text-based format), raw HTML with JSON front matter, and valid JSON. This input can be from files or the result of another directive such as the Proxy directive.

Stencil processes input and runs it through pre-defined templates. Any JSON data is placed in the .Doc.data variable for templates. The entire body of HTML input with or without front matter will is placed in the .Doc.body variable. The variable .Doc.title is either assigned by a root level JSON entry of "title" or automatically generated based on the file name.

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

### Processing HTML
Stencil can be used to inject raw HTML or any text into templates. This may be useful for integrating legacy systems that don't have a JSON API into a modern website.  The entire body of the document will be placed into the .Doc.body variable for use in your templates. 

**WARNING**: Injecting raw HTML into a template can be dangerous if the source of the HTML is from an untrusted source. Take precautions and make sure your input is trustworthy before injecting it into your template.  If you can't trust your input because you don't control it (for example, text input from a public API or website), be sure to use the [html, js, or urlquery functions](https://golang.org/pkg/text/template/#hdr-Functions) built into text/template if you can not trust the source of your input!


### Processing HTML with JSON Front Matter
In addition to processing raw HTML (or text) as outlined above, Stencil will process JSON front matter placed at the beginning of text-based inputs.  The JSON front matter must be placed between a { character on the first line and end with a } character followed by a new line.  The data in the front matter is placed in the .Doc.data variable to be used in your templates.

### Processing JSON Files and APIs
The Stencil plugin can be used to process valid JSON either from files or a live JSON API if used in conjunction with the [Proxy directive](https://caddyserver.com/docs/proxy). For Stencil to handle JSON files must either contain the .json extension or, if using Proxy, must have either a .json extension or have a MIME type of "application/json".