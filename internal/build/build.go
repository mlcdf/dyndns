package build

// Version is set at build time using -ldflags="-X 'github.com/mlcdf/dyndns/internal/build.Version=dyndns v1.0.0'
var Version = "(devel)"

// Time is set at build time using -ldflags="-X 'github.com/mlcdf/dyndns/internal/build.Time=$(date)'
var Time string

func Short() string {
	return "dyndns " + Version
}

func Long() string {
	info := Short()
	if Time != "" {
		info = info + " (release date: " + Time + ")"
	}
	return info
}
