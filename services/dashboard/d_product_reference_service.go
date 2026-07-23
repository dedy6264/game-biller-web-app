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

func AdminGetProductReferences(c echo.Context) error {
	var (
		svc = "AdminGetProductReferences"
		req models.RequestProductReferences
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductReferencesList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product references list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductReference(c echo.Context) error {
	var (
		svc = "AdminCreateProductReference"
		pr  models.ProductReference
	)
	if err := c.Bind(&pr); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pr.CreatedAt = now
	pr.UpdatedAt = now
	pr.CreatedBy = "admin"
	pr.UpdatedBy = "admin"

	_, err := repositories.CreateProductReference(connections.DBconn(), &pr)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product reference")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pr))
}

func AdminUpdateProductReference(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductReference"
		input models.ProductReference
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
	pr, err := repositories.GetProductReferenceByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product reference not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	pr.ProductReferenceCode = input.ProductReferenceCode
	pr.ProductReferenceName = input.ProductReferenceName
	pr.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pr.UpdatedBy = "admin"

	err = repositories.UpdateProductReference(db, pr)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product reference")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pr))
}

func AdminDeleteProductReference(c echo.Context) error {
	var (
		svc = "AdminDeleteProductReference"
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

	err := repositories.DeleteProductReference(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product reference")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
