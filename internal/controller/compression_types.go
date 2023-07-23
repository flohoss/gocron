package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//	@Schemes
//	@Tags		compression_types
//	@Produce	json
//	@Success	200	{array}	database.SelectOption
//	@Router		/compression_types [get]
func (c *Controller) GetCompressionTypes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.service.GetSelectOptionsCompressionTypes())
}
