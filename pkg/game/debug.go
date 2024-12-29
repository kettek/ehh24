package game

import "time"

const debug bool = false

//const debug bool = true

type profile struct {
	name     string
	start    time.Time
	end      time.Time
	duration time.Duration
}

var profiles []*profile

func getProfile(v string) *profile {
	for _, p := range profiles {
		if p.name == v {
			return p
		}
	}
	profiles = append(profiles, &profile{name: v})
	return profiles[len(profiles)-1]
}

func startProfile(v string) {
	if debug {
		getProfile(v).start = time.Now()
	}
}

func endProfile(v string) {
	if debug {
		p := getProfile(v)
		p.end = time.Now()
		p.duration = p.end.Sub(p.start)
	}
}
