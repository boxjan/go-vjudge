package controllers

type SubmitController struct {
	BaseController
}

func (c *SubmitController)SubmitPage()  {
	c.NeedLogin()
}

func (c *SubmitController)Submit()  {
	c.NeedLogin()
}
