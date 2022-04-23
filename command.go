package eelbot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command defines a slash command.
type Command struct {
	// MinArgs is the minimum number of arguments expected.
	MinArgs int

	// MaxArgs is the maximum number of arguments expected, or -1 for no limit.
	MaxArgs int

	// Summary is a single line summary of the command.
	Summary string

	// Usage is a multiple line summary of the command, with usage examples. It can use a single %s Sprintf directive
	// for the command name. If blank, Summary is used.
	Usage string

	// Eval should evaluate the command with the args provided. If a non nil error is returned, it is expected that no
	// replies were sent.
	Eval func(s Session, m *discordgo.MessageCreate, args []string) error
}

// RegisterCommand registers a new slash command. Any previously registered command with the same name will be
// overwritten. The cmd cannot be "help" as it is a reserved keyword.
func (bot *Bot) RegisterCommand(cmd string, c Command) {
	c.Summary = strings.TrimSpace(c.Summary)
	c.Usage = strings.TrimSpace(c.Usage)
	bot.cmds[strings.ToLower(cmd)] = &c
}

func (bot *Bot) createHelpCmd() {
	bot.cmds["help"] = &Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Displays the summary of available commands or details for a specific command.",
		Usage: `/%[1]s [CMD]

Displays the summary of available commands, or details for the specific CMD.

Examples:
  /%[1]s
  /%[1]s %[1]s
`,
		Eval: func(s Session, m *discordgo.MessageCreate, args []string) error {
			b := new(strings.Builder)
			if len(args) == 0 {
				cmds := make([]string, 0, len(bot.cmds))
				maxLen := 0
				for c := range bot.cmds {
					cmds = append(cmds, c)
					if len(c) > maxLen {
						maxLen = len(c)
					}
				}
				sort.Strings(cmds)
				cmdFmt := "/%-" + fmt.Sprint(maxLen) + "s: %s\n"
				b.WriteString("Available commands:\n")
				b.WriteString("```\n")
				for _, c := range cmds {
					b.WriteString(fmt.Sprintf(cmdFmt, c, bot.cmds[c].Summary))
				}
				b.WriteString("```")
			} else {
				cmd, ok := bot.cmds[args[0]]
				if !ok {
					return fmt.Errorf("unknown command: %s", args[0])
				}
				b.WriteString("```\n")
				if cmd.Usage != "" {
					b.WriteString(fmt.Sprintf(cmd.Usage, args[0]))
				} else {
					b.WriteString("/" + args[0])
					a := []string{}
					if cmd.MinArgs < 0 {
						cmd.MinArgs = 0
					}
					for i := 0; i < cmd.MinArgs; i++ {
						a = append(a, fmt.Sprintf("ARG%d", i+1))
					}
					if cmd.MaxArgs > cmd.MinArgs {
						for i := cmd.MinArgs; i < cmd.MaxArgs; i++ {
							a = append(a, fmt.Sprintf("[ARG%d]", i+1))
						}
					}
					if len(a) > 0 {
						if len(a) == 1 {
							a[0] = strings.ReplaceAll(a[0], "ARG1", "ARG")
						}
						b.WriteByte(' ')
						b.WriteString(strings.Join(a, " "))
					}
					if cmd.Summary != "" {
						b.WriteString(fmt.Sprintf("\n\n%s", cmd.Summary))
					}
				}
				b.WriteString("\n```")
			}
			s.ChannelMessageSend(m.ChannelID, b.String())
			return nil
		},
	}
}

func evalCmd(cmd string, c *Command, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) < c.MinArgs {
		return fmt.Errorf("%s requires at least %d arguments", cmd, c.MinArgs)
	}
	if c.MaxArgs >= 0 && len(args) > c.MaxArgs {
		return fmt.Errorf("%s requires at most %d arguments", cmd, c.MaxArgs)
	}
	return c.Eval(s, m, args)
}
