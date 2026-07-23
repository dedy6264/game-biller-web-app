package dashboard

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func AdminGetProductPrefixes(c echo.Context) error {
	var (
		svc = "AdminGetProductPrefixes"
		req models.RequestProductPrefixes
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductPrefixesList(connections.DBconn(), req.Start, req.Length, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product prefixes list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductPrefix(c echo.Context) error {
	var (
		svc = "AdminCreateProductPrefix"
		pp  models.ProductPrefix
	)
	if err := c.Bind(&pp); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pp.CreatedAt = now
	pp.UpdatedAt = now
	pp.CreatedBy = "admin"
	pp.UpdatedBy = "admin"

	_, err := repositories.CreateProductPrefix(connections.DBconn(), &pp)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product prefix")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pp))
}

func AdminUpdateProductPrefix(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductPrefix"
		input models.ProductPrefix
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	pp, err := repositories.GetProductPrefixByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product prefix not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	pp.ProductReferenceID = input.ProductReferenceID
	pp.PrefixNumber = input.PrefixNumber
	pp.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pp.UpdatedBy = "admin"

	err = repositories.UpdateProductPrefix(db, pp)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product prefix")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pp))
}

func AdminDeleteProductPrefix(c echo.Context) error {
	var (
		svc = "AdminDeleteProductPrefix"
		req struct {
			ID int64 `json:"id"`
		}
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if req.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	err := repositories.DeleteProductPrefix(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product prefix")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
