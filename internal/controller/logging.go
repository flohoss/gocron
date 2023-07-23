package controller

import (
	"gitlab.unjx.de/flohoss/gobackup/database"
)

func (c *Controller) createLog(log *database.Log) {
	c.service.CreateOrUpdate(log)
}
