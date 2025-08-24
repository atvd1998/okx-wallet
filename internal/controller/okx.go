package controller

import "fmt"

func (c *Controller) GetConnection() {
	c.okxService.GetConnection()
	fmt.Println("Testing connection")
}
