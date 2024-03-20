package pkg

type Log struct {
	// log level, debug, info, warn, error
	Level string `toml:"level" yaml:"level" required:"true"`

	// Filename is the file to write logs to.  Backup log files will be retained in the same directory.
	// It uses <processname>-lumberjack.log in os.TempDir() if empty.
	Filename string `toml:"filename" yaml:"filename" required:"true"`

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.
	MaxSize int `toml:"maxsize" yaml:"maxsize" deafult:"100"`

	// MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	//  Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings,
	// leap seconds, etc. The default is not to remove old log files based on age.
	MaxAge int `toml:"maxage" yaml:"maxage" default:"7"`

	// MaxBackups is the maximum number of old log files to retain.  The default is to retain all old log files (though
	// MaxAge may still cause them to get deleted.)
	MaxBackups int `toml:"maxbackups" yaml:"maxbackups" default:"5"`

	// Compress determines if the rotated log files should be compressed using gzip.
	Compress bool `toml:"compress" yaml:"compress" default:"false"`
}
