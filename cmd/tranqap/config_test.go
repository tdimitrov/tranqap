package main

import "testing"

var goodConfig = `targets:
- name: local
  host: "127.0.0.1"
  port: 22
  user: capture
  key: secret.key
  destination: pcaps
  file_pattern: trace
  file_rotation_count: 5
  use_sudo: true
  filter_port: 22`

func TestParseConfig(t *testing.T) {
	res, err := parseConfig([]byte(goodConfig))
	if err != nil {
		t.Errorf("Error parsing goodConfig: %s", err.Error())
	}

	if len(res.getTargetsList()) != 1 {
		t.Errorf("Expected the lenght of target list to be 1. Got %d", len(res.getTargetsList()))
	}

	if res.getTargetsList()[0] != "local" {
		t.Errorf("Expected the name of 1st target to be local. Got %s", res.getTargetsList()[0])
	}

	if len(res.Targets) != 1 {
		t.Errorf("Expected the lenght of targets to be 1. Got %d", len(res.Targets))
	}

	tgt := res.Targets[0]
	if *tgt.Name != "local" {
		t.Errorf("Bad name: %s", *tgt.Name)
	}
	if *tgt.Host != "127.0.0.1" {
		t.Errorf("Bad host: %s", *tgt.Host)
	}
	if *tgt.Port != 22 {
		t.Errorf("Bad port: %d", *tgt.Port)
	}
	if *tgt.User != "capture" {
		t.Errorf("Bad user: %s", *tgt.User)
	}
	if *tgt.Key != "secret.key" {
		t.Errorf("Bad key: %s", *tgt.Key)
	}
	if *tgt.Destination != "pcaps" {
		t.Errorf("Bad destination: %s", *tgt.Destination)
	}
	if *tgt.FilePattern != "trace" {
		t.Errorf("Bad file_pattern: %s", *tgt.FilePattern)
	}
	if *tgt.RotationCnt != 5 {
		t.Errorf("Bad file_rotation_count: %d", *tgt.RotationCnt)
	}
	if *tgt.UseSudo != true {
		t.Errorf("Bad use_sudo: %t", *tgt.UseSudo)
	}
	if *tgt.FilterPort != 22 {
		t.Errorf("Bad filter_port: %d", *tgt.FilterPort)
	}
}
