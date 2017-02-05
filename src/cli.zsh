zhistory() {
    while (( $# > 0 ))
    do
        case "$1" in
            (help)
                __zsh_history::history::help
                return $status
                ;;
            (edit)
                __zsh_history::history::edit
                return $status
                ;;
            (show)
                __zsh_history::history::show \
                    | less -F
                return $status
                ;;
            (*)
                print -r -- >&2 "zhistory: $1: no such arguments"
                return 1
                ;;
        esac
    done
}
