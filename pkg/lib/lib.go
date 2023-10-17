// The lib package is the set of standard libraries that are available to all
// Starlark scripts in a Starshell.
//
// In general, you should not put actual business logic in your own modules, but
// rather wrap separate libraries in a Starlark API. This allows you to keep
// your business logic separate from your Starlark API, and allows you to
// version your business logic separately from your Starlark API.
package lib
