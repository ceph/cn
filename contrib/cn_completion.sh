# bash completion for cn                                   -*- shell-script -*-

__cn_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__cn_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__cn_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__cn_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__cn_handle_reply()
{
    __cn_debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __cn_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __cn_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__cn_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__cn_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__cn_handle_flag()
{
    __cn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __cn_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __cn_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __cn_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if __cn_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__cn_handle_noun()
{
    __cn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __cn_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __cn_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__cn_handle_command()
{
    __cn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_cn_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __cn_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__cn_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __cn_handle_reply
        return
    fi
    __cn_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __cn_handle_flag
    elif __cn_contains_word "${words[c]}" "${commands[@]}"; then
        __cn_handle_command
    elif [[ $c -eq 0 ]]; then
        __cn_handle_command
    elif __cn_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __cn_handle_command
        else
            __cn_handle_noun
        fi
    else
        __cn_handle_noun
    fi
    __cn_handle_word
}

_cn_cluster_ls()
{
    last_command="cn_cluster_ls"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_start()
{
    last_command="cn_cluster_start"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--work-dir=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--work-dir=")
    flags+=("--image=")
    two_word_flags+=("-i")
    local_nonpersistent_flags+=("--image=")
    flags+=("--data=")
    two_word_flags+=("-b")
    local_nonpersistent_flags+=("--data=")
    flags+=("--size=")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--size=")
    flags+=("--flavor=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--flavor=")
    flags+=("--help")
    local_nonpersistent_flags+=("--help")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_status()
{
    last_command="cn_cluster_status"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_stop()
{
    last_command="cn_cluster_stop"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_restart()
{
    last_command="cn_cluster_restart"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_logs()
{
    last_command="cn_cluster_logs"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_purge()
{
    last_command="cn_cluster_purge"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--yes-i-am-sure")
    local_nonpersistent_flags+=("--yes-i-am-sure")
    flags+=("--all")
    local_nonpersistent_flags+=("--all")
    flags+=("--help")
    local_nonpersistent_flags+=("--help")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster_enter()
{
    last_command="cn_cluster_enter"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_cluster()
{
    last_command="cn_cluster"

    command_aliases=()

    commands=()
    commands+=("ls")
    commands+=("start")
    commands+=("status")
    commands+=("stop")
    commands+=("restart")
    commands+=("logs")
    commands+=("purge")
    commands+=("enter")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_mb()
{
    last_command="cn_s3_mb"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_rb()
{
    last_command="cn_s3_rb"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_ls()
{
    last_command="cn_s3_ls"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_la()
{
    last_command="cn_s3_la"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_put()
{
    last_command="cn_s3_put"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_get()
{
    last_command="cn_s3_get"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--skip")
    flags+=("-s")
    local_nonpersistent_flags+=("--skip")
    flags+=("--force")
    flags+=("-f")
    local_nonpersistent_flags+=("--force")
    flags+=("--continue")
    flags+=("-c")
    local_nonpersistent_flags+=("--continue")
    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_del()
{
    last_command="cn_s3_del"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_du()
{
    last_command="cn_s3_du"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_info()
{
    last_command="cn_s3_info"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_cp()
{
    last_command="cn_s3_cp"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_mv()
{
    last_command="cn_s3_mv"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3_sync()
{
    last_command="cn_s3_sync"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--debug")
    flags+=("-d")
    local_nonpersistent_flags+=("--debug")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_s3()
{
    last_command="cn_s3"

    command_aliases=()

    commands=()
    commands+=("mb")
    commands+=("rb")
    commands+=("ls")
    commands+=("la")
    commands+=("put")
    commands+=("get")
    commands+=("del")
    commands+=("du")
    commands+=("info")
    commands+=("cp")
    commands+=("mv")
    commands+=("sync")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_image_update()
{
    last_command="cn_image_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_image_ls()
{
    last_command="cn_image_ls"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--all")
    flags+=("-a")
    local_nonpersistent_flags+=("--all")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_image_show-aliases()
{
    last_command="cn_image_show-aliases"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_image()
{
    last_command="cn_image"

    command_aliases=()

    commands=()
    commands+=("update")
    commands+=("ls")
    commands+=("show-aliases")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_version()
{
    last_command="cn_version"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_kube()
{
    last_command="cn_kube"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_update-check()
{
    last_command="cn_update-check"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_flavors_ls()
{
    last_command="cn_flavors_ls"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_flavors_show()
{
    last_command="cn_flavors_show"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_flavors()
{
    last_command="cn_flavors"

    command_aliases=()

    commands=()
    commands+=("ls")
    commands+=("show")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_completion()
{
    last_command="cn_completion"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_cn_root_command()
{
    last_command="cn"

    command_aliases=()

    commands=()
    commands+=("cluster")
    commands+=("s3")
    commands+=("image")
    commands+=("version")
    commands+=("kube")
    commands+=("update-check")
    commands+=("flavors")
    commands+=("completion")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()


    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_cn()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __cn_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("cn")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __cn_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_cn cn
else
    complete -o default -o nospace -F __start_cn cn
fi

# ex: ts=4 sw=4 et filetype=sh
