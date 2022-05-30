package main

import "time"

type duration time.Duration

func (d *duration) Type() string {
	return "time.Duration"
}

func (d *duration) Set(s string) error {
	ds, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = (duration)(ds)
	return nil
}

func (d *duration) String() string {
	return (time.Duration)(*d).String()
}
