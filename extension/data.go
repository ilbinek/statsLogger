package main

import (
	"sync"
	"time"
)

type Mission struct {
	mu            sync.Mutex   `json:"-"`
	MissionName   string       `json:"missionName"`
	Worldname     string       `json:"worldname"`
	MissionAuthor string       `json:"missionAuthor"`
	MissionType   string       `json:"missionType"`
	Victory       string       `json:"victory"`
	MissionStart  string       `json:"missionStart"`
	MissionEnd    string       `json:"missionEnd"`
	Date          string       `json:"date"`
	ScoreBlue     string       `json:"scoreBlue"`
	ScoreRed      string       `json:"scoreRed"`
	ScoreGreen    string       `json:"scoreGreen"`
	Players       []*Player    `json:"players"`
	Kills         []*Kill      `json:"kills"`
	FPS           []*FPSRecord `json:"fps"`
}

func (m *Mission) AddPlayer(p *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Players = append(m.Players, p)
}

func (m *Mission) AddKill(k *Kill) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Kills = append(m.Kills, k)
}

func (m *Mission) AddFPS(f *FPSRecord) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FPS = append(m.FPS, f)
}

type Player struct {
	mu    sync.Mutex `json:"-"`
	UID   string     `json:"uid"`
	Name  string     `json:"name"`
	Side  string     `json:"side"`
	Shots int        `json:"shots"`
	Hits  int        `json:"hits"`
	Kills int        `json:"kills"`
	Squad string     `json:"squad"`
	Role  string     `json:"role"`
	Class string     `json:"class"`
}

type PlayerUpdateOptions struct {
	Name  string `json:"name"`
	Side  string `json:"side"`
	Squad string `json:"squad"`
	Role  string `json:"role"`
	Class string `json:"class"`
}

func (p *Player) Update(options PlayerUpdateOptions) {
	p.mu.Lock()
	if options.Name != "" {
		p.Name = options.Name
	}
	if options.Side != "" {
		p.Side = options.Side
	}
	if options.Squad != "" {
		p.Squad = options.Squad
	}
	if options.Role != "" {
		p.Role = options.Role
	}
	if options.Class != "" {
		p.Class = options.Class
	}
	p.mu.Unlock()
}

func (p *Player) AddShot() {
	p.mu.Lock()
	p.Shots++
	p.mu.Unlock()
}

func (p *Player) AddHit() {
	p.mu.Lock()
	p.Hits++
	p.mu.Unlock()
}

func (p *Player) AddKill() {
	p.mu.Lock()
	p.Kills++
	p.mu.Unlock()
}

type Kill struct {
	Time     string `json:"time"`
	Victim   string `json:"victim"`
	Killer   string `json:"killer"`
	Weapon   string `json:"weapon"`
	Distance string `json:"distance"`
}

type FPSRecord struct {
	TimeUTC time.Time `json:"timeUTC"`
	FPS     float64   `json:"fps"`
}
