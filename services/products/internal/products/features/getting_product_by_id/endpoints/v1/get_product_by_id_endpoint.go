package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/tracing"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/delivery"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/getting_product_by_id/dtos"
)

type getProductByIdEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewGetProductByIdEndpoint(productEndpointBase *delivery.ProductEndpointBase) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{productEndpointBase}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/:id", ep.getProductByID())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) getProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.GetProductByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.getProductByID")
		defer span.Finish()

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.Warn("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := getting_product_by_id.NewGetProductById(request.ProductId)

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			ep.Log.Warn("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		queryResult, err := mediatr.Send[*dtos.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			ep.Log.Warn("GetProductById", err)
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}