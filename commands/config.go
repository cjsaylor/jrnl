package commands

type Configuration struct {
	JournalPath          string `env:"JOURNAL_PATH"`
	JournalEditor        string `env:"JRNL_EDITOR" envDefault:"vim"`
	JournalEditorOptions string `env:"JRNL_EDITOR_OPTIONS"`
}
