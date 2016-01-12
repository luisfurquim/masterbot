package masterbot

import (
   "errors"
   "regexp"
   "github.com/wangboo/cron"
   "github.com/luisfurquim/goose"
)

var ErrReadingSSHKeys       = errors.New("Error reading SSH keys (id_dsa)")
var ErrParsingSSHKeys       = errors.New("Error parsing SSH keys (id_dsa)")
var ErrDialingToBot         = errors.New("Error failed dialing to bot")
var ErrCreatingSession      = errors.New("Error failed creating session")
var ErrReadingConfig        = errors.New("Error reading config.json")
var ErrParsingConfig        = errors.New("Error parsing config.json")
var ErrLoadingCliCerts      = errors.New("Error Loading client certificates")
var ErrFailedStartingBot    = errors.New("Error failed to starting bot")
var ErrFailedPingingBot     = errors.New("Error failed pinging bot")
var ErrStatusStoppingBot    = errors.New("Error of status stopping bot")
var ErrStoppingBot          = errors.New("Error stopping bot")
var Goose                     goose.Alert
var Kairos                   *cron.Cron
//var SessionKeepAlive string
//var SessionQueueSize int = 5

var ReBotSearchInputParameters *regexp.Regexp = regexp.MustCompile("#([\\pL0-9_]+)")
var ReBotSearchOutputData      *regexp.Regexp = regexp.MustCompile("\\$([\\pL0-9_]+)")



