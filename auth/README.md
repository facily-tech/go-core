# Auth

[![Go Reference](https://pkg.go.dev/badge/github.com/facily-tech/go-core/auth.svg)](https://pkg.go.dev/github.com/facily-tech/go-core/auth)

The purpose of the Auth package is to facilitate the authentication and authorization process using the openid connect (OIDC) standard.

## Usage

### HTTP middleware authentication against OIDC provider:

```go
zap, err := log.NewLoggerZap(log.ZapConfig{})
if err != nil {
  panic(err)
}

oidc, err := auth.New(zap, clientID, issuer)
if err != nil {
  panic(err)
}

r := chi.NewRouter()
// enable authentication middleware to all endpoints. 
// request context will be populated with scopes and roles.
r.Use(oidc.Auth)

// apply role middleware to /payment, only tokens with "charge" execute paymentHandler.
r.With(auth.HasRoleMiddleware("charge")).Get("/payment", paymentHandler)

zap.Fatal(context.Background(), "ending http", log.Error(http.ListenAndServe(":8181", r)))
```

### Validation outside transport layer

Maybe you want to use decorator pattern outside http and apply to your service.
Remember to pass request context to service calls, they'll be populated with scope and roles.

Suppose you have a service:

```go
type Payment interface{
  Create(context.Context, charge) error
}
```

You can decorate with your own type:

```go

var ErrRole = errors.New("invalid role")

type PaymentAuth struct{
  createRole string
  Payment
}

func NewPaymentAuth(createRole string, p Payment) *PaymentAuth {
  return &PaymentAuth{
    createRole: craeteRole,
    Payment: p,
  }
}

func (p *Payment) Create(ctx context.Context, c charge) error {
  if !auth.HasRole(ctx, p.createRole) {
    return errors.Wrapf(ErrRole, "missing '%s' role", p.createRole)
  }

  return p.Payment.Create(ctx, c)
}
```

Later in your service start:

```go
// Other starts before, including Payment at payment identifier
paymentAuth := NewPaymentAuth("charge", payment)

// Another service which requires Payment interface. PaymentAuth implements
// Payment with authorization.
cli := NewCli(paymentAuth)
```