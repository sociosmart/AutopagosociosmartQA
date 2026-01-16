package tasks

import "github.com/google/wire"

var TasksSet = wire.NewSet(
	ProvideSynchronizationTask,

	wire.Bind(new(SynchronizationTask), new(*synchronizationTask)),
)
