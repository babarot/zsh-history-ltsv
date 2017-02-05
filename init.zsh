#!/bin/zsh

zmodload zsh/datetime 2>/dev/null || {
  print -r -- >&2 'zsh-history-enhanced: failed loading zsh/datatime'
  return 1
}

autoload -Uz add-zsh-hook

if [[ -z $ZSH_HISTORY_FILE ]]; then
    ZSH_HISTORY_FILE="$HOME/.zsh_history_enhanced"
fi

if [[ -z $ZSH_HISTORY_FILTER ]]; then
    ZSH_HISTORY_FILTER="fzy:fzf-tmux:fzf:peco"
fi

for f in "${0:A:h}"/src/*.zsh(N-.)
do
    source "$f"
done
unset f

if [[ -n $ZSH_HISTORY_KEYBIND_GET_BY_DIR ]]; then
    zle -N "__zsh_history::history::get_by_dir"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_BY_DIR" "__zsh_history::history::get_by_dir"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_GET_ALL ]]; then
    zle -N "__zsh_history::history::get_all"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_ALL" "__zsh_history::history::get_all"
fi

add-zsh-hook precmd "__zsh_history::history::add"
