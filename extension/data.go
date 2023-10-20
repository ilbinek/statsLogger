package main

import (
	"time"
)

type Mission struct {
	ID            uint        `json:"id" gorm:"primaryKey"` // Primary key
	MissionName   string      `json:"missionName"`
	Worldname     string      `json:"worldname"`
	MissionAuthor string      `json:"missionAuthor"`
	MissionType   string      `json:"missionType"`
	Victory       string      `json:"victory"`
	MissionStart  string      `json:"missionStart"`
	MissionEnd    string      `json:"missionEnd"`
	Date          string      `json:"date"`
	ScoreBlue     string      `json:"scoreBlue"`
	ScoreRed      string      `json:"scoreRed"`
	ScoreGreen    string      `json:"scoreGreen"`
	Players       []Player    `json:"players"`
	Kills         []Kill      `json:"kills"`
	FPSRecords    []FPSRecord `json:"fps"`
}

type Player struct {
	ID        uint     `json:"id" gorm:"uniqueIndex"`
	PlayerUID string   `json:"playerUID" gorm:"index"` // Player UID
	Name      string   `json:"name"`
	Side      string   `json:"side"`
	Shots     int      `json:"shots"`
	Hits      int      `json:"hits"`
	Squad     string   `json:"squad"`
	Role      string   `json:"role"`
	Class     string   `json:"class"`
	Mission   *Mission `json:"-" gorm:"foreignKey:MissionID"`
	MissionID uint     `json:"-" gorm:"index"`
}

type Kill struct {
	ID        uint     `json:"id" gorm:"primaryKey"` // Primary key
	Time      string   `json:"time"`
	Victim    string   `json:"victim"`
	Killer    string   `json:"killer"`
	Weapon    string   `json:"weapon"`
	Distance  string   `json:"distance"`
	Mission   *Mission `json:"-"`
	MissionID uint     `json:"-" gorm:"index"`
}

type FPSRecord struct {
	TimeUTC   time.Time `json:"timeUTC" gorm:"primaryKey"` // Primary key
	FPS       float64   `json:"fps"`
	Mission   *Mission  `json:"-"`
	MissionID uint      `json:"-" gorm:"index"`
}
