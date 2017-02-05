__zsh_history::filter::reverse()
{
    awk \
        '{
            line[NR] = $0
        }
        END {
            for (i = NR; i > 0; i--) {
                print line[i]
            }
        }' 2>/dev/null
}

__zsh_history::filter::unique()
{
    awk '!a[$0]++' 2>/dev/null
}

__zsh_history::filter::grep()
{
    if [[ -z $1 ]]; then
        cat -
    else
        grep --color="always" "$1"
    fi
}

__zsh_history::filter::remove_ansi()
{
    perl -pe 's/\e\[?.*?[\@-~]//g'
}

__zsh_history::filter::awk()
{
    # Imported by https://github.com/b4b4r07/ltsv.sh
    user_awk_script=$(cat <<-'EOS'
    { print $0 }
EOS
    )

    ltsv_awk_script=$(cat <<-'EOS'
    function key(name) {
        for (i = 1; i <= NF; i++) {
            match($i, ":");
            xs[substr($i, 0, RSTART)] = substr($i, RSTART+1);
        };
        return xs[name":"];
    }
EOS
    )

    awk_scripts="${ltsv_awk_script} ${argv[1]:-$user_awk_script}"

    if [[ ! -p /dev/stdin ]]; then
        echo "Error: failed to parse: you must supply something to work with via stdin" >&2
        exit 1
    fi

    awk -F'\t' -v "dir=$PWD" \
        "$awk_scripts"
}
