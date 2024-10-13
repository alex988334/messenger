package wss_server

import (
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	hostlocal         = "localhost"
	localChatUrl      = "chat"
	host              = "orion734.ru"
	port              = "25550"
	avatarPath        = "avatar"
	versionClientPath = "version-client"
	downloadClient    = "download-client"
	aboutProject      = "about-project"

	indexFile        = "/root/messenger/pkg/index.html"
	versionFile      = "/root/messenger/pkg/messenger/version_client.txt"
	clientFile       = "/root/messenger/pkg/messenger/orion_client.txt"
	aboutProjectFile = "/root/messenger/pkg/messenger/about_project.txt"
)
