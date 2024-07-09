package auth_http_handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	auth_http_api "github.com/krissukoco/go-food-order-microservices/api/openapi/auth"
	login_usecase "github.com/krissukoco/go-food-order-microservices/internal/services/auth/usecase/auth/login"
	"github.com/krissukoco/go-food-order-microservices/internal/transport"
)

type handler struct {
	loginUc login_usecase.Usecase
}

var _ auth_http_api.ServerInterface = (*handler)(nil)

func New(
	loginUc login_usecase.Usecase,
) auth_http_api.ServerInterface {
	return &handler{loginUc}
}

// Login to gain tokens for Backoffice users
// (POST /login/backoffice)
func (h *handler) LoginBackoffice(c *fiber.Ctx) error {
	var req auth_http_api.LoginBackofficeJSONRequestBody
	httpReq := transport.HTTPRequest{
		Body: &req,
	}
	if err := transport.ParseAndValidateHTTPRequest(c, &httpReq); err != nil {
		return err
	}

	tokens, err := h.loginUc.LoginBackoffice(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, login_usecase.ErrInvalidCredentials) {
			return transport.HTTPError(c, http.StatusBadRequest, 10002, err.Error())
		}
		return err
	}

	return c.JSON(&auth_http_api.AuthTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Login as customer by providing phone number
// (POST /login/customer)
func (h *handler) LoginCustomer(c *fiber.Ctx) error {
	var req auth_http_api.LoginCustomerJSONRequestBody
	httpReq := transport.HTTPRequest{
		Body: &req,
	}
	if err := transport.ParseAndValidateHTTPRequest(c, &httpReq); err != nil {
		return err
	}

	otp, err := h.loginUc.SendCustomerOTP(c.Context(), req.PhoneNumber)
	if err != nil {
		if errors.Is(err, login_usecase.ErrInvalidCredentials) {
			return transport.HTTPError(c, http.StatusBadRequest, 10002, err.Error())
		}
		return err
	}

	return c.JSON(&auth_http_api.RequestOtpResponse{
		Id:        otp.Id.String(),
		ExpiredAt: otp.ExpiredAt.Format(time.RFC3339),
	})
}

// Verify OTP that was sent to phone
// (POST /login/customer/verify-otp/{id})
func (h *handler) CustomerVerifyOtp(c *fiber.Ctx, id uuid.UUID) error {
	var req auth_http_api.CustomerVerifyOtpJSONBody
	httpReq := transport.HTTPRequest{
		Body: &req,
	}
	if err := transport.ParseAndValidateHTTPRequest(c, &httpReq); err != nil {
		return err
	}

	tokens, err := h.loginUc.VerifyCustomerOTP(c.Context(), id, req.Otp)
	if err != nil {
		if errors.Is(err, login_usecase.ErrInvalidCredentials) {
			return transport.HTTPError(c, http.StatusBadRequest, 10002, err.Error())
		}
		return err
	}

	return c.JSON(&auth_http_api.AuthTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Login to gain tokens for POS users
// (POST /login/pos)
func (h *handler) LoginPos(c *fiber.Ctx) error {
	var req auth_http_api.LoginPosJSONRequestBody
	httpReq := transport.HTTPRequest{
		Body: &req,
	}
	if err := transport.ParseAndValidateHTTPRequest(c, &httpReq); err != nil {
		return err
	}

	tokens, err := h.loginUc.LoginPos(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, login_usecase.ErrInvalidCredentials) {
			return transport.HTTPError(c, http.StatusBadRequest, 10002, err.Error())
		}
		return err
	}

	return c.JSON(&auth_http_api.AuthTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
