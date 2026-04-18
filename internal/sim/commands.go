package sim

import (
	"sort"

	"github.com/devlikebear/wirecraft/internal/netproto"
)

type QueuedCommand struct {
	Command          netproto.Command
	ReceivedSequence uint64
}

type commandIdentity struct {
	clientID  string
	commandID string
}

func OrderQueuedCommands(commands []QueuedCommand) []QueuedCommand {
	ordered := append([]QueuedCommand(nil), commands...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left := ordered[i]
		right := ordered[j]
		if left.Command.TickHint != right.Command.TickHint {
			return left.Command.TickHint < right.Command.TickHint
		}
		if left.ReceivedSequence != right.ReceivedSequence {
			return left.ReceivedSequence < right.ReceivedSequence
		}
		if left.Command.ClientID != right.Command.ClientID {
			return left.Command.ClientID < right.Command.ClientID
		}
		return left.Command.CommandID < right.Command.CommandID
	})
	return ordered
}

func identityForCommand(command netproto.Command) commandIdentity {
	return commandIdentity{
		clientID:  command.ClientID,
		commandID: command.CommandID,
	}
}
