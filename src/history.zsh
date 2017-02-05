__zsh_history::history::add()
{
    local cmd
    if ! cmd="$(fc -ln -1)"; then
        return 1
    fi
    if [[ -o HIST_IGNORE_SPACE ]]; then
        [[ $cmd =~ ^" " ]] && return 0
    fi
    printf "date:$EPOCHSECONDS\tdir:$PWD\tcmd:$cmd\n" \
        >>|"$ZSH_HISTORY_FILE"
}

__zsh_history::history::get()
{
    local filter query="${1:?}"
    filter="$(__zsh_history::utils::get_filter "$ZSH_HISTORY_FILTER")"

    if [[ -z $filter ]]; then
        print -r -- >&2 'zsh-history-enhanced: ZSH_HISTORY_FILTER is an invalid'
        return 1
    fi

    BUFFER="$(
    cat "$ZSH_HISTORY_FILE" \
        | __zsh_history::filter::awk "$query" \
        | __zsh_history::filter::reverse \
        | __zsh_history::filter::unique \
        | __zsh_history::filter::grep "$LBUFFER" \
        | ${=filter} \
        | __zsh_history::filter::remove_ansi
    )"
    CURSOR=$#BUFFER
    zle reset-prompt
}

__zsh_history::history::get_by_dir()
{
    __zsh_history::history::get 'key("dir")==dir{print key("cmd")}'
}

__zsh_history::history::get_all()
{
    __zsh_history::history::get '{print key("cmd")}'
}

__zsh_history::history::edit()
{
    ${=EDITOR} "$ZSH_HISTORY_FILE" </dev/tty >/dev/tty
}

__zsh_history::history::show()
{
    local date dir cmd
    cat "$ZSH_HISTORY_FILE" \
        | __zsh_history::filter::awk '{print key("date"), key("dir"), key("cmd")}' \
        | while read date dir cmd; \
    do \
        printf "%s\t%s\t%s\n" \
        $(strftime "%FT%T%z" $date) "$dir" "$cmd"; \
    done
}
