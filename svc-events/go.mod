module github.com/ODIM-Project/ODIM/svc-events

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-00010101000000-000000000000
	github.com/ODIM-Project/ODIM/lib-messagebus v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-rest-client v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	gopkg.in/go-playground/validator.v9 v9.30.0
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
