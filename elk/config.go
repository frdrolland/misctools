package elk

type Configuration struct {
	StateFile             string
	CrashOnInit           bool
	MaxCrash              int
	CleanupPidfileOnCrash bool
}
