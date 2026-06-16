// Package plugins defines the plugin system for the Telegram bot.
//
// Plugins are independent modules that can intercept and handle incoming
// updates before they reach the command dispatcher. Each plugin implements
// the Plugin interface and can decide whether to "consume" the update
// (by returning true), preventing further processing.
//
// Example usage:
//
//	type MyPlugin struct{}
//
//	func (p *MyPlugin) Name() string { return "myplugin" }
//	func (p *MyPlugin) Handle(ctx context.Context, bot *telego.Bot, update *telego.Update) bool {
//	    // handle update...
//	    return true // update consumed, stop processing
//	}
package plugins

import (
	"context"

	"github.com/mymmrac/telego"
)

// Plugin defines the contract for bot plugins.
//
// Each plugin is invoked in the order they are registered. If a plugin
// returns true, subsequent plugins and the command dispatcher are skipped.
type Plugin interface {
	// Name returns the plugin identifier (used for logging and debugging).
	Name() string

	// Handle processes the incoming update.
	//
	// Returns true if the update is fully handled (consumed) and further
	// processing should be stopped. Returns false to continue processing
	// the update through the plugin chain and command dispatcher.
	Handle(ctx context.Context, bot *telego.Bot, update *telego.Update) bool
}
