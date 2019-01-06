package metrics

var Methods = methodTypes{
	Service: serviceTypes{
		Get:    serviceMethodType("Get"),
		Update: serviceMethodType("Update"),
		Add:    serviceMethodType("Add"),
		Delete: serviceMethodType("Delete"),
	},
	DB: dbTypes{
		Delete: dbMethodType("Delete"),
		Update: dbMethodType("Update"),
		Read:   dbMethodType("Read"),
		Insert: dbMethodType("Insert"),
	},
	Cache: cacheTypes{
		Hit:        cacheMethodType("Hit"),
		Miss:       cacheMethodType("Miss"),
		Invalidate: cacheMethodType("Invalidate"),
		Store:      cacheMethodType("Store"),
	},
	Transport: transportTypes{
		Put:    transportMethodType("Put"),
		Delete: transportMethodType("Delete"),
		Post:   transportMethodType("Post"),
		Get:    transportMethodType("Get"),
		Head:   transportMethodType("Head"),
	},
}

type methodTypes struct {
	Service   serviceTypes
	DB        dbTypes
	Cache     cacheTypes
	Transport transportTypes
}

type serviceMethodType string
type serviceTypes struct {
	Get    serviceMethodType
	Update serviceMethodType
	Add    serviceMethodType
	Delete serviceMethodType
}

func (c serviceMethodType) String() string {
	return "pkg.metric.service." + string(c)
}

type dbMethodType string
type dbTypes struct {
	Read   dbMethodType
	Update dbMethodType
	Insert dbMethodType
	Delete dbMethodType
}

func (c dbMethodType) String() string {
	return "pkg.metric.db." + string(c)
}

type cacheMethodType string
type cacheTypes struct {
	Hit        cacheMethodType
	Miss       cacheMethodType
	Store      cacheMethodType
	Invalidate cacheMethodType
}

func (c cacheMethodType) String() string {
	return "pkg.metric.cache." + string(c)
}

type transportMethodType string
type transportTypes struct {
	Get    transportMethodType
	Put    transportMethodType
	Post   transportMethodType
	Delete transportMethodType
	Head   transportMethodType
}

func (c transportMethodType) String() string {
	return "pkg.metric.transport." + string(c)
}

type EndPointType string

var EndPoints = endPointTypes{
	Gophers: EndPointType("Gophers"),
	Gopher:  EndPointType("Gopher"),
	Things:  EndPointType("Things"),
	Thing:   EndPointType("Thing"),
}

type endPointTypes struct {
	Gophers EndPointType
	Gopher  EndPointType
	Things  EndPointType
	Thing   EndPointType
}

func (c EndPointType) String() string {
	return "pkg.metrics.endpoints." + string(c)
}
