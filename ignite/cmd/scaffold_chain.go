package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/pkg/placeholder"
	"github.com/ignite-hq/cli/ignite/services/scaffolder"
)

const (
	flagNoDefaultModule = "no-module"
)

// NewScaffoldChain creates new command to scaffold a Comos-SDK based blockchain.
func NewScaffoldChain() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain [name]",
		Short: "Fully-featured Cosmos SDK blockchain",
		Long: `Create a new application-specific Cosmos SDK blockchain.

For example, the following command will create a blockchain called "hello" in the "hello/" directory:

  ignite scaffold chain hello

A project name can be a simple name or a URL. The name will be used as the Go module path for the project. Examples of project names:

  ignite scaffold chain foo
  ignite scaffold chain foo/bar
  ignite scaffold chain example.org/foo
  ignite scaffold chain github.com/username/foo
		
A new directory with source code files will be created in the current directory. To use a different path use the "--path" flag.

Most of the logic of your blockchain is written in custom modules. Each module effectively encapsulates an independent piece of functionality. Following the Cosmos SDK convention, custom modules are stored inside the "x/" directory. By default, Ignite creates a module with a name that matches the name of the project. To create a blockchain without a default module use the "--no-module" flag. Additional modules can be added after a project is created with "ignite scaffold module" command.

Account addresses on Cosmos SDK-based blockchains have string prefixes. For example, the Cosmos Hub blockchain uses the default "cosmos" prefix, so that addresses look like this: "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf". To use a custom address prefix use the "--address-prefix" flag. For example:

  ignite scaffold chain foo --address-prefix bar

By default when compiling a blockchain's source code Ignite creates a cache to speed up the build process. To clear the cache when building a blockchain use the "--clear-cache" flag. It is very unlikely you will ever need to use this flag.

The blockchain is using the Cosmos SDK modular blockchain framework. Learn more about Cosmos SDK on https://docs.cosmos.network`,
		Args: cobra.ExactArgs(1),
		RunE: scaffoldChainHandler,
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagSetAccountPrefixes())
	c.Flags().StringP(flagPath, "p", ".", "Create a project in a specific path")
	c.Flags().Bool(flagNoDefaultModule, false, "Create a project without a default module")

	return c
}

func scaffoldChainHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	var (
		name               = args[0]
		addressPrefix      = getAddressPrefix(cmd)
		appPath            = flagGetPath(cmd)
		noDefaultModule, _ = cmd.Flags().GetBool(flagNoDefaultModule)
	)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	appdir, err := scaffolder.Init(cacheStorage, placeholder.New(), appPath, name, addressPrefix, noDefaultModule)
	if err != nil {
		return err
	}

	s.Stop()

	path, err := relativePath(appdir)
	if err != nil {
		return err
	}

	message := `
⭐️ Successfully created a new blockchain '%[1]v'.
👉 Get started with the following commands:

 %% cd %[1]v
 %% ignite chain serve

Documentation: https://docs.ignite.com
`
	fmt.Printf(message, path)

	return nil
}
