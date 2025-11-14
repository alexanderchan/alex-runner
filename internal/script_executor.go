package runner

// BuildScriptArgs constructs the command arguments for executing a script
// based on the package manager or build tool being used.
//
// For npm/yarn/pnpm scripts:
//   - npm requires '--' separator before additional args
//   - pnpm and yarn pass args through automatically
//
// For make:
//   - args are appended directly
func BuildScriptArgs(params BuildScriptArgsParams) []string {
	var cmdArgs []string

	if params.UseRun {
		cmdArgs = []string{"run", params.ScriptName}
		// npm requires -- to separate its flags from script args
		// pnpm/yarn pass args through automatically without needing --
		if len(params.AdditionalArgs) > 0 {
			if params.Command == "npm" {
				cmdArgs = append(cmdArgs, "--")
			}
			cmdArgs = append(cmdArgs, params.AdditionalArgs...)
		}
	} else {
		// For make, just append args directly
		cmdArgs = []string{params.ScriptName}
		cmdArgs = append(cmdArgs, params.AdditionalArgs...)
	}

	return cmdArgs
}

type BuildScriptArgsParams struct {
	Command        string   // "npm", "pnpm", "yarn", "make"
	ScriptName     string   // The script/target to run
	UseRun         bool     // Whether to use "run" subcommand (for npm/pnpm/yarn)
	AdditionalArgs []string // Additional arguments to pass to the script
}

// ParseArgs separates arguments at the '--' separator
// Returns the args before '--' and the args after '--'
func ParseArgs(args []string) (beforeSep []string, afterSep []string) {
	for i, arg := range args {
		if arg == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}
