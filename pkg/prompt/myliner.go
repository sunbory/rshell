package prompt

import (
	"github.com/luckywinds/rshell/types"
	"github.com/peterh/liner"
	"log"
	"os"
	"sort"
	"strings"
)

func optCompleter(line string) (c []string) {
	l := oset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func keywordCompleter(line string) (c []string) {
	l := kset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	cl := cset.List()
	sort.Strings(cl)
	for _, value := range cl {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func cmdCompleter(line string) (c []string) {
	l := cset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func srcCompleter(line string) (c []string) {
	l := sset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func desCompleter(line string) (c []string) {
	l := dset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func authCompleter(line string) (c []string) {
	l := aset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func hostCompleter(line string) (c []string) {
	l := hset.List()
	sort.Strings(l)
	for _, value := range l {
		if strings.HasPrefix(value, strings.TrimLeft(line, " ")) {
			c = append(c, value)
		}
	}
	return
}

func choose(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func wordCompleter(line string, pos int) (head string, completions []string, tail string) {
	head = string([]rune(line)[:pos])
	tail = string([]rune(line)[pos:])

	as := strings.Fields(line)

	if len(as) == 0 {
		return "", keywordCompleter(""), tail
	} else {
		switch as[0] {
		case "sudo":
			if len(as) == 1 {
				return as[0] + " ", cmdCompleter(""), tail
			}
			if cfg.CmdSeparator != " " {
				cmds := strings.Split(strings.Join(as[1:], " "), cfg.CmdSeparator)
				cmds = choose(cmds, func(s string) bool {
					return s != ""
				})
				if len(cmds) == 1 {
					if strings.HasSuffix(line, cfg.CmdSeparator) {
						return as[0] + " " + cmds[0] + cfg.CmdSeparator, cmdCompleter(""), tail
					} else {
						return as[0] + " ", cmdCompleter(cmds[len(cmds)-1]), tail
					}
				} else {
					return as[0] + " " + strings.Join(cmds[0:len(cmds)-1], cfg.CmdSeparator) + cfg.CmdSeparator, cmdCompleter(cmds[len(cmds)-1]), tail
				}
			} else {
				return as[0] + " " + strings.Join(as[1:len(as)-1], " ") + " ", cmdCompleter(as[len(as)-1]), tail
			}
		case "download", "upload":
			if len(as) == 1 {
				return as[0] + " ", srcCompleter(""), tail
			}
			if len(as) == 2 {
				if strings.HasSuffix(line, " ") {
					return as[0] + " " + as[1] + " ", desCompleter(""), tail
				} else {
					return as[0] + " ", desCompleter(as[1]), tail
				}
			}
			if len(as) == 3 && !strings.HasSuffix(line, " ") {
				return as[0] + " " + as[1] + " ", desCompleter(""), tail
			}
		case "load":
			if len(as) == 1 {
				return as[0] + " ", optCompleter("-"), tail
			}
			if len(as) >= 1 {
				if strings.HasSuffix(line, " ") {
					return strings.Join(as, " ") + " ", optCompleter("-"), tail
				} else if strings.HasPrefix(as[len(as)-1], "-A") {
					return strings.Join(as[0:len(as)-1], " ") + " ", authCompleter(as[len(as)-1]), tail
				} else if strings.HasPrefix(as[len(as)-1], "-H"){
					return strings.Join(as[0:len(as)-1], " ") + " ", hostCompleter(as[len(as)-1]), tail
				}
			}
		default:
			cmds := strings.Split(strings.Join(as, " "), cfg.CmdSeparator)
			cmds = choose(cmds, func(s string) bool {
				return s != ""
			})
			if len(cmds) == 0 {
				return "", keywordCompleter(""), tail
			}
			if len(cmds) == 1 {
				if strings.HasSuffix(line, cfg.CmdSeparator) {
					return cmds[0] + cfg.CmdSeparator, cmdCompleter(""), tail
				} else {
					return "", keywordCompleter(cmds[0]), tail
				}
			} else {
				return strings.Join(cmds[0:len(cmds)-1], cfg.CmdSeparator) + cfg.CmdSeparator, cmdCompleter(cmds[len(cmds)-1]), tail
			}
		}
	}
	return head, nil, tail
}

func New(c types.Cfg, hostgroups types.Hostgroups) (*liner.State, error) {
	cfg = c

	for _, value := range hostgroups.Hgs {
		hset.Add(value.Groupname)
	}
	for _, value := range cfg.Mostusedcmds {
		cset.Add(value + c.CmdSeparator)
	}

	line := liner.NewLiner()
	line.SetWordCompleter(wordCompleter)
	line.SetTabCompletionStyle(liner.TabPrints)

	if f, err := os.Open(c.HistoryFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	return line, nil
}

func Prompt(l *liner.State, c types.Cfg) (string, error) {
	name, err := l.Prompt(c.PromptString)
	if err == nil {
		l.AppendHistory(name)
	} else if err == liner.ErrPromptAborted {
		return name, ErrPromptAborted
	} else {
		return "", ErrPromptError
	}

	if f, err := os.Create(c.HistoryFile); err != nil {
		log.Print("Error writing history file: ", err)
		return name, ErrPromptError
	} else {
		l.WriteHistory(f)
		f.Close()
	}

	return name, nil
}