package times

const (
	OneSecondInMilliseconds = 1000
	OneMinuteInSeconds      = 60
	OneHourInMinutes        = 60
	OneHourInMilliseconds   = OneSecondInMilliseconds * OneHourInMinutes * OneSecondInMilliseconds
	OneHourInSeconds        = 1 * OneHourInMinutes * OneMinuteInSeconds
)
