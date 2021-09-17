package config

const (
	GhNewCards                = 4212817
	GhUnderInvestigationCards = 4212819
	GhObservingCards          = 4212821
	GhCiSignalBoardProjectId  = 2093513
)

type GithubIssueCardConfig struct {
	CardName          string
	CardId            int
	OmitWithFlagShort bool
}

var GithubIssueCardConfigs = []GithubIssueCardConfig{
	{
		CardName:          "New",
		CardId:            4212817,
		OmitWithFlagShort: false,
	},
	{
		CardName:          "Investigation",
		CardId:            4212819,
		OmitWithFlagShort: false,
	},
	{
		CardName:          "Observing",
		CardId:            4212821,
		OmitWithFlagShort: false,
	},
	// {
	// 	CardName:          "CI Signal Board",
	// 	CardId:            2093513,
	// 	OmitWithFlagShort: false,
	// },
}
