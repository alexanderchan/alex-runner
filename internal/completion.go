package runner

import (
	"fmt"
	"strings"
)

// GenerateBashCompletion returns a bash completion script for alex-runner
func GenerateBashCompletion() string {
	return `# bash completion for alex-runner                        -*- shell-script -*-

_alex_runner_completion() {
    local cur prev words cword
    _init_completion || return

    # Handle -- separator (after --, don't complete)
    for ((i=1; i < cword; i++)); do
        if [[ "${words[i]}" == "--" ]]; then
            return 0
        fi
    done

    # List of all flags
    local flags=(
        -l --last
        -s --search
        --list
        --list-names
        --use-package-json
        --use-makefile
        --no-cache
        --reset
        --global-reset
        --generate-completion
        -h --help
    )

    # If previous word is a flag that expects an argument
    case "$prev" in
        -s|--search)
            # Complete with script names
            local scripts
            scripts=$(alex-runner --list-names 2>/dev/null)
            COMPREPLY=($(compgen -W "$scripts" -- "$cur"))
            return 0
            ;;
        --generate-completion)
            # Complete with shell types
            COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
            return 0
            ;;
    esac

    # If current word starts with -, complete with flags
    if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "${flags[*]}" -- "$cur"))
        return 0
    fi

    # Otherwise, complete with script names (frecency-aware)
    local scripts
    scripts=$(alex-runner --list-names 2>/dev/null)
    COMPREPLY=($(compgen -W "$scripts" -- "$cur"))
    return 0
}

complete -F _alex_runner_completion alex-runner

# If you have aliases for alex-runner, register them too:
# Example: if you have "alias rr=alex-runner" in your .bashrc, add:
# complete -F _alex_runner_completion rr
`
}

// GenerateZshCompletion returns a zsh completion script for alex-runner
func GenerateZshCompletion() string {
	return `#compdef alex-runner

# zsh completion for alex-runner

_alex_runner() {
    local context state state_descr line
    typeset -A opt_args

    _arguments -C \
        '(- *)'{-h,--help}'[Show help message]' \
        '(-l --last)'{-l,--last}'[Run the most frecent script immediately]' \
        '(-s --search)'{-s,--search}'[Search for scripts matching term]:search term:_alex_runner_scripts' \
        '--list[List all scripts with frecency scores]' \
        '--list-names[List script names only (for completion)]' \
        '--use-package-json[Only show package.json scripts]' \
        '--use-makefile[Only show Makefile targets]' \
        '--no-cache[Re-detect package manager]' \
        '--reset[Clear usage history for current directory]' \
        '--global-reset[Clear all usage history]' \
        '--generate-completion[Generate completion script]:shell:(bash zsh fish)' \
        '*: :_alex_runner_scripts' \
        && return 0
}

# Helper function to get script names (frecency-aware)
_alex_runner_scripts() {
    local -a scripts
    scripts=(${(f)"$(alex-runner --list-names 2>/dev/null)"})
    _describe 'script' scripts
}

# Make sure completion system is loaded
if ! type compdef >/dev/null 2>&1; then
    autoload -Uz compinit
    compinit
fi

# Register the completion function
compdef _alex_runner alex-runner

# If you have aliases for alex-runner, register them too:
# Example: if you have "alias rr=alex-runner" in your .zshrc, add:
# compdef _alex_runner rr
`
}

// GenerateFishCompletion returns a fish completion script for alex-runner
func GenerateFishCompletion() string {
	return `# fish completion for alex-runner

# Don't complete if we're after --
function __alex_runner_after_double_dash
    set -l tokens (commandline -opc)
    contains -- -- $tokens
end

# Get script names (frecency-aware)
function __alex_runner_scripts
    if not __alex_runner_after_double_dash
        alex-runner --list-names 2>/dev/null
    end
end

# Complete flags
complete -c alex-runner -s h -l help -d 'Show help message'
complete -c alex-runner -s l -l last -d 'Run the most frecent script immediately'
complete -c alex-runner -s s -l search -d 'Search for scripts matching term' -r
complete -c alex-runner -l list -d 'List all scripts with frecency scores'
complete -c alex-runner -l list-names -d 'List script names only (for completion)'
complete -c alex-runner -l use-package-json -d 'Only show package.json scripts'
complete -c alex-runner -l use-makefile -d 'Only show Makefile targets'
complete -c alex-runner -l no-cache -d 'Re-detect package manager'
complete -c alex-runner -l reset -d 'Clear usage history for current directory'
complete -c alex-runner -l global-reset -d 'Clear all usage history'
complete -c alex-runner -l generate-completion -d 'Generate completion script' -r -f -a 'bash zsh fish'

# Complete script names dynamically (frecency-aware)
complete -c alex-runner -f -n 'not __alex_runner_after_double_dash' -a '(__alex_runner_scripts)'

# After --search flag, suggest script names
complete -c alex-runner -f -n '__fish_seen_subcommand_from --search -s' -a '(__alex_runner_scripts)'
`
}

// GetCompletionInstallInstructions returns installation instructions for a given shell
func GetCompletionInstallInstructions(shell string) string {
	switch strings.ToLower(shell) {
	case "bash":
		return `
Bash Completion Installation:

1. Generate the completion script:
   alex-runner --generate-completion bash > ~/.alex-runner-completion.bash

2. Add to your ~/.bashrc or ~/.bash_profile:
   source ~/.alex-runner-completion.bash

3. Reload your shell or run:
   source ~/.bashrc

Alternatively, for system-wide installation:
   sudo alex-runner --generate-completion bash > /etc/bash_completion.d/alex-runner
`

	case "zsh":
		return `
Zsh Completion Installation:

1. Generate the completion script:
   alex-runner --generate-completion zsh > ~/.alex-runner-completion.zsh

2. Add to your ~/.zshrc:
   source ~/.alex-runner-completion.zsh

3. Reload your shell or run:
   source ~/.zshrc

Alternatively, if you use a completion directory:
1. Create completion directory if it doesn't exist:
   mkdir -p ~/.zsh/completions

2. Generate completion file:
   alex-runner --generate-completion zsh > ~/.zsh/completions/_alex-runner

3. Add to your ~/.zshrc (before compinit):
   fpath=(~/.zsh/completions $fpath)
   autoload -Uz compinit && compinit

4. Reload your shell:
   exec zsh
`

	case "fish":
		return `
Fish Completion Installation:

1. Create fish completions directory if it doesn't exist:
   mkdir -p ~/.config/fish/completions

2. Generate the completion script:
   alex-runner --generate-completion fish > ~/.config/fish/completions/alex-runner.fish

3. Reload completions (or restart fish):
   fish_update_completions
`

	default:
		return fmt.Sprintf("Unknown shell: %s. Supported shells: bash, zsh, fish\n", shell)
	}
}

// ListCompletions returns a formatted list of available completions
func ListCompletions() string {
	return `Available completion shells:
  bash    Generate Bash completion script
  zsh     Generate Zsh completion script
  fish    Generate Fish completion script

Usage:
  alex-runner --generate-completion bash
  alex-runner --generate-completion zsh
  alex-runner --generate-completion fish

To see installation instructions:
  alex-runner --generate-completion bash | head -20
`
}
